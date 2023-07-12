package handler

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"shop_srvs/order_srv/global"
	"shop_srvs/order_srv/model"
	"shop_srvs/order_srv/proto"
	"time"
)

// 商品服务
type OrderServer struct {
	proto.UnimplementedOrderServer
}

// 购物车信息
func (self *OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	//获取用户的购物车列表
	var shopCarts []model.ShoppingCart
	var response proto.CartItemListResponse

	if result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts); result.Error != nil {
		return nil, result.Error
	} else {
		response.Total = int32(result.RowsAffected)
	}
	for _, shopCart := range shopCarts {
		response.Data = append(response.Data, &proto.ShopCartInfoResponse{
			Id:      shopCart.ID,
			UserId:  shopCart.User,
			GoodsId: shopCart.Goods,
			Nums:    shopCart.Num,
			Checked: shopCart.Checked,
		})
	}
	return &response, nil
}

// 添加购物车
func (self *OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	//将商品添加到购物车
	//1. 原本没有 新建
	//2. 原本有 更新，合并
	var shopCart model.ShoppingCart
	result := global.DB.Where(&model.ShoppingCart{User: req.UserId, Goods: req.GoodsId}).First(&shopCart)
	if result.RowsAffected == 1 {
		//原本有改变数量就可以
		shopCart.Num += req.Nums
	} else {
		//插入操作
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Num = req.Nums
		shopCart.Checked = true
	}
	global.DB.Save(&shopCart)
	return &proto.ShopCartInfoResponse{Id: shopCart.ID}, nil
}

// 更新购物车
func (self *OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*empty.Empty, error) {
	//更新购物车,更新数量和选中状态，不能小于0
	var shopCart model.ShoppingCart
	if result := global.DB.First(&shopCart, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	shopCart.Checked = req.Checked
	if req.Nums > 0 {
		shopCart.Num = req.Nums
	}
	global.DB.Save(&shopCart)
	return &empty.Empty{}, nil
}

// 删除购物车
func (self *OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*empty.Empty, error) {
	if result := global.DB.Delete(&model.ShoppingCart{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	return &empty.Empty{}, nil
}

// 订单接口
func GenerateOrderSn(userId int32) string {
	//订单号生成规则
	//年月日时分秒+用户id+2位随机数
	now := time.Now()
	rand.Seed(now.UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(), userId, rand.Intn(90)+10)
	return orderSn
}

// 新建订单
func (self *OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	/*
				1. 从购物车中获取到选中的商品
				2. 金额自己查询，不能由前端给 ，跨微服务调用了，
			    3. 库存扣减，跨微服务
			    4. 订单的基本信息表，订单的商品信息表
			    5. 从购物车中删除
		    本地事务
			todo 事务问题！分布式事务
	*/

	var shopCarts []model.ShoppingCart
	var goodsIds []int32
	goodsNumsMap := map[int32]int32{} //商品对应的数量
	if result := global.DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Find(&shopCarts); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "购物车中没有选中的商品")
	}
	//获取id
	for _, shopCart := range shopCarts {
		goodsIds = append(goodsIds, shopCart.Goods)
		goodsNumsMap[shopCart.Goods] = shopCart.Num //获取商品的数量
	}
	//跨服务调用,商品微服务
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGetGoodsInfo{Id: goodsIds})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "批量获取商品失败")
	}
	//获取商品的价格，计算总价
	var totalPrice float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo

	for _, good := range goods.Data {
		totalPrice += good.ShopPrice * float32(goodsNumsMap[good.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Num:        goodsNumsMap[good.Id],
		})
		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     goodsNumsMap[good.Id],
		})
	}
	//跨库存微服务，扣减库存
	_, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{GoodsInvInfos: goodsInvInfo})
	if err != nil {
		return nil, status.Errorf(codes.ResourceExhausted, "扣减库存失败")
	}

	//本地事务开始
	tx := global.DB.Begin()

	//订单表生成
	//订单编号生成，20230712xxxxxx
	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.UserId),
		OrderMount:   totalPrice,
		Address:      req.Address,
		SignerName:   req.Name,
		SignerMobile: req.Mobile,
		Post:         req.Post,
		User:         req.UserId,
	}
	//保存
	if result := tx.Save(&order); result.RowsAffected == 0 {
		tx.Rollback() //回滚
		return nil, status.Errorf(codes.Internal, "订单生成失败")
	}
	for _, orderGood := range orderGoods {
		orderGood.Order = order.ID
	}
	//批量插入
	if result := tx.CreateInBatches(&orderGoods, 1000); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "订单商品生成失败")
	}
	//删除购物车
	if result := tx.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "删除购物车失败")
	}
	tx.Commit() //提交事务
	return &proto.OrderInfoResponse{
		Id:      order.ID,
		OrderSn: order.OrderSn,
		Total:   order.OrderMount,
	}, nil
}

// 订单列表
func (self *OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	//获取订单列表,
	//后台和商家查询可以区分,有id就是商家
	var orders []model.OrderInfo
	resp := &proto.OrderListResponse{}
	//关于sql语句拼接问题，如果有问题就改为指针
	var total int64
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	resp.Total = int32(total)

	//分页
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orders)
	for _, order := range orders {
		resp.Data = append(resp.Data, &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SignerMobile,
		})
	}
	return resp, nil
}

// 订单详情，需要知道商品的信息
func (self *OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var order model.OrderInfo
	var resp proto.OrderInfoDetailResponse
	//是否是当前用户的订单，由web层管理
	//如果是管理员只传 order的id，
	//如果是商家传order的id和商家的id
	//注意0值拼接问题
	if result := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&order); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	resp.OrderInfo = &proto.OrderInfoResponse{
		Id:      order.ID,
		UserId:  order.User,
		OrderSn: order.OrderSn,
		PayType: order.PayType,
		Status:  order.Status,
		Post:    order.Post,
		Total:   order.OrderMount,
		Address: order.Address,
		Name:    order.SignerName,
		Mobile:  order.SignerMobile,
	}
	var orderGoods []model.OrderGoods
	if result := global.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&orderGoods); result.Error != nil {
		return nil, result.Error
	}
	for _, orderGood := range orderGoods {
		resp.Goods = append(resp.Goods, &proto.OrderItemResponse{
			Id:         orderGood.ID,
			OrderId:    orderGood.Order,
			GoodsId:    orderGood.Goods,
			GoodsName:  orderGood.GoodsName,
			GoodsImg:   orderGood.GoodsImage,
			GoodsPrice: orderGood.GoodsPrice,
			Nums:       orderGood.Num,
		})

	}

	return &resp, nil
}

// 更新订单状态
func (self *OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*empty.Empty, error) {
	//支付后会调用此接口
	//成功会由orderSn
	if result := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	return &empty.Empty{}, nil
}

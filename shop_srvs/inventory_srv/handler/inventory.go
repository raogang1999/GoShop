package handler

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop_srvs/inventory_srv/global"
	"shop_srvs/inventory_srv/model"
	"shop_srvs/inventory_srv/proto"
)

// 商品服务
type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

// SetInv 设置库存
func (InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*empty.Empty, error) {
	//没有就新增，有就更新
	//先查询
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GetGoodsId()}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num
	global.DB.Save(&inv)
	return &empty.Empty{}, nil
}

// InvDetail 查询库存
func (InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GetGoodsId()}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有找到库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

// Sell 扣减库存
func (InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*empty.Empty, error) {
	//本地事务，要么全部成功，要么全部失败
	// 并发会出问题
	tx := global.DB.Begin() //开启事务

	for _, goodInfo := range req.GoodsInvInfos {
		var inv model.Inventory

		mutex := global.RedisLocker.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			tx.Rollback() //回滚事务
			return nil, status.Errorf(codes.Internal, "Redis锁定失败")
		}

		if result := tx.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚事务
			return nil, status.Errorf(codes.InvalidArgument, "没有找到库存信息")
		}
		//判断库存是否足够
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚事务
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)
		if ok, err := mutex.Unlock(); !ok || err != nil {
			tx.Rollback() //回滚事务
			return nil, status.Errorf(codes.Internal, "Redis解锁失败")
		}
	}

	tx.Commit() //提交事务
	return &empty.Empty{}, nil
}

// Reback 回滚库存
func (InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*empty.Empty, error) {
	//库存归还，1.超时，2.取消订单，3.下单失败

	//本地事务，要么全部成功，要么全部失败
	// 并发会出问题
	tx := global.DB.Begin() //开启事务
	for _, goodInfo := range req.GoodsInvInfos {
		var inv model.Inventory
		if result := tx.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚事务
			return nil, status.Errorf(codes.InvalidArgument, "没有找到库存信息")
		}
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() //提交事务
	return &empty.Empty{}, nil
}

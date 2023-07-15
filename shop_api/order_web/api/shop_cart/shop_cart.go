package shop_cart

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"shop_api/order_web/api"
	"shop_api/order_web/forms"
	"shop_api/order_web/global"
	"shop_api/order_web/proto"
	"strconv"
)

func List(ctx *gin.Context) {
	//获取购物车的商品
	// 通过user id
	userId, _ := ctx.Get("userId")
	// 调用后端的购物车服务
	rsp, err := global.OrderSrvClient.CartItemList(ctx, &proto.UserInfo{
		Id: int32(userId.(uint)),
	})
	if err != nil {
		zap.S().Errorw("[List] 查询购物车失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ids := make([]int32, 0)
	for _, item := range rsp.Data {
		ids = append(ids, item.GoodsId)
	}
	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	// 调用商品服务查询商品详情
	goodsRsp, err := global.GoodsSrvClient.BatchGetGoods(ctx, &proto.BatchGetGoodsInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("[List] 批量查询商品失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//响应数据格式
	/*
		{
			"total": 100,
			"goods_id": 1,
			"goods_name": "测试商品",
			"goods_image": "http://www.img.com",
			"goods_price": 100.00,
			"nums": 4,
			"check":true,
		}

	*/

	respMap := gin.H{
		"total": rsp.Total,
		"data":  make([]interface{}, 0),
	}
	for _, item := range rsp.Data {
		for _, goods := range goodsRsp.Data {
			if item.GoodsId == goods.Id {
				tmpMap := gin.H{
					"id":          item.Id,
					"goods_id":    item.GoodsId,
					"goods_name":  goods.Name,
					"goods_image": goods.GoodsFrontImage,
					"goods_price": goods.ShopPrice,
					"nums":        item.Nums,
					"check":       item.Checked,
				}
				respMap["data"] = append(respMap["data"].([]interface{}), tmpMap)
			}
		}
	}
	ctx.JSON(http.StatusOK, respMap)
}
func Delete(ctx *gin.Context) {
	//删除购物车的商品
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "url格式错误了！",
		})
		return
	}
	//获取用户id
	userId, _ := ctx.Get("userId")
	// 调用后端的购物车服务
	_, err = global.OrderSrvClient.DeleteCartItem(ctx, &proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(idInt),
	})
	if err != nil {
		zap.S().Errorw("[Delete] 删除购物车商品失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)

}

func NewCart(ctx *gin.Context) {
	//添加商品
	itemForm := forms.ShopCartItemForm{}
	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	//查询商品是否存在，
	_, err := global.GoodsSrvClient.GetGoodsDetail(ctx, &proto.GoodInfoRequest{
		Id: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[NewCart] 查询商品失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	//查询库存服务
	invRsp, err := global.InventorySrvClient.InvDetail(ctx, &proto.GoodsInvInfo{
		GoodsId: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[NewCart] 查询库存失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	if invRsp.Num < itemForm.Nums {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "库存不足",
		})
		return
	}

	//获取用户id
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CreateCartItem(ctx, &proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: itemForm.GoodsId,
		Nums:    itemForm.Nums,
	})
	if err != nil {
		zap.S().Errorw("[NewCart] 添加商品失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})

}

func UpdateCart(ctx *gin.Context) {
	//更新购物车商品
	itemForm := forms.ShopCartUpdateItemForm{}
	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	//获取商品id
	goodId := ctx.Param("id")
	goodIdInt, err := strconv.Atoi(goodId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "url格式错误了！",
		})
		return
	}

	//获取用户id
	userId, _ := ctx.Get("userId")

	request := &proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(goodIdInt),
		Nums:    itemForm.Nums,
	}
	if itemForm.Checked != nil {
		request.Checked = *itemForm.Checked
	}
	_, err = global.OrderSrvClient.UpdateCartItem(ctx, request)

	if err != nil {
		zap.S().Errorw("[UpdateCart] 更新购物车商品失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

package order

import (
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"net/http"
	"shop_api/order_web/api"
	"shop_api/order_web/forms"
	"shop_api/order_web/global"
	"shop_api/order_web/models"
	"shop_api/order_web/proto"
	"strconv"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	//判断是否为管理员
	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	request := proto.OrderFilterRequest{}
	if model.AuthorityId == 1 {
		//普通用户
		request.UserId = int32(userId.(uint))
	}
	//否则为管理
	//构造分页
	pages := ctx.DefaultQuery("p", "0")
	PagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(PagesInt)
	pageNum := ctx.DefaultQuery("pnum", "10")
	PageNumInt, _ := strconv.Atoi(pageNum)
	request.PagePerNums = int32(PageNumInt)
	//请求
	rsp, err := global.OrderSrvClient.OrderList(ctx, &request)
	if err != nil {
		zap.S().Errorw("获取订单列表失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//返回数据
	/*
		"total": 100,
		"data": [
			{
			}
		]
	*/
	respMap := gin.H{
		"total": rsp.Total,
	}
	orderList := make([]interface{}, 0)
	for _, item := range rsp.Data {
		tmpMap := map[string]interface{}{}
		tmpMap["id"] = item.Id
		tmpMap["status"] = item.Status
		tmpMap["pay_type"] = item.PayType
		tmpMap["user"] = item.UserId
		tmpMap["post"] = item.Post
		tmpMap["total"] = item.Total
		tmpMap["address"] = item.Address
		tmpMap["name"] = item.Name
		tmpMap["mobile"] = item.Mobile
		tmpMap["order_sn"] = item.OrderSn
		tmpMap["add_time"] = item.AddTime
		orderList = append(orderList, tmpMap)
	}
	respMap["data"] = orderList
	ctx.JSON(http.StatusOK, respMap)
}
func NewOrder(ctx *gin.Context) {
	orderForm := forms.CreateOrderForm{}
	if err := ctx.ShouldBind(&orderForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	//获取用户id
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CreateOrder(ctx, &proto.OrderRequest{
		UserId:  int32(userId.(uint)),
		Address: orderForm.Address,
		Name:    orderForm.Name,
		Mobile:  orderForm.Mobile,
		Post:    orderForm.Post,
	})
	if err != nil {
		zap.S().Errorw("创建订单失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	client, err := alipay.New(global.ServerConfig.AliPayInfo.AppId, global.ServerConfig.AliPayInfo.PrivateKey, false)

	if err != nil {
		zap.S().Errorw("创建支付宝客户端失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = client.LoadAliPayPublicKey(global.ServerConfig.AliPayInfo.AliPayPubKey)
	if err != nil {
		zap.S().Errorw("加载支付宝公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServerConfig.AliPayInfo.NotifyUrl //回调
	p.ReturnURL = global.ServerConfig.AliPayInfo.ReturnUrl
	p.Subject = "GoShop订单" + rsp.OrderSn
	p.OutTradeNo = rsp.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(rsp.Total), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("创建支付宝订单失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         rsp.Id,
		"alipay_url": url.String(),
	})
}

func GetOrderById(ctx *gin.Context) {
	//获取订单详情
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "url格式错误",
		})
		return
	}
	//判断是否为管理员
	request := proto.OrderRequest{}
	request.Id = int32(idInt) //订单id

	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		//普通用户
		userId, _ := ctx.Get("userId")
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.OrderSrvClient.OrderDetail(ctx, &request)
	if err != nil {
		zap.S().Errorw("获取订单详情失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//返回数据
	respMap := gin.H{}
	respMap["id"] = rsp.OrderInfo.Id
	respMap["status"] = rsp.OrderInfo.Status
	respMap["user"] = rsp.OrderInfo.UserId
	respMap["post"] = rsp.OrderInfo.Post
	respMap["total"] = rsp.OrderInfo.Total
	respMap["address"] = rsp.OrderInfo.Address
	respMap["name"] = rsp.OrderInfo.Name
	respMap["mobile"] = rsp.OrderInfo.Mobile
	respMap["pay_type"] = rsp.OrderInfo.PayType
	respMap["order_sn"] = rsp.OrderInfo.OrderSn
	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Goods {
		tmpMap := map[string]interface{}{}
		tmpMap["id"] = item.Id
		tmpMap["goods_id"] = item.GoodsId
		tmpMap["goods_price"] = item.GoodsPrice
		tmpMap["goods_name"] = item.GoodsName
		tmpMap["goods_image"] = item.GoodsImg
		tmpMap["num"] = item.Nums
		goodsList = append(goodsList, tmpMap)
	}
	respMap["goods"] = goodsList

	//生成支付宝的支付url
	client, err := alipay.New(global.ServerConfig.AliPayInfo.AppId, global.ServerConfig.AliPayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = client.LoadAliPayPublicKey((global.ServerConfig.AliPayInfo.AliPayPubKey))
	if err != nil {
		zap.S().Errorw("加载支付宝的公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServerConfig.AliPayInfo.NotifyUrl
	p.ReturnURL = global.ServerConfig.AliPayInfo.ReturnUrl
	p.Subject = "GoShop-" + rsp.OrderInfo.OrderSn
	p.OutTradeNo = rsp.OrderInfo.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(rsp.OrderInfo.Total), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("生成支付url失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	respMap["alipay_url"] = url.String()
	ctx.JSON(http.StatusOK, respMap)
}

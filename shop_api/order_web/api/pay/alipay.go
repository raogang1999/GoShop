package pay

import (
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"net/http"
	"shop_api/order_web/global"
	"shop_api/order_web/proto"
)

func Notify(ctx *gin.Context) {
	//支付宝回调通知，
	//验证是否为支付宝发送的通知
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

	noti, err := client.GetTradeNotification(ctx.Request)
	if err != nil {
		zap.S().Errorw("支付宝回调通知失败,", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	_, err = global.OrderSrvClient.UpdateOrderStatus(ctx, &proto.OrderStatus{
		OrderSn: noti.OutTradeNo,
		Status:  string(noti.TradeStatus),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	ctx.String(http.StatusOK, "success")
}

package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/order_web/api/order"
	"shop_api/order_web/api/pay"
	"shop_api/order_web/middlewares"
)

func InitOrderRoute(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("order").Use(middlewares.JWTAuth())
	{
		OrderRouter.GET("", order.List)
		OrderRouter.POST("", order.NewOrder)
		OrderRouter.GET("/:id", order.GetOrderById)

	}
	PayRouter := Router.Group("pay")
	{

		PayRouter.POST("alipay/notify", pay.Notify)
	}
}

package initialize

import (
	"github.com/gin-gonic/gin"
	"shop_api/order_web/middlewares"
	allRouter "shop_api/order_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())
	allRouter.InitHealthRouter(Router) // 注册健康检查路由
	ApiGroup := Router.Group("o/v1")
	{
		allRouter.InitOrderRoute(ApiGroup)
		allRouter.InitShopCartRouter(ApiGroup)
	}
	return Router
}

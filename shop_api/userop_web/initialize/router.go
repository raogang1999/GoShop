package initialize

import (
	"github.com/gin-gonic/gin"
	"shop_api/userop_web/middlewares"
	allRouter "shop_api/userop_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())
	allRouter.InitHealthRouter(Router) // 注册健康检查路由
	ApiGroup := Router.Group("/up/v1")
	allRouter.InitUserFavRouter(ApiGroup)
	allRouter.InitMessageRouter(ApiGroup)
	allRouter.InitAddressRouter(ApiGroup)

	return Router
}

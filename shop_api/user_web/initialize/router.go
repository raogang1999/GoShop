package initialize

import (
	"github.com/gin-gonic/gin"
	"shop_api/user_web/middlewares"
	allRouter "shop_api/user_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())
	allRouter.InitHealthRouter(Router) // 注册健康检查路由
	ApiGroup := Router.Group("u/v1")
	{
		allRouter.InitUserRouter(ApiGroup) // 注册用户路由
		allRouter.InitBaseRouter(ApiGroup) // 注册基础路由

	}
	return Router
}

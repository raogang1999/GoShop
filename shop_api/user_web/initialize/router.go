package initialize

import (
	"github.com/gin-gonic/gin"
	"shop_api/user_web/middlewares"
	userRouter "shop_api/user_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("u/v1")
	{
		userRouter.InitUserRouter(ApiGroup) // 注册用户路由
		userRouter.InitBaseRouter(ApiGroup) // 注册基础路由
	}
	return Router
}

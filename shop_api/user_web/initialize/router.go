package initialize

import (
	"github.com/gin-gonic/gin"
	userRouter "shop_api/user_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	ApiGroup := Router.Group("u/v1")
	{
		userRouter.InitUserRouter(ApiGroup) // 注册用户路由
	}
	return Router
}

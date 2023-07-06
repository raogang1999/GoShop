package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop_api/user_web/api"
)

// 为了不多个实例化，传入router
func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Info("配置用户路由，请求")
	{
		UserRouter.GET("list", api.GetUserList)
	}
}

package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/userop_web/api/message"
	"shop_api/userop_web/middlewares"
)

func InitMessageRouter(Router *gin.RouterGroup) {
	MessageRouter := Router.Group("message").Use(middlewares.JWTAuth())
	{
		MessageRouter.GET("", message.List)
		MessageRouter.POST("", message.New)
	}
}

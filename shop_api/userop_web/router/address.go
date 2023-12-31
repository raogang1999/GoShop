package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/userop_web/api/address"
	"shop_api/userop_web/middlewares"
)

func InitAddressRouter(Router *gin.RouterGroup) {
	AddressRouter := Router.Group("address")
	{
		AddressRouter.GET("", middlewares.JWTAuth(), address.List)
		AddressRouter.DELETE("/:id", middlewares.JWTAuth(), address.Delete)
		AddressRouter.POST("", middlewares.JWTAuth(), address.New)
		AddressRouter.PUT("/:id", middlewares.JWTAuth(), address.Update)
	}
}

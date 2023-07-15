package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/order_web/api/shop_cart"
	"shop_api/order_web/middlewares"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	ShopCartRouter := Router.Group("shop_cart").Use(middlewares.JWTAuth())
	{
		ShopCartRouter.GET("", shop_cart.List)             //购物车列表
		ShopCartRouter.DELETE("/:id", shop_cart.Delete)    //删除购物车
		ShopCartRouter.POST("", shop_cart.NewCart)         //添加购物车
		ShopCartRouter.PATCH("/:id", shop_cart.UpdateCart) //更新购物车
	}
}

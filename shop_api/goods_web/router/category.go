package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/goods_web/api/category"
)

func InitCategoryRouter(Router *gin.RouterGroup) {
	CategoryRouter := Router.Group("category")
	{
		CategoryRouter.GET("", category.List)
		CategoryRouter.POST("", category.New)
		CategoryRouter.GET("/:id", category.Detail)
		CategoryRouter.DELETE("/:id", category.Delete)
		CategoryRouter.PUT("/:id", category.Update)
	}
}

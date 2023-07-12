package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/goods_web/api/goods"
	"shop_api/goods_web/middlewares"
)

// 为了不多个实例化，传入router
func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	{
		GoodsRouter.GET("", goods.List)
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)               //需要权限
		GoodsRouter.GET("/:id", goods.Detail)                                                           //商品详情
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete)      //删除商品
		GoodsRouter.GET("/:id/stocks", goods.Stock)                                                     //商品库存
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus) //更新商品状态
		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)         //更新商品
	}
}

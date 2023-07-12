package initialize

import (
	"github.com/gin-gonic/gin"
	"shop_api/goods_web/middlewares"
	allRouter "shop_api/goods_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())
	allRouter.InitHealthRouter(Router) // 注册健康检查路由
	ApiGroup := Router.Group("g/v1")
	{
		allRouter.InitGoodsRouter(ApiGroup)    // 注册商品路由
		allRouter.InitCategoryRouter(ApiGroup) //注册分类路由
		allRouter.InitBannerRouter(ApiGroup)   //注册轮播图路由
		allRouter.InitBrandRouter(ApiGroup)    //注册品牌路由
	}
	return Router
}

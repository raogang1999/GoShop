package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/userop_web/api/user_fav"
	"shop_api/userop_web/middlewares"
)

func InitUserFavRouter(Router *gin.RouterGroup) {
	UserFavRouter := Router.Group("userfavs")
	{
		UserFavRouter.DELETE("/:id", middlewares.JWTAuth(), user_fav.Delete)
		UserFavRouter.GET("/:id", middlewares.JWTAuth(), user_fav.Detail)
		UserFavRouter.POST("", middlewares.JWTAuth(), user_fav.New)
		UserFavRouter.GET("", middlewares.JWTAuth(), user_fav.List)
	}
}

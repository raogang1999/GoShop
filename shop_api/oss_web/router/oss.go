package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/oss_web/api"
)

func InitOssRouter(Router *gin.RouterGroup) {
	OssRouter := Router.Group("oss")
	{
		OssRouter.GET("token", api.Token)         // 获取token
		OssRouter.POST("/callback", api.Callback) // 上传回调
	}
}

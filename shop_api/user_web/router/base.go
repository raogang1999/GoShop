package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/user_web/api"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("base")
	{
		baseRouter.GET("captcha", api.GetCaptcha)
		baseRouter.POST("send_sms", api.SendSms)
	}
}

package router

import "github.com/gin-gonic/gin"

func InitHealthRouter(Router gin.IRouter) {
	Router.GET("health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "ok",
		})
	})
}

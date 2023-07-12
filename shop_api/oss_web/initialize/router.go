package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_api/oss_web/middlewares"
	allRouter "shop_api/oss_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	Router.LoadHTMLFiles(fmt.Sprintf("oss_web/templates/index.html"))
	// 配置静态文件夹路径 第一个参数是api，第二个是文件夹路径
	Router.StaticFS("/static", http.Dir(fmt.Sprintf("oss_web/static")))
	// GET：请求方式；/hello：请求的路径
	// 当客户端以GET方法请求/hello路径时，会执行后面的匿名函数
	Router.GET("", func(c *gin.Context) {
		// c.JSON：返回JSON格式的数据
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "posts/index",
		})
	})

	Router.Use(middlewares.Cors())
	allRouter.InitHealthRouter(Router) // 注册健康检查路由

	api := Router.Group("/oss/v1")
	{
		allRouter.InitOssRouter(api) // 注册oss路由
	}
	return Router
}

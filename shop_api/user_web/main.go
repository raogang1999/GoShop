package main

import (
	"fmt"
	"go.uber.org/zap"
	"shop_api/user_web/global"
	"shop_api/user_web/initialize"
)

func main() {

	// 初始化日志
	initialize.InitLogger()
	// 初始化配置文件
	initialize.InitConfig()
	// 初始化路由
	Router := initialize.Routers()

	zap.S().Debugf("启动端口: %d", global.ServerConfig.Port)

	// 启动服务
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); nil != err {
		zap.S().Panic("启动失败: ", err.Error())
	}

}

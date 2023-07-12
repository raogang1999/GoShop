package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/inner/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"shop_api/goods_web/global"
	"shop_api/goods_web/initialize"
	"shop_api/goods_web/utils"
	"shop_api/goods_web/utils/register/consul"
	"syscall"
)

func main() {

	// 1. 初始化日志
	initialize.InitLogger()
	// 2. 初始化配置文件
	initialize.InitConfig()
	// 3. 初始化路由
	Router := initialize.Routers()
	//4. 初始化验证器
	err := initialize.InitTrans("zh")
	if err != nil {
		panic(err)
		//zap.S().Errorw("初始化验证器失败", "msg", err.Error())
	}
	//5. 初始化srv连接
	initialize.InitSrvConn()

	//获取可用端口
	viper.AutomaticEnv()
	debug := viper.GetBool("SHOP_DEBUG")
	//希望开发环境是固定端口的
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port

		}
	}

	zap.S().Debugf("启动端口: %d", global.ServerConfig.Port)
	//服务注册
	register_client := consul.NewRegistry(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV5(uuid.NamespaceOID, global.ServerConfig.Name))
	err = register_client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败: ", err.Error())
	}

	go func() {
		// 启动服务
		if err := Router.Run(fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port)); nil != err {
			zap.S().Panic("启动失败: ", err.Error())
		}

	}()
	//优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	err = register_client.DeRegister(serviceId)
	if err != nil {
		zap.S().Panic("注销失败: ", err.Error())
	} else {
		zap.S().Info("注销成功: ", serviceId)
	}
}

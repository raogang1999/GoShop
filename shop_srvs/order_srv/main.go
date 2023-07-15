package main

import (
	"flag"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/inner/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"shop_srvs/order_srv/global"
	"shop_srvs/order_srv/handler"
	"shop_srvs/order_srv/initialize"
	"shop_srvs/order_srv/proto"
	"shop_srvs/order_srv/utils"
	"shop_srvs/order_srv/utils/register/consul"
	"syscall"
)

func main() {
	IP := flag.String("ip", "192.168.112.1", "ip地址")
	Port := flag.Int("port", 0, "端口号")

	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitRedis()
	initialize.InitSrvConn() //初始化微服务客户端
	zap.S().Info(global.SeverConfig)
	flag.Parse()
	zap.S().Info("IP: ", *IP)
	zap.S().Info("Port: ", *Port)

	server := grpc.NewServer()
	//商品服务
	proto.RegisterOrderServer(server, &handler.OrderServer{})
	//健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	register_client := consul.NewRegistry(global.SeverConfig.ConsulInfo.Host, global.SeverConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV5(uuid.NamespaceOID, global.SeverConfig.Name))
	err := register_client.Register(global.SeverConfig.Host, *Port, global.SeverConfig.Name, global.SeverConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败: ", err.Error())
	}
	zap.S().Info("服务注册成功,服务端口: ", *Port)

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic(err)
	}

	go func() {
		err = server.Serve(listen)
		if err != nil {
			panic(err)
		}
	}()

	//接受终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}

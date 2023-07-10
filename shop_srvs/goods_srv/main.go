package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"shop_srvs/goods_srv/global"
	"shop_srvs/goods_srv/handler"
	"shop_srvs/goods_srv/initialize"
	"shop_srvs/goods_srv/proto"
	"shop_srvs/goods_srv/utils"
	"syscall"
)

func main() {
	IP := flag.String("ip", "192.168.112.1", "ip地址")
	Port := flag.Int("port", 50051, "端口号")

	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	zap.S().Info(global.SeverConfig)
	flag.Parse()
	zap.S().Info("IP: ", *IP)
	zap.S().Info("Port: ", *Port)

	server := grpc.NewServer()
	//商品服务
	proto.RegisterGoodsServer(server, &handler.GoodsServer{})
	//健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	//服务注册
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", global.SeverConfig.ConsulInfo.Host, global.SeverConfig.ConsulInfo.Port)

	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.SeverConfig.Host, *Port),
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	registration := new(api.AgentServiceRegistration)
	registration.Name = global.SeverConfig.Name
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = global.SeverConfig.Tags
	registration.Address = global.SeverConfig.Host
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

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
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}

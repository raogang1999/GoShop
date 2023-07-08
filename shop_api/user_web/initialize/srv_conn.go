package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop_api/user_web/global"
	"shop_api/user_web/proto"

	_ "github.com/mbobakov/grpc-consul-resolver"
)

func InitSrvConn() {
	//用户连接
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】", err.Error())
	}
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}

func InitSrvConn2() {
	//从注册中心拉去微服务的信息
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	//获取grpc服务的信息
	userSrvHost := ""
	userSrvPort := 0

	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		panic(err)
	}
	//可能有多个，可以负载均衡，但目前只用一个，
	for _, v := range data {
		userSrvHost = v.Address
		userSrvPort = v.Port
		break
	}
	if userSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接 [用户服务失败]", "msg", err.Error())
	}

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 [用户服务失败]", "msg", err.Error())
	}
	//1. 下线了，2. 服务重启了，3. 端口ip变了
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}

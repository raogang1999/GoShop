package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop_api/goods_web/global"
	"shop_api/goods_web/proto"

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
	global.GoodsSrvClient = proto.NewGoodsClient(userConn)
}

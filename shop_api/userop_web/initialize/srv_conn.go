package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop_api/userop_web/global"
	"shop_api/userop_web/proto"

	_ "github.com/mbobakov/grpc-consul-resolver"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	//商品微服务初始化
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[userop_web] 连接 [Goods_srv 服务失败]", err.Error())
	}
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	//用户操作微服务初始化
	userOpConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserOpSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[userop_web] 连接 [UserOp_srv 服务失败]", err.Error())
	}
	global.UserFavSrvClient = proto.NewUserFavClient(userOpConn)
	global.AddressSrvClient = proto.NewAddressClient(userOpConn)
	global.MessageSrvClient = proto.NewMessageClient(userOpConn)
}

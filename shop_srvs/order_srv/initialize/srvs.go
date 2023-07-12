package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop_srvs/order_srv/global"
	"shop_srvs/order_srv/proto"
)

func InitSrvConn() {
	//初始化第三方的微服务的client

	//获取配置中心的配置，知道服务的地址，
	consulInfo := global.SeverConfig.ConsulInfo
	//商品微服务
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.SeverConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)) //负载均衡
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 [商品服务失败]", err.Error())
	}
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	//库存微服务
	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.SeverConfig.InventorySrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)) //负载均衡
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 [库存服务失败]", err.Error())
	}
	global.InventorySrvClient = proto.NewInventoryClient(inventoryConn)
}

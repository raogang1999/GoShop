package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop_api/order_web/global"
	"shop_api/order_web/proto"

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
		zap.S().Fatal("[Order_web] 连接 [Goods_srv 服务失败]", err.Error())
	}
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)
	///库存微服务
	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[Order_web] 连接 [Inventory_srv 服务失败]", err.Error())
	}
	global.InventorySrvClient = proto.NewInventoryClient(inventoryConn)

	//用户连接
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.OrderSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[Order_web] 连接 [Order_srv 服务失败]", err.Error())
	}
	global.OrderSrvClient = proto.NewOrderClient(userConn)
}

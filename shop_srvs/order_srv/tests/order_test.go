package tests

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"shop_srvs/order_srv/proto"
	"testing"
)

var orderClient proto.OrderClient

// 获取userClient
func Init() {
	var err error
	conn, err := grpc.Dial("192.168.112.1:8039", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	orderClient = proto.NewOrderClient(conn)
}

// 查询订单
func TestOrderList(t *testing.T) {
	Init()
	resp, err := orderClient.OrderList(context.Background(), &proto.OrderFilterRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Total)

	for _, v := range resp.Data {
		fmt.Println(v)
	}
}

// 新建订单
func TestCreateOrder(t *testing.T) {
	Init()
	resp, err := orderClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  1,
		Address: "北京",
		Name:    "张三",
		Mobile:  "18888888888",
		Post:    "100000",
	})
	if err != nil {
		panic(err)
	}
	t.Log(resp)
}

// 订单详情
func TestOrderDetail(t *testing.T) {
	Init()
	resp, err := orderClient.OrderDetail(context.Background(), &proto.OrderRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.OrderInfo)
	fmt.Println(resp.Goods)
}

// 添加购物车
func TestCreateCartItem(t *testing.T) {
	Init()
	resp, err := orderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  1,
		GoodsId: 427,
		Nums:    1,
	})
	if err != nil {
		panic(err)
	}
	t.Log(resp)
}

// 查询购物车
func TestCartItemList(t *testing.T) {
	Init()
	resp, err := orderClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: 2,
	})
	if err != nil {
		panic(err)
	}
	t.Log(resp)
}

// 更新购物车
func TestUpdateCartItem(t *testing.T) {
	Init()
	resp, err := orderClient.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		Id:      2,
		Nums:    1,
		Checked: true,
	})
	if err != nil {
		panic(err)
	}
	t.Log(resp)
}

// 删除购物车
func TestDeleteCartItem(t *testing.T) {
	Init()
	resp, err := orderClient.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		Id: 2,
	})
	if err != nil {
		panic(err)
	}
	t.Log(resp)
}

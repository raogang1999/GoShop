package tests

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	model2 "shop_srvs/userop_srv/model"
	"shop_srvs/userop_srv/proto"
	"testing"
)

var userFavClient proto.UserFavClient
var messageClient proto.MessageClient
var addrClient proto.AddressClient

// 获取userClient
func Init(types string) {

	conn, err := grpc.Dial("192.168.112.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	if types == "userFav" {
		userFavClient = proto.NewUserFavClient(conn)
	}
	if types == "message" {
		messageClient = proto.NewMessageClient(conn)
	}
	if types == "address" {
		addrClient = proto.NewAddressClient(conn)
	}
}

// 测试地址获取
func TestGetAddressList(t *testing.T) {
	Init("address")
	resp, err := addrClient.GetAddressList(context.Background(), &proto.AddressRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

// 测试地址创建
func TestCreateAddress(t *testing.T) {
	Init("address")
	resp, err := addrClient.CreateAddress(context.Background(), &proto.AddressRequest{
		UserId:       1,
		Province:     "广东省",
		City:         "深圳市",
		District:     "南山区",
		Address:      "深圳市南山区",
		SignerName:   "张三",
		SignerMobile: "12345678901",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

// 测试地址更新
func TestUpdateAddress(t *testing.T) {
	Init("address")
	resp, err := addrClient.UpdateAddress(context.Background(), &proto.AddressRequest{
		UserId:       1,
		Id:           1,
		Province:     "广东省",
		City:         "深圳市",
		District:     "南山区",
		Address:      "深圳市南山区1",
		SignerName:   "张三",
		SignerMobile: "12345678901",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

// 测试地址删除
func TestDeleteAddress(t *testing.T) {
	Init("address")
	resp, err := addrClient.DeleteAddress(context.Background(), &proto.AddressRequest{
		Id:     1,
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

//

// 测试用户收藏
// 添加收藏
func TestCreateUserFav(t *testing.T) {
	Init("userFav")
	resp, err := userFavClient.AddUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  1,
		GoodsId: 423,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

// 获取收藏列表
func TestUserFavList(t *testing.T) {
	Init("userFav")
	resp, err := userFavClient.GetFavList(context.Background(), &proto.UserFavRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)

	}
	fmt.Println(resp)
}

// 获取详情
func TestUserFav(t *testing.T) {
	Init("userFav")
	resp, err := userFavClient.GetUserFavDetail(context.Background(), &proto.UserFavRequest{
		UserId:  1,
		GoodsId: 421,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

// 删除
func TestDeleteUserFav(t *testing.T) {
	Init("userFav")
	resp, err := userFavClient.DeleteUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  1,
		GoodsId: 421,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

// 测试消息
// 添加消息
func TestCreateMessage(t *testing.T) {
	Init("message")
	resp, err := messageClient.CreateMessage(context.Background(), &proto.MessageRequest{
		UserId:      1,
		MessageType: model2.LEAVING_MESSAGES,
		Subject:     "测试消息2",
		Message:     "这是一条测试消息",
		File:        "test.jpg",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

// 获取消息列表
func TestGetMessageList(t *testing.T) {
	Init("message")
	resp, err := messageClient.MessageList(context.Background(), &proto.MessageRequest{
		UserId: 1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

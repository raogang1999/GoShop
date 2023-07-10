package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"shop_srvs/user_srv/proto"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

// 获取userClient
func Init() {
	var err error
	conn, err = grpc.Dial("192.168.112.1:3014", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

func TestFindUserByMobile() {
	resp, err := userClient.GetUserByMobile(context.Background(), &proto.UserMobileRequest{
		Mobile: "18888888880",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

// 测试通过id用户查询
func TestFindUserById() {
	resp, err := userClient.GetUserById(context.Background(), &proto.IdRequest{
		Id: 1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

// 测试创建用户
func TestCreateUser() {
	resp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Mobile:   "12345678910",
		NickName: "Tom",
		Password: "admin123",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

// 测试更新用户
func TestUpdateUser() {
	resp, err := userClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       1,
		NickName: "测试",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

// 测试密码校验和用户列表
func TestGetUserList() {
	list, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 2,
	})
	if err != nil {
		panic(err)
	}
	for _, user := range list.Data {
		fmt.Println(user.Mobile, user.NickName, user.Password)
		//密码校验
		check_resp, err := userClient.CheckPassword(context.Background(), &proto.CheckPasswordInfo{
			Password:          "admin123",
			EncryptedPassword: user.Password,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(check_resp.Success)
	}
}

func main() {
	Init()
	//TestGetUserList()
	TestFindUserByMobile()
	TestFindUserById()
	//TestCreateUser()
	//TestUpdateUser()
	TestGetUserList()
	conn.Close()
}

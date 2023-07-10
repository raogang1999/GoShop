package tests

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"shop_srvs/goods_srv/proto"
	"testing"
)

var brandClient proto.GoodsClient
var conn *grpc.ClientConn

// 获取userClient
func Init() {
	var err error
	conn, err = grpc.Dial("192.168.112.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	brandClient = proto.NewGoodsClient(conn)
}

func TestBrandList(t *testing.T) {
	Init()
	resp, err := brandClient.BrandList(context.Background(), &proto.BrandFilterRequest{
		Page:        1,
		PagePerNums: 2,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("总计：", resp.Total)
	for _, brand := range resp.Data {
		fmt.Println(brand)
	}
}
func TestBrandCreate(t *testing.T) {
	Init()
	resp, err := brandClient.CreateBrand(context.Background(), &proto.BrandRequest{
		Name: "Ali1",
		Logo: "测试logo",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)

}

func TestDeleteBrand(t *testing.T) {
	Init()
	resp, err := brandClient.DeleteBrand(context.Background(), &proto.BrandRequest{
		Id: 1128,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func TestUpdateBrand(t *testing.T) {
	Init()
	resp, err := brandClient.UpdateBrand(context.Background(), &proto.BrandRequest{
		Id:   1110,
		Name: "测试1",
		Logo: "测试logo12",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

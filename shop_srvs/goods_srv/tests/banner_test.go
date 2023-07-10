package tests

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"shop_srvs/goods_srv/proto"
	"testing"
)

func TestBannerList(t *testing.T) {
	Init()
	list, err := brandClient.BannerList(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}
	for _, banner := range list.Data {
		t.Log(banner)
	}
}

func TestCreateBanner(t *testing.T) {
	Init()
	resp, err := brandClient.CreateBanner(context.Background(), &proto.BannerRequest{
		Image: "测试图片",
		Url:   "测试url",
		Index: 1,
	})
	if err != nil {
		panic(err)
	}
	t.Log(resp)
}

func TestDeleteBanner(t *testing.T) {
	Init()
	resp, err := brandClient.DeleteBanner(context.Background(), &proto.BannerRequest{
		Id: 2,
	})
	if err != nil {
		panic(err)
	}
	t.Log(resp)
}

func TestUpdateBanner(t *testing.T) {
	Init()
	resp, err := brandClient.UpdateBanner(context.Background(), &proto.BannerRequest{
		Id:    1,
		Image: "测试图片211",
		Url:   "测试url2",
		Index: 2,
	})
	if err != nil {
		panic(err)
	}
	t.Log(resp)
}

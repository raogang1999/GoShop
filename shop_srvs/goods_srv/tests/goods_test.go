package tests

import (
	"context"
	"fmt"
	"shop_srvs/goods_srv/proto"
	"testing"
)

func TestGetGoodsList(t *testing.T) {
	Init()
	rsp, err := brandClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		TopCategory: 130361,
		PriceMin:    90,
		//Keywords:    "深海速冻",
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rsp.Total)
	for _, v := range rsp.Data {
		fmt.Println(v)
	}
}

func TestBatchGetGoods(t *testing.T) {
	Init()
	rsp, err := brandClient.BatchGetGoods(context.Background(), &proto.BatchGetGoodsInfo{
		Id: []int32{421, 422, 423},
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rsp)
}

func TestGetGoodsDetail(t *testing.T) {
	Init()
	rsp, err := brandClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: 421,
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rsp)

}

func TestCreateGoods(t *testing.T) {
	Init()
	rsp, err := brandClient.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		Name:            "测试商品ES",
		GoodsSn:         "123456",
		Stocks:          54,
		MarketPrice:     99,
		ShopPrice:       88,
		GoodsBrief:      "测试商品ES",
		GoodsDesc:       "测试商品",
		ShipFree:        true,
		Images:          []string{"https://img.alicdn.com/imgextra/i1/2200722652651/O1CN01Q4Q4Qq1YQ1Q8QYQ1Y_!!2200722652651.jpg_430x430q90.jpg"},
		DescImages:      []string{""},
		GoodsFrontImage: "",
		IsNew:           true,
		IsHot:           true,
		OnSale:          true,
		CategoryId:      136832,
		BrandId:         1070,
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rsp)
}

func TestUpdateGoods(t *testing.T) {
	Init()
	rsp, err := brandClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              841,
		Name:            "New 商品",
		GoodsSn:         "123456asdf",
		Stocks:          99,
		MarketPrice:     99,
		ShopPrice:       88,
		GoodsBrief:      "测试商品",
		GoodsDesc:       "测试商品",
		ShipFree:        true,
		Images:          []string{"https://img.alicdn.com/imgextra/i1/2200722652651/O1CN01Q4Q4Qq1YQ1Q8QYQ1Y_!!2200722652651.jpg_430x430q90.jpg"},
		DescImages:      []string{""},
		GoodsFrontImage: "",
		IsNew:           true,
		IsHot:           true,
		OnSale:          true,
		CategoryId:      136832,
		BrandId:         1070,
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(rsp)
}

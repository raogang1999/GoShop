package tests

import (
	"context"
	"shop_srvs/goods_srv/proto"
	"testing"
)

func TestCategoryBrandList(t *testing.T) {
	Init()
	res, err := brandClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{
		Pages:       1,
		PagePerNums: 3,
	})
	if err != nil {
		panic(err)
	}
	t.Log(res)

}

func TestCategoryBrandCreate(t *testing.T) {
	Init()
	res, err := brandClient.CreateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		BrandId:    701,
		CategoryId: 136850,
	})
	if err != nil {
		panic(err)
	}
	t.Log(res)
}

func TestCategoryBrandDelete(t *testing.T) {
	Init()
	res, err := brandClient.DeleteCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		Id: 1,
	})
	if err != nil {
		panic(err)
	}
	t.Log(res)
}

func TestCategoryBrandUpdate(t *testing.T) {
	Init()
	res, err := brandClient.UpdateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		Id:         1,
		BrandId:    701,
		CategoryId: 136850,
	})
	if err != nil {
		panic(err)
	}
	t.Log(res)
}

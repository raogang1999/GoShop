package tests

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"shop_srvs/goods_srv/proto"
	"testing"
)

func TestGetAllCategorysList(t *testing.T) {
	Init()
	list, err := brandClient.GetAllCategorysList(context.Background(), &empty.Empty{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(list.JsonData)

}

func TestGetSubCategory(t *testing.T) {
	Init()
	list, err := brandClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{Id: 136604})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(list.Info)
}

func TestCreateCategory(t *testing.T) {
	Init()
	list, err := brandClient.CreateCategory(context.Background(), &proto.CategoryInfoRequest{
		Name:           "测试分类23",
		Level:          2,
		ParentCategory: 238010,
		IsTab:          false,
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(list)
}

func TestDeleteCategory(t *testing.T) {
	Init()
	list, err := brandClient.DeleteCategory(context.Background(), &proto.DeleteCategoryRequest{Id: 238010})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(list)
}

func TestUpdateCategory(t *testing.T) {
	Init()
	list, err := brandClient.UpdateCategory(context.Background(), &proto.CategoryInfoRequest{
		Id:             238011,
		Name:           "测试分类1",
		Level:          2,
		ParentCategory: 135487,
		IsTab:          true,
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(list)
}

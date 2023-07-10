package handler

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"shop_srvs/goods_srv/global"
	"shop_srvs/goods_srv/model"
	"shop_srvs/goods_srv/proto"
)

// 商品分类
//GetAllCategorysList(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*CategoryListResponse, error)
//// 获取商品子分类列表
//GetSubCategory(ctx context.Context, in *CategoryListRequest, opts ...grpc.CallOption) (*SubCategoryListResponse, error)
//// 创建商品分类
//CreateCategory(ctx context.Context, in *CategoryInfoRequest, opts ...grpc.CallOption) (*CategoryInfoResponse, error)
//// 删除商品分类
//DeleteCategory(ctx context.Context, in *DeleteCategoryRequest, opts ...grpc.CallOption) (*empty.Empty, error)
//// 更新商品分类信息
//UpdateCategory(ctx context.Context, in *CategoryInfoRequest, opts ...grpc.CallOption) (*empty.Empty, error)
// 品牌分类

// 商品分类
func (s *GoodsServer) GetAllCategorysList(context.Context, *empty.Empty) (*proto.CategoryListResponse, error) {
	/*
		[
			{
				"id": 1,
				"name": "手机",
				"level": 1,
				"parent_category_id": 0,
				"sub_category_list": [
					{
						"id": 2,
						"name": "苹果",
						"level": 2,
						"parent_category_id": 1,
						"sub_category_list": [	....
		]
	*/
	var categorys []model.Category
	//global.DB.Preload("SubCategory").Find(&categorys) //只能一级
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	b, _ := json.Marshal(&categorys)
	return &proto.CategoryListResponse{JsonData: string(b)}, nil

}

// 获取商品子分类列表
func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	categoryListResponse := proto.SubCategoryListResponse{}
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有找到商品分类信息")
	}
	categoryListResponse.Info = &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}
	var subCategorys []model.Category
	var subCategoryResponse []*proto.CategoryInfoResponse
	preloads := "SubCategory"
	if category.Level == 1 {
		preloads = "SubCategory.SubCategory"
	}
	global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Preload(preloads).Find(&subCategorys)
	for _, subCategory := range subCategorys {
		subCategoryResponse = append(subCategoryResponse, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryID,
		})
	}
	categoryListResponse.SubCategorys = subCategoryResponse
	return &categoryListResponse, nil
}

// 创建商品分类
func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{}

	category.Name = req.Name
	category.Level = req.Level
	category.IsTab = req.IsTab

	if req.Level != 1 {
		//去查询父类目是否存在

		category.ParentCategoryID = req.ParentCategory
	}
	if result := global.DB.Save(&category); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建商品分类失败")
	}

	return &proto.CategoryInfoResponse{Id: int32(category.ID)}, nil
}

// 删除商品分类
func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}

// 更新商品分类信息
func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	var category model.Category

	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}

	global.DB.Save(&category)

	return &emptypb.Empty{}, nil
}

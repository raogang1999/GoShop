package handler

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop_srvs/goods_srv/global"
	"shop_srvs/goods_srv/model"
	"shop_srvs/goods_srv/proto"
)

// 品牌
// 创建品牌
func (s *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	//不能重名
	result := global.DB.Where("name = ?", req.Name).First(&model.Brands{})
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}

	brand := &model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	global.DB.Save(&brand)
	return &proto.BrandInfoResponse{
		Id: brand.ID,
	}, nil
}

// 删除品牌
func (s *GoodsServer) DeleteBrand(ctx context.Context, in *proto.BrandRequest) (*empty.Empty, error) {
	if result := global.DB.Unscoped().Delete(&model.Brands{}, in.GetId()); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	return &empty.Empty{}, nil
}

// 更新品牌
func (s *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*empty.Empty, error) {
	brands := model.Brands{
		Name: req.Name,
	}

	if req.Name != "" {
		brands.Name = req.Name
	}
	//保证品牌名唯一
	if result := global.DB.Where("name = ?", req.Name).First(&brands); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}

	//判断品牌是否存在
	if result := global.DB.First(&brands, req.GetId()); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	if req.Logo != "" {
		brands.Logo = req.Logo
	}
	if req.Name != "" {
		brands.Name = req.Name
	}

	if result := global.DB.Save(&brands); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "更新品牌失败")
	}
	return &empty.Empty{}, nil
}

// 获取品牌列表
func (s *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	var brandListResponse proto.BrandListResponse
	var brands []model.Brands

	//分页
	result := global.DB.Scopes(Paginate(int(req.Page), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}
	var count int64
	global.DB.Model(&model.Brands{}).Count(&count)
	brandListResponse.Total = int32(count)
	var brandResponses []*proto.BrandInfoResponse
	for _, brand := range brands {
		brandResponses = append(brandResponses, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	brandListResponse.Data = brandResponses
	return &brandListResponse, nil
}

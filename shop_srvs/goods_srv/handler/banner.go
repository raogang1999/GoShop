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

// 轮播图

// 查询所有轮播图
func (s *GoodsServer) BannerList(ctx context.Context, e *empty.Empty) (*proto.BannerListResponse, error) {
	var bannerListResponse proto.BannerListResponse
	var banners []model.Banner
	result := global.DB.Find(&banners)
	bannerListResponse.Total = int32(result.RowsAffected)
	for _, banner := range banners {
		bannerListResponse.Data = append(bannerListResponse.Data, &proto.BannerResponse{
			Id:    banner.ID,
			Image: banner.Image,
			Index: banner.Index,
			Url:   banner.Url,
		})
	}
	return &bannerListResponse, nil
}

func (s *GoodsServer) CreateBanner(ctx context.Context, req *proto.BannerRequest) (*proto.BannerResponse, error) {
	var banner model.Banner
	banner.Image = req.Image
	banner.Index = req.Index
	banner.Url = req.Url

	if result := global.DB.Save(&banner); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, "创建轮播图失败")
	}
	return &proto.BannerResponse{
		Id:    banner.ID,
		Image: banner.Image,
		Index: banner.Index,
		Url:   banner.Url,
	}, nil
}
func (s *GoodsServer) DeleteBanner(ctx context.Context, req *proto.BannerRequest) (*empty.Empty, error) {

	if result := global.DB.Delete(&model.Banner{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}

	return &empty.Empty{}, nil
}
func (s *GoodsServer) UpdateBanner(ctx context.Context, req *proto.BannerRequest) (*empty.Empty, error) {
	var banner model.Banner

	//判断轮播图是否存在
	if result := global.DB.First(&banner, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}
	banner.Image = req.Image
	banner.Index = req.Index
	banner.Url = req.Url
	if result := global.DB.Save(&banner); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, "更新轮播图失败")
	}

	return &empty.Empty{}, nil
}

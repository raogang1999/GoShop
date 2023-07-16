package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop_srvs/goods_srv/global"
	"shop_srvs/goods_srv/model"
	"shop_srvs/goods_srv/proto"
)

// 商品服务
type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brand.ID,
			Name: goods.Brand.Name,
			Logo: goods.Brand.Logo,
		},
	}
}

// // 获取商品列表
func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	//关键词搜索，查询新品，查询热门，通过价格区间筛选，通过商品分类筛选
	goodListRespone := &proto.GoodsListResponse{}

	var goods []model.Goods
	localDB := global.DB.Model(model.Goods{})

	//es复合查询
	q := elastic.NewBoolQuery()

	if req.Keywords != "" {
		//关键词搜索
		localDB = localDB.Where("name like ?", "%"+req.Keywords+"%")
		q.Must(elastic.NewMultiMatchQuery(req.Keywords, "name", "goods_brief"))

	}
	if req.IsHot {
		//查询热门
		localDB = localDB.Where("is_hot = ?", req.IsHot)
		q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.IsNew {
		//查询新品
		localDB = localDB.Where("is_new = ?", req.IsNew)
		q.Filter(elastic.NewTermQuery("is_new", req.IsNew))
	}
	if req.PriceMin > 0 {
		//通过价格区间筛选
		//localDB = localDB.Where("shop_price >= ?", req.PriceMin)
		q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}
	if req.PriceMax > 0 {
		//通过价格区间筛选
		//localDB = localDB.Where("shop_price <= ?", req.PriceMax)
		q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}
	if req.Brand > 0 {
		//localDB = localDB.Where("brand_id = ?", req.Brand)
		q.Filter(elastic.NewTermQuery("brand_id", req.Brand))
	}
	//子查询
	//通过商品分类筛选
	var subQuery string
	categoryIds := make([]interface{}, 0)
	if req.TopCategory > 0 {
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}
		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory)
		}
		type Result struct {
			Id int32 `json:"id"`
		}
		var results []Result
		global.DB.Raw(subQuery).Pluck("id", &categoryIds).Scan(&results)
		for _, re := range results {
			categoryIds = append(categoryIds, re.Id)
		}
		// 生成 terms查询
		q.Filter(elastic.NewTermsQuery("category_id", categoryIds...))
	}

	//分页处理
	if req.Pages == 0 {
		req.Pages = 1
	}
	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}

	result, err := global.EsClient.Search().Index(model.EsGoods{}.GetIndexName()).Query(q).From(int(req.Pages)).Size(int(req.PagePerNums)).Do(ctx)

	if err != nil {
		return nil, err
	}
	goodListRespone.Total = int32(result.Hits.TotalHits.Value)
	goodsIds := make([]int32, 0)
	for _, hit := range result.Hits.Hits {
		goods := model.Goods{}
		_ = json.Unmarshal(hit.Source, &goods)
		goodsIds = append(goodsIds, goods.ID) //获取到商品id
	}
	//
	sqlResults := localDB.Preload("Category").Preload("Brand").Find(&goods, goodsIds)
	if sqlResults.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	for _, good := range goods {
		data := ModelToResponse(good)
		goodListRespone.Data = append(goodListRespone.Data, &data)
	}

	return goodListRespone, nil

}

func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGetGoodsInfo) (*proto.GoodsListResponse, error) {
	goodsListResponse := &proto.GoodsListResponse{}
	var goods []model.Goods

	result := global.DB.Where(req.Id).Find(&goods)
	for _, good := range goods {
		data := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &data)
	}
	goodsListResponse.Total = int32(result.RowsAffected)
	return goodsListResponse, nil
}
func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods
	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	response := ModelToResponse(goods)
	return &response, nil
}
func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}

	//图片文件怎么上传的？
	//oss
	var goods model.Goods
	goods.Brand = brand
	goods.BrandID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale
	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "创建商品失败")
	}
	tx.Commit()
	return &proto.GoodsInfoResponse{
		Id: goods.ID,
	}, nil
}
func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*empty.Empty, error) {
	if result := global.DB.Delete(&model.Goods{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &empty.Empty{}, nil
}
func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*empty.Empty, error) {
	//查询商品存在
	var goods model.Goods
	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	//查询分类存在
	var category model.Category
	if result := global.DB.First(&category, goods.CategoryID); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	//查询品牌存在
	var brand model.Brands
	if result := global.DB.First(&brand, goods.BrandID); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	//更新商品
	goods.Brand = brand
	goods.BrandID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale

	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "更新商品失败")
	}
	tx.Commit()

	return &empty.Empty{}, nil
}

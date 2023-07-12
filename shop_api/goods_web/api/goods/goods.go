package goods

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"shop_api/goods_web/api"
	"shop_api/goods_web/forms"
	"shop_api/goods_web/global"
	"shop_api/goods_web/proto"
	"strconv"
)

// 更新商品
func Update(ctx *gin.Context) {
	//获取表单
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	//获取id
	id := ctx.Param("id")
	goodsId, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	//调用grpc接口
	if _, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              int32(goodsId),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}

// 商品状态
func UpdateStatus(ctx *gin.Context) {
	//获取表单
	goodsForm := forms.GoodsStatusForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	//获取id
	id := ctx.Param("id")
	goodsId, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	//调用grpc接口
	if _, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:     int32(goodsId),
		IsHot:  *goodsForm.IsHot,
		IsNew:  *goodsForm.IsNew,
		OnSale: *goodsForm.OnSale,
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})

}

// 库存
func Stock(ctx *gin.Context) {
	goodsId := ctx.Param("id")
	_, err := strconv.Atoi(goodsId)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	//Todo 商品库存
	return
}

// 删除
func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{
		Id: int32(i),
	})
	if err != nil {
		zap.S().Errorw("删除商品失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
	return
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	goodsId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "id 错误",
		})
		return
	}
	rsp, err := global.GoodsSrvClient.GetGoodsDetail(ctx, &proto.GoodInfoRequest{
		Id: int32(goodsId),
	})
	if err != nil {
		zap.S().Errorw("获取商品详情失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//返回数据

	//手动封装返回的数据
	responseInfo := map[string]interface{}{
		"id":          rsp.Id,
		"name":        rsp.Name,
		"goods_brief": rsp.GoodsBrief,
		"desc":        rsp.GoodsDesc,
		"ship_free":   rsp.ShipFree,
		"images":      rsp.Images,
		"desc_images": rsp.DescImages,
		"front_image": rsp.GoodsFrontImage,
		"shop_price":  rsp.ShopPrice,
		"category": map[string]interface{}{
			"id":   rsp.Category.Id,
			"name": rsp.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   rsp.Brand.Id,
			"name": rsp.Brand.Name,
			"logo": rsp.Brand.Logo,
		},
		"is_hot":  rsp.IsHot,
		"is_new":  rsp.IsNew,
		"on_sale": rsp.OnSale,
	}
	ctx.JSON(http.StatusOK, responseInfo)
}

func New(ctx *gin.Context) {
	//新建商品，
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	goodsClient := global.GoodsSrvClient
	rsp, err := goodsClient.CreateGoods(ctx, &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		GoodsDesc:       goodsForm.GoodsDesc,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,

		CategoryId: goodsForm.CategoryId,
		BrandId:    goodsForm.Brand,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//ToDo 商品的库存 分布式事务
	//返回结果
	ctx.JSON(http.StatusOK, rsp)

}

func List(ctx *gin.Context) {
	//商品的列表

	//构建请求
	request := &proto.GoodsFilterRequest{}

	//参数获取
	priceMin := ctx.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	request.PriceMin = int32(priceMinInt)
	priceMax := ctx.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	request.PriceMax = int32(priceMaxInt)

	isHot := ctx.DefaultQuery("ih", "0")
	if isHot == "1" {
		request.IsHot = true
	}
	isNew := ctx.DefaultQuery("in", "0")
	if isNew == "1" {
		request.IsNew = true
	}
	isTab := ctx.DefaultQuery("it", "0")
	if isTab == "1" {
		request.IsTab = true
	}

	categoryId := ctx.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	request.TopCategory = int32(categoryIdInt)

	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	pagePerNums := ctx.DefaultQuery("pnum", "0")
	pagePerNumsInt, _ := strconv.Atoi(pagePerNums)
	request.PagePerNums = int32(pagePerNumsInt)

	keywords := ctx.DefaultQuery("q", "")
	request.Keywords = keywords

	brandId := ctx.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	request.Brand = int32(brandIdInt)

	// 请求商品的服务srv
	rsp, err := global.GoodsSrvClient.GoodsList(ctx, request)
	if err != nil {
		zap.S().Errorw("[List] 查询 【商品列表失败】", "msg", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return //终止
	}

	//返回数据
	resMap := map[string]interface{}{
		"total": rsp.Total,
	}
	//手动封装返回的数据
	goodsList := make([]interface{}, 0)
	for _, value := range rsp.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"category": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}
	resMap["data"] = goodsList
	ctx.JSON(http.StatusOK, resMap)

}

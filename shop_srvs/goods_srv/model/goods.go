package model

import (
	"context"
	"gorm.io/gorm"
	"shop_srvs/goods_srv/global"
	"strconv"
)

/*
1. 分类信息表
类目级别，1，2，3
属性的值尽量不为null
*/
type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(20);not null;comment:'分类名称,不能为空'" json:"name"`
	ParentCategoryID *int32      `json:"parent"` //解决零值更新问题
	ParentCategory   *Category   `json:"-"`      //忽略
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID;default:NULL" json:"sub_category"`

	Level int32 `gorm:"type:int;not null;default:1;comment:'分类级别，1，2，3'" json:"level"`
	IsTab bool  `gorm:"type:tinyint;default:false;comment:'是否在首页进行tab展示，false不展示，true展示'" json:"is_tab"`
}

/*
2. 品牌信息表
品牌名称
品牌logo
*/
type Brands struct {
	BaseModel
	//主键自增
	Name string `gorm:"type:varchar(50);not null;comment:'品牌名称'"`
	Logo string `gorm:"type:varchar(200);not null;default:'';comment:'品牌logo'"`
}

/*
3. 品牌与分类关系表，多对多
*/
type GoodsCategoryBrand struct {
	BaseModel
	//外键id,联合索引
	CategoryID int32 `gorm:"type:int;index:idx_category_brand;"`
	Category   Category

	BrandID int32 `gorm:"type:int;index:idx_category_brand;"`
	Brand   Brands
}

func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

/*4. 轮播图
轮播图名称
轮播图图片
*/

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null;comment:'轮播图图片'"`
	Url   string `gorm:"type:varchar(200);not null;comment:'轮播图链接'"`
	Index int32  `gorm:"type:int;not null;default:1;comment:'轮播图顺序'"`
}

/*
5. 商品表,
对于图片需要自定义类型，
自定义类型需要实现gorm的Scanner和Valuer接口
在base.go中实现
*/
type Goods struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;not null;comment:'分类外键id'"`
	Category   Category
	BrandID    int32 `gorm:"type:int;not null;comment:'品牌外键id'"`
	Brand      Brands

	OnSale   bool `gorm:"type:tinyint;not null;default:false;comment:'是否上架，false表示未上架，true表示上架'"`
	ShipFree bool `gorm:"type:tinyint;not null;default:false;comment:'是否包邮，false表示不包邮，true表示包邮'"`
	IsNew    bool `gorm:"type:tinyint;not null;default:false;comment:'是否新品，false表示不是新品，true表示新品'"`
	IsHot    bool `gorm:"type:tinyint;not null;default:false;comment:'是否热卖，false表示不是热卖，true表示热卖'"`

	Name    string `gorm:"type:varchar(100);not null;comment:'商品名称'"`
	GoodsSn string `gorm:"type:varchar(50);not null;comment:'商品编号'"`

	ClickNum    int32   `gorm:"type:int;not null;default:0;comment:'点击数'"`
	SoldNum     int32   `gorm:"type:int;not null;default:0;comment:'商品销售量'"`
	FavNum      int32   `gorm:"type:int;not null;default:0;comment:'收藏数'"`
	MarketPrice float32 `gorm:"type:decimal(10,2);not null;default:0.00;comment:'市场价'"`
	ShopPrice   float32 `gorm:"type:decimal(10,2);not null;default:0.00;comment:'本店价'"`

	GoodsBrief      string   `gorm:"type:varchar(200);not null;comment:'商品简介'"`
	Images          GormList `gorm:"type:varchar(5000);not null;comment:'商品轮播图'"`
	DescImages      GormList `gorm:"type:varchar(10000);not null;comment:'商品描述图片'"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null;comment:'商品封面图'"`
}

func (g *Goods) AfterCreate(tx *gorm.DB) (err error) {
	//创建商品的时候，需要创建商品的sku
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandID:     g.BrandID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		GoodsBrief:  g.GoodsBrief,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		ShopPrice:   g.ShopPrice,
	}
	_, err = global.EsClient.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(string(esModel.ID)).Do(context.Background())
	return
}

func (g *Goods) AfterUpdate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandID:     g.BrandID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	_, err = global.EsClient.Update().Index(esModel.GetIndexName()).
		Doc(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (g *Goods) AfterDelete(tx *gorm.DB) (err error) {
	_, err = global.EsClient.Delete().Index(EsGoods{}.GetIndexName()).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

package main

import (
	"context"
	"github.com/olivere/elastic/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"shop_srvs/goods_srv/global"
	"shop_srvs/goods_srv/model"
	"strconv"
	"time"
)

func main() {
	MySQL2ES()
}
func MySQL() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:root@tcp(192.168.120.172:3306)/shop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"
	//设置全局logger，打印每个sql语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //取消USERS这种命名方式
		},
	})
	if err != nil {
		panic(err)
	}

	//创建表

	//迁移
	_ = db.AutoMigrate(
		&model.Category{},
		&model.Brands{},
		&model.Banner{},
		&model.GoodsCategoryBrand{},
	)
	_ = db.AutoMigrate(&model.Goods{})
}

func MySQL2ES() {
	dsn := "root:root@tcp(192.168.120.172:3306)/shop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"
	//设置全局logger，打印每个sql语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //取消USERS这种命名方式
		},
	})
	if err != nil {
		panic(err)
	}
	host := "http://192.168.120.172:9200"
	logger := log.New(log.Writer(), "GoShop", log.LstdFlags)

	global.EsClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false), elastic.SetInfoLog(logger))
	if err != nil {
		panic(err)
	}
	//从mysql拿去所有的数据，并创建到es中
	var goods []model.Goods
	db.Find(&goods)
	for _, g := range goods {
		esModel := model.EsGoods{
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
		_, err = global.EsClient.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

}

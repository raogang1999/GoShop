package initialize

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"shop_srvs/goods_srv/global"
	"shop_srvs/goods_srv/model"
)

func InitEs() {
	//初始化连接
	host := fmt.Sprintf("http://%s:%d", global.SeverConfig.EsInfo.Host, global.SeverConfig.EsInfo.Port)
	logger := log.New(log.Writer(), "GoShop", log.LstdFlags)
	var err error
	global.EsClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false), elastic.SetInfoLog(logger))
	if err != nil {
		panic(err)
	}

	//新建mapping和index
	//先查询是否存在index
	exists, err := global.EsClient.IndexExists(model.EsGoods{}.GetIndexName()).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exists {
		//创建索引
		createIndex, err := global.EsClient.CreateIndex(model.EsGoods{}.GetIndexName()).BodyString(model.EsGoods{}.GetMapping()).Do(context.Background())
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			panic("创建索引失败")
		}
	}
}

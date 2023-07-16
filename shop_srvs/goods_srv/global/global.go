package global

import (
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
	"shop_srvs/goods_srv/config"
)

var (
	DB          *gorm.DB
	EsClient    *elastic.Client
	SeverConfig *config.ServerConfig = &config.ServerConfig{}
	NacosConfig *config.NacosConfig  = &config.NacosConfig{}
)

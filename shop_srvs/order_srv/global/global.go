package global

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"shop_srvs/order_srv/config"
	"shop_srvs/order_srv/proto"
)

var (
	DB          *gorm.DB
	RedisDB     *redis.Client
	RedisLocker *redsync.Redsync
	SeverConfig *config.ServerConfig = &config.ServerConfig{}
	NacosConfig *config.NacosConfig  = &config.NacosConfig{}

	// 微服务全局
	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)

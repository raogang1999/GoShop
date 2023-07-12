package global

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"shop_srvs/inventory_srv/config"
)

var (
	DB          *gorm.DB
	RedisDB     *redis.Client
	RedisLocker *redsync.Redsync
	SeverConfig *config.ServerConfig = &config.ServerConfig{}
	NacosConfig *config.NacosConfig  = &config.NacosConfig{}
)

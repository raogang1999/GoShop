package global

import (
	"gorm.io/gorm"
	"shop_srvs/user_srv/config"
)

var (
	DB          *gorm.DB
	SeverConfig *config.ServerConfig = &config.ServerConfig{}
	NacosConfig *config.NacosConfig  = &config.NacosConfig{}
)

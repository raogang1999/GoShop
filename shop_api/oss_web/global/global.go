package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop_api/oss_web/config"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	Trans        ut.Translator
)

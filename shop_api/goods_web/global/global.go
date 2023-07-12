package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop_api/goods_web/config"
	"shop_api/goods_web/proto"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	Trans        ut.Translator

	GoodsSrvClient proto.GoodsClient
)

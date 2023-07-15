package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop_api/userop_web/config"
	"shop_api/userop_web/proto"
)

var (
	ServerConfig     *config.ServerConfig = &config.ServerConfig{}
	NacosConfig      *config.NacosConfig  = &config.NacosConfig{}
	Trans            ut.Translator
	MessageSrvClient proto.MessageClient
	AddressSrvClient proto.AddressClient
	UserFavSrvClient proto.UserFavClient
	GoodsSrvClient   proto.GoodsClient
)

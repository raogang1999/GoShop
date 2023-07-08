package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop_api/user_web/config"
	"shop_api/user_web/proto"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}
	Trans        ut.Translator

	UserSrvClient proto.UserClient
)

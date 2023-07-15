package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop_api/order_web/config"
	"shop_api/order_web/proto"
)

var (
	ServerConfig       *config.ServerConfig = &config.ServerConfig{}
	NacosConfig        *config.NacosConfig  = &config.NacosConfig{}
	Trans              ut.Translator
	OrderSrvClient     proto.OrderClient
	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)

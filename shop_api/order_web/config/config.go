package config

type ServerConfig struct {
	Name             string       `mapstructure:"name" json:"name"`
	Host             string       `mapstructure:"host" json:"host"`
	Port             int          `mapstructure:"port" json:"port"`
	Tags             []string     `mapstructure:"tags" json:"tags"`
	GoodsSrvInfo     SrvConfig    `mapstructure:"goods_srv" json:"goods_srv"`
	OrderSrvInfo     SrvConfig    `mapstructure:"order_srv" json:"order_srv"`
	InventorySrvInfo SrvConfig    `mapstructure:"inventory_srv" json:"inventory_srv"`
	JWTInfo          JWTConfig    `mapstructure:"jwt" json:"jwt"`
	AliSmsConfig     AliSmsConfig `mapstructure:"sms" json:"sms"`
	ConsulInfo       ConsulConfig `mapstructure:"consul" json:"consul"`
	AliPayInfo       AliPayConfig `mapstructure:"alipay" json:"alipay"`
}

// 阿里云支付配置
type AliPayConfig struct {
	AppId        string `mapstructure:"app_id" json:"app_id"`
	PrivateKey   string `mapstructure:"private_key" json:"private_key"`
	AliPayPubKey string `mapstructure:"ali_pub_key" json:"ali_pub_key"`
	NotifyUrl    string `mapstructure:"notify_url" json:"notify_url"`
	ReturnUrl    string `mapstructure:"return_url" json:"return_url"`
}
type SrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type AliSmsConfig struct {
	AccessKeyId  string `mapstructure:"key" json:"key"`
	AccessSecret string `mapstructure:"secret" json:"secret" `
	Expire       int    `mapstructure:"expire" json:"expire"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	Group     string `mapstructure:"group"`
	DataId    string `mapstructure:"data_id"`
	Password  string `mapstructure:"password"`
}

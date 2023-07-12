package config

//MySQL配置

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"` //数据库名
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}
type ServerConfig struct {
	Name             string             `mapstructure:"name" json:"name"` //给consul用
	Tags             []string           `mapstructure:"tags" json:"tags"`
	Host             string             `mapstructure:"host" json:"host"` //服务启动的host
	MysqlInfo        MysqlConfig        `mapstructure:"mysql" json:"mysql"`
	ConsulInfo       ConsulConfig       `mapstructure:"consul" json:"consul"`
	RedisInfo        RedisConfig        `mapstructure:"redis" json:"redis"`
	GoodsSrvInfo     GoodsSrvConfig     `mapstructure:"goods_srv" json:"goods_srv"`
	InventorySrvInfo InventorySrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
}

// 商品微服务的配置
type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

// 库存微服务的配置
type InventorySrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	Group     string `mapstructure:"group"`
	DataId    string `mapstructure:"data_id"`
	Password  string `mapstructure:"password"`
}

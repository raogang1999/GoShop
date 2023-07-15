package config

type ServerConfig struct {
	Name          string       `mapstructure:"name" json:"name"`
	Host          string       `mapstructure:"host" json:"host"`
	Port          int          `mapstructure:"port" json:"port"`
	Tags          []string     `mapstructure:"tags" json:"tags"`
	GoodsSrvInfo  SrvConfig    `mapstructure:"goods_srv" json:"goods_srv"`
	UserOpSrvInfo SrvConfig    `mapstructure:"userop_srv" json:"userop_srv"`
	JWTInfo       JWTConfig    `mapstructure:"jwt" json:"jwt"`
	ConsulInfo    ConsulConfig `mapstructure:"consul" json:"consul"`
}

type SrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
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

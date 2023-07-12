package config

type ServerConfig struct {
	Name         string        `mapstructure:"name" json:"name"`
	Host         string        `mapstructure:"host" json:"host"`
	Port         int           `mapstructure:"port" json:"port"`
	Tags         []string      `mapstructure:"tags" json:"tags"`
	UserSrvInfo  UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	JWTInfo      JWTConfig     `mapstructure:"jwt" json:"jwt"`
	AliSmsConfig AliSmsConfig  `mapstructure:"sms" json:"sms"`
	RedisConfig  RedisConfig   `mapstructure:"redis" json:"redis"`
	ConsulInfo   ConsulConfig  `mapstructure:"consul" json:"consul"`
}

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
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
type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
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

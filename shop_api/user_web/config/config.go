package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type ServerConfig struct {
	Name         string        `mapstructure:"name"`
	Port         int           `mapstructure:"port"`
	UserSrvInfo  UserSrvConfig `mapstructure:"user_srv"`
	JWTInfo      JWTConfig     `mapstructure:"jwt"`
	AliSmsConfig AliSmsConfig  `mapstructure:"sms"`
	RedisConfig  RedisConfig   `mapstructure:"redis"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type AliSmsConfig struct {
	AccessKeyId  string `mapstructure:"key"`
	AccessSecret string `mapstructure:"secret"`
	Expire       int    `mapstructure:"expire"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

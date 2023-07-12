package config

type ServerConfig struct {
	Name        string         `mapstructure:"name" json:"name"`
	Host        string         `mapstructure:"host" json:"host"`
	Port        int            `mapstructure:"port" json:"port"`
	Tags        []string       `mapstructure:"tags" json:"tags"`
	UserSrvInfo GoodsSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	JWTInfo     JWTConfig      `mapstructure:"jwt" json:"jwt"`
	ConsulInfo  ConsulConfig   `mapstructure:"consul" json:"consul"`
	OssInfo     OssConfig      `mapstructure:"oss" json:"oss"`
}

type GoodsSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
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

type OssConfig struct {
	Host        string `mapstructure:"host" json:"host"`
	ApiKey      string `mapstructure:"key" json:"key"`
	ApiSecret   string `mapstructure:"secret" json:"secret"`
	CallBackUrl string `mapstructure:"callback_url" json:"callback_url"`
	UploadDir   string `mapstructure:"upload_dir" json:"upload_dir"`
}

package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"shop_api/user_web/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	debug := GetEnvInfo("SHOP_DEBUG")
	configPrefix := "config"
	configFileName := fmt.Sprintf("user_web/%s_product.yaml", configPrefix)
	if debug {
		configFileName = fmt.Sprintf("user_web/%s_debug.yaml", configPrefix)
	}
	v := viper.New()
	v.SetConfigFile(configFileName)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息%v", global.ServerConfig)
	//动态监听变化
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Infof("配置文件修改了...", in.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&global.ServerConfig)
	})

}

package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"time"
)

type ServerConfig struct {
	Name        string `mapstructure:"name"`
	Port        int    `mapstructure:"port"`
	MysqlConfig `mapstructure:"mysql"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func main() {

	debug := GetEnvInfo("SHOP_DEBUG")
	configFileName := "config_product.yaml"
	if debug {
		configFileName = "config_debug.yaml"
	}
	v := viper.New()
	v.SetConfigFile(configFileName)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	var serverConfig ServerConfig
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	fmt.Println(serverConfig)

	//动态监听变化
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...", in.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&serverConfig)
		fmt.Println(serverConfig)
	})
	time.Sleep(10 * time.Second)
}

package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
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

	if err := v.Unmarshal(global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息%v", global.NacosConfig)
	//动态监听变化
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Infof("配置文件修改了...", in.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.NacosConfig)
	})

	//从nacos中读取信息
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}
	// 创建动态配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})

	err = json.Unmarshal([]byte(content), global.ServerConfig)
	if err != nil {
		zap.S().Errorw("json unmarshal failed", "msg", err.Error())
	}
	zap.S().Infof("服务器配置信息%v", global.ServerConfig)
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
		OnChange: func(namespace, group, dataId, data string) {
			zap.S().Info("配置文件修改了...")
			err = json.Unmarshal([]byte(content), global.ServerConfig)
			if err != nil {
				zap.S().Errorw("json unmarshal failed", "msg", err.Error())
			}
			zap.S().Infof("服务器配置信息%v", global.ServerConfig)
		},
	})
	if err != nil {
		panic(err)
	}

}

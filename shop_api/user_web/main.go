package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"shop_api/user_web/global"
	"shop_api/user_web/initialize"

	shop_validator "shop_api/user_web/validator"
)

func main() {

	// 1. 初始化日志
	initialize.InitLogger()
	// 2. 初始化配置文件
	initialize.InitConfig()
	// 3. 初始化路由
	Router := initialize.Routers()
	//4. 初始化验证器
	err := initialize.InitTrans("zh")

	if err != nil {
		panic(err)
		//zap.S().Errorw("初始化验证器失败", "msg", err.Error())
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", shop_validator.ValidateMobile)
		//翻译
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法手机号码！", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	zap.S().Debugf("启动端口: %d", global.ServerConfig.Port)

	// 启动服务
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); nil != err {
		zap.S().Panic("启动失败: ", err.Error())
	}

}

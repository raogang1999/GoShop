package api

import (
	"context"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"net/http"
	"shop_api/user_web/forms"
	"shop_api/user_web/global"
	"strings"
	"time"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func GenerateSmsCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}
func SendSms(ctx *gin.Context) {
	sendSmsForm := forms.SendSmSForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	client, _err := CreateClient(tea.String(global.ServerConfig.AliSmsConfig.AccessKeyId), tea.String(global.ServerConfig.AliSmsConfig.AccessSecret))
	if _err != nil {
		return
	}

	smsCode := GenerateSmsCode(6)
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{}
	sendSmsRequest.PhoneNumbers = tea.String(sendSmsForm.Mobile)
	sendSmsRequest.SignName = tea.String("智能火灾监管系统")
	sendSmsRequest.TemplateCode = tea.String("SMS_188990505")
	sendSmsRequest.TemplateParam = tea.String("{\"code\":" + smsCode + "}")

	runtime := &util.RuntimeOptions{}
	options, err := client.SendSmsWithOptions(sendSmsRequest, runtime)
	fmt.Println(options.Body.Code)
	if err != nil {
		return
	}
	//注册码还需要被校验
	//保存验证码 手机号+验证码 map[string] string
	//redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisConfig.Host, global.ServerConfig.RedisConfig.Port),
		Password: "", // 没有密码，默认值
		DB:       0,  // 默认DB 0
	})

	rdb.Set(context.Background(), sendSmsForm.Mobile, smsCode, time.Duration(global.ServerConfig.AliSmsConfig.Expire)*time.Second)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}

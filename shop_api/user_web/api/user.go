package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"shop_api/user_web/forms"
	"shop_api/user_web/global"
	"shop_api/user_web/global/response"
	"shop_api/user_web/middlewares"
	"shop_api/user_web/models"
	"shop_api/user_web/proto"
	"strconv"
	"strings"
	"time"
)

func HandlerGrpcErrorToHttp(err error, c *gin.Context) {
	//将GRP的状态码转换成HTTP的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "服务不可用",
				})

			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误" + e.Message(),
				})
			}

		}
	}
}

func GetUserList(ctx *gin.Context) {

	claims, _ := ctx.Get("claims")
	curUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户：%d", curUser.ID)

	//调用接口

	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("pSize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 [用户列表失败]")
		HandlerGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)

	for _, value := range rsp.Data {
		user := response.UserResponse{
			Id:       value.Id,
			Nickname: value.NickName,
			Birthday: response.JsonTime(time.Unix(int64(value.Birthday), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}

		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)

	zap.S().Debug("获取用户列表页")

}

func removeTopStruct(fileds map[string]string) map[string]string {
	res := map[string]string{}
	for filed, err := range fileds {
		res[filed[strings.IndexAny(filed, ".")+1:]] = err
	}
	return res
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"msg": removeTopStruct(errs.Translate(global.Trans)),
	})
}

// 登录
func PasswordLogin(c *gin.Context) {
	passwordLoginForm := forms.PasswordLoginForm{}
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	//验证码
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	//登录逻辑
	if rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.UserMobileRequest{Mobile: passwordLoginForm.Mobile}); err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				c.JSON(http.StatusBadRequest, gin.H{
					"mobile": "用户不存在",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"mobile": "登录失败",
				})
			}
			return
		}

	} else {
		//校验密码
		if pasRsp, pasErr := global.UserSrvClient.CheckPassword(context.Background(), &proto.CheckPasswordInfo{
			Password:          passwordLoginForm.Password,
			EncryptedPassword: rsp.Password,
		}); pasErr != nil {
			//网络问题
			c.JSON(http.StatusInternalServerError, gin.H{
				"mobile": "登录失败",
			})
		} else {
			if pasRsp.Success {
				//生成token
				new_jwt := middlewares.NewJWT()

				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					Nickname:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               //生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //过期时间
						Issuer:    "GoShop",
					},
				}
				token, err := new_jwt.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}
				//返回结果
				c.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nik_name":   rsp.NickName,
					"token":      token,
					"expired_at": claims.StandardClaims.ExpiresAt * 1000,
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "密码错误",
				})
			}

		}
	}
}

// 注册
func Register(ctx *gin.Context) {
	//表单
	registerForm := forms.RegisterForm{}
	if err := ctx.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	//短信验证码校验
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisConfig.Host, global.ServerConfig.RedisConfig.Port),
		Password: "", // 没有密码，默认值
		DB:       0,  // 默认DB 0
	})
	value, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
		return
	}
	//验证码验证
	if value != registerForm.Code {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}

	//注册逻辑
	if rsp, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		Mobile:   registerForm.Mobile,
		Password: registerForm.Password,
	}); err != nil {
		HandlerGrpcErrorToHttp(err, ctx)
		zap.S().Errorw("[Register] 注册 [用户服务失败]", "msg", err.Error())
	} else {
		//注册成功
		//生成token，返回x-token
		new_jwt := middlewares.NewJWT()
		claims := models.CustomClaims{
			ID:          uint(rsp.Id),
			Nickname:    rsp.NickName,
			AuthorityId: uint(rsp.Role),
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix(),               //生效时间
				ExpiresAt: time.Now().Unix() + 60*60*24*30, //过期时间
				Issuer:    "GoShop",
			},
		}
		token, err := new_jwt.CreateToken(claims)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "生成token失败",
			})
			return
		}
		//返回结果
		ctx.JSON(http.StatusOK, gin.H{
			"id":         rsp.Id,
			"nik_name":   rsp.NickName,
			"token":      token,
			"expired_at": claims.StandardClaims.ExpiresAt * 1000,
		})
	}
}

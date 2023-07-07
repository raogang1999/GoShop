package forms

type SendSmSForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号格式
	//动态验证码类型，登录或者注册
	Type string `form:"type" json:"type" binding:"required,oneof=login register"`
}

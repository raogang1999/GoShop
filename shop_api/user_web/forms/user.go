package forms

type PasswordLoginForm struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号格式
	Password  string `form:"password" json:"password" binding:"required,min=6,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,len=5"`
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

type RegisterForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号格式
	Password string `form:"password" json:"password" binding:"required,min=6,max=20"`
	Code     string `form:"code" json:"code" binding:"required,len=6"`
}

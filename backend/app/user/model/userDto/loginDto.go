package userDto

type CaptchaVO struct {
	Id  string `json:"uuid"`
	Img string `json:"img"`
}

type LoginBody struct {
	Username string `json:"username" binding:"required"` //用户名
	Password string `json:"password" binding:"required"` //密码
	Code     string `json:"code" binding:"required"`     //验证码
	Uuid     string `json:"uuid" binding:"required"`     //uuid
}

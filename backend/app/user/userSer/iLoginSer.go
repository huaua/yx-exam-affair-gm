package userSer

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"
)

type ILoginService interface {
	Login(user *userDo.SysUser) *string

	GenerateCode() *userDto.CaptchaVO

	VerifyCaptcha(id, base64 string) bool

	Logout(token string)
}

package userSerImpl

import (
	"image/color"
	"strconv"

	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/mojocn/base64Captcha"
	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/cst"
	"github.com/sealsee/web-base/public/jwt"
	"github.com/sealsee/web-base/public/utils/redis"
	"github.com/sealsee/web-base/public/utils/token"
)

type LoginService struct {
}

var loginService *LoginService

func init() {
	loginService = &LoginService{}
}

func NewLoginService() *LoginService {
	return loginService
}

// 生成driver，高，宽，背景文字的干扰，画线条数，背景颜色的指针，字体
var driver = base64Captcha.NewDriverMath(38, 106, 0, 0, &color.RGBA{0, 0, 0, 0}, nil, []string{"wqy-microhei.ttc"})
var store = base64Captcha.DefaultMemStore

func (loginService *LoginService) GenerateCode() *userDto.CaptchaVO {
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64, err := captcha.Generate()
	if err != nil {
		panic(err)
	}
	return &userDto.CaptchaVO{Id: id, Img: b64}
}

func (loginService *LoginService) Login(user *userDo.SysUser) *string {

	sessionUser := context.NewSessionUser()
	// 将相关用户信息放入session
	sessionUser.UserId = user.UserId
	sessionUser.UserName = user.UserName
	sessionUser.Avatar = user.Avatar
	sessionUser.NickName = user.NickName
	sessionUser.DeptId = user.DeptId

	sessionUser.Token = strconv.FormatInt(user.UserId, 10)
	tokenStr := jwt.GenToken(sessionUser.Token)
	token.RefreshToken(sessionUser)
	return &tokenStr
}

func (loginService *LoginService) VerifyCaptcha(id, base64 string) bool {
	return store.Verify(id, base64, true)
}

func (loginService *LoginService) Logout(token string) {
	redis.Del(cst.LoginTokenKey + token)
}

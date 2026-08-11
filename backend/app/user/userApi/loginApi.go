package userApi

import (
	"yx-exam-affair-gm/app/user/model/userDto"
	"yx-exam-affair-gm/app/user/userSer"
	"yx-exam-affair-gm/app/user/userSer/userSerImpl"

	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/jwt"
	"github.com/sealsee/web-base/public/route"
	"github.com/sealsee/web-base/public/setting"
	"github.com/sealsee/web-base/public/utils/crypto"
	"github.com/sealsee/web-base/public/utils/token"
	"github.com/sealsee/web-base/public/web"
)

type LoginController struct {
}

var loginService userSer.ILoginService

func init() {
	route.UseDefaultGroup().
		AddPOSTHandel("/login", login, false, "登录").
		AddGETHandel("/captcha", getCaptcha, false, "获取验证码").
		AddPOSTHandel("/logout", logout, false, "退出登录").
		AddGETHandel("/get_info", getInfo, true, "获取登录用户信息").
		AddGETHandel("/get_routers", getRouters, true, "获取路由信息")

	loginService = userSerImpl.NewLoginService()
}

func getCaptcha(c *gin.Context) {
	web.NewJsonResult(c).AddData("captcha", loginService.GenerateCode()).Render()
}

func login(c *gin.Context) {
	json := web.NewJsonResult(c)

	var loginBody userDto.LoginBody
	if err := c.ShouldBindJSON(&loginBody); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	// 生产环境必须验证码校验；非生产环境，可用"8888"模拟通过
	if setting.Conf.Mode == "prod" || (setting.Conf.Mode != "prod" && loginBody.Code != "8888") {
		captcha := loginService.VerifyCaptcha(loginBody.Uuid, loginBody.Code)
		if !captcha {
			json.SetErrs(errs.CAPTCHA_ERR).Render()
			return
		}
	}
	user := sysUserService.GetSysUserByName(loginBody.Username)
	if user == nil {
		json.SetErrs(errs.USER_PASSWORD_ERR).Render()
		return
	} else if user.State == common.Disable {
		json.SetErrsWithArgs(errs.USER_DISABLED_ERR, loginBody.Username).Render()
		return
	} else if !crypto.CheckPwd(loginBody.Password, user.Password) {
		json.SetErrs(errs.USER_PASSWORD_ERR).Render()
		return
	}

	json.AddData("token", loginService.Login(user)).Render()
}

func logout(c *gin.Context) {
	json := web.NewJsonResult(c)

	authHeader := c.Request.Header.Get("Authorization")
	if len(authHeader) == 0 {
		json.Render()
		return
	}
	loginUser, _ := jwt.ParseToken(authHeader)
	if loginUser != nil {
		loginService.Logout(loginUser.Token)
	}
	web.NewJsonResult(c).Render()
}

// 登录后，首先调用此方法获取用户信息
func getInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	sessionUser := ctx.GetUser()

	// 获取用户角色
	roles := sysRoleService.GetRolePermission(sessionUser.UserId)

	// 刷新sessionUser里的Roles
	sessionUser.Roles = roles
	// 本系统未用到权限控制，暂全局设为-1
	sessionUser.Permissions = []string{"-1"}
	go token.RefreshToken(sessionUser)

	json.AddData("user", sessionUser).
		AddData("roles", roles).
		Render()
}

// 注意：调用此方法前提，必须调用过getInfo
func getRouters(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	sessionUser := ctx.GetUser()
	// 获取用户菜单列表和权限
	menus, permissions := sysMenuService.ListMenuTreeByUserId(sessionUser.UserId)

	go func() {
		// 刷新sessionUser里的Roles
		sessionUser.Permissions = permissions
		// 如果Roles为空，则可能是并发导致这里查询不是最新sessionUser
		if len(sessionUser.Roles) == 0 {
			roles := sysRoleService.GetRolePermission(sessionUser.UserId)
			sessionUser.Roles = roles
		}
		token.RefreshToken(sessionUser)
	}()

	routers := sysMenuService.BuildRouters(menus)
	json.AddData("routers", routers).
		AddData("permissions", permissions).
		Render()
}

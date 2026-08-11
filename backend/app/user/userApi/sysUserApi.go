package userApi

import (
	"regexp"
	"strings"
	"unicode/utf8"
	"yx-exam-affair-gm/app/cst"
	"yx-exam-affair-gm/app/hint"
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"
	"yx-exam-affair-gm/app/user/userSer"
	"yx-exam-affair-gm/app/user/userSer/userSerImpl"

	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/route"
	"github.com/sealsee/web-base/public/utils/crypto"
	"github.com/sealsee/web-base/public/utils/token"
	"github.com/sealsee/web-base/public/web"
)

var sysUserService userSer.ISysUserService

func init() {
	route.AddGroup("/api/auth/sys_user", "用户").
		AddGETHandel("/get", getUser, true, "查询用户信息").
		AddPOSTHandel("/list", listUser, true, "查询用户列表").
		AddPOSTHandel("/add", addUser, true, "新增用户").
		AddPOSTHandel("/edit", editUser, true, "修改用户").
		AddPOSTHandel("/reset_pwd", resetPwd, true, "重置密码").
		AddPOSTHandel("/change_state", changeUserState, true, "修改状态").
		AddPOSTHandel("/del", delUser, true, "删除用户").
		AddPOSTHandel("/export", exportUser, true, "导出用户").
		AddGETHandel("/get_profile", getProfile, true, "个人信息").
		AddPOSTHandel("/update_profile", updateProfile, true, "更新个人信息").
		AddPOSTHandel("/change_pwd", changePwd, true, "更换个人密码").
		AddPOSTHandel("/avatar", avatar, true, "更新个人头像")
	sysUserService = userSerImpl.NewSysUserService()
}

func getUser(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	targetUserId := ctx.QueryInt64("userId")
	// 查询当前用户所拥有的所有角色列表
	roles := sysRoleService.GetSysRoleByUserId(ctx.GetUserId())
	json.AddData("roles", roles)
	if targetUserId < 1 {
		json.Render()
		return
	}
	// 查询目标用户信息
	user := sysUserService.GetSysUserById(targetUserId)
	user.Password = "" // 密码不需要返回
	// 查询目标用户的角色ids
	roleIds := sysRoleService.GetRoleIdsByUserId(targetUserId)
	json.AddData("user", user).
		AddData("roleIds", roleIds).Render()
}

func listUser(c *gin.Context) {
	json := web.NewJsonResult(c)
	userQuery := new(userDto.SysUserQuery)
	if err := c.ShouldBindJSON(userQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	page := userQuery.GetPage()
	list := sysUserService.ListSysUser(userQuery, page)
	json.SetPageList(list, page).Render()
}

func addUser(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	userDO := new(userDo.SysUser)
	if err := c.ShouldBindJSON(userDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	// 账号长度校验
	userName := strings.TrimSpace(userDO.UserName)
	userNameLen := utf8.RuneCountInString(userName)
	if userName == "" || userNameLen < 3 || userNameLen > 18 {
		json.SetErrs(hint.SYS_USERNAME_FORMAT_ERR).Render()
		return
	}

	// 校验密码
	if userDO.Password == "" || !isValidPassword(userDO.Password) {
		json.SetErrs(hint.SYS_PASSWORD_FORMAT_ERR).Render()
		return
	}

	// 姓名长度校验
	nickName := strings.TrimSpace(userDO.NickName)
	nickNameLen := utf8.RuneCountInString(nickName)
	if nickName == "" || nickNameLen < 2 || nickNameLen > 15 {
		json.SetErrs(hint.SYS_NICKNAME_FORMAT_ERR).Render()
		return
	}

	// 校验角色、学院
	if userDO.RoleKey == "" {
		json.SetErrs(hint.SYS_USER_NO_ROLEKEY).Render()
		return
	}
	if userDO.RoleKey == cst.SYS_ROLE_KEY_DEPARTMENT && userDO.DeptId < 1 {
		json.SetErrs(hint.SYS_USER_NO_DEPARTMENT).Render()
		return
	}
	// 招办角色的部门id固定为-100
	if userDO.RoleKey == cst.SYS_ROLE_KEY_ADMISSIONS {
		userDO.DeptId = cst.BIZ_DEPTID_FOR_ADMISSIONS
	}
	// 角色id转换
	userDO.RoleIds = append(userDO.RoleIds, cst.RoleMap[userDO.RoleKey])

	userDO.State = common.OK
	userDO.UserType = cst.USER_TYPE_SYS
	userDO.SetCreateBy(ctx.GetUserId())
	userDO.SetUpdateBy(ctx.GetUserId())
	success := sysUserService.InsertSysUser(userDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func isValidPassword(password string) bool {
	// 校验密码是否由6~18位字符组成 并且同时包含大小写字母、数字
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	isCorrectLength := utf8.RuneCountInString(password) >= 6 && utf8.RuneCountInString(password) <= 18

	return hasLower && hasUpper && hasDigit && isCorrectLength
}

func editUser(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	sessionUser := ctx.GetUser()

	userDO := new(userDo.SysUser)
	if err := c.ShouldBindJSON(userDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if userDO.UserId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	// 姓名长度校验
	nickName := strings.TrimSpace(userDO.NickName)
	nickNameLen := utf8.RuneCountInString(nickName)
	if nickName == "" || nickNameLen < 2 || nickNameLen > 15 {
		json.SetErrs(hint.SYS_NICKNAME_FORMAT_ERR).Render()
		return
	}

	userEditDO := new(userDo.SysUser)
	userEditDO.UserId = userDO.UserId
	userEditDO.NickName = userDO.NickName

	// 是学院修改
	if userDO.DeptId > 0 {
		// 如果角色变了，--不允许！
		userDB := sysUserService.GetSysUserById(userDO.UserId)
		if userDB.DeptId == cst.BIZ_DEPTID_FOR_ADMISSIONS {
			json.SetErrs(errs.PARAM_INVALID_ERR).Render()
			return
		}
		userEditDO.DeptId = userDO.DeptId
	}

	userEditDO.SetUpdateBy(sessionUser.UserId)
	success := sysUserService.UpdateUserSingle(userEditDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func resetPwd(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	userDO := new(userDo.SysUser)
	if err := c.ShouldBindJSON(userDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if userDO.UserId < 1 || userDO.Password == "" {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if !isValidPassword(userDO.Password) {
		json.SetErrs(hint.SYS_PASSWORD_FORMAT_ERR).Render()
		return
	}
	// 比较用户的新老密码是否一致
	userDB := sysUserService.GetSysUserById(userDO.UserId)
	if crypto.CheckPwd(userDO.Password, userDB.Password) {
		json.SetErrs(hint.SYS_SAME_PASSWORD_ERR).Render()
		return
	}
	success := sysUserService.ResetPwd(userDO.UserId, userDO.Password, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func changeUserState(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	userDO := new(userDo.SysUser)
	if err := c.ShouldBindJSON(userDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if userDO.UserId < 1 || userDO.State == "" {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := sysUserService.UpdateUserState(userDO.UserId, userDO.State, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func delUser(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	userQuery := new(userDto.SysUserQuery)
	if err := c.ShouldBindJSON(userQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if userQuery.UserId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := sysUserService.DeleteSysUserById(userQuery.UserId, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func exportUser(c *gin.Context) {
	json := web.NewJsonResult(c)
	userQuery := new(userDto.SysUserQuery)
	if err := c.ShouldBindJSON(userQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data := sysUserService.ExportSysUser(userQuery)
	web.DataPackageExcel(c, data)
}

func getProfile(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	userId := ctx.GetUserId()
	// 查询当前用户信息
	user := sysUserService.GetSysUserById(userId)
	// 查询当前用户所拥有的所有角色列表
	roles := sysRoleService.GetSysRoleByUserId(userId)
	roleGroup := make([]string, 0)
	for _, r := range roles {
		roleGroup = append(roleGroup, r.RoleName)
	}
	json.AddData("user", user).
		AddData("roleGroup", strings.Join(roleGroup, ",")).Render()
}

func updateProfile(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	userDO := new(userDo.SysUser)
	if err := c.ShouldBindJSON(userDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if userDO.UserId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	userDO.SetUpdateBy(ctx.GetUserId())
	success := sysUserService.UpdateUserProfile(userDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func changePwd(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	user := new(userDto.SysUserQuery)
	if err := c.ShouldBindJSON(user); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if user.OldPassword == "" || user.NewPassword == "" {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if !isValidPassword(user.NewPassword) {
		json.SetErrs(hint.SYS_PASSWORD_FORMAT_ERR).Render()
		return
	}
	userId := ctx.GetUserId()
	// 查询当前用户信息
	userDB := sysUserService.GetSysUserById(userId)
	if !crypto.CheckPwd(user.OldPassword, userDB.Password) {
		json.SetErrs(hint.SYS_OLD_PASSWORD_ERR).Render()
		return
	}
	if crypto.CheckPwd(user.NewPassword, userDB.Password) {
		json.SetErrs(hint.SYS_SAME_PASSWORD_ERR).Render()
		return
	}
	// 更新密码
	success := sysUserService.ResetPwd(userId, user.NewPassword, userId)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func avatar(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	sessionUser := ctx.GetUser()
	userDO := new(userDo.SysUser)
	if err := c.ShouldBindJSON(userDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if userDO.Avatar == "" {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := sysUserService.UpdateUserAvatar(sessionUser.UserId, userDO.Avatar, sessionUser.UserId)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	// 更新session信息
	if sessionUser.Token != "" {
		sessionUser.Avatar = userDO.Avatar
		token.RefreshToken(sessionUser)
	}
	json.Render()
}

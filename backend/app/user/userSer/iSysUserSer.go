package userSer

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
)

type ISysUserService interface {
	GetSysUserById(id int64) *userDo.SysUser
	GetSysUserByName(username string) *userDo.SysUser
	ListSysUser(user *userDto.SysUserQuery, page *page.Page) []*userDo.SysUser
	InsertSysUser(user *userDo.SysUser) bool
	// 更新用户及角色
	UpdateSysUserAndRole(user *userDo.SysUser) bool
	// 更新用户信息-只更新user表
	UpdateUserSingle(user *userDo.SysUser) bool
	// 重置密码
	ResetPwd(userId int64, password string, oper int64) bool
	// 修改用户状态
	UpdateUserState(userId int64, state string, oper int64) bool
	// 更新用户个人信息
	UpdateUserProfile(user *userDo.SysUser) bool
	// 更新用户头像
	UpdateUserAvatar(userId int64, avatar string, oper int64) bool
	DeleteSysUserById(id int64, deleteBy int64) bool
	ExportSysUser(user *userDto.SysUserQuery) []byte
}

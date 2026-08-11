package userSer

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
)

type ISysRoleService interface {
	GetSysRoleById(id int64) *userDo.SysRole
	ListSysRole(role *userDto.SysRoleQuery, page *page.Page) []*userDo.SysRole
	InsertSysRole(role *userDo.SysRole) bool
	UpdateSysRoleById(role *userDo.SysRole) bool
	// 更新角色及对应菜单
	UpdateSysRoleAndMenu(role *userDo.SysRole) bool
	// 修改角色状态
	UpdateRoleState(roleId int64, state string, oper int64) bool
	DeleteSysRoleById(id int64, deleteBy int64) bool
	// 获取用户的角色ids
	GetRoleIdsByUserId(userId int64) []int64
	// 获取角色数据权限
	GetRolePermission(userId int64) []string
	// 获取用户拥有的角色列表
	GetSysRoleByUserId(userId int64) []*userDo.SysRole
}

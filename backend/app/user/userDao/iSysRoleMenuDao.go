package userDao

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

type ISysRoleMenuDao interface {
	Count(q *userDto.SysRoleMenuQuery) int
	List(q *userDto.SysRoleMenuQuery, page *page.Page) []*userDo.SysRoleMenu
	ListWithCondition(q *userDto.SysRoleMenuQuery, page *page.Page, condition string, args ...interface{}) []*userDo.SysRoleMenu
	// 查询菜单使用数量
	CheckMenuExistRole(menuId int64) int
	// 批量新增角色菜单信息
	BatchInsertRoleMenu(gtx tx.GTx, roleMenuList []*userDo.SysRoleMenu) int
	// 通过角色ID删除角色和菜单关联
	DeleteRoleMenuByRoleId(gtx tx.GTx, roleId int64) int
	// 批量删除角色菜单关联信息
	BatchDeleteRoleMenu(gtx tx.GTx, roleIds []int64) int
}

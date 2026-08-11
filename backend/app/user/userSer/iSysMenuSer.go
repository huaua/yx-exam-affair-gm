package userSer

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
)

type ISysMenuService interface {
	GetSysMenuById(id int64) *userDo.SysMenu
	ListSysMenu(menu *userDto.SysMenuQuery, page *page.Page) []*userDo.SysMenu
	// 查询菜单列表
	ListSysMenuByUserId(userId int64) []*userDo.SysMenu
	BuildMenuTreeSelect(menus []*userDo.SysMenu) []*userDto.TreeSelect
	// 获取菜单和权限
	ListMenuTreeByUserId(userId int64) (list []*userDo.SysMenu, permissions []string)
	BuildRouters(menus []*userDo.SysMenu) []*userDto.RouterVO
	// 根据角色ID查询菜单树信息
	ListMenuIdByRoleId(roleId int64, menuCheckStrictly bool) []int64
	InsertSysMenu(menu *userDo.SysMenu) bool
	UpdateSysMenuById(menu *userDo.SysMenu) bool
	DeleteSysMenuById(id int64, deleteBy int64) bool
	ExportSysMenu(menu *userDto.SysMenuQuery) []byte
}

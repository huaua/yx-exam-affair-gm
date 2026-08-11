package userDao

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

type ISysMenuDao interface {
	GetById(id int64) *userDo.SysMenu
	Count(menu *userDto.SysMenuQuery) int
	CountWithCondition(menu *userDto.SysMenuQuery, condition string, args ...interface{}) int
	List(menu *userDto.SysMenuQuery, page *page.Page) []*userDo.SysMenu
	ListWithCondition(q *userDto.SysMenuQuery, page *page.Page, condition string, args ...interface{}) []*userDo.SysMenu
	ListMenuTreeAll() []*userDo.SysMenu
	Insert(gtx tx.GTx, menu *userDo.SysMenu) bool
	BatchInsert(gtx tx.GTx, menuList []*userDo.SysMenu) bool
	Update(gtx tx.GTx, uData *userDo.SysMenu, uWhere *userDo.SysMenu) bool
	DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool
}

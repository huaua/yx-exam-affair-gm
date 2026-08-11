package userDao

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

type ISysRoleDao interface {
	GetById(id int64) *userDo.SysRole
	Count(role *userDto.SysRoleQuery) int
	CountWithCondition(role *userDto.SysRoleQuery, condition string, args ...interface{}) int
	ListAll() []*userDo.SysRole
	List(role *userDto.SysRoleQuery, page *page.Page) []*userDo.SysRole
	ListWithCondition(q *userDto.SysRoleQuery, page *page.Page, condition string, args ...interface{}) []*userDo.SysRole
	Insert(gtx tx.GTx, role *userDo.SysRole) bool
	BatchInsert(gtx tx.GTx, roleList []*userDo.SysRole) bool
	Update(gtx tx.GTx, uData *userDo.SysRole, uWhere *userDo.SysRole) bool
	DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool
}

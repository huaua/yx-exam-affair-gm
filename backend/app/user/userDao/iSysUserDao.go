package userDao

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

type ISysUserDao interface {
	GetById(id int64) *userDo.SysUser
	GetByName(username string) *userDo.SysUser
	Count(user *userDto.SysUserQuery) int
	CountWithCondition(userQuery *userDto.SysUserQuery, condition string, args ...interface{}) int
	List(user *userDto.SysUserQuery, page *page.Page) []*userDo.SysUser
	Insert(gtx tx.GTx, user *userDo.SysUser) bool
	BatchInsert(gtx tx.GTx, userList []*userDo.SysUser) bool
	Update(gtx tx.GTx, uData *userDo.SysUser, uWhere *userDo.SysUser) bool
	DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool
}

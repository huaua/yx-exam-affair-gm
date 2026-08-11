package userDao

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

type ISysDictTypeDao interface {
	GetById(id int64) *userDo.SysDictType
	Count(dictType *userDto.SysDictTypeQuery) int
	List(dictType *userDto.SysDictTypeQuery, page *page.Page) []*userDo.SysDictType
	Insert(gtx tx.GTx, dictType *userDo.SysDictType) bool
	BatchInsert(gtx tx.GTx, dictTypeList []*userDo.SysDictType) bool
	Update(gtx tx.GTx, uData *userDo.SysDictType, uWhere *userDo.SysDictType) bool
	DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool
}

package userDao

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

type ISysDictDataDao interface {
	GetById(id int64) *userDo.SysDictData
	GetByType(dictType string) []*userDo.SysDictData
	Count(dictData *userDto.SysDictDataQuery) int
	List(dictData *userDto.SysDictDataQuery, page *page.Page) []*userDo.SysDictData
	Insert(gtx tx.GTx, dictData *userDo.SysDictData) bool
	BatchInsert(gtx tx.GTx, dictDataList []*userDo.SysDictData) bool
	Update(gtx tx.GTx, uData *userDo.SysDictData, uWhere *userDo.SysDictData) bool
	DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool
	GetCountByLabelAndValueAndType(where *userDo.SysDictData) int
	UpdateNew(gtx tx.Db, uWhere *userDo.SysDictData, uData *userDo.SysDictData) int
}

package userDao

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

type ISysDataInoutDao interface {
	GetById(id int64) *userDo.SysDataInout
	Count(dataInout *userDto.SysDataInoutQuery) int
	List(dataInout *userDto.SysDataInoutQuery, page *page.Page) []*userDo.SysDataInout
	ListWithCondition(q *userDto.SysDataInoutQuery, page *page.Page, c string, args ...interface{}) []*userDo.SysDataInout
	Insert(gtx tx.GTx, dataInout *userDo.SysDataInout) bool
	BatchInsert(gtx tx.GTx, dataInoutList []*userDo.SysDataInout) bool
	Update(gtx tx.GTx, uData *userDo.SysDataInout, uWhere *userDo.SysDataInout) bool
	DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool
}

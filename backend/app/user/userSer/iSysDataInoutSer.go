package userSer

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
)

type ISysDataInoutService interface {
	GetSysDataInoutById(id int64) *userDo.SysDataInout
	ListSysDataInout(dataInout *userDto.SysDataInoutQuery, page *page.Page) []*userDo.SysDataInout
	ListSysDataInoutWithCondition(dataInoutQuery *userDto.SysDataInoutQuery, page *page.Page, condition string, args ...interface{}) []*userDo.SysDataInout
	InsertSysDataInout(dataInout *userDo.SysDataInout) bool
	UpdateSysDataInoutById(dataInout *userDo.SysDataInout) bool
	DeleteSysDataInoutById(id int64, deleteBy int64) bool
}

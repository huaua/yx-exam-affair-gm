package userSer

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
)

type ISysDictTypeService interface {
	GetSysDictTypeById(id int64) *userDo.SysDictType
	ListSysDictType(dictType *userDto.SysDictTypeQuery, page *page.Page) []*userDo.SysDictType
	InsertSysDictType(dictType *userDo.SysDictType) bool
	UpdateSysDictTypeById(dictType *userDo.SysDictType) bool
	DeleteSysDictTypeById(id int64, deleteBy int64) bool
	ExportSysDictType(dictType *userDto.SysDictTypeQuery) []byte
}

package userSer

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
)

type ISysDictDataService interface {
	GetSysDictDataById(id int64) *userDo.SysDictData
	GetSysDictDataByType(dictType string) []*userDo.SysDictData
	ListSysDictData(dictData *userDto.SysDictDataQuery, page *page.Page) []*userDo.SysDictData
	InsertSysDictData(dictData *userDo.SysDictData) bool
	UpdateSysDictDataById(dictData *userDo.SysDictData) bool
	DeleteSysDictDataById(id int64, deleteBy int64) bool
	ExportSysDictData(dictData *userDto.SysDictDataQuery) []byte
	UpdateDicData(dictData *userDo.SysDictData) bool
}

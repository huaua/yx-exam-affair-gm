package userDao

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

type ISysConfigDao interface {
	GetById(id int64) *userDo.SysConfig
	GetByKey(key string) *userDo.SysConfig
	Count(config *userDto.SysConfigQuery) int
	CountWithCondition(config *userDto.SysConfigQuery, condition string, args ...interface{}) int
	List(config *userDto.SysConfigQuery, page *page.Page) []*userDo.SysConfig
	ListWithCondition(q *userDto.SysConfigQuery, p *page.Page, condition string, args ...interface{}) []*userDo.SysConfig
	ListMap(config *userDto.SysConfigQuery, page *page.Page) []map[string]interface{}
	Insert(gtx tx.GTx, config *userDo.SysConfig) bool
	BatchInsert(gtx tx.GTx, configList []*userDo.SysConfig) bool
	Update(gtx tx.GTx, uData *userDo.SysConfig, uWhere *userDo.SysConfig) bool
	DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool
}

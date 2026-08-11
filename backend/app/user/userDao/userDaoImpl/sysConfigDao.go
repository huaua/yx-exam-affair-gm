package userDaoImpl

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/ds"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/query"
	"github.com/sealsee/web-base/public/ds/tx"
)

var sysConfigDaoImpl *SysConfigDaoImpl

func init() {
	sysConfigDaoImpl = &SysConfigDaoImpl{}
}

type SysConfigDaoImpl struct{}

func NewSysConfigDao() *SysConfigDaoImpl {
	return sysConfigDaoImpl
}

func (config *SysConfigDaoImpl) GetById(id int64) *userDo.SysConfig {
	if id < 1 {
		return nil
	}
	data := new(userDo.SysConfig)
	where := new(userDo.SysConfig)
	where.ConfigId = id
	where.Deleted = common.Normal
	res := ds.GetGDB().Find(data, where)
	if res.Error != nil {
		panic(res.Error)
	}
	if res.RowsAffected < 1 {
		return nil
	}
	return data
}

func (config *SysConfigDaoImpl) GetByKey(key string) *userDo.SysConfig {
	if key == "" {
		return nil
	}
	data := new(userDo.SysConfig)
	where := new(userDo.SysConfig)
	where.ConfigKey = key
	where.Deleted = common.Normal
	res := ds.GetGDB().Find(data, where)
	if res.Error != nil {
		panic(res.Error)
	}
	if res.RowsAffected < 1 {
		return nil
	}
	return data
}

func (config *SysConfigDaoImpl) Count(configQuery *userDto.SysConfigQuery) int {
	configQuery.Deleted = common.Normal
	return query.ExecGetQueryCount[userDto.SysConfigQuery, userDo.SysConfig](configQuery)
}

func (config *SysConfigDaoImpl) CountWithCondition(configQuery *userDto.SysConfigQuery, condition string, args ...interface{}) int {
	configQuery.Deleted = common.Normal
	return query.ExecGetQueryCountWithCondition[userDo.SysConfig](configQuery, condition, args...)
}

func (config *SysConfigDaoImpl) List(configQuery *userDto.SysConfigQuery, page *page.Page) []*userDo.SysConfig {
	configQuery.Deleted = common.Normal
	return query.ExecQueryList[userDto.SysConfigQuery, userDo.SysConfig](configQuery, page)
}

func (config *SysConfigDaoImpl) ListWithCondition(q *userDto.SysConfigQuery, p *page.Page, condition string, args ...interface{}) []*userDo.SysConfig {
	q.Deleted = common.Normal
	return query.ExecQueryListWithCondition[userDo.SysConfig](q, p, condition, args...)
}

func (config *SysConfigDaoImpl) ListMap(configQuery *userDto.SysConfigQuery, page *page.Page) []map[string]interface{} {
	configQuery.Deleted = common.Normal
	return query.ExecQueryListMapWithCondition[userDo.SysConfig](configQuery, page, "1=?", "1")
}

func (config *SysConfigDaoImpl) Insert(gtx tx.GTx, configDO *userDo.SysConfig) bool {
	configDO.Deleted = common.Normal
	res := gtx.Create(configDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (config *SysConfigDaoImpl) BatchInsert(gtx tx.GTx, configList []*userDo.SysConfig) bool {
	res := gtx.CreateInBatches(configList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (config *SysConfigDaoImpl) Update(gtx tx.GTx, uData *userDo.SysConfig, uWhere *userDo.SysConfig) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (config *SysConfigDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(userDo.SysConfig)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(userDo.SysConfig)
	uWhere.ConfigId = id
	uWhere.Deleted = common.Normal
	return config.Update(gtx, uData, uWhere)
}

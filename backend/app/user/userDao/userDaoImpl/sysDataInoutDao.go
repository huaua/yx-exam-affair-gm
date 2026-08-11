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

var sysDataInoutDaoImpl *SysDataInoutDaoImpl

func init() {
	sysDataInoutDaoImpl = &SysDataInoutDaoImpl{}
}

type SysDataInoutDaoImpl struct{}

func NewSysDataInoutDao() *SysDataInoutDaoImpl {
	return sysDataInoutDaoImpl
}

func (dataInout *SysDataInoutDaoImpl) GetById(id int64) *userDo.SysDataInout {
	if id < 1 {
		return nil
	}
	data := new(userDo.SysDataInout)
	where := new(userDo.SysDataInout)
	where.Id = id
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

func (dataInout *SysDataInoutDaoImpl) Count(dataInoutQuery *userDto.SysDataInoutQuery) int {
	dataInoutQuery.Deleted = common.Normal
	return query.ExecGetQueryCount[userDto.SysDataInoutQuery, userDo.SysDataInout](dataInoutQuery)
}

func (dataInout *SysDataInoutDaoImpl) List(dataInoutQuery *userDto.SysDataInoutQuery, page *page.Page) []*userDo.SysDataInout {
	dataInoutQuery.Deleted = common.Normal
	return query.ExecQueryList[userDto.SysDataInoutQuery, userDo.SysDataInout](dataInoutQuery, page)
}

func (dataInout *SysDataInoutDaoImpl) ListWithCondition(q *userDto.SysDataInoutQuery, page *page.Page, c string, args ...interface{}) []*userDo.SysDataInout {
	q.Deleted = common.Normal
	return query.ExecQueryListWithCondition[userDo.SysDataInout](q, page, c, args...)
}

func (dataInout *SysDataInoutDaoImpl) Insert(gtx tx.GTx, dataInoutDO *userDo.SysDataInout) bool {
	dataInoutDO.Deleted = common.Normal
	res := gtx.Create(dataInoutDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (dataInout *SysDataInoutDaoImpl) BatchInsert(gtx tx.GTx, dataInoutList []*userDo.SysDataInout) bool {
	res := gtx.CreateInBatches(dataInoutList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (dataInout *SysDataInoutDaoImpl) Update(gtx tx.GTx, uData *userDo.SysDataInout, uWhere *userDo.SysDataInout) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (dataInout *SysDataInoutDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(userDo.SysDataInout)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(userDo.SysDataInout)
	uWhere.Id = id
	uWhere.Deleted = common.Normal
	return dataInout.Update(gtx, uData, uWhere)
}

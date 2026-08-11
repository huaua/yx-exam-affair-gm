package bizDaoImpl

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/ds"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/query"
	"github.com/sealsee/web-base/public/ds/tx"
)

var eaDeptInfoDaoImpl *EaDeptInfoDaoImpl

func init() {
	eaDeptInfoDaoImpl = &EaDeptInfoDaoImpl{}
}

/*
 * desc: 学院基础数据层实现
 * author: 蓝鲸
 * date: 2024/08/14
 */
type EaDeptInfoDaoImpl struct{}

func NewEaDeptInfoDao() *EaDeptInfoDaoImpl {
	return eaDeptInfoDaoImpl
}

func (deptInfo *EaDeptInfoDaoImpl) GetById(id int64) *bizDo.EaDeptInfo {
	if id < 1 {
		return nil
	}
	data := new(bizDo.EaDeptInfo)
	where := new(bizDo.EaDeptInfo)
	where.DeptId = id
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

func (deptInfo *EaDeptInfoDaoImpl) Count(q *bizDto.EaDeptInfoQuery) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCount[bizDto.EaDeptInfoQuery, bizDo.EaDeptInfo](q)
}

func (deptInfo *EaDeptInfoDaoImpl) CountWithCondition(q *bizDto.EaDeptInfoQuery, condition interface{}, args ...interface{}) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCountWithCondition[bizDo.EaDeptInfo](q, condition, args...)
}

func (deptInfo *EaDeptInfoDaoImpl) List(q *bizDto.EaDeptInfoQuery, p *page.Page) []*bizDo.EaDeptInfo {
	q.Deleted = common.Normal
	return query.ExecQueryList[bizDto.EaDeptInfoQuery, bizDo.EaDeptInfo](q, p)
}

func (deptInfo *EaDeptInfoDaoImpl) ListWithCondition(q *bizDto.EaDeptInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaDeptInfo {
	q.Deleted = common.Normal
	return query.ExecQueryListWithCondition[bizDo.EaDeptInfo](q, p, condition, args...)
}

func (deptInfo *EaDeptInfoDaoImpl) ListMapWithCondition(q *bizDto.EaDeptInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
	q.Deleted = common.Normal
	return query.ExecQueryListMapWithCondition[bizDo.EaDeptInfo](q, p, condition, args...)
}

func (deptInfo *EaDeptInfoDaoImpl) Insert(gtx tx.GTx, deptInfoDO *bizDo.EaDeptInfo) bool {
	deptInfoDO.Deleted = common.Normal
	res := gtx.Create(deptInfoDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (deptInfo *EaDeptInfoDaoImpl) BatchInsert(gtx tx.GTx, deptInfoList []*bizDo.EaDeptInfo) bool {
	res := gtx.CreateInBatches(deptInfoList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (deptInfo *EaDeptInfoDaoImpl) Update(gtx tx.GTx, uData *bizDo.EaDeptInfo, uWhere *bizDo.EaDeptInfo) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (deptInfo *EaDeptInfoDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(bizDo.EaDeptInfo)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(bizDo.EaDeptInfo)
	uWhere.DeptId = id
	uWhere.Deleted = common.Normal
	return deptInfo.Update(gtx, uData, uWhere)
}

func (deptInfo *EaDeptInfoDaoImpl) UpdateNew(gtx tx.Db, uWhere *bizDo.EaDeptInfo, uData *bizDo.EaDeptInfo) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (deptInfo *EaDeptInfoDaoImpl) UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaDeptInfo, uData *bizDo.EaDeptInfo, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

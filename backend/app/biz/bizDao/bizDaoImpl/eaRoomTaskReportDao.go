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

var eaRoomTaskReportDaoImpl *EaRoomTaskReportDaoImpl

func init() {
	eaRoomTaskReportDaoImpl = &EaRoomTaskReportDaoImpl{}
}

/*
 * desc: 任务教室征集关联数据层实现
 * author: 蓝鲸
 * date: 2024/08/14
 */
type EaRoomTaskReportDaoImpl struct{}

func NewEaRoomTaskReportDao() *EaRoomTaskReportDaoImpl {
	return eaRoomTaskReportDaoImpl
}

func (roomTaskReport *EaRoomTaskReportDaoImpl) GetById(id int64) *bizDo.EaRoomTaskReport {
	if id < 1 {
		return nil
	}
	data := new(bizDo.EaRoomTaskReport)
	where := new(bizDo.EaRoomTaskReport)
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

func (roomTaskReport *EaRoomTaskReportDaoImpl) Count(q *bizDto.EaRoomTaskReportQuery) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCount[bizDto.EaRoomTaskReportQuery, bizDo.EaRoomTaskReport](q)
}

func (roomTaskReport *EaRoomTaskReportDaoImpl) CountWithCondition(q *bizDto.EaRoomTaskReportQuery, condition interface{}, args ...interface{}) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCountWithCondition[bizDo.EaRoomTaskReport](q, condition, args...)
}

func (roomTaskReport *EaRoomTaskReportDaoImpl) List(q *bizDto.EaRoomTaskReportQuery, p *page.Page) []*bizDo.EaRoomTaskReport {
	q.Deleted = common.Normal
	return query.ExecQueryList[bizDto.EaRoomTaskReportQuery, bizDo.EaRoomTaskReport](q, p)
}

func (roomTaskReport *EaRoomTaskReportDaoImpl) ListWithCondition(q *bizDto.EaRoomTaskReportQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaRoomTaskReport {
	q.Deleted = common.Normal
	return query.ExecQueryListWithCondition[bizDo.EaRoomTaskReport](q, p, condition, args...)
}

func (roomTaskReport *EaRoomTaskReportDaoImpl) ListMapWithCondition(q *bizDto.EaRoomTaskReportQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
	q.Deleted = common.Normal
	return query.ExecQueryListMapWithCondition[bizDo.EaRoomTaskReport](q, p, condition, args...)
}

func (roomTaskReport *EaRoomTaskReportDaoImpl) Insert(gtx tx.GTx, taskRoomInfoDO *bizDo.EaRoomTaskReport) bool {
	taskRoomInfoDO.Deleted = common.Normal
	res := gtx.Create(taskRoomInfoDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (roomTaskReport *EaRoomTaskReportDaoImpl) BatchInsert(gtx tx.GTx, taskRoomInfoList []*bizDo.EaRoomTaskReport) bool {
	res := gtx.CreateInBatches(taskRoomInfoList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (roomTaskReport *EaRoomTaskReportDaoImpl) Update(gtx tx.GTx, uData *bizDo.EaRoomTaskReport, uWhere *bizDo.EaRoomTaskReport) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (roomTaskReport *EaRoomTaskReportDaoImpl) DeleteById(gtx tx.GTx, id int64) bool {
	rlt := gtx.Where("id = ?", id).Delete(&bizDo.EaRoomTaskReport{})
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return int(rlt.RowsAffected) >= 0
}

func (roomTaskReport *EaRoomTaskReportDaoImpl) UpdateNew(gtx tx.Db, uWhere *bizDo.EaRoomTaskReport, uData *bizDo.EaRoomTaskReport) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (roomTaskReport *EaRoomTaskReportDaoImpl) UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaRoomTaskReport, uData *bizDo.EaRoomTaskReport, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

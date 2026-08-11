package bizDaoImpl

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	"yx-exam-affair-gm/app/cst"

	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/ds"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/query"
	"github.com/sealsee/web-base/public/ds/tx"
)

var eaTchTaskReportDaoImpl *EaTchTaskReportDaoImpl

func init() {
	eaTchTaskReportDaoImpl = &EaTchTaskReportDaoImpl{}
}

/*
 * desc: 老师任务分配数据层实现
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchTaskReportDaoImpl struct{}

func NewEaTchTaskReportDao() *EaTchTaskReportDaoImpl {
	return eaTchTaskReportDaoImpl
}

func (tchTaskReport *EaTchTaskReportDaoImpl) GetById(id int64) *bizDo.EaTchTaskReport {
	if id < 1 {
		return nil
	}
	data := new(bizDo.EaTchTaskReport)
	where := new(bizDo.EaTchTaskReport)
	where.ReportId = id
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

func (tchTaskReport *EaTchTaskReportDaoImpl) Count(q *bizDto.EaTchTaskReportQuery) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCount[bizDto.EaTchTaskReportQuery, bizDo.EaTchTaskReport](q)
}

func (tchTaskReport *EaTchTaskReportDaoImpl) CountWithCondition(q *bizDto.EaTchTaskReportQuery, condition interface{}, args ...interface{}) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCountWithCondition[bizDo.EaTchTaskReport](q, condition, args...)
}

func (tchTaskReport *EaTchTaskReportDaoImpl) List(q *bizDto.EaTchTaskReportQuery, p *page.Page) []*bizDo.EaTchTaskReport {
	q.Deleted = common.Normal
	return query.ExecQueryList[bizDto.EaTchTaskReportQuery, bizDo.EaTchTaskReport](q, p)
}

func (tchTaskReport *EaTchTaskReportDaoImpl) ListWithCondition(q *bizDto.EaTchTaskReportQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaTchTaskReport {
	q.Deleted = common.Normal
	return query.ExecQueryListWithCondition[bizDo.EaTchTaskReport](q, p, condition, args...)
}

func (tchTaskReport *EaTchTaskReportDaoImpl) ListMapWithCondition(q *bizDto.EaTchTaskReportQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
	q.Deleted = common.Normal
	return query.ExecQueryListMapWithCondition[bizDo.EaTchTaskReport](q, p, condition, args...)
}

func (tchTaskReport *EaTchTaskReportDaoImpl) Insert(gtx tx.GTx, tchTaskReportDO *bizDo.EaTchTaskReport) bool {
	tchTaskReportDO.Deleted = common.Normal
	res := gtx.Create(tchTaskReportDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (tchTaskReport *EaTchTaskReportDaoImpl) BatchInsert(gtx tx.GTx, tchTaskReportList []*bizDo.EaTchTaskReport) bool {
	res := gtx.CreateInBatches(tchTaskReportList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (tchTaskReport *EaTchTaskReportDaoImpl) Update(gtx tx.GTx, uData *bizDo.EaTchTaskReport, uWhere *bizDo.EaTchTaskReport) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (tchTaskReport *EaTchTaskReportDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(bizDo.EaTchTaskReport)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(bizDo.EaTchTaskReport)
	uWhere.ReportId = id
	uWhere.Deleted = common.Normal
	return tchTaskReport.Update(gtx, uData, uWhere)
}

func (tchTaskReport *EaTchTaskReportDaoImpl) DeleteByIds(gtx tx.Db, ids []interface{}, deleteBy int64) bool {
	uData := new(bizDo.EaTchTaskReport)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(bizDo.EaTchTaskReport)
	uWhere.AddIn("report_id", ids...)
	uWhere.Deleted = common.Normal
	return tchTaskReport.UpdateNew(gtx, uWhere, uData) > 0
}

func (tchTaskReport *EaTchTaskReportDaoImpl) DeleteByTaskDetailId(gtx tx.GTx, taskDetailId, deleteBy int64) bool {
	uData := new(bizDo.EaTchTaskReport)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(bizDo.EaTchTaskReport)
	uWhere.TaskDetailId = taskDetailId
	uWhere.Deleted = common.Normal
	return tchTaskReport.Update(gtx, uData, uWhere)
}

func (tchTaskReport *EaTchTaskReportDaoImpl) DeleteTempSavedTchReportOne(gtx tx.GTx, deptId, taskId, taskDetailId, tchId int64) bool {
	if deptId == 0 || taskId == 0 || taskDetailId == 0 {
		return false
	}
	uWhere := new(bizDo.EaTchTaskReport)
	uWhere.DeptId = deptId
	uWhere.TchTaskId = taskId
	uWhere.TaskDetailId = taskDetailId
	uWhere.TchId = tchId
	uWhere.State = cst.TEMP_SAVE
	rlt := gtx.Where(uWhere).Delete(&bizDo.EaTchTaskReport{})
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return int(rlt.RowsAffected) > 0
}

func (tchTaskReport *EaTchTaskReportDaoImpl) DeleteTempSavedTchReport(gtx tx.GTx, deptId, taskId int64) bool {
	if deptId == 0 || taskId == 0 {
		return false
	}
	uWhere := new(bizDo.EaTchTaskReport)
	uWhere.DeptId = deptId
	uWhere.TchTaskId = taskId
	uWhere.State = cst.TEMP_SAVE
	rlt := gtx.Where(uWhere).Delete(&bizDo.EaTchTaskReport{})
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return int(rlt.RowsAffected) > 0
}

func (tchTaskReport *EaTchTaskReportDaoImpl) UpdateNew(gtx tx.Db, uWhere *bizDo.EaTchTaskReport, uData *bizDo.EaTchTaskReport) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (tchTaskReport *EaTchTaskReportDaoImpl) UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaTchTaskReport, uData *bizDo.EaTchTaskReport, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

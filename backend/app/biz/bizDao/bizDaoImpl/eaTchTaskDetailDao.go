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

var eaTchTaskDetailDaoImpl *EaTchTaskDetailDaoImpl

func init() {
	eaTchTaskDetailDaoImpl = &EaTchTaskDetailDaoImpl{}
}

/*
 * desc: 监考老师征集任务明细数据层实现
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchTaskDetailDaoImpl struct{}

func NewEaTchTaskDetailDao() *EaTchTaskDetailDaoImpl {
	return eaTchTaskDetailDaoImpl
}

func (tchTaskDetail *EaTchTaskDetailDaoImpl) GetById(id int64) *bizDo.EaTchTaskDetail {
	if id < 1 {
		return nil
	}
	data := new(bizDo.EaTchTaskDetail)
	where := new(bizDo.EaTchTaskDetail)
	where.TaskDetailId = id
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

func (tchTaskDetail *EaTchTaskDetailDaoImpl) Count(q *bizDto.EaTchTaskDetailQuery) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCount[bizDto.EaTchTaskDetailQuery, bizDo.EaTchTaskDetail](q)
}

func (tchTaskDetail *EaTchTaskDetailDaoImpl) CountWithCondition(q *bizDto.EaTchTaskDetailQuery, condition interface{}, args ...interface{}) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCountWithCondition[bizDo.EaTchTaskDetail](q, condition, args...)
}

func (tchTaskDetail *EaTchTaskDetailDaoImpl) List(q *bizDto.EaTchTaskDetailQuery, p *page.Page) []*bizDo.EaTchTaskDetail {
	q.Deleted = common.Normal
	return query.ExecQueryList[bizDto.EaTchTaskDetailQuery, bizDo.EaTchTaskDetail](q, p)
}

func (tchTaskDetail *EaTchTaskDetailDaoImpl) ListWithCondition(q *bizDto.EaTchTaskDetailQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaTchTaskDetail {
	q.Deleted = common.Normal
	return query.ExecQueryListWithCondition[bizDo.EaTchTaskDetail](q, p, condition, args...)
}

func (tchTaskDetail *EaTchTaskDetailDaoImpl) ListMapWithCondition(q *bizDto.EaTchTaskDetailQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
	q.Deleted = common.Normal
	return query.ExecQueryListMapWithCondition[bizDo.EaTchTaskDetail](q, p, condition, args...)
}

func (tchTaskDetail *EaTchTaskDetailDaoImpl) Insert(gtx tx.GTx, tchTaskDetailDO *bizDo.EaTchTaskDetail) bool {
	tchTaskDetailDO.Deleted = common.Normal
	res := gtx.Create(tchTaskDetailDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (tchTaskDetail *EaTchTaskDetailDaoImpl) BatchInsert(gtx tx.GTx, tchTaskDetailList []*bizDo.EaTchTaskDetail) bool {
	res := gtx.CreateInBatches(tchTaskDetailList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (tchTaskDetail *EaTchTaskDetailDaoImpl) Update(gtx tx.GTx, uData *bizDo.EaTchTaskDetail, uWhere *bizDo.EaTchTaskDetail) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (tchTaskDetail *EaTchTaskDetailDaoImpl) UpdateReportNum(gtx tx.GTx, tchTaskDetailId int64, reportNum int) bool {
	res := gtx.Raw("update ea_tch_task_detail set collect_num = collect_num + ? where task_detail_id = ? and deleted = 1", reportNum, tchTaskDetailId).Scan(nil)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected) > 0
}

func (tchTaskDetail *EaTchTaskDetailDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(bizDo.EaTchTaskDetail)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(bizDo.EaTchTaskDetail)
	uWhere.TaskDetailId = id
	uWhere.Deleted = common.Normal
	return tchTaskDetail.Update(gtx, uData, uWhere)
}

func (tchTaskDetail *EaTchTaskDetailDaoImpl) UpdateNew(gtx tx.Db, uWhere *bizDo.EaTchTaskDetail, uData *bizDo.EaTchTaskDetail) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (tchTaskDetail *EaTchTaskDetailDaoImpl) UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaTchTaskDetail, uData *bizDo.EaTchTaskDetail, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

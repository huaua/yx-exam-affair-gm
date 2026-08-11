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

var eaTchTaskInfoDaoImpl *EaTchTaskInfoDaoImpl

func init() {
	eaTchTaskInfoDaoImpl = &EaTchTaskInfoDaoImpl{}
}

/*
 * desc: 监考老师征集数据层实现
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchTaskInfoDaoImpl struct{}

func NewEaTchTaskInfoDao() *EaTchTaskInfoDaoImpl {
	return eaTchTaskInfoDaoImpl
}

func (tchTaskInfo *EaTchTaskInfoDaoImpl) GetById(id int64) *bizDo.EaTchTaskInfo {
	if id < 1 {
		return nil
	}
	data := new(bizDo.EaTchTaskInfo)
	where := new(bizDo.EaTchTaskInfo)
	where.TchTaskId = id
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

func (tchTaskInfo *EaTchTaskInfoDaoImpl) Count(q *bizDto.EaTchTaskInfoQuery) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCount[bizDto.EaTchTaskInfoQuery, bizDo.EaTchTaskInfo](q)
}

func (tchTaskInfo *EaTchTaskInfoDaoImpl) CountWithCondition(q *bizDto.EaTchTaskInfoQuery, condition interface{}, args ...interface{}) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCountWithCondition[bizDo.EaTchTaskInfo](q, condition, args...)
}

func (tchTaskInfo *EaTchTaskInfoDaoImpl) List(q *bizDto.EaTchTaskInfoQuery, p *page.Page) []*bizDo.EaTchTaskInfo {
	q.Deleted = common.Normal
	return query.ExecQueryList[bizDto.EaTchTaskInfoQuery, bizDo.EaTchTaskInfo](q, p)
}

func (tchTaskInfo *EaTchTaskInfoDaoImpl) ListWithCondition(q *bizDto.EaTchTaskInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaTchTaskInfo {
	q.Deleted = common.Normal
	return query.ExecQueryListWithCondition[bizDo.EaTchTaskInfo](q, p, condition, args...)
}

func (tchTaskInfo *EaTchTaskInfoDaoImpl) ListMapWithCondition(q *bizDto.EaTchTaskInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
	q.Deleted = common.Normal
	return query.ExecQueryListMapWithCondition[bizDo.EaTchTaskInfo](q, p, condition, args...)
}

func (tchTaskInfo *EaTchTaskInfoDaoImpl) Insert(gtx tx.GTx, tchTaskInfoDO *bizDo.EaTchTaskInfo) bool {
	tchTaskInfoDO.Deleted = common.Normal
	res := gtx.Create(tchTaskInfoDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (tchTaskInfo *EaTchTaskInfoDaoImpl) BatchInsert(gtx tx.GTx, tchTaskInfoList []*bizDo.EaTchTaskInfo) bool {
	res := gtx.CreateInBatches(tchTaskInfoList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (tchTaskInfo *EaTchTaskInfoDaoImpl) Update(gtx tx.GTx, uData *bizDo.EaTchTaskInfo, uWhere *bizDo.EaTchTaskInfo) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (tchTaskInfo *EaTchTaskInfoDaoImpl) UpdateReportNum(gtx tx.GTx, tchTaskId int64, reportNum int) bool {
	res := gtx.Raw("update ea_tch_task_info set collect_num = collect_num + ? where tch_task_id = ? and deleted = 1", reportNum, tchTaskId).Scan(nil)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected) > 0
}

func (tchTaskInfo *EaTchTaskInfoDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(bizDo.EaTchTaskInfo)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(bizDo.EaTchTaskInfo)
	uWhere.TchTaskId = id
	uWhere.Deleted = common.Normal
	return tchTaskInfo.Update(gtx, uData, uWhere)
}

func (tchTaskInfo *EaTchTaskInfoDaoImpl) UpdateNew(gtx tx.Db, uWhere *bizDo.EaTchTaskInfo, uData *bizDo.EaTchTaskInfo) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (tchTaskInfo *EaTchTaskInfoDaoImpl) UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaTchTaskInfo, uData *bizDo.EaTchTaskInfo, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

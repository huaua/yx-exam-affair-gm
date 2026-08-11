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

var eaRoomTaskInfoDaoImpl *EaRoomTaskInfoDaoImpl

func init() {
	eaRoomTaskInfoDaoImpl = &EaRoomTaskInfoDaoImpl{}
}

/*
 * desc: 征集任务数据层实现
 * author: 蓝鲸
 * date: 2024/08/14
 */
type EaRoomTaskInfoDaoImpl struct{}

func NewEaRoomTaskInfoDao() *EaRoomTaskInfoDaoImpl {
	return eaRoomTaskInfoDaoImpl
}

func (roomTaskInfo *EaRoomTaskInfoDaoImpl) GetById(id int64) *bizDo.EaRoomTaskInfo {
	if id < 1 {
		return nil
	}
	data := new(bizDo.EaRoomTaskInfo)
	where := new(bizDo.EaRoomTaskInfo)
	where.TaskId = id
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

func (roomTaskInfo *EaRoomTaskInfoDaoImpl) Count(q *bizDto.EaRoomTaskInfoQuery) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCount[bizDto.EaRoomTaskInfoQuery, bizDo.EaRoomTaskInfo](q)
}

func (roomTaskInfo *EaRoomTaskInfoDaoImpl) CountWithCondition(q *bizDto.EaRoomTaskInfoQuery, condition interface{}, args ...interface{}) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCountWithCondition[bizDo.EaRoomTaskInfo](q, condition, args...)
}

func (roomTaskInfo *EaRoomTaskInfoDaoImpl) List(q *bizDto.EaRoomTaskInfoQuery, p *page.Page) []*bizDo.EaRoomTaskInfo {
	q.Deleted = common.Normal
	return query.ExecQueryList[bizDto.EaRoomTaskInfoQuery, bizDo.EaRoomTaskInfo](q, p)
}

func (roomTaskInfo *EaRoomTaskInfoDaoImpl) ListWithCondition(q *bizDto.EaRoomTaskInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaRoomTaskInfo {
	q.Deleted = common.Normal
	return query.ExecQueryListWithCondition[bizDo.EaRoomTaskInfo](q, p, condition, args...)
}

func (roomTaskInfo *EaRoomTaskInfoDaoImpl) ListMapWithCondition(q *bizDto.EaRoomTaskInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
	q.Deleted = common.Normal
	return query.ExecQueryListMapWithCondition[bizDo.EaRoomTaskInfo](q, p, condition, args...)
}

func (roomTaskInfo *EaRoomTaskInfoDaoImpl) Insert(gtx tx.GTx, taskInfoDO *bizDo.EaRoomTaskInfo) bool {
	taskInfoDO.Deleted = common.Normal
	res := gtx.Create(taskInfoDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (roomTaskInfo *EaRoomTaskInfoDaoImpl) BatchInsert(gtx tx.GTx, taskInfoList []*bizDo.EaRoomTaskInfo) bool {
	res := gtx.CreateInBatches(taskInfoList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (roomTaskInfo *EaRoomTaskInfoDaoImpl) Update(gtx tx.GTx, uData *bizDo.EaRoomTaskInfo, uWhere *bizDo.EaRoomTaskInfo) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (roomTaskInfo *EaRoomTaskInfoDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(bizDo.EaRoomTaskInfo)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(bizDo.EaRoomTaskInfo)
	uWhere.TaskId = id
	uWhere.Deleted = common.Normal
	return roomTaskInfo.Update(gtx, uData, uWhere)
}

func (roomTaskInfo *EaRoomTaskInfoDaoImpl) UpdateNew(gtx tx.Db, uWhere *bizDo.EaRoomTaskInfo, uData *bizDo.EaRoomTaskInfo) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (roomTaskInfo *EaRoomTaskInfoDaoImpl) UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaRoomTaskInfo, uData *bizDo.EaRoomTaskInfo, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (roomTaskInfo *EaRoomTaskInfoDaoImpl) UpdateCollectNumAndConfirmNum(gtx tx.Db, taskId int64, collectNum int, confirmNum int) int {
	res := gtx.Raw("update ea_room_task_info set collect_num = collect_num + ?, comfirm_num = comfirm_num + ? where task_id = ? and deleted = 1", collectNum, confirmNum, taskId).Scan(nil)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

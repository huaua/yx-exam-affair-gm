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

var eaJudgeTaskInfoDaoImpl *EaJudgeTaskInfoDaoImpl

func init() {
    eaJudgeTaskInfoDaoImpl = &EaJudgeTaskInfoDaoImpl{}
}

/*
 * desc: 评委征集任务数据层实现
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeTaskInfoDaoImpl struct{}

func NewEaJudgeTaskInfoDao() *EaJudgeTaskInfoDaoImpl {
    return eaJudgeTaskInfoDaoImpl
}

func (judgeTaskInfo *EaJudgeTaskInfoDaoImpl) GetById(id int64) *bizDo.EaJudgeTaskInfo {
    if id < 1 {
        return nil
    }
    data := new(bizDo.EaJudgeTaskInfo)
    where := new(bizDo.EaJudgeTaskInfo)
    where.JudgeTaskId = id
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

func (judgeTaskInfo *EaJudgeTaskInfoDaoImpl) Count(q *bizDto.EaJudgeTaskInfoQuery) int {
    q.Deleted = common.Normal
    return query.ExecGetQueryCount[bizDto.EaJudgeTaskInfoQuery, bizDo.EaJudgeTaskInfo](q)
}

func (judgeTaskInfo *EaJudgeTaskInfoDaoImpl) CountWithCondition(q *bizDto.EaJudgeTaskInfoQuery, condition interface{}, args ...interface{}) int {
    q.Deleted = common.Normal
    return query.ExecGetQueryCountWithCondition[bizDo.EaJudgeTaskInfo](q, condition, args...)
}

func (judgeTaskInfo *EaJudgeTaskInfoDaoImpl) List(q *bizDto.EaJudgeTaskInfoQuery, p *page.Page) []*bizDo.EaJudgeTaskInfo {
    q.Deleted = common.Normal
    return query.ExecQueryList[bizDto.EaJudgeTaskInfoQuery, bizDo.EaJudgeTaskInfo](q, p)
}

func (judgeTaskInfo *EaJudgeTaskInfoDaoImpl) ListWithCondition(q *bizDto.EaJudgeTaskInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaJudgeTaskInfo {
    q.Deleted = common.Normal
    return query.ExecQueryListWithCondition[bizDo.EaJudgeTaskInfo](q, p, condition, args...)
}

func (judgeTaskInfo *EaJudgeTaskInfoDaoImpl) ListMapWithCondition(q *bizDto.EaJudgeTaskInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
    q.Deleted = common.Normal
    return query.ExecQueryListMapWithCondition[bizDo.EaJudgeTaskInfo](q, p, condition, args...)
}

func (judgeTaskInfo *EaJudgeTaskInfoDaoImpl) Insert(gtx tx.GTx, judgeTaskInfoDO *bizDo.EaJudgeTaskInfo) bool {
    judgeTaskInfoDO.Deleted = common.Normal
    res := gtx.Create(judgeTaskInfoDO)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (judgeTaskInfo *EaJudgeTaskInfoDaoImpl) BatchInsert(gtx tx.GTx, judgeTaskInfoList []*bizDo.EaJudgeTaskInfo) bool {
    res := gtx.CreateInBatches(judgeTaskInfoList, 100)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

func (judgeTaskInfo *EaJudgeTaskInfoDaoImpl) Update(gtx tx.GTx, uData *bizDo.EaJudgeTaskInfo, uWhere *bizDo.EaJudgeTaskInfo) bool {
    uWhere.Deleted = common.Normal
    res := gtx.Where(uWhere).Updates(uData)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

func (judgeTaskInfo *EaJudgeTaskInfoDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
    uData := new(bizDo.EaJudgeTaskInfo)
    uData.SetDeleteBy(deleteBy)
    uWhere := new(bizDo.EaJudgeTaskInfo)
    uWhere.JudgeTaskId = id
    uWhere.Deleted = common.Normal
    return judgeTaskInfo.Update(gtx, uData, uWhere)
}

func (judgeTaskInfo *EaJudgeTaskInfoDaoImpl) UpdateNew(gtx tx.Db, uWhere *bizDo.EaJudgeTaskInfo, uData *bizDo.EaJudgeTaskInfo) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (judgeTaskInfo *EaJudgeTaskInfoDaoImpl) UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaJudgeTaskInfo, uData *bizDo.EaJudgeTaskInfo, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}
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

var eaJudgeExpertsAppraisalDaoImpl *EaJudgeExpertsAppraisalDaoImpl

func init() {
    eaJudgeExpertsAppraisalDaoImpl = &EaJudgeExpertsAppraisalDaoImpl{}
}

/*
 * desc: 评委评价数据层实现
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeExpertsAppraisalDaoImpl struct{}

func NewEaJudgeExpertsAppraisalDao() *EaJudgeExpertsAppraisalDaoImpl {
    return eaJudgeExpertsAppraisalDaoImpl
}

func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalDaoImpl) GetById(id int64) *bizDo.EaJudgeExpertsAppraisal {
    if id < 1 {
        return nil
    }
    data := new(bizDo.EaJudgeExpertsAppraisal)
    where := new(bizDo.EaJudgeExpertsAppraisal)
    where.AppraisalId = id
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

func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalDaoImpl) Count(q *bizDto.EaJudgeExpertsAppraisalQuery) int {
    q.Deleted = common.Normal
    return query.ExecGetQueryCount[bizDto.EaJudgeExpertsAppraisalQuery, bizDo.EaJudgeExpertsAppraisal](q)
}

func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalDaoImpl) CountWithCondition(q *bizDto.EaJudgeExpertsAppraisalQuery, condition interface{}, args ...interface{}) int {
    q.Deleted = common.Normal
    return query.ExecGetQueryCountWithCondition[bizDo.EaJudgeExpertsAppraisal](q, condition, args...)
}

func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalDaoImpl) List(q *bizDto.EaJudgeExpertsAppraisalQuery, p *page.Page) []*bizDo.EaJudgeExpertsAppraisal {
    q.Deleted = common.Normal
    return query.ExecQueryList[bizDto.EaJudgeExpertsAppraisalQuery, bizDo.EaJudgeExpertsAppraisal](q, p)
}

func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalDaoImpl) ListWithCondition(q *bizDto.EaJudgeExpertsAppraisalQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaJudgeExpertsAppraisal {
    q.Deleted = common.Normal
    return query.ExecQueryListWithCondition[bizDo.EaJudgeExpertsAppraisal](q, p, condition, args...)
}

func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalDaoImpl) ListMapWithCondition(q *bizDto.EaJudgeExpertsAppraisalQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
    q.Deleted = common.Normal
    return query.ExecQueryListMapWithCondition[bizDo.EaJudgeExpertsAppraisal](q, p, condition, args...)
}

func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalDaoImpl) Insert(gtx tx.GTx, judgeExpertsAppraisalDO *bizDo.EaJudgeExpertsAppraisal) bool {
    judgeExpertsAppraisalDO.Deleted = common.Normal
    res := gtx.Create(judgeExpertsAppraisalDO)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalDaoImpl) BatchInsert(gtx tx.GTx, judgeExpertsAppraisalList []*bizDo.EaJudgeExpertsAppraisal) bool {
    res := gtx.CreateInBatches(judgeExpertsAppraisalList, 100)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalDaoImpl) Update(gtx tx.GTx, uData *bizDo.EaJudgeExpertsAppraisal, uWhere *bizDo.EaJudgeExpertsAppraisal) bool {
    uWhere.Deleted = common.Normal
    res := gtx.Where(uWhere).Updates(uData)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
    uData := new(bizDo.EaJudgeExpertsAppraisal)
    uData.SetDeleteBy(deleteBy)
    uWhere := new(bizDo.EaJudgeExpertsAppraisal)
    uWhere.AppraisalId = id
    uWhere.Deleted = common.Normal
    return judgeExpertsAppraisal.Update(gtx, uData, uWhere)
}

func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalDaoImpl) UpdateNew(gtx tx.Db, uWhere *bizDo.EaJudgeExpertsAppraisal, uData *bizDo.EaJudgeExpertsAppraisal) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalDaoImpl) UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaJudgeExpertsAppraisal, uData *bizDo.EaJudgeExpertsAppraisal, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}
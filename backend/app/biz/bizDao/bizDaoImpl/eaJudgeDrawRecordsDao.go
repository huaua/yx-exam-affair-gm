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

var eaJudgeDrawRecordsDaoImpl *EaJudgeDrawRecordsDaoImpl

func init() {
    eaJudgeDrawRecordsDaoImpl = &EaJudgeDrawRecordsDaoImpl{}
}

/*
 * desc: 评委抽取规则数据层实现
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeDrawRecordsDaoImpl struct{}

func NewEaJudgeDrawRecordsDao() *EaJudgeDrawRecordsDaoImpl {
    return eaJudgeDrawRecordsDaoImpl
}

func (judgeDrawRecords *EaJudgeDrawRecordsDaoImpl) GetById(id int64) *bizDo.EaJudgeDrawRecords {
    if id < 1 {
        return nil
    }
    data := new(bizDo.EaJudgeDrawRecords)
    where := new(bizDo.EaJudgeDrawRecords)
    where.RulesId = id
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

func (judgeDrawRecords *EaJudgeDrawRecordsDaoImpl) Count(q *bizDto.EaJudgeDrawRecordsQuery) int {
    q.Deleted = common.Normal
    return query.ExecGetQueryCount[bizDto.EaJudgeDrawRecordsQuery, bizDo.EaJudgeDrawRecords](q)
}

func (judgeDrawRecords *EaJudgeDrawRecordsDaoImpl) CountWithCondition(q *bizDto.EaJudgeDrawRecordsQuery, condition interface{}, args ...interface{}) int {
    q.Deleted = common.Normal
    return query.ExecGetQueryCountWithCondition[bizDo.EaJudgeDrawRecords](q, condition, args...)
}

func (judgeDrawRecords *EaJudgeDrawRecordsDaoImpl) List(q *bizDto.EaJudgeDrawRecordsQuery, p *page.Page) []*bizDo.EaJudgeDrawRecords {
    q.Deleted = common.Normal
    return query.ExecQueryList[bizDto.EaJudgeDrawRecordsQuery, bizDo.EaJudgeDrawRecords](q, p)
}

func (judgeDrawRecords *EaJudgeDrawRecordsDaoImpl) ListWithCondition(q *bizDto.EaJudgeDrawRecordsQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaJudgeDrawRecords {
    q.Deleted = common.Normal
    return query.ExecQueryListWithCondition[bizDo.EaJudgeDrawRecords](q, p, condition, args...)
}

func (judgeDrawRecords *EaJudgeDrawRecordsDaoImpl) ListMapWithCondition(q *bizDto.EaJudgeDrawRecordsQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
    q.Deleted = common.Normal
    return query.ExecQueryListMapWithCondition[bizDo.EaJudgeDrawRecords](q, p, condition, args...)
}

func (judgeDrawRecords *EaJudgeDrawRecordsDaoImpl) Insert(gtx tx.GTx, judgeDrawRecordsDO *bizDo.EaJudgeDrawRecords) bool {
    judgeDrawRecordsDO.Deleted = common.Normal
    res := gtx.Create(judgeDrawRecordsDO)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (judgeDrawRecords *EaJudgeDrawRecordsDaoImpl) BatchInsert(gtx tx.GTx, judgeDrawRecordsList []*bizDo.EaJudgeDrawRecords) bool {
    res := gtx.CreateInBatches(judgeDrawRecordsList, 100)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

func (judgeDrawRecords *EaJudgeDrawRecordsDaoImpl) Update(gtx tx.GTx, uData *bizDo.EaJudgeDrawRecords, uWhere *bizDo.EaJudgeDrawRecords) bool {
    uWhere.Deleted = common.Normal
    res := gtx.Where(uWhere).Updates(uData)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

func (judgeDrawRecords *EaJudgeDrawRecordsDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
    uData := new(bizDo.EaJudgeDrawRecords)
    uData.SetDeleteBy(deleteBy)
    uWhere := new(bizDo.EaJudgeDrawRecords)
    uWhere.RulesId = id
    uWhere.Deleted = common.Normal
    return judgeDrawRecords.Update(gtx, uData, uWhere)
}

func (judgeDrawRecords *EaJudgeDrawRecordsDaoImpl) UpdateNew(gtx tx.Db, uWhere *bizDo.EaJudgeDrawRecords, uData *bizDo.EaJudgeDrawRecords) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (judgeDrawRecords *EaJudgeDrawRecordsDaoImpl) UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaJudgeDrawRecords, uData *bizDo.EaJudgeDrawRecords, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}
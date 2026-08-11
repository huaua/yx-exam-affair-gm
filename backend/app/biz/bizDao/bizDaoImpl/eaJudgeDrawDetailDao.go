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

var eaJudgeDrawDetailDaoImpl *EaJudgeDrawDetailDaoImpl

func init() {
    eaJudgeDrawDetailDaoImpl = &EaJudgeDrawDetailDaoImpl{}
}

/*
 * desc: 抽取评委明细数据层实现
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeDrawDetailDaoImpl struct{}

func NewEaJudgeDrawDetailDao() *EaJudgeDrawDetailDaoImpl {
    return eaJudgeDrawDetailDaoImpl
}

func (judgeDrawDetail *EaJudgeDrawDetailDaoImpl) GetById(id int64) *bizDo.EaJudgeDrawDetail {
    if id < 1 {
        return nil
    }
    data := new(bizDo.EaJudgeDrawDetail)
    where := new(bizDo.EaJudgeDrawDetail)
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

func (judgeDrawDetail *EaJudgeDrawDetailDaoImpl) Count(q *bizDto.EaJudgeDrawDetailQuery) int {
    q.Deleted = common.Normal
    return query.ExecGetQueryCount[bizDto.EaJudgeDrawDetailQuery, bizDo.EaJudgeDrawDetail](q)
}

func (judgeDrawDetail *EaJudgeDrawDetailDaoImpl) CountWithCondition(q *bizDto.EaJudgeDrawDetailQuery, condition interface{}, args ...interface{}) int {
    q.Deleted = common.Normal
    return query.ExecGetQueryCountWithCondition[bizDo.EaJudgeDrawDetail](q, condition, args...)
}

func (judgeDrawDetail *EaJudgeDrawDetailDaoImpl) List(q *bizDto.EaJudgeDrawDetailQuery, p *page.Page) []*bizDo.EaJudgeDrawDetail {
    q.Deleted = common.Normal
    return query.ExecQueryList[bizDto.EaJudgeDrawDetailQuery, bizDo.EaJudgeDrawDetail](q, p)
}

func (judgeDrawDetail *EaJudgeDrawDetailDaoImpl) ListWithCondition(q *bizDto.EaJudgeDrawDetailQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaJudgeDrawDetail {
    q.Deleted = common.Normal
    return query.ExecQueryListWithCondition[bizDo.EaJudgeDrawDetail](q, p, condition, args...)
}

func (judgeDrawDetail *EaJudgeDrawDetailDaoImpl) ListMapWithCondition(q *bizDto.EaJudgeDrawDetailQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
    q.Deleted = common.Normal
    return query.ExecQueryListMapWithCondition[bizDo.EaJudgeDrawDetail](q, p, condition, args...)
}

func (judgeDrawDetail *EaJudgeDrawDetailDaoImpl) Insert(gtx tx.GTx, judgeDrawDetailDO *bizDo.EaJudgeDrawDetail) bool {
    judgeDrawDetailDO.Deleted = common.Normal
    res := gtx.Create(judgeDrawDetailDO)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (judgeDrawDetail *EaJudgeDrawDetailDaoImpl) BatchInsert(gtx tx.GTx, judgeDrawDetailList []*bizDo.EaJudgeDrawDetail) bool {
    res := gtx.CreateInBatches(judgeDrawDetailList, 100)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

func (judgeDrawDetail *EaJudgeDrawDetailDaoImpl) Update(gtx tx.GTx, uData *bizDo.EaJudgeDrawDetail, uWhere *bizDo.EaJudgeDrawDetail) bool {
    uWhere.Deleted = common.Normal
    res := gtx.Where(uWhere).Updates(uData)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

func (judgeDrawDetail *EaJudgeDrawDetailDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
    uData := new(bizDo.EaJudgeDrawDetail)
    uData.SetDeleteBy(deleteBy)
    uWhere := new(bizDo.EaJudgeDrawDetail)
    uWhere.Id = id
    uWhere.Deleted = common.Normal
    return judgeDrawDetail.Update(gtx, uData, uWhere)
}

func (judgeDrawDetail *EaJudgeDrawDetailDaoImpl) UpdateNew(gtx tx.Db, uWhere *bizDo.EaJudgeDrawDetail, uData *bizDo.EaJudgeDrawDetail) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (judgeDrawDetail *EaJudgeDrawDetailDaoImpl) UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaJudgeDrawDetail, uData *bizDo.EaJudgeDrawDetail, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}
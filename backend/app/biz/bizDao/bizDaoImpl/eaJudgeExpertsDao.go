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

var eaJudgeExpertsDaoImpl *EaJudgeExpertsDaoImpl

func init() {
	eaJudgeExpertsDaoImpl = &EaJudgeExpertsDaoImpl{}
}

/*
 * desc: 评委专家库数据层实现
 * author: init code
 * date: 2025/06/06
 */
type EaJudgeExpertsDaoImpl struct{}

func NewEaJudgeExpertsDao() *EaJudgeExpertsDaoImpl {
	return eaJudgeExpertsDaoImpl
}

func (judgeExperts *EaJudgeExpertsDaoImpl) GetById(id int64) *bizDo.EaJudgeExperts {
	if id < 1 {
		return nil
	}
	data := new(bizDo.EaJudgeExperts)
	where := new(bizDo.EaJudgeExperts)
	where.ExpertId = id
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

func (judgeExperts *EaJudgeExpertsDaoImpl) Count(q *bizDto.EaJudgeExpertsQuery) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCount[bizDto.EaJudgeExpertsQuery, bizDo.EaJudgeExperts](q)
}

func (judgeExperts *EaJudgeExpertsDaoImpl) CountWithCondition(q *bizDto.EaJudgeExpertsQuery, condition interface{}, args ...interface{}) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCountWithCondition[bizDo.EaJudgeExperts](q, condition, args...)
}

func (judgeExperts *EaJudgeExpertsDaoImpl) List(q *bizDto.EaJudgeExpertsQuery, p *page.Page) []*bizDo.EaJudgeExperts {
	q.Deleted = common.Normal
	return query.ExecQueryList[bizDto.EaJudgeExpertsQuery, bizDo.EaJudgeExperts](q, p)
}

func (judgeExperts *EaJudgeExpertsDaoImpl) ListWithCondition(q *bizDto.EaJudgeExpertsQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaJudgeExperts {
	q.Deleted = common.Normal
	return query.ExecQueryListWithCondition[bizDo.EaJudgeExperts](q, p, condition, args...)
}

func (judgeExperts *EaJudgeExpertsDaoImpl) ListMapWithCondition(q *bizDto.EaJudgeExpertsQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
	q.Deleted = common.Normal
	return query.ExecQueryListMapWithCondition[bizDo.EaJudgeExperts](q, p, condition, args...)
}

func (judgeExperts *EaJudgeExpertsDaoImpl) Insert(gtx tx.GTx, judgeExpertsDO *bizDo.EaJudgeExperts) bool {
	judgeExpertsDO.Deleted = common.Normal
	res := gtx.Create(judgeExpertsDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (judgeExperts *EaJudgeExpertsDaoImpl) BatchInsert(gtx tx.GTx, judgeExpertsList []*bizDo.EaJudgeExperts) bool {
	res := gtx.CreateInBatches(judgeExpertsList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (judgeExperts *EaJudgeExpertsDaoImpl) Update(gtx tx.GTx, uData *bizDo.EaJudgeExperts, uWhere *bizDo.EaJudgeExperts) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (judgeExperts *EaJudgeExpertsDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(bizDo.EaJudgeExperts)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(bizDo.EaJudgeExperts)
	uWhere.ExpertId = id
	uWhere.Deleted = common.Normal
	return judgeExperts.Update(gtx, uData, uWhere)
}

func (judgeExperts *EaJudgeExpertsDaoImpl) UpdateNew(gtx tx.Db, uWhere *bizDo.EaJudgeExperts, uData *bizDo.EaJudgeExperts) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (judgeExperts *EaJudgeExpertsDaoImpl) UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaJudgeExperts, uData *bizDo.EaJudgeExperts, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

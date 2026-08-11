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

var eaTchInfoDaoImpl *EaTchInfoDaoImpl

func init() {
	eaTchInfoDaoImpl = &EaTchInfoDaoImpl{}
}

/*
 * desc: 教职工数据层实现
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchInfoDaoImpl struct{}

func NewEaTchInfoDao() *EaTchInfoDaoImpl {
	return eaTchInfoDaoImpl
}

func (tchInfo *EaTchInfoDaoImpl) GetById(id int64) *bizDo.EaTchInfo {
	if id < 1 {
		return nil
	}
	data := new(bizDo.EaTchInfo)
	where := new(bizDo.EaTchInfo)
	where.TchId = id
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

func (tchInfo *EaTchInfoDaoImpl) Count(q *bizDto.EaTchInfoQuery) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCount[bizDto.EaTchInfoQuery, bizDo.EaTchInfo](q)
}

func (tchInfo *EaTchInfoDaoImpl) CountWithCondition(q *bizDto.EaTchInfoQuery, condition interface{}, args ...interface{}) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCountWithCondition[bizDo.EaTchInfo](q, condition, args...)
}

func (tchInfo *EaTchInfoDaoImpl) List(q *bizDto.EaTchInfoQuery, p *page.Page) []*bizDo.EaTchInfo {
	q.Deleted = common.Normal
	return query.ExecQueryList[bizDto.EaTchInfoQuery, bizDo.EaTchInfo](q, p)
}

func (tchInfo *EaTchInfoDaoImpl) ListWithCondition(q *bizDto.EaTchInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaTchInfo {
	q.Deleted = common.Normal
	return query.ExecQueryListWithCondition[bizDo.EaTchInfo](q, p, condition, args...)
}

func (tchInfo *EaTchInfoDaoImpl) ListMapWithCondition(q *bizDto.EaTchInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
	q.Deleted = common.Normal
	return query.ExecQueryListMapWithCondition[bizDo.EaTchInfo](q, p, condition, args...)
}

func (tchInfo *EaTchInfoDaoImpl) Insert(gtx tx.GTx, tchInfoDO *bizDo.EaTchInfo) bool {
	tchInfoDO.Deleted = common.Normal
	res := gtx.Create(tchInfoDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (tchInfo *EaTchInfoDaoImpl) BatchInsert(gtx tx.GTx, tchInfoList []*bizDo.EaTchInfo) bool {
	res := gtx.CreateInBatches(tchInfoList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (tchInfo *EaTchInfoDaoImpl) Update(gtx tx.GTx, uData *bizDo.EaTchInfo, uWhere *bizDo.EaTchInfo) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (tchInfo *EaTchInfoDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(bizDo.EaTchInfo)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(bizDo.EaTchInfo)
	uWhere.TchId = id
	uWhere.Deleted = common.Normal
	return tchInfo.Update(gtx, uData, uWhere)
}

func (tchInfo *EaTchInfoDaoImpl) UpdateNew(gtx tx.Db, uWhere *bizDo.EaTchInfo, uData *bizDo.EaTchInfo) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (tchInfo *EaTchInfoDaoImpl) UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaTchInfo, uData *bizDo.EaTchInfo, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (tchInfo *EaTchInfoDaoImpl) ListAuditHistory(tchId int64) []*bizDo.EaTchInfoAuditHistory {
	if tchId < 1 {
		return nil
	}
	var list []*bizDo.EaTchInfoAuditHistory
	res := ds.GetGDB().Where("tch_id = ?", tchId).Order("change_time DESC").Find(&list)
	if res.Error != nil {
		panic(res.Error)
	}
	return list
}

func (tchInfo *EaTchInfoDaoImpl) GetLatestAuditHistory(tchId int64) *bizDo.EaTchInfoAuditHistory {
	if tchId < 1 {
		return nil
	}
	data := new(bizDo.EaTchInfoAuditHistory)
	res := ds.GetGDB().Where("tch_id = ? AND audit_status = ?", tchId, cst.TCH_AUDIT_PASSED).
		Order("change_time DESC").First(data)
	if res.Error != nil {
		return nil
	}
	if res.RowsAffected < 1 {
		return nil
	}
	return data
}

func (tchInfo *EaTchInfoDaoImpl) RefreshInvigilationExperienceByTchId(gtx tx.GTx, tchId int64) bool {
	if tchId < 1 {
		return false
	}
	// 1.判断是否还存在未删除的监考分配记录
	var cnt int64
	if err := gtx.Raw("select count(1) from ea_arrange_tch where deleted = 1 and tch_id = ?", tchId).Scan(&cnt).Error; err != nil {
		panic(err)
	}
	// 2.有记录则置为有经验，否则置为无经验
	experience := cst.TCH_NO_INVIGILATION_EXPERIENCE
	if cnt > 0 {
		experience = cst.TCH_HAS_INVIGILATION_EXPERIENCE
	}
	// 3.仅当与当前值不同时写回，避免无效更新触发 update_by 自动覆盖
	res := gtx.Raw("update ea_tch_info set invigilation_experience = ? where tch_id = ? and (invigilation_experience is null or invigilation_experience <> ?)",
		experience, tchId, experience).Scan(nil)
	if res.Error != nil {
		panic(res.Error)
	}
	return true
}

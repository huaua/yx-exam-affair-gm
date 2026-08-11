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

var eaArrangeTchDaoImpl *EaArrangeTchDaoImpl

func init() {
	eaArrangeTchDaoImpl = &EaArrangeTchDaoImpl{}
}

/*
 * desc: 监考老师分配数据层实现
 * author: init code
 * date: 2024/09/24
 */
type EaArrangeTchDaoImpl struct{}

func NewEaArrangeTchDao() *EaArrangeTchDaoImpl {
	return eaArrangeTchDaoImpl
}

func (arrangeTch *EaArrangeTchDaoImpl) GetById(id int64) *bizDo.EaArrangeTch {
	if id < 1 {
		return nil
	}
	data := new(bizDo.EaArrangeTch)
	where := new(bizDo.EaArrangeTch)
	where.TchArrangeId = id
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

func (arrangeTch *EaArrangeTchDaoImpl) Count(q *bizDto.EaArrangeTchQuery) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCount[bizDto.EaArrangeTchQuery, bizDo.EaArrangeTch](q)
}

func (arrangeTch *EaArrangeTchDaoImpl) CountWithCondition(q *bizDto.EaArrangeTchQuery, condition interface{}, args ...interface{}) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCountWithCondition[bizDo.EaArrangeTch](q, condition, args...)
}

func (arrangeTch *EaArrangeTchDaoImpl) List(q *bizDto.EaArrangeTchQuery, p *page.Page) []*bizDo.EaArrangeTch {
	q.Deleted = common.Normal
	return query.ExecQueryList[bizDto.EaArrangeTchQuery, bizDo.EaArrangeTch](q, p)
}

func (arrangeTch *EaArrangeTchDaoImpl) ListWithCondition(q *bizDto.EaArrangeTchQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaArrangeTch {
	q.Deleted = common.Normal
	return query.ExecQueryListWithCondition[bizDo.EaArrangeTch](q, p, condition, args...)
}

func (arrangeTch *EaArrangeTchDaoImpl) ListMapWithCondition(q *bizDto.EaArrangeTchQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
	q.Deleted = common.Normal
	return query.ExecQueryListMapWithCondition[bizDo.EaArrangeTch](q, p, condition, args...)
}

func (arrangeTch *EaArrangeTchDaoImpl) Insert(gtx tx.GTx, arrangeTchDO *bizDo.EaArrangeTch) bool {
	arrangeTchDO.Deleted = common.Normal
	res := gtx.Create(arrangeTchDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (arrangeTch *EaArrangeTchDaoImpl) BatchInsert(gtx tx.GTx, arrangeTchList []*bizDo.EaArrangeTch) bool {
	res := gtx.CreateInBatches(arrangeTchList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (arrangeTch *EaArrangeTchDaoImpl) Update(gtx tx.GTx, uData *bizDo.EaArrangeTch, uWhere *bizDo.EaArrangeTch) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (arrangeTch *EaArrangeTchDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(bizDo.EaArrangeTch)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(bizDo.EaArrangeTch)
	uWhere.TchArrangeId = id
	uWhere.Deleted = common.Normal
	return arrangeTch.Update(gtx, uData, uWhere)
}

func (arrangeTch *EaArrangeTchDaoImpl) DeleteByRoomReportId(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(bizDo.EaArrangeTch)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(bizDo.EaArrangeTch)
	uWhere.RoomReportId = id
	uWhere.Deleted = common.Normal
	return arrangeTch.Update(gtx, uData, uWhere)
}

func (arrangeTch *EaArrangeTchDaoImpl) DeleteByTchReportId(gtx tx.GTx, roomReportId int64, tchReportId int64, deleteBy int64) bool {
	uData := new(bizDo.EaArrangeTch)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(bizDo.EaArrangeTch)
	uWhere.RoomReportId = roomReportId
	uWhere.TchReportId = tchReportId
	uWhere.Deleted = common.Normal
	return arrangeTch.Update(gtx, uData, uWhere)
}

func (arrangeTch *EaArrangeTchDaoImpl) UpdateNew(gtx tx.Db, uWhere *bizDo.EaArrangeTch, uData *bizDo.EaArrangeTch) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (arrangeTch *EaArrangeTchDaoImpl) UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaArrangeTch, uData *bizDo.EaArrangeTch, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

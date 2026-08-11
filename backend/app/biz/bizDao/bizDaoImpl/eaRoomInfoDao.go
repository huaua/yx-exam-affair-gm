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

var eaRoomInfoDaoImpl *EaRoomInfoDaoImpl

func init() {
	eaRoomInfoDaoImpl = &EaRoomInfoDaoImpl{}
}

/*
 * desc: 教室数据层实现
 * author: 蓝鲸
 * date: 2024/08/14
 */
type EaRoomInfoDaoImpl struct{}

func NewEaRoomInfoDao() *EaRoomInfoDaoImpl {
	return eaRoomInfoDaoImpl
}

func (roomInfo *EaRoomInfoDaoImpl) GetById(id int64) *bizDo.EaRoomInfo {
	if id < 1 {
		return nil
	}
	data := new(bizDo.EaRoomInfo)
	where := new(bizDo.EaRoomInfo)
	where.RoomId = id
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

func (roomInfo *EaRoomInfoDaoImpl) Count(q *bizDto.EaRoomInfoQuery) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCount[bizDto.EaRoomInfoQuery, bizDo.EaRoomInfo](q)
}

func (roomInfo *EaRoomInfoDaoImpl) CountWithCondition(q *bizDto.EaRoomInfoQuery, condition interface{}, args ...interface{}) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCountWithCondition[bizDo.EaRoomInfo](q, condition, args...)
}

func (roomInfo *EaRoomInfoDaoImpl) List(q *bizDto.EaRoomInfoQuery, p *page.Page) []*bizDo.EaRoomInfo {
	q.Deleted = common.Normal
	return query.ExecQueryList[bizDto.EaRoomInfoQuery, bizDo.EaRoomInfo](q, p)
}

func (roomInfo *EaRoomInfoDaoImpl) ListWithCondition(q *bizDto.EaRoomInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaRoomInfo {
	q.Deleted = common.Normal
	return query.ExecQueryListWithCondition[bizDo.EaRoomInfo](q, p, condition, args...)
}

func (roomInfo *EaRoomInfoDaoImpl) ListMapWithCondition(q *bizDto.EaRoomInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
	q.Deleted = common.Normal
	return query.ExecQueryListMapWithCondition[bizDo.EaRoomInfo](q, p, condition, args...)
}

func (roomInfo *EaRoomInfoDaoImpl) Insert(gtx tx.GTx, classroomInfoDO *bizDo.EaRoomInfo) bool {
	classroomInfoDO.Deleted = common.Normal
	res := gtx.Create(classroomInfoDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (roomInfo *EaRoomInfoDaoImpl) BatchInsert(gtx tx.GTx, classroomInfoList []*bizDo.EaRoomInfo) bool {
	res := gtx.CreateInBatches(classroomInfoList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (roomInfo *EaRoomInfoDaoImpl) Update(gtx tx.GTx, uData *bizDo.EaRoomInfo, uWhere *bizDo.EaRoomInfo) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (roomInfo *EaRoomInfoDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(bizDo.EaRoomInfo)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(bizDo.EaRoomInfo)
	uWhere.RoomId = id
	uWhere.Deleted = common.Normal
	return roomInfo.Update(gtx, uData, uWhere)
}

func (roomInfo *EaRoomInfoDaoImpl) UpdateNew(gtx tx.Db, uWhere *bizDo.EaRoomInfo, uData *bizDo.EaRoomInfo) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (roomInfo *EaRoomInfoDaoImpl) UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaRoomInfo, uData *bizDo.EaRoomInfo, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

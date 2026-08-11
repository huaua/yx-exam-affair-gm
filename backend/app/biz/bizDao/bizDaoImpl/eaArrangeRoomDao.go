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

var eaArrangeRoomDaoImpl *EaArrangeRoomDaoImpl

func init() {
    eaArrangeRoomDaoImpl = &EaArrangeRoomDaoImpl{}
}

/*
 * desc: 考场编排详情数据层实现
 * author: init code
 * date: 2024/09/24
 */
type EaArrangeRoomDaoImpl struct{}

func NewEaArrangeRoomDao() *EaArrangeRoomDaoImpl {
    return eaArrangeRoomDaoImpl
}

func (arrangeRoom *EaArrangeRoomDaoImpl) GetById(id int64) *bizDo.EaArrangeRoom {
    if id < 1 {
        return nil
    }
    data := new(bizDo.EaArrangeRoom)
    where := new(bizDo.EaArrangeRoom)
    where.RoomArrangeId = id
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

func (arrangeRoom *EaArrangeRoomDaoImpl) Count(q *bizDto.EaArrangeRoomQuery) int {
    q.Deleted = common.Normal
    return query.ExecGetQueryCount[bizDto.EaArrangeRoomQuery, bizDo.EaArrangeRoom](q)
}

func (arrangeRoom *EaArrangeRoomDaoImpl) CountWithCondition(q *bizDto.EaArrangeRoomQuery, condition interface{}, args ...interface{}) int {
    q.Deleted = common.Normal
    return query.ExecGetQueryCountWithCondition[bizDo.EaArrangeRoom](q, condition, args...)
}

func (arrangeRoom *EaArrangeRoomDaoImpl) List(q *bizDto.EaArrangeRoomQuery, p *page.Page) []*bizDo.EaArrangeRoom {
    q.Deleted = common.Normal
    return query.ExecQueryList[bizDto.EaArrangeRoomQuery, bizDo.EaArrangeRoom](q, p)
}

func (arrangeRoom *EaArrangeRoomDaoImpl) ListWithCondition(q *bizDto.EaArrangeRoomQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaArrangeRoom {
    q.Deleted = common.Normal
    return query.ExecQueryListWithCondition[bizDo.EaArrangeRoom](q, p, condition, args...)
}

func (arrangeRoom *EaArrangeRoomDaoImpl) ListMapWithCondition(q *bizDto.EaArrangeRoomQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
    q.Deleted = common.Normal
    return query.ExecQueryListMapWithCondition[bizDo.EaArrangeRoom](q, p, condition, args...)
}

func (arrangeRoom *EaArrangeRoomDaoImpl) Insert(gtx tx.GTx, arrangeRoomDO *bizDo.EaArrangeRoom) bool {
    arrangeRoomDO.Deleted = common.Normal
    res := gtx.Create(arrangeRoomDO)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (arrangeRoom *EaArrangeRoomDaoImpl) BatchInsert(gtx tx.GTx, arrangeRoomList []*bizDo.EaArrangeRoom) bool {
    res := gtx.CreateInBatches(arrangeRoomList, 100)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

func (arrangeRoom *EaArrangeRoomDaoImpl) Update(gtx tx.GTx, uData *bizDo.EaArrangeRoom, uWhere *bizDo.EaArrangeRoom) bool {
    uWhere.Deleted = common.Normal
    res := gtx.Where(uWhere).Updates(uData)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

func (arrangeRoom *EaArrangeRoomDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
    uData := new(bizDo.EaArrangeRoom)
    uData.SetDeleteBy(deleteBy)
    uWhere := new(bizDo.EaArrangeRoom)
    uWhere.RoomArrangeId = id
    uWhere.Deleted = common.Normal
    return arrangeRoom.Update(gtx, uData, uWhere)
}

func (arrangeRoom *EaArrangeRoomDaoImpl) UpdateNew(gtx tx.Db, uWhere *bizDo.EaArrangeRoom, uData *bizDo.EaArrangeRoom) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (arrangeRoom *EaArrangeRoomDaoImpl) UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaArrangeRoom, uData *bizDo.EaArrangeRoom, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}
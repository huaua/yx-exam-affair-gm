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

var eaArrangeProfInfoDaoImpl *EaArrangeProfInfoDaoImpl

func init() {
    eaArrangeProfInfoDaoImpl = &EaArrangeProfInfoDaoImpl{}
}

/*
 * desc: 编排专业数据层实现
 * author: init code
 * date: 2024/09/24
 */
type EaArrangeProfInfoDaoImpl struct{}

func NewEaArrangeProfInfoDao() *EaArrangeProfInfoDaoImpl {
    return eaArrangeProfInfoDaoImpl
}

func (arrangeProfInfo *EaArrangeProfInfoDaoImpl) GetById(id int64) *bizDo.EaArrangeProfInfo {
    if id < 1 {
        return nil
    }
    data := new(bizDo.EaArrangeProfInfo)
    where := new(bizDo.EaArrangeProfInfo)
    where.InfoId = id
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

func (arrangeProfInfo *EaArrangeProfInfoDaoImpl) Count(q *bizDto.EaArrangeProfInfoQuery) int {
    q.Deleted = common.Normal
    return query.ExecGetQueryCount[bizDto.EaArrangeProfInfoQuery, bizDo.EaArrangeProfInfo](q)
}

func (arrangeProfInfo *EaArrangeProfInfoDaoImpl) CountWithCondition(q *bizDto.EaArrangeProfInfoQuery, condition interface{}, args ...interface{}) int {
    q.Deleted = common.Normal
    return query.ExecGetQueryCountWithCondition[bizDo.EaArrangeProfInfo](q, condition, args...)
}

func (arrangeProfInfo *EaArrangeProfInfoDaoImpl) List(q *bizDto.EaArrangeProfInfoQuery, p *page.Page) []*bizDo.EaArrangeProfInfo {
    q.Deleted = common.Normal
    return query.ExecQueryList[bizDto.EaArrangeProfInfoQuery, bizDo.EaArrangeProfInfo](q, p)
}

func (arrangeProfInfo *EaArrangeProfInfoDaoImpl) ListWithCondition(q *bizDto.EaArrangeProfInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaArrangeProfInfo {
    q.Deleted = common.Normal
    return query.ExecQueryListWithCondition[bizDo.EaArrangeProfInfo](q, p, condition, args...)
}

func (arrangeProfInfo *EaArrangeProfInfoDaoImpl) ListMapWithCondition(q *bizDto.EaArrangeProfInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
    q.Deleted = common.Normal
    return query.ExecQueryListMapWithCondition[bizDo.EaArrangeProfInfo](q, p, condition, args...)
}

func (arrangeProfInfo *EaArrangeProfInfoDaoImpl) Insert(gtx tx.GTx, arrangeProfInfoDO *bizDo.EaArrangeProfInfo) bool {
    arrangeProfInfoDO.Deleted = common.Normal
    res := gtx.Create(arrangeProfInfoDO)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (arrangeProfInfo *EaArrangeProfInfoDaoImpl) BatchInsert(gtx tx.GTx, arrangeProfInfoList []*bizDo.EaArrangeProfInfo) bool {
    res := gtx.CreateInBatches(arrangeProfInfoList, 100)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

func (arrangeProfInfo *EaArrangeProfInfoDaoImpl) Update(gtx tx.GTx, uData *bizDo.EaArrangeProfInfo, uWhere *bizDo.EaArrangeProfInfo) bool {
    uWhere.Deleted = common.Normal
    res := gtx.Where(uWhere).Updates(uData)
    if res.Error != nil {
        panic(res.Error)
    }
    return res.RowsAffected > 0
}

func (arrangeProfInfo *EaArrangeProfInfoDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
    uData := new(bizDo.EaArrangeProfInfo)
    uData.SetDeleteBy(deleteBy)
    uWhere := new(bizDo.EaArrangeProfInfo)
    uWhere.InfoId = id
    uWhere.Deleted = common.Normal
    return arrangeProfInfo.Update(gtx, uData, uWhere)
}

func (arrangeProfInfo *EaArrangeProfInfoDaoImpl) UpdateNew(gtx tx.Db, uWhere *bizDo.EaArrangeProfInfo, uData *bizDo.EaArrangeProfInfo) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (arrangeProfInfo *EaArrangeProfInfoDaoImpl) UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaArrangeProfInfo, uData *bizDo.EaArrangeProfInfo, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}
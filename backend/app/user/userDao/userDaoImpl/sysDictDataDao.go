package userDaoImpl

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/ds"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/query"
	"github.com/sealsee/web-base/public/ds/tx"
)

var sysDictDataDaoImpl *SysDictDataDaoImpl

func init() {
	sysDictDataDaoImpl = &SysDictDataDaoImpl{}
}

type SysDictDataDaoImpl struct{}

func NewSysDictDataDao() *SysDictDataDaoImpl {
	return sysDictDataDaoImpl
}

func (dictData *SysDictDataDaoImpl) GetById(id int64) *userDo.SysDictData {
	if id < 1 {
		return nil
	}
	data := new(userDo.SysDictData)
	where := new(userDo.SysDictData)
	where.DictCode = id
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

func (dictData *SysDictDataDaoImpl) GetByType(dictType string) []*userDo.SysDictData {
	if dictType == "" {
		return nil
	}
	q := new(userDto.SysDictDataQuery)
	q.Deleted = common.Normal
	q.DictType = dictType
	q.State = common.OK
	q.PageSize = page.MAX_PAGE_SIZE
	q.AddOrderAsc("dict_sort")
	return query.ExecQueryList[userDto.SysDictDataQuery, userDo.SysDictData](q, q.GetPage())
}

func (dictData *SysDictDataDaoImpl) Count(dictDataQuery *userDto.SysDictDataQuery) int {
	dictDataQuery.Deleted = common.Normal
	return query.ExecGetQueryCount[userDto.SysDictDataQuery, userDo.SysDictData](dictDataQuery)
}

func (dictData *SysDictDataDaoImpl) List(dictDataQuery *userDto.SysDictDataQuery, page *page.Page) []*userDo.SysDictData {
	dictDataQuery.Deleted = common.Normal
	return query.ExecQueryList[userDto.SysDictDataQuery, userDo.SysDictData](dictDataQuery, page)
}

func (dictData *SysDictDataDaoImpl) Insert(gtx tx.GTx, dictDataDO *userDo.SysDictData) bool {
	dictDataDO.Deleted = common.Normal
	res := gtx.Create(dictDataDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (dictData *SysDictDataDaoImpl) BatchInsert(gtx tx.GTx, dictDataList []*userDo.SysDictData) bool {
	res := gtx.CreateInBatches(dictDataList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (dictData *SysDictDataDaoImpl) Update(gtx tx.GTx, uData *userDo.SysDictData, uWhere *userDo.SysDictData) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (dictData *SysDictDataDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(userDo.SysDictData)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(userDo.SysDictData)
	uWhere.DictCode = id
	uWhere.Deleted = common.Normal
	return dictData.Update(gtx, uData, uWhere)
}

func (dictData *SysDictDataDaoImpl) GetCountByLabelAndValueAndType(where *userDo.SysDictData) int {
	dictDataQuery := &userDto.SysDictDataQuery{DictType: where.DictType}
	dictDataQuery.Deleted = common.Normal
	if where.DictCode > 1 {
		return query.ExecGetQueryCountWithCondition[userDo.SysDictData](dictDataQuery, "(dict_label = ? OR dict_value = ?) AND dict_code != ?", where.DictLabel, where.DictValue, where.DictCode)
	} else {
		return query.ExecGetQueryCountWithCondition[userDo.SysDictData](dictDataQuery, "(dict_label = ? OR dict_value = ?)", where.DictLabel, where.DictValue)
	}
}

func (dictData *SysDictDataDaoImpl) UpdateNew(gtx tx.Db, uWhere *userDo.SysDictData, uData *userDo.SysDictData) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

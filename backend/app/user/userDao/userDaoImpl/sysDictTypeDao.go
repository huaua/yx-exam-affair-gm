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

var sysDictTypeDaoImpl *SysDictTypeDaoImpl

func init() {
	sysDictTypeDaoImpl = &SysDictTypeDaoImpl{}
}

type SysDictTypeDaoImpl struct{}

func NewSysDictTypeDao() *SysDictTypeDaoImpl {
	return sysDictTypeDaoImpl
}

func (dictType *SysDictTypeDaoImpl) GetById(id int64) *userDo.SysDictType {
	if id < 1 {
		return nil
	}
	data := new(userDo.SysDictType)
	where := new(userDo.SysDictType)
	where.DictId = id
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

func (dictType *SysDictTypeDaoImpl) Count(dictTypeQuery *userDto.SysDictTypeQuery) int {
	dictTypeQuery.Deleted = common.Normal
	return query.ExecGetQueryCount[userDto.SysDictTypeQuery, userDo.SysDictType](dictTypeQuery)
}

func (dictType *SysDictTypeDaoImpl) List(dictTypeQuery *userDto.SysDictTypeQuery, page *page.Page) []*userDo.SysDictType {
	q := new(userDto.SysDictTypeQuery)
	q.Deleted = common.Normal
	var condition string
	var args []interface{}
	if len(dictTypeQuery.DictName) > 0 && len(dictTypeQuery.DictType) > 0 {
		condition = "dict_name like ? and dict_type like ?"
		return query.ExecQueryListWithCondition[userDo.SysDictType](q, page, condition, "%"+dictTypeQuery.DictName+"%", "%"+dictTypeQuery.DictType+"%")
	} else {
		if len(dictTypeQuery.DictName) > 0 {
			condition = "dict_name like ?"
			args = append(args, "%"+dictTypeQuery.DictName+"%")
		} else if len(dictTypeQuery.DictType) > 0 {
			condition = "dict_type like ?"
			args = append(args, "%"+dictTypeQuery.DictType+"%")
		}
		if len(condition) > 0 {
			return query.ExecQueryListWithCondition[userDo.SysDictType](q, page, condition, args)
		} else {
			return query.ExecQueryList[userDto.SysDictTypeQuery, userDo.SysDictType](q, page)
		}
	}
}

func (dictType *SysDictTypeDaoImpl) Insert(gtx tx.GTx, dictTypeDO *userDo.SysDictType) bool {
	dictTypeDO.Deleted = common.Normal
	res := gtx.Create(dictTypeDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (dictType *SysDictTypeDaoImpl) BatchInsert(gtx tx.GTx, dictTypeList []*userDo.SysDictType) bool {
	res := gtx.CreateInBatches(dictTypeList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (dictType *SysDictTypeDaoImpl) Update(gtx tx.GTx, uData *userDo.SysDictType, uWhere *userDo.SysDictType) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (dictType *SysDictTypeDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(userDo.SysDictType)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(userDo.SysDictType)
	uWhere.DictId = id
	uWhere.Deleted = common.Normal
	return dictType.Update(gtx, uData, uWhere)
}

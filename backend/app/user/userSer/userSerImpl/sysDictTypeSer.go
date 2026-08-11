package userSerImpl

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"
	"yx-exam-affair-gm/app/user/userDao"
	"yx-exam-affair-gm/app/user/userDao/userDaoImpl"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
	"github.com/sealsee/web-base/public/utils/file/excel"
)

type SysDictTypeService struct{}

var sysDictTypeService *SysDictTypeService
var sysDictTypeDao userDao.ISysDictTypeDao

func init() {
	sysDictTypeService = &SysDictTypeService{}
	sysDictTypeDao = userDaoImpl.NewSysDictTypeDao()
}

func NewSysDictTypeService() *SysDictTypeService {
	return sysDictTypeService
}

func (dictType *SysDictTypeService) GetSysDictTypeById(id int64) *userDo.SysDictType {
	return sysDictTypeDao.GetById(id)
}

func (dictType *SysDictTypeService) ListSysDictType(dictTypeQuery *userDto.SysDictTypeQuery, page *page.Page) []*userDo.SysDictType {
	return sysDictTypeDao.List(dictTypeQuery, page)
}

func (dictType *SysDictTypeService) InsertSysDictType(dictTypeDO *userDo.SysDictType) bool {
	return tx.ExecGTx(func(gtx tx.GTx) {
		sysDictTypeDao.Insert(gtx, dictTypeDO)
	})
}

func (dictType *SysDictTypeService) UpdateSysDictTypeById(dictTypeDO *userDo.SysDictType) bool {
	uWhere := &userDo.SysDictType{DictId: dictTypeDO.DictId}
	dictTypeDO.DictId = 0
	return tx.ExecGTx(func(gtx tx.GTx) {
		sysDictTypeDao.Update(gtx, dictTypeDO, uWhere)
	})
}

func (dictType *SysDictTypeService) DeleteSysDictTypeById(id int64, deleteBy int64) bool {
	return tx.ExecGTx(func(gtx tx.GTx) {
		sysDictTypeDao.DeleteById(gtx, id, deleteBy)
	})
}

// 注：仅支持1000条以内的简单同步导出 by codegen
func (dictType *SysDictTypeService) ExportSysDictType(dictTypeQuery *userDto.SysDictTypeQuery) []byte {
	nPage := page.NewPage()
	nPage.PageSize = page.MAX_PAGE_SIZE
	dataList := sysDictTypeDao.List(dictTypeQuery, nPage)
	row_title := []interface{}{"字典主键", "字典名称", "字典类型", "状态", "创建者", "创建时间", "更新者", "更新时间", "备注"}
	dataSlice := make([][]interface{}, 0)
	dataSlice = append(dataSlice, row_title)
	for _, v := range dataList {
		dataSlice = append(dataSlice, []interface{}{v.DictId, v.DictName, v.DictType, v.State, v.CreateBy, v.CreateTime, v.UpdateBy, v.UpdateTime, v.Remark})
	}
	return excel.ExportExcel(dataSlice)
}

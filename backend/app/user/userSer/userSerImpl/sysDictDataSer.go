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

type SysDictDataService struct{}

var sysDictDataService *SysDictDataService
var sysDictDataDao userDao.ISysDictDataDao

func init() {
	sysDictDataService = &SysDictDataService{}
	sysDictDataDao = userDaoImpl.NewSysDictDataDao()
}

func NewSysDictDataService() *SysDictDataService {
	return sysDictDataService
}

func (dictData *SysDictDataService) GetSysDictDataById(id int64) *userDo.SysDictData {
	return sysDictDataDao.GetById(id)
}

func (dictData *SysDictDataService) GetSysDictDataByType(dictType string) []*userDo.SysDictData {
	return sysDictDataDao.GetByType(dictType)
}

func (dictData *SysDictDataService) ListSysDictData(dictDataQuery *userDto.SysDictDataQuery, page *page.Page) []*userDo.SysDictData {
	return sysDictDataDao.List(dictDataQuery, page)
}

func (dictData *SysDictDataService) InsertSysDictData(dictDataDO *userDo.SysDictData) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysDictDataDao.Insert(gtx, dictDataDO)
	})
	return
}

func (dictData *SysDictDataService) UpdateSysDictDataById(dictDataDO *userDo.SysDictData) (success bool) {
	dictDataDO.SetNullableCols("remark")
	tx.ExecGTx(func(gtx tx.Db) {
		affectedRows := sysDictDataDao.UpdateNew(gtx, &userDo.SysDictData{DictCode: dictDataDO.DictCode}, dictDataDO)
		success = affectedRows > 0
	})
	return
}

func (dictData *SysDictDataService) DeleteSysDictDataById(id int64, deleteBy int64) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysDictDataDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

// 注：仅支持1000条以内的简单同步导出 by codegen
func (dictData *SysDictDataService) ExportSysDictData(dictDataQuery *userDto.SysDictDataQuery) []byte {
	nPage := page.NewPage()
	nPage.PageSize = page.MAX_PAGE_SIZE
	dataList := sysDictDataDao.List(dictDataQuery, nPage)
	row_title := []interface{}{"字典编码", "字典排序", "字典标签", "字典键值", "字典类型", "样式属性", "表格回显样式", "是否默认", "状态", "创建者", "创建时间", "更新者", "更新时间", "备注"}
	dataSlice := make([][]interface{}, 0)
	dataSlice = append(dataSlice, row_title)
	for _, v := range dataList {
		dataSlice = append(dataSlice, []interface{}{v.DictCode, v.DictSort, v.DictLabel, v.DictValue, v.DictType, v.CssClass, v.ListClass, v.IsDefault, v.State, v.CreateBy, v.CreateTime, v.UpdateBy, v.UpdateTime, v.Remark})
	}
	return excel.ExportExcel(dataSlice)
}

func (dictData *SysDictDataService) UpdateDicData(dictDataDO *userDo.SysDictData) (success bool) {
	uWhere := &userDo.SysDictData{DictCode: dictDataDO.DictCode}
	dictDataDO.DictCode = 0
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysDictDataDao.Update(gtx, dictDataDO, uWhere)
	})
	return
}

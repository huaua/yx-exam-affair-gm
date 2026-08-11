package bizSerImpl

import (
	"fmt"
	"github.com/sealsee/web-base/public/cst/common"
	"unicode/utf8"
	"yx-exam-affair-gm/app/biz/bizDao"
	"yx-exam-affair-gm/app/biz/bizDao/bizDaoImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
	"github.com/sealsee/web-base/public/utils/file/excel"
)

/*
 * desc: 学院基础服务层实现
 * author: 蓝鲸
 * date: 2024/08/14
 */
type EaDeptInfoService struct{}

var eaDeptInfoService *EaDeptInfoService
var eaDeptInfoDao bizDao.IEaDeptInfoDao

func init() {
	eaDeptInfoService = &EaDeptInfoService{}
	eaDeptInfoDao = bizDaoImpl.NewEaDeptInfoDao()
}

func NewEaDeptInfoService() *EaDeptInfoService {
	return eaDeptInfoService
}

func (deptInfo *EaDeptInfoService) GetEaDeptInfoById(id int64) *bizDo.EaDeptInfo {
	return eaDeptInfoDao.GetById(id)
}

func (deptInfo *EaDeptInfoService) ListEaDeptInfo(q *bizDto.EaDeptInfoQuery, p *page.Page) []*bizDo.EaDeptInfo {
	return eaDeptInfoDao.List(q, p)
}

func (deptInfo *EaDeptInfoService) ListAllEaDeptInfo(q *bizDto.EaDeptInfoQuery) []*bizDo.EaDeptInfo {
	listDepartAll := make([]*bizDo.EaDeptInfo, 0)
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE
	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		list := eaDeptInfoDao.List(q, pageInfo)
		if len(list) == 0 {
			break
		}
		listDepartAll = append(listDepartAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}
	return listDepartAll
}

func (deptInfo *EaDeptInfoService) GetAllDeptMap() map[int64]*bizDo.EaDeptInfo {
	deptMap := make(map[int64]*bizDo.EaDeptInfo)
	departAll := deptInfo.ListAllEaDeptInfo(&bizDto.EaDeptInfoQuery{})
	for _, i := range departAll {
		deptMap[i.DeptId] = i
	}
	return deptMap
}

func (deptInfo *EaDeptInfoService) BatchInsertEaDeptInfo(deptInfoList []*bizDo.EaDeptInfo, optId int64) (success bool, errInfo string) {

	// 获取所有已经存在的院校代码列表
	deptListInDb := deptInfo.ListAllEaDeptInfo(new(bizDto.EaDeptInfoQuery))

	// 遍历，校验字段长度和院系代码是否重复
	deptCodeMap := make(map[string]*bizDo.EaDeptInfo)
	for _, v := range deptListInDb {
		deptCodeMap[v.DeptCode] = v
	}
	for _, deptInfoDO := range deptInfoList {
		if len(deptInfoDO.DeptCode) <= 0 || len(deptInfoDO.DeptName) <= 0 {
			errInfo = "学院代码、学院名称不能为空"
			return false, errInfo
		}
		if utf8.RuneCountInString(deptInfoDO.DeptCode) > 15 {
			errInfo = "学院代码只能输入1-15个字"
			return false, errInfo
		}
		if utf8.RuneCountInString(deptInfoDO.DeptName) < 2 || utf8.RuneCountInString(deptInfoDO.DeptName) > 15 {
			errInfo = "学院名称只能输入2-15个字"
			return false, errInfo
		}
		if len(deptInfoDO.CampusName) > 0 {
			if utf8.RuneCountInString(deptInfoDO.CampusName) < 2 || utf8.RuneCountInString(deptInfoDO.CampusName) > 15 {
				errInfo = "所在校区只能输入2-15个字"
				return false, errInfo
			}
		}
		if deptCodeMap[deptInfoDO.DeptCode] != nil {
			errInfo = fmt.Sprintf("学院代码%s重复，请检查数据", deptInfoDO.DeptCode)
			return false, errInfo
		}

		deptCodeMap[deptInfoDO.DeptCode] = deptInfoDO
		deptInfoDO.Deleted = common.Normal
		deptInfoDO.SetCreateBy(optId)
		deptInfoDO.SetUpdateBy(optId)
	}

	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaDeptInfoDao.BatchInsert(gtx, deptInfoList)
	})
	return success, ""
}

func (deptInfo *EaDeptInfoService) UpdateEaDeptInfoById(deptInfoDO *bizDo.EaDeptInfo) (success bool, errInfo string) {
	if utf8.RuneCountInString(deptInfoDO.DeptName) < 2 || utf8.RuneCountInString(deptInfoDO.DeptName) > 15 {
		errInfo = "学院名称只能输入2-15个字"
		return false, errInfo
	}
	if len(deptInfoDO.CampusName) > 0 {
		if utf8.RuneCountInString(deptInfoDO.CampusName) < 2 || utf8.RuneCountInString(deptInfoDO.CampusName) > 15 {
			errInfo = "所在校区只能输入2-15个字"
			return false, errInfo
		}
	}
	uWhere := &bizDo.EaDeptInfo{DeptId: deptInfoDO.DeptId}
	deptInfoDO.DeptId = 0
	deptInfoDO.SetNullableCols("campus_name")
	tx.ExecGTx(func(gtx tx.Db) {
		success = eaDeptInfoDao.UpdateNew(gtx, uWhere, deptInfoDO) > 0
	})
	return true, ""
}

func (deptInfo *EaDeptInfoService) DeleteEaDeptInfoById(id int64, deleteBy int64) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaDeptInfoDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

// 导出服务 -- start
type EaDeptInfoExportStruct struct {
	Query *bizDto.EaDeptInfoQuery
}

func (h *EaDeptInfoExportStruct) Title() string {
	return "学院基础导出"
}

func (h *EaDeptInfoExportStruct) HeaderColumn() []string {
	return []string{"学院id,dept_id", "学院代码,dept_code", "学院名称,dept_name", "所在校区,campus_name", "是否删除：1-未删除；2-已删除；,deleted", "备注,remark", "创建人ID,create_by", "创建时间,create_time", "更新人ID,update_by", "更新时间,update_time"}
}

func (h *EaDeptInfoExportStruct) Rows(p *page.Page) []map[string]interface{} {
	return eaDeptInfoDao.ListMapWithCondition(h.Query, p, nil, nil)
}

func (h *EaDeptInfoExportStruct) Finish(url string) {
}

// 导出服务 -- end
func (deptInfo *EaDeptInfoService) ExportEaDeptInfo(q *bizDto.EaDeptInfoQuery) ([]byte, error) {
	return excel.NewExcel().ExportSync(&EaDeptInfoExportStruct{q})
}

package bizApi

import (
	"yx-exam-affair-gm/app/biz/bizSer"
	"yx-exam-affair-gm/app/biz/bizSer/bizSerImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	"yx-exam-affair-gm/app/cst"

	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/route"
	"github.com/sealsee/web-base/public/web"
)

var eaDeptInfoService bizSer.IEaDeptInfoService

func init() {
	route.AddGroup("/api/auth/ea_dept_info", "学院基础").
		AddPOSTHandel("/get", getDeptInfo, true, "查询学院基础详情").
		AddPOSTHandel("/list", listDeptInfo, true, "查询学院基础列表").
		AddPOSTHandelPowerR("/add", addDeptInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "新增学院基础").
		AddPOSTHandelPowerR("/edit", editDeptInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "修改学院基础").
		AddPOSTHandelPowerR("/del", delDeptInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "删除学院基础").
		AddPOSTHandelPowerR("/export", exportDeptInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "导出学院基础")

	eaDeptInfoService = bizSerImpl.NewEaDeptInfoService()
}

func getDeptInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaDeptInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.DeptId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	deptInfoDB := eaDeptInfoService.GetEaDeptInfoById(q.DeptId)
	json.SetObj(deptInfoDB).Render()
}

func listDeptInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaDeptInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	sessionDeptId := ctx.GetUser().DeptId
	if cst.BIZ_DEPTID_FOR_ADMISSIONS != sessionDeptId {
		// 获取当前用户的所属学院id，招办为-100
		q.DeptId = sessionDeptId
	}
	p := q.GetPage()
	list := eaDeptInfoService.ListEaDeptInfo(q, p)
	json.SetPageList(list, p).Render()
}

func addDeptInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	var departList []*bizDo.EaDeptInfo
	if err := c.ShouldBindJSON(&departList); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, errInfo := eaDeptInfoService.BatchInsertEaDeptInfo(departList, ctx.GetUserId())
	if !success {
		if len(errInfo) > 0 {
			json.SetErr(errInfo).Render()
			return
		}
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func editDeptInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	deptInfoDO := new(bizDo.EaDeptInfo)
	if err := c.ShouldBindJSON(deptInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if deptInfoDO.DeptId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	deptInfo2Update := new(bizDo.EaDeptInfo)
	deptInfo2Update.DeptId = deptInfoDO.DeptId
	deptInfo2Update.DeptName = deptInfoDO.DeptName
	deptInfo2Update.CampusName = deptInfoDO.CampusName
	deptInfo2Update.SetUpdateBy(ctx.GetUserId())
	success, errInfo := eaDeptInfoService.UpdateEaDeptInfoById(deptInfo2Update)
	if !success {
		if len(errInfo) > 0 {
			json.SetErr(errInfo).Render()
			return
		}
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func delDeptInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaDeptInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.DeptId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := eaDeptInfoService.DeleteEaDeptInfoById(q.DeptId, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func exportDeptInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaDeptInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data, err := eaDeptInfoService.ExportEaDeptInfo(q)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, data)
}

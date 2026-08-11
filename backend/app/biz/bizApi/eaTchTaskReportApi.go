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

var eaTchTaskReportService bizSer.IEaTchTaskReportService

func init() {
	route.AddGroup("/api/auth/ea_tch_task_report", "老师任务分配").
		AddPOSTHandel("/get", getTchTaskReport, true, "查询老师任务分配详情").
		AddPOSTHandelPowerR("/list", listTchTaskReport, true, cst.SYS_ROLE_KEY_DEPARTMENT, "查询老师任务分配列表").
		AddPOSTHandel("/temp_list", tempListTchTaskReport, true, "查询已经暂存的上报记录").
		AddPOSTHandel("/temp_save", tempSaveTchTaskReport, true, "暂存上报监考老师").
		AddPOSTHandel("/temp_clear", tempClearTchTaskReport, true, "清除暂存上报监考老师").
		AddPOSTHandel("/add", addTchTaskReport, true, "新增老师任务分配").
		AddPOSTHandel("/edit", editTchTaskReport, true, "修改老师任务分配").
		AddPOSTHandel("/del", delTchTaskReport, true, "取消上报").
		AddPOSTHandel("/export", exportTchTaskReport, true, "导出老师任务分配").
		AddPOSTHandel("/list_tch", listEnabledTchInfo, true, "查询可上报的教职工列表").
		AddPOSTHandel("/check_cancel", checkCancelResult, true, "检查是否可以取消上报").
		AddPOSTHandel("/replace_report", replaceReport, true, "替换上报").
		AddPOSTHandelPowerR("/list_by_task", listTchTaskReport, true, cst.SYS_ROLE_KEY_ADMISSIONS, "征集老师审核").
		AddPOSTHandelPowerR("/audit", auditTchTaskReport, true, cst.SYS_ROLE_KEY_ADMISSIONS, "审核征集状态").
		AddPOSTHandelPowerR("/reset", resetTchTaskReport, true, cst.SYS_ROLE_KEY_ADMISSIONS, "重置征集状态")

	eaTchTaskReportService = bizSerImpl.NewEaTchTaskReportService()
}

func getTchTaskReport(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaTchTaskReportQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.ReportId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	tchTaskReportDB := eaTchTaskReportService.GetEaTchTaskReportById(q.ReportId)
	json.SetObj(tchTaskReportDB).Render()
}

func listTchTaskReport(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaTchTaskReportQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	// 一个条件都没有的时候，就不返回数据
	if q.TchTaskId == 0 && q.TaskTime.IsZero() && q.DeptId == 0 && q.TchName == "" && q.JobNo == "" {
		json.Render()
		return
	}
	p := q.GetPage()
	user := ctx.GetUser()
	deptId := user.DeptId
	if deptId > 0 {
		q.DeptId = deptId
	}
	list := eaTchTaskReportService.ListEaTchTaskReport(q, p)
	json.SetPageList(list, p).Render()
}

func addTchTaskReport(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	var tchTaskReportList []*bizDo.EaTchTaskReport
	if err := c.ShouldBindJSON(&tchTaskReportList); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	if len(tchTaskReportList) <= 0 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	user := ctx.GetUser()
	userId := user.UserId
	deptId := user.DeptId
	if cst.BIZ_DEPTID_FOR_ADMISSIONS == deptId {
		json.SetErrs(errs.REFUSE_VISIT_ERR).Render()
		return
	}

	success, errMsg := eaTchTaskReportService.BatchInsertEaTchTaskReport(deptId, tchTaskReportList, userId)
	if !success {
		if len(errMsg) > 0 {
			json.SetErr(errMsg).Render()
			return
		}
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func editTchTaskReport(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	tchTaskReportDO := new(bizDto.EditTchTaskReportQuery)
	if err := c.ShouldBindJSON(tchTaskReportDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	if tchTaskReportDO.TaskDetailId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	sessionUser := ctx.GetUser()
	tchTaskReportDO.DeptId = sessionUser.DeptId
	success, errMsg := eaTchTaskReportService.EditEaTchTaskReport(tchTaskReportDO, sessionUser.UserId)
	if !success {
		if len(errMsg) > 0 {
			json.SetErr(errMsg).Render()
			return
		}
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func delTchTaskReport(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaTchTaskReportQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.ReportId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, err := eaTchTaskReportService.DeleteEaTchTaskReportById(q.ReportId, ctx.GetUserId())
	if !success {
		if err.Invalid() {
			json.SetErrs(err).Render()
			return
		}
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func exportTchTaskReport(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaTchTaskReportQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.TchTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data, err := eaTchTaskReportService.ExportEaTchTaskReport(q)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, data)
}

func listEnabledTchInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaTchTaskReportQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.TchTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	deptId := ctx.GetUser().DeptId
	if cst.BIZ_DEPTID_FOR_ADMISSIONS == deptId {
		json.SetErrs(errs.REFUSE_VISIT_ERR).Render()
		return
	}

	q.DeptId = deptId

	list, errMsg := eaTchTaskReportService.ListEnabledEaTchInfo(q)
	if len(errMsg) > 0 {
		json.SetErr(errMsg).Render()
		return
	}
	json.SetList(list).Render()
}

func checkCancelResult(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaTchTaskReportQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.ReportId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, _ := eaTchTaskReportService.CheckCancelResult(q.ReportId)
	json.SetObj(success).Render()
	//if !success {
	//	json.SetObj(success).Render()
	//	return
	//}
	//
	//json.Render()
}

func replaceReport(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDo.EaTchTaskReport)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.TaskDetailId < 1 || q.TchId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	q.SetUpdateBy(ctx.GetUserId())
	success, err := eaTchTaskReportService.ReplaceReport(q)
	if !success {
		if err.Invalid() {
			json.SetErrs(err).Render()
			return
		}
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func tempListTchTaskReport(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDo.EaTchTaskReport)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.TchTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	user := ctx.GetUser()
	deptId := user.DeptId
	if cst.BIZ_DEPTID_FOR_ADMISSIONS == deptId {
		json.SetErrs(errs.REFUSE_VISIT_ERR).Render()
		return
	}
	list := eaTchTaskReportService.TempListEaTchTaskReport(deptId, q.TchTaskId)
	json.SetList(list).Render()
}

func tempSaveTchTaskReport(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	var tchTaskReportList []*bizDo.EaTchTaskReport
	if err := c.ShouldBindJSON(&tchTaskReportList); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	if len(tchTaskReportList) <= 0 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	user := ctx.GetUser()
	userId := user.UserId
	deptId := user.DeptId
	if cst.BIZ_DEPTID_FOR_ADMISSIONS == deptId {
		json.SetErrs(errs.REFUSE_VISIT_ERR).Render()
		return
	}

	success, errMsg := eaTchTaskReportService.TempSaveEaTchTaskReport(deptId, tchTaskReportList, userId)
	if !success {
		if len(errMsg) > 0 {
			json.SetErr(errMsg).Render()
			return
		}
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func tempClearTchTaskReport(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	var tchTaskReportList []*bizDo.EaTchTaskReport
	if err := c.ShouldBindJSON(&tchTaskReportList); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	if len(tchTaskReportList) <= 0 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	user := ctx.GetUser()
	userId := user.UserId
	deptId := user.DeptId
	if cst.BIZ_DEPTID_FOR_ADMISSIONS == deptId {
		json.SetErrs(errs.REFUSE_VISIT_ERR).Render()
		return
	}

	success, errMsg := eaTchTaskReportService.TempClearEaTchTaskReport(deptId, tchTaskReportList, userId)
	if !success {
		if len(errMsg) > 0 {
			json.SetErr(errMsg).Render()
			return
		}
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func auditTchTaskReport(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.AuditTchTaskReportQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.ReportId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, err := eaTchTaskReportService.AuditTchTaskReport(q.ReportId, q.Action, ctx.GetUserId())
	if !success {
		if err.Invalid() {
			json.SetErrs(err).Render()
			return
		}
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func resetTchTaskReport(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.ResetTchTaskReportQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.ReportId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, err := eaTchTaskReportService.ResetTchTaskReport(q.ReportId, ctx.GetUserId())
	if !success {
		if err.Invalid() {
			json.SetErrs(err).Render()
			return
		}
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

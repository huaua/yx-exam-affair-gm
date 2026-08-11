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

var eaTchInfoService bizSer.IEaTchInfoService

func init() {
	route.AddGroup("/api/auth/ea_tch_info", "教职工").
		AddPOSTHandel("/get", getTchInfo, true, "查询教职工详情").
		AddPOSTHandel("/list", listTchInfo, true, "查询教职工列表").
		AddPOSTHandel("/add", addTchInfo, true, "新增教职工").
		AddPOSTHandel("/edit", editTchInfo, true, "修改教职工").
		AddPOSTHandelPowerR("/del", delTchInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "删除教职工").
		AddPOSTHandel("/export", exportTchInfo, true, "导出教职工").
		AddPOSTHandel("/import", importTchInfo, true, "导入教职工").
		AddPOSTHandelPowerR("/audit", auditTchInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "审核教职工").
		AddPOSTHandelPowerR("/change_enabled", changeTchEnabled, true, cst.SYS_ROLE_KEY_ADMISSIONS, "切换启用禁用状态").
		AddPOSTHandelPowerR("/change_enabled_batch", changeTchEnabledBatch, true, cst.SYS_ROLE_KEY_ADMISSIONS, "批量切换启用禁用状态").
		AddPOSTHandel("/arrange_records", listInvigilationRecords, true, "查询教职工监考记录").
		AddPOSTHandel("/audit_history", listAuditHistory, true, "查询教职工审核历史").
		AddPOSTHandel("/latest_audit_history", getLatestAuditHistory, true, "查询最近一次审核通过的快照").
		AddPOSTHandel("/submit", submitTchInfo, true, "学院提交教职工至待审核").
		AddPOSTHandel("/submit_batch", submitTchInfoBatch, true, "学院批量提交教职工至待审核").
		AddPOSTHandelPowerR("/next_pending", getNextPendingTchInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "查询下一位待审核教职工")

	eaTchInfoService = bizSerImpl.NewEaTchInfoService()
}

func getTchInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaTchInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.TchId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	tchInfoDB := eaTchInfoService.GetEaTchInfoById(q.TchId)
	if tchInfoDB == nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	deptId := context.NewUserContext(c).GetUser().DeptId
	if deptId > 0 && (tchInfoDB.DeptId != deptId || tchInfoDB.EnabledStatus != cst.TCH_ENABLED) {
		json.SetErrs(errs.REFUSE_VISIT_ERR).Render()
		return
	}
	json.SetObj(tchInfoDB).Render()
}

func listTchInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaTchInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	p := q.GetPage()

	if ctx.GetUser().DeptId > 0 {
		q.DeptId = ctx.GetUser().DeptId
		q.EnabledStatus = cst.TCH_ENABLED
		q.CollegeMode = true
	}

	list := eaTchInfoService.ListEaTchInfo(q, p)
	json.SetPageList(list, p).Render()
}

func addTchInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	tchInfoDO := new(bizDo.EaTchInfo)
	if err := c.ShouldBindJSON(tchInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	userInfo := ctx.GetUser()
	userId := userInfo.UserId
	deptId := userInfo.DeptId

	tchInfoDO.DateFrom = cst.ADDBYADMISSIONS
	if deptId > 0 {
		if tchInfoDO.SubmitAction != cst.TCH_SUBMIT_ACTION_DRAFT && tchInfoDO.SubmitAction != cst.TCH_SUBMIT_ACTION_SUBMIT {
			json.SetErr("请选择暂存或提交操作").Render()
			return
		}
		tchInfoDO.DeptId = deptId
		tchInfoDO.DateFrom = cst.ADDBYDEPT
		if tchInfoDO.SubmitAction == cst.TCH_SUBMIT_ACTION_SUBMIT {
			tchInfoDO.AuditStatus = cst.TCH_AUDIT_PENDING
		} else {
			tchInfoDO.AuditStatus = cst.TCH_AUDIT_UNSUBMITTED
		}
	}
	tchInfoDO.EnabledStatus = cst.TCH_ENABLED

	tchInfoDO.SetCreateBy(userId)
	tchInfoDO.SetUpdateBy(userId)
	success, errMsg := eaTchInfoService.InsertEaTchInfo(tchInfoDO)
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

func editTchInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	tchInfoDO := new(bizDo.EaTchInfo)
	if err := c.ShouldBindJSON(tchInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if tchInfoDO.TchId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	userInfo := ctx.GetUser()
	userId := userInfo.UserId
	deptId := userInfo.DeptId
	if deptId > 0 && tchInfoDO.SubmitAction != cst.TCH_SUBMIT_ACTION_DRAFT && tchInfoDO.SubmitAction != cst.TCH_SUBMIT_ACTION_SUBMIT {
		json.SetErr("请选择暂存或提交操作").Render()
		return
	}
	tchInfoDO.SetUpdateBy(userId)
	success, errMsg := eaTchInfoService.UpdateEaTchInfoById(tchInfoDO, deptId)
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

func delTchInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaTchInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.TchId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	userInfo := ctx.GetUser()
	userId := userInfo.UserId
	deptId := userInfo.DeptId
	success, errMsg := eaTchInfoService.DeleteEaTchInfoById(q.TchId, userId, deptId)
	if !success {
		if errMsg.Invalid() {
			json.SetErrs(errMsg).Render()
			return
		}
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func exportTchInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaTchInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if ctx.GetUser().DeptId > 0 {
		q.DeptId = ctx.GetUser().DeptId
		q.EnabledStatus = cst.TCH_ENABLED
		q.CollegeMode = true
	}
	data, err := eaTchInfoService.ExportEaTchInfo(q)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, data)
}

func importTchInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaTchInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	user := ctx.GetUser()
	q.OperId = user.UserId
	if user.DeptId > 0 {
		q.DeptId = user.DeptId
		q.CollegeMode = true
	}
	success, errMsg := eaTchInfoService.ImportEaTchInfo(q)
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

func auditTchInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	req := new(bizDto.EaTchInfoAuditReq)
	if err := c.ShouldBindJSON(req); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if req.TchId <= 0 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, errMsg := eaTchInfoService.AuditEaTchInfo(req.TchId, req.Action, req.Reason, ctx.GetUserId())
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

func getNextPendingTchInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaTchInfoQuery)
	// 兼容旧前端的空 body：解析失败时按无筛选条件处理
	if err := c.ShouldBindJSON(q); err != nil {
		q = new(bizDto.EaTchInfoQuery)
	}
	// 学院用户仅能在本学院范围内流转
	if deptId := ctx.GetUser().DeptId; deptId > 0 {
		q.DeptId = deptId
	}
	json.SetObj(eaTchInfoService.GetNextPendingTchInfo(q)).Render()
}

func changeTchEnabled(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	req := new(bizDto.EaTchInfoChangeEnabledReq)
	if err := c.ShouldBindJSON(req); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if req.TchId <= 0 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, errMsg := eaTchInfoService.ChangeEnabledStatus(req.TchId, req.EnabledStatus, ctx.GetUserId())
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

func changeTchEnabledBatch(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	req := new(bizDto.EaTchInfoChangeEnabledBatchReq)
	if err := c.ShouldBindJSON(req); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if len(req.TchIds) <= 0 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	affected, errMsg := eaTchInfoService.ChangeEnabledStatusBatch(req.TchIds, req.EnabledStatus, ctx.GetUserId())
	if errMsg != "" {
		json.SetErr(errMsg).Render()
		return
	}
	json.SetObj(map[string]int{"affected": affected}).Render()
}

func listInvigilationRecords(c *gin.Context) {
	json := web.NewJsonResult(c)
	req := new(bizDto.EaTchInfoArrangeRecordReq)
	if err := c.ShouldBindJSON(req); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if req.TchId <= 0 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	list := eaTchInfoService.ListInvigilationRecords(req.TchId)
	json.SetList(list).Render()
}

func listAuditHistory(c *gin.Context) {
	json := web.NewJsonResult(c)
	req := new(bizDto.EaTchInfoArrangeRecordReq)
	if err := c.ShouldBindJSON(req); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if req.TchId <= 0 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	list := eaTchInfoService.ListAuditHistory(req.TchId)
	json.SetList(list).Render()
}

func getLatestAuditHistory(c *gin.Context) {
	json := web.NewJsonResult(c)
	req := new(bizDto.EaTchInfoArrangeRecordReq)
	if err := c.ShouldBindJSON(req); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if req.TchId <= 0 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	json.SetObj(eaTchInfoService.GetLatestAuditHistory(req.TchId)).Render()
}

func submitTchInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	req := new(bizDto.EaTchInfoArrangeRecordReq)
	if err := c.ShouldBindJSON(req); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if req.TchId <= 0 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	user := ctx.GetUser()
	success, errMsg := eaTchInfoService.SubmitForAudit(req.TchId, user.DeptId, user.UserId)
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

// 学院批量提交教职工至待审核
type tchSubmitBatchReq struct {
	TchIds []int64 `json:"tchIds"`
}

func submitTchInfoBatch(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	req := new(tchSubmitBatchReq)
	if err := c.ShouldBindJSON(req); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	affected, msg := eaTchInfoService.SubmitForAuditBatch(req.TchIds, ctx.GetUser().DeptId, ctx.GetUserId())
	if msg != "" {
		json.SetErr(msg).Render()
		return
	}
	json.SetObj(map[string]int{"affected": affected}).Render()
}

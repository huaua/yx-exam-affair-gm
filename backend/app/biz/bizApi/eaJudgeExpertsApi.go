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

var eaJudgeExpertsService bizSer.IEaJudgeExpertsService

func init() {
	route.AddGroup("/api/auth/ea_judge_experts", "评委专家库").
		AddPOSTHandel("/get", getJudgeExperts, true, "查询评委专家库详情").
		AddPOSTHandel("/list", listJudgeExperts, true, "查询评委专家库列表").
		AddPOSTHandel("/add", addJudgeExperts, true, "新增评委专家库").
		AddPOSTHandel("/edit", editJudgeExperts, true, "修改评委专家库").
		AddPOSTHandel("/submit", submitJudgeExperts, true, "学院提交专家审核").
		AddPOSTHandel("/submit_batch", submitJudgeExpertsBatch, true, "学院批量提交专家审核").
		AddPOSTHandelPowerR("/audit", auditJudgeExperts, true, cst.SYS_ROLE_KEY_ADMISSIONS, "招办审核专家").
		AddPOSTHandelPowerR("/next_pending", getNextPendingJudgeExpert, true, cst.SYS_ROLE_KEY_ADMISSIONS, "查询下一位待审核专家").
		AddPOSTHandel("/latest_audit_history", getLatestJudgeExpertsAuditHistory, true, "查询最近一次审核通过快照").
		AddPOSTHandelPowerR("/evaluate", evaluateJudgeExperts, true, cst.SYS_ROLE_KEY_ADMISSIONS, "评价专家").
		AddPOSTHandelPowerR("/status", changeJudgeExpertsStatus, true, cst.SYS_ROLE_KEY_ADMISSIONS, "修改专家状态").
		AddPOSTHandel("/del", delJudgeExperts, true, "删除评委专家库").
		AddPOSTHandel("/import", importJudgeExperts, true, "导入评委专家库").
		AddPOSTHandel("/export", exportJudgeExperts, true, "导出评委专家库")
	route.AddGroup("/api/auth/ea_judge_experts_ban", "禁止抽取专家管理").
		AddPOSTHandelPowerR("/list", listBannedJudgeExperts, true, cst.SYS_ROLE_KEY_ADMISSIONS, "查询禁止抽取专家").
		AddPOSTHandelPowerR("/candidates", listBanCandidates, true, cst.SYS_ROLE_KEY_ADMISSIONS, "查询可禁止抽取专家").
		AddPOSTHandelPowerR("/add", addBannedJudgeExpert, true, cst.SYS_ROLE_KEY_ADMISSIONS, "新增禁止抽取专家").
		AddPOSTHandelPowerR("/del", delBannedJudgeExperts, true, cst.SYS_ROLE_KEY_ADMISSIONS, "删除禁止抽取专家").
		AddPOSTHandelPowerR("/del_all", delAllBannedJudgeExperts, true, cst.SYS_ROLE_KEY_ADMISSIONS, "删除全部禁止抽取专家").
		AddPOSTHandelPowerR("/import", importBannedJudgeExperts, true, cst.SYS_ROLE_KEY_ADMISSIONS, "导入禁止抽取专家").
		AddPOSTHandelPowerR("/export", exportBannedJudgeExperts, true, cst.SYS_ROLE_KEY_ADMISSIONS, "导出禁止抽取专家")

	eaJudgeExpertsService = bizSerImpl.NewEaJudgeExpertsService()
}

func listBannedJudgeExperts(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaJudgeExpertsQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	q.BannedOnly = true
	p := q.GetPage()
	json.SetPageList(eaJudgeExpertsService.ListEaJudgeExperts(q, p), p).Render()
}

func listBanCandidates(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaJudgeExpertsQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	// 只返回“启用”且“审核通过”且未被禁止抽取的专家
	q.AllowDraw = cst.YES
	q.AuditStatus = cst.EXPERT_AUDIT_PASSED
	q.ExcludeBanned = true
	p := q.GetPage()
	json.SetPageList(eaJudgeExpertsService.ListEaJudgeExperts(q, p), p).Render()
}

func addBannedJudgeExpert(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeExpertsQuery)
	if err := c.ShouldBindJSON(q); err != nil || q.ExpertId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if success, msg := eaJudgeExpertsService.BanEaJudgeExpert(q.ExpertId, ctx.GetUserId()); !success {
		json.SetErr(msg).Render()
		return
	}
	json.Render()
}

func delBannedJudgeExperts(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeExpertsIdsRequest)
	if err := c.ShouldBindJSON(q); err != nil || len(q.ExpertIds) < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if success, msg := eaJudgeExpertsService.UnbanEaJudgeExperts(q.ExpertIds, ctx.GetUserId()); !success {
		json.SetErr(msg).Render()
		return
	}
	json.Render()
}

func delAllBannedJudgeExperts(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	if success, msg := eaJudgeExpertsService.UnbanAllEaJudgeExperts(ctx.GetUserId()); !success {
		json.SetErr(msg).Render()
		return
	}
	json.Render()
}

func importBannedJudgeExperts(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeExpertsQuery)
	if err := c.ShouldBindJSON(q); err != nil || len(q.Files) < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	q.OperId = ctx.GetUserId()
	success, msg := eaJudgeExpertsService.ImportBannedEaJudgeExperts(q)
	if !success {
		json.SetErr(msg).Render()
		return
	}
	json.SetObj(msg).Render()
}

func exportBannedJudgeExperts(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaJudgeExpertsQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	q.BannedOnly = true
	data, err := eaJudgeExpertsService.ExportEaJudgeExperts(q)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, data)
}

func getJudgeExperts(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaJudgeExpertsQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.ExpertId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	judgeExpertsDB := eaJudgeExpertsService.GetEaJudgeExpertsById(q.ExpertId)
	json.SetObj(judgeExpertsDB).Render()
}

func listJudgeExperts(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeExpertsQuery)
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
	list := eaJudgeExpertsService.ListEaJudgeExperts(q, p)
	json.SetPageList(list, p).Render()
}

func addJudgeExperts(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	judgeExpertsDO := new(bizDo.EaJudgeExperts)
	if err := c.ShouldBindJSON(judgeExpertsDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	judgeExpertsDO.SetCreateBy(ctx.GetUserId())
	judgeExpertsDO.SetUpdateBy(ctx.GetUserId())
	// 获取当前用户的所属学院id，招办为-100
	deptId := ctx.GetUser().DeptId
	if deptId == cst.BIZ_DEPTID_FOR_ADMISSIONS {
		judgeExpertsDO.DataFrom = cst.ADDBYADMISSIONS
	} else {
		judgeExpertsDO.DataFrom = cst.ADDBYDEPT
		judgeExpertsDO.DeptId = deptId
		judgeExpertsDO.ExpertCategory = cst.EXPERT_CATEGORY_INTERNAL
		judgeExpertsDO.UnitName = ""
	}
	success, errInfo := eaJudgeExpertsService.InsertEaJudgeExperts(judgeExpertsDO)
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

func editJudgeExperts(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	judgeExpertsDO := new(bizDo.EaJudgeExperts)
	if err := c.ShouldBindJSON(judgeExpertsDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if judgeExpertsDO.ExpertId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	judgeExpertsDO.SetUpdateBy(ctx.GetUserId())
	success, errInfo := eaJudgeExpertsService.UpdateEaJudgeExpertsById(judgeExpertsDO, ctx.GetUser().DeptId)
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

func submitJudgeExperts(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeExpertsQuery)
	if err := c.ShouldBindJSON(q); err != nil || q.ExpertId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, msg := eaJudgeExpertsService.SubmitForAudit(q.ExpertId, ctx.GetUser().DeptId, ctx.GetUserId())
	if !success {
		json.SetErr(msg).Render()
		return
	}
	json.Render()
}

// 学院批量提交专家审核
type judgeExpertSubmitBatchReq struct {
	ExpertIds []int64 `json:"expertIds"`
}

func submitJudgeExpertsBatch(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	req := new(judgeExpertSubmitBatchReq)
	if err := c.ShouldBindJSON(req); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	affected, msg := eaJudgeExpertsService.SubmitForAuditBatch(req.ExpertIds, ctx.GetUser().DeptId, ctx.GetUserId())
	if msg != "" {
		json.SetErr(msg).Render()
		return
	}
	json.SetObj(map[string]int{"affected": affected}).Render()
}

func auditJudgeExperts(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeExpertsAuditRequest)
	if err := c.ShouldBindJSON(q); err != nil || q.ExpertId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, msg := eaJudgeExpertsService.AuditEaJudgeExperts(q.ExpertId, q.Action, q.Reason, ctx.GetUserId())
	if !success {
		json.SetErr(msg).Render()
		return
	}
	json.Render()
}

func getNextPendingJudgeExpert(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeExpertsQuery)
	// 兼容旧前端的空 body：解析失败时按无筛选条件处理
	if err := c.ShouldBindJSON(q); err != nil {
		q = new(bizDto.EaJudgeExpertsQuery)
	}
	// 非招办用户仅能在本学院范围内流转
	if sessionDeptId := ctx.GetUser().DeptId; cst.BIZ_DEPTID_FOR_ADMISSIONS != sessionDeptId {
		q.DeptId = sessionDeptId
	}
	json.SetObj(eaJudgeExpertsService.GetNextPendingExpert(q)).Render()
}

func getLatestJudgeExpertsAuditHistory(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaJudgeExpertsQuery)
	if err := c.ShouldBindJSON(q); err != nil || q.ExpertId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	json.SetObj(eaJudgeExpertsService.GetLatestAuditHistory(q.ExpertId)).Render()
}

func evaluateJudgeExperts(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeExpertsEvaluationRequest)
	if err := c.ShouldBindJSON(q); err != nil || q.ExpertId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, msg := eaJudgeExpertsService.EvaluateEaJudgeExperts(q.ExpertId, q.Evaluation, ctx.GetUserId())
	if !success {
		json.SetErr(msg).Render()
		return
	}
	json.Render()
}

func changeJudgeExpertsStatus(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeExpertsStatusRequest)
	if err := c.ShouldBindJSON(q); err != nil || q.ExpertId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, msg := eaJudgeExpertsService.ChangeAllowDraw(q.ExpertId, q.AllowDraw, ctx.GetUserId())
	if !success {
		json.SetErr(msg).Render()
		return
	}
	json.Render()
}

func delJudgeExperts(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeExpertsQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.ExpertId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := eaJudgeExpertsService.DeleteEaJudgeExpertsById(q.ExpertId, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func importJudgeExperts(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaJudgeExpertsQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	ctx := context.NewUserContext(c)
	// 获取当前用户的所属学院id，招办为-100
	deptId := ctx.GetUser().DeptId
	if deptId == cst.BIZ_DEPTID_FOR_ADMISSIONS {
		q.DataFrom = cst.ADDBYADMISSIONS
	} else {
		q.DataFrom = cst.ADDBYDEPT
	}
	q.DeptId = ctx.GetUser().DeptId
	q.OperId = ctx.GetUserId()
	// 导入数据
	success, msg := eaJudgeExpertsService.ImportEaJudgeExperts(q)
	if !success {
		json.SetErr(msg).Render()
		return
	}
	json.SetObj(msg).Render()
}

func exportJudgeExperts(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeExpertsQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	deptId := ctx.GetUser().DeptId
	if deptId != cst.BIZ_DEPTID_FOR_ADMISSIONS {
		q.DeptId = deptId
		q.CollegeView = true
	}
	data, err := eaJudgeExpertsService.ExportEaJudgeExperts(q)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, data)
}

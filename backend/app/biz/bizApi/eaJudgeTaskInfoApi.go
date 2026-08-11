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

var eaJudgeTaskInfoService bizSer.IEaJudgeTaskInfoService

func init() {
	route.AddGroup("/api/auth/ea_judge_task_info", "评委征集任务").
		AddPOSTHandelPowerR("/get", getJudgeTaskInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "查看评委征集任务详情").
		AddPOSTHandelPowerR("/list", listJudgeTaskInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "查询评委征集任务列表").
		AddPOSTHandel("/list_all", listAllJudgeTaskInfo, true, "查询所有评委征集任务下拉列表").
		AddPOSTHandelPowerR("/add", addJudgeTaskInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "新增评委征集任务").
		AddPOSTHandelPowerR("/edit", editJudgeTaskInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "修改评委征集任务").
		AddPOSTHandelPowerR("/del", delJudgeTaskInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "删除评委征集任务").
		AddPOSTHandelPowerR("/extract", extractJudgeTask, true, cst.SYS_ROLE_KEY_ADMISSIONS, "不按擅长科目进行抽取").
		AddPOSTHandelPowerR("/extract_by_subject", extractJudgeTaskBySubject, true, cst.SYS_ROLE_KEY_ADMISSIONS, "按擅长科目进行抽取").
		AddPOSTHandelPowerR("/cancel_extract", cancelExtract, true, cst.SYS_ROLE_KEY_ADMISSIONS, "取消抽取").
		AddPOSTHandelPowerR("/append_rounds", appendJudgeTaskRounds, true, cst.SYS_ROLE_KEY_ADMISSIONS, "追加评分轮次").
		AddPOSTHandelPowerR("/publish", publishJudgeTask, true, cst.SYS_ROLE_KEY_ADMISSIONS, "发布学院上报任务").
		AddPOSTHandelPowerR("/unpublish", unpublishJudgeTask, true, cst.SYS_ROLE_KEY_ADMISSIONS, "取消发布学院上报任务").
		AddPOSTHandelPowerR("/cancel_import", cancelImportReports, true, cst.SYS_ROLE_KEY_ADMISSIONS, "取消导入学院上报要求").
		AddPOSTHandelPowerR("/requirement/list", listReportRequirements, true, cst.SYS_ROLE_KEY_ADMISSIONS, "查询学院上报要求列表").
		AddPOSTHandelPowerR("/requirement/import", importReportRequirements, true, cst.SYS_ROLE_KEY_ADMISSIONS, "导入学院上报要求").
		AddPOSTHandelPowerR("/requirement/edit", editReportRequirement, true, cst.SYS_ROLE_KEY_ADMISSIONS, "修改单个上报要求").
		AddPOSTHandelPowerR("/requirement/batch_edit", batchEditReportRequirements, true, cst.SYS_ROLE_KEY_ADMISSIONS, "批量修改上报要求").
		AddPOSTHandelPowerR("/requirement/del", delReportRequirement, true, cst.SYS_ROLE_KEY_ADMISSIONS, "删除上报要求")

	eaJudgeTaskInfoService = bizSerImpl.NewEaJudgeTaskInfoService()
}

func appendJudgeTaskRounds(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeTaskInfoQuery)
	if c.ShouldBindJSON(q) != nil || q.JudgeTaskId < 1 || q.ScoreRounds < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, msg := eaJudgeTaskInfoService.AppendScoreRounds(q.JudgeTaskId, q.ScoreRounds, ctx.GetUserId())
	if !success {
		json.SetErr(msg).Render()
		return
	}
	json.Render()
}

func publishJudgeTask(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeTaskInfoQuery)
	if c.ShouldBindJSON(q) != nil || q.JudgeTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, msg := eaJudgeTaskInfoService.PublishCollegeReportTask(q.JudgeTaskId, ctx.GetUserId())
	if !success {
		json.SetErr(msg).Render()
		return
	}
	json.Render()
}

func unpublishJudgeTask(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeTaskInfoQuery)
	if c.ShouldBindJSON(q) != nil || q.JudgeTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, msg := eaJudgeTaskInfoService.UnpublishCollegeReportTask(q.JudgeTaskId, ctx.GetUserId())
	if !success {
		json.SetErr(msg).Render()
		return
	}
	json.Render()
}

func cancelImportReports(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeTaskInfoQuery)
	if c.ShouldBindJSON(q) != nil || q.JudgeTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, msg := eaJudgeTaskInfoService.CancelImportReports(q.JudgeTaskId, ctx.GetUserId())
	if !success {
		json.SetErr(msg).Render()
		return
	}
	json.Render()
}

func getJudgeTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaJudgeTaskInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.JudgeTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	judgeTaskInfoDB := eaJudgeTaskInfoService.GetEaJudgeTaskInfoById(q.JudgeTaskId)
	json.SetList(judgeTaskInfoDB).Render()
}

func listJudgeTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaJudgeTaskInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	p := q.GetPage()
	list := eaJudgeTaskInfoService.ListEaJudgeTaskInfo(q, p)
	json.SetPageList(list, p).Render()
}

func listAllJudgeTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaJudgeTaskInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	list := eaJudgeTaskInfoService.ListAllEaJudgeTaskInfo(q)
	json.SetList(list).Render()
}

func addJudgeTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	judgeTaskInfoDO := new(bizDo.EaJudgeTaskInfo)
	if err := c.ShouldBindJSON(judgeTaskInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	judgeTaskInfoDO.SetCreateBy(ctx.GetUserId())
	judgeTaskInfoDO.SetUpdateBy(ctx.GetUserId())
	success, errInfo := eaJudgeTaskInfoService.InsertEaJudgeTaskInfo(judgeTaskInfoDO)
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

func editJudgeTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	judgeTaskInfoDO := new(bizDo.EaJudgeTaskInfo)
	if err := c.ShouldBindJSON(judgeTaskInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if judgeTaskInfoDO.JudgeTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	judgeTaskInfoDO.SetUpdateBy(ctx.GetUserId())
	success, errInfo := eaJudgeTaskInfoService.UpdateEaJudgeTaskInfoById(judgeTaskInfoDO)
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

func delJudgeTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeTaskInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.JudgeTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, errInfo := eaJudgeTaskInfoService.DeleteEaJudgeTaskInfoById(q.JudgeTaskId, ctx.GetUserId())
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

func extractJudgeTask(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	extractJudgeTaskDTO := new(bizDto.ExtractJudgeTaskDTO)
	if err := c.ShouldBindJSON(extractJudgeTaskDTO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if extractJudgeTaskDTO.JudgeTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, errInfo := eaJudgeTaskInfoService.ExtractJudgeTask(extractJudgeTaskDTO, ctx.GetUserId())
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

func extractJudgeTaskBySubject(c *gin.Context) {

	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	var extractJudgeTaskDTOList []*bizDto.ExtractJudgeTaskDTO
	if err := c.ShouldBindJSON(&extractJudgeTaskDTOList); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if len(extractJudgeTaskDTOList) < 1 || extractJudgeTaskDTOList[0].JudgeTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, errInfo := eaJudgeTaskInfoService.ExtractJudgeTaskBySubject(extractJudgeTaskDTOList, ctx.GetUserId())
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

func cancelExtract(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeTaskInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.JudgeTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, errInfo := eaJudgeTaskInfoService.CancelExtract(q.JudgeTaskId, ctx.GetUserId())
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

func listReportRequirements(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaJudgeTaskReportRequirementQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.JudgeTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	requirementList := eaJudgeTaskInfoService.ListReportRequirements(q.JudgeTaskId)
	json.SetList(requirementList).Render()
}

func importReportRequirements(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	req := new(bizDto.EaJudgeTaskReportRequirementQuery)
	if err := c.ShouldBindJSON(req); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if req.JudgeTaskId < 1 || len(req.Files) < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	req.OperId = ctx.GetUserId()
	success, errInfo := eaJudgeTaskInfoService.ImportReportRequirements(req)
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

func editReportRequirement(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	req := new(bizDo.EaJudgeTaskReportRequirement)
	if err := c.ShouldBindJSON(req); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if req.RequirementId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, errInfo := eaJudgeTaskInfoService.UpdateReportRequirement(req.RequirementId, req.ExpertNum, ctx.GetUserId())
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

func batchEditReportRequirements(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	req := new(bizDto.BatchUpdateReportRequirementReq)
	if err := c.ShouldBindJSON(req); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if len(req.Items) < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, errInfo := eaJudgeTaskInfoService.BatchUpdateReportRequirements(req.Items, ctx.GetUserId())
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

func delReportRequirement(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	req := new(bizDto.EaJudgeTaskReportRequirementQuery)
	if err := c.ShouldBindJSON(req); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if req.RequirementId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, errInfo := eaJudgeTaskInfoService.DeleteReportRequirement(req.RequirementId, ctx.GetUserId())
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

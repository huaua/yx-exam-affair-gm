package bizApi

import (
	"yx-exam-affair-gm/app/biz/bizSer"
	"yx-exam-affair-gm/app/biz/bizSer/bizSerImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/route"
	"github.com/sealsee/web-base/public/web"
)

var eaJudgeExpertsAppraisalService bizSer.IEaJudgeExpertsAppraisalService

func init() {
	route.AddGroup("/api/auth/ea_judge_experts_appraisal", "评委评价").
		AddPOSTHandel("/get", getJudgeExpertsAppraisal, true, "查询评委评价详情").
		AddPOSTHandel("/list", listJudgeExpertsAppraisal, true, "查询评委评价列表").
		AddPOSTHandel("/add", addJudgeExpertsAppraisal, true, "新增评委评价").
		AddPOSTHandel("/edit", editJudgeExpertsAppraisal, true, "修改评委评价").
		AddPOSTHandel("/del", delJudgeExpertsAppraisal, true, "删除评委评价").
		AddPOSTHandel("/export", exportJudgeExpertsAppraisal, true, "导出评委评价")

	eaJudgeExpertsAppraisalService = bizSerImpl.NewEaJudgeExpertsAppraisalService()
}

func getJudgeExpertsAppraisal(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaJudgeExpertsAppraisalQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.AppraisalId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	judgeExpertsAppraisalDB := eaJudgeExpertsAppraisalService.GetEaJudgeExpertsAppraisalById(q.AppraisalId)
	json.SetObj(judgeExpertsAppraisalDB).Render()
}

func listJudgeExpertsAppraisal(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaJudgeExpertsAppraisalQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	p := q.GetPage()
	list := eaJudgeExpertsAppraisalService.ListEaJudgeExpertsAppraisal(q, p)
	json.SetPageList(list, p).Render()
}

func addJudgeExpertsAppraisal(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	judgeExpertsAppraisalDO := new(bizDo.EaJudgeExpertsAppraisal)
	if err := c.ShouldBindJSON(judgeExpertsAppraisalDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if judgeExpertsAppraisalDO.ExpertId < 1 || judgeExpertsAppraisalDO.JudgeTaskId < 1 ||
		len(judgeExpertsAppraisalDO.AppraisalContent) < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	judgeExpertsAppraisalDO.SetCreateBy(ctx.GetUserId())
	judgeExpertsAppraisalDO.SetUpdateBy(ctx.GetUserId())
	success, errInfo := eaJudgeExpertsAppraisalService.InsertEaJudgeExpertsAppraisal(judgeExpertsAppraisalDO)
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

func editJudgeExpertsAppraisal(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	judgeExpertsAppraisalDO := new(bizDo.EaJudgeExpertsAppraisal)
	if err := c.ShouldBindJSON(judgeExpertsAppraisalDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if judgeExpertsAppraisalDO.AppraisalId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	judgeExpertsAppraisalDO.SetUpdateBy(ctx.GetUserId())
	success := eaJudgeExpertsAppraisalService.UpdateEaJudgeExpertsAppraisalById(judgeExpertsAppraisalDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func delJudgeExpertsAppraisal(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeExpertsAppraisalQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.AppraisalId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := eaJudgeExpertsAppraisalService.DeleteEaJudgeExpertsAppraisalById(q.AppraisalId, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func exportJudgeExpertsAppraisal(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaJudgeExpertsAppraisalQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data, err := eaJudgeExpertsAppraisalService.ExportEaJudgeExpertsAppraisal(q)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, data)
}

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

var eaJudgeDrawDetailService bizSer.IEaJudgeDrawDetailService

func init() {
	route.AddGroup("/api/auth/ea_judge_draw_detail", "抽取评委明细").
		AddPOSTHandel("/list_by_role", listJudgeDrawDetail, true, "查询当前角色可查看的抽取评委列表").
		AddPOSTHandelPowerR("/list", listAllJudgeDrawDetail, true, cst.SYS_ROLE_KEY_ADMISSIONS, "查询所有抽取评委列表").
		AddPOSTHandel("/edit", editJudgeDrawDetail, true, "评委应答").
		AddPOSTHandelPowerR("/reset", resetJudgeDrawDetail, true, cst.SYS_ROLE_KEY_ADMISSIONS, "重置应答").
		AddPOSTHandelPowerR("/list_by_expert", listByExpert, true, cst.SYS_ROLE_KEY_ADMISSIONS, "根据专家查询参与的评分任务列表").
		AddPOSTHandel("/export", exportJudgeDrawDetail, true, "导出抽取评委明细")

	eaJudgeDrawDetailService = bizSerImpl.NewEaJudgeDrawDetailService()
}

func listJudgeDrawDetail(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaJudgeDrawDetailQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.JudgeTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	q.DeptId = ctx.GetUser().DeptId
	p := q.GetPage()
	list := eaJudgeDrawDetailService.ListEaJudgeDrawDetail(q, p)
	json.SetPageList(list, p).Render()
}

func listAllJudgeDrawDetail(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaJudgeDrawDetailQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.JudgeTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	p := q.GetPage()
	list := eaJudgeDrawDetailService.ListEaJudgeDrawDetail(q, p)
	json.SetPageList(list, p).Render()
}

func listByExpert(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaJudgeDrawDetailQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.ExpertId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	q.ConfirmState = cst.ExpertConfirmStatusMap[cst.ATTEND]
	list := eaJudgeDrawDetailService.ListAllEaJudgeDrawDetail(q, true)
	json.SetList(list).Render()
}

func editJudgeDrawDetail(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	judgeDrawDetailDO := new(bizDo.EaJudgeDrawDetail)
	if err := c.ShouldBindJSON(judgeDrawDetailDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if judgeDrawDetailDO.Id < 1 || !cst.ExpertConfirmStatusValueIsValid(judgeDrawDetailDO.ConfirmState) {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	judgeDrawDetailDO.SetUpdateBy(ctx.GetUserId())
	success, errInfo := eaJudgeDrawDetailService.UpdateEaJudgeDrawDetailById(judgeDrawDetailDO)
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

func resetJudgeDrawDetail(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	judgeDrawDetailDO := new(bizDo.EaJudgeDrawDetail)
	if err := c.ShouldBindJSON(judgeDrawDetailDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if judgeDrawDetailDO.Id < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	judgeDrawDetailDO.ConfirmState = cst.ExpertConfirmStatusMap[cst.NOT_CONFIRM]
	judgeDrawDetailDO.SetUpdateBy(ctx.GetUserId())
	success, errInfo := eaJudgeDrawDetailService.UpdateEaJudgeDrawDetailById(judgeDrawDetailDO)
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

func exportJudgeDrawDetail(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaJudgeDrawDetailQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data, err := eaJudgeDrawDetailService.ExportEaJudgeDrawDetail(q)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, data)
}

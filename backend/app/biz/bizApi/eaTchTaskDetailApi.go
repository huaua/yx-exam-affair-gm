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

var eaTchTaskDetailService bizSer.IEaTchTaskDetailService

func init() {
	route.AddGroup("/api/auth/ea_tch_task_detail", "监考老师征集任务明细").
		AddPOSTHandel("/get", getTchTaskDetail, true, "查询监考老师征集任务明细详情").
		AddPOSTHandel("/list", listTchTaskDetail, true, "查询监考老师征集任务明细列表").
		AddPOSTHandelPowerR("/add", addTchTaskDetail, true, cst.SYS_ROLE_KEY_ADMISSIONS, "新增监考老师征集任务明细").
		AddPOSTHandelPowerR("/edit", editTchTaskDetail, true, cst.SYS_ROLE_KEY_ADMISSIONS, "修改监考老师征集任务明细").
		AddPOSTHandelPowerR("/del", delTchTaskDetail, true, cst.SYS_ROLE_KEY_ADMISSIONS, "删除监考老师征集任务明细").
		AddPOSTHandel("/export", exportTchTaskDetail, true, "导出监考老师征集任务明细").
		AddPOSTHandelPowerR("/list_by_day", listTchTaskDetailByDay, true, cst.SYS_ROLE_KEY_ADMISSIONS, "查询监考老师征集任务列表按日分组")

	eaTchTaskDetailService = bizSerImpl.NewEaTchTaskDetailService()
}

func getTchTaskDetail(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaTchTaskDetailQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.TaskDetailId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	tchTaskDetailDB := eaTchTaskDetailService.GetEaTchTaskDetailById(q.TaskDetailId)
	json.SetObj(tchTaskDetailDB).Render()
}

func listTchTaskDetail(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaTchTaskDetailQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	p := q.GetPage()

	if q.OnlyQueryNeedReportTimeList {
		ctx := context.NewUserContext(c)
		deptId := ctx.GetUser().DeptId
		if cst.BIZ_DEPTID_FOR_ADMISSIONS == deptId {
			json.SetErrs(errs.REFUSE_VISIT_ERR).Render()
			return
		}
		q.DeptId = deptId
	}
	list := eaTchTaskDetailService.ListEaTchTaskDetail(q, p)
	json.SetPageList(list, p).Render()
}

func listTchTaskDetailByDay(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaTchTaskDetailQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.TchTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	list := eaTchTaskDetailService.ListEaTchTaskDetailByDay(q)
	json.SetList(list).Render()
}

func addTchTaskDetail(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	tchTaskDetailDO := new(bizDo.EaTchTaskDetail)
	if err := c.ShouldBindJSON(tchTaskDetailDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	tchTaskDetailDO.SetCreateBy(ctx.GetUserId())
	tchTaskDetailDO.SetUpdateBy(ctx.GetUserId())
	success := eaTchTaskDetailService.InsertEaTchTaskDetail(tchTaskDetailDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func editTchTaskDetail(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	tchTaskDetailDO := new(bizDo.EaTchTaskDetail)
	if err := c.ShouldBindJSON(tchTaskDetailDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if tchTaskDetailDO.TaskDetailId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	tchTaskDetailDO.SetUpdateBy(ctx.GetUserId())
	success := eaTchTaskDetailService.UpdateEaTchTaskDetailById(tchTaskDetailDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func delTchTaskDetail(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaTchTaskDetailQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.TaskDetailId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := eaTchTaskDetailService.DeleteEaTchTaskDetailById(q.TaskDetailId, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func exportTchTaskDetail(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaTchTaskDetailQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data, err := eaTchTaskDetailService.ExportEaTchTaskDetail(q)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, data)
}

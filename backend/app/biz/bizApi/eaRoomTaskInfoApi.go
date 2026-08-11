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

var eaRoomTaskInfoService bizSer.IEaRoomTaskInfoService

func init() {
	route.AddGroup("/api/auth/ea_task_info", "征集任务").
		AddPOSTHandel("/get", getTaskInfo, true, "查询征集任务详情").
		AddPOSTHandel("/list", listTaskInfo, true, "查询征集任务列表").
		AddPOSTHandel("/list_all", listAllPublishTaskInfo, true, "查询所有已发布的任务列表").
		AddPOSTHandelPowerR("/add", addTaskInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "新增征集任务").
		AddPOSTHandelPowerR("/edit", editTaskInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "修改征集任务").
		AddPOSTHandelPowerR("/change", editTaskInfoState, true, cst.SYS_ROLE_KEY_ADMISSIONS, "发布/取消发布征集任务").
		AddPOSTHandelPowerR("/del", delTaskInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "删除征集任务").
		AddPOSTHandelPowerR("/export", exportTaskInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "导出征集任务")

	eaRoomTaskInfoService = bizSerImpl.NewEaRoomTaskInfoService()
}

func getTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaRoomTaskInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.TaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	taskInfoDB := eaRoomTaskInfoService.GetEaTaskInfoById(q.TaskId)
	json.SetObj(taskInfoDB).Render()
}

func listTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaRoomTaskInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	p := q.GetPage()
	list := eaRoomTaskInfoService.ListEaTaskInfo(q, p)
	json.SetPageList(list, p).Render()
}

func listAllPublishTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaRoomTaskInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	q.State = cst.PUBLISH
	list := eaRoomTaskInfoService.ListAllEaTaskInfo(q)
	json.SetList(list).Render()
}

func addTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	taskInfoDO := new(bizDo.EaRoomTaskInfo)
	if err := c.ShouldBindJSON(taskInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	taskInfoDO.SetCreateBy(ctx.GetUserId())
	taskInfoDO.SetUpdateBy(ctx.GetUserId())
	success, errInfo := eaRoomTaskInfoService.InsertEaTaskInfo(taskInfoDO)
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

func editTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	taskInfoDO := new(bizDo.EaRoomTaskInfo)
	if err := c.ShouldBindJSON(taskInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if taskInfoDO.TaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	taskInfo2Edit := new(bizDo.EaRoomTaskInfo)
	taskInfo2Edit.TaskId = taskInfoDO.TaskId
	taskInfo2Edit.TaskName = taskInfoDO.TaskName
	taskInfo2Edit.TaskDay = taskInfoDO.TaskDay
	taskInfo2Edit.NeedNum = taskInfoDO.NeedNum
	taskInfo2Edit.LevelCode = taskInfoDO.LevelCode
	taskInfo2Edit.SetUpdateBy(ctx.GetUserId())
	success, errInfo := eaRoomTaskInfoService.UpdateEaTaskInfoById(taskInfo2Edit, false)
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

func editTaskInfoState(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	taskInfoDO := new(bizDo.EaRoomTaskInfo)
	if err := c.ShouldBindJSON(taskInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if taskInfoDO.TaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	taskInfo2Update := new(bizDo.EaRoomTaskInfo)
	taskInfo2Update.TaskId = taskInfoDO.TaskId
	taskInfo2Update.State = taskInfoDO.State
	taskInfo2Update.SetUpdateBy(ctx.GetUserId())
	success, errInfo := eaRoomTaskInfoService.UpdateEaTaskInfoById(taskInfo2Update, true)
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

func delTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaRoomTaskInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.TaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := eaRoomTaskInfoService.DeleteEaTaskInfoById(q.TaskId, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func exportTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaRoomTaskInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data, err := eaRoomTaskInfoService.ExportEaTaskInfo(q)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, data)
}

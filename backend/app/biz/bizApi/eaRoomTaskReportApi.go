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

var eaRoomTaskReportService bizSer.IEaRoomTaskReportService

func init() {
	route.AddGroup("/api/auth/ea_task_room_info", "任务教室征集关联").
		AddPOSTHandel("/get", getTaskRoomInfo, true, "查询任务教室征集关联详情").
		AddPOSTHandel("/list", listTaskRoomInfo, true, "查询绑定任务的教室列表").
		AddPOSTHandel("/list_all", listAllTaskRoomInfo, true, "获取任务下所有可使用的考场列表").
		AddPOSTHandel("/add", addTaskRoomInfo, true, "考场上报").
		AddPOSTHandel("/edit", editTaskRoomInfo, true, "确认上报").
		AddPOSTHandel("/del", delTaskRoomInfo, true, "取消上报").
		AddPOSTHandel("/reset", resetTaskRoomInfo, true, "重置征集确认").
		AddPOSTHandel("/reset_batch", resetBatchTaskRoomInfo, true, "批量重置征集确认").
		AddPOSTHandel("/confirm", confirmTaskRoomInfo, true, "征集确认").
		AddPOSTHandel("/export", exportTaskRoomInfo, true, "导出任务教室征集关联")

	eaRoomTaskReportService = bizSerImpl.NewEaRoomTaskReportService()
}

func getTaskRoomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaRoomTaskReportQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.Id < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	taskRoomInfoDB := eaRoomTaskReportService.GetEaTaskRoomInfoById(q.Id)
	json.SetObj(taskRoomInfoDB).Render()
}

func listTaskRoomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaRoomTaskReportQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.TaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.OnlyQuerySelfDept {
		ctx := context.NewUserContext(c)
		q.DeptId = ctx.GetUser().DeptId
	}
	list := eaRoomTaskReportService.ListAllEaTaskRoomInfo(q)
	json.SetList(list).Render()
}

func listAllTaskRoomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaRoomTaskReportQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.TaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	q.DeptId = ctx.GetUser().DeptId
	p := q.GetPage()
	// list := eaRoomTaskReportService.ListEaCanUseRoomListByTaskId(q.TaskId, ctx.GetUser().DeptId, p)
	enListNum, roomList := eaRoomTaskReportService.ListCanUseRoomListByTaskId(q, p)
	json.SetObj(enListNum)
	json.SetPageList(roomList, p).Render()
}

func addTaskRoomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	taskRoomInfoDO := new(bizDo.EaRoomTaskReport)
	if err := c.ShouldBindJSON(taskRoomInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if taskRoomInfoDO.RoomId < 1 || taskRoomInfoDO.TaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	taskRoomInfoDO.SetCreateBy(ctx.GetUserId())
	taskRoomInfoDO.SetUpdateBy(ctx.GetUserId())
	success, errInfo := eaRoomTaskReportService.InsertEaTaskRoomInfo(taskRoomInfoDO)
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

func editTaskRoomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	taskRoomInfoDO := new(bizDto.EaRoomConfirmDto)
	if err := c.ShouldBindJSON(taskRoomInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if taskRoomInfoDO.TaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	success, errInfo := eaRoomTaskReportService.BatchUpdateEaTaskRoomInfo(taskRoomInfoDO.RoomList, taskRoomInfoDO.TaskId, ctx.GetUserId())
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

func delTaskRoomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaRoomTaskReportQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.RoomId < 1 || q.TaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	success, errInfo := eaRoomTaskReportService.DeleteEaTaskRoomInfo(q)
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

func resetTaskRoomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaRoomTaskReportQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.Id < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	success, errInfo := eaRoomTaskReportService.ResetEaTaskRoomInfoById(q.Id, ctx.GetUserId())
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

func resetBatchTaskRoomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	dto := new(bizDto.EaRoomConfirmDto)
	if err := c.ShouldBindJSON(dto); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if dto.TaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if len(dto.RoomList) <= 0 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	success, errInfo := eaRoomTaskReportService.BatchResetEaTaskRoomInfo(dto.RoomList, dto.TaskId, ctx.GetUserId())
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

func confirmTaskRoomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaRoomTaskReportQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.Id < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	success, errInfo := eaRoomTaskReportService.ConfirmEaTaskRoomInfoById(q.Id, int64(q.TchNum), ctx.GetUserId())
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

func exportTaskRoomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaRoomTaskReportQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data, err := eaRoomTaskReportService.ExportEaTaskRoomInfo(q)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, data)
}

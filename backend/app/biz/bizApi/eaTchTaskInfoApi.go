package bizApi

import (
	"strings"
	"unicode/utf8"
	"yx-exam-affair-gm/app/biz/bizSer"
	"yx-exam-affair-gm/app/biz/bizSer/bizSerImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	"yx-exam-affair-gm/app/cst"
	"yx-exam-affair-gm/app/hint"

	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/route"
	"github.com/sealsee/web-base/public/utils/timeUtils"
	"github.com/sealsee/web-base/public/web"
)

var eaTchTaskInfoService bizSer.IEaTchTaskInfoService

func init() {
	route.AddGroup("/api/auth/ea_tch_task_info", "监考老师征集").
		AddPOSTHandel("/get", getTchTaskInfo, true, "查询监考老师征集详情").
		AddPOSTHandel("/list", listTchTaskInfo, true, "查询监考老师征集列表").
		AddPOSTHandelPowerR("/list_all_by_dept", listAllTchTaskInfoByDept, true, cst.SYS_ROLE_KEY_DEPARTMENT, "查询所有对学院有征集任务的列表").
		AddPOSTHandelPowerR("/add", addTchTaskInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "新增监考老师征集").
		AddPOSTHandelPowerR("/edit", editTchTaskInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "修改监考老师征集").
		AddPOSTHandelPowerR("/del", delTchTaskInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "删除监考老师征集").
		AddPOSTHandel("/export", exportTchTaskInfo, true, "导出监考老师征集").
		AddPOSTHandel("/check_info", checkTchTaskInfo, true, "校验老师征集任务信息").
		AddPOSTHandelPowerR("/change_state", changeTchTaskState, true, cst.SYS_ROLE_KEY_ADMISSIONS, "修改监考老师征集状态")

	eaTchTaskInfoService = bizSerImpl.NewEaTchTaskInfoService()
}

func getTchTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaTchTaskInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.TchTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	tchTaskInfoDB := eaTchTaskInfoService.GetEaTchTaskInfoById(q.TchTaskId)
	json.SetObj(tchTaskInfoDB).Render()
}

func listTchTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaTchTaskInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	p := q.GetPage()
	list := eaTchTaskInfoService.ListEaTchTaskInfo(q, p)
	json.SetPageList(list, p).Render()
}

func listAllTchTaskInfoByDept(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	deptId := ctx.GetUser().DeptId
	list := eaTchTaskInfoService.ListAllEaTchTaskInfoByDept(deptId)
	json.SetList(list).Render()
}

func checkTchTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	tchTaskInfoDO := new(bizDo.EaTchTaskInfo)
	if err := c.ShouldBindJSON(tchTaskInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	// 任务名称长度校验
	tchTaskName := strings.TrimSpace(tchTaskInfoDO.TchTaskName)
	tchTaskNameLen := utf8.RuneCountInString(tchTaskName)
	if tchTaskName == "" || tchTaskNameLen < 2 || tchTaskNameLen > 30 {
		json.SetErrs(hint.BIZ_TCH_TASK_NAME_LEN_ERR).Render()
		return
	}
	// 征集日期时间校验，起止范围最大7天
	if tchTaskInfoDO.TaskStartTime.IsZero() || tchTaskInfoDO.TaskEndTime.IsZero() {
		json.SetErrs(hint.BIZ_TCH_TASK_TIME_EMPTY_ERR).Render()
		return
	}
	diffHour := timeUtils.GetHourDiffer(tchTaskInfoDO.TaskStartTime.String(), tchTaskInfoDO.TaskEndTime.String())
	if diffHour < 0 || diffHour > (24*7) {
		json.SetErrs(hint.BIZ_TCH_TASK_TIME_DIFF_ERR).Render()
		return
	}
	// 根据名称查询是否存在记录
	q := new(bizDto.EaTchTaskInfoQuery)
	q.TchTaskName = tchTaskName
	if eaTchTaskInfoService.CountEaTchTaskInfo(q) > 0 {
		json.SetErrs(hint.BIZ_TCH_TASK_NAME_EXIST_ERR).Render()
		return
	}
	json.Render()
}

func addTchTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	tchTaskInfoDO := new(bizDo.EaTchTaskInfo)
	if err := c.ShouldBindJSON(tchTaskInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	// 任务名称长度校验
	tchTaskName := strings.TrimSpace(tchTaskInfoDO.TchTaskName)
	tchTaskNameLen := utf8.RuneCountInString(tchTaskName)
	if tchTaskName == "" || tchTaskNameLen < 2 || tchTaskNameLen > 30 {
		json.SetErrs(hint.BIZ_TCH_TASK_NAME_LEN_ERR).Render()
		return
	}
	// 征集日期时间校验，起止范围最大7天
	if tchTaskInfoDO.TaskStartTime.IsZero() || tchTaskInfoDO.TaskEndTime.IsZero() {
		json.SetErrs(hint.BIZ_TCH_TASK_TIME_EMPTY_ERR).Render()
		return
	}
	diffHour := timeUtils.GetHourDiffer(tchTaskInfoDO.TaskStartTime.String(), tchTaskInfoDO.TaskEndTime.String())
	if diffHour < 0 || diffHour > (24*7) {
		json.SetErrs(hint.BIZ_TCH_TASK_TIME_DIFF_ERR).Render()
		return
	}
	tchTaskInfoDO.SetCreateBy(ctx.GetUserId())
	tchTaskInfoDO.SetUpdateBy(ctx.GetUserId())
	err := eaTchTaskInfoService.InsertEaTchTaskInfo(tchTaskInfoDO)
	if err.Invalid() {
		json.SetErrs(err).Render()
		return
	}
	json.Render()
}

func changeTchTaskState(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	tchTaskInfoDO := new(bizDo.EaTchTaskInfo)
	if err := c.ShouldBindJSON(tchTaskInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if tchTaskInfoDO.TchTaskId < 1 || tchTaskInfoDO.State == "" ||
		(tchTaskInfoDO.State != cst.YES1 && tchTaskInfoDO.State != cst.NO2) {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := eaTchTaskInfoService.UpdateTchTaskState(tchTaskInfoDO.TchTaskId, tchTaskInfoDO.State, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func editTchTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	tchTaskInfoDO := new(bizDo.EaTchTaskInfo)
	if err := c.ShouldBindJSON(tchTaskInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if tchTaskInfoDO.TchTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	tchTaskInfoDO.SetUpdateBy(ctx.GetUserId())
	err := eaTchTaskInfoService.UpdateEaTchTaskInfoById(tchTaskInfoDO)
	if err.Invalid() {
		json.SetErrs(err).Render()
		return
	}
	json.Render()
}

func delTchTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaTchTaskInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.TchTaskId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := eaTchTaskInfoService.DeleteEaTchTaskInfoById(q.TchTaskId, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func exportTchTaskInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaTchTaskInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data, err := eaTchTaskInfoService.ExportEaTchTaskInfo(q)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, data)
}

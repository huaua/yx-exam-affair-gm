package bizApi

import (
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
	"github.com/sealsee/web-base/public/web"
)

var eaRoomInfoService bizSer.IEaRoomInfoService

func init() {
	route.AddGroup("/api/auth/ea_classroom_info", "教室").
		AddPOSTHandel("/get", getClassroomInfo, true, "查询教室详情").
		AddPOSTHandel("/list", listClassroomInfo, true, "查询教室列表").
		AddPOSTHandel("/add", addClassroomInfo, true, "新增教室").
		AddPOSTHandel("/edit", editClassroomInfo, true, "修改教室").
		AddPOSTHandel("/change", editClassroomInfoStatus, true, "启用/禁用教室").
		AddPOSTHandel("/del", delClassroomInfo, true, "删除教室").
		AddPOSTHandel("/export", exportClassroomInfo, true, "导出教室").
		AddPOSTHandel("/import", importClassroomInfo, true, "导入教室")

	eaRoomInfoService = bizSerImpl.NewEaRoomInfoService()
}

func getClassroomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	deptId := ctx.GetUser().DeptId

	q := new(bizDto.EaRoomInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.RoomId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	classroomInfoDB := eaRoomInfoService.GetEaClassroomInfoById(q.RoomId)
	if classroomInfoDB == nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if classroomInfoDB.DeptId != deptId {
		json.SetErrs(hint.SYS_NO_PERMISSION).Render()
		return
	}
	json.SetObj(classroomInfoDB).Render()
}

func listClassroomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaRoomInfoQuery)
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
	list := eaRoomInfoService.ListEaClassroomInfo(q, p)
	json.SetPageList(list, p).Render()
}

func addClassroomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	classroomInfoDO := new(bizDo.EaRoomInfo)
	if err := c.ShouldBindJSON(classroomInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	classroomInfoDO.SetCreateBy(ctx.GetUserId())
	classroomInfoDO.SetUpdateBy(ctx.GetUserId())
	// 获取当前用户的所属学院id，招办为-100
	classroomInfoDO.DeptId = ctx.GetUser().DeptId
	success, errInfo := eaRoomInfoService.InsertEaClassroomInfo(classroomInfoDO)
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

func editClassroomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	classroomInfoDO := new(bizDo.EaRoomInfo)
	if err := c.ShouldBindJSON(classroomInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if classroomInfoDO.RoomId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	classroomInfo2Update := new(bizDo.EaRoomInfo)
	classroomInfo2Update.RoomId = classroomInfoDO.RoomId
	classroomInfo2Update.CampusName = classroomInfoDO.CampusName
	classroomInfo2Update.RoomName = classroomInfoDO.RoomName
	classroomInfo2Update.GroupCapacity = classroomInfoDO.GroupCapacity
	classroomInfo2Update.Capacity = classroomInfoDO.Capacity
	classroomInfo2Update.GroupNum = classroomInfoDO.GroupNum
	classroomInfo2Update.SetUpdateBy(ctx.GetUserId())
	success, errInfo := eaRoomInfoService.UpdateEaClassroomInfoById(classroomInfo2Update, false, ctx.GetUser().DeptId)
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

func editClassroomInfoStatus(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	classroomInfoDO := new(bizDo.EaRoomInfo)
	if err := c.ShouldBindJSON(classroomInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if classroomInfoDO.RoomId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	classroomInfo2Update := new(bizDo.EaRoomInfo)
	classroomInfo2Update.RoomId = classroomInfoDO.RoomId
	classroomInfo2Update.State = classroomInfoDO.State
	classroomInfo2Update.SetUpdateBy(ctx.GetUserId())
	success, errInfo := eaRoomInfoService.UpdateEaClassroomInfoById(classroomInfo2Update, true, ctx.GetUser().DeptId)
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

func delClassroomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaRoomInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.RoomId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := eaRoomInfoService.DeleteEaClassroomInfoById(q.RoomId, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func exportClassroomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaRoomInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data, err := eaRoomInfoService.ExportEaClassroomInfo(q)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, data)
}

func importClassroomInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaRoomInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	ctx := context.NewUserContext(c)

	q.DeptId = ctx.GetUser().DeptId
	q.OperId = ctx.GetUserId()
	// 导入数据
	success, msg := eaRoomInfoService.ImportEaClassroomInfo(q)
	if !success {
		json.SetErr(msg).Render()
		return
	}
	json.SetObj(msg).Render()
}

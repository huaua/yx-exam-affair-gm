package bizApi

import (
	"yx-exam-affair-gm/app/biz/bizSer"
	"yx-exam-affair-gm/app/biz/bizSer/bizSerImpl"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	"yx-exam-affair-gm/app/cst"

	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/route"
	"github.com/sealsee/web-base/public/web"
)

var eaArrangeRoomService bizSer.IEaArrangeRoomService

func init() {
	route.AddGroup("/api/auth/ea_arrange_room", "考场编排详情").
		AddPOSTHandel("/get", getArrangeRoom, true, "查询考场编排详情详情").
		AddPOSTHandelPowerR("/list", listArrangeRoom, true, cst.SYS_ROLE_KEY_ADMISSIONS, "查询考场编排详情列表").
		AddPOSTHandel("/add", addArrangeRoom, true, "新增考场编排详情").
		AddPOSTHandelPowerR("/edit", editArrangeRoom, true, cst.SYS_ROLE_KEY_ADMISSIONS, "修改考场编排详情").
		AddPOSTHandel("/del", delArrangeRoom, true, "删除考场编排详情").
		AddPOSTHandel("/export", exportArrangeRoom, true, "导出考场编排详情")

	eaArrangeRoomService = bizSerImpl.NewEaArrangeRoomService()
}

func getArrangeRoom(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaArrangeRoomQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.RoomArrangeId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	arrangeRoomDB := eaArrangeRoomService.GetEaArrangeRoomById(q.RoomArrangeId)
	json.SetObj(arrangeRoomDB).Render()
}

func listArrangeRoom(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaArrangeRoomQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.InfoId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	list := eaArrangeRoomService.ListEaArrangeRoom(q)
	json.SetList(list).Render()
}

func addArrangeRoom(c *gin.Context) {
	json := web.NewJsonResult(c)
	//ctx := context.NewUserContext(c)
	//arrangeRoomAddDto := new(bizDto.EaArrangeRoomAddDto)
	//if err := c.ShouldBindJSON(arrangeRoomAddDto); err != nil {
	//	json.SetErrs(errs.PARAM_INVALID_ERR).Render()
	//	return
	//}
	//
	//success, errInfo := eaArrangeRoomService.BatchInsertArrangeRoom(arrangeRoomAddDto, ctx.GetUserId())
	//if !success {
	//	if errInfo.Invalid() {
	//		json.SetErrs(errInfo).Render()
	//		return
	//	}
	//	json.SetErrs(errs.HANDLE_ERR).Render()
	//	return
	//}
	json.Render()
}

func editArrangeRoom(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	arrangeRoomAddDto := new(bizDto.EaArrangeRoomAddDto)
	if err := c.ShouldBindJSON(arrangeRoomAddDto); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	success, errInfo := eaArrangeRoomService.BatchUpdateArrangeRoom(arrangeRoomAddDto, ctx.GetUserId())
	if !success {
		if errInfo.Invalid() {
			json.SetErrs(errInfo).Render()
			return
		}
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func delArrangeRoom(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaArrangeRoomQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.RoomArrangeId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := eaArrangeRoomService.DeleteEaArrangeRoomById(q.RoomArrangeId, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func exportArrangeRoom(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaArrangeRoomQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data, err := eaArrangeRoomService.ExportEaArrangeRoom(q)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, data)
}

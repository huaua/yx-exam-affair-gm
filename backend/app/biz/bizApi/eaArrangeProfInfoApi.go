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

var eaArrangeProfInfoService bizSer.IEaArrangeProfInfoService

func init() {
	route.AddGroup("/api/auth/ea_arrange_prof_info", "编排专业").
		AddPOSTHandelPowerR("/get", getArrangeProfInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "查询编排专业详情").
		AddPOSTHandelPowerR("/list", listArrangeProfInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "查询编排专业列表").
		AddPOSTHandelPowerR("/list_available_room", listAvailableRoomInTask, true, cst.SYS_ROLE_KEY_ADMISSIONS, "查询任务下所有可编排的教室列表").
		AddPOSTHandel("/add", addArrangeProfInfo, true, "新增编排专业").
		AddPOSTHandelPowerR("/cancel_arrange", editArrangeProfInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "取消编排").
		AddPOSTHandel("/del", delArrangeProfInfo, true, "删除编排专业").
		AddPOSTHandel("/export", exportArrangeProfInfo, true, "导出编排专业").
		AddPOSTHandelPowerR("/import", importArrangeProfInfo, true, cst.SYS_ROLE_KEY_ADMISSIONS, "导入编排专业")

	eaArrangeProfInfoService = bizSerImpl.NewEaArrangeProfInfoService()
}

func getArrangeProfInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaArrangeProfInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.InfoId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	arrangeProfInfoDB := eaArrangeProfInfoService.GetEaArrangeProfInfoById(q.InfoId)
	json.SetObj(arrangeProfInfoDB).Render()
}

func listArrangeProfInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaArrangeProfInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	p := q.GetPage()
	list := eaArrangeProfInfoService.ListEaArrangeProfInfo(q, p)
	json.SetPageList(list, p).Render()
}

func addArrangeProfInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	arrangeProfInfoDO := new(bizDo.EaArrangeProfInfo)
	if err := c.ShouldBindJSON(arrangeProfInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	arrangeProfInfoDO.SetCreateBy(ctx.GetUserId())
	arrangeProfInfoDO.SetUpdateBy(ctx.GetUserId())
	success := eaArrangeProfInfoService.InsertEaArrangeProfInfo(arrangeProfInfoDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func editArrangeProfInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	arrangeProfInfoDO := new(bizDo.EaArrangeProfInfo)
	if err := c.ShouldBindJSON(arrangeProfInfoDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if arrangeProfInfoDO.InfoId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	arrangeProfInfoDO.SetUpdateBy(ctx.GetUserId())
	success, errInfo := eaArrangeProfInfoService.CancelArrangeById(arrangeProfInfoDO)
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

func delArrangeProfInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaArrangeProfInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.InfoId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := eaArrangeProfInfoService.DeleteEaArrangeProfInfoById(q.InfoId, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func exportArrangeProfInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaArrangeProfInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data, err := eaArrangeProfInfoService.ExportEaArrangeProfInfo(q)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, data)
}

func importArrangeProfInfo(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaArrangeProfInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	ctx := context.NewUserContext(c)

	q.OperId = ctx.GetUserId()
	success, err := eaArrangeProfInfoService.ImportEaArrangeProfInfo(q)
	if !success {
		json.SetErr(err).Render()
		return
	}
	json.Render()
}

func listAvailableRoomInTask(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaArrangeProfInfoQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	if q.InfoId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	list := eaArrangeProfInfoService.ListAvailableRoomInTask(q)
	json.SetList(list).Render()
}

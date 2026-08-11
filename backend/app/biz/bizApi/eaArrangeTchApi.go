package bizApi

import (
	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/route"
	"github.com/sealsee/web-base/public/web"
	"yx-exam-affair-gm/app/biz/bizSer"
	"yx-exam-affair-gm/app/biz/bizSer/bizSerImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	"yx-exam-affair-gm/app/cst"
)

var eaArrangeTchService bizSer.IEaArrangeTchService

func init() {
	route.AddGroup("/api/auth/ea_arrange_tch", "监考老师分配").
		AddPOSTHandel("/get", getArrangeTch, true, "查询监考老师分配详情").
		AddPOSTHandel("/list", listArrangeTch, true, "查询监考老师分配列表-按考场").
		AddPOSTHandel("/list_by_prof", listArrangeTchByProf, true, "查询监考老师分配列表-按专业").
		AddPOSTHandel("/add", addArrangeTch, true, "新增监考老师分配").
		AddPOSTHandelPowerR("/edit", editArrangeTch, true, cst.SYS_ROLE_KEY_ADMISSIONS, "修改监考老师分配").
		AddPOSTHandel("/del", delArrangeTch, true, "取消监考老师分配").
		AddPOSTHandel("/del_single", delSingleArrangeTch, true, "取消单个监考老师分配").
		AddPOSTHandelPowerR("/export", exportArrangeTch, true, cst.SYS_ROLE_KEY_ADMISSIONS, "导出监考老师分配")

	eaArrangeTchService = bizSerImpl.NewEaArrangeTchService()
}

func getArrangeTch(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaArrangeTchQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.TchArrangeId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	arrangeTchDB := eaArrangeTchService.GetEaArrangeTchById(q.TchArrangeId)
	json.SetObj(arrangeTchDB).Render()
}

func listArrangeTch(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaArrangeTchQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	p := q.GetPage()
	list := eaArrangeTchService.ListEaArrangeTch(q, p)
	json.SetPageList(list, p).Render()
}

func listArrangeTchByProf(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaArrangeTchQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	p := q.GetPage()
	list := eaArrangeTchService.ListEaArrangeTchByProf(q, p)
	json.SetPageList(list, p).Render()
}

func addArrangeTch(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	arrangeTchDO := new(bizDo.EaArrangeTch)
	if err := c.ShouldBindJSON(arrangeTchDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	arrangeTchDO.SetCreateBy(ctx.GetUserId())
	arrangeTchDO.SetUpdateBy(ctx.GetUserId())
	success := eaArrangeTchService.InsertEaArrangeTch(arrangeTchDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func editArrangeTch(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	var saveDtoList []*bizDto.EaArrangeTchSaveDto
	if err := c.ShouldBindJSON(&saveDtoList); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if len(saveDtoList) < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success, errInfo := eaArrangeTchService.BatchUpdateEaArrangeTch(saveDtoList, ctx.GetUserId())
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

func delArrangeTch(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaArrangeTchQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.RoomReportId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := eaArrangeTchService.DelArrangeTchByRoomReportId(q.RoomReportId, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func delSingleArrangeTch(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	q := new(bizDto.EaArrangeTchQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.RoomReportId < 1 || q.TchReportId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := eaArrangeTchService.DelSingleArrangeTch(q.RoomReportId, q.TchReportId, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func exportArrangeTch(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(bizDto.EaArrangeTchQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data, err := eaArrangeTchService.ExportEaArrangeTch(q)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, data)
}

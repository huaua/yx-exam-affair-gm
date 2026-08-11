package userApi

import (
	"yx-exam-affair-gm/app/user/model/userDto"
	"yx-exam-affair-gm/app/user/userSer"
	"yx-exam-affair-gm/app/user/userSer/userSerImpl"

	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/route"
	"github.com/sealsee/web-base/public/web"
)

var sysDataInoutService userSer.ISysDataInoutService

func init() {
	route.AddGroup("/api/auth/sys_data_inout", "导入导出记录").
		AddPOSTHandel("/list", listDataInout, true, "查询导入导出记录列表")

	sysDataInoutService = userSerImpl.NewSysDataInoutService()
}

func listDataInout(c *gin.Context) {
	json := web.NewJsonResult(c)
	dataInoutQuery := new(userDto.SysDataInoutQuery)
	if err := c.ShouldBindJSON(dataInoutQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	page := dataInoutQuery.GetPage()
	list := sysDataInoutService.ListSysDataInout(dataInoutQuery, page)
	json.SetPageList(list, page).Render()
}

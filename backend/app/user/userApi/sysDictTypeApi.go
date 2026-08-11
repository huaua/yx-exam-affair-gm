package userApi

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"
	"yx-exam-affair-gm/app/user/userSer"
	"yx-exam-affair-gm/app/user/userSer/userSerImpl"

	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/route"
	"github.com/sealsee/web-base/public/web"
)

var sysDictTypeService userSer.ISysDictTypeService

func init() {
	route.AddGroup("/api/auth/sys_dict_type", "字典类型").
		AddGETHandel("/get", getDictType, true, "查询字典类型信息").
		AddPOSTHandel("/list", listDictType, true, "查询字典类型列表").
		AddGETHandel("/option_select", optionSelect, true, "获取字典选择框列表").
		AddPOSTHandel("/add", addDictType, true, "新增字典类型").
		AddPOSTHandel("/edit", editDictType, true, "修改字典类型").
		AddPOSTHandel("/refresh_cache", refreshDictCache, true, "刷新字典缓存").
		AddPOSTHandel("/del", delDictType, true, "删除字典类型").
		AddPOSTHandel("/export", exportDictType, true, "导出字典类型")

	sysDictTypeService = userSerImpl.NewSysDictTypeService()
}

func getDictType(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	id := ctx.QueryInt64("dictId")
	if id < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	dictTypeDB := sysDictTypeService.GetSysDictTypeById(id)
	json.SetObj(dictTypeDB).Render()
}

func listDictType(c *gin.Context) {
	json := web.NewJsonResult(c)
	dictTypeQuery := new(userDto.SysDictTypeQuery)
	if err := c.ShouldBindJSON(dictTypeQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	page := dictTypeQuery.GetPage()
	list := sysDictTypeService.ListSysDictType(dictTypeQuery, page)
	json.SetPageList(list, page).Render()
}

func optionSelect(c *gin.Context) {
	json := web.NewJsonResult(c)
	q := new(userDto.SysDictTypeQuery)
	q.PageSize = page.MAX_PAGE_SIZE
	list := sysDictTypeService.ListSysDictType(q, q.GetPage())
	json.SetList(list).Render()
}

func addDictType(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	dictTypeDO := new(userDo.SysDictType)
	if err := c.ShouldBindJSON(dictTypeDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	dictTypeDO.SetCreateBy(ctx.GetUserId())
	dictTypeDO.SetUpdateBy(ctx.GetUserId())
	success := sysDictTypeService.InsertSysDictType(dictTypeDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func editDictType(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	dictTypeDO := new(userDo.SysDictType)
	if err := c.ShouldBindJSON(dictTypeDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if dictTypeDO.DictId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	dictTypeDO.SetUpdateBy(ctx.GetUserId())
	success := sysDictTypeService.UpdateSysDictTypeById(dictTypeDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func refreshDictCache(c *gin.Context) {
	json := web.NewJsonResult(c)
	// TODO 刷新redis缓存： 1、清除所有字典缓存 2、重新加载所有字典缓存
	json.Render()
}

func delDictType(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	dictTypeQuery := new(userDto.SysDictTypeQuery)
	if err := c.ShouldBindJSON(dictTypeQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if dictTypeQuery.DictId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := sysDictTypeService.DeleteSysDictTypeById(dictTypeQuery.DictId, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func exportDictType(c *gin.Context) {
	json := web.NewJsonResult(c)
	dictTypeQuery := new(userDto.SysDictTypeQuery)
	if err := c.ShouldBind(dictTypeQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data := sysDictTypeService.ExportSysDictType(dictTypeQuery)
	web.DataPackageExcel(c, data)
}

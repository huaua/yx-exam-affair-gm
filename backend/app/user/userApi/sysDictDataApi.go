package userApi

import (
	"unicode/utf8"
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"
	"yx-exam-affair-gm/app/user/userSer"
	"yx-exam-affair-gm/app/user/userSer/userSerImpl"

	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/route"
	"github.com/sealsee/web-base/public/web"
)

var sysDictDataService userSer.ISysDictDataService

func init() {
	route.AddGroup("/api/auth/sys_dict_data", "字典数据").
		AddGETHandel("/get", getDictData, true, "查询字典数据信息").
		AddGETHandel("/get_dicts", getDictsByType, true, "根据类型查字典数据").
		AddPOSTHandel("/list", listDictData, true, "查询字典数据列表").
		AddPOSTHandel("/add", addDictData, true, "新增字典数据").
		AddPOSTHandel("/edit", editDictData, true, "修改字典数据").
		AddPOSTHandel("/del", delDictData, true, "删除字典数据").
		AddPOSTHandel("/export", exportDictData, true, "导出字典数据").
		AddPOSTHandel("/change_status", changeDictStatus, true, "修改字典数状态")

	sysDictDataService = userSerImpl.NewSysDictDataService()
}

func getDictData(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	id := ctx.QueryInt64("dictCode")
	if id < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	dictDataDB := sysDictDataService.GetSysDictDataById(id)
	json.SetObj(dictDataDB).Render()
}

func getDictsByType(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	dictType := ctx.Query("dictType")
	if dictType == "" {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	list := sysDictDataService.GetSysDictDataByType(dictType)
	json.SetList(list).Render()
}

func listDictData(c *gin.Context) {
	json := web.NewJsonResult(c)
	dictDataQuery := new(userDto.SysDictDataQuery)
	if err := c.ShouldBindJSON(dictDataQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	// 判断字典类型
	if len(dictDataQuery.DictType) <= 0 {
		json.SetErr("字典类型不可为空！").Render()
		return
	}
	page := dictDataQuery.GetPage()
	list := sysDictDataService.ListSysDictData(dictDataQuery, page)
	json.SetPageList(list, page).Render()
}

func addDictData(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	dictDataDO := new(userDo.SysDictData)
	if err := c.ShouldBindJSON(dictDataDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	// 基础校验
	if len(dictDataDO.DictLabel) > 64 || len(dictDataDO.DictValue) > 64 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	dictDataDO.SetCreateBy(ctx.GetUserId())
	dictDataDO.SetUpdateBy(ctx.GetUserId())
	success := sysDictDataService.InsertSysDictData(dictDataDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func editDictData(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	dictDataDO := new(userDo.SysDictData)
	if err := c.ShouldBindJSON(dictDataDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	// 基础校验
	if utf8.RuneCountInString(dictDataDO.DictLabel) > 64 || utf8.RuneCountInString(dictDataDO.DictValue) > 64 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if dictDataDO.DictCode < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	dictDataDO.SetUpdateBy(ctx.GetUserId())
	success := sysDictDataService.UpdateSysDictDataById(dictDataDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func delDictData(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	dictDataQuery := new(userDto.SysDictDataQuery)
	if err := c.ShouldBindJSON(dictDataQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if dictDataQuery.DictCode < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := sysDictDataService.DeleteSysDictDataById(dictDataQuery.DictCode, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func exportDictData(c *gin.Context) {
	json := web.NewJsonResult(c)
	dictDataQuery := new(userDto.SysDictDataQuery)
	if err := c.ShouldBind(dictDataQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data := sysDictDataService.ExportSysDictData(dictDataQuery)
	web.DataPackageExcel(c, data)
}

func changeDictStatus(c *gin.Context) {
	json := web.NewJsonResult(c)
	dictDataQuery := new(userDto.SysDictDataQuery)
	ctx := context.NewUserContext(c)
	if err := c.ShouldBind(dictDataQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	dictDataDO := new(userDo.SysDictData)
	dictDataDO.DictCode = dictDataQuery.DictCode
	dictDataDO.State = dictDataQuery.State
	dictDataDO.UpdateBy = ctx.GetUserId()
	success := sysDictDataService.UpdateDicData(dictDataDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

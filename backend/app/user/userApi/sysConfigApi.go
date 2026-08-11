package userApi

import (
	"unicode/utf8"
	"yx-exam-affair-gm/app/hint"
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

var sysConfigService userSer.ISysConfigService

func init() {
	route.AddGroup("/api/auth/sys_config", "参数配置").
		AddGETHandel("/get", getConfig, true, "查询参数配置详情").
		AddPOSTHandel("/list", listConfig, true, "查询参数配置列表").
		AddPOSTHandel("/add", addConfig, true, "新增参数配置").
		AddPOSTHandel("/edit", editConfig, true, "修改参数配置").
		AddPOSTHandel("/refresh_cache", refreshConfigCache, true, "刷新参数缓存").
		AddPOSTHandel("/del", delConfig, true, "删除参数配置").
		AddPOSTHandel("/export", exportConfig, true, "导出参数配置")

	sysConfigService = userSerImpl.NewSysConfigService()
}

func getConfig(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	configId := ctx.QueryInt64("configId")
	if configId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	configDB := sysConfigService.GetSysConfigById(configId)
	json.SetObj(configDB).Render()
}

func listConfig(c *gin.Context) {
	json := web.NewJsonResult(c)
	configQuery := new(userDto.SysConfigQuery)
	if err := c.ShouldBindJSON(configQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	page := configQuery.GetPage()
	list := sysConfigService.ListSysConfig(configQuery, page)
	returnList := make([]*userDo.SysConfig, 0)
	// 如果列表中单项的id小于100，则不显示
	for _, item := range list {
		if item.ConfigId < 100 {
			continue
		}
		returnList = append(returnList, item)
	}
	json.SetPageList(returnList, page).Render()
}

func addConfig(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	configDO := new(userDo.SysConfig)
	if err := c.ShouldBindJSON(configDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	// 长度校验
	if utf8.RuneCountInString(configDO.ConfigName) > 100 || utf8.RuneCountInString(configDO.ConfigValue) > 500 || utf8.RuneCountInString(configDO.Remark) > 500 {
		json.SetErrs(hint.SYS_LENGTH_ERR).Render()
		return
	}

	configDO.ConfigType = "N"
	configDO.SetCreateBy(ctx.GetUserId())
	configDO.SetUpdateBy(ctx.GetUserId())
	success := sysConfigService.InsertSysConfig(configDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func editConfig(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	configDO := new(userDo.SysConfig)
	if err := c.ShouldBindJSON(configDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if configDO.ConfigId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}

	// 长度校验
	if utf8.RuneCountInString(configDO.ConfigName) > 100 || utf8.RuneCountInString(configDO.ConfigValue) > 500 || utf8.RuneCountInString(configDO.Remark) > 500 {
		json.SetErrs(hint.SYS_LENGTH_ERR).Render()
		return
	}
	configDO.SetUpdateBy(ctx.GetUserId())
	success := sysConfigService.UpdateSysConfigById(configDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func refreshConfigCache(c *gin.Context) {
	json := web.NewJsonResult(c)
	// TODO 刷新redis缓存： 1、清除所有参数缓存 2、重新加载所有参数缓存
	json.Render()
}

func delConfig(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	configQuery := new(userDto.SysConfigQuery)
	if err := c.ShouldBindJSON(configQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if configQuery.ConfigId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := sysConfigService.DeleteSysConfigById(configQuery.ConfigId, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func exportConfig(c *gin.Context) {
	json := web.NewJsonResult(c)
	configQuery := new(userDto.SysConfigQuery)
	if err := c.ShouldBind(configQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	// data := sysConfigService.ExportSysConfig(configQuery)
	data, err := sysConfigService.ExportSysConfigNew(configQuery)
	if err != nil {
		json.SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, data)
}

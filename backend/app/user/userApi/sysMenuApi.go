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

var sysMenuService userSer.ISysMenuService

func init() {
	route.AddGroup("/api/auth/sys_menu", "菜单权限").
		AddGETHandel("/get", getMenu, true, "查询菜单权限信息").
		AddGETHandel("/treeselect", getTreeselect, true, "获取菜单下拉树").
		AddPOSTHandel("/list", listMenu, true, "查询菜单权限列表").
		AddPOSTHandel("/add", addMenu, true, "新增菜单权限").
		AddPOSTHandel("/edit", editMenu, true, "修改菜单权限").
		AddPOSTHandel("/del", delMenu, true, "删除菜单权限").
		AddGETHandel("/export", exportMenu, true, "导出菜单权限")

	sysMenuService = userSerImpl.NewSysMenuService()
}

func getMenu(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	id := ctx.QueryInt64("menuId")
	if id < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	menuDB := sysMenuService.GetSysMenuById(id)
	json.SetObj(menuDB).Render()
}

func getTreeselect(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	// 查询菜单列表
	menuList := sysMenuService.ListSysMenuByUserId(ctx.GetUserId())
	treeselect := sysMenuService.BuildMenuTreeSelect(menuList)
	json.AddData("treeselect", treeselect).
		Render()
}

func listMenu(c *gin.Context) {
	json := web.NewJsonResult(c)
	menuQuery := new(userDto.SysMenuQuery)
	if err := c.ShouldBindJSON(menuQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	p := menuQuery.GetPage()
	p.PageSize = page.MAX_PAGE_SIZE
	list := sysMenuService.ListSysMenu(menuQuery, p)
	json.SetPageList(list, p).Render()
}

func addMenu(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	menuDO := new(userDo.SysMenu)
	if err := c.ShouldBindJSON(menuDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	menuDO.SetCreateBy(ctx.GetUserId())
	menuDO.SetUpdateBy(ctx.GetUserId())
	success := sysMenuService.InsertSysMenu(menuDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func editMenu(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	menuDO := new(userDo.SysMenu)
	if err := c.ShouldBindJSON(menuDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if menuDO.MenuId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	menuDO.SetUpdateBy(ctx.GetUserId())
	success := sysMenuService.UpdateSysMenuById(menuDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func delMenu(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	menuQuery := new(userDto.SysMenuQuery)
	if err := c.ShouldBindJSON(menuQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if menuQuery.MenuId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := sysMenuService.DeleteSysMenuById(menuQuery.MenuId, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func exportMenu(c *gin.Context) {
	json := web.NewJsonResult(c)
	menuQuery := new(userDto.SysMenuQuery)
	if err := c.ShouldBindJSON(menuQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	data := sysMenuService.ExportSysMenu(menuQuery)
	web.DataPackageExcel(c, data)
}

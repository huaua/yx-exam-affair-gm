package userApi

import (
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

var sysRoleService userSer.ISysRoleService

func init() {
	route.AddGroup("/api/auth/sys_role", "角色").
		AddGETHandel("/get", getRole, true, "查询角色详情").
		AddGETHandel("/get_with_menutree", getRoleWithMenuTree, true, "获取角色菜单详情").
		AddPOSTHandel("/list", listRole, true, "查询角色列表").
		AddPOSTHandel("/add", addRole, true, "新增角色").
		AddPOSTHandel("/edit", editRole, true, "修改角色").
		AddPOSTHandel("/change_state", changeRoleState, true, "修改状态").
		AddPOSTHandel("/del", delRole, true, "删除角色")

	sysRoleService = userSerImpl.NewSysRoleService()
}

func getRole(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	id := ctx.QueryInt64("roleId")
	if id < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	roleDB := sysRoleService.GetSysRoleById(id)
	json.SetObj(roleDB).Render()
}

func getRoleWithMenuTree(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	id := ctx.QueryInt64("roleId")
	if id < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	roleDB := sysRoleService.GetSysRoleById(id)
	// 查询菜单列表
	menuList := sysMenuService.ListSysMenuByUserId(ctx.GetUserId())
	treeselect := sysMenuService.BuildMenuTreeSelect(menuList)
	// 查询角色所拥有的菜单项
	checkedKeys := sysMenuService.ListMenuIdByRoleId(id, *roleDB.MenuCheckStrictly)

	json.AddData("role", roleDB).
		AddData("menus", treeselect).
		AddData("checkedKeys", checkedKeys).
		Render()
}

func listRole(c *gin.Context) {
	json := web.NewJsonResult(c)
	roleQuery := new(userDto.SysRoleQuery)
	if err := c.ShouldBindJSON(roleQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	page := roleQuery.GetPage()
	list := sysRoleService.ListSysRole(roleQuery, page)
	json.SetPageList(list, page).Render()
}

func addRole(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	roleDO := new(userDo.SysRole)
	if err := c.ShouldBindJSON(roleDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	roleDO.SetCreateBy(ctx.GetUserId())
	roleDO.SetUpdateBy(ctx.GetUserId())
	success := sysRoleService.InsertSysRole(roleDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func editRole(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	roleDO := new(userDo.SysRole)
	if err := c.ShouldBindJSON(roleDO); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if roleDO.RoleId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	roleDO.SetUpdateBy(ctx.GetUserId())
	success := sysRoleService.UpdateSysRoleAndMenu(roleDO)
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func changeRoleState(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	role := new(userDo.SysRole)
	if err := c.ShouldBindJSON(role); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if role.RoleId < 1 || role.State == "" {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := sysRoleService.UpdateRoleState(role.RoleId, role.State, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

func delRole(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	roleQuery := new(userDto.SysRoleQuery)
	if err := c.ShouldBindJSON(roleQuery); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if roleQuery.RoleId < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	success := sysRoleService.DeleteSysRoleById(roleQuery.RoleId, ctx.GetUserId())
	if !success {
		json.SetErrs(errs.HANDLE_ERR).Render()
		return
	}
	json.Render()
}

package userDto

import "github.com/sealsee/web-base/public/basemodel"

type SysRoleMenuQuery struct {
	RoleId int64 `form:"roleId" json:"roleId,string,omitempty"` // 角色ID
	MenuId int64 `form:"menuId" json:"menuId,string,omitempty"` // 菜单ID
	basemodel.BaseEntityQuery
}

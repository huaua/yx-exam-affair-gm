package userDo

type SysRoleMenu struct {
	RoleId int64 `json:"roleId,string,omitempty"` // 角色ID
	MenuId int64 `json:"menuId,string,omitempty"` // 菜单ID
}

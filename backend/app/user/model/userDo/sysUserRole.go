package userDo

type SysUserRole struct {
	UserId int64 `json:"userId,string,omitempty"` // 用户ID
	RoleId int64 `json:"roleId,string,omitempty"` // 角色ID
}

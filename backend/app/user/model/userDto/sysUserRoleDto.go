package userDto

import "github.com/sealsee/web-base/public/basemodel"

type SysUserRoleQuery struct {
	UserId int64 `form:"userId" json:"userId,string,omitempty"` // 用户ID
	RoleId int64 `form:"roleId" json:"roleId,string,omitempty"` // 角色ID
	basemodel.BaseEntityQuery
}

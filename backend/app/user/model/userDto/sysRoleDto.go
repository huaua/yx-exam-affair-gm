package userDto

import "github.com/sealsee/web-base/public/basemodel"

type SysRoleQuery struct {
	RoleId            int64  `gorm:"primaryKey" form:"roleId" json:"roleId,string,omitempty"` // 主键 角色ID
	RoleName          string `form:"roleName" json:"roleName,omitempty"`                      // 角色名称
	RoleKey           string `form:"roleKey" json:"roleKey,omitempty"`                        // 角色权限字符串
	RoleSort          int    `form:"roleSort" json:"roleSort,omitempty"`                      // 显示顺序
	DataScope         string `form:"dataScope" json:"dataScope,omitempty"`                    // 数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）
	MenuCheckStrictly int    `form:"menuCheckStrictly" json:"menuCheckStrictly,omitempty"`    // 菜单树选择项是否关联显示
	DeptCheckStrictly int    `form:"deptCheckStrictly" json:"deptCheckStrictly,omitempty"`    // 部门树选择项是否关联显示
	State             string `form:"state" json:"state,omitempty"`                            // 角色状态（1正常 2停用）
	Remark            string `form:"remark" json:"remark,omitempty"`                          // 备注
	basemodel.BaseEntityQuery
}

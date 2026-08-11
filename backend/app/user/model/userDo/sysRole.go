package userDo

import "github.com/sealsee/web-base/public/basemodel"

type SysRole struct {
	RoleId            int64   `gorm:"primaryKey" json:"roleId,string,omitempty"` // 主键 角色ID
	RoleName          string  `json:"roleName,omitempty"`                        // 角色名称
	RoleKey           string  `json:"roleKey,omitempty"`                         // 角色权限字符串
	RoleSort          int     `json:"roleSort,omitempty"`                        // 显示顺序
	DataScope         string  `json:"dataScope,omitempty"`                       // 数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）
	MenuCheckStrictly *bool   `json:"menuCheckStrictly,omitempty"`               // 菜单树选择项是否关联显示
	State             string  `json:"state,omitempty"`                           // 角色状态（1正常 2停用）
	Remark            string  `json:"remark,omitempty"`                          // 备注
	MenuIds           []int64 `gorm:"-" json:"menuIds"`                          // 菜单Ids
	basemodel.BaseEntity
}

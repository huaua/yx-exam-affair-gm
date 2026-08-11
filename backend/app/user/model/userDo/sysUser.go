package userDo

import (
	"github.com/sealsee/web-base/public/basemodel"
)

type SysUser struct {
	UserId      int64   `gorm:"primaryKey" json:"userId,string,omitempty"` // 主键 用户ID
	DeptId      int64   `json:"deptId,string,omitempty"`                   // 部门ID
	UserName    string  `json:"userName,omitempty"`                        // 用户账号
	NickName    string  `json:"nickName,omitempty"`                        // 用户昵称
	UserType    string  `json:"userType,omitempty"`                        // 用户类型（00系统用户）
	Email       string  `json:"email,omitempty"`                           // 用户邮箱
	Phonenumber string  `json:"phonenumber,omitempty"`                     // 手机号码
	Sex         string  `json:"sex,omitempty"`                             // 用户性别（男 女 未知）
	Avatar      string  `json:"avatar,omitempty"`                          // 头像地址
	Password    string  `json:"password,omitempty"`                        // 密码
	State       string  `json:"state,omitempty"`                           // 帐号状态（1正常 0停用）
	Remark      string  `json:"remark,omitempty"`                          // 备注
	RoleIds     []int64 `gorm:"-" json:"roleIds"`                          // 角色Ids
	RoleKey     string  `gorm:"-" json:"roleKey"`                          // 角色代号
	DeptName    string  `gorm:"-" json:"deptName"`                         // 学院名称
	basemodel.BaseEntity
}

package userDto

import (
	"github.com/sealsee/web-base/public/basemodel"
)

type SysUserQuery struct {
	UserId      int64  `gorm:"primaryKey" json:"userId,string,omitempty"` // 主键 用户ID
	DeptId      int    `form:"deptId" json:"deptId,omitempty"`            // 部门ID
	UserName    string `form:"userName" json:"userName,omitempty"`        // 用户账号
	NickName    string `form:"nickName" json:"nickName,omitempty"`        // 用户昵称
	UserType    string `form:"userType" json:"userType,omitempty"`        // 用户类型（00系统用户）
	Email       string `form:"email" json:"email,omitempty"`              // 用户邮箱
	Phonenumber string `form:"phonenumber" json:"phonenumber,omitempty"`  // 手机号码
	Sex         string `form:"sex" json:"sex,omitempty"`                  // 用户性别（男 女 未知）
	State       string `form:"state" json:"state,omitempty"`              // 帐号状态（1正常 0停用）
	Remark      string `form:"remark" json:"remark,omitempty"`            // 备注
	OldPassword string `gorm:"-" json:"oldPassword,omitempty"`            // 修改密码-原密码
	NewPassword string `gorm:"-" json:"newPassword,omitempty"`            // 修改密码-新密码
	basemodel.BaseEntityQuery
}

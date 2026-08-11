package userDto

import "github.com/sealsee/web-base/public/basemodel"

type SysDictTypeQuery struct {
	DictId   int64  `gorm:"primaryKey" json:"dictId,string,omitempty"` // 主键 字典主键
	DictName string `form:"dictName" json:"dictName,omitempty"`        // 字典名称
	DictType string `form:"dictType" json:"dictType,omitempty"`        // 字典类型
	State    string `form:"state" json:"state,omitempty"`              // 状态（1正常 0停用）
	Remark   string `form:"remark" json:"remark,omitempty"`            // 备注
	basemodel.BaseEntityQuery
}

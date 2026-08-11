package userDo

import "github.com/sealsee/web-base/public/basemodel"

type SysDictType struct {
	DictId   int64  `gorm:"primaryKey" json:"dictId,string,omitempty"` // 主键 字典主键
	DictName string `json:"dictName,omitempty"`                        // 字典名称
	DictType string `json:"dictType,omitempty"`                        // 字典类型
	State    string `json:"state,omitempty"`                           // 状态（1-正常 0-停用）
	Remark   string `json:"remark,omitempty"`                          // 备注
	basemodel.BaseEntity
}

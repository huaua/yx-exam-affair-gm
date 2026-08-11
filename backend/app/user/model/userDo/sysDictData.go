package userDo

import (
	"github.com/sealsee/web-base/public/basemodel"
)

type SysDictData struct {
	DictCode  int64  `gorm:"primaryKey" json:"dictCode,string,omitempty"` // 主键 字典编码
	DictSort  int    `json:"dictSort,omitempty"`                          // 字典排序
	DictLabel string `json:"dictLabel,omitempty"`                         // 字典标签
	DictValue string `json:"dictValue,omitempty"`                         // 字典键值
	DictType  string `json:"dictType,omitempty"`                          // 字典类型
	CssClass  string `json:"cssClass,omitempty"`                          // 样式属性（其他样式扩展）
	ListClass string `json:"listClass,omitempty"`                         // 表格回显样式
	IsDefault string `json:"isDefault,omitempty"`                         // 是否默认（Y是 N否）
	State     string `json:"state,omitempty"`                             // 状态（1正常 0停用）
	Remark    string `json:"remark,omitempty"`                            // 备注
	basemodel.BaseEntity
}

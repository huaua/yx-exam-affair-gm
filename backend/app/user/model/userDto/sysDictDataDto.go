package userDto

import "github.com/sealsee/web-base/public/basemodel"

type SysDictDataQuery struct {
	DictCode  int64  `gorm:"primaryKey" json:"dictCode,string,omitempty"` // 主键 字典编码
	DictSort  int    `form:"dictSort" json:"dictSort,omitempty"`          // 字典排序
	DictLabel string `form:"dictLabel" json:"dictLabel,omitempty"`        // 字典标签
	DictValue string `form:"dictValue" json:"dictValue,omitempty"`        // 字典键值
	DictType  string `form:"dictType" json:"dictType,omitempty"`          // 字典类型
	CssClass  string `form:"cssClass" json:"cssClass,omitempty"`          // 样式属性（其他样式扩展）
	ListClass string `form:"listClass" json:"listClass,omitempty"`        // 表格回显样式
	IsDefault string `form:"isDefault" json:"isDefault,omitempty"`        // 是否默认（Y是 N否）
	State     string `form:"state" json:"state,omitempty"`                // 状态（1正常 0停用）
	Remark    string `form:"remark" json:"remark,omitempty"`              // 备注
	basemodel.BaseEntityQuery
}

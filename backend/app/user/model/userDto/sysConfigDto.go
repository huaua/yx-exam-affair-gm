package userDto

import (
	"github.com/sealsee/web-base/public/basemodel"
)

type SysConfigQuery struct {
	ConfigId    int64  `gorm:"primaryKey" form:"configId" json:"configId,string,omitempty"` // 主键 参数主键
	ConfigName  string `form:"configName" json:"configName,omitempty"`                      // 参数名称
	ConfigKey   string `form:"configKey" json:"configKey,omitempty"`                        // 参数键名
	ConfigValue string `form:"configValue" json:"configValue,omitempty"`                    // 参数键值
	ConfigType  string `form:"configType" json:"configType,omitempty"`                      // 系统内置（Y是 N否）
	Remark      string `form:"remark" json:"remark,omitempty"`                              // 备注
	basemodel.BaseEntityQuery
}

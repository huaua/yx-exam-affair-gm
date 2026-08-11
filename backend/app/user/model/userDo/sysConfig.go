package userDo

import "github.com/sealsee/web-base/public/basemodel"

type SysConfig struct {
	ConfigId    int64  `gorm:"primaryKey" json:"configId,string,omitempty"` // 主键 参数主键
	ConfigName  string `json:"configName,omitempty"`                        // 参数名称
	ConfigKey   string `json:"configKey,omitempty"`                         // 参数键名
	ConfigValue string `json:"configValue,omitempty"`                       // 参数键值
	ConfigType  string `json:"configType,omitempty"`                        // 系统内置（Y是 N否）
	Remark      string `json:"remark,omitempty"`                            // 备注
	basemodel.BaseEntity
}

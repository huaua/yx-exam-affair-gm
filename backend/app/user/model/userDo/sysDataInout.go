package userDo

import (
	"time"

	"github.com/sealsee/web-base/public/basemodel"
)

type SysDataInout struct {
	Id         int64     `gorm:"primaryKey" json:"id,string,omitempty"` // 主键 主键
	OperType   int       `json:"operType,omitempty"`                    // 操作类型 1-导入 2-更新 3-导出
	BizType    int       `json:"bizType,omitempty"`                     // 业务类型 1-高考改革数据 2-非高考改革数据 3-中专数据 4-专升本数据 5-第二学位数据
	BatchNo    string    `json:"batchNo,omitempty"`                     // 操作批次号
	Param      string    `json:"param,omitempty"`                       // 参数信息
	FileIn     string    `json:"fileIn,omitempty"`                      // 上传文件地址，多个压缩成zip文件
	FileInName string    `json:"fileInName,omitempty"`                  // 上传文件名称
	FileOut    string    `json:"fileOut,omitempty"`                     // 下载文件地址，多个压缩zip文件
	OperTime   time.Time `json:"operTime,omitempty"`                    // 操作时间
	State      string    `json:"state,omitempty"`                       // 状态 1-成功 2-失败
	Remark     string    `json:"remark,omitempty"`                      // 备注
	basemodel.BaseEntity
}

package bizDto

import (
	"github.com/sealsee/web-base/public/basemodel"
)

/*
 * desc: 编排专业数据查询实体
 * author: init code
 * date: 2024/09/24
 */
type EaArrangeProfInfoQuery struct {
	InfoId         int64              `gorm:"primaryKey" form:"infoId" json:"infoId,string,omitempty"` // 主键 数据ID
	RoomTaskId     int64              `form:"roomTaskId" json:"roomTaskId,string,omitempty"`           // 考场任务ID
	TaskDay        basemodel.BaseTime `form:"taskDay" json:"taskDay,string,omitempty"`                 // 征集日期
	DeptName       string             `form:"deptName" json:"deptName,omitempty"`                      // 学院名称
	ProfName       string             `form:"profName" json:"profName,omitempty"`                      // 专业名称
	DirectionName  string             `form:"directionName" json:"directionName,omitempty"`            // 方向名称
	ArrangeType    int                `form:"arrangeType" json:"arrangeType,string,omitempty"`         // 安排类型: 1-按组 2-按位
	StuNum         int                `form:"stuNum" json:"stuNum,string,omitempty"`                   // 考生数量
	ArrangeStuNum  int                `form:"arrangeStuNum" json:"arrangeStuNum,string,omitempty"`     // 已编排考生数量
	ArrangeRoomNum int                `form:"arrangeRoomNum" json:"arrangeRoomNum,string,omitempty"`   // 已编排考场数量
	ArrangeState   string             `form:"arrangeState" json:"arrangeState,omitempty"`              // 考场编排状态: 1-待编排 2-已编排
	Remark         string             `form:"remark" json:"remark,omitempty"`                          // 备注
	Files          map[string]string  `gorm:"-" form:"files" json:"files,omitempty"`                   // 导入文件 map[文件名]地址
	OperId         int64              `gorm:"-" form:"operId" json:"operId,omitempty"`                 // 操作人userid
	BizType        int                `gorm:"-" form:"bizType" json:"bizType,omitempty"`               //
	basemodel.BaseEntityQuery
}

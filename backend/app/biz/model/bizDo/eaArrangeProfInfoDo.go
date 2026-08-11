package bizDo

import (
	"github.com/sealsee/web-base/public/basemodel"
)

/*
 * desc: 编排专业数据实体
 * author: init code
 * date: 2024/09/24
 */
type EaArrangeProfInfo struct {
	InfoId              int64              `gorm:"primaryKey" json:"infoId,string,omitempty"` // 主键 数据ID
	RoomTaskId          int64              `json:"roomTaskId,string,omitempty"`               // 考场任务ID
	TaskDay             basemodel.BaseTime `json:"taskDay,omitempty"`                         // 征集日期
	DeptName            string             `json:"deptName,omitempty"`                        // 学院名称
	ProfName            string             `json:"profName,omitempty"`                        // 专业名称
	DirectionName       string             `json:"directionName,omitempty"`                   // 方向名称
	ArrangeType         int                `json:"arrangeType,string,omitempty"`              // 安排类型: 1-按组 2-按位
	StuNum              int                `json:"stuNum,string,omitempty"`                   // 考生数量
	ArrangeStuNum       int                `json:"arrangeStuNum,string,omitempty"`            // 已编排考生数量
	ArrangeRoomNum      int                `json:"arrangeRoomNum,string,omitempty"`           // 已编排考场数量
	ArrangeState        string             `json:"arrangeState,omitempty"`                    // 考场编排状态: 1-待编排 2-已编排
	Remark              string             `json:"remark,omitempty"`                          // 备注
	MaxArrangeRoomNoStr string             `gorm:"-" json:"maxArrangeRoomNoStr,omitempty"`    // 最大编排考场号
	MaxArrangeRoomNo    int                `gorm:"-" json:"maxArrangeRoomNo,omitempty"`       // 最大编排考场号
	RoomTaskName        string             `gorm:"-" json:"roomTaskName,omitempty"`           // 考场任务名称
	basemodel.BaseEntity
}

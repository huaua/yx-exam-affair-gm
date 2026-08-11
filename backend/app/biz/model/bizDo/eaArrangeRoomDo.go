package bizDo

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 考场编排详情数据实体
 * author: init code
 * date: 2024/09/24
 */
type EaArrangeRoom struct {
	RoomArrangeId   int64  `gorm:"primaryKey" json:"roomArrangeId,string,omitempty"`      // 主键 考场安排ID
	InfoId          int64  `json:"infoId,string,omitempty"`                               // 数据ID
	ArrangeType     int    `json:"arrangeType,string,omitempty"`                          // 安排类型: 1-按组 2-按位
	RoomReportId    int64  `json:"roomReportId,string,omitempty"`                         // 考场上报ID
	RoomTaskId      int64  `json:"roomTaskId,string,omitempty"`                           // 任务ID
	RoomId          int64  `json:"roomId,string,omitempty"`                               // 教室ID
	RoomName        string `json:"roomName,omitempty"`                                    // 教室名称
	StuNum          int    `json:"stuNum,string,omitempty"`                               // 考生数量
	RoomCapacity    int    `json:"roomCapacity,string,omitempty"`                         // 考场实际容量
	ArrangeNoStr    string `json:"arrangeNoStr,omitempty"`                                // 考场安排编号(字符串)
	ArrangeNo       int    `json:"arrangeNo,string,omitempty"`                            // 考场安排编号(数值)
	NeedTchNum      int    `json:"needTchNum,string,omitempty"`                           // 需要监考老师数量
	AssignmentState string `json:"assignmentState,omitempty"`                             // 分配状态 1-未分配 2-已分配
	Remark          string `json:"remark,omitempty"`                                      // 备注
	GroupNum        int    `gorm:"-" json:"groupNum,string,omitempty"`                    // 组数
	GroupCapacity   int    `gorm:"-" json:"groupCapacity,string,omitempty"`               // 按组容量
	Capacity        int    `gorm:"-" json:"capacity,string,omitempty"`                    // 按位容量
	DeptName        string `gorm:"-" form:"deptName" json:"deptName,omitempty"`           // 学院名称
	ProfName        string `gorm:"-" form:"profName" json:"profName,omitempty"`           // 专业名称
	DirectionName   string `gorm:"-" form:"directionName" json:"directionName,omitempty"` // 方向名称
	basemodel.BaseEntity
}

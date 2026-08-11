package bizDto

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 考场编排详情数据查询实体
 * author: init code
 * date: 2024/09/24
 */
type EaArrangeRoomQuery struct {
	RoomArrangeId   int64  `gorm:"primaryKey" form:"roomArrangeId" json:"roomArrangeId,string,omitempty"` // 主键 考场安排ID
	InfoId          int64  `form:"infoId" json:"infoId,string,omitempty"`                                 // 数据ID
	ArrangeType     int    `json:"arrangeType,string,omitempty"`                                          // 安排类型: 1-按组 2-按位
	RoomReportId    int64  `form:"roomReportId" json:"roomReportId,string,omitempty"`                     // 考场上报ID
	RoomTaskId      int64  `form:"roomTaskId" json:"roomTaskId,string,omitempty"`                         // 任务ID
	RoomId          int64  `form:"roomId" json:"roomId,string,omitempty"`                                 // 教室ID
	RoomName        string `form:"roomName" json:"roomName,omitempty"`                                    // 教室名称
	StuNum          int    `form:"stuNum" json:"stuNum,string,omitempty"`                                 // 考生数量
	RoomCapacity    int    `form:"roomCapacity" json:"roomCapacity,string,omitempty"`                     // 考场实际容量
	ArrangeNoStr    string `form:"arrangeNoStr" json:"arrangeNoStr,omitempty"`                            // 考场安排编号(字符串)
	ArrangeNo       int    `form:"arrangeNo" json:"arrangeNo,string,omitempty"`                           // 考场安排编号(数值)
	NeedTchNum      int    `form:"needTchNum" json:"needTchNum,string,omitempty"`                         // 需要监考老师数量
	AssignmentState string `form:"assignmentState" json:"assignmentState,omitempty"`                      // 分配状态 1-未分配 2-已分配
	Remark          string `form:"remark" json:"remark,omitempty"`                                        // 备注
	DeptId          int64  `gorm:"-" form:"deptId" json:"deptId,string,omitempty"`                        // 所属学院id
	ProfName        string `gorm:"-" form:"profName" json:"profName,omitempty"`                           // 专业名称
	DirectionName   string `gorm:"-" form:"directionName" json:"directionName,omitempty"`                 // 方向名称
	basemodel.BaseEntityQuery
}

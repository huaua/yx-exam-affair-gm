package bizDto

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 监考老师分配数据查询实体
 * author: init code
 * date: 2024/09/24
 */
type EaArrangeTchQuery struct {
	TchArrangeId    int64  `gorm:"primaryKey" form:"tchArrangeId" json:"tchArrangeId,string,omitempty"` // 主键 老师安排ID
	RoomReportId    int64  `form:"roomReportId" json:"roomReportId,string,omitempty"`                   // 考场安排ID
	TchReportId     int64  `form:"tchReportId" json:"tchReportId,string,omitempty"`                     // 教职工上报ID
	TchTaskId       int64  `form:"tchTaskId" json:"tchTaskId,string,omitempty"`                         // 教师征集任务ID
	TchId           int64  `form:"tchId" json:"tchId,string,omitempty"`                                 // 教职工ID
	TchName         string `form:"tchName" json:"tchName,omitempty"`                                    // 教职工姓名（冗余）
	Remark          string `form:"remark" json:"remark,omitempty"`                                      // 备注
	RoomTaskId      int64  `form:"roomTaskId" json:"roomTaskId,string,omitempty"`                       // 任务ID
	RoomArrangeId   int64  `form:"roomArrangeId" json:"roomArrangeId,string,omitempty"`                 // 考场安排ID
	AssignmentState string `form:"assignmentState" json:"assignmentState,omitempty"`                    // 分配状态 1-未分配 2-已分配
	RoomName        string `form:"roomName" json:"roomName,omitempty"`                                  // 教室名称
	DeptName        string `form:"deptName" json:"deptName,omitempty"`                                  // 所属学院
	ProfName        string `form:"profName" json:"profName,omitempty"`                                  // 专业名称
	DirectionName   string `form:"directionName" json:"directionName,omitempty"`                        // 方向名称
	basemodel.BaseEntityQuery
}

package bizDto

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 考场监考分配结果查询实体-按考场分组
 * author: liyu
 * date: 2024/09/29
 */
type EaArrangeTchResultDto struct {
	RoomArrangeId      int64    `form:"roomArrangeId" json:"roomArrangeId,string,omitempty"`    // 考场安排ID
	RoomReportId       int64    `form:"roomReportId" json:"roomReportId,string,omitempty"`      // 考场上报ID
	RoomId             int64    `form:"roomId" json:"roomId,string,omitempty"`                  // 教室ID
	RoomName           string   `form:"roomName" json:"roomName,omitempty"`                     // 教室名称
	ProfDirection      string   `form:"profDirection" json:"profDirection,omitempty"`           // 专业方向
	TaskDay            string   `form:"taskDay" json:"taskDay,omitempty"`                       // 考场日期
	TchNum             int      `form:"tchNum" json:"tchNum,string,omitempty"`                  // 需要的监考老师数量
	ArrangeTchReportId []string `form:"arrangeTchReportId" json:"arrangeTchReportId,omitempty"` // 编排的上报老师id
	ArrangeTchName     []string `form:"arrangeTchName" json:"arrangeTchName,omitempty"`         // 编排的老师名称
	ArrangeNoStr       string   `form:"arrangeNoStr" json:"arrangeNoStr,omitempty"`             // 考场安排编号(字符串)
	AssignmentState    string   `form:"assignmentState" json:"assignmentState,omitempty"`       // 分配状态 1-未分配 2-已分配
}

type EaArrangeTchDto struct {
	TchReportId int64  `form:"tchReportId" json:"tchReportId,string,omitempty"` // 教职工上报ID
	TchName     string `form:"tchName" json:"tchName,omitempty"`                // 教职工姓名
}

/*
 * desc: 考场监考分配结果查询实体-按专业分组
 * author: liyu
 * date: 2024/09/30
 */
type EaArrangeProfResultDto struct {
	InfoId          int64              `form:"infoId" json:"infoId,string,omitempty"`            // 专业数据ID
	DeptId          int64              `form:"deptId" json:"deptId,string,omitempty"`            // 所属学院id
	DeptName        string             `form:"deptName" json:"deptName,omitempty"`               // 学院名称
	ProfName        string             `form:"profName" json:"profName,omitempty"`               // 专业名称
	DirectionName   string             `form:"directionName" json:"directionName,omitempty"`     // 方向名称
	ArrangeType     int                `form:"arrangeType" json:"arrangeType,string,omitempty"`  // 安排类型: 1-按组 2-按位
	RoomName        string             `form:"roomName" json:"roomName,omitempty"`               // 教室名称
	GroupNum        int                `json:"groupNum,string,omitempty"`                        // 组数
	StuNum          int                `form:"stuNum" json:"stuNum,string,omitempty"`            // 考生数量
	ArrangeNoStr    string             `form:"arrangeNoStr" json:"arrangeNoStr,omitempty"`       // 考场安排编号(字符串)
	TchNum          int                `form:"tchNum" json:"tchNum,string,omitempty"`            // 需要的监考老师数量
	ArrangeTchName  string             `form:"arrangeTchName" json:"arrangeTchName,omitempty"`   // 编排的老师名称
	AssignmentState string             `form:"assignmentState" json:"assignmentState,omitempty"` // 分配状态 1-未分配 2-已分配
	DeptCode        string             `form:"deptCode" json:"deptCode,omitempty"`               // 学院代码
	TaskDay         basemodel.BaseTime `form:"taskDay" json:"taskDay,omitempty"`                 // 征集日期
}

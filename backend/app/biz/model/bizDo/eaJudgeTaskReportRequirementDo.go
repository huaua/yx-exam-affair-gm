package bizDo

import "github.com/sealsee/web-base/public/basemodel"

// EaJudgeTaskReportRequirement 学院按评分方向上报专家的数量要求。
// 覆盖 Deleted 为 int64：软删用 Unix 时间戳，避免唯一键冲突（参见 update_260803_judge_task_report_requirement_deleted.sql）。
type EaJudgeTaskReportRequirement struct {
	RequirementId int64  `gorm:"primaryKey" json:"requirementId,string,omitempty"`
	JudgeTaskId   int64  `json:"judgeTaskId,string,omitempty"`
	DirectionCode string `json:"directionCode,omitempty"`
	DirectionName string `json:"directionName,omitempty"`
	DeptId        int64  `json:"deptId,string,omitempty"`
	DeptName      string `json:"deptName,omitempty"`
	ExpertNum     int    `json:"expertNum,string,omitempty"`
	ReportedNum   int    `json:"reportedNum,string,omitempty"`
	Deleted       int64  `gorm:"column:deleted" json:"-"`
	basemodel.BaseEntity
}

func (EaJudgeTaskReportRequirement) TableName() string { return "ea_judge_task_report_requirement" }

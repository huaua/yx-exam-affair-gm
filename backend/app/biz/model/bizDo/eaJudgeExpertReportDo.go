package bizDo

import (
	"time"

	"github.com/sealsee/web-base/public/basemodel"
)

// EaJudgeExpertReport 学院按任务研究方向上报的专家记录。
type EaJudgeExpertReport struct {
	ReportId      int64     `gorm:"primaryKey" json:"reportId,string,omitempty"`
	JudgeTaskId   int64     `json:"judgeTaskId,string,omitempty"`
	RequirementId int64     `json:"requirementId,string,omitempty"`
	DeptId        int64     `json:"deptId,string,omitempty"`
	ExpertId      int64     `json:"expertId,string,omitempty"`
	SubmitState   int       `json:"submitState,string,omitempty"` // 1-未提交 2-已提交
	AuditState    int       `json:"auditState,string,omitempty"`  // 1-待审核 2-审核通过 3-审核不通过
	ReplacedFrom  int64     `json:"replacedFrom,string,omitempty"`
	DeleteTime    *time.Time `gorm:"column:delete_time" json:"deleteTime,omitempty"`
	basemodel.BaseEntity
}

func (EaJudgeExpertReport) TableName() string { return "ea_judge_expert_report" }

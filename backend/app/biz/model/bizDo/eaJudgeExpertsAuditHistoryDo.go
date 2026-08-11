package bizDo

import "time"

// EaJudgeExpertsAuditHistory 保存专家上一次审核通过数据的快照。
type EaJudgeExpertsAuditHistory struct {
	HistoryId    int64     `gorm:"primaryKey" json:"historyId,string,omitempty"`
	ExpertId     int64     `json:"expertId,string,omitempty"`
	DeptId       int64     `json:"deptId,string,omitempty"`
	AuditStatus  int       `json:"auditStatus,string,omitempty"`
	SnapshotData string    `gorm:"type:longtext" json:"snapshotData,omitempty"`
	ChangeBy     int64     `json:"changeBy,string,omitempty"`
	ChangeTime   time.Time `json:"changeTime,omitempty"`
}

func (EaJudgeExpertsAuditHistory) TableName() string {
	return "ea_judge_experts_audit_history"
}

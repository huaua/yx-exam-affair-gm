package bizDo

import "time"

// EaTchInfoAuditHistory 保存教职工审核通过后的修改前快照。
type EaTchInfoAuditHistory struct {
	HistoryId    int64     `gorm:"primaryKey" json:"historyId,string,omitempty"`
	TchId        int64     `json:"tchId,string,omitempty"`
	DeptId       int64     `json:"deptId,string,omitempty"`
	AuditStatus  int       `json:"auditStatus,string,omitempty"`
	SnapshotData string    `gorm:"type:longtext" json:"snapshotData,omitempty"`
	ChangeBy     int64     `json:"changeBy,string,omitempty"`
	ChangeTime   time.Time `json:"changeTime,omitempty"`
}

func (EaTchInfoAuditHistory) TableName() string {
	return "ea_tch_info_audit_history"
}

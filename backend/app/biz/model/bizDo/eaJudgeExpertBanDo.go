package bizDo

import "time"

// EaJudgeExpertBan records experts who must be excluded from judge extraction.
type EaJudgeExpertBan struct {
	BanId      int64     `gorm:"primaryKey" json:"banId,string,omitempty"`
	ExpertId   int64     `json:"expertId,string,omitempty"`
	IdCard     string    `gorm:"column:id_card" json:"idCard,omitempty"` // 证件号
	CreateBy   int64     `json:"createBy,string,omitempty"`
	CreateTime time.Time `json:"createTime,omitempty"`
	UpdateBy   int64     `json:"updateBy,string,omitempty"` // 更新者
	UpdateTime time.Time `json:"updateTime,omitempty"`      // 更新时间
	Deleted    int       `json:"-"`                         // 1正常 2逻辑删除
}

func (EaJudgeExpertBan) TableName() string { return "ea_judge_expert_ban" }

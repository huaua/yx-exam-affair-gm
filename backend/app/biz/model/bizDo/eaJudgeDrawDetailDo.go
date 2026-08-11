package bizDo

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 抽取评委明细数据实体
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeDrawDetail struct {
	Id             int64  `gorm:"primaryKey" json:"id,string,omitempty"` // 主键 关联ID
	ExpertId       int64  `json:"expertId,string,omitempty"`             // 专家ID
	ExpertName     string `json:"expertName,omitempty"`                  // 专家姓名
	DeptId         int64  `json:"deptId,string,omitempty"`               // 专家部门ID
	JudgeTaskId    int64  `json:"judgeTaskId,string,omitempty"`          // 任务ID
	JudgeTaskName  string `gorm:"-" json:"judgeTaskName,omitempty"`      // 任务名称
	TaskTimeRemark string `json:"taskTimeRemark,omitempty"`              // 评分日期
	DrawTimes      int    `json:"drawTimes,string,omitempty"`            // 抽取次数
	SubjectName    string `json:"subjectName,omitempty"`                 // 科目名称
	ConfirmState   int    `json:"confirmState,string,omitempty"`         // 确认状态 1-未应答 2-确认参加 3-不参加
	DrawBySubject  int    `json:"drawBySubject,string,omitempty"`        // 是否按擅长科目抽取 1-是 2-否
	Remark         string `json:"remark,omitempty"`                      // 备注
	Sex            string `gorm:"-" json:"sex,omitempty"`                // 性别（男或女）
	PhoneNo        string `gorm:"-" json:"phoneNo,omitempty"`            // 手机号
	DeptName       string `gorm:"-" json:"deptName,omitempty"`           // 专家来源-所属学院
	Type           string `gorm:"-" json:"type,omitempty"`               // 专家类型，多选字典值
	TypeStr        string `gorm:"-" json:"typeStr,omitempty"`            // 专家类型 1-博士 2-硕士 3-本科 4-附中
	GoodSubjects   string `gorm:"-" json:"goodSubjects,omitempty"`       // 擅长科目
	basemodel.BaseEntity
}

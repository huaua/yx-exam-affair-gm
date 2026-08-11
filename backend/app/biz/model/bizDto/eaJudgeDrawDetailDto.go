package bizDto

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 抽取评委明细数据查询实体
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeDrawDetailQuery struct {
	Id             int64  `gorm:"primaryKey" form:"id" json:"id,string,omitempty"`   // 主键 关联ID
	ExpertId       int64  `form:"expertId" json:"expertId,string,omitempty"`         // 专家ID
	ExpertName     string `form:"expertName" json:"expertName,omitempty"`            // 专家姓名
	DeptId         int64  `json:"deptId,string,omitempty"`                           // 专家部门ID
	JudgeTaskId    int64  `form:"judgeTaskId" json:"judgeTaskId,string,omitempty"`   // 任务ID
	TaskTimeRemark string `form:"taskTimeRemark" json:"taskTimeRemark,omitempty"`    // 评分日期
	DrawTimes      int    `form:"drawTimes" json:"drawTimes,string,omitempty"`       // 抽取次数
	SubjectName    string `form:"subjectName" json:"subjectName,omitempty"`          // 科目名称
	ConfirmState   int    `form:"confirmState" json:"confirmState,string,omitempty"` // 确认状态 1-未应答 2-确认参加 3-不参加
	DrawBySubject  int    `json:"drawBySubject,omitempty"`                           // 是否按擅长科目抽取 1-是 2-否
	Remark         string `form:"remark" json:"remark,omitempty"`                    // 备注
	basemodel.BaseEntityQuery
}

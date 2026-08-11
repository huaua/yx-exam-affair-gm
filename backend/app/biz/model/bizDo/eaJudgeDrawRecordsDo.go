package bizDo

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 评委抽取规则数据实体
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeDrawRecords struct {
	RulesId              int64  `gorm:"primaryKey" json:"rulesId,string,omitempty"` // 主键 规则ID
	JudgeTaskId          int64  `json:"judgeTaskId,string,omitempty"`               // 任务ID
	DrawTimes            int    `json:"drawTimes,string,omitempty"`                 // 抽取次数
	SubjectName          string `json:"subjectName,omitempty"`                      // 科目名称
	RequiredNum          int    `json:"requiredNum,string,omitempty"`               // 该轮次/科目需要的评委数
	OnCampusTotalNum     int    `json:"onCampusTotalNum,string,omitempty"`          // 校内抽取数
	OnCampusAttendNum    int    `json:"onCampusAttendNum,string,omitempty"`         // 校内参加评委数
	OnCampusUnattendNum  int    `json:"onCampusUnattendNum,string,omitempty"`       // 校内不参加评委数
	OffCampusTotalNum    int    `json:"offCampusTotalNum,string,omitempty"`         // 校外抽取数
	OffCampusAttendNum   int    `json:"offCampusAttendNum,string,omitempty"`        // 校外参加评委数
	OffCampusUnattendNum int    `json:"offCampusUnattendNum,string,omitempty"`      // 校外不参加评委数
	Remark               string `json:"remark,omitempty"`                           // 备注
	basemodel.BaseEntity
}

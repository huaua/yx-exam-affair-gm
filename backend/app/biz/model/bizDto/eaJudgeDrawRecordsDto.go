package bizDto

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 评委抽取规则数据查询实体
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeDrawRecordsQuery struct {
	RulesId              int64  `gorm:"primaryKey" form:"rulesId" json:"rulesId,string,omitempty"`         // 主键 规则ID
	JudgeTaskId          int64  `form:"judgeTaskId" json:"judgeTaskId,string,omitempty"`                   // 任务ID
	DrawTimes            int    `form:"drawTimes" json:"drawTimes,string,omitempty"`                       // 抽取次数
	SubjectName          string `form:"subjectName" json:"subjectName,omitempty"`                          // 科目名称
	OnCampusTotalNum     int    `form:"onCampusTotalNum" json:"onCampusTotalNum,string,omitempty"`         // 校内抽取数
	OnCampusAttendNum    int    `form:"onCampusAttendNum" json:"onCampusAttendNum,string,omitempty"`       // 校内参加评委数
	OnCampusUnattendNum  int    `form:"onCampusUnattendNum" json:"onCampusUnattendNum,string,omitempty"`   // 校内不参加评委数
	OffCampusTotalNum    int    `form:"offCampusTotalNum" json:"offCampusTotalNum,string,omitempty"`       // 校外抽取数
	OffCampusAttendNum   int    `form:"offCampusAttendNum" json:"offCampusAttendNum,string,omitempty"`     // 校外参加评委数
	OffCampusUnattendNum int    `form:"offCampusUnattendNum" json:"offCampusUnattendNum,string,omitempty"` // 校外不参加评委数
	Remark               string `form:"remark" json:"remark,omitempty"`                                    // 备注
	basemodel.BaseEntityQuery
}

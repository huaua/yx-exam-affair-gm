package bizDto

import (
	"github.com/sealsee/web-base/public/basemodel"
)

/*
 * desc: 评委征集任务数据查询实体
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeTaskInfoQuery struct {
	JudgeTaskId          int64              `gorm:"primaryKey" form:"judgeTaskId" json:"judgeTaskId,string,omitempty"` // 主键 评委征集任务ID
	JudgeTaskName        string             `form:"judgeTaskName" json:"judgeTaskName,omitempty"`                      // 评委任务名称
	TaskYear            int                `form:"taskYear" json:"taskYear,string,omitempty"`
	TaskMode            int                `form:"taskMode" json:"taskMode,string,omitempty"`
	JudgeTaskType        int                `form:"judgeTaskType" json:"judgeTaskType,string,omitempty"`               // 评委任务专家类型 1-博士 2-硕士 3-本科 4-附中
	TaskStartTime        basemodel.BaseTime `form:"taskStartTime" json:"taskStartTime,string,omitempty"`               // 评分开始时间
	TaskEndTime          basemodel.BaseTime `form:"taskEndTime" json:"taskEndTime,string,omitempty"`                   // 评分结束时间
	DrawSubject          int                `form:"drawSubject" json:"drawSubject,string,omitempty"`                   // 是否按擅长科目抽取: 1-是 2-否
	State                int                `form:"state" json:"state,string,omitempty"`                               // 抽取状态: 1-未开始 2-抽取中 3-已完成
	DrawTimes            int                `form:"drawTimes" json:"drawTimes,string,omitempty"`                       // 抽取次数
	ScoreRounds          int                `form:"scoreRounds" json:"scoreRounds,string,omitempty"`
	SubjectCount         int                `form:"subjectCount" json:"subjectCount,string,omitempty"`
	ExpertsPerRound      int                `form:"expertsPerRound" json:"expertsPerRound,string,omitempty"`
	PublishState         int                `form:"publishState" json:"publishState,string,omitempty"`
	OnCampusTotalNum     int                `form:"onCampusTotalNum" json:"onCampusTotalNum,string,omitempty"`         // 校内评委需求量
	OnCampusAttendNum    int                `form:"onCampusAttendNum" json:"onCampusAttendNum,string,omitempty"`       // 校内参加评委数
	OnCampusUnattendNum  int                `json:"onCampusUnattendNum,string,omitempty"`                              // 校内不参加评委数
	OffCampusTotalNum    int                `form:"offCampusTotalNum" json:"offCampusTotalNum,string,omitempty"`       // 校外评委需求量
	OffCampusAttendNum   int                `form:"offCampusAttendNum" json:"offCampusAttendNum,string,omitempty"`     // 校外参加评委数
	OffCampusUnattendNum int                `json:"offCampusUnattendNum,string,omitempty"`                             // 校外不参加评委数
	Remark               string             `form:"remark" json:"remark,omitempty"`                                    // 备注
	basemodel.BaseEntityQuery
}

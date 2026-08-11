package bizDo

import (
	"github.com/sealsee/web-base/public/basemodel"
)

/*
 * desc: 评委征集任务数据实体
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeTaskInfo struct {
	JudgeTaskId          int64              `gorm:"primaryKey" json:"judgeTaskId,string,omitempty"` // 主键 评委征集任务ID
	JudgeTaskName        string             `json:"judgeTaskName,omitempty"`                        // 评委任务名称
	TaskYear             int                `json:"taskYear,string,omitempty"`                      // 任务年份
	TaskMode             int                `json:"taskMode,string,omitempty"`                      // 任务模式: 1-直接抽取 2-学院上报
	JudgeTaskType        int                `json:"judgeTaskType,string,omitempty"`                 // 评委任务专家类型 1-博士 2-硕士 3-本科 4-附中
	JudgeTaskTypeStr     string             `gorm:"-" json:"judgeTaskTypeStr,omitempty"`            // 评委任务专家类型描述
	TaskStartTime        basemodel.BaseTime `json:"taskStartTime,omitempty"`                        // 评分开始时间
	TaskEndTime          basemodel.BaseTime `json:"taskEndTime,omitempty"`                          // 评分结束时间
	DrawSubject          int                `json:"drawSubject,string,omitempty"`                   // 是否按擅长科目抽取: 1-是 2-否
	State                int                `json:"state,string,omitempty"`                         // 抽取状态: 1-未开始 2-抽取中 3-已完成
	DrawTimes            int                `json:"drawTimes,string,omitempty"`                     // 抽取次数
	ScoreRounds          int                `json:"scoreRounds,string,omitempty"`                   // 评分轮次
	SubjectCount         int                `json:"subjectCount,string,omitempty"`                  // 科目数量
	ExpertsPerRound      int                `json:"expertsPerRound,string,omitempty"`               // 每轮专家数
	PublishState         int                `json:"publishState,string,omitempty"`                  // 发布状态: 1-未发布 2-已发布
	StateStr             string             `gorm:"-" json:"stateStr,omitempty"`                    // 抽取状态描述
	OnCampusTotalNum     int                `json:"onCampusTotalNum,string,omitempty"`              // 校内评委需求量
	OnCampusAttendNum    int                `json:"onCampusAttendNum,string,omitempty"`             // 校内参加评委数
	OnCampusUnattendNum  int                `json:"onCampusUnattendNum,string,omitempty"`           // 校内不参加评委数
	OffCampusTotalNum    int                `json:"offCampusTotalNum,string,omitempty"`             // 校外评委需求量
	OffCampusAttendNum   int                `json:"offCampusAttendNum,string,omitempty"`            // 校外参加评委数
	OffCampusUnattendNum int                `json:"offCampusUnattendNum,string,omitempty"`          // 校外不参加评委数
	ReportCount          int                `gorm:"-" json:"reportCount,string,omitempty"`             // 上报数（学院上报模式下的已上报记录数）
	RequirementCount     int                `gorm:"-" json:"requirementCount,string,omitempty"`        // 需求数（直接抽取=轮次*科目*每轮专家数；上报模式=各方向专家需求总和）
	DrawCount            int                `gorm:"-" json:"drawCount,string,omitempty"`               // 已抽取评委数（直接抽取模式下 ea_judge_draw_selection 评分专家数量）
	Remark               string             `json:"remark,omitempty"`                               // 备注
	basemodel.BaseEntity
}

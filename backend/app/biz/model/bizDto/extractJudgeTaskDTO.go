package bizDto

/*
 * desc: 评委抽取传输实体
 * author: 蓝鲸
 * date: 2025/06/11
 */
type ExtractJudgeTaskDTO struct {
	JudgeTaskId      int64  `form:"judgeTaskId" json:"judgeTaskId,string,omitempty"`           // 任务id
	SubjectName      string `form:"subjectName" json:"subjectName,omitempty"`                  // 科目名称
	OnCampusNeedNum  int    `form:"onCampusNeedNum" json:"onCampusNeedNum,string,omitempty"`   // 校内评委需求量
	OffCampusNeedNum int    `form:"offCampusNeedNum" json:"offCampusNeedNum,string,omitempty"` // 校外评委需求量
}

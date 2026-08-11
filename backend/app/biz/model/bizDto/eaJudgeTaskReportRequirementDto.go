package bizDto

type EaJudgeTaskReportRequirementQuery struct {
	RequirementId int64             `json:"requirementId,string,omitempty"`
	JudgeTaskId   int64             `json:"judgeTaskId,string,omitempty"`
	ExpertNum     int               `json:"expertNum,string,omitempty"`
	Files         map[string]string `json:"files,omitempty"`
	OperId        int64             `json:"operId,omitempty"`
}

// BatchUpdateReqItem 批量修改上报专家数的单条请求
type BatchUpdateReqItem struct {
	RequirementId int64 `json:"requirementId,string,omitempty"`
	ExpertNum     int   `json:"expertNum,string,omitempty"`
}

// BatchUpdateReportRequirementReq 批量修改上报专家数请求
type BatchUpdateReportRequirementReq struct {
	Items []BatchUpdateReqItem `json:"items"`
}

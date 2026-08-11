package bizDto

/*
 * desc: 考场监考分配结果查询实体
 * author: liyu
 * date: 2024/09/04
 */
type EaArrangeTchSaveDto struct {
	RoomTaskId      int64    `form:"roomTaskId" json:"roomTaskId,string,omitempty"`     // 考场任务id
	RoomReportId    int64    `form:"roomReportId" json:"roomReportId,string,omitempty"` // 考场上报ID
	TchReportIdList []string `form:"tchReportIdList" json:"tchReportIdList,omitempty"`  // 教职工上报id
}

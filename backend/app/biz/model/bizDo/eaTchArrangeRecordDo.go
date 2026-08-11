package bizDo

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 教职工监考记录视图数据（倒序展示）。
 * author: 蓝鲸
 * date: 2026/07/23
 */
type EaTchArrangeRecord struct {
	TchArrangeId  int64              `json:"tchArrangeId,string,omitempty"` // 监考分配ID
	TchReportId   int64              `json:"tchReportId,string,omitempty"`  // 教职工上报ID
	TchTaskId     int64              `json:"tchTaskId,string,omitempty"`    // 教师征集任务ID
	TchTaskName   string             `json:"tchTaskName,omitempty"`         // 考试任务名称
	TaskTime      basemodel.BaseTime `json:"taskTime,omitempty"`            // 监考日期
	RoomReportId  int64              `json:"roomReportId,string,omitempty"` // 考场安排ID
	RoomName      string             `json:"roomName,omitempty"`            // 考场名称
	ArrangeNoStr  string             `json:"arrangeNoStr,omitempty"`        // 考场编号
	TaskStartTime basemodel.BaseTime `json:"taskStartTime,omitempty"`       // 任务开始时间（用于排序）
}

func (EaTchArrangeRecord) TableName() string {
	return "ea_arrange_tch"
}

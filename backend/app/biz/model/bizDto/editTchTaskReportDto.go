package bizDto

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
)

/*
 * desc: 老师任务分配数据编辑实体
 * author: 蓝鲸
 * date: 2024/09/06
 */
type EditTchTaskReportQuery struct {
	TaskDetailId int64              `form:"taskDetailId" json:"taskDetailId,string,omitempty"` // 教师征集详情ID
	DeptId       int64              `form:"deptId" json:"deptId,string,omitempty"`             // 学院id
	TchList      []*bizDo.EaTchInfo `form:"tchList" json:"tchList,omitempty"`                  // 教职工信息列表
}

package bizDo

import (
	"github.com/sealsee/web-base/public/basemodel"
)

/*
 * desc: 老师任务分配数据实体
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchTaskReport struct {
	ReportId             int64              `gorm:"primaryKey" json:"reportId,string,omitempty"` // 主键 上报ID
	TchTaskId            int64              `json:"tchTaskId,string,omitempty"`                  // 征集任务ID
	TaskDetailId         int64              `json:"taskDetailId,string,omitempty"`               // 教师征集详情ID
	TaskTime             basemodel.BaseTime `json:"taskTime,omitempty"`                          // 征集日期
	TchId                int64              `json:"tchId,string,omitempty"`                      // 教职工ID
	TchName              string             `json:"tchName,omitempty"`                           // 姓名
	Sex                  string             `json:"sex,omitempty"`                               // 性别（男或女）
	JobNo                string             `json:"jobNo,omitempty"`                             // 工号
	DeptId               int64              `json:"deptId,string,omitempty"`                     // 学院id
	State                string             `json:"state,omitempty"`                             // 征集状态: 1-已上报,待审核 2-待分配 3-不通过 4-已分配
	ReplacedFrom         int64              `json:"replacedFrom,string,omitempty"`               // 被替换的原上报记录ID
	Remark               string             `json:"remark,omitempty"`                            // 备注
	ReportTaskDetailList []*EaTchTaskDetail `gorm:"-" json:"reportTaskDetailList,omitempty"`     // 上报征集详情列表
	basemodel.BaseEntity
}

// 暂存上报记录对象，学院上报时回显使用
type EaTchTaskReportTmp struct {
	TchId     int64                `json:"tchId,string,omitempty"` // 教职工ID
	TchName   string               `json:"tchName,omitempty"`      // 姓名
	Sex       string               `json:"sex,omitempty"`          // 性别（男或女）
	JobNo     string               `json:"jobNo,omitempty"`        // 工号
	Remark    string               `json:"remark,omitempty"`       // 备注
	TaskTimes []basemodel.BaseTime `json:"taskTimes,omitempty"`    // 征集日期
}

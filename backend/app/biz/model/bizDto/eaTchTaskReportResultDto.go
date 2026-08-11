package bizDto

import (
	"github.com/sealsee/web-base/public/basemodel"
)

/*
 * desc: 老师任务分配数据查询实体
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchTaskReportResultDTO struct {
	Sex                   string             `form:"sex" json:"sex,omitempty"`                                    // 性别
	Duties                string             `form:"duties" json:"duties,omitempty"`                              // 职务
	ProfOffice            string             `form:"profOffice" json:"profOffice,omitempty"`                      // 专业或行政
	ReportId              int64              `form:"reportId" json:"reportId,string,omitempty"`                   // 主键 上报ID
	TchTaskId             int64              `form:"tchTaskId" json:"tchTaskId,string,omitempty"`                 // 征集任务ID
	TchTaskName           string             `form:"tchTaskName" json:"tchTaskName,omitempty"`                    // 征集任务名称
	TaskDetailId          int64              `form:"taskDetailId" json:"taskDetailId,string,omitempty"`           // 教师征集详情ID
	TaskTime              basemodel.BaseTime `form:"taskTime" json:"taskTime,omitempty"`                          // 征集日期
	TchId                 int64              `form:"tchId" json:"tchId,string,omitempty"`                         // 教职工ID
	TchName               string             `form:"tchName" json:"tchName,omitempty"`                           // 姓名
	JobNo                 string             `form:"jobNo" json:"jobNo,omitempty"`                                // 工号
	IdCard                string             `form:"idCard" json:"idCard,omitempty"`                              // 证件号
	Education             string             `form:"education" json:"education,omitempty"`                        // 专业背景
	InvigilationExperience string             `form:"invigilationExperience" json:"invigilationExperience,omitempty"` // 监考经验 1-有 2-无
	InvigilationCount     int                `form:"invigilationCount" json:"invigilationCount,omitempty"`        // 监考次数
	Remark                string             `form:"remark" json:"remark,omitempty"`                              // 备注
	DeptId                int64              `form:"deptId" json:"deptId,string,omitempty"`                       // 学院id
	DeptName              string             `form:"deptName" json:"deptName,omitempty"`                          // 学院名称
	State                 string             `form:"state" json:"state,omitempty"`                                // 征集状态: 1-已上报,待审核 2-待分配 3-不通过 4-已分配
}

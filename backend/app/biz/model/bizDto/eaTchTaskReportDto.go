package bizDto

import (
	"github.com/sealsee/web-base/public/basemodel"
)

/*
 * desc: 审核征集状态请求（招办）
 * author: 2026/07/23
 */
type AuditTchTaskReportQuery struct {
	ReportId int64  `form:"reportId" json:"reportId,string,omitempty"` // 上报ID
	Action   string `form:"action" json:"action,omitempty"`           // pass / reject
}

/*
 * desc: 重置征集状态请求（招办）
 * author: 2026/07/23
 */
type ResetTchTaskReportQuery struct {
	ReportId int64 `form:"reportId" json:"reportId,string,omitempty"` // 上报ID
}

/*
 * desc: 老师任务分配数据查询实体
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchTaskReportQuery struct {
	ReportId            int64              `gorm:"primaryKey" form:"reportId" json:"reportId,string,omitempty"` // 主键 上报ID
	TchTaskId           int64              `form:"tchTaskId" json:"tchTaskId,string,omitempty"`                 // 征集任务ID
	TaskDetailId        int64              `form:"taskDetailId" json:"taskDetailId,string,omitempty"`           // 教师征集详情ID
	TaskTime            basemodel.BaseTime `form:"taskTime" json:"taskTime,omitempty"`                          // 征集日期
	TchId               int64              `form:"tchId" json:"tchId,string,omitempty"`                         // 教职工ID
	TchName             string             `form:"tchName" json:"tchName,omitempty"`                            // 姓名
	JobNo               string             `form:"jobNo" json:"jobNo,omitempty"`                                // 工号
	Sex                 string             `form:"sex" json:"sex,omitempty"`                                    // 性别（男/女）
	InvigilationExperience string           `form:"invigilationExperience" json:"invigilationExperience,omitempty"` // 监考经验 1-有 2-无
	TchState            string             `form:"tchState" json:"tchState,omitempty"`                          // 教职工状态 1-已删除 2-未删除
	DeptId              int64              `form:"deptId" json:"deptId,string,omitempty"`                       // 学院id
	State               string             `form:"state" json:"state,omitempty"`                                // 征集状态: 1-已上报,待审核 2-待分配 3-不通过 4-已分配
	Remark              string             `form:"remark" json:"remark,omitempty"`                              // 备注
	TchNameOrJobNo      string             `gorm:"-" form:"tchNameOrJobNo" json:"tchNameOrJobNo,omitempty"`     //
	basemodel.BaseEntityQuery
}

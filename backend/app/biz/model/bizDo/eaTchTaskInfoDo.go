package bizDo

import (
	"github.com/sealsee/web-base/public/basemodel"
)

/*
 * desc: 监考老师征集数据实体
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchTaskInfo struct {
	TchTaskId      int64              `gorm:"primaryKey" json:"tchTaskId,string,omitempty"` // 主键 教师征集任务ID
	TchTaskName    string             `json:"tchTaskName,omitempty"`                        // 考试任务名称
	TaskTimeRemark string             `json:"taskTimeRemark,omitempty"`                     // 征集日期说明
	TaskStartTime  basemodel.BaseTime `json:"taskStartTime,omitempty"`                      // 征集开始时间
	TaskEndTime    basemodel.BaseTime `json:"taskEndTime,omitempty"`                        // 征集结束时间
	NeedNum        int                `json:"needNum,string,omitempty"`                     // 总征集数量
	CollectNum     int                `json:"collectNum,string,omitempty"`                  // 已征集数量
	AuditPassNum   int                `json:"auditPassNum,string,omitempty"`                // 审核通过数
	State          string             `json:"state,omitempty"`                              // 发布状态: 1-未发布 2-已发布
	Remark         string             `json:"remark,omitempty"`                             // 备注
	TchTaskDetail  []*EaTchTaskDetail `gorm:"-" json:"tchTaskDetail"`                       // 监考老师征集任务详情
	basemodel.BaseEntity
}

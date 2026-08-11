package bizDto

import (
	"github.com/sealsee/web-base/public/basemodel"
)

/*
 * desc: 监考老师征集数据查询实体
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchTaskInfoQuery struct {
	TchTaskId      int64              `gorm:"primaryKey" form:"tchTaskId" json:"tchTaskId,string,omitempty"` // 主键 教师征集任务ID
	TchTaskName    string             `form:"tchTaskName" json:"tchTaskName,omitempty"`                      // 考试任务名称
	TaskTimeRemark string             `form:"taskTimeRemark" json:"taskTimeRemark,omitempty"`                // 征集日期说明
	TaskStartTime  basemodel.BaseTime `form:"taskStartTime" json:"taskStartTime,omitempty"`                  // 征集开始时间
	TaskEndTime    basemodel.BaseTime `form:"taskEndTime" json:"taskEndTime,omitempty"`                      // 征集结束时间
	NeedNum        int                `form:"needNum" json:"needNum,string,omitempty"`                       // 总征集数量
	CollectNum     int                `form:"collectNum" json:"collectNum,string,omitempty"`                 // 已征集数量
	State          string             `form:"state" json:"state,omitempty"`                                  // 发布状态: 1-未发布 2-已发布
	Remark         string             `form:"remark" json:"remark,omitempty"`                                // 备注
	basemodel.BaseEntityQuery
}

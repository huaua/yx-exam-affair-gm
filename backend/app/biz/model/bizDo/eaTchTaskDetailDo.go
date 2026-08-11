package bizDo

import (
	"github.com/sealsee/web-base/public/basemodel"
)

/*
 * desc: 监考老师征集任务明细数据实体
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchTaskDetail struct {
	TaskDetailId int64              `gorm:"primaryKey" json:"taskDetailId,string,omitempty"` // 主键 教师征集详情ID
	TchTaskId    int64              `json:"tchTaskId,string,omitempty"`                      // 教师征集任务ID
	TaskTime     basemodel.BaseTime `json:"taskTime,omitempty"`                              // 征集日期
	DeptId       int64              `json:"deptId,string,omitempty"`                         // 学院id
	NeedNum      int                `json:"needNum,string,omitempty"`                        // 征集人数
	CollectNum   int                `json:"collectNum,string,omitempty"`                     // 已征集数量
	ForceFlag    string             `json:"forceFlag,omitempty"`                             // 是否强制征集: 1-强制 2-不强制
	Remark       string             `json:"remark,omitempty"`                                // 备注
	DeptName     string             `gorm:"-" json:"deptName,omitempty"`                     // 学院名称
	basemodel.BaseEntity
}

package bizDto

import (
	"github.com/sealsee/web-base/public/basemodel"
)

/*
 * desc: 监考老师征集任务明细数据查询实体
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchTaskDetailQuery struct {
	TaskDetailId                int64              `gorm:"primaryKey" form:"taskDetailId" json:"taskDetailId,string,omitempty"`               // 主键 教师征集详情ID
	TchTaskId                   int64              `form:"tchTaskId" json:"tchTaskId,string,omitempty"`                                       // 教师征集任务ID
	TaskTime                    basemodel.BaseTime `form:"taskTime" json:"taskTime,omitempty"`                                                // 征集日期
	DeptId                      int64              `form:"deptId" json:"deptId,string,omitempty"`                                             // 学院id
	NeedNum                     int                `form:"needNum" json:"needNum,string,omitempty"`                                           // 征集人数
	CollectNum                  int                `form:"collectNum" json:"collectNum,string,omitempty"`                                     // 已征集数量
	ForceFlag                   string             `form:"forceFlag" json:"forceFlag,omitempty"`                                              // 是否强制征集: 1-强制 2-不强制
	Remark                      string             `form:"remark" json:"remark,omitempty"`                                                    // 备注
	OnlyQueryNeedReportTimeList bool               `gorm:"-" form:"onlyQueryNeedReportTimeList" json:"onlyQueryNeedReportTimeList,omitempty"` // 仅查询当前学院需要上报的时间列表
	basemodel.BaseEntityQuery
}

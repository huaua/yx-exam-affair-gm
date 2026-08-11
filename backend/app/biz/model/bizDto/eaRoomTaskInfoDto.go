package bizDto

import (
	"github.com/sealsee/web-base/public/basemodel"
)

/*
 * desc: 征集任务数据查询实体
 * author: 蓝鲸
 * date: 2024/08/14
 */
type EaRoomTaskInfoQuery struct {
	TaskId     int64  `gorm:"primaryKey" form:"taskId" json:"taskId,string,omitempty"` // 主键 任务ID
	TaskName   string `form:"taskName" json:"taskName,omitempty"`                      // 任务名称
	TaskDay    string `form:"taskDay" json:"taskDay,omitempty"`                        // 征集日期
	LevelCode  string `form:"levelCode" json:"levelCode,omitempty"`                    // 所属层次
	NeedNum    int    `form:"needNum" json:"needNum,string,omitempty"`                 // 需要考场数量
	CollectNum int    `form:"collectNum" json:"collectNum,string,omitempty"`           // 已征集数量
	ComfirmNum int    `form:"comfirmNum" json:"comfirmNum,string,omitempty"`           // 已确认数量
	State      string `form:"state" json:"state,omitempty"`                            // 发布状态: 1-未发布 2-已发布
	Remark     string `form:"remark" json:"remark,omitempty"`                          // 备注
	basemodel.BaseEntityQuery
}

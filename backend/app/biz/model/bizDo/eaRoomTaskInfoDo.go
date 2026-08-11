package bizDo

import (
	"github.com/sealsee/web-base/public/basemodel"
)

/*
 * desc: 征集任务数据实体
 * author: 蓝鲸
 * date: 2024/08/14
 */
type EaRoomTaskInfo struct {
	TaskId     int64              `gorm:"primaryKey" json:"taskId,string,omitempty"` // 主键 任务ID
	TaskName   string             `json:"taskName,omitempty"`                        // 任务名称
	TaskDay    basemodel.BaseTime `json:"taskDay,omitempty"`                         // 征集日期
	LevelCode  string             `json:"levelCode,omitempty"`                       // 所属层次
	NeedNum    int                `json:"needNum,string,omitempty"`                  // 需要考场数量
	CollectNum int                `json:"collectNum,string,omitempty"`               // 已征集数量
	ComfirmNum int                `json:"comfirmNum,string,omitempty"`               // 已确认数量
	State      string             `json:"state,omitempty"`                           // 发布状态: 1-未发布 2-已发布
	Remark     string             `json:"remark,omitempty"`                          // 备注
	basemodel.BaseEntity
}

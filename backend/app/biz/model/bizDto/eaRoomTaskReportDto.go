package bizDto

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 任务教室征集关联数据查询实体
 * author: 蓝鲸
 * date: 2024/08/14
 */
type EaRoomTaskReportQuery struct {
	Id                int64  `gorm:"primaryKey" form:"id" json:"id,string,omitempty"`               // 主键 id
	TaskId            int64  `form:"taskId" json:"taskId,string,omitempty"`                         // 任务ID
	RoomId            int64  `form:"roomId" json:"roomId,string,omitempty"`                         // 教室ID
	RoomName          string `form:"roomName" json:"roomName,omitempty"`                            // 教室名称
	TchNum            int    `form:"tchNum" json:"tchNum,string,omitempty"`                         // 监考老师数量
	RoomType          string `form:"roomType" json:"roomType,omitempty"`                            // 教室类型 1-按组 2-按位
	DeptId            int64  `form:"deptId" json:"deptId,string,omitempty"`                         // 所属学院id
	State             string `form:"state" json:"state,omitempty"`                                  // 确认状态: 1-未确认 2-已确认
	Remark            string `form:"remark" json:"remark,omitempty"`                                // 备注
	OnlyQuerySelfDept bool   `gorm:"-" form:"onlyQuerySelfDept" json:"onlyQuerySelfDept,omitempty"` // 仅查询自己部门的
	PublishState      string `gorm:"-" form:"publishState" json:"publishState,omitempty"`           // 教室使用状态 1-未上报 2-已上报未确认 3-已上报已确认
	basemodel.BaseEntityQuery
}

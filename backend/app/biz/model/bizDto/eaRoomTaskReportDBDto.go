package bizDto

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 任务教室征集关联数据实体-用于raw sql执行结果映射
 * author: 鲤鱼
 * date: 2024/09/18
 */
type EaRoomTaskReportDBDto struct {
	Id            int64  `gorm:"primaryKey" json:"id,string,omitempty"` // 主键 id
	TaskId        int64  `json:"taskId,string,omitempty"`               // 任务ID
	TaskDay       string `json:"taskDay,omitempty"`                     // 征集日期
	RoomId        int64  `json:"roomId,string,omitempty"`               // 教室ID
	RoomName      string `json:"roomName,omitempty"`                    // 教室名称
	DeptId        int64  `json:"deptId,string,omitempty"`               // 所属学院id
	TchNum        int    `json:"tchNum,string,omitempty"`               // 监考老师数量
	State         string `json:"state,omitempty"`                       // 确认状态: 1-未确认 2-已确认
	Remark        string `json:"remark,omitempty"`                      // 备注
	CampusName    string `json:"campusName,omitempty"`                  // 所在校区
	LevelCode     string `json:"levelCode,omitempty"`                   // 所属层次
	GroupNum      int    `json:"groupNum,string,omitempty"`             // 组数
	GroupCapacity int    `json:"groupCapacity,string,omitempty"`        // 按组容量
	Capacity      int    `json:"capacity,string,omitempty"`             // 按位容量
	basemodel.BaseEntity
}

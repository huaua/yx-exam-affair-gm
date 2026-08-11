package bizDo

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 教室数据实体
 * author: 蓝鲸
 * date: 2024/08/14
 */
type EaRoomInfo struct {
	RoomId        int64  `gorm:"primaryKey" json:"roomId,string,omitempty"` // 主键 教室ID
	RoomName      string `json:"roomName,omitempty"`                        // 教室名称
	CampusName    string `json:"campusName,omitempty"`                      // 所在校区
	LevelCode     string `json:"levelCode,omitempty"`                       // 所属层次
	DeptId        int64  `json:"deptId,string,omitempty"`                   // 所属学院id
	GroupNum      int    `json:"groupNum,string,omitempty"`                 // 组数
	GroupCapacity int    `json:"groupCapacity,string,omitempty"`            // 按组容量
	Capacity      int    `json:"capacity,string,omitempty"`                 // 按位容量
	State         string `json:"state,omitempty"`                           // 启用状态: 1-正常 2-停用
	Remark        string `json:"remark,omitempty"`                          // 备注
	DeptName      string `gorm:"-" json:"deptName,omitempty"`               // 学院名称
	basemodel.BaseEntity
}

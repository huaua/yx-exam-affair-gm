package bizDto

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 教室数据查询实体
 * author: 蓝鲸
 * date: 2024/08/14
 */
type EaRoomInfoQuery struct {
	RoomId     int64             `gorm:"primaryKey" form:"roomId" json:"roomId,string,omitempty"` // 主键 教室ID
	RoomName   string            `form:"roomName" json:"roomName,omitempty"`                      // 教室名称
	CampusName string            `form:"campusName" json:"campusName,omitempty"`                  // 所在校区
	LevelCode  string            `form:"levelCode" json:"levelCode,omitempty"`                    // 所属层次
	RoomType   string            `form:"roomType" json:"roomType,omitempty"`                      // 教室类型
	DeptId     int64             `form:"deptId" json:"deptId,string,omitempty"`                   // 所属学院id
	GroupNum   int               `form:"groupNum" json:"groupNum,string,omitempty"`               // 组数
	Capacity   int               `form:"capacity" json:"capacity,string,omitempty"`               // 容量
	State      string            `form:"state" json:"state,omitempty"`                            // 启用状态: 1-正常 2-停用
	Remark     string            `form:"remark" json:"remark,omitempty"`                          // 备注'
	Files      map[string]string `gorm:"-" form:"files" json:"files,omitempty"`                   // 导入文件 map[文件名]地址
	OperId     int64             `gorm:"-" form:"operId" json:"operId,omitempty"`                 // 操作人userid
	BizType    int               `gorm:"-" form:"bizType" json:"bizType,omitempty"`               //
	basemodel.BaseEntityQuery
}

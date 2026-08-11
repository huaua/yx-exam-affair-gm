package bizDo

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 监考老师分配数据实体
 * author: init code
 * date: 2024/09/24
 */
type EaArrangeTch struct {
	TchArrangeId int64  `gorm:"primaryKey" json:"tchArrangeId,string,omitempty"` // 主键 老师安排ID
	RoomTaskId   int64  `json:"roomTaskId,string,omitempty"`                     // 考场任务ID
	RoomReportId int64  `json:"roomReportId,string,omitempty"`                   // 考场安排ID
	TchReportId  int64  `json:"tchReportId,string,omitempty"`                    // 教职工上报ID
	TchTaskId    int64  `json:"tchTaskId,string,omitempty"`                      // 教师征集任务ID
	TchId        int64  `json:"tchId,string,omitempty"`                          // 教职工ID
	TchName      string `json:"tchName,omitempty"`                               // 教职工姓名（冗余）
	Remark       string `json:"remark,omitempty"`                                // 备注
	basemodel.BaseEntity
}

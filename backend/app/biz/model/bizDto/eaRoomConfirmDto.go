package bizDto

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
)

/*
 * desc: 考场征集列表显示实体
 * author: 蓝鲸
 * date: 2024/08/16
 */
type EaRoomConfirmDto struct {
	TaskId   int64                     `form:"taskId" json:"taskId,string,omitempty"` // 已征集任务数
	RoomList []*bizDo.EaRoomTaskReport `form:"roomList" json:"roomList,omitempty"`    // 教室列表
}

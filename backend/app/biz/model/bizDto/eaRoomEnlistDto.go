package bizDto

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
)

/*
 * desc: 考场征集列表显示实体
 * author: 蓝鲸
 * date: 2024/08/16
 */
type EaRoomEnlistDto struct {
	PublishState     string `json:"publishState,omitempty"` // 教室使用状态 1-未上报 2-已上报未确认 3-已上报已确认
	bizDo.EaRoomInfo        // 教室列表
}

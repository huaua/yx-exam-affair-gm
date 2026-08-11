package bizDto

import "yx-exam-affair-gm/app/biz/model/bizDo"

/*
 * desc: 考场编排详情数据查询实体
 * author: init code
 * date: 2024/09/24
 */
type EaArrangeRoomAddDto struct {
	InfoId          int64                  `form:"infoId" json:"infoId,string,omitempty"`            // 数据ID
	ArrangeRoomList []*bizDo.EaArrangeRoom `form:"arrangeRoomList" json:"arrangeRoomList,omitempty"` // 备注
}

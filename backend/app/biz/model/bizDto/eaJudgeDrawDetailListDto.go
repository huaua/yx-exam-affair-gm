package bizDto

import "yx-exam-affair-gm/app/biz/model/bizDo"

/*
 * desc: 抽取评委明细数据查询实体
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeDrawDetailList struct {
	AttendNum   int                        `json:"attendNum,string,omitempty"`   // 确认参加人数
	UnattendNum int                        `json:"unattendNum,string,omitempty"` // 确认不参加人数
	NotSureNum  int                        `json:"notSureNum,string,omitempty"`  // 待定人数
	List        []*bizDo.EaJudgeDrawDetail `json:"list,omitempty"`               // 评委列表
}

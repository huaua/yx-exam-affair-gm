package bizSer

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
)

/*
 * desc: 评委抽取规则服务层接口
 * author: init code
 * date: 2025/06/11
 */
type IEaJudgeDrawRecordsService interface {

	/*
	 * desc: 查询所有
	 * author：蓝鲸
	 * date: 2025/06/12
	 */
	ListAllEaJudgeDrawRecords(q *bizDto.EaJudgeDrawRecordsQuery) []*bizDo.EaJudgeDrawRecords
}

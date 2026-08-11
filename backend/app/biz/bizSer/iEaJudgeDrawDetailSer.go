package bizSer

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
)

/*
 * desc: 抽取评委明细服务层接口
 * author: init code
 * date: 2025/06/11
 */
type IEaJudgeDrawDetailService interface {

	/*
	 * desc: 查询满足条件的记录列表
	 * author: init code
	 * date: 2025/06/11
	 */
	ListEaJudgeDrawDetail(q *bizDto.EaJudgeDrawDetailQuery, p *page.Page) *bizDto.EaJudgeDrawDetailList

	/*
	 * desc: 查询满足条件的所有记录列表
	 * author: init code
	 * date: 2025/06/11
	 */
	ListAllEaJudgeDrawDetail(q *bizDto.EaJudgeDrawDetailQuery, needQueryTaskInfo bool) []*bizDo.EaJudgeDrawDetail

	/*
	 * desc: 更新数据
	 * author: init code
	 * date: 2025/06/11
	 */
	UpdateEaJudgeDrawDetailById(judgeDrawDetail *bizDo.EaJudgeDrawDetail) (bool, string)

	/*
	 * desc: 导出数据
	 * author: init code
	 * date: 2025/06/11
	 */
	ExportEaJudgeDrawDetail(q *bizDto.EaJudgeDrawDetailQuery) ([]byte, error)
}

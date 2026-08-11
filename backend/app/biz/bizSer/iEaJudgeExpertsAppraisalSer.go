package bizSer

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
)

/*
 * desc: 评委评价服务层接口
 * author: init code
 * date: 2025/06/11
 */
type IEaJudgeExpertsAppraisalService interface {

	/*
	 * desc: 根据id获取数据
	 * author: init code
	 * date: 2025/06/11
	 */
	GetEaJudgeExpertsAppraisalById(id int64) *bizDo.EaJudgeExpertsAppraisal

	/*
	 * desc: 查询满足条件的记录列表
	 * author: init code
	 * date: 2025/06/11
	 */
	ListEaJudgeExpertsAppraisal(q *bizDto.EaJudgeExpertsAppraisalQuery, p *page.Page) []*bizDo.EaJudgeExpertsAppraisal

	/*
	 * desc: 插入数据
	 * author: init code
	 * date: 2025/06/11
	 */
	InsertEaJudgeExpertsAppraisal(judgeExpertsAppraisal *bizDo.EaJudgeExpertsAppraisal) (bool, string)

	/*
	 * desc: 更新数据
	 * author: init code
	 * date: 2025/06/11
	 */
	UpdateEaJudgeExpertsAppraisalById(judgeExpertsAppraisal *bizDo.EaJudgeExpertsAppraisal) bool

	/*
	 * desc: 删除数据
	 * author: init code
	 * date: 2025/06/11
	 */
	DeleteEaJudgeExpertsAppraisalById(id int64, deleteBy int64) bool

	/*
	 * desc: 导出数据
	 * author: init code
	 * date: 2025/06/11
	 */
	ExportEaJudgeExpertsAppraisal(q *bizDto.EaJudgeExpertsAppraisalQuery) ([]byte, error)
}

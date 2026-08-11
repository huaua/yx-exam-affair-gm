package bizDao

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

/*
 * desc: 评委专家库数据层接口
 * author: init code
 * date: 2025/06/06
 */
type IEaJudgeExpertsDao interface {

	/*
	 * desc: 根据id获取数据
	 * author: init code
	 * date: 2025/06/06
	 */
	GetById(id int64) *bizDo.EaJudgeExperts

	/*
	 * desc: 查询满足条件的条数
	 * author: init code
	 * date: 2025/06/06
	 */
	Count(q *bizDto.EaJudgeExpertsQuery) int

	/*
	 * desc: 查询满足条件的条数（支持特殊条件）
	 * author: init code
	 * date: 2025/06/06
	 */
	CountWithCondition(q *bizDto.EaJudgeExpertsQuery, condition interface{}, args ...interface{}) int

	/*
	 * desc: 查询满足条件的记录列表
	 * author: init code
	 * date: 2025/06/06
	 */
	List(q *bizDto.EaJudgeExpertsQuery, p *page.Page) []*bizDo.EaJudgeExperts

	/*
	 * desc: 查询满足条件的记录列表（支持特殊条件）
	 * author: init code
	 * date: 2025/06/06
	 */
	ListWithCondition(q *bizDto.EaJudgeExpertsQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaJudgeExperts

	/*
	 * desc: 查询满足条件的记录列表map（支持特殊条件）
	 * author: init code
	 * date: 2025/06/06
	 */
	ListMapWithCondition(q *bizDto.EaJudgeExpertsQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any

	/*
	 * desc: 新增记录
	 * author: init code
	 * date: 2025/06/06
	 */
	Insert(gtx tx.GTx, judgeExperts *bizDo.EaJudgeExperts) bool

	/*
	 * desc: 批量新增记录
	 * author: init code
	 * date: 2025/06/06
	 */
	BatchInsert(gtx tx.GTx, judgeExpertsList []*bizDo.EaJudgeExperts) bool

	/*
	 * desc: 更新记录
	 * author: init code
	 * date: 2025/06/06
	 */
	Update(gtx tx.GTx, uData *bizDo.EaJudgeExperts, uWhere *bizDo.EaJudgeExperts) bool

	/*
	 * desc: 删除记录
	 * author: init code
	 * date: 2025/06/06
	 */
	DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool

	/*
	 * desc: 更新（支持清空）
	 * author: init code
	 * date: 2025/06/06
	 */
	UpdateNew(gtx tx.Db, uWhere *bizDo.EaJudgeExperts, uData *bizDo.EaJudgeExperts) int

	/*
	 * desc: 特殊条件更新（支持清空）
	 * author: init code
	 * date: 2025/06/06
	 */
	UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaJudgeExperts, uData *bizDo.EaJudgeExperts, condition interface{}, args ...interface{}) int
}

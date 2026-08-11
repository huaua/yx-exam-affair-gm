package bizDao

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

/*
 * desc: 老师任务分配数据层接口
 * author: 蓝鲸
 * date: 2024/09/04
 */
type IEaTchTaskReportDao interface {

	/*
	 * desc: 根据id获取数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	GetById(id int64) *bizDo.EaTchTaskReport

	/*
	 * desc: 查询满足条件的条数
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	Count(q *bizDto.EaTchTaskReportQuery) int

	/*
	 * desc: 查询满足条件的条数（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	CountWithCondition(q *bizDto.EaTchTaskReportQuery, condition interface{}, args ...interface{}) int

	/*
	 * desc: 查询满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	List(q *bizDto.EaTchTaskReportQuery, p *page.Page) []*bizDo.EaTchTaskReport

	/*
	 * desc: 查询满足条件的记录列表（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ListWithCondition(q *bizDto.EaTchTaskReportQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaTchTaskReport

	/*
	 * desc: 查询满足条件的记录列表map（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ListMapWithCondition(q *bizDto.EaTchTaskReportQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any

	/*
	 * desc: 新增记录
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	Insert(gtx tx.GTx, tchTaskReport *bizDo.EaTchTaskReport) bool

	/*
	 * desc: 批量新增记录
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	BatchInsert(gtx tx.GTx, tchTaskReportList []*bizDo.EaTchTaskReport) bool

	/*
	 * desc: 更新记录
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	Update(gtx tx.GTx, uData *bizDo.EaTchTaskReport, uWhere *bizDo.EaTchTaskReport) bool

	/*
	 * desc: 删除记录
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool

	/*
	 * desc:批量删除
	 * author: 蓝鲸
	 * date: 2024/09/06
	 */
	DeleteByIds(gtx tx.Db, ids []interface{}, deleteBy int64) bool

	/*
	 * desc:根据监考征集任务明细id删除已征集的记录
	 * author: 鲤鱼
	 * date: 2024/09/06
	 */
	DeleteByTaskDetailId(gtx tx.GTx, taskDetailId, deleteBy int64) bool

	/*
	 * desc: 删除学院暂存的上报记录, 按时间单个清除
	 * author: 鲤鱼
	 * date: 2024/12/20
	 */
	DeleteTempSavedTchReportOne(gtx tx.GTx, deptId, taskId, taskDetailId, tchId int64) bool

	/*
	 * desc: 删除学院暂存的上报记录, 按任务整个清除
	 * author: 鲤鱼
	 * date: 2024/12/20
	 */
	DeleteTempSavedTchReport(gtx tx.GTx, deptId, taskId int64) bool

	/*
	 * desc: 更新（支持清空）
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	UpdateNew(gtx tx.Db, uWhere *bizDo.EaTchTaskReport, uData *bizDo.EaTchTaskReport) int

	/*
	 * desc: 特殊条件更新（支持清空）
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaTchTaskReport, uData *bizDo.EaTchTaskReport, condition interface{}, args ...interface{}) int
}

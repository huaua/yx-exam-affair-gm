package bizDao

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

/*
 * desc: 监考老师征集任务明细数据层接口
 * author: 蓝鲸
 * date: 2024/09/04
 */
type IEaTchTaskDetailDao interface {

	/*
	 * desc: 根据id获取数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	GetById(id int64) *bizDo.EaTchTaskDetail

	/*
	 * desc: 查询满足条件的条数
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	Count(q *bizDto.EaTchTaskDetailQuery) int

	/*
	 * desc: 查询满足条件的条数（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	CountWithCondition(q *bizDto.EaTchTaskDetailQuery, condition interface{}, args ...interface{}) int

	/*
	 * desc: 查询满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	List(q *bizDto.EaTchTaskDetailQuery, p *page.Page) []*bizDo.EaTchTaskDetail

	/*
	 * desc: 查询满足条件的记录列表（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ListWithCondition(q *bizDto.EaTchTaskDetailQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaTchTaskDetail

	/*
	 * desc: 查询满足条件的记录列表map（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ListMapWithCondition(q *bizDto.EaTchTaskDetailQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any

	/*
	 * desc: 新增记录
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	Insert(gtx tx.GTx, tchTaskDetail *bizDo.EaTchTaskDetail) bool

	/*
	 * desc: 批量新增记录
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	BatchInsert(gtx tx.GTx, tchTaskDetailList []*bizDo.EaTchTaskDetail) bool

	/*
	 * desc: 更新记录
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	Update(gtx tx.GTx, uData *bizDo.EaTchTaskDetail, uWhere *bizDo.EaTchTaskDetail) bool

	/*
	 * desc: 更新上报老师数量
	 * author: 蓝鲸
	 * date: 2024/09/06
	 */
	UpdateReportNum(gtx tx.GTx, tchTaskDetailId int64, reportNum int) bool

	/*
	 * desc: 删除记录
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool

	/*
	 * desc: 更新（支持清空）
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	UpdateNew(gtx tx.Db, uWhere *bizDo.EaTchTaskDetail, uData *bizDo.EaTchTaskDetail) int

	/*
	 * desc: 特殊条件更新（支持清空）
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaTchTaskDetail, uData *bizDo.EaTchTaskDetail, condition interface{}, args ...interface{}) int
}

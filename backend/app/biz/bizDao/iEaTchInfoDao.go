package bizDao

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

/*
 * desc: 教职工数据层接口
 * author: 蓝鲸
 * date: 2024/09/04
 */
type IEaTchInfoDao interface {

	/*
	 * desc: 根据id获取数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	GetById(id int64) *bizDo.EaTchInfo

	/*
	 * desc: 查询满足条件的条数
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	Count(q *bizDto.EaTchInfoQuery) int

	/*
	 * desc: 查询满足条件的条数（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	CountWithCondition(q *bizDto.EaTchInfoQuery, condition interface{}, args ...interface{}) int

	/*
	 * desc: 查询满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	List(q *bizDto.EaTchInfoQuery, p *page.Page) []*bizDo.EaTchInfo

	/*
	 * desc: 查询满足条件的记录列表（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ListWithCondition(q *bizDto.EaTchInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaTchInfo

	/*
	 * desc: 查询满足条件的记录列表map（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ListMapWithCondition(q *bizDto.EaTchInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any

	/*
	 * desc: 新增记录
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	Insert(gtx tx.GTx, tchInfo *bizDo.EaTchInfo) bool

	/*
	 * desc: 批量新增记录
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	BatchInsert(gtx tx.GTx, tchInfoList []*bizDo.EaTchInfo) bool

	/*
	 * desc: 更新记录
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	Update(gtx tx.GTx, uData *bizDo.EaTchInfo, uWhere *bizDo.EaTchInfo) bool

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
	UpdateNew(gtx tx.Db, uWhere *bizDo.EaTchInfo, uData *bizDo.EaTchInfo) int

	/*
	 * desc: 特殊条件更新（支持清空）
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaTchInfo, uData *bizDo.EaTchInfo, condition interface{}, args ...interface{}) int

	/*
	 * desc: 查询教职工的审核历史快照列表（按时间倒序）
	 * author: 蓝鲸
	 * date: 2026/07/23
	 */
	ListAuditHistory(tchId int64) []*bizDo.EaTchInfoAuditHistory

	/*
	 * desc: 查询教职工最近一次审核通过的快照（用于编辑对比）
	 * author: 蓝鲸
	 * date: 2026/07/23
	 */
	GetLatestAuditHistory(tchId int64) *bizDo.EaTchInfoAuditHistory

	/*
	 * desc: 根据教职工id重新计算并写回监考经验字段
	 * desc: 1-有监考记录 2-无监考记录
	 * author: 蓝鲸
	 * date: 2026/07/23
	 */
	RefreshInvigilationExperienceByTchId(gtx tx.GTx, tchId int64) bool
}

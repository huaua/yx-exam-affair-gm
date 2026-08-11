package bizDao

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

/*
 * desc: 监考老师分配数据层接口
 * author: init code
 * date: 2024/09/24
 */
type IEaArrangeTchDao interface {

	/*
	 * desc: 根据id获取数据
	 * author: init code
	 * date: 2024/09/24
	 */
	GetById(id int64) *bizDo.EaArrangeTch

	/*
	 * desc: 查询满足条件的条数
	 * author: init code
	 * date: 2024/09/24
	 */
	Count(q *bizDto.EaArrangeTchQuery) int

	/*
	 * desc: 查询满足条件的条数（支持特殊条件）
	 * author: init code
	 * date: 2024/09/24
	 */
	CountWithCondition(q *bizDto.EaArrangeTchQuery, condition interface{}, args ...interface{}) int

	/*
	 * desc: 查询满足条件的记录列表
	 * author: init code
	 * date: 2024/09/24
	 */
	List(q *bizDto.EaArrangeTchQuery, p *page.Page) []*bizDo.EaArrangeTch

	/*
	 * desc: 查询满足条件的记录列表（支持特殊条件）
	 * author: init code
	 * date: 2024/09/24
	 */
	ListWithCondition(q *bizDto.EaArrangeTchQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaArrangeTch

	/*
	 * desc: 查询满足条件的记录列表map（支持特殊条件）
	 * author: init code
	 * date: 2024/09/24
	 */
	ListMapWithCondition(q *bizDto.EaArrangeTchQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any

	/*
	 * desc: 新增记录
	 * author: init code
	 * date: 2024/09/24
	 */
	Insert(gtx tx.GTx, arrangeTch *bizDo.EaArrangeTch) bool

	/*
	 * desc: 批量新增记录
	 * author: init code
	 * date: 2024/09/24
	 */
	BatchInsert(gtx tx.GTx, arrangeTchList []*bizDo.EaArrangeTch) bool

	/*
	 * desc: 更新记录
	 * author: init code
	 * date: 2024/09/24
	 */
	Update(gtx tx.GTx, uData *bizDo.EaArrangeTch, uWhere *bizDo.EaArrangeTch) bool

	/*
	 * desc: 删除记录
	 * author: init code
	 * date: 2024/09/24
	 */
	DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool

	/*
	 * desc: 根据教室上报id删除记录
	 * author: 蓝鲸
	 * date: 2024/09/30
	 */
	DeleteByRoomReportId(gtx tx.GTx, id int64, deleteBy int64) bool

	/*
	 * desc: 根据教职工上报id删除单条分配记录
	 * author: system
	 * date: 2026/07/24
	 */
	DeleteByTchReportId(gtx tx.GTx, roomReportId int64, tchReportId int64, deleteBy int64) bool

	/*
	 * desc: 更新（支持清空）
	 * author: init code
	 * date: 2024/09/24
	 */
	UpdateNew(gtx tx.Db, uWhere *bizDo.EaArrangeTch, uData *bizDo.EaArrangeTch) int

	/*
	 * desc: 特殊条件更新（支持清空）
	 * author: init code
	 * date: 2024/09/24
	 */
	UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaArrangeTch, uData *bizDo.EaArrangeTch, condition interface{}, args ...interface{}) int
}

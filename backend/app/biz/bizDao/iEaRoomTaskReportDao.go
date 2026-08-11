package bizDao

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

/*
 * desc: 任务教室征集关联数据层接口
 * author: 蓝鲸
 * date: 2024/08/14
 */
type IEaRoomTaskReportDao interface {

	/*
	 * desc: 根据id获取数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	GetById(id int64) *bizDo.EaRoomTaskReport

	/*
	 * desc: 查询满足条件的条数
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	Count(q *bizDto.EaRoomTaskReportQuery) int

	/*
	 * desc: 查询满足条件的条数（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	CountWithCondition(q *bizDto.EaRoomTaskReportQuery, condition interface{}, args ...interface{}) int

	/*
	 * desc: 查询满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	List(q *bizDto.EaRoomTaskReportQuery, p *page.Page) []*bizDo.EaRoomTaskReport

	/*
	 * desc: 查询满足条件的记录列表（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ListWithCondition(q *bizDto.EaRoomTaskReportQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaRoomTaskReport

	/*
	 * desc: 查询满足条件的记录列表map（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ListMapWithCondition(q *bizDto.EaRoomTaskReportQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any

	/*
	 * desc: 新增记录
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	Insert(gtx tx.GTx, taskRoomInfo *bizDo.EaRoomTaskReport) bool

	/*
	 * desc: 批量新增记录
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	BatchInsert(gtx tx.GTx, taskRoomInfoList []*bizDo.EaRoomTaskReport) bool

	/*
	 * desc: 更新记录
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	Update(gtx tx.GTx, uData *bizDo.EaRoomTaskReport, uWhere *bizDo.EaRoomTaskReport) bool

	/*
	 * desc: 删除记录
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	DeleteById(gtx tx.GTx, id int64) bool

	/*
	 * desc: 更新（支持清空）
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	UpdateNew(gtx tx.Db, uWhere *bizDo.EaRoomTaskReport, uData *bizDo.EaRoomTaskReport) int

	/*
	 * desc: 特殊条件更新（支持清空）
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaRoomTaskReport, uData *bizDo.EaRoomTaskReport, condition interface{}, args ...interface{}) int
}

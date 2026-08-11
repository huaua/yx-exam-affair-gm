package bizDao

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

/*
 * desc: 教室数据层接口
 * author: 蓝鲸
 * date: 2024/08/14
 */
type IEaRoomInfoDao interface {

	/*
	 * desc: 根据id获取数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	GetById(id int64) *bizDo.EaRoomInfo

	/*
	 * desc: 查询满足条件的条数
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	Count(q *bizDto.EaRoomInfoQuery) int

	/*
	 * desc: 查询满足条件的条数（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	CountWithCondition(q *bizDto.EaRoomInfoQuery, condition interface{}, args ...interface{}) int

	/*
	 * desc: 查询满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	List(q *bizDto.EaRoomInfoQuery, p *page.Page) []*bizDo.EaRoomInfo

	/*
	 * desc: 查询满足条件的记录列表（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ListWithCondition(q *bizDto.EaRoomInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaRoomInfo

	/*
	 * desc: 查询满足条件的记录列表map（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ListMapWithCondition(q *bizDto.EaRoomInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any

	/*
	 * desc: 新增记录
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	Insert(gtx tx.GTx, classroomInfo *bizDo.EaRoomInfo) bool

	/*
	 * desc: 批量新增记录
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	BatchInsert(gtx tx.GTx, classroomInfoList []*bizDo.EaRoomInfo) bool

	/*
	 * desc: 更新记录
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	Update(gtx tx.GTx, uData *bizDo.EaRoomInfo, uWhere *bizDo.EaRoomInfo) bool

	/*
	 * desc: 删除记录
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool

	/*
	 * desc: 更新（支持清空）
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	UpdateNew(gtx tx.Db, uWhere *bizDo.EaRoomInfo, uData *bizDo.EaRoomInfo) int

	/*
	 * desc: 特殊条件更新（支持清空）
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaRoomInfo, uData *bizDo.EaRoomInfo, condition interface{}, args ...interface{}) int
}

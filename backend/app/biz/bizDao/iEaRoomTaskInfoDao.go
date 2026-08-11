package bizDao

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

/*
 * desc: 征集任务数据层接口
 * author: 蓝鲸
 * date: 2024/08/14
 */
type IEaRoomTaskInfoDao interface {

	/*
	 * desc: 根据id获取数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	GetById(id int64) *bizDo.EaRoomTaskInfo

	/*
	 * desc: 查询满足条件的条数
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	Count(q *bizDto.EaRoomTaskInfoQuery) int

	/*
	 * desc: 查询满足条件的条数（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	CountWithCondition(q *bizDto.EaRoomTaskInfoQuery, condition interface{}, args ...interface{}) int

	/*
	 * desc: 查询满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	List(q *bizDto.EaRoomTaskInfoQuery, p *page.Page) []*bizDo.EaRoomTaskInfo

	/*
	 * desc: 查询满足条件的记录列表（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ListWithCondition(q *bizDto.EaRoomTaskInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaRoomTaskInfo

	/*
	 * desc: 查询满足条件的记录列表map（支持特殊条件）
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ListMapWithCondition(q *bizDto.EaRoomTaskInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any

	/*
	 * desc: 新增记录
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	Insert(gtx tx.GTx, taskInfo *bizDo.EaRoomTaskInfo) bool

	/*
	 * desc: 批量新增记录
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	BatchInsert(gtx tx.GTx, taskInfoList []*bizDo.EaRoomTaskInfo) bool

	/*
	 * desc: 更新记录
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	Update(gtx tx.GTx, uData *bizDo.EaRoomTaskInfo, uWhere *bizDo.EaRoomTaskInfo) bool

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
	UpdateNew(gtx tx.Db, uWhere *bizDo.EaRoomTaskInfo, uData *bizDo.EaRoomTaskInfo) int

	/*
	 * desc: 特殊条件更新（支持清空）
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaRoomTaskInfo, uData *bizDo.EaRoomTaskInfo, condition interface{}, args ...interface{}) int

	/*
	 * desc: 更新上报人数和确认人数
	 * author: 蓝鲸
	 * date: 2024/08/16
	 */
	UpdateCollectNumAndConfirmNum(gtx tx.Db, taskId int64, collectNum int, confirmNum int) int
}

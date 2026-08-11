package bizSer

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
)

/*
 * desc: 征集任务服务层接口
 * author: 蓝鲸
 * date: 2024/08/14
 */
type IEaRoomTaskInfoService interface {

	/*
	 * desc: 根据id获取数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	GetEaTaskInfoById(id int64) *bizDo.EaRoomTaskInfo

	/*
	 * desc: 查询满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ListEaTaskInfo(q *bizDto.EaRoomTaskInfoQuery, p *page.Page) []*bizDo.EaRoomTaskInfo

	/*
	 * desc: 查询所有满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ListAllEaTaskInfo(q *bizDto.EaRoomTaskInfoQuery) []*bizDo.EaRoomTaskInfo

	/*
	 * desc: 插入数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	InsertEaTaskInfo(taskInfo *bizDo.EaRoomTaskInfo) (bool, string)

	/*
	 * desc: 更新数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	UpdateEaTaskInfoById(taskInfo *bizDo.EaRoomTaskInfo, onlyChangeStatus bool) (bool, string)

	/*
	 * desc: 删除数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	DeleteEaTaskInfoById(id int64, deleteBy int64) bool

	/*
	 * desc: 导出数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ExportEaTaskInfo(q *bizDto.EaRoomTaskInfoQuery) ([]byte, error)
}

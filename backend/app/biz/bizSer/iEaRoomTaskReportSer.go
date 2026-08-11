package bizSer

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
)

/*
 * desc: 任务教室征集关联服务层接口
 * author: 蓝鲸
 * date: 2024/08/14
 */
type IEaRoomTaskReportService interface {

	/*
	 * desc: 根据id获取数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	GetEaTaskRoomInfoById(id int64) *bizDo.EaRoomTaskReport

	/*
	 * desc: 查询所有满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/09/25
	 */
	ListAllRoomReport(q *bizDto.EaRoomTaskReportQuery) []*bizDo.EaRoomTaskReport

	/*
	 * desc: 查询所有满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ListAllEaTaskRoomInfo(q *bizDto.EaRoomTaskReportQuery) []*bizDto.EaRoomTaskReportDBDto

	/*
	 * desc: 查询任务下所有可选择的教室-重构
	 * author: 鲤鱼
	 * date: 2024/09/14
	 */
	ListCanUseRoomListByTaskId(q *bizDto.EaRoomTaskReportQuery, p *page.Page) (int, []*bizDto.EaRoomEnlistDto)

	/*
	 * desc: 插入数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	InsertEaTaskRoomInfo(taskRoomInfo *bizDo.EaRoomTaskReport) (bool, string)

	/*
	 * desc: 更新数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	UpdateEaTaskRoomInfoById(taskRoomInfo *bizDo.EaRoomTaskReport) bool

	/*
	 * desc: 批量更新数据
	 * author：蓝鲸
	 * date：2024/08/19
	 */
	BatchUpdateEaTaskRoomInfo(list []*bizDo.EaRoomTaskReport, taskId int64, updateBy int64) (bool, string)

	/*
	 * desc: 根据教室id更新教室名称
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	UpdateRoomNameById(classroomInfoDO *bizDo.EaRoomInfo, updateBy int64) bool

	/*
	 * desc: 取消上报
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	DeleteEaTaskRoomInfo(query *bizDto.EaRoomTaskReportQuery) (bool, string)

	/*
	 * desc: 删除数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	DeleteEaTaskRoomInfoById(id int64) bool

	/*
	 * desc: 导出数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ExportEaTaskRoomInfo(q *bizDto.EaRoomTaskReportQuery) ([]byte, error)

	/*
	 * desc: 重置单条征集确认（已征集→已上报，清空监考老师数量）
	 * author: 蓝鲸
	 * date: 2026/07/24
	 */
	ResetEaTaskRoomInfoById(id int64, updateBy int64) (bool, string)

	/*
	 * desc: 批量重置征集确认（已征集→已上报，清空监考老师数量）
	 * author: 蓝鲸
	 * date: 2026/07/24
	 */
	BatchResetEaTaskRoomInfo(list []*bizDo.EaRoomTaskReport, taskId int64, updateBy int64) (bool, string)

	/*
	 * desc: 单条征集确认（已上报→已征集，设置监考老师数量）
	 * author: 蓝鲸
	 * date: 2026/07/24
	 */
	ConfirmEaTaskRoomInfoById(id int64, tchNum int64, updateBy int64) (bool, string)
}

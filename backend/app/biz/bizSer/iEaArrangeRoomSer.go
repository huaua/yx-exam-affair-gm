package bizSer

import (
	"github.com/sealsee/web-base/public/errs"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
)

/*
 * desc: 考场编排详情服务层接口
 * author: init code
 * date: 2024/09/24
 */
type IEaArrangeRoomService interface {

	/*
	 * desc: 根据id获取数据
	 * author: init code
	 * date: 2024/09/24
	 */
	GetEaArrangeRoomById(id int64) *bizDo.EaArrangeRoom

	/*
	 * desc: 根据任务下最大的考场编号
	 * author: 蓝鲸
	 * date: 2024/09/26
	 */
	GetMaxArrangeRoomNo(roomTaskId int64) *bizDo.EaArrangeRoom

	/*
	 * desc: 查询满足条件的记录列表
	 * author: init code
	 * date: 2024/09/24
	 */
	ListEaArrangeRoom(q *bizDto.EaArrangeRoomQuery) []*bizDo.EaArrangeRoom

	/*
	 * desc: 查询教室分配记录列表
	 * author: 蓝鲸
	 * date: 2024/09/24
	 */
	//ListAssignment(q *bizDto.EaArrangeRoomQuery) []*bizDo.EaArrangeRoom

	/*
	 * desc: 查询所有满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/09/25
	 */
	ListAllArrangeRoom(q *bizDto.EaArrangeRoomQuery) []*bizDo.EaArrangeRoom

	/*
	 * desc: 插入数据
	 * author: init code
	 * date: 2024/09/24
	 */
	InsertEaArrangeRoom(arrangeRoom *bizDo.EaArrangeRoom) bool

	/*
	 * desc: 批量新增数据
	 * author: 蓝鲸
	 * date: 2024/09/26
	 */
	BatchInsertArrangeRoom(arrangeRoomAddDto *bizDto.EaArrangeRoomAddDto, optId int64) (bool, errs.ERROR)

	/*
	 * desc: 更新数据
	 * author: init code
	 * date: 2024/09/24
	 */
	UpdateEaArrangeRoomById(arrangeRoom *bizDo.EaArrangeRoom) bool

	/*
	 * desc: 批量更新数据
	 * author: 蓝鲸
	 * date: 2024/09/26
	 */
	BatchUpdateArrangeRoom(arrangeRoomAddDto *bizDto.EaArrangeRoomAddDto, optId int64) (bool, errs.ERROR)

	/*
	 * desc: 删除数据
	 * author: init code
	 * date: 2024/09/24
	 */
	DeleteEaArrangeRoomById(id int64, deleteBy int64) bool

	/*
	 * desc: 导出数据
	 * author: init code
	 * date: 2024/09/24
	 */
	ExportEaArrangeRoom(q *bizDto.EaArrangeRoomQuery) ([]byte, error)
}

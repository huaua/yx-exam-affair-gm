package bizSer

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
)

/*
 * desc: 教室服务层接口
 * author: 蓝鲸
 * date: 2024/08/14
 */
type IEaRoomInfoService interface {

	/*
	 * desc: 根据id获取数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	GetEaClassroomInfoById(id int64) *bizDo.EaRoomInfo

	/*
	 * desc: 查询满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ListEaClassroomInfo(q *bizDto.EaRoomInfoQuery, p *page.Page) []*bizDo.EaRoomInfo

	/*
	 * desc: 查询所有满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ListAllEaClassroomInfo(q *bizDto.EaRoomInfoQuery) []*bizDo.EaRoomInfo

	/*
	 * desc: 插入数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	InsertEaClassroomInfo(classroomInfo *bizDo.EaRoomInfo) (bool, string)

	/*
	 * desc: 更新数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	UpdateEaClassroomInfoById(classroomInfo *bizDo.EaRoomInfo, onlyChangeStatus bool, deptId int64) (bool, string)

	/*
	 * desc: 删除数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	DeleteEaClassroomInfoById(id int64, deleteBy int64) bool

	/*
	 * desc: 导出数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ExportEaClassroomInfo(q *bizDto.EaRoomInfoQuery) ([]byte, error)

	/*
	 * desc: 导入数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ImportEaClassroomInfo(q *bizDto.EaRoomInfoQuery) (bool, string)
}

package bizSer

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
)

/*
 * desc: 监考老师征集任务明细服务层接口
 * author: 蓝鲸
 * date: 2024/09/04
 */
type IEaTchTaskDetailService interface {

	/*
	 * desc: 根据id获取数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	GetEaTchTaskDetailById(id int64) *bizDo.EaTchTaskDetail

	/*
	 * desc: 查询满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ListEaTchTaskDetail(q *bizDto.EaTchTaskDetailQuery, p *page.Page) []*bizDo.EaTchTaskDetail

	/*
	 * desc: 查询所有满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ListAllEaTchTaskDetail(q *bizDto.EaTchTaskDetailQuery) []*bizDo.EaTchTaskDetail

	/*
	 * desc: 查询满足条件的记录列表-按日分组
	 * author: 鲤鱼
	 * date: 2024/09/05
	 */
	ListEaTchTaskDetailByDay(q *bizDto.EaTchTaskDetailQuery) []*map[string]any

	/*
	 * desc: 插入数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	InsertEaTchTaskDetail(tchTaskDetail *bizDo.EaTchTaskDetail) bool

	/*
	 * desc: 更新数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	UpdateEaTchTaskDetailById(tchTaskDetail *bizDo.EaTchTaskDetail) bool

	/*
	 * desc: 更新已上报人数
	 * author: 蓝鲸
	 * date: 2024/09/06
	 */
	UpdateReportedTchNum(tchTaskDetailId int64, reportedTchNum int) bool

	/*
	 * desc: 删除数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	DeleteEaTchTaskDetailById(id int64, deleteBy int64) bool

	/*
	 * desc: 导出数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ExportEaTchTaskDetail(q *bizDto.EaTchTaskDetailQuery) ([]byte, error)
}

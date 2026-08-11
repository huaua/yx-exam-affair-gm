package bizSer

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/errs"
)

/*
 * desc: 监考老师征集服务层接口
 * author: 蓝鲸
 * date: 2024/09/04
 */
type IEaTchTaskInfoService interface {

	/*
	 * desc: 根据id获取数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	GetEaTchTaskInfoById(id int64) *bizDo.EaTchTaskInfo

	/*
	 * desc: 查询满足条件的条数
	 * author: 鲤鱼
	 * date: 2024/09/04
	 */
	CountEaTchTaskInfo(q *bizDto.EaTchTaskInfoQuery) int

	/*
	 * desc: 查询满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ListEaTchTaskInfo(q *bizDto.EaTchTaskInfoQuery, p *page.Page) []*bizDo.EaTchTaskInfo

	/*
	 * desc: 查询所有满足条件的
	 * author: 蓝鲸
	 * date: 2024/09/05
	 */
	ListAllEaTchTaskInfo(q *bizDto.EaTchTaskInfoQuery) []*bizDo.EaTchTaskInfo

	/*
	 * desc: 查询所有满足条件的
	 * author: 蓝鲸
	 * date: 2024/09/05
	 */
	ListAllEaTchTaskInfoByDept(deptId int64) []*bizDo.EaTchTaskInfo

	/*
	 * desc: 插入数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	InsertEaTchTaskInfo(tchTaskInfo *bizDo.EaTchTaskInfo) errs.ERROR

	/*
	 * desc: 修改发布状态
	 * author: 鲤鱼
	 * date: 2024/09/05
	 */
	UpdateTchTaskState(tchTaskId int64, state string, oper int64) bool

	/*
	 * desc: 更新数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	UpdateEaTchTaskInfoById(tchTaskInfo *bizDo.EaTchTaskInfo) errs.ERROR

	/*
	 * desc: 删除数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	DeleteEaTchTaskInfoById(id int64, deleteBy int64) bool

	/*
	 * desc: 导出数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ExportEaTchTaskInfo(q *bizDto.EaTchTaskInfoQuery) ([]byte, error)
}

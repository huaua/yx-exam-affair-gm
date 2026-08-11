package bizSer

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
)

/*
 * desc: 学院基础服务层接口
 * author: 蓝鲸
 * date: 2024/08/14
 */
type IEaDeptInfoService interface {

	/*
	 * desc: 根据id获取数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	GetEaDeptInfoById(id int64) *bizDo.EaDeptInfo

	/*
	 * desc: 查询满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ListEaDeptInfo(q *bizDto.EaDeptInfoQuery, p *page.Page) []*bizDo.EaDeptInfo

	/*
	 * desc: 获取所有已有学院列表
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ListAllEaDeptInfo(q *bizDto.EaDeptInfoQuery) []*bizDo.EaDeptInfo

	/*
	 * desc: 获取所有学院id对应的map
	 * author: 蓝鲸
	 * date: 2024/08/20
	 */
	GetAllDeptMap() map[int64]*bizDo.EaDeptInfo

	/*
	 * desc: 插入数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	BatchInsertEaDeptInfo(deptInfoList []*bizDo.EaDeptInfo, optId int64) (bool, string)

	/*
	 * desc: 更新数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	UpdateEaDeptInfoById(deptInfo *bizDo.EaDeptInfo) (bool, string)

	/*
	 * desc: 删除数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	DeleteEaDeptInfoById(id int64, deleteBy int64) bool

	/*
	 * desc: 导出数据
	 * author: 蓝鲸
	 * date: 2024/08/14
	 */
	ExportEaDeptInfo(q *bizDto.EaDeptInfoQuery) ([]byte, error)
}

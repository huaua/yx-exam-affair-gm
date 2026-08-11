package bizSer

import (
	"github.com/sealsee/web-base/public/errs"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
)

/*
 * desc: 监考老师分配服务层接口
 * author: init code
 * date: 2024/09/24
 */
type IEaArrangeTchService interface {

	/*
	 * desc: 根据id获取数据
	 * author: init code
	 * date: 2024/09/24
	 */
	GetEaArrangeTchById(id int64) *bizDo.EaArrangeTch

	/*
	 * desc: 查询满足条件的记录列表-按考场
	 * author: 鲤鱼
	 * date: 2024/09/29
	 */
	ListEaArrangeTch(q *bizDto.EaArrangeTchQuery, p *page.Page) []*bizDto.EaArrangeTchResultDto

	/*
	 * desc: 查询所有满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/09/29
	 */
	ListAllArrangeTch(q *bizDto.EaArrangeTchQuery) []*bizDo.EaArrangeTch

	/*
	 * desc: 查询满足条件的记录列表-按专业
	 * author: 鲤鱼
	 * date: 2024/09/30
	 */
	ListEaArrangeTchByProf(q *bizDto.EaArrangeTchQuery, p *page.Page) []*bizDto.EaArrangeProfResultDto

	/*
	 * desc: 插入数据
	 * author: init code
	 * date: 2024/09/24
	 */
	InsertEaArrangeTch(arrangeTch *bizDo.EaArrangeTch) bool

	/*
	 * desc: 更新数据
	 * author: init code
	 * date: 2024/09/24
	 */
	UpdateEaArrangeTchById(arrangeTch *bizDo.EaArrangeTch) bool

	/*
	 * desc: 批量新增或更新数据
	 * author: init code
	 * date: 2024/09/24
	 */
	BatchUpdateEaArrangeTch(saveDtoList []*bizDto.EaArrangeTchSaveDto, optId int64) (bool, errs.ERROR)

	/*
	 * desc: 删除数据
	 * author: init code
	 * date: 2024/09/24
	 */
	DeleteEaArrangeTchById(id int64, deleteBy int64) bool

	/*
	 * desc: 取消分配
	 * author: 蓝鲸
	 * date: 2024/09/30
	 */
	DelArrangeTchByRoomReportId(roomReportId int64, optId int64) bool

	/*
	 * desc: 取消单个监考老师分配
	 * author: system
	 * date: 2026/07/24
	 */
	DelSingleArrangeTch(roomReportId int64, tchReportId int64, optId int64) bool

	/*
	 * desc: 导出数据
	 * author: init code
	 * date: 2024/09/24
	 */
	ExportEaArrangeTch(q *bizDto.EaArrangeTchQuery) ([]byte, error)
}

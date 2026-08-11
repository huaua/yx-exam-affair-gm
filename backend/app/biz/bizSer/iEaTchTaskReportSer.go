package bizSer

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/errs"

	"github.com/sealsee/web-base/public/ds/page"
)

/*
 * desc: 老师任务分配服务层接口
 * author: 蓝鲸
 * date: 2024/09/04
 */
type IEaTchTaskReportService interface {

	/*
	 * desc: 根据id获取数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	GetEaTchTaskReportById(id int64) *bizDo.EaTchTaskReport

	/*
	 * desc: 查询满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ListEaTchTaskReport(q *bizDto.EaTchTaskReportQuery, p *page.Page) []*bizDto.EaTchTaskReportResultDTO

	/*
	 * desc: 查询所有满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/09/05
	 */
	ListAllEaTchTaskReport(q *bizDto.EaTchTaskReportQuery) []*bizDo.EaTchTaskReport

	/*
	 * desc: 查询暂存上报记录
	 * author: liyu
	 * date: 2024/12/23
	 */
	TempListEaTchTaskReport(deptId, tchTaskId int64) []*bizDo.EaTchTaskReportTmp
	/*
	 * desc: 暂存上报记录
	 * author: liyu
	 * date: 2024/12/20
	 */
	TempSaveEaTchTaskReport(deptId int64, tchTaskReportList []*bizDo.EaTchTaskReport, optId int64) (bool, string)

	/*
	 * desc: 清除暂存记录
	 * author: liyu
	 * date: 2024/12/20
	 */
	TempClearEaTchTaskReport(deptId int64, tchTaskReportList []*bizDo.EaTchTaskReport, optId int64) (bool, string)

	/*
	 * desc: 插入数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	InsertEaTchTaskReport(tchTaskReport *bizDo.EaTchTaskReport) bool

	/*
	 * desc: 批量新增数据
	 * author: 蓝鲸
	 * date: 2024/09/05
	 */
	BatchInsertEaTchTaskReport(deptId int64, tchTaskReportList []*bizDo.EaTchTaskReport, optId int64) (bool, string)

	/*
	 * desc: 更新数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	UpdateEaTchTaskReportById(tchTaskReport *bizDo.EaTchTaskReport) bool

	/*
	 * desc: 编辑数据
	 * author: 蓝鲸
	 * date: 2024/09/06
	 */
	EditEaTchTaskReport(editInfo *bizDto.EditTchTaskReportQuery, optId int64) (bool, string)

	/*
	 * desc: 删除数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	DeleteEaTchTaskReportById(id int64, deleteBy int64) (success bool, err errs.ERROR)

	/*
	 * desc: 导出数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ExportEaTchTaskReport(q *bizDto.EaTchTaskReportQuery) ([]byte, error)

	/*
	 * desc: 查询可上报的教职工列表
	 * author: 蓝鲸
	 * date: 2024/09/05
	 */
	ListEnabledEaTchInfo(q *bizDto.EaTchTaskReportQuery) ([]*bizDo.EaTchInfo, string)

	/*
	 * desc: 校验是否可以直接取消
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	CheckCancelResult(id int64) (success bool, err errs.ERROR)

	/*
	 * desc: 校验是否可以直接取消
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ReplaceReport(tchTaskReport *bizDo.EaTchTaskReport) (success bool, err errs.ERROR)

	/*
	 * desc: 审核征集状态（招办）— pass→待分配, reject→不通过
	 * author: 2026/07/23
	 */
	AuditTchTaskReport(reportId int64, action string, optId int64) (success bool, err errs.ERROR)

	/*
	 * desc: 重置征集状态（招办）— 已审/已驳回 重置回 待审核
	 * author: 2026/07/23
	 */
	ResetTchTaskReport(reportId int64, optId int64) (success bool, err errs.ERROR)
}

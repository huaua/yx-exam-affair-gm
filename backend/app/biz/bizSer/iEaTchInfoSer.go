package bizSer

import (
	"github.com/sealsee/web-base/public/errs"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
)

/*
 * desc: 教职工服务层接口
 * author: 蓝鲸
 * date: 2024/09/04
 */
type IEaTchInfoService interface {

	/*
	 * desc: 根据id获取数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	GetEaTchInfoById(id int64) *bizDo.EaTchInfo
	GetNextPendingTchInfo(q *bizDto.EaTchInfoQuery) *bizDo.EaTchInfo

	/*
	 * desc: 查询满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ListEaTchInfo(q *bizDto.EaTchInfoQuery, p *page.Page) []*bizDo.EaTchInfo

	/*
	 * desc: 查询所有满足条件的记录列表
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ListAllEaTchInfo(q *bizDto.EaTchInfoQuery, condition interface{}, args ...interface{}) []*bizDo.EaTchInfo

	/*
	 * desc: 插入数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	InsertEaTchInfo(tchInfo *bizDo.EaTchInfo) (bool, string)

	/*
	 * desc: 更新数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	UpdateEaTchInfoById(tchInfo *bizDo.EaTchInfo, deptId int64) (success bool, errInfo string)

	/*
	 * desc: 删除数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	DeleteEaTchInfoById(id int64, deleteBy int64, deptId int64) (success bool, errMsg errs.ERROR)

	/*
	 * desc: 导出数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ExportEaTchInfo(q *bizDto.EaTchInfoQuery) ([]byte, error)

	/*
	 * desc: 导入数据
	 * author: 蓝鲸
	 * date: 2024/09/04
	 */
	ImportEaTchInfo(q *bizDto.EaTchInfoQuery) (success bool, errInfo string)

	/*
	 * desc: 招办审核教职工（通过/不通过）
	 * author: 蓝鲸
	 * date: 2026/07/23
	 */
	AuditEaTchInfo(tchId int64, action string, reason string, auditBy int64) (bool, string)

	/*
	 * desc: 招办修改教职工的启用状态
	 * author: 蓝鲸
	 * date: 2026/07/23
	 */
	ChangeEnabledStatus(tchId int64, enabledStatus int, optBy int64) (bool, string)

	/*
	 * desc: 招办批量修改教职工的启用状态
	 * author: 蓝鲸
	 * date: 2026/07/23
	 */
	ChangeEnabledStatusBatch(tchIds []int64, enabledStatus int, optBy int64) (int, string)

	/*
	 * desc: 获取教职工的监考记录列表（倒序）
	 * author: 蓝鲸
	 * date: 2026/07/23
	 */
	ListInvigilationRecords(tchId int64) []*bizDo.EaTchArrangeRecord

	/*
	 * desc: 获取教职工最近一次审核通过的快照（用于编辑对比）
	 * author: 蓝鲸
	 * date: 2026/07/23
	 */
	GetLatestAuditHistory(tchId int64) *bizDo.EaTchInfoAuditHistory

	/*
	 * desc: 获取教职工所有审核历史
	 * author: 蓝鲸
	 * date: 2026/07/23
	 */
	ListAuditHistory(tchId int64) []*bizDo.EaTchInfoAuditHistory

	/*
	 * desc: 学院用户将未提交的教职工数据提交至待审核状态
	 * author: 蓝鲸
	 * date: 2026/07/23
	 */
	SubmitForAudit(tchId int64, deptId int64, submitBy int64) (bool, string)
	SubmitForAuditBatch(ids []int64, deptId int64, submitBy int64) (int, string)
}

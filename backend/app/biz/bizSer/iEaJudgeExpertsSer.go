package bizSer

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
)

/*
 * desc: 评委专家库服务层接口
 * author: init code
 * date: 2025/06/06
 */
type IEaJudgeExpertsService interface {

	/*
	 * desc: 根据id获取数据
	 * author: init code
	 * date: 2025/06/06
	 */
	GetEaJudgeExpertsById(id int64) *bizDo.EaJudgeExperts
	GetLatestAuditHistory(id int64) *bizDo.EaJudgeExpertsAuditHistory
	GetNextPendingExpert(q *bizDto.EaJudgeExpertsQuery) *bizDo.EaJudgeExperts

	/*
	 * desc: 查询满足条件的记录列表
	 * author: init code
	 * date: 2025/06/06
	 */
	ListEaJudgeExperts(q *bizDto.EaJudgeExpertsQuery, p *page.Page) []*bizDo.EaJudgeExperts

	/*
	 * desc: 插入数据
	 * author: init code
	 * date: 2025/06/06
	 */
	InsertEaJudgeExperts(judgeExperts *bizDo.EaJudgeExperts) (bool, string)

	/*
	 * desc: 更新数据
	 * author: init code
	 * date: 2025/06/06
	 */
	UpdateEaJudgeExpertsById(judgeExperts *bizDo.EaJudgeExperts, sessionDeptId int64) (bool, string)
	SubmitForAudit(id int64, sessionDeptId int64, submitBy int64) (bool, string)
	SubmitForAuditBatch(ids []int64, sessionDeptId int64, submitBy int64) (int, string)
	AuditEaJudgeExperts(id int64, action string, reason string, auditBy int64) (bool, string)
	EvaluateEaJudgeExperts(id int64, evaluation string, evaluateBy int64) (bool, string)
	ChangeAllowDraw(id int64, allowDraw int, updateBy int64) (bool, string)
	BanEaJudgeExpert(id int64, updateBy int64) (bool, string)
	UnbanEaJudgeExperts(ids []int64, updateBy int64) (bool, string)
	UnbanAllEaJudgeExperts(updateBy int64) (bool, string)
	ImportBannedEaJudgeExperts(q *bizDto.EaJudgeExpertsQuery) (bool, string)

	/*
	 * desc: 删除数据
	 * author: init code
	 * date: 2025/06/06
	 */
	DeleteEaJudgeExpertsById(id int64, deleteBy int64) bool

	// 导入评委专家库
	ImportEaJudgeExperts(q *bizDto.EaJudgeExpertsQuery) (bool, string)

	/*
	 * desc: 导出数据
	 * author: init code
	 * date: 2025/06/06
	 */
	ExportEaJudgeExperts(q *bizDto.EaJudgeExpertsQuery) ([]byte, error)
}

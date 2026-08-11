package bizSer

import (
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
)

/*
 * desc: 评委征集任务服务层接口
 * author: init code
 * date: 2025/06/11
 */
type IEaJudgeTaskInfoService interface {

	/*
	 * desc: 根据id获取数据
	 * author: init code
	 * date: 2025/06/11
	 */
	GetEaJudgeTaskInfoById(id int64) []*bizDo.EaJudgeDrawRecords

	/*
	 * desc: 查询满足条件的记录列表
	 * author: init code
	 * date: 2025/06/11
	 */
	ListEaJudgeTaskInfo(q *bizDto.EaJudgeTaskInfoQuery, p *page.Page) []*bizDo.EaJudgeTaskInfo

	/*
	 * desc: 获取所有数据
	 * author: 蓝鲸
	 * date: 2025/06/16
	 */
	ListAllEaJudgeTaskInfo(q *bizDto.EaJudgeTaskInfoQuery) []*bizDo.EaJudgeTaskInfo

	/*
	 * desc: 插入数据
	 * author: init code
	 * date: 2025/06/11
	 */
	InsertEaJudgeTaskInfo(judgeTaskInfo *bizDo.EaJudgeTaskInfo) (bool, string)

	/*
	 * desc: 更新数据
	 * author: init code
	 * date: 2025/06/11
	 */
	UpdateEaJudgeTaskInfoById(judgeTaskInfo *bizDo.EaJudgeTaskInfo) (bool, string)

	/*
	 * desc: 删除数据
	 * author: init code
	 * date: 2025/06/11
	 */
	DeleteEaJudgeTaskInfoById(id int64, deleteBy int64) (bool, string)

	/*
	 * desc: 不按擅长科目抽取
	 * author: 蓝鲸
	 * date: 2025/06/11
	 */
	ExtractJudgeTask(extractJudgeTask *bizDto.ExtractJudgeTaskDTO, optId int64) (bool, string)

	/*
	 * desc: 按擅长科目抽取
	 * author: 蓝鲸
	 * date: 2025/06/11
	 */
	ExtractJudgeTaskBySubject(extractJudgeTaskList []*bizDto.ExtractJudgeTaskDTO, optId int64) (bool, string)

	/*
	 * desc: 取消抽取
	 * author: 蓝鲸
	 * date: 2025/06/11
	 */
	CancelExtract(id int64, optId int64) (bool, string)
	AppendScoreRounds(id int64, rounds int, optId int64) (bool, string)
	PublishCollegeReportTask(id int64, optId int64) (bool, string)
	UnpublishCollegeReportTask(id int64, optId int64) (bool, string)
	CancelImportReports(id int64, optId int64) (bool, string)

	ListReportRequirements(taskId int64) []*bizDo.EaJudgeTaskReportRequirement
	ImportReportRequirements(q *bizDto.EaJudgeTaskReportRequirementQuery) (bool, string)
	UpdateReportRequirement(id int64, expertNum int, optId int64) (bool, string)
	BatchUpdateReportRequirements(items []bizDto.BatchUpdateReqItem, optId int64) (bool, string)
	DeleteReportRequirement(id int64, optId int64) (bool, string)
}

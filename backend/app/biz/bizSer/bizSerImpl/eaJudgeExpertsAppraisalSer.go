package bizSerImpl

import (
	"fmt"
	"time"
	"yx-exam-affair-gm/app/biz/bizDao"
	"yx-exam-affair-gm/app/biz/bizDao/bizDaoImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
	"github.com/sealsee/web-base/public/utils/file/excel"
)

/*
 * desc: 评委评价服务层实现
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeExpertsAppraisalService struct{}

var eaJudgeExpertsAppraisalService *EaJudgeExpertsAppraisalService
var eaJudgeExpertsAppraisalDao bizDao.IEaJudgeExpertsAppraisalDao

func init() {
	eaJudgeExpertsAppraisalService = &EaJudgeExpertsAppraisalService{}
	eaJudgeExpertsAppraisalDao = bizDaoImpl.NewEaJudgeExpertsAppraisalDao()
}

func NewEaJudgeExpertsAppraisalService() *EaJudgeExpertsAppraisalService {
	return eaJudgeExpertsAppraisalService
}

func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalService) GetEaJudgeExpertsAppraisalById(id int64) *bizDo.EaJudgeExpertsAppraisal {
	return eaJudgeExpertsAppraisalDao.GetById(id)
}

func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalService) ListEaJudgeExpertsAppraisal(q *bizDto.EaJudgeExpertsAppraisalQuery, p *page.Page) []*bizDo.EaJudgeExpertsAppraisal {
	return eaJudgeExpertsAppraisalDao.List(q, p)
}

func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalService) InsertEaJudgeExpertsAppraisal(judgeExpertsAppraisalDO *bizDo.EaJudgeExpertsAppraisal) (success bool, errInfo string) {

	// 获取任务详情
	dbTaskInfo := eaJudgeTaskInfoDao.GetById(judgeExpertsAppraisalDO.JudgeTaskId)
	if dbTaskInfo == nil {
		return false, "任务不存在"
	}

	// 校验当前评委是否参与了此次评分任务
	exists := eaJudgeDrawDetailDao.Count(&bizDto.EaJudgeDrawDetailQuery{ExpertId: judgeExpertsAppraisalDO.ExpertId,
		JudgeTaskId: judgeExpertsAppraisalDO.JudgeTaskId}) > 0
	if !exists {
		return false, "当前评委未参与此次评分任务"
	}

	// 构建参数
	judgeExpertsAppraisalDO.JudgeTaskName = dbTaskInfo.JudgeTaskName
	judgeExpertsAppraisalDO.TaskTimeRemark = fmt.Sprintf("%v-%v", dbTaskInfo.TaskStartTime.FormatString(time.DateOnly), dbTaskInfo.TaskEndTime.FormatString(time.DateOnly))

	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaJudgeExpertsAppraisalDao.Insert(gtx, judgeExpertsAppraisalDO)
	})
	return
}

func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalService) UpdateEaJudgeExpertsAppraisalById(judgeExpertsAppraisalDO *bizDo.EaJudgeExpertsAppraisal) (success bool) {
	uWhere := &bizDo.EaJudgeExpertsAppraisal{AppraisalId: judgeExpertsAppraisalDO.AppraisalId}
	judgeExpertsAppraisalDO.AppraisalId = 0
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaJudgeExpertsAppraisalDao.Update(gtx, judgeExpertsAppraisalDO, uWhere)
	})
	return
}

func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalService) DeleteEaJudgeExpertsAppraisalById(id int64, deleteBy int64) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaJudgeExpertsAppraisalDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

// 导出服务 -- start
type EaJudgeExpertsAppraisalExportStruct struct {
	Query *bizDto.EaJudgeExpertsAppraisalQuery
}

func (h *EaJudgeExpertsAppraisalExportStruct) Title() string {
	return "评委评价导出"
}

func (h *EaJudgeExpertsAppraisalExportStruct) HeaderColumn() []string {
	return []string{"评价ID,appraisal_id", "专家ID,expert_id", "任务ID,judge_task_id", "评委任务名称,judge_task_name", "评分日期,task_time_remark", "评价内容,appraisal_content", "删除标记：1-正常的；2-已删除的；,deleted", "备注,remark", "创建者,create_by", "创建时间,create_time", "更新者,update_by", "更新时间,update_time"}
}

func (h *EaJudgeExpertsAppraisalExportStruct) Rows(p *page.Page) []map[string]interface{} {
	return eaJudgeExpertsAppraisalDao.ListMapWithCondition(h.Query, p, nil, nil)
}

func (h *EaJudgeExpertsAppraisalExportStruct) Finish(url string) {
}

// 导出服务 -- end
func (judgeExpertsAppraisal *EaJudgeExpertsAppraisalService) ExportEaJudgeExpertsAppraisal(q *bizDto.EaJudgeExpertsAppraisalQuery) ([]byte, error) {
	return excel.NewExcel().ExportSync(&EaJudgeExpertsAppraisalExportStruct{q})
}

package bizSerImpl

import (
	"yx-exam-affair-gm/app/biz/bizDao"
	"yx-exam-affair-gm/app/biz/bizDao/bizDaoImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	"yx-exam-affair-gm/app/cst"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
	"github.com/sealsee/web-base/public/utils/file/excel"
)

/*
 * desc: 抽取评委明细服务层实现
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeDrawDetailService struct{}

var eaJudgeDrawDetailService *EaJudgeDrawDetailService
var eaJudgeDrawDetailDao bizDao.IEaJudgeDrawDetailDao

func init() {
	eaJudgeDrawDetailService = &EaJudgeDrawDetailService{}
	eaJudgeDrawDetailDao = bizDaoImpl.NewEaJudgeDrawDetailDao()
}

func NewEaJudgeDrawDetailService() *EaJudgeDrawDetailService {
	return eaJudgeDrawDetailService
}

func (judgeDrawDetail *EaJudgeDrawDetailService) ListEaJudgeDrawDetail(q *bizDto.EaJudgeDrawDetailQuery, p *page.Page) *bizDto.EaJudgeDrawDetailList {
	detailList := &bizDto.EaJudgeDrawDetailList{}

	// 统计当前任务下，当前部门所有评委的确认情况
	queryAllExperts := &bizDto.EaJudgeDrawDetailQuery{
		JudgeTaskId: q.JudgeTaskId,
		DeptId:      q.DeptId,
	}
	queryAllExperts.AddLikeAll("expert_name", q.ExpertName)
	judgeAllListInDept := judgeDrawDetail.ListAllEaJudgeDrawDetail(queryAllExperts, false)
	totalNeedConfirmNum := len(judgeAllListInDept)
	for _, detail := range judgeAllListInDept {
		if detail.ConfirmState == cst.ExpertConfirmStatusMap[cst.ATTEND] {
			detailList.AttendNum += 1
		} else if detail.ConfirmState == cst.ExpertConfirmStatusMap[cst.UNATTEND] {
			detailList.UnattendNum += 1
		}
	}
	detailList.NotSureNum = totalNeedConfirmNum - detailList.AttendNum - detailList.UnattendNum

	// 获取任务下评委列表
	q.AddLikeAll("expert_name", q.ExpertName)
	list := eaJudgeDrawDetailDao.List(q, p)
	if len(list) < 1 {
		return detailList
	}
	detailList.List = list

	// 根据id列表获取专家信息
	expertIds := make([]interface{}, 0)
	for _, v := range list {
		expertIds = append(expertIds, v.ExpertId)
	}
	if len(expertIds) < 1 {
		return detailList
	}
	query := &bizDto.EaJudgeExpertsQuery{}
	query.AddIn("expert_id", expertIds...)
	experts := eaJudgeExpertsService.ListAllEaJudgeExperts(query)
	expertMap := make(map[int64]*bizDo.EaJudgeExperts)
	for _, v := range experts {
		expertMap[v.ExpertId] = v
	}

	// 获取所有学院列表
	deptMap := eaDeptInfoService.GetAllDeptMap()

	// 组装专家数据
	for _, v := range list {
		expertInfo := expertMap[v.ExpertId]
		if expertInfo == nil {
			continue
		}
		v.GoodSubjects = expertInfo.GoodSubjects
		v.Sex = expertInfo.Sex
		v.Type = expertInfo.Type
		v.PhoneNo = expertInfo.PhoneNo
		if deptMap[v.DeptId] != nil {
			v.DeptName = deptMap[v.DeptId].DeptName
		} else {
			v.DeptName = "--"
		}
		v.TypeStr = v.Type
	}

	return detailList
}

func (judgeDrawDetail *EaJudgeDrawDetailService) ListAllEaJudgeDrawDetail(q *bizDto.EaJudgeDrawDetailQuery, needQueryTaskInfo bool) []*bizDo.EaJudgeDrawDetail {
	listAll := make([]*bizDo.EaJudgeDrawDetail, 0)
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE
	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		list := eaJudgeDrawDetailDao.List(q, pageInfo)
		if len(list) == 0 {
			break
		}
		listAll = append(listAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}
	if needQueryTaskInfo {
		taskIdList := make([]interface{}, 0)
		for _, v := range listAll {
			taskIdList = append(taskIdList, v.JudgeTaskId)
		}
		taskQuery := &bizDto.EaJudgeTaskInfoQuery{}
		taskQuery.AddIn("judge_task_id", taskIdList...)
		taskInfoList := eaJudgeTaskInfoService.ListAllEaJudgeTaskInfo(taskQuery)
		taskInfoMap := make(map[int64]*bizDo.EaJudgeTaskInfo)
		for _, v := range taskInfoList {
			taskInfoMap[v.JudgeTaskId] = v
		}
		for _, v := range listAll {
			taskInfo := taskInfoMap[v.JudgeTaskId]
			if taskInfo == nil {
				continue
			}
			v.JudgeTaskName = taskInfo.JudgeTaskName
		}
	}
	return listAll
}

func (judgeDrawDetail *EaJudgeDrawDetailService) UpdateEaJudgeDrawDetailById(judgeDrawDetailDO *bizDo.EaJudgeDrawDetail) (success bool, errInfo string) {

	// 获取当前评委信息
	id := judgeDrawDetailDO.Id
	dbJudgeDrawDetail := eaJudgeDrawDetailDao.GetById(id)
	if dbJudgeDrawDetail == nil {
		return false, "当前操作数据不存在"
	}
	if judgeDrawDetailDO.ConfirmState != cst.ExpertConfirmStatusMap[cst.NOT_CONFIRM] &&
		dbJudgeDrawDetail.ConfirmState != cst.ExpertConfirmStatusMap[cst.NOT_CONFIRM] {
		return false, "当前操作数据已确认，请重新选择"
	}
	if judgeDrawDetailDO.ConfirmState == cst.ExpertConfirmStatusMap[cst.NOT_CONFIRM] &&
		dbJudgeDrawDetail.ConfirmState == cst.ExpertConfirmStatusMap[cst.NOT_CONFIRM] {
		return false, "当前操作数据未确认，请重新选择"
	}
	isOnCampExpert := dbJudgeDrawDetail.DeptId != cst.BIZ_DEPTID_FOR_ADMISSIONS

	type deltaCounter struct {
		onCampAttend    int
		onCampUnattend  int
		offCampAttend   int
		offCampUnattend int
	}

	var delta deltaCounter

	if isOnCampExpert {
		if judgeDrawDetailDO.ConfirmState == cst.ExpertConfirmStatusMap[cst.NOT_CONFIRM] {
			if dbJudgeDrawDetail.ConfirmState == cst.ExpertConfirmStatusMap[cst.ATTEND] {
				delta.onCampAttend = -1
			} else if dbJudgeDrawDetail.ConfirmState == cst.ExpertConfirmStatusMap[cst.UNATTEND] {
				delta.onCampUnattend = -1
			}
		} else if judgeDrawDetailDO.ConfirmState == cst.ExpertConfirmStatusMap[cst.ATTEND] {
			delta.onCampAttend = 1
		} else if judgeDrawDetailDO.ConfirmState == cst.ExpertConfirmStatusMap[cst.UNATTEND] {
			delta.onCampUnattend = 1
		}
	} else {
		if judgeDrawDetailDO.ConfirmState == cst.ExpertConfirmStatusMap[cst.NOT_CONFIRM] {
			if dbJudgeDrawDetail.ConfirmState == cst.ExpertConfirmStatusMap[cst.ATTEND] {
				delta.offCampAttend = -1
			} else if dbJudgeDrawDetail.ConfirmState == cst.ExpertConfirmStatusMap[cst.UNATTEND] {
				delta.offCampUnattend = -1
			}
		} else if judgeDrawDetailDO.ConfirmState == cst.ExpertConfirmStatusMap[cst.ATTEND] {
			delta.offCampAttend = 1
		} else if judgeDrawDetailDO.ConfirmState == cst.ExpertConfirmStatusMap[cst.UNATTEND] {
			delta.offCampUnattend = 1
		}
	}

	uDate := &bizDo.EaJudgeDrawDetail{ConfirmState: judgeDrawDetailDO.ConfirmState}
	uWhere := &bizDo.EaJudgeDrawDetail{Id: id}

	tx.ExecGTx(func(gtx tx.Db) {
		success = eaJudgeDrawDetailDao.Update(gtx, uDate, uWhere)

		// 安全地构造 SQL 更新语句
		sqlBase := "on_campus_attend_num = on_campus_attend_num + ?, " +
			"on_campus_unattend_num = on_campus_unattend_num + ?, " +
			"off_campus_attend_num = off_campus_attend_num + ?, " +
			"off_campus_unattend_num = off_campus_unattend_num + ?, " +
			"update_by = ?, " + "update_time = now()"

		// 更新任务表
		sql2UpdateTask := "UPDATE ea_judge_task_info SET " + sqlBase + " WHERE judge_task_id = ? and deleted = 1"
		gtx.Exec(sql2UpdateTask,
			delta.onCampAttend, delta.onCampUnattend, delta.offCampAttend, delta.offCampUnattend, judgeDrawDetailDO.UpdateBy, dbJudgeDrawDetail.JudgeTaskId)

		// 更新记录表
		sql2UpdateRecord := "UPDATE ea_judge_draw_records SET " + sqlBase +
			" WHERE judge_task_id = ? AND subject_name = ? AND draw_times = ? and deleted = 1"
		gtx.Exec(sql2UpdateRecord,
			delta.onCampAttend, delta.onCampUnattend, delta.offCampAttend, delta.offCampUnattend,
			judgeDrawDetailDO.UpdateBy, dbJudgeDrawDetail.JudgeTaskId, dbJudgeDrawDetail.SubjectName, dbJudgeDrawDetail.DrawTimes)
	})

	return
}

func (judgeDrawDetail *EaJudgeDrawDetailService) DeleteEaJudgeDrawDetailById(id int64, deleteBy int64) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaJudgeDrawDetailDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

// 导出服务 -- start
type EaJudgeDrawDetailExportStruct struct {
	Query *bizDto.EaJudgeDrawDetailQuery
}

func (h *EaJudgeDrawDetailExportStruct) Title() string {
	return "抽取评委明细导出"
}

func (h *EaJudgeDrawDetailExportStruct) HeaderColumn() []string {
	return []string{"关联ID,id", "专家ID,expert_id", "专家姓名,expert_name", "任务ID,judge_task_id", "评分日期,task_time_remark", "抽取次数,draw_times", "科目编码 1-素描 2-速写 3-色彩,subject_code", "科目名称,subject_name", "确认状态 1-未应答 2-确认参加 3-不参加,confirm_state", "删除标记：1-正常的；2-已删除的；,deleted", "备注,remark", "创建者,create_by", "创建时间,create_time", "更新者,update_by", "更新时间,update_time"}
}

func (h *EaJudgeDrawDetailExportStruct) Rows(p *page.Page) []map[string]interface{} {
	return eaJudgeDrawDetailDao.ListMapWithCondition(h.Query, p, nil, nil)
}

func (h *EaJudgeDrawDetailExportStruct) Finish(url string) {
}

// 导出服务 -- end
func (judgeDrawDetail *EaJudgeDrawDetailService) ExportEaJudgeDrawDetail(q *bizDto.EaJudgeDrawDetailQuery) ([]byte, error) {
	return excel.NewExcel().ExportSync(&EaJudgeDrawDetailExportStruct{q})
}

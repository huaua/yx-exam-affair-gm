package bizSerImpl

import (
	"strconv"
	"time"
	"yx-exam-affair-gm/app/biz/bizDao"
	"yx-exam-affair-gm/app/biz/bizDao/bizDaoImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	"yx-exam-affair-gm/app/cst"
	"yx-exam-affair-gm/app/hint"

	"github.com/sealsee/web-base/public/ds/query"
	"github.com/sealsee/web-base/public/errs"

	"github.com/sealsee/web-base/public/basemodel"
	"github.com/sealsee/web-base/public/cst/common"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

/*
 * desc: 老师任务分配服务层实现
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchTaskReportService struct{}

var eaTchTaskReportService *EaTchTaskReportService
var eaTchTaskReportDao bizDao.IEaTchTaskReportDao

func init() {
	eaTchTaskReportService = &EaTchTaskReportService{}
	eaTchInfoService = &EaTchInfoService{}
	eaTchTaskDetailService = &EaTchTaskDetailService{}
	eaTchTaskReportDao = bizDaoImpl.NewEaTchTaskReportDao()
}

func NewEaTchTaskReportService() *EaTchTaskReportService {
	return eaTchTaskReportService
}

func (tchTaskReport *EaTchTaskReportService) GetEaTchTaskReportById(id int64) *bizDo.EaTchTaskReport {
	return eaTchTaskReportDao.GetById(id)
}

func (tchTaskReport *EaTchTaskReportService) ListEaTchTaskReport(q *bizDto.EaTchTaskReportQuery, p *page.Page) []*bizDto.EaTchTaskReportResultDTO {

	sql := `select a.report_id, a.tch_task_id, a.task_detail_id, a.task_time, a.tch_id, b.tch_name, b.job_no,
			b.id_card, b.education, b.invigilation_experience,
			(select count(1) from ea_arrange_tch ar where ar.deleted = 1 and ar.tch_id = a.tch_id) as invigilation_count,
			a.remark, a.dept_id, d.dept_name, a.state, b.sex, b.duties, b.prof_office
			from ea_tch_task_report a
			left join ea_tch_info b on a.tch_id = b.tch_id
			left join ea_dept_info d on a.dept_id = d.dept_id
			where `
	condition, args := (&basemodel.BaseEntityQuery{}).
		AddEq("a.tch_task_id", q.TchTaskId).
		AddEq("a.task_time", q.TaskTime).
		AddEq("a.deleted", 1).AddNot("a.state", cst.TEMP_SAVE).
		AddEq("a.state", q.State).
		AddEq("a.task_detail_id", q.TaskDetailId).
		AddEq("a.dept_id", q.DeptId).
		AddEq("b.sex", q.Sex).
		AddEq("b.invigilation_experience", q.InvigilationExperience).
		AddLikeAll("b.job_no", q.JobNo).
		AddLikeAll("b.tch_name", q.TchName).
		GetConditionArgs()
	listNew := query.RawSqlQueryListWithPage[bizDto.EaTchTaskReportResultDTO](p, sql+condition, args...)
	return listNew
}

func (tchTaskReport *EaTchTaskReportService) ListAllEaTchTaskReport(q *bizDto.EaTchTaskReportQuery) []*bizDo.EaTchTaskReport {
	listAll := make([]*bizDo.EaTchTaskReport, 0)
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE
	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		list := eaTchTaskReportDao.List(q, pageInfo)
		if len(list) == 0 {
			break
		}
		listAll = append(listAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}
	return listAll
}

func (tchTaskReport *EaTchTaskReportService) TempListEaTchTaskReport(deptId, tchTaskId int64) []*bizDo.EaTchTaskReportTmp {
	// 获取学院暂存的上报操作记录
	tchTaskReportQuery := &bizDto.EaTchTaskReportQuery{}
	tchTaskReportQuery.DeptId = deptId
	tchTaskReportQuery.TchTaskId = tchTaskId
	tchTaskReportQuery.State = cst.TEMP_SAVE
	tchTaskReportList := tchTaskReport.ListAllEaTchTaskReport(tchTaskReportQuery)

	reportTmpList := make([]*bizDo.EaTchTaskReportTmp, 0)
	reportTmpMap := make(map[int64]*bizDo.EaTchTaskReportTmp)
	for _, tchReport := range tchTaskReportList {
		if reportTmpMap[tchReport.TchId] == nil {
			reportTmp := &bizDo.EaTchTaskReportTmp{}
			reportTmp.TchId = tchReport.TchId
			reportTmp.TchName = tchReport.TchName
			reportTmp.Sex = tchReport.Sex
			reportTmp.JobNo = tchReport.JobNo
			reportTmp.Remark = tchReport.Remark
			taskTimes := make([]basemodel.BaseTime, 0)
			taskTimes = append(taskTimes, tchReport.TaskTime)
			reportTmp.TaskTimes = taskTimes
			reportTmpMap[tchReport.TchId] = reportTmp
		} else {
			reportMapTmp := reportTmpMap[tchReport.TchId]
			reportMapTmp.TaskTimes = append(reportMapTmp.TaskTimes, tchReport.TaskTime)
			// reportTmpMap[tchReport.TchId] = reportMapTmp // 上面引用生效
		}
	}
	for _, v := range reportTmpMap {
		reportTmpList = append(reportTmpList, v)
	}
	return reportTmpList
}

func (tchTaskReport *EaTchTaskReportService) TempSaveEaTchTaskReport(deptId int64, tchTaskReportList []*bizDo.EaTchTaskReport, optId int64) (bool, string) {
	tchTaskId := tchTaskReportList[0].ReportTaskDetailList[0].TchTaskId

	// 校验任务发布状态
	dbTchTaskInfo := eaTchTaskInfoDao.GetById(tchTaskId)
	if dbTchTaskInfo == nil {
		return false, "任务不存在，请检查！"
	}
	if dbTchTaskInfo.State != cst.PUBLISH {
		return false, "当前任务未发布，请先发布任务！"
	}

	// 遍历要插入的教师列表，如果时间冲突，则不允许上报
	tchTaskReportList2Insert := make([]*bizDo.EaTchTaskReport, 0)
	for _, tch := range tchTaskReportList {
		for _, taskTime := range tch.ReportTaskDetailList {
			temp := &bizDo.EaTchTaskReport{
				TchTaskId:    taskTime.TchTaskId,
				TaskDetailId: taskTime.TaskDetailId,
				TaskTime:     taskTime.TaskTime,
				TchId:        tch.TchId,
				TchName:      tch.TchName,
				Sex:          tch.Sex,
				JobNo:        tch.JobNo,
				Remark:       tch.Remark,
				DeptId:       deptId,
				State:        cst.TEMP_SAVE,
			}
			temp.Deleted = common.Normal
			temp.SetCreateBy(optId)
			temp.SetUpdateBy(optId)
			tchTaskReportList2Insert = append(tchTaskReportList2Insert, temp)
		}
	}

	success := tx.ExecGTx(func(gtx tx.GTx) {
		eaTchTaskReportDao.BatchInsert(gtx, tchTaskReportList2Insert)
	})

	return success, ""
}

func (tchTaskReport *EaTchTaskReportService) TempClearEaTchTaskReport(deptId int64, tchTaskReportList []*bizDo.EaTchTaskReport, optId int64) (bool, string) {
	tchTaskId := tchTaskReportList[0].ReportTaskDetailList[0].TchTaskId

	// 校验任务发布状态
	dbTchTaskInfo := eaTchTaskInfoDao.GetById(tchTaskId)
	if dbTchTaskInfo == nil {
		return false, "任务不存在，请检查！"
	}
	if dbTchTaskInfo.State != cst.PUBLISH {
		return false, "当前任务未发布，请先发布任务！"
	}

	// 遍历要插入的教师列表，如果时间冲突，则不允许上报
	tchTaskReportList2Clear := make([]*bizDo.EaTchTaskReport, 0)
	for _, tch := range tchTaskReportList {
		for _, taskTime := range tch.ReportTaskDetailList {
			temp := &bizDo.EaTchTaskReport{
				TchTaskId:    taskTime.TchTaskId,
				TaskDetailId: taskTime.TaskDetailId,
				TaskTime:     taskTime.TaskTime,
				TchId:        tch.TchId,
				TchName:      tch.TchName,
				Sex:          tch.Sex,
				JobNo:        tch.JobNo,
				Remark:       tch.Remark,
				DeptId:       deptId,
				State:        cst.TEMP_SAVE,
			}
			temp.Deleted = common.Normal
			temp.SetCreateBy(optId)
			temp.SetUpdateBy(optId)
			tchTaskReportList2Clear = append(tchTaskReportList2Clear, temp)
		}
	}

	success := tx.ExecGTx(func(gtx tx.GTx) {
		for _, clearTch := range tchTaskReportList2Clear {
			eaTchTaskReportDao.DeleteTempSavedTchReportOne(gtx, deptId, tchTaskId, clearTch.TaskDetailId, clearTch.TchId)
		}
	})

	return success, ""
}

func (tchTaskReport *EaTchTaskReportService) InsertEaTchTaskReport(tchTaskReportDO *bizDo.EaTchTaskReport) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaTchTaskReportDao.Insert(gtx, tchTaskReportDO)
	})
	return
}

func (tchTaskReport *EaTchTaskReportService) BatchInsertEaTchTaskReport(deptId int64, tchTaskReportList []*bizDo.EaTchTaskReport, optId int64) (bool, string) {
	tchTaskId := tchTaskReportList[0].ReportTaskDetailList[0].TchTaskId

	// 校验任务发布状态
	dbTchTaskInfo := eaTchTaskInfoDao.GetById(tchTaskId)
	if dbTchTaskInfo == nil {
		return false, "任务不存在，请检查！"
	}
	if dbTchTaskInfo.State != cst.PUBLISH {
		return false, "当前任务未发布，请先发布任务！"
	}

	// 获取任务下所有时间列表
	taskDetails := eaTchTaskDetailService.ListEaTchTaskDetail(&bizDto.EaTchTaskDetailQuery{DeptId: deptId, TchTaskId: tchTaskId}, page.NewPage())
	if len(taskDetails) <= 0 {
		return false, "当前任务未在任何时间设置对此学院的征集要求，请检查！"
	}
	timeList := make([]interface{}, 0)
	for _, v := range taskDetails {
		timeList = append(timeList, v.TaskTime)
	}

	// 获取时间列表下所有的教师列表
	tchAndTimeMap := make(map[string]bool)
	tchTaskReportQuery := &bizDto.EaTchTaskReportQuery{}
	tchTaskReportQuery.AddIn("task_time", timeList...)
	tchTaskReportQuery.AddNot("state", cst.TEMP_SAVE)
	tchTaskHasReportList := tchTaskReport.ListAllEaTchTaskReport(tchTaskReportQuery)
	if len(tchTaskHasReportList) > 0 {
		for _, v := range tchTaskHasReportList {
			tchAndTimeMap[strconv.FormatInt(v.TchId, 10)+"_"+v.TaskTime.String()] = true
		}
	}

	// 遍历要插入的教师列表，如果时间冲突，则不允许上报
	tchTaskReportList2Insert := make([]*bizDo.EaTchTaskReport, 0)
	for _, tch := range tchTaskReportList {
		for _, taskTime := range tch.ReportTaskDetailList {
			if tchAndTimeMap[strconv.FormatInt(tch.TchId, 10)+"_"+taskTime.TaskTime.String()] {
				return false, "监考老师：" + tch.TchName + "，在时间：" + taskTime.TaskTime.String() + "，已被分配任务，请检查！"
			}
			temp := &bizDo.EaTchTaskReport{
				TchTaskId:    taskTime.TchTaskId,
				TaskDetailId: taskTime.TaskDetailId,
				TaskTime:     taskTime.TaskTime,
				TchId:        tch.TchId,
				TchName:      tch.TchName,
				Sex:          tch.Sex,
				JobNo:        tch.JobNo,
				Remark:       tch.Remark,
				DeptId:       deptId,
				State:        cst.TCH_REPORT_PENDING_AUDIT,
			}
			temp.Deleted = common.Normal
			temp.SetCreateBy(optId)
			temp.SetUpdateBy(optId)
			tchTaskReportList2Insert = append(tchTaskReportList2Insert, temp)
		}
	}

	// 将类型按照任务时间分组计数
	tchTaskDetailNumMap := make(map[int64]int)
	for _, v := range tchTaskReportList2Insert {
		tchTaskDetailNumMap[v.TaskDetailId]++
	}

	success := tx.ExecGTx(func(gtx tx.GTx) {
		// 先删除暂存记录 state='A'的记录
		eaTchTaskReportDao.DeleteTempSavedTchReport(gtx, deptId, tchTaskId)
		// 再批量插入
		eaTchTaskReportDao.BatchInsert(gtx, tchTaskReportList2Insert)
		// 循环更新每一天的计数
		for key, value := range tchTaskDetailNumMap {
			eaTchTaskDetailService.UpdateReportedTchNum(key, value)
		}
	})

	return success, ""
}

func (tchTaskReport *EaTchTaskReportService) UpdateEaTchTaskReportById(tchTaskReportDO *bizDo.EaTchTaskReport) (success bool) {
	uWhere := &bizDo.EaTchTaskReport{ReportId: tchTaskReportDO.ReportId}
	tchTaskReportDO.ReportId = 0
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaTchTaskReportDao.Update(gtx, tchTaskReportDO, uWhere)
	})
	return
}

func (tchTaskReport *EaTchTaskReportService) UpdateEaTchTaskReportByTchId(tchTaskReportDO *bizDo.EaTchTaskReport) (success bool) {
	uWhere := &bizDo.EaTchTaskReport{TchId: tchTaskReportDO.TchId}
	tchTaskReportDO.TchId = 0
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaTchTaskReportDao.Update(gtx, tchTaskReportDO, uWhere)
	})
	return
}

func (tchTaskReport *EaTchTaskReportService) EditEaTchTaskReport(editInfo *bizDto.EditTchTaskReportQuery, optId int64) (bool, string) {
	// 获取详情信息
	dbTaskDetail := eaTchTaskDetailService.GetEaTchTaskDetailById(editInfo.TaskDetailId)
	if dbTaskDetail == nil {
		return false, "获取任务详情失败，请检查！"
	}

	// 获取当前日期所有上报过的教职工列表
	tchAndTimeMap := make(map[int64]bool)
	tchTaskReportQuery := &bizDto.EaTchTaskReportQuery{TaskTime: dbTaskDetail.TaskTime}
	tchTaskReportQuery.AddNot("state", cst.TEMP_SAVE)
	tchTaskHasReportList := tchTaskReport.ListAllEaTchTaskReport(tchTaskReportQuery)
	if len(tchTaskHasReportList) > 0 {
		for _, v := range tchTaskHasReportList {
			tchAndTimeMap[v.TchId] = true
		}
	}

	// 数据分组，并进行相应的判断
	thisTimeReportedTchIdMap := make(map[int64]bool)
	for _, tch := range editInfo.TchList {
		thisTimeReportedTchIdMap[tch.TchId] = true
	}

	// 组装待新增的列表
	list2add := make([]*bizDo.EaTchTaskReport, 0)
	for _, tch := range editInfo.TchList {
		// 检查日期和已上报的是否冲突
		if tchAndTimeMap[tch.TchId] {
			return false, tch.TchName + "在【" + dbTaskDetail.TaskTime.String() + "】，已被上报其他任务，请检查！"
		}

		temp := &bizDo.EaTchTaskReport{
			TchTaskId:    dbTaskDetail.TchTaskId,
			TaskDetailId: dbTaskDetail.TaskDetailId,
			TaskTime:     dbTaskDetail.TaskTime,
			TchId:        tch.TchId,
			TchName:      tch.TchName,
			Sex:          tch.Sex,
			JobNo:        tch.JobNo,
			Remark:       tch.Remark,
			DeptId:       editInfo.DeptId,
			State:        cst.TCH_REPORT_PENDING_AUDIT,
		}
		temp.Deleted = common.Normal
		temp.SetCreateBy(optId)
		temp.SetUpdateBy(optId)
		list2add = append(list2add, temp)

	}

	success := tx.ExecGTx(func(gtx tx.Db) {
		if len(list2add) > 0 {
			eaTchTaskReportDao.BatchInsert(gtx, list2add)
		}
		if len(list2add) > 0 {
			eaTchTaskDetailService.UpdateReportedTchNum(editInfo.TaskDetailId, len(list2add))
		}
	})

	return success, ""
}

func (tchTaskReport *EaTchTaskReportService) DeleteEaTchTaskReportById(id int64, deleteBy int64) (success bool, err errs.ERROR) {
	dbTchTaskReport := eaTchTaskReportDao.GetById(id)
	success = false
	if dbTchTaskReport == nil || dbTchTaskReport.TaskDetailId < 1 {
		err = hint.BIZ_TCH_TASK_NOT_EXIST_ERR
		return
	}

	if dbTchTaskReport.State == cst.TCH_REPORT_ASSIGNED {
		err = hint.BIZ_TCH_HAS_ARRANGED_CANCEL_ERR
		return
	}

	dbTaskDetail := eaTchTaskDetailService.GetEaTchTaskDetailById(dbTchTaskReport.TaskDetailId)
	if dbTaskDetail == nil {
		err = hint.BIZ_TCH_TASK_DETAIL_NOT_EXIST_ERR
		return
	}

	if dbTaskDetail.ForceFlag == cst.YES1 && dbTaskDetail.CollectNum <= dbTaskDetail.NeedNum {
		err = hint.BIZ_NEED_REPLACE_TCH_ERR
		return
	}
	success = tx.ExecGTx(func(gtx tx.GTx) {
		eaTchTaskReportDao.DeleteById(gtx, id, deleteBy)
		eaTchTaskDetailService.UpdateReportedTchNum(dbTaskDetail.TaskDetailId, -1)
	})
	return
}

func (tchTaskReport *EaTchTaskReportService) ExportEaTchTaskReport(q *bizDto.EaTchTaskReportQuery) ([]byte, error) {
	return exportReportTchData(q)
}

func (tchTaskReport *EaTchTaskReportService) ListEnabledEaTchInfo(q *bizDto.EaTchTaskReportQuery) ([]*bizDo.EaTchInfo, string) {

	// 获取所有教职工列表
	tchQuery := &bizDto.EaTchInfoQuery{
		DeptId:        q.DeptId,
		EnabledStatus: cst.TCH_ENABLED,
		AuditStatus:   cst.TCH_AUDIT_PASSED,
	}
	var allTchInfo []*bizDo.EaTchInfo
	if len(q.TchNameOrJobNo) > 0 {
		likeCondition := "%" + q.TchNameOrJobNo + "%"
		condition := "tch_name like ? or job_no like ?"
		allTchInfo = eaTchInfoService.ListAllEaTchInfo(tchQuery, condition, likeCondition, likeCondition)
	} else {
		allTchInfo = eaTchInfoService.ListAllEaTchInfo(tchQuery, nil, nil)
	}

	if len(allTchInfo) <= 0 {
		return allTchInfo, ""
	}

	// 获取指定任务下的所有日期列表
	timeList := make([]interface{}, 0)
	if !q.TaskTime.IsZero() {
		timeList = append(timeList, q.TaskTime)
	} else {
		// 一个教职工征集任务的时间跨度最大为7天，所有这里就写死页数了
		taskDetails := eaTchTaskDetailService.ListEaTchTaskDetail(&bizDto.EaTchTaskDetailQuery{
			TchTaskId: q.TchTaskId,
			DeptId:    q.DeptId,
		}, page.NewPage())
		if len(taskDetails) <= 0 {
			return nil, "当前任务未在任何时间设置对此学院的征集要求，请检查！"
		}
		for _, v := range taskDetails {
			timeList = append(timeList, v.TaskTime)
		}
	}

	// 根据时间列表，查询所有在时间范围内已经上报的教职工列表
	tchUsedTimeListMap := make(map[int64][]basemodel.BaseTime)
	query := &bizDto.EaTchTaskReportQuery{}
	query.AddIn("task_time", timeList...)
	tchTaskReportList := tchTaskReport.ListAllEaTchTaskReport(query)
	if len(tchTaskReportList) <= 0 {
		return allTchInfo, ""
	}
	tchTaskTmpSavedMap := make(map[int64]*bizDo.EaTchTaskReport)
	for _, v := range tchTaskReportList {
		if tchUsedTimeListMap[v.TchId] == nil {
			tchUsedTimeListMap[v.TchId] = make([]basemodel.BaseTime, 0)
		}
		tchUsedTimeListMap[v.TchId] = append(tchUsedTimeListMap[v.TchId], v.TaskTime)
		if v.State == cst.TEMP_SAVE && v.TchTaskId == q.TchTaskId {
			tchTaskTmpSavedMap[v.TchId] = v
		}
	}

	// 为每个教师设置哪些时间段已上报
	filterTchInfo := make([]*bizDo.EaTchInfo, 0)
	for _, v := range allTchInfo {
		if tchUsedTimeListMap[v.TchId] != nil {

			// 如果任务下的日期数量和上报的日期数量一致，则需要过滤该条数据
			if len(tchUsedTimeListMap[v.TchId]) >= len(timeList) {
				continue
			}
			// 如果暂存列表有这个老师，也过滤该条记录
			if tchTaskTmpSavedMap[v.TchId] != nil {
				continue
			}

			v.ReportTaskTimeList = tchUsedTimeListMap[v.TchId]
		}
		filterTchInfo = append(filterTchInfo, v)
	}
	return filterTchInfo, ""
}

func (tchTaskReport *EaTchTaskReportService) CheckCancelResult(id int64) (success bool, err errs.ERROR) {
	success = false
	dbTchTaskReport := eaTchTaskReportDao.GetById(id)
	if dbTchTaskReport == nil || dbTchTaskReport.TaskDetailId < 1 {
		err = hint.BIZ_TCH_TASK_NOT_EXIST_ERR
		return
	}

	if dbTchTaskReport.State == cst.TCH_REPORT_ASSIGNED {
		err = hint.BIZ_NEED_REPLACE_TCH_ERR
		return
	}

	dbTaskDetail := eaTchTaskDetailService.GetEaTchTaskDetailById(dbTchTaskReport.TaskDetailId)
	if dbTaskDetail == nil {
		err = hint.BIZ_TCH_TASK_DETAIL_NOT_EXIST_ERR
		return
	}

	if dbTaskDetail.ForceFlag == cst.YES1 && dbTaskDetail.CollectNum <= dbTaskDetail.NeedNum {
		err = hint.BIZ_NEED_REPLACE_TCH_ERR
		return
	}

	success = true
	return
}

func (tchTaskReport *EaTchTaskReportService) ReplaceReport(tchTaskReport2Replace *bizDo.EaTchTaskReport) (success bool, err errs.ERROR) {

	if tchTaskReport2Replace.ReportId < 1 || tchTaskReport2Replace.TchId < 1 || tchTaskReport2Replace.TaskDetailId < 1 {
		err = errs.PARAM_INVALID_ERR
		return
	}

	dbTchInfo := eaTchInfoService.GetEaTchInfoById(tchTaskReport2Replace.TchId)
	if dbTchInfo == nil {
		err = hint.BIZ_TCH_NOT_EXIST_ERR
		return
	}

	dbTaskDetail := eaTchTaskDetailDao.GetById(tchTaskReport2Replace.TaskDetailId)
	if dbTaskDetail == nil || dbTaskDetail.TaskTime.IsZero() {
		err = hint.BIZ_TCH_TASK_DETAIL_NOT_EXIST_ERR
		return
	}

	dbTchReport := eaTchTaskReportService.GetEaTchTaskReportById(tchTaskReport2Replace.ReportId)
	if dbTchReport == nil {
		err = hint.BIZ_TCH_REPORT_NOT_EXIST_ERR.Format(tchTaskReport2Replace.ReportId)
		return
	}

	// 校验日期是否冲突
	exists := eaTchTaskReportDao.Count(&bizDto.EaTchTaskReportQuery{
		TchId:    tchTaskReport2Replace.TchId,
		TaskTime: dbTaskDetail.TaskTime,
	}) > 0
	if exists {
		err = hint.BIZ_TCH_HAS_REPORT_THIS_TIME_ERR.Format(dbTchInfo.TchName, dbTchInfo.JobNo, dbTaskDetail.TaskTime.FormatString(time.DateOnly))
		return
	}

	// 替换采用“原记录软删除 + 新增替换记录”的方式，避免覆盖原老师及原审核状态，
	// 并通过 replaced_from 保留完整替换链路。
	newReport := &bizDo.EaTchTaskReport{
		TchTaskId:    dbTchReport.TchTaskId,
		TaskDetailId: dbTchReport.TaskDetailId,
		TaskTime:     dbTchReport.TaskTime,
		TchId:        dbTchInfo.TchId,
		TchName:      dbTchInfo.TchName,
		Sex:          dbTchInfo.Sex,
		JobNo:        dbTchInfo.JobNo,
		DeptId:       dbTchReport.DeptId,
		State:        cst.TCH_REPORT_PENDING_AUDIT,
		ReplacedFrom: dbTchReport.ReportId,
		Remark:       dbTchReport.Remark,
	}
	newReport.SetCreateBy(tchTaskReport2Replace.UpdateBy)
	newReport.SetUpdateBy(tchTaskReport2Replace.UpdateBy)
	oldTchId := dbTchReport.TchId
	newTchId := newReport.TchId
	success = tx.ExecGTx(func(gtx tx.GTx) {
		if !eaTchTaskReportDao.DeleteById(gtx, dbTchReport.ReportId, tchTaskReport2Replace.UpdateBy) {
			panic("soft delete replaced teacher report failed")
		}
		if !eaTchTaskReportDao.Insert(gtx, newReport) {
			panic("insert replacement teacher report failed")
		}

		// 如果已经编排，当前有效编排需改为引用新上报记录；历史编排记录不修改。
		if dbTchReport.State == cst.TCH_REPORT_ASSIGNED {
			arrangeTch := &bizDo.EaArrangeTch{
				TchReportId: newReport.ReportId,
				TchId:       newReport.TchId,
				TchName:     newReport.TchName,
			}
			arrangeTch.SetUpdateBy(newReport.UpdateBy)
			uArrangeTchWhere := &bizDo.EaArrangeTch{TchReportId: dbTchReport.ReportId}
			if !eaArrangeTchDao.Update(gtx, arrangeTch, uArrangeTchWhere) {
				panic("update replacement teacher arrangement failed")
			}
		}

		// 重新计算原/新教师的监考经验
		if oldTchId > 0 {
			eaTchInfoDao.RefreshInvigilationExperienceByTchId(gtx, oldTchId)
		}
		if newTchId > 0 && newTchId != oldTchId {
			eaTchInfoDao.RefreshInvigilationExperienceByTchId(gtx, newTchId)
		}
	})
	return
}

/*
 * desc: 招办审核征集状态
 * action="pass" → TCH_REPORT_PENDING_ASSIGN（待分配）
 * action="reject" → TCH_REPORT_REJECTED（不通过）
 * 仅当状态为 已上报,待审核（"1"）时可审核
 */
func (tchTaskReport *EaTchTaskReportService) AuditTchTaskReport(reportId int64, action string, optId int64) (success bool, err errs.ERROR) {
	success = false
	if reportId < 1 {
		err = errs.PARAM_INVALID_ERR
		return
	}
	if action != cst.TCH_AUDIT_ACTION_PASS && action != cst.TCH_AUDIT_ACTION_REJECT {
		err = errs.PARAM_INVALID_ERR
		return
	}

	dbReport := eaTchTaskReportDao.GetById(reportId)
	if dbReport == nil {
		err = hint.BIZ_TCH_REPORT_NOT_EXIST_ERR.Format(reportId)
		return
	}
	if dbReport.State != cst.TCH_REPORT_PENDING_AUDIT {
		err = hint.BIZ_TCH_HAS_ARRANGED_CANCEL_ERR
		return
	}

	var newState string
	if action == cst.TCH_AUDIT_ACTION_PASS {
		newState = cst.TCH_REPORT_PENDING_ASSIGN
	} else {
		newState = cst.TCH_REPORT_REJECTED
	}

	uData := &bizDo.EaTchTaskReport{State: newState}
	uData.SetUpdateBy(optId)
	uWhere := &bizDo.EaTchTaskReport{ReportId: reportId}
	success = tx.ExecGTx(func(gtx tx.GTx) {
		eaTchTaskReportDao.Update(gtx, uData, uWhere)
	})
	return
}

/*
 * desc: 招办重置征集状态
 * 仅当状态为 待分配（"2"）/ 不通过（"3"）时可重置回 已上报,待审核（"1"）
 */
func (tchTaskReport *EaTchTaskReportService) ResetTchTaskReport(reportId int64, optId int64) (success bool, err errs.ERROR) {
	success = false
	if reportId < 1 {
		err = errs.PARAM_INVALID_ERR
		return
	}

	dbReport := eaTchTaskReportDao.GetById(reportId)
	if dbReport == nil {
		err = hint.BIZ_TCH_REPORT_NOT_EXIST_ERR.Format(reportId)
		return
	}
	if dbReport.State != cst.TCH_REPORT_PENDING_ASSIGN && dbReport.State != cst.TCH_REPORT_REJECTED {
		err = hint.BIZ_TCH_HAS_ARRANGED_CANCEL_ERR
		return
	}

	uData := &bizDo.EaTchTaskReport{State: cst.TCH_REPORT_PENDING_AUDIT}
	uData.SetUpdateBy(optId)
	uWhere := &bizDo.EaTchTaskReport{ReportId: reportId}
	success = tx.ExecGTx(func(gtx tx.GTx) {
		eaTchTaskReportDao.Update(gtx, uData, uWhere)
	})
	return
}

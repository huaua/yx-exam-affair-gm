package bizSerImpl

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"yx-exam-affair-gm/app/biz/bizDao"
	"yx-exam-affair-gm/app/biz/bizDao/bizDaoImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	"yx-exam-affair-gm/app/cst"

	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/ds"
	"github.com/sealsee/web-base/public/utils/redis"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"

	"gorm.io/gorm"
)

/*
 * desc: 评委征集任务服务层实现
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeTaskInfoService struct{}

var eaJudgeTaskInfoService *EaJudgeTaskInfoService
var eaJudgeTaskInfoDao bizDao.IEaJudgeTaskInfoDao

func init() {
	eaJudgeTaskInfoService = &EaJudgeTaskInfoService{}
	eaJudgeTaskInfoDao = bizDaoImpl.NewEaJudgeTaskInfoDao()
}

func NewEaJudgeTaskInfoService() *EaJudgeTaskInfoService {
	return eaJudgeTaskInfoService
}

func (judgeTaskInfo *EaJudgeTaskInfoService) GetEaJudgeTaskInfoById(id int64) []*bizDo.EaJudgeDrawRecords {
	records := eaJudgeDrawRecordsService.ListAllEaJudgeDrawRecords(&bizDto.EaJudgeDrawRecordsQuery{JudgeTaskId: id})
	for _, v := range records {
		if len(v.SubjectName) < 1 {
			v.SubjectName = cst.ALL
		}
	}
	return records
}

func (judgeTaskInfo *EaJudgeTaskInfoService) ListEaJudgeTaskInfo(q *bizDto.EaJudgeTaskInfoQuery, p *page.Page) []*bizDo.EaJudgeTaskInfo {
	if len(q.JudgeTaskName) > 0 {
		q.AddLikeAll("judge_task_name", q.JudgeTaskName)
		q.JudgeTaskName = ""
	}
	list := eaJudgeTaskInfoDao.List(q, p)
	// 收集所有任务ID，批量查询上报数
	taskIds := make([]int64, 0, len(list))
	for _, v := range list {
		taskIds = append(taskIds, v.JudgeTaskId)
	}
	reportCountMap := make(map[int64]int64, len(taskIds))
	if len(taskIds) > 0 {
		var results []struct {
			TaskId int64
			Cnt    int64
		}
		ds.GetGDB().Table("ea_judge_expert_report").
			Select("judge_task_id as task_id, count(1) as cnt").
			Where("judge_task_id IN ? AND submit_state = 2 AND deleted = 1", taskIds).
			Group("judge_task_id").
			Scan(&results)
		for _, r := range results {
			reportCountMap[r.TaskId] = r.Cnt
		}
	}
	for _, v := range list {
		v.JudgeTaskTypeStr = cst.GetExpertTypeLabelByValue(v.JudgeTaskType)
		v.StateStr = cst.GetJudgeTaskStatusLabelByValue(v.State)
		v.ReportCount = int(reportCountMap[v.JudgeTaskId])
	}

	// 直接抽取模式：批量查询 ea_judge_draw_records 表的 required_num 总和
	drawCountMap := make(map[int64]int64, len(taskIds))
	// 上报模式：批量查询 ea_judge_task_report_requirement 表的 expert_num 总和
	requirementCountMap := make(map[int64]int64, len(taskIds))
	if len(taskIds) > 0 {
		var drawResults []struct {
			TaskId int64
			Cnt    int64
		}
		ds.GetGDB().Table("ea_judge_draw_records").
			Select("judge_task_id as task_id, COALESCE(SUM(required_num), 0) as cnt").
			Where("judge_task_id IN ? AND deleted = 1", taskIds).
			Group("judge_task_id").
			Scan(&drawResults)
		for _, r := range drawResults {
			drawCountMap[r.TaskId] = r.Cnt
		}

		var reqResults []struct {
			TaskId int64
			Cnt    int64
		}
		ds.GetGDB().Table("ea_judge_task_report_requirement").
			Select("judge_task_id as task_id, COALESCE(SUM(expert_num), 0) as cnt").
			Where("judge_task_id IN ? AND deleted = 1", taskIds).
			Group("judge_task_id").
			Scan(&reqResults)
		for _, r := range reqResults {
			requirementCountMap[r.TaskId] = r.Cnt
		}
	}
	// 直接抽取模式：批量查询实际已抽取的评分专家数
	drawExpertCountMap := make(map[int64]int64, len(taskIds))
	if len(taskIds) > 0 {
		var drawExpertResults []struct {
			TaskId int64
			Cnt    int64
		}
		ds.GetGDB().Table("ea_judge_draw_selection").
			Select("judge_task_id as task_id, count(1) as cnt").
			Where("judge_task_id IN ? AND expert_role = 1 AND submit_state = 2 AND deleted = 1", taskIds).
			Group("judge_task_id").
			Scan(&drawExpertResults)
		for _, r := range drawExpertResults {
			drawExpertCountMap[r.TaskId] = r.Cnt
		}
	}

	for _, v := range list {
		if v.TaskMode == 1 {
			v.RequirementCount = int(drawCountMap[v.JudgeTaskId])
			v.DrawCount = int(drawExpertCountMap[v.JudgeTaskId])
		} else if v.TaskMode == 2 {
			v.RequirementCount = int(requirementCountMap[v.JudgeTaskId])
		}
	}
	return list
}

func (judgeTaskInfo *EaJudgeTaskInfoService) ListAllEaJudgeTaskInfo(q *bizDto.EaJudgeTaskInfoQuery) []*bizDo.EaJudgeTaskInfo {
	listAll := make([]*bizDo.EaJudgeTaskInfo, 0)
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE
	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		list := eaJudgeTaskInfoDao.List(q, pageInfo)
		if len(list) == 0 {
			break
		}
		listAll = append(listAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}
	return listAll
}

func (judgeTaskInfo *EaJudgeTaskInfoService) InsertEaJudgeTaskInfo(judgeTaskInfoDO *bizDo.EaJudgeTaskInfo) (success bool, errInfo string) {

	// 参数必填校验
	if utf8.RuneCountInString(judgeTaskInfoDO.JudgeTaskName) < 1 || utf8.RuneCountInString(judgeTaskInfoDO.JudgeTaskName) > 50 ||
		judgeTaskInfoDO.TaskYear < 1900 || judgeTaskInfoDO.TaskMode < 1 || judgeTaskInfoDO.TaskMode > 2 ||
		!cst.ExpertValueIsValid(judgeTaskInfoDO.JudgeTaskType) ||
		judgeTaskInfoDO.TaskStartTime.IsZero() || judgeTaskInfoDO.TaskEndTime.IsZero() {
		return false, "任务年份、任务名称（1-50字符）、任务类型、任务模式和评分日期均须正确填写"
	}
	if judgeTaskInfoDO.TaskMode == 1 && (judgeTaskInfoDO.ScoreRounds < 1 || judgeTaskInfoDO.SubjectCount < 1 || judgeTaskInfoDO.ExpertsPerRound < 1) {
		return false, "评分轮次、科目数量、每轮专家数必须为大于0的正整数"
	}

	// 任务名称不允许重复
	exit := eaJudgeTaskInfoDao.Count(&bizDto.EaJudgeTaskInfoQuery{JudgeTaskName: judgeTaskInfoDO.JudgeTaskName, TaskYear: judgeTaskInfoDO.TaskYear}) > 0
	if exit {
		return false, "任务名称已经存在"
	}

	if cst.YES != judgeTaskInfoDO.DrawSubject {
		judgeTaskInfoDO.DrawSubject = cst.NO
	}

	judgeTaskInfoDO.State = cst.JudgeTaskStatusMap[cst.NOT_START]
	if judgeTaskInfoDO.TaskMode == 2 {
		judgeTaskInfoDO.PublishState = 1
	} else {
		judgeTaskInfoDO.PublishState = 2
	}
	judgeTaskInfoDO.DrawTimes = 0
	txSuccess := tx.ExecGTx(func(gtx tx.GTx) {
		success = eaJudgeTaskInfoDao.Insert(gtx, judgeTaskInfoDO)
		if success && judgeTaskInfoDO.TaskMode == 1 {
			success = insertMissingDrawRecords(gtx, judgeTaskInfoDO, 1, judgeTaskInfoDO.ScoreRounds, judgeTaskInfoDO.CreateBy)
		}
	})
	if !txSuccess {
		success = false
	}
	return
}

// insertMissingDrawRecords 按“评分轮次 × 科目数量”补齐评委抽取规则。
// 科目名称由招办后续在评委抽取管理中填写，新增时仅预置每科所需专家数。
func insertMissingDrawRecords(gtx tx.GTx, task *bizDo.EaJudgeTaskInfo, startRound int, endRound int, optId int64) bool {
	if task == nil || task.JudgeTaskId < 1 || task.SubjectCount < 1 || task.ExpertsPerRound < 1 || startRound < 1 || endRound < startRound {
		return false
	}
	records := make([]*bizDo.EaJudgeDrawRecords, 0, (endRound-startRound+1)*task.SubjectCount)
	for round := startRound; round <= endRound; round++ {
		var existing int64
		gtx.Raw("SELECT COUNT(1) FROM ea_judge_draw_records WHERE judge_task_id = ? AND draw_times = ? AND deleted = 1", task.JudgeTaskId, round).Scan(&existing)
		for subjectIndex := int(existing); subjectIndex < task.SubjectCount; subjectIndex++ {
			record := &bizDo.EaJudgeDrawRecords{
				JudgeTaskId: task.JudgeTaskId,
				DrawTimes:   round,
				SubjectName: "",
				RequiredNum: task.ExpertsPerRound,
			}
			record.Deleted = common.Normal
			record.SetCreateBy(optId)
			record.SetUpdateBy(optId)
			records = append(records, record)
		}
	}
	if len(records) == 0 {
		return true
	}
	return eaJudgeDrawRecordsDao.BatchInsert(gtx, records)
}

func (judgeTaskInfo *EaJudgeTaskInfoService) UpdateEaJudgeTaskInfoById(judgeTaskInfoDO *bizDo.EaJudgeTaskInfo) (success bool, errInfo string) {
	// 状态未未开始才能修改
	dbTaskInfo := eaJudgeTaskInfoDao.GetById(judgeTaskInfoDO.JudgeTaskId)
	if dbTaskInfo == nil || dbTaskInfo.State != cst.JudgeTaskStatusMap[cst.NOT_START] {
		return false, "当前任务状态已经进行抽取动作，不能修改"
	}

	// 任务名称不允许重复
	if utf8.RuneCountInString(judgeTaskInfoDO.JudgeTaskName) < 1 || utf8.RuneCountInString(judgeTaskInfoDO.JudgeTaskName) > 50 || judgeTaskInfoDO.TaskStartTime.IsZero() || judgeTaskInfoDO.TaskEndTime.IsZero() {
		return false, "只能修改任务名称（1-50字符）和评分日期"
	}
	if dbTaskInfo.JudgeTaskName != judgeTaskInfoDO.JudgeTaskName {
		exit := eaJudgeTaskInfoDao.Count(&bizDto.EaJudgeTaskInfoQuery{JudgeTaskName: judgeTaskInfoDO.JudgeTaskName, TaskYear: dbTaskInfo.TaskYear}) > 0
		if exit {
			return false, "任务名称已经存在"
		}
	}

	uWhere := &bizDo.EaJudgeTaskInfo{JudgeTaskId: judgeTaskInfoDO.JudgeTaskId}
	uData := &bizDo.EaJudgeTaskInfo{JudgeTaskName: judgeTaskInfoDO.JudgeTaskName, TaskStartTime: judgeTaskInfoDO.TaskStartTime, TaskEndTime: judgeTaskInfoDO.TaskEndTime}
	uData.SetUpdateBy(judgeTaskInfoDO.UpdateBy)
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaJudgeTaskInfoDao.Update(gtx, uData, uWhere)
	})
	return
}

func (judgeTaskInfo *EaJudgeTaskInfoService) DeleteEaJudgeTaskInfoById(id int64, deleteBy int64) (success bool, errInfo string) {
	dbTaskInfo := eaJudgeTaskInfoDao.GetById(id)
	if dbTaskInfo == nil || dbTaskInfo.State != cst.JudgeTaskStatusMap[cst.NOT_START] {
		return false, "当前任务状态已经进行抽取动作，不能删除"
	}
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaJudgeTaskInfoDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

func (judgeTaskInfo *EaJudgeTaskInfoService) AppendScoreRounds(id int64, rounds int, optId int64) (success bool, errInfo string) {
	dbTask := eaJudgeTaskInfoDao.GetById(id)
	if dbTask == nil || dbTask.TaskMode != 1 {
		return false, "只有直接抽取模式可以追加轮次"
	}
	if rounds < 1 {
		return false, "追加轮次必须为大于0的正整数"
	}
	newTotalRounds := dbTask.ScoreRounds + rounds
	uData := &bizDo.EaJudgeTaskInfo{ScoreRounds: newTotalRounds}
	uData.SetUpdateBy(optId)
	txSuccess := tx.ExecGTx(func(gtx tx.GTx) {
		success = eaJudgeTaskInfoDao.Update(gtx, uData, &bizDo.EaJudgeTaskInfo{JudgeTaskId: id})
		if success {
			// 同时补齐历史缺失轮次，并插入本次追加的轮次科目数据。
			dbTask.ScoreRounds = newTotalRounds
			success = insertMissingDrawRecords(gtx, dbTask, 1, newTotalRounds, optId)
		}
	})
	if !txSuccess {
		success = false
	}
	return
}

func (judgeTaskInfo *EaJudgeTaskInfoService) PublishCollegeReportTask(id int64, optId int64) (success bool, errInfo string) {
	dbTask := eaJudgeTaskInfoDao.GetById(id)
	if dbTask == nil || dbTask.TaskMode != 2 {
		return false, "只有学院上报模式可以发布"
	}
	if dbTask.PublishState == 2 {
		return false, "任务已发布"
	}
	var count int64
	ds.GetGDB().Table("ea_judge_task_report_requirement").Where("judge_task_id = ? AND deleted = 1", id).Count(&count)
	if count < 1 {
		return false, "请先导入学院上报要求后再发布"
	}
	uData := &bizDo.EaJudgeTaskInfo{PublishState: 2}
	uData.SetUpdateBy(optId)
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaJudgeTaskInfoDao.Update(gtx, uData, &bizDo.EaJudgeTaskInfo{JudgeTaskId: id})
	})
	return
}

func (judgeTaskInfo *EaJudgeTaskInfoService) ExtractJudgeTask(extractJudgeTask *bizDto.ExtractJudgeTaskDTO, optId int64) (success bool, errInfo string) {
	// 状态检测
	dbJudgeTaskInfo := eaJudgeTaskInfoDao.GetById(extractJudgeTask.JudgeTaskId)
	if dbJudgeTaskInfo == nil {
		return false, "任务不存在，请检查！"
	}
	if dbJudgeTaskInfo.DrawSubject == cst.YES {
		return false, "当前任务设置为按擅长科目抽取，不能进行当前抽取操作"
	}
	if dbJudgeTaskInfo.State == cst.JudgeTaskStatusMap[cst.EXTRACTING] {
		return false, "当前任务状态为抽取中，不能进行抽取"
	}

	totalNum := extractJudgeTask.OnCampusNeedNum + extractJudgeTask.OffCampusNeedNum
	// 抽取的校内评委数和校外评委数，合计不能为0
	if totalNum <= 0 {
		return false, "校内评委数和校外评委数，合计不能为0"
	}

	// 获取专家库中所有的专家
	query := new(bizDto.EaJudgeExpertsQuery)
	query.AllowDraw = cst.YES
	query.AuditStatus = cst.EXPERT_AUDIT_PASSED
	query.Type = strconv.Itoa(dbJudgeTaskInfo.JudgeTaskType)
	expertsAll := eaJudgeExpertsService.ListAllEaJudgeExperts(query)
	expertsAll = excludeBannedJudgeExperts(expertsAll)
	if len(expertsAll) <= 0 {
		return false, "没有可抽取的专家"
	}

	// 打乱专家顺序
	rand.Shuffle(len(expertsAll), func(i, j int) {
		expertsAll[i], expertsAll[j] = expertsAll[j], expertsAll[i]
	})

	// 获取当前任务下所有已经应答的专家
	q := new(bizDto.EaJudgeDrawDetailQuery)
	q.JudgeTaskId = extractJudgeTask.JudgeTaskId
	confirmStateList := make([]any, 0)
	confirmStateList = append(confirmStateList, cst.ExpertConfirmStatusMap[cst.ATTEND], cst.ExpertConfirmStatusMap[cst.UNATTEND])
	q.AddIn("confirm_state", confirmStateList...)
	confirmedExperts := eaJudgeDrawDetailService.ListAllEaJudgeDrawDetail(q, false)
	confirmedExpertsMap := make(map[int64]bool)
	for _, v := range confirmedExperts {
		confirmedExpertsMap[v.ExpertId] = true
	}

	// 抽取的评委总数不得超过可抽取的专家数量
	//if len(expertsAll)-len(confirmedExperts) < totalNum {
	//	return false, "抽取的评委总数不得超过可抽取的专家数量"
	//}

	// 将可抽取的专家按照校内校外分组
	success, errInfo, expertsOnCampus, expertsOffCampus := _validateDrawRules(extractJudgeTask, expertsAll, confirmedExpertsMap)
	if !success {
		return false, errInfo
	}

	// 抽取过程中加锁
	key := fmt.Sprintf("extract_judge_task_lock:%d", extractJudgeTask.JudgeTaskId)
	if redis.Exists(key) {
		return false, "当前任务正在抽取中，请稍后再试"
	}
	redis.SetString(key, "1", 5*time.Minute)
	dbJudgeTaskInfo.DrawTimes += 1
	dbJudgeTaskInfo.State = cst.JudgeTaskStatusMap[cst.EXTRACTING]
	dbJudgeTaskInfo.OnCampusTotalNum = dbJudgeTaskInfo.OnCampusAttendNum
	dbJudgeTaskInfo.OffCampusTotalNum = dbJudgeTaskInfo.OffCampusAttendNum
	updateJudgeTaskInfo, recordList, detailList := _buildData(extractJudgeTask, optId, expertsOnCampus, expertsOffCampus, cst.NO, dbJudgeTaskInfo)

	// 操作数据库
	tx.ExecGTx(func(gtx tx.GTx) {

		// 更新任务表
		uWhere := &bizDo.EaJudgeTaskInfo{JudgeTaskId: dbJudgeTaskInfo.JudgeTaskId}
		updateJudgeTaskInfo.JudgeTaskId = 0
		success = eaJudgeTaskInfoDao.Update(gtx, updateJudgeTaskInfo, uWhere)

		// 插入抽取记录表
		success = eaJudgeDrawRecordsDao.BatchInsert(gtx, recordList)

		// 插入抽取明细表
		success = eaJudgeDrawDetailDao.BatchInsert(gtx, detailList)

	})
	redis.Del(key)
	return
}

func (judgeTaskInfo *EaJudgeTaskInfoService) ExtractJudgeTaskBySubject(extractJudgeTaskList []*bizDto.ExtractJudgeTaskDTO, optId int64) (success bool, errInfo string) {
	// 校验数量
	totalNum := 0
	for index, extractJudgeTask := range extractJudgeTaskList {
		totalNum = totalNum + extractJudgeTask.OnCampusNeedNum + extractJudgeTask.OffCampusNeedNum
		// 抽取的校内评委数和校外评委数，合计不能为0
		if index == len(extractJudgeTaskList)-1 && totalNum <= 0 {
			return false, "校内评委数和校外评委数，合计不能为0"
		}
	}
	judgeTaskId := extractJudgeTaskList[0].JudgeTaskId
	dbJudgeTaskInfo := eaJudgeTaskInfoDao.GetById(judgeTaskId)
	if dbJudgeTaskInfo == nil {
		return false, "任务不存在，请检查！"
	}
	if dbJudgeTaskInfo.DrawSubject == cst.NO {
		return false, "当前任务设置为不按擅长科目抽取，不能进行当前抽取操作"
	}
	if dbJudgeTaskInfo.State == cst.JudgeTaskStatusMap[cst.EXTRACTING] {
		return false, "当前任务状态为抽取中，不能进行抽取"
	}

	// 获取专家库中所有的专家
	query := new(bizDto.EaJudgeExpertsQuery)
	query.AllowDraw = cst.YES
	query.AuditStatus = cst.EXPERT_AUDIT_PASSED
	query.Type = strconv.Itoa(dbJudgeTaskInfo.JudgeTaskType)
	expertsAll := eaJudgeExpertsService.ListAllEaJudgeExperts(query)
	expertsAll = excludeBannedJudgeExperts(expertsAll)
	if len(expertsAll) <= 0 {
		return false, "没有可抽取的专家"
	}

	// 打乱专家顺序
	rand.Shuffle(len(expertsAll), func(i, j int) {
		expertsAll[i], expertsAll[j] = expertsAll[j], expertsAll[i]
	})

	// 获取当前任务下所有已经应答的专家
	q := new(bizDto.EaJudgeDrawDetailQuery)
	q.JudgeTaskId = judgeTaskId
	confirmStateList := make([]any, 0)
	confirmStateList = append(confirmStateList, cst.ExpertConfirmStatusMap[cst.ATTEND], cst.ExpertConfirmStatusMap[cst.UNATTEND])
	q.AddIn("confirm_state", confirmStateList...)
	confirmedExperts := eaJudgeDrawDetailService.ListAllEaJudgeDrawDetail(q, false)
	confirmedExpertsMap := make(map[int64]bool)
	for _, v := range confirmedExperts {
		confirmedExpertsMap[v.ExpertId] = true
	}

	// 抽取过程中加锁
	key := fmt.Sprintf("extract_judge_task_lock:%d", judgeTaskId)
	if redis.Exists(key) {
		return false, "当前任务正在抽取中，请稍后再试"
	}
	redis.SetString(key, "1", 5*time.Minute)
	drawRecordList := make([]*bizDo.EaJudgeDrawRecords, 0)
	drawDetailList := make([]*bizDo.EaJudgeDrawDetail, 0)

	// 抽取次数+1
	dbJudgeTaskInfo.DrawTimes += 1
	dbJudgeTaskInfo.State = cst.JudgeTaskStatusMap[cst.EXTRACTING]
	dbJudgeTaskInfo.OnCampusTotalNum = dbJudgeTaskInfo.OnCampusAttendNum
	dbJudgeTaskInfo.OffCampusTotalNum = dbJudgeTaskInfo.OffCampusAttendNum

	for _, extractJudgeTask := range extractJudgeTaskList {
		// 如果科目没有设置数量，就不用处理了
		if extractJudgeTask.OnCampusNeedNum+extractJudgeTask.OffCampusNeedNum == 0 {
			continue
		}
		success, errInfo, expertsOnCampus, expertsOffCampus := _validateDrawRules(extractJudgeTask, expertsAll, confirmedExpertsMap)
		if !success {
			redis.Del(key)
			return false, errInfo
		}
		updateJudgeTaskInfo, recordList, detailList := _buildData(extractJudgeTask, optId, expertsOnCampus, expertsOffCampus, cst.YES, dbJudgeTaskInfo)
		for _, detail := range detailList {
			confirmedExpertsMap[detail.ExpertId] = true
		}
		drawRecordList = append(drawRecordList, recordList...)
		drawDetailList = append(drawDetailList, detailList...)
		dbJudgeTaskInfo = updateJudgeTaskInfo
	}

	// 操作数据库
	tx.ExecGTx(func(gtx tx.GTx) {

		// 更新任务表
		uWhere := &bizDo.EaJudgeTaskInfo{JudgeTaskId: dbJudgeTaskInfo.JudgeTaskId}
		dbJudgeTaskInfo.JudgeTaskId = 0
		success = eaJudgeTaskInfoDao.Update(gtx, dbJudgeTaskInfo, uWhere)

		// 插入抽取记录表
		success = eaJudgeDrawRecordsDao.BatchInsert(gtx, drawRecordList)

		// 插入抽取明细表
		success = eaJudgeDrawDetailDao.BatchInsert(gtx, drawDetailList)

	})
	redis.Del(key)
	return
}

/*
 * desc：按规则校验抽取条件
 * author： 蓝鲸
 * date： 2025-06-12
 */
func excludeBannedJudgeExperts(experts []*bizDo.EaJudgeExperts) []*bizDo.EaJudgeExperts {
	if len(experts) == 0 {
		return experts
	}
	var bans []bizDo.EaJudgeExpertBan
	ds.GetGDB().Where("deleted = ?", common.Normal).Find(&bans)
	bannedIds := make(map[int64]bool, len(bans))
	for _, ban := range bans {
		bannedIds[ban.ExpertId] = true
	}
	result := make([]*bizDo.EaJudgeExperts, 0, len(experts))
	for _, expert := range experts {
		if !bannedIds[expert.ExpertId] {
			result = append(result, expert)
		}
	}
	return result
}

func _validateDrawRules(extractJudgeTask *bizDto.ExtractJudgeTaskDTO, expertsAll []*bizDo.EaJudgeExperts,
	confirmedExpertsMap map[int64]bool) (
	success bool, errInfo string, expertsOnCampus []*bizDo.EaJudgeExperts, expertsOffCampus []*bizDo.EaJudgeExperts) {
	success = true
	msgStart := ""
	name := extractJudgeTask.SubjectName
	if len(name) > 0 {
		msgStart = "科目：" + name + "，"
	}

	// 根据科目筛选专家库
	expertsCanBeDrawed := make([]*bizDo.EaJudgeExperts, 0)
	if len(extractJudgeTask.SubjectName) > 0 {
		for _, expert := range expertsAll {
			goodSubjects := expert.GoodSubjects
			if len(goodSubjects) < 1 {
				continue
			}
			for _, subject := range strings.Split(goodSubjects, ",") {
				if subject == extractJudgeTask.SubjectName {
					expertsCanBeDrawed = append(expertsCanBeDrawed, expert)
					break
				}
			}
		}
	} else {
		expertsCanBeDrawed = append(expertsCanBeDrawed, expertsAll...)
	}

	// 将可抽取的专家按照校内校外分组
	for _, v := range expertsCanBeDrawed {
		if confirmedExpertsMap[v.ExpertId] {
			continue
		}
		if v.DeptId == cst.BIZ_DEPTID_FOR_ADMISSIONS {
			expertsOffCampus = append(expertsOffCampus, v)
		} else {
			expertsOnCampus = append(expertsOnCampus, v)
		}
	}

	// 分别判断校内和校外的专家数量是否满足抽取数
	if len(expertsOnCampus) < extractJudgeTask.OnCampusNeedNum {
		return false, msgStart + "校内可抽取的评委数不足", nil, nil
	}
	if len(expertsOffCampus) < extractJudgeTask.OffCampusNeedNum {
		return false, msgStart + "校外可抽取的评委数不足", nil, nil
	}
	return
}

/*
 * desc：抽取数据构建
 * author： 蓝鲸
 * date： 2025-06-12
 */
func _buildData(extractJudgeTask *bizDto.ExtractJudgeTaskDTO, optId int64, expertsOnCampus []*bizDo.EaJudgeExperts,
	expertsOffCampus []*bizDo.EaJudgeExperts, drawBySubject int, dbJudgeTaskInfo *bizDo.EaJudgeTaskInfo) (
	updateJudgeTaskInfo *bizDo.EaJudgeTaskInfo, drawRecordList []*bizDo.EaJudgeDrawRecords, drawDetailList []*bizDo.EaJudgeDrawDetail) {
	dbJudgeTaskInfo.OnCampusTotalNum = dbJudgeTaskInfo.OnCampusTotalNum + extractJudgeTask.OnCampusNeedNum
	dbJudgeTaskInfo.OffCampusTotalNum = dbJudgeTaskInfo.OffCampusTotalNum + extractJudgeTask.OffCampusNeedNum
	updateJudgeTaskInfo = dbJudgeTaskInfo

	// 构建抽取任务记录列表
	drawRecord := &bizDo.EaJudgeDrawRecords{
		JudgeTaskId:          extractJudgeTask.JudgeTaskId,
		SubjectName:          extractJudgeTask.SubjectName,
		DrawTimes:            dbJudgeTaskInfo.DrawTimes,
		OnCampusTotalNum:     extractJudgeTask.OnCampusNeedNum,
		OnCampusAttendNum:    0,
		OnCampusUnattendNum:  0,
		OffCampusTotalNum:    extractJudgeTask.OffCampusNeedNum,
		OffCampusAttendNum:   0,
		OffCampusUnattendNum: 0,
	}
	drawRecord.SetCreateBy(optId)
	drawRecord.SetUpdateBy(optId)
	drawRecord.Deleted = common.Normal
	drawRecordList = append(drawRecordList, drawRecord)

	// 构建抽取明细列表
	for i := 0; i < extractJudgeTask.OnCampusNeedNum; i++ {
		experts := expertsOnCampus[i]
		temp := &bizDo.EaJudgeDrawDetail{
			JudgeTaskId:    extractJudgeTask.JudgeTaskId,
			SubjectName:    extractJudgeTask.SubjectName,
			TaskTimeRemark: time.Time(dbJudgeTaskInfo.TaskStartTime).Format("2006-01-02") + "至" + time.Time(dbJudgeTaskInfo.TaskEndTime).Format("2006-01-02"),
			DrawTimes:      dbJudgeTaskInfo.DrawTimes,
			ExpertId:       experts.ExpertId,
			ExpertName:     experts.Name,
			DeptId:         experts.DeptId,
			ConfirmState:   cst.ExpertConfirmStatusMap[cst.NOT_CONFIRM],
			DrawBySubject:  drawBySubject,
		}
		temp.SetCreateBy(optId)
		temp.SetUpdateBy(optId)
		temp.Deleted = common.Normal
		drawDetailList = append(drawDetailList, temp)
	}
	for i := 0; i < extractJudgeTask.OffCampusNeedNum; i++ {
		experts := expertsOffCampus[i]
		temp := &bizDo.EaJudgeDrawDetail{
			JudgeTaskId:    extractJudgeTask.JudgeTaskId,
			SubjectName:    extractJudgeTask.SubjectName,
			TaskTimeRemark: time.Time(dbJudgeTaskInfo.TaskStartTime).Format("2006-01-02") + "至" + time.Time(dbJudgeTaskInfo.TaskEndTime).Format("2006-01-02"),
			DrawTimes:      dbJudgeTaskInfo.DrawTimes,
			ExpertId:       experts.ExpertId,
			ExpertName:     experts.Name,
			DeptId:         experts.DeptId,
			ConfirmState:   cst.ExpertConfirmStatusMap[cst.NOT_CONFIRM],
			DrawBySubject:  drawBySubject,
		}
		temp.SetCreateBy(optId)
		temp.SetUpdateBy(optId)
		temp.Deleted = common.Normal
		drawDetailList = append(drawDetailList, temp)
	}
	return
}

func (judgeTaskInfo *EaJudgeTaskInfoService) CancelExtract(id int64, optId int64) (success bool, errInfo string) {
	// 状态检测
	dbJudgeTaskInfo := eaJudgeTaskInfoDao.GetById(id)
	if dbJudgeTaskInfo == nil {
		return false, "任务不存在，请检查！"
	}
	if dbJudgeTaskInfo.State != cst.JudgeTaskStatusMap[cst.EXTRACTING] {
		return false, "当前任务状态不允许取消抽取"
	}

	tx.ExecGTx(func(gtx tx.GTx) {

		// 更新任务表
		uWhere := &bizDo.EaJudgeTaskInfo{JudgeTaskId: id}
		judgeTaskInfo := &bizDo.EaJudgeTaskInfo{
			State: cst.JudgeTaskStatusMap[cst.FINISHED],
		}
		judgeTaskInfo.SetUpdateBy(optId)
		success = eaJudgeTaskInfoDao.Update(gtx, judgeTaskInfo, uWhere)

		// 删除没有确认的评委抽取详情
		uWhere4Detail := &bizDo.EaJudgeDrawDetail{
			ConfirmState: cst.ExpertConfirmStatusMap[cst.NOT_CONFIRM],
			JudgeTaskId:  id,
		}
		uWhere4Detail.Deleted = common.Normal
		uDate := new(bizDo.EaJudgeDrawDetail)
		uDate.Deleted = common.Deleted
		uDate.SetUpdateBy(optId)
		eaJudgeDrawDetailDao.Update(gtx, uDate, uWhere4Detail)
	})

	return
}

// UnpublishCollegeReportTask 取消发布（将已发布状态回退为未发布）
func (judgeTaskInfo *EaJudgeTaskInfoService) UnpublishCollegeReportTask(id int64, optId int64) (success bool, errInfo string) {
	dbTask := eaJudgeTaskInfoDao.GetById(id)
	if dbTask == nil {
		return false, "任务不存在"
	}
	if dbTask.TaskMode != 2 {
		return false, "只有学院上报模式可以取消发布"
	}
	if dbTask.PublishState != 2 {
		return false, "当前任务未发布，无法取消"
	}
	uData := &bizDo.EaJudgeTaskInfo{PublishState: 1}
	uData.SetUpdateBy(optId)
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaJudgeTaskInfoDao.Update(gtx, uData, &bizDo.EaJudgeTaskInfo{JudgeTaskId: id})
	})
	return
}

// CancelImportReports 取消导入（软删已导入的上报要求数据）。
// deleted 使用 UNIX_TIMESTAMP() 时间戳，避免 uk_task_direction_dept 唯一键冲突。
func (judgeTaskInfo *EaJudgeTaskInfoService) CancelImportReports(id int64, optId int64) (success bool, errInfo string) {
	dbTask := eaJudgeTaskInfoDao.GetById(id)
	if dbTask == nil {
		return false, "任务不存在"
	}
	if dbTask.TaskMode != 2 {
		return false, "只有学院上报模式可以取消导入"
	}
	err := ds.GetGDB().Transaction(func(tx *gorm.DB) error {
		return tx.Exec(
			"UPDATE ea_judge_task_report_requirement SET deleted=UNIX_TIMESTAMP(),update_by=?,update_time=NOW() WHERE judge_task_id=? AND deleted=1",
			optId, id,
		).Error
	})
	if err != nil {
		return false, "取消导入失败：" + err.Error()
	}
	return true, ""
}

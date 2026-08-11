package bizSerImpl

import (
	"fmt"
	"time"
	"yx-exam-affair-gm/app/biz/bizDao"
	"yx-exam-affair-gm/app/biz/bizDao/bizDaoImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	"yx-exam-affair-gm/app/cst"
	"yx-exam-affair-gm/app/hint"

	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/ds"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/utils/file/excel"
	"go.uber.org/zap"
)

/*
 * desc: 监考老师征集服务层实现
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchTaskInfoService struct{}

var eaTchTaskInfoService *EaTchTaskInfoService
var eaTchTaskInfoDao bizDao.IEaTchTaskInfoDao

func init() {
	eaTchTaskInfoService = &EaTchTaskInfoService{}
	eaTchTaskInfoDao = bizDaoImpl.NewEaTchTaskInfoDao()
	eaTchTaskReportDao = bizDaoImpl.NewEaTchTaskReportDao()
}

func NewEaTchTaskInfoService() *EaTchTaskInfoService {
	return eaTchTaskInfoService
}

func (tchTaskInfo *EaTchTaskInfoService) GetEaTchTaskInfoById(id int64) *bizDo.EaTchTaskInfo {
	return eaTchTaskInfoDao.GetById(id)
}

func (tchTaskInfo *EaTchTaskInfoService) CountEaTchTaskInfo(q *bizDto.EaTchTaskInfoQuery) int {
	return eaTchTaskInfoDao.Count(q)
}

func (tchTaskInfo *EaTchTaskInfoService) ListEaTchTaskInfo(q *bizDto.EaTchTaskInfoQuery, p *page.Page) []*bizDo.EaTchTaskInfo {
	q.AddLikeAll("tch_task_name", q.TchTaskName)
	list := eaTchTaskInfoDao.List(q, p)
	if len(list) == 0 {
		return list
	}
	ids := make([]int64, 0, len(list))
	for _, item := range list {
		ids = append(ids, item.TchTaskId)
	}
	var auditPassCounts []struct {
		TchTaskId int64 `gorm:"column:tch_task_id"`
		Count     int   `gorm:"column:count"`
	}
	sql := `SELECT tch_task_id, COUNT(*) AS count FROM ea_tch_task_report
	        WHERE tch_task_id IN (?) AND state IN (?, ?) AND deleted = ?
	        GROUP BY tch_task_id`
	ds.GetGDB().Raw(sql, ids, cst.TCH_REPORT_PENDING_ASSIGN, cst.TCH_REPORT_ASSIGNED, common.Normal).Scan(&auditPassCounts)
	countMap := make(map[int64]int, len(auditPassCounts))
	for _, c := range auditPassCounts {
		countMap[c.TchTaskId] = c.Count
	}
	for _, item := range list {
		item.AuditPassNum = countMap[item.TchTaskId]
	}
	return list
}

func (tchTaskInfo *EaTchTaskInfoService) ListAllEaTchTaskInfo(q *bizDto.EaTchTaskInfoQuery) []*bizDo.EaTchTaskInfo {
	listAll := make([]*bizDo.EaTchTaskInfo, 0)
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE
	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		list := eaTchTaskInfoDao.List(q, pageInfo)
		if len(list) == 0 {
			break
		}
		listAll = append(listAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}
	return listAll
}

/*
 * desc: 数组去重
 * author: 蓝鲸
 * date: 2024/09/05
 */
func unique(slice []int64) []int64 {
	keys := make(map[int64]bool)
	list := make([]int64, 0)

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

func (tchTaskInfo *EaTchTaskInfoService) ListAllEaTchTaskInfoByDept(deptId int64) []*bizDo.EaTchTaskInfo {

	// 查询所有需要上报的任务列表
	taskDetailList := eaTchTaskDetailService.ListAllEaTchTaskDetail(&bizDto.EaTchTaskDetailQuery{DeptId: deptId})
	if len(taskDetailList) < 1 {
		return nil
	}

	// 遍历得到所有任务id
	taskIdMap := make(map[int64][]*bizDo.EaTchTaskDetail)
	taskIdList := make([]interface{}, 0)
	for _, v := range taskDetailList {
		taskIdList = append(taskIdList, v.TchTaskId)
		if taskIdMap[v.TchTaskId] == nil {
			taskIdMap[v.TchTaskId] = make([]*bizDo.EaTchTaskDetail, 0)
		}
		taskIdMap[v.TchTaskId] = append(taskIdMap[v.TchTaskId], v)
	}

	// 查询所有已启用的任务列表
	taskQuery := &bizDto.EaTchTaskInfoQuery{
		State: cst.PUBLISH,
	}
	taskQuery.AddIn("tch_task_id", taskIdList...)

	taskInfoList := eaTchTaskInfoService.ListAllEaTchTaskInfo(taskQuery)
	for _, v := range taskInfoList {
		allNeedNum := 0
		allReportNum := 0
		if len(taskIdMap[v.TchTaskId]) < 1 {
			continue
		}
		for _, v2 := range taskIdMap[v.TchTaskId] {
			allNeedNum += v2.NeedNum
			allReportNum += v2.CollectNum
		}
		v.CollectNum = allReportNum
		v.NeedNum = allNeedNum
	}
	return taskInfoList
}

func (tchTaskInfo *EaTchTaskInfoService) InsertEaTchTaskInfo(tchTaskInfoDO *bizDo.EaTchTaskInfo) (err errs.ERROR) {
	// 遍历每个任务详情列表
	taskDetailList := tchTaskInfoDO.TchTaskDetail
	if len(taskDetailList) == 0 {
		return hint.BIZ_TCH_TASK_DETAIL_EMPTY_ERR
	}
	allNeedNum := 0
	taskDetailMap := make(map[string]*bizDo.EaTchTaskDetail)
	for _, taskDetail := range taskDetailList {
		if taskDetail.TaskTime.IsZero() || taskDetail.DeptId == 0 || taskDetail.NeedNum == 0 {
			timeStr := taskDetail.TaskTime.FormatString("01月02日")
			return hint.BIZ_TCH_TASK_DETAIL_INFO_ERR.Format(timeStr, taskDetail.DeptName, taskDetail.NeedNum)
		}
		allNeedNum += taskDetail.NeedNum
		// 转map存储 日期_学院id
		key := fmt.Sprintf("%v_%v", taskDetail.TaskTime.FormatString(time.DateOnly), taskDetail.DeptId)
		// 顺便检验传来的列表里是否有重复记录
		if taskDetailMap[key] != nil {
			return hint.BIZ_TCH_TASK_DETAIL_DEPT_DUPLICATE_ERR.Format(taskDetail.TaskTime.FormatString("01月02日"), taskDetail.DeptName)
		}
		taskDetailMap[key] = taskDetail
	}
	tchTaskInfoDO.NeedNum = allNeedNum
	tchTaskInfoDO.TaskTimeRemark = fmt.Sprintf("%v-%v", tchTaskInfoDO.TaskStartTime.FormatString(time.DateOnly), tchTaskInfoDO.TaskEndTime.FormatString(time.DateOnly))
	tx.ExecGTx(func(gtx tx.GTx) {
		// 先插入征集任务表
		tchTaskInfoDO.State = cst.UNPUBLISH
		tchTaskInfoDO.CollectNum = 0
		success := eaTchTaskInfoDao.Insert(gtx, tchTaskInfoDO)
		if success && tchTaskInfoDO.TchTaskId > 0 {
			// 再批量插入征集详情记录
			for _, taskDetail := range taskDetailList {
				taskDetail.TchTaskId = tchTaskInfoDO.TchTaskId
				taskDetail.CollectNum = 0
				taskDetail.Deleted = common.Normal
				taskDetail.SetCreateBy(taskDetail.CreateBy)
				taskDetail.SetUpdateBy(taskDetail.CreateBy)
			}
			eaTchTaskDetailDao.BatchInsert(gtx, taskDetailList)
		}
	})
	return
}

func (tchTaskInfo *EaTchTaskInfoService) UpdateTchTaskState(tchTaskId int64, state string, oper int64) (success bool) {
	uWhere := &bizDo.EaTchTaskInfo{TchTaskId: tchTaskId}
	uSet := &bizDo.EaTchTaskInfo{State: state}
	uSet.SetUpdateBy(oper)
	// 数据库执行事务
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaTchTaskInfoDao.Update(gtx, uSet, uWhere)
	})
	return
}

func (tchTaskInfo *EaTchTaskInfoService) UpdateEaTchTaskInfoById(tchTaskInfoDO *bizDo.EaTchTaskInfo) (err errs.ERROR) {
	// 查询任务信息
	taskInfoDB := eaTchTaskInfoDao.GetById(tchTaskInfoDO.TchTaskId)
	if taskInfoDB == nil {
		return errs.PARAM_INVALID_ERR
	}
	if taskInfoDB.State == cst.PUBLISH {
		return hint.BIZ_TCH_TASK_STATE_PUBLISHED_ERR
	}
	operId := tchTaskInfoDO.UpdateBy

	// 遍历每个任务详情列表
	taskDetailList := tchTaskInfoDO.TchTaskDetail
	if len(taskDetailList) == 0 {
		return hint.BIZ_TCH_TASK_DETAIL_EMPTY_ERR
	}
	// 任务-总需人数
	allNeedNum := 0
	taskDetailMap := make(map[string]*bizDo.EaTchTaskDetail)
	// 新增列表 没有detailId的记录
	taskDetailListAdd := make([]*bizDo.EaTchTaskDetail, 0)
	for _, taskDetail := range taskDetailList {
		if taskDetail.TaskTime.IsZero() || taskDetail.DeptId == 0 {
			return errs.PARAM_INVALID_ERR
		}
		if taskDetail.DeptName == "" {
			dept := eaDeptInfoDao.GetById(taskDetail.DeptId)
			taskDetail.DeptName = dept.DeptName
		}
		timeStr := taskDetail.TaskTime.FormatString("01月02日")
		if taskDetail.NeedNum == 0 {
			return hint.BIZ_TCH_TASK_DETAIL_INFO_ERR.Format(timeStr, taskDetail.DeptName, taskDetail.NeedNum)
		}
		allNeedNum += taskDetail.NeedNum
		// 征集人数不改变，新增为0，修改不动
		taskDetail.CollectNum = 0
		// 没有主键，为新增记录
		if taskDetail.TaskDetailId == 0 {
			taskDetail.TchTaskId = tchTaskInfoDO.TchTaskId
			taskDetail.Deleted = common.Normal
			taskDetail.SetCreateBy(operId)
			taskDetail.SetUpdateBy(operId)
			taskDetailListAdd = append(taskDetailListAdd, taskDetail)
		}
		// 转map存储 日期_学院id
		key := fmt.Sprintf("%v_%v", taskDetail.TaskTime.FormatString(time.DateOnly), taskDetail.DeptId)
		// 顺便检验传来的列表里是否有重复记录
		if taskDetailMap[key] != nil {
			return hint.BIZ_TCH_TASK_DETAIL_DEPT_DUPLICATE_ERR.Format(timeStr, taskDetail.DeptName)
		}
		taskDetailMap[key] = taskDetail
	}

	// 数据库查询该任务现有的所有明细记录, 找出新增、修改、删除列表
	q := &bizDto.EaTchTaskDetailQuery{TchTaskId: tchTaskInfoDO.TchTaskId}
	taskDetailDBList := eaTchTaskDetailDao.List(q, &page.Page{CurPage: 1, PageSize: page.MAX_PAGE_SIZE})
	// 修改记录列表，数据库里已存在
	taskDetailListUpdate := make([]*bizDo.EaTchTaskDetail, 0)
	// 删除列表，数据库有，传过来的新列表里没有
	taskDetailListDelete := make([]*bizDo.EaTchTaskDetail, 0)
	// 任务-已征人数, 删除列表处理后需要释放对应人数
	allCollectNum := 0
	for _, taskDetailDB := range taskDetailDBList {
		key := fmt.Sprintf("%v_%v", taskDetailDB.TaskTime.FormatString(time.DateOnly), taskDetailDB.DeptId)
		taskDetail := taskDetailMap[key]
		// 修改
		if taskDetail != nil {
			// 数据异常，修改必须有主键id
			if taskDetail.TaskDetailId == 0 {
				return hint.BIZ_TCH_TASK_DETAIL_INFO_ERR.Format(taskDetail.TaskTime.FormatString("01月02日"),
					taskDetail.DeptName, taskDetail.NeedNum)
			}
			taskDetail.SetUpdateBy(operId)
			taskDetailListUpdate = append(taskDetailListUpdate, taskDetail)
			// 总共已征集人数，就只有这里被修改的记录
			allCollectNum += taskDetailDB.CollectNum
		} else {
			// 删除
			taskDetailListDelete = append(taskDetailListDelete, taskDetailDB)
		}
	}

	// 数据库处理
	tx.ExecGTx(func(gtx tx.GTx) bool {
		// 更新征集任务总需人数、已征集人数
		uTaskWhere := &bizDo.EaTchTaskInfo{TchTaskId: tchTaskInfoDO.TchTaskId}
		uTaskData := &bizDo.EaTchTaskInfo{NeedNum: allNeedNum, CollectNum: allCollectNum}
		uTaskData.SetUpdateBy(operId)
		eaTchTaskInfoDao.Update(gtx, uTaskData, uTaskWhere)
		// 遍历修改
		for _, upItem := range taskDetailListUpdate {
			if upItem.TaskDetailId == 0 {
				zap.L().Error("征集明细修改数据异常!")
				return false
			}
			uDetailWhere := &bizDo.EaTchTaskDetail{TaskDetailId: upItem.TaskDetailId}
			upItem.TchTaskId = 0
			upItem.TaskDetailId = 0
			eaTchTaskDetailDao.Update(gtx, upItem, uDetailWhere)
		}
		// 遍历删除
		for _, delItem := range taskDetailListDelete {
			eaTchTaskDetailDao.DeleteById(gtx, delItem.TaskDetailId, operId)
			// 清除上报人数
			eaTchTaskReportDao.DeleteByTaskDetailId(gtx, delItem.TaskDetailId, operId)
		}
		// 批量新增
		eaTchTaskDetailDao.BatchInsert(gtx, taskDetailListAdd)
		return true
	})
	return
}

func (tchTaskInfo *EaTchTaskInfoService) DeleteEaTchTaskInfoById(id int64, deleteBy int64) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaTchTaskInfoDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

// 导出服务 -- start
type EaTchTaskInfoExportStruct struct {
	Query *bizDto.EaTchTaskInfoQuery
}

func (h *EaTchTaskInfoExportStruct) Title() string {
	return "监考老师征集导出"
}

func (h *EaTchTaskInfoExportStruct) HeaderColumn() []string {
	return []string{"教师征集任务ID,tch_task_id", "考试任务名称,tch_task_name", "征集日期说明,task_time_remark", "征集开始时间,task_start_time", "征集结束时间,task_end_time", "总征集数量,need_num", "已征集数量,collect_num", "发布状态: 1-未发布 2-已发布,state", "删除标记：1-正常的；2-已删除的；,deleted", "备注,remark", "创建者,create_by", "创建时间,create_time", "更新者,update_by", "更新时间,update_time"}
}

func (h *EaTchTaskInfoExportStruct) Rows(p *page.Page) []map[string]interface{} {
	return eaTchTaskInfoDao.ListMapWithCondition(h.Query, p, nil, nil)
}

func (h *EaTchTaskInfoExportStruct) Finish(url string) {
}

// 导出服务 -- end
func (tchTaskInfo *EaTchTaskInfoService) ExportEaTchTaskInfo(q *bizDto.EaTchTaskInfoQuery) ([]byte, error) {
	return excel.NewExcel().ExportSync(&EaTchTaskInfoExportStruct{q})
}

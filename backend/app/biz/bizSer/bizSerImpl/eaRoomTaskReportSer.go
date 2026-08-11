package bizSerImpl

import (
	"strconv"
	"yx-exam-affair-gm/app/biz/bizDao"
	"yx-exam-affair-gm/app/biz/bizDao/bizDaoImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	"yx-exam-affair-gm/app/cst"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/cst/common"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/query"
	"github.com/sealsee/web-base/public/ds/tx"
	"github.com/sealsee/web-base/public/utils/file/excel"
)

/*
 * desc: 任务教室征集关联服务层实现
 * author: 蓝鲸
 * date: 2024/08/14
 */
type EaRoomTaskReportService struct{}

var eaRoomTaskReportService *EaRoomTaskReportService
var eaRoomTaskReportDao bizDao.IEaRoomTaskReportDao

func init() {
	eaRoomTaskReportService = &EaRoomTaskReportService{}
	eaRoomTaskReportDao = bizDaoImpl.NewEaRoomTaskReportDao()
}

func NewEaRoomTaskReportService() *EaRoomTaskReportService {
	return eaRoomTaskReportService
}

func (taskRoomInfo *EaRoomTaskReportService) GetEaTaskRoomInfoById(id int64) *bizDo.EaRoomTaskReport {
	return eaRoomTaskReportDao.GetById(id)
}

// func (taskRoomInfo *EaRoomTaskReportService) ListEaTaskRoomInfo(q *bizDto.EaRoomTaskReportQuery) []*bizDo.EaRoomTaskReport {
// 	taskRoomInfoList := taskRoomInfo.ListAllEaTaskRoomInfo(q)
// 	if len(taskRoomInfoList) <= 0 {
// 		return nil
// 	}

// 	roomIdList := make([]int64, 0)
// 	for _, taskRoomInfo := range taskRoomInfoList {
// 		roomIdList = append(roomIdList, taskRoomInfo.RoomId)
// 	}
// 	roomInfoMap := make(map[int64]*bizDo.EaRoomInfo)
// 	roomBaseInfoList := eaRoomInfoDao.ListWithCondition(&bizDto.EaRoomInfoQuery{}, &page.Page{CurPage: 1, PageSize: page.MAX_PAGE_SIZE}, "room_id in ?", roomIdList)
// 	if len(roomBaseInfoList) <= 0 {
// 		return taskRoomInfoList
// 	}
// 	for _, roomInfo := range roomBaseInfoList {
// 		roomInfoMap[roomInfo.RoomId] = roomInfo
// 	}
// 	for _, taskRoomInfo := range taskRoomInfoList {
// 		if roomInfoMap[taskRoomInfo.RoomId] == nil {
// 			continue
// 		}
// 		taskRoomInfo.CampusName = roomInfoMap[taskRoomInfo.RoomId].CampusName
// 		taskRoomInfo.GroupCapacity = roomInfoMap[taskRoomInfo.RoomId].GroupCapacity
// 		taskRoomInfo.LevelCode = roomInfoMap[taskRoomInfo.RoomId].LevelCode
// 		taskRoomInfo.DeptId = roomInfoMap[taskRoomInfo.RoomId].DeptId
// 		taskRoomInfo.Capacity = roomInfoMap[taskRoomInfo.RoomId].Capacity
// 		taskRoomInfo.GroupNum = roomInfoMap[taskRoomInfo.RoomId].GroupNum
// 	}
// 	return taskRoomInfoList
// }

func (taskRoomInfo *EaRoomTaskReportService) ListAllEaTaskRoomInfo(q *bizDto.EaRoomTaskReportQuery) []*bizDto.EaRoomTaskReportDBDto {
	listDepartAll := make([]*bizDto.EaRoomTaskReportDBDto, 0)
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE
	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		var list []*bizDto.EaRoomTaskReportDBDto
		// list = eaRoomTaskReportDao.List(q, pageInfo)
		sql := `select a.*, b.campus_name, b.level_code
				from ea_room_task_report a, ea_room_info b 
				where a.deleted = ? and b.deleted = 1 and a.room_id = b.room_id`
		if q.State == cst.YES1 {
			sql += ` and b.state = '1'`
		}
		q.SetAlias("a")
		q.AddLikeAll("room_name", q.RoomName)
		list = query.RawSqlQueryListWithPageWhere[bizDto.EaRoomTaskReportDBDto](q, pageInfo, sql, 1)
		if len(list) == 0 {
			break
		}
		listDepartAll = append(listDepartAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}
	return listDepartAll
}

func (taskRoomInfo *EaRoomTaskReportService) ListAllRoomReport(q *bizDto.EaRoomTaskReportQuery) []*bizDo.EaRoomTaskReport {
	listDepartAll := make([]*bizDo.EaRoomTaskReport, 0)
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE
	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		list := eaRoomTaskReportDao.List(q, pageInfo)
		if len(list) == 0 {
			break
		}
		listDepartAll = append(listDepartAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}
	return listDepartAll
}

func (taskRoomInfo *EaRoomTaskReportService) ListCanUseRoomListByTaskId(q *bizDto.EaRoomTaskReportQuery, p *page.Page) (enListNum int, roomList []*bizDto.EaRoomEnlistDto) {
	// 获取任务信息
	taskInfo := eaRoomTaskInfoService.GetEaTaskInfoById(q.TaskId)
	if taskInfo == nil {
		return
	}
	// 查询当前任务的已上报数量
	sql := `select count(1) from ea_room_task_report a 
			left join ea_room_info b on b.deleted = 1 and a.room_id = b.room_id
			where a.deleted = 1 and a.task_id = ? and a.dept_id = ?`
	enListNum = query.RawSqlQuery[int](sql, q.TaskId, q.DeptId)

	// 查询可使用的考场列表
	sql2 := `select r.* from (
				select a.*, case b.state when '1' then '2' when '2' then '3' else '1' end as publish_state 
				from ea_room_info a 
				left join ea_room_task_report b 
					on b.deleted = 1 and b.room_id = a.room_id and b.task_id = ? and b.dept_id = a.dept_id
				where a.deleted = 1 and a.dept_id = ? and a.level_code = ?
				and a.room_id not in (
					select c.room_id from ea_room_task_report c where c.deleted = 1 and c.task_id <> ? and c.task_day = ?
				)
			) r`
	// 查询条件
	if q.PublishState != "" {
		sql2 += " where r.publish_state = " + q.PublishState
	}
	roomList = query.RawSqlQueryListWithPage[bizDto.EaRoomEnlistDto](p, sql2, q.TaskId, q.DeptId, taskInfo.LevelCode, q.TaskId, taskInfo.TaskDay)
	return
}

func (taskRoomInfo *EaRoomTaskReportService) InsertEaTaskRoomInfo(taskRoomInfoDO *bizDo.EaRoomTaskReport) (success bool, errInfo string) {

	// 获取当前任务
	taskInfo := eaRoomTaskInfoService.GetEaTaskInfoById(taskRoomInfoDO.TaskId)
	if taskInfo == nil {
		return false, "所选任务不存在"
	}
	if taskInfo.State != cst.PUBLISH {
		return false, "所选任务未发布"
	}

	// 获取考场详情
	classroomInfo := eaRoomInfoService.GetEaClassroomInfoById(taskRoomInfoDO.RoomId)
	if classroomInfo == nil {
		return false, "所选教室不存在"
	}
	taskRoomInfoDO.TaskDay = taskInfo.TaskDay
	taskRoomInfoDO.DeptId = classroomInfo.DeptId
	taskRoomInfoDO.RoomName = classroomInfo.RoomName
	taskRoomInfoDO.Capacity = classroomInfo.Capacity
	taskRoomInfoDO.GroupNum = classroomInfo.GroupNum
	taskRoomInfoDO.GroupCapacity = classroomInfo.GroupCapacity

	if taskInfo.LevelCode != classroomInfo.LevelCode {
		return false, "所选教室与当前任务层次不匹配"
	}

	// 检查当前考场是否已经被上报给相同日期的考试任务
	taskRoomQuery := &bizDto.EaRoomTaskReportQuery{
		RoomId: taskRoomInfoDO.RoomId,
	}
	taskRoomInfos := taskRoomInfo.ListAllEaTaskRoomInfo(taskRoomQuery)
	if len(taskRoomInfos) > 0 {

		// 获取此考场所有已经被上报的任务
		taskIdList := make([]int64, 0)
		for _, taskRoomInfo := range taskRoomInfos {
			// 如果当前教室已经被上报过该日程了，则不允许重复上报
			if taskRoomInfo.TaskId == taskRoomInfoDO.TaskId {
				return false, "此任务下已经上报过当前教室，请检查重试"
			}
			taskIdList = append(taskIdList, taskRoomInfo.TaskId)
		}
		taskMap := make(map[int64]*bizDo.EaRoomTaskInfo)
		taskInfoList := eaRoomTaskInfoDao.ListWithCondition(&bizDto.EaRoomTaskInfoQuery{}, &page.Page{CurPage: 1, PageSize: page.MAX_PAGE_SIZE}, "task_id in (?)", taskIdList)
		if len(taskInfoList) > 0 {
			for _, taskInfo := range taskInfoList {
				taskMap[taskInfo.TaskId] = taskInfo
			}
		}

		for _, taskRoomInfo := range taskRoomInfos {
			if taskMap[taskRoomInfo.TaskId] == nil || taskMap[taskRoomInfo.TaskId].TaskDay.IsZero() {
				return false, "数据有误，请刷新重试"
			}
			if taskMap[taskRoomInfo.TaskId].TaskDay == taskInfo.TaskDay {
				return false, "教室已被其他相同日期的考试任务使用，请检查数据"
			}
		}
	}

	taskRoomInfoDO.State = common.OK
	taskRoomInfoDO.Deleted = common.Normal
	return tx.ExecGTx(func(gtx tx.Db) {
		eaRoomTaskReportDao.Insert(gtx, taskRoomInfoDO)
		eaRoomTaskInfoDao.UpdateCollectNumAndConfirmNum(gtx, taskRoomInfoDO.TaskId, 1, 0)
	}), ""
}

func (taskRoomInfo *EaRoomTaskReportService) UpdateEaTaskRoomInfoById(taskRoomInfoDO *bizDo.EaRoomTaskReport) (success bool) {
	uWhere := &bizDo.EaRoomTaskReport{Id: taskRoomInfoDO.Id}
	taskRoomInfoDO.Id = 0
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaRoomTaskReportDao.Update(gtx, taskRoomInfoDO, uWhere)
	})
	return
}

func (taskRoomInfo *EaRoomTaskReportService) BatchUpdateEaTaskRoomInfo(list []*bizDo.EaRoomTaskReport, taskId int64, updateBy int64) (bool, string) {

	// 获取当前任务下所有已经确认了的考场
	dbConfirmedTaskRoomInfoList := eaRoomTaskReportService.ListAllEaTaskRoomInfo(&bizDto.EaRoomTaskReportQuery{
		TaskId: taskId,
		State:  cst.CONFIRM,
	})
	list2UpdateConfirm := make([]*bizDo.EaRoomTaskReport, 0)
	taskRoomId2UpdateMap := make(map[int64]bool)
	for _, taskRoomInfo := range list {
		if taskRoomInfo.Id <= 0 {
			continue
		}
		dbRoomReport := eaRoomTaskReportDao.GetById(taskRoomInfo.Id)
		if dbRoomReport == nil {
			return false, "不存在已上报的数据，请检查"
		}
		if dbRoomReport.State == cst.CONFIRM && dbRoomReport.TchNum != taskRoomInfo.TchNum && checkRoomReportHasArrange(taskRoomInfo.Id) {
			return false, dbRoomReport.RoomName + "已经被编排，不允许修改"
		}
		taskRoomId2UpdateMap[taskRoomInfo.Id] = true
		temp := &bizDo.EaRoomTaskReport{
			Id:     taskRoomInfo.Id,
			TchNum: taskRoomInfo.TchNum,
			State:  cst.CONFIRM,
		}
		temp.SetUpdateBy(updateBy)
		list2UpdateConfirm = append(list2UpdateConfirm, temp)
	}

	list2UpdateUnConfirm := make([]*bizDo.EaRoomTaskReport, 0)
	for _, i := range dbConfirmedTaskRoomInfoList {
		if taskRoomId2UpdateMap[i.Id] {
			continue
		}
		if checkRoomReportHasArrange(i.Id) {
			return false, i.RoomName + "已经被编排，不允许清除"
		}
		temp := &bizDo.EaRoomTaskReport{
			Id:     i.Id,
			TchNum: 0,
			State:  cst.UNCONFIRM,
		}
		temp.SetUpdateBy(updateBy)
		temp.SetToZeroCols("tch_num")
		list2UpdateUnConfirm = append(list2UpdateUnConfirm, temp)
	}

	if len(list2UpdateConfirm) <= 0 && len(list2UpdateUnConfirm) <= 0 {
		return true, ""
	}

	eaTaskInfo2UpdateWhere := &bizDo.EaRoomTaskInfo{
		TaskId: taskId,
	}
	eaTaskInfo2Update := &bizDo.EaRoomTaskInfo{
		TaskId:     taskId,
		ComfirmNum: len(list2UpdateConfirm),
	}
	eaTaskInfo2Update.SetToZeroCols("comfirm_num")

	return tx.ExecGTx(func(gtx tx.Db) {
		if len(list2UpdateConfirm) > 0 {
			for _, taskRoomInfo := range list2UpdateConfirm {
				eaRoomTaskReportDao.Update(gtx, taskRoomInfo, &bizDo.EaRoomTaskReport{Id: taskRoomInfo.Id})
			}
		}
		if len(list2UpdateUnConfirm) > 0 {
			for _, taskRoomInfo := range list2UpdateUnConfirm {
				eaRoomTaskReportDao.UpdateNew(gtx, &bizDo.EaRoomTaskReport{Id: taskRoomInfo.Id}, taskRoomInfo)
			}
		}
		eaRoomTaskInfoDao.UpdateNew(gtx, eaTaskInfo2UpdateWhere, eaTaskInfo2Update)
	}), ""
}

func checkRoomReportHasArrange(roomReportId int64) bool {
	return eaArrangeRoomDao.Count(&bizDto.EaArrangeRoomQuery{RoomReportId: roomReportId}) > 0
}

func (taskRoomInfo *EaRoomTaskReportService) UpdateRoomNameById(classroomInfoDO *bizDo.EaRoomInfo, updateBy int64) (success bool) {
	uWhere := &bizDo.EaRoomTaskReport{RoomId: classroomInfoDO.RoomId}
	taskRoomInfoDO := &bizDo.EaRoomTaskReport{RoomName: classroomInfoDO.RoomName, DeptId: classroomInfoDO.DeptId}
	taskRoomInfoDO.SetUpdateBy(updateBy)
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaRoomTaskReportDao.Update(gtx, taskRoomInfoDO, uWhere)
	})
	return
}

func (taskRoomInfo *EaRoomTaskReportService) DeleteEaTaskRoomInfo(query *bizDto.EaRoomTaskReportQuery) (bool, string) {

	// 更具任务id和教室id获取关联信息
	taskRoomQuery := &bizDto.EaRoomTaskReportQuery{
		RoomId: query.RoomId,
		TaskId: query.TaskId,
	}
	taskRoomList := eaRoomTaskReportDao.List(taskRoomQuery, &page.Page{
		CurPage:  1,
		PageSize: 1,
	})
	if len(taskRoomList) < 1 {
		return false, "所选教室未上报到此任务"
	}
	if taskRoomList[0].State == cst.CONFIRM {
		return false, "教室已被确认，无法取消"
	}

	return tx.ExecGTx(func(gtx tx.Db) {
		eaRoomTaskReportDao.DeleteById(gtx, taskRoomList[0].Id)
		eaRoomTaskInfoDao.UpdateCollectNumAndConfirmNum(gtx, query.TaskId, -1, 0)
	}), ""
}

func (taskRoomInfo *EaRoomTaskReportService) DeleteEaTaskRoomInfoById(id int64) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaRoomTaskReportDao.DeleteById(gtx, id)
	})
	return
}

// 导出服务 -- start
type EaTaskRoomInfoExportStruct struct {
	Query *bizDto.EaRoomTaskReportQuery
}

func (h *EaTaskRoomInfoExportStruct) Title() string {
	return "考场征集数据导出"
}

func (h *EaTaskRoomInfoExportStruct) HeaderColumn() []string {
	return []string{"任务名称,taskName", "教室名称,room_name", "所属校区,campusName", "层次,levelName", "组数,group_num", "按组容量,group_capacity", "按位容量,capacity", "所属学院,deptName", "监考老师数量,tch_num", "审核状态,stateStr"}
}

func (h *EaTaskRoomInfoExportStruct) Rows(p *page.Page) []map[string]interface{} {
	result := eaRoomTaskReportDao.ListMapWithCondition(h.Query, p, nil, nil)
	roomMap := make(map[int64]*bizDo.EaRoomInfo, 0)
	deptMap := make(map[int64]*bizDo.EaDeptInfo, 0)
	roomTaskInfo := eaRoomTaskInfoDao.GetById(h.Query.TaskId)
	// 获取任务对应的层次名称
	levelName := ""
	if len(roomTaskInfo.LevelCode) > 0 {
		levelQuery := &userDto.SysDictDataQuery{DictType: cst.BIZ_LEVEL, DictValue: roomTaskInfo.LevelCode, State: strconv.Itoa(cst.YES)}
		levelQuery.PageSize = 1
		levelList := sysDicDataService.ListSysDictData(levelQuery, &page.Page{CurPage: 1, PageSize: 1})
		if len(levelList) > 0 {
			levelName = levelList[0].DictLabel
		}
	}
	for _, res := range result {
		res["taskName"] = roomTaskInfo.TaskName
		res["levelName"] = levelName
		if _, ok := res["state"]; ok {
			switch res["state"] {
			case "1", int64(1):
				res["stateStr"] = "已上报"
				// 仅上报未征集时，监考老师数量显示 --
				res["tch_num"] = "--"
			case "2", int64(2):
				res["stateStr"] = "已征集"
			}
		}
		if _, ok := res["room_id"]; ok {
			roomId := res["room_id"].(int64)
			if roomMap[roomId] == nil {
				roomMap[roomId] = eaRoomInfoDao.GetById(roomId)
			}
			roomInfo := roomMap[roomId]
			if roomInfo != nil {
				res["campusName"] = roomInfo.CampusName
			}
		}
		if _, ok := res["dept_id"]; ok {
			deptId := res["dept_id"].(int64)
			if deptId == cst.BIZ_DEPTID_FOR_ADMISSIONS {
				res["deptName"] = "公共教室"
			} else {
				if deptMap[deptId] == nil {
					deptMap[deptId] = eaDeptInfoDao.GetById(deptId)
				}
				deptInfo := deptMap[deptId]
				if deptInfo != nil {
					res["deptName"] = deptInfo.DeptName
				}
			}
		}
	}
	return result
}

func (h *EaTaskRoomInfoExportStruct) Finish(url string) {
}

// 导出服务 -- end
func (taskRoomInfo *EaRoomTaskReportService) ExportEaTaskRoomInfo(q *bizDto.EaRoomTaskReportQuery) ([]byte, error) {
	return excel.NewExcel().ExportSync(&EaTaskRoomInfoExportStruct{q})
}

// ResetEaTaskRoomInfoById 重置单条征集确认：已征集(2)→已上报(1)，清空监考老师数量
func (taskRoomInfo *EaRoomTaskReportService) ResetEaTaskRoomInfoById(id int64, updateBy int64) (bool, string) {
	dbRoom := eaRoomTaskReportDao.GetById(id)
	if dbRoom == nil {
		return false, "数据不存在，请刷新重试"
	}
	if dbRoom.State != cst.CONFIRM {
		return false, "只有已征集状态才能重置"
	}
	if checkRoomReportHasArrange(id) {
		return false, dbRoom.RoomName + "已经被编排，不允许重置"
	}

	updateDO := &bizDo.EaRoomTaskReport{
		Id:     id,
		State:  cst.UNCONFIRM,
		TchNum: 0,
	}
	updateDO.SetUpdateBy(updateBy)
	updateDO.SetToZeroCols("tch_num")

	return tx.ExecGTx(func(gtx tx.Db) {
		eaRoomTaskReportDao.UpdateNew(gtx, &bizDo.EaRoomTaskReport{Id: id}, updateDO)
		eaRoomTaskInfoDao.UpdateCollectNumAndConfirmNum(gtx, dbRoom.TaskId, 0, -1)
	}), ""
}

// BatchResetEaTaskRoomInfo 批量重置征集确认：已征集(2)→已上报(1)，清空监考老师数量
func (taskRoomInfo *EaRoomTaskReportService) BatchResetEaTaskRoomInfo(list []*bizDo.EaRoomTaskReport, taskId int64, updateBy int64) (bool, string) {
	list2Update := make([]*bizDo.EaRoomTaskReport, 0)
	resetCount := 0
	for _, item := range list {
		if item.Id <= 0 {
			continue
		}
		dbRoom := eaRoomTaskReportDao.GetById(item.Id)
		if dbRoom == nil {
			return false, "存在不存在的考场数据，请刷新重试"
		}
		if dbRoom.State != cst.CONFIRM {
			return false, dbRoom.RoomName + "不是已征集状态，无法重置"
		}
		if checkRoomReportHasArrange(item.Id) {
			return false, dbRoom.RoomName + "已经被编排，不允许重置"
		}
		temp := &bizDo.EaRoomTaskReport{
			Id:     item.Id,
			State:  cst.UNCONFIRM,
			TchNum: 0,
		}
		temp.SetUpdateBy(updateBy)
		temp.SetToZeroCols("tch_num")
		list2Update = append(list2Update, temp)
		resetCount++
	}

	if len(list2Update) <= 0 {
		return true, ""
	}

	return tx.ExecGTx(func(gtx tx.Db) {
		for _, item := range list2Update {
			eaRoomTaskReportDao.UpdateNew(gtx, &bizDo.EaRoomTaskReport{Id: item.Id}, item)
		}
		eaRoomTaskInfoDao.UpdateCollectNumAndConfirmNum(gtx, taskId, 0, -resetCount)
	}), ""
}

// ConfirmEaTaskRoomInfoById 单条征集确认：已上报(1)→已征集(2)，设置监考老师数量
func (taskRoomInfo *EaRoomTaskReportService) ConfirmEaTaskRoomInfoById(id int64, tchNum int64, updateBy int64) (bool, string) {
	dbRoom := eaRoomTaskReportDao.GetById(id)
	if dbRoom == nil {
		return false, "数据不存在，请刷新重试"
	}
	if dbRoom.State != cst.UNCONFIRM {
		return false, "只有已上报状态才能确认征集"
	}

	updateDO := &bizDo.EaRoomTaskReport{
		Id:     id,
		State:  cst.CONFIRM,
		TchNum: int(tchNum),
	}
	updateDO.SetUpdateBy(updateBy)

	return tx.ExecGTx(func(gtx tx.Db) {
		eaRoomTaskReportDao.UpdateNew(gtx, &bizDo.EaRoomTaskReport{Id: id}, updateDO)
		eaRoomTaskInfoDao.UpdateCollectNumAndConfirmNum(gtx, dbRoom.TaskId, 0, 1)
	}), ""
}

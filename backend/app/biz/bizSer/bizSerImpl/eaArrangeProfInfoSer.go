package bizSerImpl

import (
	"encoding/json"
	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/utils/snowflake"
	"strconv"
	"time"
	"unicode/utf8"
	"yx-exam-affair-gm/app/biz/bizDao"
	"yx-exam-affair-gm/app/biz/bizDao/bizDaoImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	"yx-exam-affair-gm/app/cst"
	"yx-exam-affair-gm/app/hint"
	"yx-exam-affair-gm/app/user/model/userDo"
	cmutils "yx-exam-affair-gm/app/utils"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
	"github.com/sealsee/web-base/public/utils/file/excel"
)

/*
 * desc: 编排专业服务层实现
 * author: init code
 * date: 2024/09/24
 */
type EaArrangeProfInfoService struct{}

var eaArrangeProfInfoService *EaArrangeProfInfoService
var eaArrangeProfInfoDao bizDao.IEaArrangeProfInfoDao

func init() {
	eaArrangeProfInfoService = &EaArrangeProfInfoService{}
	eaArrangeProfInfoDao = bizDaoImpl.NewEaArrangeProfInfoDao()
}

func NewEaArrangeProfInfoService() *EaArrangeProfInfoService {
	return eaArrangeProfInfoService
}

func (arrangeProfInfo *EaArrangeProfInfoService) GetEaArrangeProfInfoById(id int64) *bizDo.EaArrangeProfInfo {
	dbProfInfo := eaArrangeProfInfoDao.GetById(id)
	if dbProfInfo == nil {
		return nil
	}
	dbRoomTaskInfo := eaRoomTaskInfoService.GetEaTaskInfoById(dbProfInfo.RoomTaskId)
	if dbRoomTaskInfo == nil {
		return dbProfInfo
	}
	dbProfInfo.RoomTaskName = dbRoomTaskInfo.TaskName

	// 查询当前任务下最大考场编号
	maxArrangeNoRoomInfo := eaArrangeRoomService.GetMaxArrangeRoomNo(dbProfInfo.RoomTaskId)
	if maxArrangeNoRoomInfo != nil {
		dbProfInfo.MaxArrangeRoomNoStr = maxArrangeNoRoomInfo.ArrangeNoStr
		dbProfInfo.MaxArrangeRoomNo = maxArrangeNoRoomInfo.ArrangeNo
	}
	return dbProfInfo
}

func (arrangeProfInfo *EaArrangeProfInfoService) ListEaArrangeProfInfo(q *bizDto.EaArrangeProfInfoQuery, p *page.Page) []*bizDo.EaArrangeProfInfo {
	q.AddLikeAll("dept_name", q.DeptName)
	q.AddLikeAll("prof_name", q.ProfName)
	q.AddLikeAll("direction_name", q.DirectionName)
	return eaArrangeProfInfoDao.List(q, p)
}

func (arrangeProfInfo *EaArrangeProfInfoService) InsertEaArrangeProfInfo(arrangeProfInfoDO *bizDo.EaArrangeProfInfo) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaArrangeProfInfoDao.Insert(gtx, arrangeProfInfoDO)
	})
	return
}

func (arrangeProfInfo *EaArrangeProfInfoService) UpdateEaArrangeProfInfoById(arrangeProfInfoDO *bizDo.EaArrangeProfInfo) (success bool) {
	uWhere := &bizDo.EaArrangeProfInfo{InfoId: arrangeProfInfoDO.InfoId}
	arrangeProfInfoDO.InfoId = 0
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaArrangeProfInfoDao.Update(gtx, arrangeProfInfoDO, uWhere)
	})
	return
}

func (arrangeProfInfo *EaArrangeProfInfoService) CancelArrangeById(arrangeProfInfoDO *bizDo.EaArrangeProfInfo) (success bool, errInfo errs.ERROR) {
	success = false

	// 查询专业编排信息
	dbArrangeProfInfo := eaArrangeProfInfoDao.GetById(arrangeProfInfoDO.InfoId)
	if dbArrangeProfInfo == nil {
		return success, hint.BIZ_NO_ARRANGE_ROOM_ERR
	}

	// 查询日程下以及分配的列表
	allArrangeRoomReportIdList := make([]interface{}, 0)
	allArrangeRoomList := eaArrangeRoomService.ListAllArrangeRoom(&bizDto.EaArrangeRoomQuery{
		InfoId: dbArrangeProfInfo.InfoId,
	})
	for _, arrangeRoomInfo := range allArrangeRoomList {
		allArrangeRoomReportIdList = append(allArrangeRoomReportIdList, arrangeRoomInfo.RoomReportId)
	}
	if len(allArrangeRoomReportIdList) > 0 {
		arrangeTchQuery := &bizDto.EaArrangeTchQuery{}
		arrangeTchQuery.AddIn("room_report_id", allArrangeRoomReportIdList...)
		listEaArrangeTch := eaArrangeTchService.ListEaArrangeTch(arrangeTchQuery, page.NewPage())
		if len(listEaArrangeTch) > 0 {
			return success, hint.BIZ_HAS_TCH_ARRANGE_ERR
		}
	}

	uWhere := &bizDo.EaArrangeProfInfo{InfoId: arrangeProfInfoDO.InfoId}

	updateArrangeProfInfo := &bizDo.EaArrangeProfInfo{ArrangeState: cst.UNARRANGE, ArrangeStuNum: 0, ArrangeRoomNum: 0}
	updateArrangeProfInfo.SetToZeroCols("arrange_stu_num", "arrange_room_num")
	updateArrangeProfInfo.SetUpdateBy(arrangeProfInfoDO.UpdateBy)
	success = tx.ExecGTx(func(gtx tx.Db) {
		eaArrangeProfInfoDao.UpdateNew(gtx, uWhere, updateArrangeProfInfo)
		for _, arrangeRoomInfo := range allArrangeRoomList {
			eaArrangeRoomDao.DeleteById(gtx, arrangeRoomInfo.RoomArrangeId, arrangeProfInfoDO.UpdateBy)
		}
	})
	return
}

func (arrangeProfInfo *EaArrangeProfInfoService) DeleteEaArrangeProfInfoById(id int64, deleteBy int64) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaArrangeProfInfoDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

// 导出服务 -- start
type EaArrangeProfInfoExportStruct struct {
	Query *bizDto.EaArrangeProfInfoQuery
}

func (h *EaArrangeProfInfoExportStruct) Title() string {
	return "编排专业导出"
}

func (h *EaArrangeProfInfoExportStruct) HeaderColumn() []string {
	return []string{"数据ID,info_id", "考场任务ID,room_task_id", "征集日期,task_day", "所属学院id,dept_id", "学院名称,dept_name", "专业名称,prof_name", "方向名称,direction_name", "安排类型: 1-按组 2-按位,arrange_type", "考生数量,stu_num", "已编排考生数量,arrange_stu_num", "已编排考场数量,arrange_room_num", "考场编排状态: 1-待编排 2-已编排,arrange_state", "删除标记：1-正常的；2-已删除的；,deleted", "备注,remark", "创建者,create_by", "创建时间,create_time", "更新者,update_by", "更新时间,update_time"}
}

func (h *EaArrangeProfInfoExportStruct) Rows(p *page.Page) []map[string]interface{} {
	return eaArrangeProfInfoDao.ListMapWithCondition(h.Query, p, nil, nil)
}

func (h *EaArrangeProfInfoExportStruct) Finish(url string) {
}

// 导出服务 -- end
func (arrangeProfInfo *EaArrangeProfInfoService) ExportEaArrangeProfInfo(q *bizDto.EaArrangeProfInfoQuery) ([]byte, error) {
	return excel.NewExcel().ExportSync(&EaArrangeProfInfoExportStruct{q})
}

// 导入服务 -- start
type HandlerImpProfInfoStruct struct {
	BizType          int
	BatchNo          string
	ExistRoomInfoMap map[string]*bizDo.EaArrangeProfInfo
	Query            *bizDto.EaArrangeProfInfoQuery
	List             []*bizDo.EaArrangeProfInfo
	Msg              string
	Err              string
	ExcelHeader      []string
}

func (h *HandlerImpProfInfoStruct) Headers(headers []string) {
	if len(headers) != len(h.ExcelHeader) {
		h.Err += "表头数据错误!"
		return
	}
	for i, header := range headers {
		if header != h.ExcelHeader[i] {
			h.Err += "表头数据错误!"
		}
	}
}

func (h *HandlerImpProfInfoStruct) Row(mp *map[string]string) excel.RunStatus {
	// 判断表头数据是否错误
	if h.Err != "" {
		return excel.Exit
	}
	m := *mp
	rowIdx := strconv.Itoa(cmutils.StrToInt(m["rid"]) + 1)
	proInfo2Insert := new(bizDo.EaArrangeProfInfo)

	// 取值
	roomTaskId := h.Query.RoomTaskId
	taskDay := h.Query.TaskDay
	deptName := m[h.ExcelHeader[0]]
	profName := m[h.ExcelHeader[1]]
	directionName := m[h.ExcelHeader[2]]
	arrangeTypeStr := m[h.ExcelHeader[3]]
	stuNum := m[h.ExcelHeader[4]]

	if len(deptName) <= 0 {
		h.Err = "第" + rowIdx + "行，学院名称不能为空！"
		return excel.Exit
	}
	//if h.DeptNameMap[deptName] == nil {
	//	h.Err = "第" + rowIdx + "行，学院不存在"
	//	return excel.Exit
	//}
	//deptId := h.DeptNameMap[deptName].DeptId
	proInfo2Insert.RoomTaskId = roomTaskId
	proInfo2Insert.TaskDay = taskDay
	//proInfo2Insert.DeptId = deptId
	proInfo2Insert.DeptName = deptName

	if len(profName) <= 0 {
		h.Err = "第" + rowIdx + "行，专业名称不能为空！"
		return excel.Exit
	}
	if utf8.RuneCountInString(profName) > 64 {
		h.Err = "第" + rowIdx + "行，专业名称不能超过64个字符！"
		return excel.Exit
	}
	proInfo2Insert.ProfName = profName

	if len(directionName) > 0 {
		if utf8.RuneCountInString(directionName) > 64 {
			h.Err = "第" + rowIdx + "行，方向名称不能超过64个字符！"
			return excel.Exit
		}
		proInfo2Insert.DirectionName = directionName
	}

	// 唯一校验
	key := strconv.FormatInt(roomTaskId, 10) + "_" + deptName + "_" + profName + "_" + directionName
	if h.ExistRoomInfoMap[key] != nil {
		h.Err = "第" + rowIdx + "行，专业方向信息已存在！"
		return excel.Exit
	}
	h.ExistRoomInfoMap[key] = proInfo2Insert

	if len(arrangeTypeStr) <= 0 {
		h.Err = "第" + rowIdx + "行，类型不能为空！"
		return excel.Exit
	}
	if arrangeTypeStr == cst.RoomTypeMap[cst.BIZ_TYPE_GROUP] {
		proInfo2Insert.ArrangeType = int(cst.BIZ_TYPE_GROUP)
	} else if arrangeTypeStr == cst.RoomTypeMap[cst.BIZ_TYPE_PERSON] {
		proInfo2Insert.ArrangeType = int(cst.BIZ_TYPE_PERSON)
	} else {
		h.Err = "第" + rowIdx + "行，类型的值只能为按组或者按位！"
		return excel.Exit
	}

	if len(stuNum) > 0 {
		if !cmutils.IsPositiveInteger(stuNum) {
			h.Err = "第" + rowIdx + "行，考生数只能是正整数！"
			return excel.Exit
		}

		proInfo2Insert.StuNum = cmutils.StrToInt(stuNum)
	}

	// 创建人、时间
	proInfo2Insert.Deleted = common.Normal
	proInfo2Insert.ArrangeRoomNum = 0
	proInfo2Insert.ArrangeStuNum = 0
	proInfo2Insert.ArrangeState = cst.UNARRANGE
	proInfo2Insert.SetCreateBy(h.Query.OperId)
	proInfo2Insert.SetUpdateBy(h.Query.OperId)
	// 放入list
	h.List = append(h.List, proInfo2Insert)
	return excel.Next
}

func (h *HandlerImpProfInfoStruct) After() {
	if len(h.Err) > 0 {
		return
	}
	// 判断是否有导入数据
	if len(h.List) < 1 {
		h.Err = "没有导入数据！"
		return
	}
	h.Msg += "成功导入 " + strconv.Itoa(len(h.List)) + " 条记录!\n"
	// 拼装导入记录
	dataInout := new(userDo.SysDataInout)
	dataInout.OperType = common.OPER_TYPE_IMPORT_DATA
	dataInout.BizType = h.BizType
	dataInout.BatchNo = h.BatchNo
	// 压缩导入的多个文件
	if len(h.Query.Files) > 1 {
		dataInout.FileIn = cmutils.CompressZipFiles(h.Query.Files)
		dataInout.FileInName = "多文件"
	} else {
		for k, v := range h.Query.Files {
			dataInout.FileIn = v
			dataInout.FileInName = k
		}
	}
	// 传参json化
	h.Query.Files = nil // 文件不放在条件参数里记录
	p, _ := json.Marshal(h.Query)
	pm := string(p)
	if len(pm) > 1000 {
		dataInout.Param = pm[:1000] //截取1000个字符长度
	} else {
		dataInout.Param = pm
	}
	dataInout.OperTime = time.Now()
	dataInout.State = common.OK
	dataInout.Remark = h.Msg
	dataInout.SetCreateBy(h.Query.OperId)
	dataInout.SetUpdateBy(h.Query.OperId)
	// 数据库执行事务
	success := tx.ExecGTx(func(gtx tx.GTx) {
		// 最后批量插入导入的数据，并新增一条导入记录
		eaArrangeProfInfoDao.BatchInsert(gtx, h.List)
		sysDataInoutDao.Insert(gtx, dataInout)
	})
	if !success {
		h.Err = "导入失败，请检查是否与已有数据重复！"
	}
}

// 导入服务  --  end

func (arrangeProfInfo *EaArrangeProfInfoService) ImportEaArrangeProfInfo(q *bizDto.EaArrangeProfInfoQuery) (bool, string) {

	// 查询任务详情
	roomTask := eaRoomTaskInfoService.GetEaTaskInfoById(q.RoomTaskId)
	if roomTask == nil {
		return false, "考场任务不存在！"
	}
	q.TaskDay = roomTask.TaskDay

	// 使用雪花算法生成全局唯一id，作为本次导入操作的批次号
	currentBatchNo := strconv.FormatInt(snowflake.GenID(), 10)

	// 查询所有学院信息放入map
	//deptNameMap := make(map[string]*bizDo.EaDeptInfo)
	//allDeptList := eaDeptInfoService.ListAllEaDeptInfo(&bizDto.EaDeptInfoQuery{})
	//for _, dept := range allDeptList {
	//	deptNameMap[dept.DeptName] = dept
	//}

	// 查询所有已存在的专业方向信息
	existProfInfoMap := make(map[string]*bizDo.EaArrangeProfInfo)
	ExistProfInfoList := arrangeProfInfo.ListAllProfInfo(&bizDto.EaArrangeProfInfoQuery{})
	for _, prof := range ExistProfInfoList {
		key := strconv.FormatInt(prof.RoomTaskId, 10) + "_" + prof.DeptName + "_" + prof.ProfName + "_" + prof.DirectionName
		existProfInfoMap[key] = prof
	}

	// 遍历文件map  文件名->文件地址...
	excelHeader := []string{"学院", "专业", "方向", "类型", "考生数"}
	handlerImpStruct := &HandlerImpProfInfoStruct{BizType: int(cst.BIZ_ARRANGE_PROF_INFO), BatchNo: currentBatchNo,
		Query: q, ExcelHeader: excelHeader, ExistRoomInfoMap: existProfInfoMap}
	for _, v := range q.Files {
		// 根据url导入excel处理
		excel.NewExcel().ImportWithUrl(v, handlerImpStruct)
	}
	if len(handlerImpStruct.Err) > 0 {
		return false, handlerImpStruct.Err
	}
	return true, handlerImpStruct.Msg
}

func (deptInfo *EaArrangeProfInfoService) ListAllProfInfo(q *bizDto.EaArrangeProfInfoQuery) []*bizDo.EaArrangeProfInfo {
	listDepartAll := make([]*bizDo.EaArrangeProfInfo, 0)
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE
	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		list := eaArrangeProfInfoDao.List(q, pageInfo)
		if len(list) == 0 {
			break
		}
		listDepartAll = append(listDepartAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}
	return listDepartAll
}

func (deptInfo *EaArrangeProfInfoService) ListAvailableRoomInTask(q *bizDto.EaArrangeProfInfoQuery) []*bizDo.EaRoomTaskReport {

	// 查询当前待编排任务下的详情信息
	dbProfInfo := eaArrangeProfInfoService.GetEaArrangeProfInfoById(q.InfoId)
	if dbProfInfo == nil || dbProfInfo.RoomTaskId < 1 || dbProfInfo.ArrangeType < 1 {
		return make([]*bizDo.EaRoomTaskReport, 0)
	}

	// 查询任务下所有已确认的教室信息
	reportQuery := &bizDto.EaRoomTaskReportQuery{
		TaskId: dbProfInfo.RoomTaskId,
		State:  cst.CONFIRM,
	}
	if dbProfInfo.ArrangeType == int(cst.BIZ_TYPE_GROUP) {
		reportQuery.AddGt("group_capacity", 0)
	} else {
		reportQuery.AddGt("capacity", 0)
	}
	reportedRoomList := eaRoomTaskReportService.ListAllRoomReport(reportQuery)
	if len(reportedRoomList) < 1 {
		return make([]*bizDo.EaRoomTaskReport, 0)
	}
	curInfoRoomMap := make(map[int64]bool)
	for _, room := range reportedRoomList {
		curInfoRoomMap[room.RoomId] = true
	}

	// 查询当前方向任务下所有已编排的教室已使用容量
	roomUsedNumMap := make(map[int64]int)
	groupRoomUsedNumMap := make(map[int64]int)
	allRoomUsedNumMap := make(map[int64]int)
	groupAllRoomUsedNumMap := make(map[int64]int)
	allArrangeRoom := eaArrangeRoomService.ListAllArrangeRoom(&bizDto.EaArrangeRoomQuery{RoomTaskId: dbProfInfo.RoomTaskId})
	for _, room := range allArrangeRoom {
		if room.ArrangeType == int(cst.BIZ_TYPE_GROUP) {
			if room.InfoId == dbProfInfo.InfoId {
				groupRoomUsedNumMap[room.RoomId] = groupRoomUsedNumMap[room.RoomId] + room.StuNum
			}
			groupAllRoomUsedNumMap[room.RoomId] = groupAllRoomUsedNumMap[room.RoomId] + room.StuNum
		} else {
			if room.InfoId == dbProfInfo.InfoId {
				roomUsedNumMap[room.RoomId] = roomUsedNumMap[room.RoomId] + room.StuNum
			}
			allRoomUsedNumMap[room.RoomId] = allRoomUsedNumMap[room.RoomId] + room.StuNum
		}
	}

	// 遍历组装数据
	showReportedRoomList := make([]*bizDo.EaRoomTaskReport, 0)
	for _, room := range reportedRoomList {

		// 不显示任务下的已编排的教室
		if dbProfInfo.ArrangeType == int(cst.BIZ_TYPE_GROUP) {
			if groupRoomUsedNumMap[room.RoomId] > 0 {
				continue
			}
		} else {
			if roomUsedNumMap[room.RoomId] > 0 {
				continue
			}
		}

		unUsedCapacity := room.Capacity - allRoomUsedNumMap[room.RoomId]
		groupUnUsedCapacity := room.GroupCapacity - groupAllRoomUsedNumMap[room.RoomId]
		if dbProfInfo.ArrangeType == int(cst.BIZ_TYPE_GROUP) {

			// 如果当前教室，已经按位分配了容量，则不显示
			if allRoomUsedNumMap[room.RoomId] > 0 {
				continue
			}
			if groupUnUsedCapacity < 1 {
				continue
			}
		} else {
			if groupAllRoomUsedNumMap[room.RoomId] > 0 {
				continue
			}
			if unUsedCapacity < 1 {
				continue
			}
		}

		room.GroupCapacity = groupUnUsedCapacity
		room.Capacity = unUsedCapacity
		showReportedRoomList = append(showReportedRoomList, room)
	}
	return showReportedRoomList
}

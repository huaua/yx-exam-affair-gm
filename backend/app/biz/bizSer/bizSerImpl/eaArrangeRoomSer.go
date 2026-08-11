package bizSerImpl

import (
	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/errs"
	"strconv"
	"strings"
	"yx-exam-affair-gm/app/biz/bizDao"
	"yx-exam-affair-gm/app/biz/bizDao/bizDaoImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	"yx-exam-affair-gm/app/cst"
	"yx-exam-affair-gm/app/hint"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
	"github.com/sealsee/web-base/public/utils/file/excel"
)

/*
 * desc: 考场编排详情服务层实现
 * author: init code
 * date: 2024/09/24
 */
type EaArrangeRoomService struct{}

var eaArrangeRoomService *EaArrangeRoomService
var eaArrangeRoomDao bizDao.IEaArrangeRoomDao

func init() {
	eaArrangeRoomService = &EaArrangeRoomService{}
	eaArrangeRoomDao = bizDaoImpl.NewEaArrangeRoomDao()
}

func NewEaArrangeRoomService() *EaArrangeRoomService {
	return eaArrangeRoomService
}

func (arrangeRoom *EaArrangeRoomService) GetEaArrangeRoomById(id int64) *bizDo.EaArrangeRoom {
	return eaArrangeRoomDao.GetById(id)
}

func (arrangeRoom *EaArrangeRoomService) GetMaxArrangeRoomNo(roomTaskId int64) *bizDo.EaArrangeRoom {
	arrangeRoomQuery := &bizDto.EaArrangeRoomQuery{
		RoomTaskId: roomTaskId,
	}
	arrangeRoomQuery.AddOrderDesc("arrange_no")
	arrangeRooms := eaArrangeRoomDao.List(arrangeRoomQuery, page.NewPage())
	if len(arrangeRooms) < 1 {
		return nil
	}
	return arrangeRooms[0]
}

func (arrangeRoom *EaArrangeRoomService) ListEaArrangeRoom(q *bizDto.EaArrangeRoomQuery) []*bizDo.EaArrangeRoom {

	// 查询当前待编排任务下的详情信息
	dbProfInfo := eaArrangeProfInfoService.GetEaArrangeProfInfoById(q.InfoId)
	if dbProfInfo == nil || dbProfInfo.RoomTaskId < 1 || dbProfInfo.ArrangeType < 1 {
		return make([]*bizDo.EaArrangeRoom, 0)
	}

	// 查询已编排的教室列表
	listAllArrangeRoom := arrangeRoom.ListAllArrangeRoom(&bizDto.EaArrangeRoomQuery{
		RoomTaskId: dbProfInfo.RoomTaskId,
	})
	if len(listAllArrangeRoom) < 1 {
		return make([]*bizDo.EaArrangeRoom, 0)
	}

	// 组装教室已用容量
	roomUsedNumMap := make(map[int64]int)
	groupRoomUsedNumMap := make(map[int64]int)
	for _, room := range listAllArrangeRoom {
		if room.ArrangeType == int(cst.BIZ_TYPE_GROUP) {
			groupRoomUsedNumMap[room.RoomId] = groupRoomUsedNumMap[room.RoomId] + room.StuNum
		} else {
			roomUsedNumMap[room.RoomId] = roomUsedNumMap[room.RoomId] + room.StuNum
		}
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
		return make([]*bizDo.EaArrangeRoom, 0)
	}
	roomReportMap := make(map[int64]*bizDo.EaRoomTaskReport)
	for _, room := range reportedRoomList {
		roomReportMap[room.RoomId] = room
	}

	newListAllArrangeRoom := make([]*bizDo.EaArrangeRoom, 0)
	for _, arrangeRoom := range listAllArrangeRoom {
		roomReportCapacity := 0
		if arrangeRoom.InfoId != q.InfoId {
			continue
		}
		arrangeRoom.GroupNum = roomReportMap[arrangeRoom.RoomId].GroupNum
		unUsedCapacity := 0
		if dbProfInfo.ArrangeType == int(cst.BIZ_TYPE_GROUP) {
			roomReportCapacity = roomReportMap[arrangeRoom.RoomId].GroupCapacity
			arrangeRoom.GroupCapacity = roomReportCapacity - groupRoomUsedNumMap[arrangeRoom.RoomId] + arrangeRoom.StuNum
			arrangeRoom.Capacity = roomReportMap[arrangeRoom.RoomId].Capacity - roomUsedNumMap[arrangeRoom.RoomId]
			unUsedCapacity = roomReportCapacity - groupRoomUsedNumMap[arrangeRoom.RoomId]
		} else {
			roomReportCapacity = roomReportMap[arrangeRoom.RoomId].Capacity
			arrangeRoom.GroupCapacity = roomReportMap[arrangeRoom.RoomId].GroupCapacity - groupRoomUsedNumMap[arrangeRoom.RoomId]
			arrangeRoom.Capacity = roomReportCapacity - roomUsedNumMap[arrangeRoom.RoomId] + arrangeRoom.StuNum
			unUsedCapacity = roomReportCapacity - roomUsedNumMap[arrangeRoom.RoomId]
		}

		if unUsedCapacity < arrangeRoom.StuNum {
			// 上报总容量 - 已用容量 < 当前分配容量，显示当前分配容量
			if unUsedCapacity > 0 {
				arrangeRoom.RoomCapacity = arrangeRoom.StuNum + unUsedCapacity
			} else {
				arrangeRoom.RoomCapacity = arrangeRoom.StuNum
			}
		} else {
			// 上报总容量 - 已用容量 > 当前分配容量，显示上报总容量 - 已用容量 + 已分配容量
			arrangeRoom.RoomCapacity = unUsedCapacity + arrangeRoom.StuNum
		}
		newListAllArrangeRoom = append(newListAllArrangeRoom, arrangeRoom)
	}

	return newListAllArrangeRoom
}

//func (arrangeRoom *EaArrangeRoomService) ListAssignment(q *bizDto.EaArrangeRoomQuery) []*bizDo.EaArrangeRoom {
//	newSql := ""
//}

func (arrangeRoom *EaArrangeRoomService) ListAllArrangeRoom(q *bizDto.EaArrangeRoomQuery) []*bizDo.EaArrangeRoom {
	listDepartAll := make([]*bizDo.EaArrangeRoom, 0)
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE
	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		list := eaArrangeRoomDao.List(q, pageInfo)
		if len(list) == 0 {
			break
		}
		listDepartAll = append(listDepartAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}
	return listDepartAll
}

func (arrangeRoom *EaArrangeRoomService) InsertEaArrangeRoom(arrangeRoomDO *bizDo.EaArrangeRoom) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaArrangeRoomDao.Insert(gtx, arrangeRoomDO)
	})
	return
}

func (arrangeRoom *EaArrangeRoomService) BatchInsertArrangeRoom(arrangeRoomAddDto *bizDto.EaArrangeRoomAddDto, optId int64) (success bool, errInfo errs.ERROR) {
	if arrangeRoomAddDto.InfoId < 1 || len(arrangeRoomAddDto.ArrangeRoomList) < 1 {
		return false, errs.PARAM_INVALID_ERR
	}

	// 获取任务详情
	dbProfInfo := eaArrangeProfInfoService.GetEaArrangeProfInfoById(arrangeRoomAddDto.InfoId)
	if dbProfInfo == nil || dbProfInfo.RoomTaskId < 1 || dbProfInfo.ArrangeType < 1 {
		return false, hint.BIZ_NO_ARRANGE_ROOM_ERR
	}

	if dbProfInfo.ArrangeState == cst.ARRANGED {
		return false, hint.BIZ_ARRANGE_ROOM_ADD_ERR
	}

	stuNum := 0
	for index, arrangeRoom := range arrangeRoomAddDto.ArrangeRoomList {
		arrangeRoom.InfoId = dbProfInfo.InfoId

		if len(arrangeRoom.ArrangeNoStr) < 1 {
			return false, hint.BIZ_ARRANGE_ROOM_NO_FORMAT_ERR.Format(index + 1)
		}
		trimmedArrangeNoStr := strings.TrimLeft(arrangeRoom.ArrangeNoStr, "0")
		if len(trimmedArrangeNoStr) < 1 {
			return false, hint.BIZ_ARRANGE_ROOM_NO_FORMAT_ERR.Format(index + 1)
		}
		arrangeNoInt, err := strconv.Atoi(trimmedArrangeNoStr)
		if err != nil {
			return false, hint.BIZ_ARRANGE_ROOM_NO_FORMAT_ERR.Format(index + 1)
		}
		arrangeRoom.ArrangeNo = arrangeNoInt
		arrangeRoom.SetCreateBy(optId)
		arrangeRoom.SetUpdateBy(optId)
		stuNum = stuNum + arrangeRoom.StuNum
	}
	roomNum := len(arrangeRoomAddDto.ArrangeRoomList)
	profInfo := &bizDo.EaArrangeProfInfo{
		StuNum:         stuNum,
		ArrangeRoomNum: roomNum,
		ArrangeState:   cst.ARRANGED,
	}
	profInfo.SetUpdateBy(optId)
	profInfoWhere := &bizDo.EaArrangeProfInfo{
		InfoId: dbProfInfo.InfoId,
	}

	success = tx.ExecGTx(func(gtx tx.GTx) {
		eaArrangeRoomDao.BatchInsert(gtx, arrangeRoomAddDto.ArrangeRoomList)
		eaArrangeProfInfoDao.Update(gtx, profInfo, profInfoWhere)
	})
	return
}

func (arrangeRoom *EaArrangeRoomService) UpdateEaArrangeRoomById(arrangeRoomDO *bizDo.EaArrangeRoom) (success bool) {
	uWhere := &bizDo.EaArrangeRoom{RoomArrangeId: arrangeRoomDO.RoomArrangeId}
	arrangeRoomDO.RoomArrangeId = 0
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaArrangeRoomDao.Update(gtx, arrangeRoomDO, uWhere)
	})
	return
}

func (arrangeRoom *EaArrangeRoomService) BatchUpdateArrangeRoom(arrangeRoomAddDto *bizDto.EaArrangeRoomAddDto, optId int64) (success bool, errInfo errs.ERROR) {
	if arrangeRoomAddDto.InfoId < 1 || len(arrangeRoomAddDto.ArrangeRoomList) < 1 {
		return false, errs.PARAM_INVALID_ERR
	}

	// 获取任务详情
	dbProfInfo := eaArrangeProfInfoService.GetEaArrangeProfInfoById(arrangeRoomAddDto.InfoId)
	if dbProfInfo == nil || dbProfInfo.RoomTaskId < 1 || dbProfInfo.ArrangeType < 1 {
		return false, hint.BIZ_NO_ARRANGE_ROOM_ERR
	}

	// 查询所有当前专业方向任务下已编排的教室列表
	allArrangeRoomList := make([]*bizDo.EaArrangeRoom, 0)
	arrangeRoomIdMap := make(map[int64]*bizDo.EaArrangeRoom)
	dbArrangeNoMap := make(map[int]bool)
	arrangeRoomQuery := &bizDto.EaArrangeRoomQuery{RoomTaskId: dbProfInfo.RoomTaskId}
	allArrangeRoomListTemp := eaArrangeRoomService.ListAllArrangeRoom(arrangeRoomQuery)
	for _, arrangeRoom := range allArrangeRoomListTemp {
		if arrangeRoom.InfoId == dbProfInfo.InfoId {
			arrangeRoomIdMap[arrangeRoom.RoomId] = arrangeRoom
			allArrangeRoomList = append(allArrangeRoomList, arrangeRoom)
		}
		dbArrangeNoMap[arrangeRoom.ArrangeNo] = true
	}

	updateArrangeRoomList := make([]*bizDo.EaArrangeRoom, 0)
	insertArrangeRoomList := make([]*bizDo.EaArrangeRoom, 0)
	delArrangeRoomList := make([]*bizDo.EaArrangeRoom, 0)

	// 遍历本次需要保存的教室数据
	stuNum := 0
	curArrRoomIdMap := make(map[int64]bool)
	for index, arrangeRoom := range arrangeRoomAddDto.ArrangeRoomList {
		stuNum = stuNum + arrangeRoom.StuNum
		curArrRoomIdMap[arrangeRoom.RoomId] = true

		if len(arrangeRoom.ArrangeNoStr) < 1 {
			return false, hint.BIZ_ARRANGE_ROOM_NO_FORMAT_ERR.Format(index + 1)
		}
		trimmedArrangeNoStr := strings.TrimLeft(arrangeRoom.ArrangeNoStr, "0")
		if len(trimmedArrangeNoStr) < 1 {
			return false, hint.BIZ_ARRANGE_ROOM_NO_FORMAT_ERR.Format(index + 1)
		}
		arrangeNoInt, err := strconv.Atoi(trimmedArrangeNoStr)
		if err != nil {
			return false, hint.BIZ_ARRANGE_ROOM_NO_FORMAT_ERR.Format(index + 1)
		}
		arrangeRoom.ArrangeNo = arrangeNoInt

		// 如果教室之前存在，则只更新数据
		if arrangeRoomIdMap[arrangeRoom.RoomId] != nil {
			arrangeRoom.RoomArrangeId = arrangeRoomIdMap[arrangeRoom.RoomId].RoomArrangeId
			arrangeRoom.RoomArrangeId = arrangeRoomIdMap[arrangeRoom.RoomId].RoomArrangeId
			updateArrangeRoomList = append(updateArrangeRoomList, arrangeRoom)
			continue
		}

		arrangeRoom.InfoId = dbProfInfo.InfoId
		arrangeRoom.ArrangeType = dbProfInfo.ArrangeType
		arrangeRoom.AssignmentState = cst.UNASSIGNMENT
		arrangeRoom.Deleted = common.Normal
		arrangeRoom.SetCreateBy(optId)
		arrangeRoom.SetUpdateBy(optId)
		insertArrangeRoomList = append(insertArrangeRoomList, arrangeRoom)
	}

	for _, arrangeRoom := range allArrangeRoomList {
		if _, ok := curArrRoomIdMap[arrangeRoom.RoomId]; !ok {
			delArrangeRoomList = append(delArrangeRoomList, arrangeRoom)
			dbArrangeNoMap[arrangeRoom.ArrangeNo] = false
		}
	}

	// 校验编号是否重复
	for _, arrangeRoom := range updateArrangeRoomList {
		dbRoom := arrangeRoomIdMap[arrangeRoom.RoomId]
		if dbRoom == nil {
			continue
		}
		dbArrangeNoMap[dbRoom.ArrangeNo] = false
	}
	for _, arrangeRoom := range updateArrangeRoomList {
		dbRoom := arrangeRoomIdMap[arrangeRoom.RoomId]
		if dbRoom == nil {
			continue
		}

		// 如果编号和原来库里的一致，则不校验
		if dbRoom.ArrangeNo == arrangeRoom.ArrangeNo {
			dbArrangeNoMap[arrangeRoom.ArrangeNo] = true
			continue
		}

		if dbArrangeNoMap[arrangeRoom.ArrangeNo] {
			return false, hint.BIZ_ARRANGE_ROOM_NO_REPEAT_ERR.Format(arrangeRoom.ArrangeNoStr)
		}
		dbArrangeNoMap[arrangeRoom.ArrangeNo] = true
	}
	for _, arrangeRoom := range insertArrangeRoomList {
		if _, ok := dbArrangeNoMap[arrangeRoom.ArrangeNo]; ok {
			return false, hint.BIZ_ARRANGE_ROOM_NO_REPEAT_ERR.Format(arrangeRoom.ArrangeNoStr)
		}
	}

	// 封装专业方向编排更新对象
	roomNum := len(arrangeRoomAddDto.ArrangeRoomList)
	profInfo := &bizDo.EaArrangeProfInfo{
		ArrangeStuNum:  stuNum,
		ArrangeRoomNum: roomNum,
		ArrangeState:   cst.ARRANGED,
	}
	profInfo.SetUpdateBy(optId)
	profInfo.SetToZeroCols("arrange_stu_num", "arrange_room_num")
	profInfoWhere := &bizDo.EaArrangeProfInfo{
		InfoId: dbProfInfo.InfoId,
	}

	// 操作数据库
	success = tx.ExecGTx(func(gtx tx.Db) {
		if len(insertArrangeRoomList) > 0 {
			eaArrangeRoomDao.BatchInsert(gtx, insertArrangeRoomList)
		}
		for _, delArrRoom := range delArrangeRoomList {
			eaArrangeRoomDao.DeleteById(gtx, delArrRoom.RoomArrangeId, optId)
		}
		for _, updateArrRoom := range updateArrangeRoomList {
			eaArrangeRoomDao.Update(gtx, updateArrRoom, &bizDo.EaArrangeRoom{RoomArrangeId: updateArrRoom.RoomArrangeId})
		}
		eaArrangeProfInfoDao.UpdateNew(gtx, profInfoWhere, profInfo)
	})
	return
}

func (arrangeRoom *EaArrangeRoomService) DeleteEaArrangeRoomById(id int64, deleteBy int64) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaArrangeRoomDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

// 导出服务 -- start
type EaArrangeRoomExportStruct struct {
	Query *bizDto.EaArrangeRoomQuery
}

func (h *EaArrangeRoomExportStruct) Title() string {
	return "考场编排详情导出"
}

func (h *EaArrangeRoomExportStruct) HeaderColumn() []string {
	return []string{"考场安排ID,room_arrange_id", "数据ID,info_id", "考场上报ID,room_report_id", "任务ID,room_task_id", "教室ID,room_id", "教室名称,room_name", "考生数量,stu_num", "考场实际容量,room_capacity", "考场安排编号(字符串),arrange_no_str", "考场安排编号(数值),arrange_no", "需要监考老师数量,need_tch_num", "删除标记：1-正常的；2-已删除的；,deleted", "备注,remark", "创建者,create_by", "创建时间,create_time", "更新者,update_by", "更新时间,update_time"}
}

func (h *EaArrangeRoomExportStruct) Rows(p *page.Page) []map[string]interface{} {
	return eaArrangeRoomDao.ListMapWithCondition(h.Query, p, nil, nil)
}

func (h *EaArrangeRoomExportStruct) Finish(url string) {
}

// 导出服务 -- end
func (arrangeRoom *EaArrangeRoomService) ExportEaArrangeRoom(q *bizDto.EaArrangeRoomQuery) ([]byte, error) {
	return excel.NewExcel().ExportSync(&EaArrangeRoomExportStruct{q})
}

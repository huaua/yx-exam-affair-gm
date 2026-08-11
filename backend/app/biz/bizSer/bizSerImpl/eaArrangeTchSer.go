package bizSerImpl

import (
	"fmt"
	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/errs"
	"strconv"
	"yx-exam-affair-gm/app/biz/bizDao"
	"yx-exam-affair-gm/app/biz/bizDao/bizDaoImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	"yx-exam-affair-gm/app/cst"
	"yx-exam-affair-gm/app/hint"
	cmutils "yx-exam-affair-gm/app/utils"

	"github.com/sealsee/web-base/public/basemodel"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/query"
	"github.com/sealsee/web-base/public/ds/tx"
)

/*
 * desc: 监考老师分配服务层实现
 * author: init code
 * date: 2024/09/24
 */
type EaArrangeTchService struct{}

var eaArrangeTchService *EaArrangeTchService
var eaArrangeTchDao bizDao.IEaArrangeTchDao

func init() {
	eaArrangeTchService = &EaArrangeTchService{}
	eaArrangeTchDao = bizDaoImpl.NewEaArrangeTchDao()
}

func NewEaArrangeTchService() *EaArrangeTchService {
	return eaArrangeTchService
}

func (arrangeTch *EaArrangeTchService) GetEaArrangeTchById(id int64) *bizDo.EaArrangeTch {
	return eaArrangeTchDao.GetById(id)
}

func (arrangeTch *EaArrangeTchService) ListEaArrangeTch(q *bizDto.EaArrangeTchQuery, p *page.Page) []*bizDto.EaArrangeTchResultDto {
	if q.RoomTaskId == 0 {
		return nil
	}
	sql := `select
                ANY_VALUE(r.room_arrange_id) as room_arrange_id,
                r.room_report_id,
                ANY_VALUE(r.room_id) as room_id,
                ANY_VALUE(r.room_name) as room_name,
                group_concat(
					CONCAT(
						r.dept_name,
						CONCAT('-', r.prof_name),
						IF(r.direction_name IS NOT NULL AND r.direction_name != '', CONCAT('-', r.direction_name), '')
					) SEPARATOR '\n' ) as prof_direction,
                ANY_VALUE(r.task_day) as task_day,
                ANY_VALUE(r.tch_num) as tch_num,
                group_concat( r.arrange_no_str SEPARATOR "/" ) as arrange_no_str,
                ANY_VALUE(r.assignment_state) as assignment_state
            from (
                select a.room_arrange_id, a.room_task_id, a.room_report_id, a.room_id, a.room_name, 
                        a.arrange_no_str, a.assignment_state, b.dept_name, b.prof_name, 
                        b.direction_name, c.task_day, c.tch_num
                from ea_arrange_room a 
                left join ea_arrange_prof_info b on b.deleted = 1 and b.info_id = a.info_id
                left join ea_room_task_report c on c.deleted = 1 and c.id = a.room_report_id
                where a.deleted = 1 and a.room_task_id = %v%v
            ) r
            where %v
            GROUP BY r.room_report_id
            ORDER BY arrange_no_str`
	hasRoomArrangeId := ""
	if q.RoomArrangeId != 0 {
		hasRoomArrangeId = fmt.Sprintf(" and a.room_arrange_id = %v", q.RoomArrangeId)
	}
	condition, args := (&basemodel.BaseEntityQuery{}).
		AddEq("1", 1).
		AddEq("r.assignment_state", q.AssignmentState).
		AddLikeAll("r.dept_name", q.DeptName).
		AddLikeAll("r.room_name", q.RoomName).
		AddLikeAll("r.prof_direction", q.ProfName).
		AddLikeAll("r.prof_direction", q.DirectionName).
		GetConditionArgs()
	sql = fmt.Sprintf(sql, q.RoomTaskId, hasRoomArrangeId, condition)
	listRoom := query.RawSqlQueryListWithPage[bizDto.EaArrangeTchResultDto](p, sql, args...)

	// 遍历查询监考老师分配列表，组装老师id、名称数据
	for _, room := range listRoom {
		sql2 := `select a.tch_report_id, a.tch_name from ea_arrange_tch a where a.deleted = 1 and a.room_report_id = ?`
		listTch := query.RawSqlQueryList[bizDto.EaArrangeTchDto](sql2, room.RoomReportId)
		for _, tch := range listTch {
			room.ArrangeTchReportId = append(room.ArrangeTchReportId, strconv.Itoa(int(tch.TchReportId)))
			room.ArrangeTchName = append(room.ArrangeTchName, tch.TchName)
		}
	}
	return listRoom
}

func (arrangeTch *EaArrangeTchService) ListEaArrangeTchByProf(q *bizDto.EaArrangeTchQuery, p *page.Page) []*bizDto.EaArrangeProfResultDto {
	if q.RoomTaskId == 0 {
		return nil
	}
	sql := `select
                r.info_id,
                ANY_VALUE(r.dept_name) as dept_name,
                ANY_VALUE(r.prof_name) as prof_name, 
                ANY_VALUE(r.direction_name) as direction_name,
                ANY_VALUE(r.arrange_type) as arrange_type,
                ANY_VALUE(r.group_num) as group_num,
                ANY_VALUE(r.stu_num) as stu_num,
                ANY_VALUE(r.room_name) as room_name,
                r.arrange_no_str,
                ANY_VALUE(r.tch_num) as tch_num,
                ANY_VALUE(r.task_day) as task_day,
                group_concat( r.tch_name SEPARATOR " " ) as arrange_tch_name,
                ANY_VALUE(r.assignment_state) as assignment_state
            from (
                select a.room_report_id, a.room_id, a.room_name, a.arrange_no_str, a.assignment_state, 
                       a.stu_num, b.info_id, b.dept_name, b.prof_name, b.direction_name, b.task_day,
                       b.arrange_type, c.group_num, c.group_capacity, c.capacity,
                       c.tch_num, d.tch_arrange_id, d.tch_task_id, d.tch_report_id, d.tch_id, d.tch_name
                from ea_arrange_room a 
                left join ea_arrange_prof_info b on b.deleted = 1 and b.info_id = a.info_id
                left join ea_room_task_report c on c.deleted = 1 and c.id = a.room_report_id
                left join ea_arrange_tch d on d.deleted = 1 and d.room_report_id = a.room_report_id
                where a.deleted = 1 and a.room_task_id = %v
            ) r
            where %v
            GROUP BY r.info_id, r.arrange_no_str 
            ORDER BY arrange_no_str`

	condition, args := (&basemodel.BaseEntityQuery{}).
		AddEq("1", 1).
		AddEq("r.assignment_state", q.AssignmentState).
		AddLikeAll("r.dept_name", q.DeptName).
		AddLikeAll("r.room_name", q.RoomName).
		AddLikeAll("r.prof_name", q.ProfName).
		AddLikeAll("r.direction_name", q.DirectionName).
		AddLikeAll("r.arrange_tch_name", q.TchName).
		GetConditionArgs()
	sql = fmt.Sprintf(sql, q.RoomTaskId, condition)
	return query.RawSqlQueryListWithPage[bizDto.EaArrangeProfResultDto](p, sql, args...)
}

func (arrangeTch *EaArrangeTchService) ListAllArrangeTch(q *bizDto.EaArrangeTchQuery) []*bizDo.EaArrangeTch {
	listAll := make([]*bizDo.EaArrangeTch, 0)
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE

	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		list := eaArrangeTchDao.List(q, pageInfo)
		if len(list) == 0 {
			break
		}
		listAll = append(listAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}
	return listAll
}

func (arrangeTch *EaArrangeTchService) InsertEaArrangeTch(arrangeTchDO *bizDo.EaArrangeTch) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaArrangeTchDao.Insert(gtx, arrangeTchDO)
	})
	return
}

func (arrangeTch *EaArrangeTchService) UpdateEaArrangeTchById(arrangeTchDO *bizDo.EaArrangeTch) (success bool) {
	uWhere := &bizDo.EaArrangeTch{TchArrangeId: arrangeTchDO.TchArrangeId}
	arrangeTchDO.TchArrangeId = 0
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaArrangeTchDao.Update(gtx, arrangeTchDO, uWhere)
	})
	return
}

func (arrangeTch *EaArrangeTchService) BatchUpdateEaArrangeTch(saveDtoList []*bizDto.EaArrangeTchSaveDto, optId int64) (success bool, errInfo errs.ERROR) {
	success = false

	// 获取教室上报详情
	dbRoomReport := eaRoomTaskReportService.GetEaTaskRoomInfoById(saveDtoList[0].RoomReportId)
	if dbRoomReport == nil {
		return false, hint.BIZ_ROOM_TASK_NOT_EXIST_ERR
	}

	// 查询当前考场任务下所有的上报用户id
	dbTchReportIdList := make([]int64, 0)
	dbTchReportIdMap := make(map[int64]bool)
	tchReportIdMap := make(map[string]bool)
	dbArrangeTchList := arrangeTch.ListAllArrangeTch(&bizDto.EaArrangeTchQuery{RoomTaskId: saveDtoList[0].RoomTaskId})
	for _, arrangeTch := range dbArrangeTchList {
		if !dbTchReportIdMap[arrangeTch.TchReportId] {
			dbTchReportIdList = append(dbTchReportIdList, arrangeTch.TchReportId)
			dbTchReportIdMap[arrangeTch.TchReportId] = true
		}

		tchReportIdMap[strconv.FormatInt(arrangeTch.RoomReportId, 10)+"_"+strconv.FormatInt(arrangeTch.TchReportId, 10)] = true
	}

	roomReportIdList := make([]int64, 0)
	tch2ArrangeList := make([]*bizDo.EaArrangeTch, 0)
	for _, saveDto := range saveDtoList {
		if saveDto.RoomReportId == 0 {
			continue
		}
		roomReportIdList = append(roomReportIdList, saveDto.RoomReportId)
		for _, tchReportIdStr := range saveDto.TchReportIdList {
			tchReportId := cmutils.StrToInt64(tchReportIdStr)
			dbTchReport := eaTchTaskReportService.GetEaTchTaskReportById(tchReportId)
			if dbTchReport == nil {
				errInfo = hint.BIZ_TCH_REPORT_NOT_EXIST_ERR.Format(tchReportIdStr)
				return
			}
			// 如果老师已经上报过，并且不是和之前的不是同一个考场上报记录，则提示不允许重复上报
			if _, ok := tchReportIdMap[strconv.FormatInt(saveDto.RoomReportId, 10)+"_"+tchReportIdStr]; !ok {
				if dbTchReport.State == cst.TCH_REPORT_ASSIGNED {
					errInfo = hint.BIZ_TCH_HAS_ARRANGED_ERR.Format(dbTchReport.TchName)
					return
				}
			}

			if dbTchReport.TaskTime != dbRoomReport.TaskDay {
				errInfo = hint.BIZ_TASK_DAY_NOT_MATCH_ERR.Format(dbTchReport.TchName)
				return
			}
			arrangeTch := &bizDo.EaArrangeTch{
				RoomTaskId:   saveDto.RoomTaskId,
				RoomReportId: saveDto.RoomReportId,
				TchReportId:  tchReportId,
				TchTaskId:    dbTchReport.TchTaskId,
				TchId:        dbTchReport.TchId,
				TchName:      dbTchReport.TchName,
			}
			arrangeTch.Deleted = common.Normal
			arrangeTch.SetCreateBy(optId)
			arrangeTch.SetUpdateBy(optId)
			tch2ArrangeList = append(tch2ArrangeList, arrangeTch)
		}
	}

	success = tx.ExecGTx(func(gtx tx.GTx) {

		// 更新原来已分配的上报状态未已上报
		for _, tchReportId := range dbTchReportIdList {
			updateTchReport := &bizDo.EaTchTaskReport{State: cst.ESCALATED}
			updateTchReport.SetUpdateBy(optId)
			eaTchTaskReportDao.Update(gtx, updateTchReport, &bizDo.EaTchTaskReport{ReportId: tchReportId})
		}

		// 更新教职工上报状态为已分配
		for _, tch2Arrange := range tch2ArrangeList {
			updateTchReport := &bizDo.EaTchTaskReport{State: cst.TCH_REPORT_ASSIGNED}
			updateTchReport.SetUpdateBy(optId)
			eaTchTaskReportDao.Update(gtx, updateTchReport, &bizDo.EaTchTaskReport{ReportId: tch2Arrange.TchReportId})
		}

		for _, roomReportId := range roomReportIdList {

			// 删除所有当前教室编排下的分配记录
			eaArrangeTchDao.DeleteByRoomReportId(gtx, roomReportId, optId)
		}

		for _, saveDto := range saveDtoList {
			// 更新教室分配状态
			updateRoomArrange := &bizDo.EaArrangeRoom{}
			updateRoomArrange.SetUpdateBy(optId)
			if len(saveDto.TchReportIdList) > 0 {
				updateRoomArrange.AssignmentState = cst.ASSIGNED
			} else {
				updateRoomArrange.AssignmentState = cst.ESCALATED
			}
			eaArrangeRoomDao.Update(gtx, updateRoomArrange, &bizDo.EaArrangeRoom{RoomReportId: saveDto.RoomReportId})
		}

		// 批量新增
		eaArrangeTchDao.BatchInsert(gtx, tch2ArrangeList)

		// 回写监考经验：当前批次涉及的tch_id都需要重新计算
		refreshTchIdSet := make(map[int64]struct{})
		for _, tch2Arrange := range tch2ArrangeList {
			refreshTchIdSet[tch2Arrange.TchId] = struct{}{}
		}
		// 原已分配列表中存在的tch_id（用于新分配后该tch仍有分配记录的情况）
		for _, dbTch := range dbArrangeTchList {
			refreshTchIdSet[dbTch.TchId] = struct{}{}
		}
		for tchId := range refreshTchIdSet {
			eaTchInfoDao.RefreshInvigilationExperienceByTchId(gtx, tchId)
		}
	})
	return
}

func (arrangeTch *EaArrangeTchService) DeleteEaArrangeTchById(id int64, deleteBy int64) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaArrangeTchDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

func (arrangeTch *EaArrangeTchService) DelArrangeTchByRoomReportId(roomReportId int64, optId int64) (success bool) {
	// 查询已经分配的教职工上报列表
	reportTchIdList := make([]int64, 0)
	reportTchIdMap := make(map[int64]bool)
	dbArrangeTchList := arrangeTch.ListAllArrangeTch(&bizDto.EaArrangeTchQuery{RoomReportId: roomReportId})
	// 记录本次被取消分配涉及的教职工id，用于重新回算监考经验
	affectedTchIdSet := make(map[int64]struct{})
	for _, arrangeTch := range dbArrangeTchList {
		affectedTchIdSet[arrangeTch.TchId] = struct{}{}
		if _, ok := reportTchIdMap[arrangeTch.TchReportId]; !ok {
			reportTchIdMap[arrangeTch.TchReportId] = true
			reportTchIdList = append(reportTchIdList, arrangeTch.TchReportId)
		}
	}
	success = tx.ExecGTx(func(gtx tx.GTx) {
		eaArrangeTchDao.DeleteByRoomReportId(gtx, roomReportId, optId)
		if len(reportTchIdList) > 0 {
			for _, reportTchId := range reportTchIdList {
				report := &bizDo.EaTchTaskReport{State: cst.ESCALATED}
				report.SetUpdateBy(optId)
				eaTchTaskReportDao.Update(gtx, report, &bizDo.EaTchTaskReport{ReportId: reportTchId})
			}
		}
		roomReport := &bizDo.EaArrangeRoom{AssignmentState: cst.ESCALATED}
		roomReport.SetUpdateBy(optId)
		eaArrangeRoomDao.Update(gtx, roomReport, &bizDo.EaArrangeRoom{RoomReportId: roomReportId})
		// 重新计算受影响教职工的监考经验
		for tchId := range affectedTchIdSet {
			eaTchInfoDao.RefreshInvigilationExperienceByTchId(gtx, tchId)
		}
	})
	return
}

func (arrangeTch *EaArrangeTchService) ExportEaArrangeTch(q *bizDto.EaArrangeTchQuery) ([]byte, error) {
	return exportArrangeData(q)
}

// DelSingleArrangeTch 取消单个监考老师分配
// 1. 删除 ea_arrange_tch 单条记录
// 2. 将该老师上报状态回退为 ESCALATED（待分配）
// 3. 检查该教室是否还有其他已分配老师，无则回退教室状态为 ESCALATED
// 4. 重新计算该老师的监考经验
func (arrangeTch *EaArrangeTchService) DelSingleArrangeTch(roomReportId int64, tchReportId int64, optId int64) (success bool) {
	// 查询待删除的分配记录，获取 tchId 用于经验回算
	q := &bizDto.EaArrangeTchQuery{RoomReportId: roomReportId}
	dbList := arrangeTch.ListAllArrangeTch(q)
	var targetTchId int64
	for _, item := range dbList {
		if item.TchReportId == tchReportId {
			targetTchId = item.TchId
			break
		}
	}
	if targetTchId == 0 {
		return false
	}

	success = tx.ExecGTx(func(gtx tx.GTx) {
		// 1. 删除单条分配记录
		eaArrangeTchDao.DeleteByTchReportId(gtx, roomReportId, tchReportId, optId)

		// 2. 回退老师上报状态为待分配(ESCALATED)
		report := &bizDo.EaTchTaskReport{State: cst.ESCALATED}
		report.SetUpdateBy(optId)
		eaTchTaskReportDao.Update(gtx, report, &bizDo.EaTchTaskReport{ReportId: tchReportId})

		// 3. 检查该教室是否还有剩余已分配老师
		var remainingCnt int64
		gtx.Raw("select count(1) from ea_arrange_tch where deleted=1 and room_report_id=?", roomReportId).Scan(&remainingCnt)
		if remainingCnt < 1 {
			// 无剩余分配，教室状态回退为待分配
			roomUpdate := &bizDo.EaArrangeRoom{AssignmentState: cst.ESCALATED}
			roomUpdate.SetUpdateBy(optId)
			eaArrangeRoomDao.Update(gtx, roomUpdate, &bizDo.EaArrangeRoom{RoomReportId: roomReportId})
		}

		// 4. 重新计算该老师的监考经验
		eaTchInfoDao.RefreshInvigilationExperienceByTchId(gtx, targetTchId)
	})
	return
}

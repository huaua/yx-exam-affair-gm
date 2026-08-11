package bizSerImpl

import (
	"github.com/sealsee/web-base/public/cst/common"
	"strconv"
	"unicode/utf8"
	"yx-exam-affair-gm/app/biz/bizDao"
	"yx-exam-affair-gm/app/biz/bizDao/bizDaoImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	"yx-exam-affair-gm/app/cst"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
	"github.com/sealsee/web-base/public/utils/file/excel"
)

/*
 * desc: 征集任务服务层实现
 * author: 蓝鲸
 * date: 2024/08/14
 */
type EaRoomTaskInfoService struct{}

var eaRoomTaskInfoService *EaRoomTaskInfoService
var eaRoomTaskInfoDao bizDao.IEaRoomTaskInfoDao

func init() {
	eaRoomTaskInfoService = &EaRoomTaskInfoService{}
	eaRoomTaskInfoDao = bizDaoImpl.NewEaRoomTaskInfoDao()
}

func NewEaRoomTaskInfoService() *EaRoomTaskInfoService {
	return eaRoomTaskInfoService
}

func (roomTaskInfo *EaRoomTaskInfoService) GetEaTaskInfoById(id int64) *bizDo.EaRoomTaskInfo {
	return eaRoomTaskInfoDao.GetById(id)
}

func (roomTaskInfo *EaRoomTaskInfoService) ListEaTaskInfo(q *bizDto.EaRoomTaskInfoQuery, p *page.Page) []*bizDo.EaRoomTaskInfo {
	q.AddLikeAll("task_name", q.TaskName)
	return eaRoomTaskInfoDao.List(q, p)
}

func (roomTaskInfo *EaRoomTaskInfoService) ListAllEaTaskInfo(q *bizDto.EaRoomTaskInfoQuery) []*bizDo.EaRoomTaskInfo {
	listDepartAll := make([]*bizDo.EaRoomTaskInfo, 0)
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE
	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		list := eaRoomTaskInfoDao.List(q, pageInfo)
		if len(list) == 0 {
			break
		}
		listDepartAll = append(listDepartAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}
	return listDepartAll
}

func (roomTaskInfo *EaRoomTaskInfoService) InsertEaTaskInfo(taskInfoDO *bizDo.EaRoomTaskInfo) (success bool, errInfo string) {
	if len(taskInfoDO.TaskName) <= 0 || taskInfoDO.TaskDay.IsZero() || taskInfoDO.LevelCode == "" || taskInfoDO.NeedNum <= 0 {
		return false, "任务名称，任务日期，层次，考场数量不能为空"
	}

	success, errInfo = _validate(taskInfoDO, true)
	if !success {
		return false, errInfo
	}

	taskInfoDO.State = common.OK
	taskInfoDO.Deleted = common.Normal
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaRoomTaskInfoDao.Insert(gtx, taskInfoDO)
	})
	return
}

func _validate(taskInfoDO *bizDo.EaRoomTaskInfo, needValidateTaskName bool) (success bool, errInfo string) {
	if len(taskInfoDO.TaskName) > 0 {
		if utf8.RuneCountInString(taskInfoDO.TaskName) < 2 || utf8.RuneCountInString(taskInfoDO.TaskName) > 30 {
			return false, "任务名称长度必须在2-30之间"
		}

		if needValidateTaskName {
			// 检查任务名称是否唯一
			exit := eaRoomTaskInfoDao.Count(&bizDto.EaRoomTaskInfoQuery{TaskName: taskInfoDO.TaskName}) > 0
			if exit {
				return false, "任务名称已经存在"
			}
		}
	}

	if len(taskInfoDO.LevelCode) > 0 {
		// 层级校验
		levelQuery := &userDto.SysDictDataQuery{DictType: cst.BIZ_LEVEL, DictValue: taskInfoDO.LevelCode, State: strconv.Itoa(cst.YES)}
		levelQuery.PageSize = 1
		levelList := sysDicDataService.ListSysDictData(levelQuery, &page.Page{CurPage: 1, PageSize: 1})
		if len(levelList) <= 0 {
			return false, "所属层次不正确"
		}
	}

	return true, ""
}

func (roomTaskInfo *EaRoomTaskInfoService) UpdateEaTaskInfoById(taskInfoDO *bizDo.EaRoomTaskInfo, onlyChangeStatus bool) (success bool, errInfo string) {

	// 判断当前任务是否处于未发布状态且已征集数量为0
	dbTaskInfo := eaRoomTaskInfoDao.GetById(taskInfoDO.TaskId)
	if dbTaskInfo == nil {
		return false, "当前任务不存在"
	}

	if !onlyChangeStatus {
		if dbTaskInfo.State != common.OK || dbTaskInfo.CollectNum > 0 {
			return false, "当前任务处于已发布状态或者已征集数量不为0，不能修改"
		}
		if len(taskInfoDO.TaskName) <= 0 || taskInfoDO.TaskDay.IsZero() || taskInfoDO.LevelCode == "" || taskInfoDO.NeedNum <= 0 {
			return false, "任务名称，任务日期，层次，考场数量不能为空"
		}
		if len(taskInfoDO.TaskName) > 0 && dbTaskInfo.TaskName != taskInfoDO.TaskName {
			success, errInfo = _validate(taskInfoDO, true)
			if !success {
				return false, errInfo
			}
		} else {
			success, errInfo = _validate(taskInfoDO, false)
			if !success {
				return false, errInfo
			}
		}
	}

	uWhere := &bizDo.EaRoomTaskInfo{TaskId: taskInfoDO.TaskId}
	taskInfoDO.TaskId = 0
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaRoomTaskInfoDao.Update(gtx, taskInfoDO, uWhere)
	})
	return
}

func (roomTaskInfo *EaRoomTaskInfoService) DeleteEaTaskInfoById(id int64, deleteBy int64) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaRoomTaskInfoDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

// 导出服务 -- start
type EaRoomTaskInfoExportStruct struct {
	Query *bizDto.EaRoomTaskInfoQuery
}

func (h *EaRoomTaskInfoExportStruct) Title() string {
	return "征集任务导出"
}

func (h *EaRoomTaskInfoExportStruct) HeaderColumn() []string {
	return []string{"任务ID,task_id", "任务名称,task_name", "征集日期,task_day", "所属层次,level_code", "需要考场数量,need_num", "已征集数量,collect_num", "已确认数量,comfirm_num", "发布状态: 1-未发布 2-已发布,state", "删除标记：1-正常的；2-已删除的；,deleted", "备注,remark", "创建者,create_by", "创建时间,create_time", "更新者,update_by", "更新时间,update_time"}
}

func (h *EaRoomTaskInfoExportStruct) Rows(p *page.Page) []map[string]interface{} {
	return eaRoomTaskInfoDao.ListMapWithCondition(h.Query, p, nil, nil)
}

func (h *EaRoomTaskInfoExportStruct) Finish(url string) {
}

// 导出服务 -- end
func (roomTaskInfo *EaRoomTaskInfoService) ExportEaTaskInfo(q *bizDto.EaRoomTaskInfoQuery) ([]byte, error) {
	return excel.NewExcel().ExportSync(&EaRoomTaskInfoExportStruct{q})
}

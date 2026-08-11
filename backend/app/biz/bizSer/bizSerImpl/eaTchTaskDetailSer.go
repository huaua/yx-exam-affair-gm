package bizSerImpl

import (
	"sort"
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
 * desc: 监考老师征集任务明细服务层实现
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchTaskDetailService struct{}

var eaTchTaskDetailService *EaTchTaskDetailService
var eaTchTaskDetailDao bizDao.IEaTchTaskDetailDao

func init() {
	eaTchTaskDetailService = &EaTchTaskDetailService{}
	eaTchTaskDetailDao = bizDaoImpl.NewEaTchTaskDetailDao()
}

func NewEaTchTaskDetailService() *EaTchTaskDetailService {
	return eaTchTaskDetailService
}

func (tchTaskDetail *EaTchTaskDetailService) GetEaTchTaskDetailById(id int64) *bizDo.EaTchTaskDetail {
	return eaTchTaskDetailDao.GetById(id)
}

func (tchTaskDetail *EaTchTaskDetailService) ListEaTchTaskDetail(q *bizDto.EaTchTaskDetailQuery, p *page.Page) []*bizDo.EaTchTaskDetail {
	return eaTchTaskDetailDao.List(q, p)
}

func (tchTaskDetail *EaTchTaskDetailService) ListEaTchTaskDetailByDay(q *bizDto.EaTchTaskDetailQuery) []*map[string]any {
	taskDetailList := eaTchTaskDetailDao.List(q, &page.Page{CurPage: 1, PageSize: page.MAX_PAGE_SIZE})
	taskDetailDayMap := make(map[string][]*bizDo.EaTchTaskDetail)
	// 遍历列表，按日期分组
	for _, taskDetail := range taskDetailList {
		dayStr := taskDetail.TaskTime.FormatString(time.DateTime)
		taskDetailDayList := taskDetailDayMap[dayStr]
		if taskDetailDayList == nil {
			taskDetailDayList = make([]*bizDo.EaTchTaskDetail, 0)
		}
		taskDetailDayList = append(taskDetailDayList, taskDetail)
		taskDetailDayMap[dayStr] = taskDetailDayList
	}
	// 按日期排序
	rankCols := make([]string, 0)
	for key := range taskDetailDayMap {
		rankCols = append(rankCols, key)
	}
	sort.Strings(rankCols)
	// 格式转换为前端方便使用的k-v形式
	taskDetailDayMapList := make([]*map[string]any, 0)
	for _, day := range rankCols {
		taskDetailDayMapList = append(taskDetailDayMapList, &map[string]any{"day": day, "list": taskDetailDayMap[day]})
	}
	return taskDetailDayMapList
}

func (tchTaskDetail *EaTchTaskDetailService) ListAllEaTchTaskDetail(q *bizDto.EaTchTaskDetailQuery) []*bizDo.EaTchTaskDetail {
	listAll := make([]*bizDo.EaTchTaskDetail, 0)
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE
	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		list := eaTchTaskDetailDao.List(q, pageInfo)
		if len(list) == 0 {
			break
		}
		listAll = append(listAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}
	return listAll
}

func (tchTaskDetail *EaTchTaskDetailService) InsertEaTchTaskDetail(tchTaskDetailDO *bizDo.EaTchTaskDetail) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaTchTaskDetailDao.Insert(gtx, tchTaskDetailDO)
	})
	return
}

func (tchTaskDetail *EaTchTaskDetailService) UpdateEaTchTaskDetailById(tchTaskDetailDO *bizDo.EaTchTaskDetail) (success bool) {
	uWhere := &bizDo.EaTchTaskDetail{TaskDetailId: tchTaskDetailDO.TaskDetailId}
	tchTaskDetailDO.TaskDetailId = 0
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaTchTaskDetailDao.Update(gtx, tchTaskDetailDO, uWhere)
	})
	return
}

func (tchTaskDetail *EaTchTaskDetailService) UpdateReportedTchNum(tchTaskDetailId int64, reportedTchNum int) bool {

	// 获取详情
	dbTchTaskDetail := tchTaskDetail.GetEaTchTaskDetailById(tchTaskDetailId)
	if tchTaskDetail == nil {
		return false
	}

	return tx.ExecGTx(func(gtx tx.GTx) {
		eaTchTaskDetailDao.UpdateReportNum(gtx, tchTaskDetailId, reportedTchNum)
		eaTchTaskInfoDao.UpdateReportNum(gtx, dbTchTaskDetail.TchTaskId, reportedTchNum)
	})
}

func (tchTaskDetail *EaTchTaskDetailService) DeleteEaTchTaskDetailById(id int64, deleteBy int64) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaTchTaskDetailDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

// 导出服务 -- start
type EaTchTaskDetailExportStruct struct {
	Query *bizDto.EaTchTaskDetailQuery
}

func (h *EaTchTaskDetailExportStruct) Title() string {
	return "监考老师征集任务明细导出"
}

func (h *EaTchTaskDetailExportStruct) HeaderColumn() []string {
	return []string{"教师征集详情ID,task_detail_id", "教师征集任务ID,tch_task_id", "征集日期,task_time", "学院id,dept_id", "征集人数,need_num", "已征集数量,collect_num", "是否强制征集: 1-强制 2-不强制,force_flag", "删除标记：1-正常的；2-已删除的；,deleted", "备注,remark", "创建者,create_by", "创建时间,create_time", "更新者,update_by", "更新时间,update_time"}
}

func (h *EaTchTaskDetailExportStruct) Rows(p *page.Page) []map[string]interface{} {
	return eaTchTaskDetailDao.ListMapWithCondition(h.Query, p, nil, nil)
}

func (h *EaTchTaskDetailExportStruct) Finish(url string) {
}

// 导出服务 -- end
func (tchTaskDetail *EaTchTaskDetailService) ExportEaTchTaskDetail(q *bizDto.EaTchTaskDetailQuery) ([]byte, error) {
	return excel.NewExcel().ExportSync(&EaTchTaskDetailExportStruct{q})
}

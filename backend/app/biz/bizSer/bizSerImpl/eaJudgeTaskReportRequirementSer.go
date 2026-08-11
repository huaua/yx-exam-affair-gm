package bizSerImpl

import (
	"fmt"
	"strconv"
	"strings"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	cmutils "yx-exam-affair-gm/app/utils"

	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/ds"
	"github.com/sealsee/web-base/public/ds/tx"
	"github.com/sealsee/web-base/public/utils/file/excel"
)

type judgeRequirementImportHandler struct {
	TaskId      int64
	OperId      int64
	HeadersList []string
	Rows        []*bizDo.EaJudgeTaskReportRequirement
	Err         string
}

func (h *judgeRequirementImportHandler) Headers(headers []string) {
	if len(headers) != len(h.HeadersList) {
		h.Err = "导入模板表头错误"
		return
	}
	for i := range headers {
		if strings.TrimSpace(headers[i]) != h.HeadersList[i] {
			h.Err = "导入模板表头错误"
			return
		}
	}
}

func (h *judgeRequirementImportHandler) Row(mp *map[string]string) excel.RunStatus {
	if h.Err != "" {
		return excel.Exit
	}
	m := *mp
	rowNo := strconv.Itoa(cmutils.StrToInt(m["rid"]) + 1)
	code := strings.TrimSpace(m[h.HeadersList[0]])
	name := strings.TrimSpace(m[h.HeadersList[1]])
	dept := strings.TrimSpace(m[h.HeadersList[2]])
	numText := strings.TrimSpace(m[h.HeadersList[3]])
	if code == "" || name == "" || dept == "" {
		h.Err = "第" + rowNo + "行，方向代码、方向名称、上报学院均不能为空"
		return excel.Exit
	}
	if !cmutils.IsPositiveInteger(numText) {
		h.Err = "第" + rowNo + "行，上报专家数必须为大于0的正整数"
		return excel.Exit
	}
	var did int64
	ds.GetGDB().Raw("SELECT dept_id FROM ea_dept_info WHERE dept_name = ? AND deleted = 1 LIMIT 1", dept).Scan(&did)
	if did == 0 {
		h.Err = "第" + rowNo + "行，上报学院【" + dept + "】不存在"
		return excel.Exit
	}
	r := &bizDo.EaJudgeTaskReportRequirement{JudgeTaskId: h.TaskId, DirectionCode: code, DirectionName: name, DeptId: did, DeptName: dept, ExpertNum: cmutils.StrToInt(numText)}
	r.Deleted = common.Normal
	r.SetCreateBy(h.OperId)
	r.SetUpdateBy(h.OperId)
	h.Rows = append(h.Rows, r)
	return excel.Next
}

func (h *judgeRequirementImportHandler) After() {}

func (s *EaJudgeTaskInfoService) ListReportRequirements(taskId int64) []*bizDo.EaJudgeTaskReportRequirement {
	list := make([]*bizDo.EaJudgeTaskReportRequirement, 0)
	ds.GetGDB().Where("judge_task_id = ? AND deleted = 1", taskId).Order("direction_code, dept_name").Find(&list)
	return list
}

func (s *EaJudgeTaskInfoService) ImportReportRequirements(q *bizDto.EaJudgeTaskReportRequirementQuery) (bool, string) {
	task := eaJudgeTaskInfoDao.GetById(q.JudgeTaskId)
	if task == nil || task.TaskMode != 2 {
		return false, "学院上报任务不存在"
	}
	if task.PublishState == 2 {
		return false, "任务已发布，不能继续导入"
	}
	h := &judgeRequirementImportHandler{TaskId: q.JudgeTaskId, OperId: q.OperId, HeadersList: []string{"方向代码", "方向名称", "上报学院", "上报专家数"}}
	for _, url := range q.Files {
		excel.NewExcel().ImportWithUrl(url, h)
		if h.Err != "" {
			return false, h.Err
		}
	}
	if len(h.Rows) == 0 {
		return false, "导入文件中没有有效数据"
	}
	// 去重检查：同一任务下，direction_code + dept_name 组合不能已存在
	existingKeys := make(map[string]bool, len(h.Rows))
	var existingList []struct {
		DirectionCode string
		DeptName      string
	}
	ds.GetGDB().Table("ea_judge_task_report_requirement").
		Select("direction_code, dept_name").
		Where("judge_task_id = ? AND deleted = 1", q.JudgeTaskId).
		Scan(&existingList)
	for _, e := range existingList {
		existingKeys[e.DirectionCode+"|"+e.DeptName] = true
	}
	var dupRows []string
	for i, row := range h.Rows {
		key := row.DirectionCode + "|" + row.DeptName
		if existingKeys[key] {
			dupRows = append(dupRows, strconv.Itoa(i+1))
		}
	}
	if len(dupRows) > 0 {
		return false, "第" + strings.Join(dupRows, ",") + "行数据已存在，不允许重复导入"
	}
	success := tx.ExecGTx(func(gtx tx.GTx) {
		for _, row := range h.Rows {
			gtx.Create(row)
		}
	})
	if !success {
		return false, "导入失败"
	}
	return true, ""
}

func (s *EaJudgeTaskInfoService) UpdateReportRequirement(id int64, expertNum int, optId int64) (bool, string) {
	if expertNum < 1 {
		return false, "上报专家数必须为大于0的正整数"
	}
	var row bizDo.EaJudgeTaskReportRequirement
	if ds.GetGDB().Where("requirement_id = ? AND deleted = 1", id).First(&row).RowsAffected < 1 {
		return false, "上报要求不存在"
	}
	var reportedNum int64
	ds.GetGDB().Raw("SELECT COUNT(1) FROM ea_judge_expert_report WHERE requirement_id=? AND submit_state=2 AND audit_state IN (1,2) AND deleted=1", id).Scan(&reportedNum)
	if reportedNum > int64(expertNum) {
		return false, "上报专家数不能小于当前已上报数量"
	}
	res := ds.GetGDB().Model(&row).Updates(map[string]interface{}{"expert_num": expertNum, "update_by": optId})
	return res.RowsAffected > 0, ""
}

// BatchUpdateReportRequirements 批量修改上报专家数（单次请求更新多条）
func (s *EaJudgeTaskInfoService) BatchUpdateReportRequirements(items []bizDto.BatchUpdateReqItem, optId int64) (bool, string) {
	if len(items) == 0 {
		return false, "没有需要保存的数据"
	}
	// 收集所有ID，批量查询已存在的记录
	ids := make([]int64, len(items))
	for i, item := range items {
		if item.ExpertNum < 1 {
			return false, "上报专家数必须为大于0的正整数"
		}
		ids[i] = item.RequirementId
	}
	var existingRows []bizDo.EaJudgeTaskReportRequirement
	ds.GetGDB().Where("requirement_id IN ? AND deleted = 1", ids).Find(&existingRows)
	rowMap := make(map[int64]*bizDo.EaJudgeTaskReportRequirement, len(existingRows))
	for i := range existingRows {
		rowMap[existingRows[i].RequirementId] = &existingRows[i]
	}
	// 校验：检查不存在 / reportedNum > expertNum
	for _, item := range items {
		row, ok := rowMap[item.RequirementId]
		if !ok {
			return false, "部分上报要求数据不存在，请刷新后重试"
		}
		var reportedNum int64
		ds.GetGDB().Raw("SELECT COUNT(1) FROM ea_judge_expert_report WHERE requirement_id=? AND submit_state=2 AND audit_state IN (1,2) AND deleted=1", row.RequirementId).Scan(&reportedNum)
		if reportedNum > int64(item.ExpertNum) {
			return false, fmt.Sprintf("【%s-%s】的上报专家数不能小于当前已上报数量(%d)", row.DirectionName, row.DeptName, reportedNum)
		}
	}
	// 批量更新
	success := tx.ExecGTx(func(gtx tx.GTx) {
		for _, item := range items {
			gtx.Exec("UPDATE ea_judge_task_report_requirement SET expert_num=?, update_by=?, update_time=NOW() WHERE requirement_id=? AND deleted=1",
				item.ExpertNum, optId, item.RequirementId)
		}
	})
	if !success {
		return false, "保存失败"
	}
	return true, ""
}

func (s *EaJudgeTaskInfoService) DeleteReportRequirement(id int64, optId int64) (bool, string) {
	var row bizDo.EaJudgeTaskReportRequirement
	if ds.GetGDB().Where("requirement_id = ? AND deleted = 1", id).First(&row).RowsAffected < 1 {
		return false, "上报要求不存在"
	}
	if row.ReportedNum > 0 {
		return false, "当前方向已存在专家上报数据，不允许删除"
	}
	res := ds.GetGDB().Exec("UPDATE ea_judge_task_report_requirement SET deleted=UNIX_TIMESTAMP(),update_by=?,update_time=NOW() WHERE requirement_id=? AND deleted=1", optId, id)
	if res.Error != nil {
		return false, res.Error.Error()
	}
	return res.RowsAffected > 0, ""
}

package bizApi

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"yx-exam-affair-gm/app/cst"

	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/ds"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/route"
	"github.com/sealsee/web-base/public/web"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type collegeReportReq struct {
	JudgeTaskId    int64  `json:"judgeTaskId,string"`
	RequirementId  int64  `json:"requirementId,string"`
	ExpertId       int64  `json:"expertId,string"`
	ReportId       int64  `json:"reportId,string"`
	NewExpertId    int64  `json:"newExpertId,string"`
	Name           string `json:"name"`
	ExpertCategory int    `json:"expertCategory,string"`
	Title          string `json:"title"`
	GoodSubjects   string `json:"goodSubjects"`
	DeptId         int64  `json:"deptId,string"`
	DirectionCode  string `json:"directionCode"`
	Sex            string `json:"sex"`
	AuditState     int    `json:"auditState,string"`
}

type collegeReportTask struct {
	JudgeTaskId   int64     `json:"judgeTaskId,string"`
	JudgeTaskName string    `json:"judgeTaskName"`
	TaskName      string    `json:"taskName"`
	TaskYear      int       `json:"taskYear,string"`
	JudgeTaskType int       `json:"judgeTaskType,string"`
	TaskStartTime time.Time `json:"taskStartTime"`
	TaskEndTime   time.Time `json:"taskEndTime"`
}

type collegeReportDirection struct {
	RequirementId int64     `json:"requirementId,string"`
	JudgeTaskId   int64     `json:"judgeTaskId,string"`
	DirectionCode string    `json:"directionCode"`
	DirectionName string    `json:"directionName"`
	ExpertNum     int       `json:"expertNum,string"`
	ReportedNum   int       `json:"reportedNum,string"`
	ApprovedNum   int       `json:"approvedNum,string"`
	RejectedNum   int       `json:"rejectedNum,string"`
	DraftNum      int       `json:"draftNum,string"`
	TaskName      string    `json:"taskName"`
	TaskStartTime time.Time `json:"taskStartTime"`
	TaskEndTime   time.Time `json:"taskEndTime"`
}

type collegeReportExpert struct {
	ReportId       int64     `json:"reportId,string"`
	ExpertId       int64     `json:"expertId,string"`
	Name           string    `json:"name"`
	Sex            string    `json:"sex"`
	PhoneNo        string    `json:"phoneNo"`
	IdCard         string    `json:"idCard"`
	ExpertCategory int       `json:"expertCategory,string"`
	DeptName       string    `json:"deptName"`
	UnitName       string    `json:"unitName"`
	Type           string    `json:"type"`
	Title          string    `json:"title"`
	InitEdu        int       `json:"initEdu,string"`
	InitMajor      string    `json:"initMajor"`
	FinalEdu       int       `json:"finalEdu,string"`
	FinalMajor     string    `json:"finalMajor"`
	GoodSubjects   string    `json:"goodSubjects"`
	BankCard       string    `json:"bankCard"`
	Intro          string    `json:"intro"`
	Evaluation     string    `json:"evaluation"`
	SubmitState    int       `json:"submitState,string"`
	AuditState     int       `json:"auditState,string"`
	RequirementId  int64     `json:"requirementId,string"`
	JudgeTaskId    int64     `json:"judgeTaskId,string"`
	DirectionCode  string    `json:"directionCode"`
	DirectionName  string    `json:"directionName"`
	TaskName       string    `json:"taskName"`
	TaskStartTime  time.Time `json:"taskStartTime"`
	TaskEndTime    time.Time `json:"taskEndTime"`
}

func init() {
	route.AddGroup("/api/auth/ea_judge_expert_report", "评委专家上报").
		AddPOSTHandelPowerR("/tasks", collegeReportTasks, true, cst.SYS_ROLE_KEY_DEPARTMENT, "学院可上报任务").
		AddPOSTHandelPowerR("/directions", collegeReportDirections, true, cst.SYS_ROLE_KEY_DEPARTMENT, "学院上报方向").
		AddPOSTHandelPowerR("/experts", collegeReportExperts, true, cst.SYS_ROLE_KEY_DEPARTMENT, "学院可用专家").
		AddPOSTHandelPowerR("/save", saveCollegeReportDraft, true, cst.SYS_ROLE_KEY_DEPARTMENT, "暂存上报专家").
		AddPOSTHandelPowerR("/remove", removeCollegeReportDraft, true, cst.SYS_ROLE_KEY_DEPARTMENT, "清除未提交专家").
		AddPOSTHandelPowerR("/submit", submitCollegeReport, true, cst.SYS_ROLE_KEY_DEPARTMENT, "提交上报专家").
		AddPOSTHandelPowerR("/reported", collegeReportedExperts, true, cst.SYS_ROLE_KEY_DEPARTMENT, "查看已上报专家").
		AddPOSTHandelPowerR("/replace", replaceCollegeReportedExpert, true, cst.SYS_ROLE_KEY_DEPARTMENT, "替换待审核专家").
		AddPOSTHandelPowerR("/export", exportCollegeReportedExperts, true, cst.SYS_ROLE_KEY_DEPARTMENT, "导出上报专家").
		AddPOSTHandelPowerR("/audit/tasks", admissionsReportTasks, true, cst.SYS_ROLE_KEY_ADMISSIONS, "上报审核任务").
		AddPOSTHandelPowerR("/audit/directions", admissionsReportDirections, true, cst.SYS_ROLE_KEY_ADMISSIONS, "上报审核评分方向").
		AddPOSTHandelPowerR("/audit/list", admissionsReportList, true, cst.SYS_ROLE_KEY_ADMISSIONS, "上报审核列表").
		AddPOSTHandelPowerR("/audit/export", exportAdmissionsReportedExperts, true, cst.SYS_ROLE_KEY_ADMISSIONS, "导出上报审核").
		AddPOSTHandelPowerR("/audit/operate", admissionsReportOperate, true, cst.SYS_ROLE_KEY_ADMISSIONS, "审核上报专家").
		AddPOSTHandelPowerR("/audit/reset", admissionsReportReset, true, cst.SYS_ROLE_KEY_ADMISSIONS, "重置上报审核")
}

// admissionsReportTasks 返回全部学院上报模式任务，不受发布状态限制。
func admissionsReportTasks(c *gin.Context) {
	var list []collegeReportTask
	ds.GetGDB().Raw(`SELECT t.judge_task_id,t.judge_task_name,
		CONCAT(t.task_year,'年',t.judge_task_name) task_name,t.task_year,t.judge_task_type,t.task_start_time,t.task_end_time
		FROM ea_judge_task_info t WHERE t.deleted=1 AND t.task_mode=2
		ORDER BY t.task_year DESC,t.judge_task_id DESC`).Scan(&list)
	web.NewJsonResult(c).SetList(list).Render()
}

// admissionsReportDirections 返回所选任务已导入的全部评分方向。
func admissionsReportDirections(c *gin.Context) {
	q, ok := bindCollegeReportReq(c)
	if !ok {
		return
	}
	var list []collegeReportDirection
	if q.JudgeTaskId > 0 {
		ds.GetGDB().Raw(`SELECT requirement_id,judge_task_id,direction_code,direction_name,expert_num,reported_num
			FROM ea_judge_task_report_requirement WHERE judge_task_id=? AND deleted=1 ORDER BY direction_code,requirement_id`, q.JudgeTaskId).Scan(&list)
	}
	web.NewJsonResult(c).SetList(list).Render()
}

func admissionsReportList(c *gin.Context) {
	q, ok := bindCollegeReportReq(c)
	if !ok {
		return
	}
	if q.JudgeTaskId <= 0 {
		web.NewJsonResult(c).SetList([]collegeReportExpert{}).Render()
		return
	}
	args := []interface{}{q.JudgeTaskId}
	sql := `SELECT p.report_id,p.requirement_id,p.judge_task_id,e.expert_id,e.name,e.sex,e.phone_no,e.id_card,e.expert_category,
		COALESCE(d.dept_name,'') dept_name,e.unit_name,e.type,e.title,e.init_edu,e.init_major,e.final_edu,e.final_major,
		e.good_subjects,e.bank_card,e.intro,e.evaluation,p.submit_state,p.audit_state,r.direction_code,r.direction_name,
		CONCAT(t.task_year,'年',t.judge_task_name) task_name,t.task_start_time,t.task_end_time
		FROM ea_judge_expert_report p
		JOIN ea_judge_experts e ON e.expert_id=p.expert_id
		JOIN ea_judge_task_report_requirement r ON r.requirement_id=p.requirement_id AND r.deleted=1
		JOIN ea_judge_task_info t ON t.judge_task_id=p.judge_task_id AND t.deleted=1
		LEFT JOIN ea_dept_info d ON d.dept_id=p.dept_id
		WHERE p.judge_task_id=? AND p.submit_state=2 AND p.deleted=1`
	if q.DeptId > 0 {
		sql += " AND p.dept_id=?"
		args = append(args, q.DeptId)
	}
	if strings.TrimSpace(q.DirectionCode) != "" {
		sql += " AND r.direction_code=?"
		args = append(args, strings.TrimSpace(q.DirectionCode))
	}
	if strings.TrimSpace(q.Name) != "" {
		sql += " AND e.name LIKE ?"
		args = append(args, "%"+strings.TrimSpace(q.Name)+"%")
	}
	if strings.TrimSpace(q.Sex) != "" {
		sql += " AND e.sex=?"
		args = append(args, strings.TrimSpace(q.Sex))
	}
	if strings.TrimSpace(q.Title) != "" {
		sql += " AND e.title=?"
		args = append(args, strings.TrimSpace(q.Title))
	}
	if q.AuditState > 0 {
		sql += " AND p.audit_state=?"
		args = append(args, q.AuditState)
	}
	sql += " ORDER BY r.direction_code,p.report_id"
	var list []collegeReportExpert
	ds.GetGDB().Raw(sql, args...).Scan(&list)
	web.NewJsonResult(c).SetList(list).Render()
}

func exportAdmissionsReportedExperts(c *gin.Context) {
	q, ok := bindCollegeReportReq(c)
	if !ok {
		return
	}
	if q.JudgeTaskId <= 0 {
		web.NewJsonResult(c).SetErr("请选择任务").Render()
		return
	}
	args := []interface{}{q.JudgeTaskId}
	sql := `SELECT CONCAT(t.task_year,'年',t.judge_task_name) task_name,
		CONCAT(r.direction_code,r.direction_name) direction,
		CONCAT(DATE(t.task_start_time),'至',DATE(t.task_end_time)) score_date,
		e.name,e.sex,e.phone_no,e.id_card,COALESCE(d.dept_name,'') dept_name,e.expert_category,
		e.init_edu,e.init_major,e.final_edu,e.final_major,e.good_subjects,e.bank_card,e.intro,e.evaluation,p.audit_state
		FROM ea_judge_expert_report p
		JOIN ea_judge_experts e ON e.expert_id=p.expert_id
		JOIN ea_judge_task_report_requirement r ON r.requirement_id=p.requirement_id AND r.deleted=1
		JOIN ea_judge_task_info t ON t.judge_task_id=p.judge_task_id AND t.deleted=1
		LEFT JOIN ea_dept_info d ON d.dept_id=p.dept_id
		WHERE p.judge_task_id=? AND p.submit_state=2 AND p.deleted=1`
	if q.DeptId > 0 {
		sql += " AND p.dept_id=?"
		args = append(args, q.DeptId)
	}
	if strings.TrimSpace(q.DirectionCode) != "" {
		sql += " AND r.direction_code=?"
		args = append(args, strings.TrimSpace(q.DirectionCode))
	}
	if strings.TrimSpace(q.Name) != "" {
		sql += " AND e.name LIKE ?"
		args = append(args, "%"+strings.TrimSpace(q.Name)+"%")
	}
	if strings.TrimSpace(q.Sex) != "" {
		sql += " AND e.sex=?"
		args = append(args, strings.TrimSpace(q.Sex))
	}
	if strings.TrimSpace(q.Title) != "" {
		sql += " AND e.title=?"
		args = append(args, strings.TrimSpace(q.Title))
	}
	if q.AuditState > 0 {
		sql += " AND p.audit_state=?"
		args = append(args, q.AuditState)
	}
	sql += " ORDER BY r.direction_code,p.report_id"
	var rows []map[string]interface{}
	ds.GetGDB().Raw(sql, args...).Scan(&rows)
	// 数据库中 intro 字段为 text，map 转出来是 []byte，excelize 不识别而写入空字符串。
	// 统一将 []byte 转换为 string，保证专家介绍等长文本列能正常写出。
	for _, row := range rows {
		for k, v := range row {
			if b, ok := v.([]byte); ok {
				row[k] = string(b)
			}
		}
	}
	f := excelize.NewFile()
	sheet := "评委上报审核"
	f.SetSheetName("Sheet1", sheet)
	headers := []string{"任务名称", "评分方向", "评分日期", "姓名", "性别", "手机号", "身份证号", "所属学院", "专家类别", "初始学历", "所学专业", "最终学历", "所学专业", "擅长科目", "银行卡号", "专家介绍", "专家评价", "审核状态"}
	keys := []string{"task_name", "direction", "score_date", "name", "sex", "phone_no", "id_card", "dept_name", "expert_category", "init_edu", "init_major", "final_edu", "final_major", "good_subjects", "bank_card", "intro", "evaluation", "audit_state"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	for ri, row := range rows {
		for ci, key := range keys {
			value := row[key]
			switch key {
			case "audit_state":
				switch fmt.Sprint(value) {
				case "1":
					value = "已上报，待审核"
				case "2":
					value = "审核通过"
				case "3":
					value = "审核不通过"
				}
			case "expert_category":
				switch fmt.Sprint(value) {
				case "1":
					value = "校内"
				case "2":
					value = "校外"
				default:
					value = ""
				}
			case "init_edu", "final_edu":
				if n, e := strconv.Atoi(fmt.Sprint(value)); e == nil {
					value = cst.GetEduTypeLabelByValue(n)
				}
			}
			cell, _ := excelize.CoordinatesToCellName(ci+1, ri+2)
			f.SetCellValue(sheet, cell, value)
		}
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		web.NewJsonResult(c).SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, buf.Bytes())
}

func updateReportAuditState(c *gin.Context, target int, requireReviewed bool) {
	q, ok := bindCollegeReportReq(c)
	if !ok {
		return
	}
	if q.ReportId <= 0 || (target != 1 && target != 2 && target != 3) {
		web.NewJsonResult(c).SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	var row struct {
		RequirementId int64
		AuditState    int
	}
	if ds.GetGDB().Raw("SELECT requirement_id,audit_state FROM ea_judge_expert_report WHERE report_id=? AND submit_state=2 AND deleted=1", q.ReportId).Scan(&row).RowsAffected == 0 {
		web.NewJsonResult(c).SetErr("上报记录不存在").Render()
		return
	}
	if requireReviewed && row.AuditState == 1 {
		web.NewJsonResult(c).SetErr("待审核数据无需重置").Render()
		return
	}
	if !requireReviewed && row.AuditState != 1 {
		web.NewJsonResult(c).SetErr("该数据已审核，请先重置后再操作").Render()
		return
	}
	uid := context.NewUserContext(c).GetUserId()
	err := ds.GetGDB().Transaction(func(tx *gorm.DB) error {
		if e := tx.Exec("UPDATE ea_judge_expert_report SET audit_state=?,update_by=?,update_time=NOW() WHERE report_id=? AND submit_state=2 AND deleted=1", target, uid, q.ReportId).Error; e != nil {
			return e
		}
		return tx.Exec(`UPDATE ea_judge_task_report_requirement SET reported_num=(SELECT COUNT(1) FROM ea_judge_expert_report
			WHERE requirement_id=? AND submit_state=2 AND deleted=1),update_by=?,update_time=NOW() WHERE requirement_id=?`, row.RequirementId, uid, row.RequirementId).Error
	})
	if err != nil {
		web.NewJsonResult(c).SetErr("操作失败").Render()
		return
	}
	web.NewJsonResult(c).Render()
}

func admissionsReportOperate(c *gin.Context) {
	q := new(collegeReportReq)
	if c.ShouldBindJSON(q) != nil || q.ReportId <= 0 || (q.AuditState != 2 && q.AuditState != 3) {
		web.NewJsonResult(c).SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	// 复用统一状态更新逻辑需要重新绑定请求，因此直接执行同等校验与更新。
	var row struct {
		RequirementId int64
		AuditState    int
	}
	if ds.GetGDB().Raw("SELECT requirement_id,audit_state FROM ea_judge_expert_report WHERE report_id=? AND submit_state=2 AND deleted=1", q.ReportId).Scan(&row).RowsAffected == 0 {
		web.NewJsonResult(c).SetErr("上报记录不存在").Render()
		return
	}
	if row.AuditState != 1 {
		web.NewJsonResult(c).SetErr("该数据已审核，请先重置后再操作").Render()
		return
	}
	uid := context.NewUserContext(c).GetUserId()
	err := ds.GetGDB().Transaction(func(tx *gorm.DB) error {
		if e := tx.Exec("UPDATE ea_judge_expert_report SET audit_state=?,update_by=?,update_time=NOW() WHERE report_id=?", q.AuditState, uid, q.ReportId).Error; e != nil {
			return e
		}
		return tx.Exec(`UPDATE ea_judge_task_report_requirement SET reported_num=(SELECT COUNT(1) FROM ea_judge_expert_report WHERE requirement_id=? AND submit_state=2 AND deleted=1),update_by=?,update_time=NOW() WHERE requirement_id=?`, row.RequirementId, uid, row.RequirementId).Error
	})
	if err != nil {
		web.NewJsonResult(c).SetErr("审核失败").Render()
		return
	}
	web.NewJsonResult(c).Render()
}

func admissionsReportReset(c *gin.Context) { updateReportAuditState(c, 1, true) }

func bindCollegeReportReq(c *gin.Context) (*collegeReportReq, bool) {
	q := new(collegeReportReq)
	if c.ShouldBindJSON(q) != nil {
		web.NewJsonResult(c).SetErrs(errs.PARAM_INVALID_ERR).Render()
		return q, false
	}
	return q, true
}

func collegeDept(c *gin.Context) (int64, int64) {
	ctx := context.NewUserContext(c)
	return ctx.GetUser().DeptId, ctx.GetUserId()
}

func collegeReportTasks(c *gin.Context) {
	deptId, _ := collegeDept(c)
	var list []collegeReportTask
	ds.GetGDB().Raw(`SELECT DISTINCT t.judge_task_id,t.judge_task_name,
		CONCAT(t.task_year,'年',t.judge_task_name) task_name,t.task_year,t.judge_task_type,t.task_start_time,t.task_end_time
		FROM ea_judge_task_info t JOIN ea_judge_task_report_requirement r ON r.judge_task_id=t.judge_task_id AND r.deleted=1
		WHERE t.deleted=1 AND t.task_mode=2 AND t.publish_state=2 AND r.dept_id=?
		ORDER BY t.task_year DESC,t.judge_task_id DESC`, deptId).Scan(&list)
	web.NewJsonResult(c).SetList(list).Render()
}

func collegeReportDirections(c *gin.Context) {
	q, ok := bindCollegeReportReq(c)
	if !ok {
		return
	}
	deptId, _ := collegeDept(c)
	var list []collegeReportDirection
	sql := `SELECT r.requirement_id,r.judge_task_id,r.direction_code,r.direction_name,r.expert_num,
		SUM(CASE WHEN p.deleted=1 AND p.submit_state=2 THEN 1 ELSE 0 END) reported_num,
		SUM(CASE WHEN p.deleted=1 AND p.submit_state=2 AND p.audit_state=2 THEN 1 ELSE 0 END) approved_num,
		SUM(CASE WHEN p.deleted=1 AND p.submit_state=2 AND p.audit_state=3 THEN 1 ELSE 0 END) rejected_num,
		SUM(CASE WHEN p.deleted=1 AND p.submit_state=1 THEN 1 ELSE 0 END) draft_num,
		CONCAT(t.task_year,'年',t.judge_task_name) task_name,t.task_start_time,t.task_end_time
		FROM ea_judge_task_report_requirement r
		JOIN ea_judge_task_info t ON t.judge_task_id=r.judge_task_id AND t.deleted=1 AND t.task_mode=2 AND t.publish_state=2
		LEFT JOIN ea_judge_expert_report p ON p.requirement_id=r.requirement_id AND p.dept_id=r.dept_id AND p.deleted=1
		WHERE r.deleted=1 AND r.dept_id=?`
	args := []interface{}{deptId}
	if q.JudgeTaskId > 0 {
		sql += " AND r.judge_task_id=?"
		args = append(args, q.JudgeTaskId)
	}
	sql += " GROUP BY r.requirement_id, t.judge_task_id, t.task_year, t.judge_task_name, t.task_start_time, t.task_end_time ORDER BY t.task_year DESC, t.judge_task_id DESC, r.direction_code"
	ds.GetGDB().Raw(sql, args...).Scan(&list)
	web.NewJsonResult(c).SetList(list).Render()
}

func validCollegeRequirement(requirementId, deptId int64) (collegeReportDirection, int, error) {
	var row collegeReportDirection
	var taskType int
	r := ds.GetGDB().Raw(`SELECT r.requirement_id,r.judge_task_id,r.direction_code,r.direction_name,r.expert_num,t.judge_task_type
		FROM ea_judge_task_report_requirement r JOIN ea_judge_task_info t ON t.judge_task_id=r.judge_task_id
		WHERE r.requirement_id=? AND r.dept_id=? AND r.deleted=1 AND t.deleted=1 AND t.task_mode=2 AND t.publish_state=2`, requirementId, deptId).
		Row()
	if err := r.Scan(&row.RequirementId, &row.JudgeTaskId, &row.DirectionCode, &row.DirectionName, &row.ExpertNum, &taskType); err != nil {
		return row, 0, fmt.Errorf("上报要求不存在或任务尚未发布")
	}
	return row, taskType, nil
}

func collegeReportExperts(c *gin.Context) {
	q, ok := bindCollegeReportReq(c)
	if !ok {
		return
	}
	deptId, _ := collegeDept(c)
	req, taskType, err := validCollegeRequirement(q.RequirementId, deptId)
	if err != nil {
		web.NewJsonResult(c).SetErr(err.Error()).Render()
		return
	}
	args := []interface{}{q.RequirementId, deptId, taskType, q.RequirementId}
	sql := `SELECT COALESCE(p.report_id,0) report_id,e.expert_id,e.name,e.sex,e.phone_no,e.id_card,e.expert_category,
		COALESCE(d.dept_name,'') dept_name,e.unit_name,e.type,e.title,e.init_edu,e.init_major,e.final_edu,e.final_major,
		e.good_subjects,e.bank_card,e.intro,e.evaluation,COALESCE(p.submit_state,0) submit_state,COALESCE(p.audit_state,0) audit_state
		FROM ea_judge_experts e LEFT JOIN ea_dept_info d ON d.dept_id=e.dept_id
		LEFT JOIN ea_judge_expert_report p ON p.requirement_id=? AND p.dept_id=? AND p.expert_id=e.expert_id AND p.deleted=1
		WHERE e.deleted=1 AND e.dept_id=? AND e.audit_status=3 AND e.allow_draw=1 AND FIND_IN_SET(?,e.type)>0
		AND NOT EXISTS(SELECT 1 FROM ea_judge_expert_ban b WHERE b.expert_id=e.expert_id AND b.deleted=1)
		AND NOT EXISTS(SELECT 1 FROM ea_judge_expert_report used WHERE used.requirement_id=? AND used.expert_id=e.expert_id AND used.deleted=1 AND used.submit_state=2 AND used.audit_state IN (1,2))`
	// 专家所属学院必须是当前学院；同一专家可以在同一任务的不同研究方向重复上报。
	args = []interface{}{q.RequirementId, deptId, deptId, taskType, q.RequirementId}
	if strings.TrimSpace(q.Name) != "" {
		sql += " AND e.name LIKE ?"
		args = append(args, "%"+strings.TrimSpace(q.Name)+"%")
	}
	if q.ExpertCategory > 0 {
		sql += " AND e.expert_category=?"
		args = append(args, q.ExpertCategory)
	}
	if strings.TrimSpace(q.Title) != "" {
		sql += " AND e.title LIKE ?"
		args = append(args, "%"+strings.TrimSpace(q.Title)+"%")
	}
	if strings.TrimSpace(q.GoodSubjects) != "" {
		sql += " AND e.good_subjects LIKE ?"
		args = append(args, "%"+strings.TrimSpace(q.GoodSubjects)+"%")
	}
	sql += " ORDER BY p.submit_state DESC,e.expert_id"
	var list []collegeReportExpert
	ds.GetGDB().Raw(sql, args...).Scan(&list)
	_ = req
	web.NewJsonResult(c).SetList(list).Render()
}

func saveCollegeReportDraft(c *gin.Context) {
	q, ok := bindCollegeReportReq(c)
	if !ok {
		return
	}
	deptId, userId := collegeDept(c)
	req, taskType, err := validCollegeRequirement(q.RequirementId, deptId)
	if err != nil {
		web.NewJsonResult(c).SetErr(err.Error()).Render()
		return
	}
	var expertCount int64
	ds.GetGDB().Raw("SELECT COUNT(1) FROM ea_judge_experts WHERE expert_id=? AND dept_id=? AND deleted=1 AND audit_status=3 AND allow_draw=1 AND FIND_IN_SET(?,type)>0", q.ExpertId, deptId, taskType).Scan(&expertCount)
	if expertCount == 0 {
		web.NewJsonResult(c).SetErr("该专家不满足本任务上报要求").Render()
		return
	}
	var activeCount int64
	ds.GetGDB().Raw("SELECT COUNT(1) FROM ea_judge_expert_report WHERE requirement_id=? AND dept_id=? AND deleted=1 AND (submit_state=1 OR audit_state IN (1,2))", q.RequirementId, deptId).Scan(&activeCount)
	if activeCount >= int64(req.ExpertNum) {
		web.NewJsonResult(c).SetErr("上报专家数不能超过要求数量").Render()
		return
	}
	// 保留历史：deleted 字段存 Unix 时间戳（bigint），未删除=1，删除=UNIX_TIMESTAMP()。
	// 同一专家多次清除会生成不同时间戳的已删记录，普通唯一索引 (requirement_id,expert_id,deleted) 天然不冲突。
	res := ds.GetGDB().Exec(`INSERT INTO ea_judge_expert_report(judge_task_id,requirement_id,dept_id,expert_id,submit_state,audit_state,deleted,create_by,create_time,update_by,update_time)
		VALUES(?,?,?,?,1,1,1,?,NOW(),?,NOW()) ON DUPLICATE KEY UPDATE submit_state=IF(audit_state=3,1,submit_state),audit_state=IF(audit_state=3,1,audit_state),deleted=1,update_by=VALUES(update_by),update_time=NOW()`, req.JudgeTaskId, req.RequirementId, deptId, q.ExpertId, userId, userId)
	if res.Error != nil {
		web.NewJsonResult(c).SetErr("专家已选择或暂存失败").Render()
		return
	}
	web.NewJsonResult(c).Render()
}

func removeCollegeReportDraft(c *gin.Context) {
	q, ok := bindCollegeReportReq(c)
	if !ok {
		return
	}
	deptId, userId := collegeDept(c)
	if q.ReportId <= 0 {
		web.NewJsonResult(c).SetErr("缺少上报记录编号").Render()
		return
	}
	var row struct {
		SubmitState   int
		Deleted       int64
		DeptId        int64
		RequirementId int64
		ExpertId      int64
	}
	r := ds.GetGDB().Raw("SELECT submit_state,deleted,dept_id,requirement_id,expert_id FROM ea_judge_expert_report WHERE report_id=?", q.ReportId).Row()
	if err := r.Scan(&row.SubmitState, &row.Deleted, &row.DeptId, &row.RequirementId, &row.ExpertId); err != nil {
		web.NewJsonResult(c).SetErr("清除记录不存在").Render()
		return
	}
	if row.Deleted != 1 {
		web.NewJsonResult(c).SetErr("该记录已删除，请刷新页面").Render()
		return
	}
	if row.DeptId != deptId {
		web.NewJsonResult(c).SetErr("只能清除本学院数据").Render()
		return
	}
	if row.SubmitState != 1 {
		web.NewJsonResult(c).SetErr("只能清除未提交的专家").Render()
		return
	}
	// 保留历史：deleted 字段存 Unix 时间戳，每次删除时间戳不同，普通唯一索引天然不冲突。
	res := ds.GetGDB().Exec("UPDATE ea_judge_expert_report SET deleted=UNIX_TIMESTAMP(),update_by=?,update_time=NOW() WHERE report_id=? AND dept_id=? AND submit_state=1 AND deleted=1", userId, q.ReportId, deptId)
	if res.Error != nil {
		web.NewJsonResult(c).SetErr("清除失败：" + res.Error.Error()).Render()
		return
	}
	if res.RowsAffected == 0 {
		web.NewJsonResult(c).SetErr("清除失败，请刷新后重试").Render()
		return
	}
	web.NewJsonResult(c).Render()
}

func submitCollegeReport(c *gin.Context) {
	q, ok := bindCollegeReportReq(c)
	if !ok {
		return
	}
	deptId, userId := collegeDept(c)
	req, _, err := validCollegeRequirement(q.RequirementId, deptId)
	if err != nil {
		web.NewJsonResult(c).SetErr(err.Error()).Render()
		return
	}
	var count int64
	ds.GetGDB().Raw("SELECT COUNT(1) FROM ea_judge_expert_report WHERE requirement_id=? AND dept_id=? AND deleted=1 AND (submit_state=1 OR (submit_state=2 AND audit_state IN (1,2)))", q.RequirementId, deptId).Scan(&count)
	if count != int64(req.ExpertNum) {
		web.NewJsonResult(c).SetErr(fmt.Sprintf("当前有效上报数为%d，必须与要求上报数%d一致", count, req.ExpertNum)).Render()
		return
	}
	res := ds.GetGDB().Exec("UPDATE ea_judge_expert_report SET submit_state=2,audit_state=1,update_by=?,update_time=NOW() WHERE requirement_id=? AND dept_id=? AND submit_state=1 AND deleted=1", userId, q.RequirementId, deptId)
	if res.Error != nil {
		web.NewJsonResult(c).SetErr("提交失败").Render()
		return
	}
	ds.GetGDB().Exec(`UPDATE ea_judge_task_report_requirement SET reported_num=(SELECT COUNT(1) FROM ea_judge_expert_report
		WHERE requirement_id=? AND submit_state=2 AND deleted=1),update_by=?,update_time=NOW() WHERE requirement_id=?`, q.RequirementId, userId, q.RequirementId)
	web.NewJsonResult(c).Render()
}

func collegeReportedExperts(c *gin.Context) {
	q, ok := bindCollegeReportReq(c)
	if !ok {
		return
	}
	deptId, _ := collegeDept(c)
	var list []collegeReportExpert
	ds.GetGDB().Raw(`SELECT p.report_id,e.expert_id,e.name,e.sex,e.phone_no,e.id_card,e.expert_category,COALESCE(d.dept_name,'') dept_name,
		e.unit_name,e.type,e.title,e.init_edu,e.init_major,e.final_edu,e.final_major,e.good_subjects,e.bank_card,e.intro,e.evaluation,p.submit_state,p.audit_state
		FROM ea_judge_expert_report p JOIN ea_judge_experts e ON e.expert_id=p.expert_id LEFT JOIN ea_dept_info d ON d.dept_id=e.dept_id
		WHERE p.requirement_id=? AND p.dept_id=? AND p.submit_state=2 AND p.deleted=1 ORDER BY p.report_id`, q.RequirementId, deptId).Scan(&list)
	web.NewJsonResult(c).SetList(list).Render()
}

func replaceCollegeReportedExpert(c *gin.Context) {
	q, ok := bindCollegeReportReq(c)
	if !ok {
		return
	}
	deptId, userId := collegeDept(c)
	var old struct {
		ReportId, RequirementId, JudgeTaskId int64
		AuditState                           int
	}
	if ds.GetGDB().Raw("SELECT report_id,requirement_id,judge_task_id,audit_state FROM ea_judge_expert_report WHERE report_id=? AND dept_id=? AND submit_state=2 AND deleted=1", q.ReportId, deptId).Scan(&old).RowsAffected == 0 || old.AuditState != 1 {
		web.NewJsonResult(c).SetErr("只能替换已上报且待审核的专家").Render()
		return
	}
	_, taskType, err := validCollegeRequirement(old.RequirementId, deptId)
	if err != nil {
		web.NewJsonResult(c).SetErr(err.Error()).Render()
		return
	}
	var n int64
	ds.GetGDB().Raw("SELECT COUNT(1) FROM ea_judge_experts WHERE expert_id=? AND dept_id=? AND deleted=1 AND audit_status=3 AND allow_draw=1 AND FIND_IN_SET(?,type)>0", q.NewExpertId, deptId, taskType).Scan(&n)
	if n == 0 {
		web.NewJsonResult(c).SetErr("替换专家不满足任务要求").Render()
		return
	}
	ds.GetGDB().Raw("SELECT COUNT(1) FROM ea_judge_expert_report WHERE requirement_id=? AND expert_id=? AND dept_id=? AND deleted=1", old.RequirementId, q.NewExpertId, deptId).Scan(&n)
	if n > 0 {
		web.NewJsonResult(c).SetErr("所选专家已被选择或已上报，不能替换").Render()
		return
	}
	err = ds.GetGDB().Transaction(func(tx *gorm.DB) error {
		if e := tx.Exec("UPDATE ea_judge_expert_report SET deleted=UNIX_TIMESTAMP(),update_by=?,update_time=NOW() WHERE report_id=?", userId, old.ReportId).Error; e != nil {
			return e
		}
		return tx.Exec(`INSERT INTO ea_judge_expert_report(judge_task_id,requirement_id,dept_id,expert_id,submit_state,audit_state,replaced_from,deleted,create_by,create_time,update_by,update_time)
			VALUES(?,?,?,?,2,1,?,1,?,NOW(),?,NOW()) ON DUPLICATE KEY UPDATE submit_state=2,audit_state=1,replaced_from=VALUES(replaced_from),deleted=1,update_by=VALUES(update_by),update_time=NOW()`, old.JudgeTaskId, old.RequirementId, deptId, q.NewExpertId, old.ReportId, userId, userId).Error
	})
	if err != nil {
		web.NewJsonResult(c).SetErr("替换失败，所选专家可能已上报").Render()
		return
	}
	web.NewJsonResult(c).Render()
}

func exportCollegeReportedExperts(c *gin.Context) {
	q, ok := bindCollegeReportReq(c)
	if !ok {
		return
	}
	deptId, _ := collegeDept(c)
	var rows []map[string]interface{}
	ds.GetGDB().Raw(`SELECT CONCAT(t.task_year,'年',t.judge_task_name) task_name,r.direction_name,
		CONCAT(DATE(t.task_start_time),'至',DATE(t.task_end_time)) score_date,e.name,e.sex,e.phone_no,e.id_card,e.type,
		e.init_edu,e.init_major,e.final_edu,e.final_major,e.good_subjects,e.bank_card,e.intro,e.evaluation,p.audit_state
		FROM ea_judge_expert_report p JOIN ea_judge_task_info t ON t.judge_task_id=p.judge_task_id
		JOIN ea_judge_task_report_requirement r ON r.requirement_id=p.requirement_id JOIN ea_judge_experts e ON e.expert_id=p.expert_id
		WHERE p.judge_task_id=? AND p.dept_id=? AND p.submit_state=2 AND p.deleted=1 ORDER BY r.direction_code,p.report_id`, q.JudgeTaskId, deptId).Scan(&rows)
	// 数据库中 intro 字段为 text，map 转出来是 []byte，excelize 不识别而写入空字符串。
	// 统一将 []byte 转换为 string，保证专家介绍等长文本列能正常写出。
	for _, row := range rows {
		for k, v := range row {
			if b, ok := v.([]byte); ok {
				row[k] = string(b)
			}
		}
	}
	f := excelize.NewFile()
	sheet := "专家上报"
	f.SetSheetName("Sheet1", sheet)
	headers := []string{"任务名称", "评分方向", "评分日期", "姓名", "性别", "手机号", "证件号", "专家类型", "初始学历", "所学专业", "最终学历", "所学专业", "擅长科目", "银行卡号", "专家介绍", "专家评价", "审核状态"}
	keys := []string{"task_name", "direction_name", "score_date", "name", "sex", "phone_no", "id_card", "type", "init_edu", "init_major", "final_edu", "final_major", "good_subjects", "bank_card", "intro", "evaluation", "audit_state"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	for ri, row := range rows {
		for ci, key := range keys {
			value := row[key]
			if key == "audit_state" {
				switch fmt.Sprint(value) {
				case "1":
					value = "已上报，待审核"
				case "2":
					value = "审核通过"
				case "3":
					value = "审核不通过"
				}
			}
			if key == "init_edu" || key == "final_edu" {
				if n, e := strconv.Atoi(fmt.Sprint(value)); e == nil {
					value = cst.GetEduTypeLabelByValue(n)
				}
			}
			cell, _ := excelize.CoordinatesToCellName(ci+1, ri+2)
			f.SetCellValue(sheet, cell, value)
		}
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		web.NewJsonResult(c).SetErr(err.Error()).Render()
		return
	}
	web.DataPackageExcel(c, buf.Bytes())
}

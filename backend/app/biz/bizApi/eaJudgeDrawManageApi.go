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
)

type drawManageTask struct {
	JudgeTaskId     int64     `json:"judgeTaskId,string"`
	TaskName        string    `json:"taskName"`
	JudgeTaskName   string    `json:"judgeTaskName"`
	TaskYear        int       `json:"taskYear,string"`
	ScoreRounds     int       `json:"scoreRounds,string"`
	SubjectCount    int       `json:"subjectCount,string"`
	ExpertsPerRound int       `json:"expertsPerRound,string"`
	JudgeTaskType   int       `json:"judgeTaskType,string"`
	TaskStartTime   time.Time `json:"taskStartTime"`
	TaskEndTime     time.Time `json:"taskEndTime"`
}
type drawManageRow struct {
	RulesId      int64  `json:"rulesId,string"`
	JudgeTaskId  int64  `json:"judgeTaskId,string"`
	DrawTimes    int    `json:"drawTimes,string"`
	SubjectName  string `json:"subjectName"`
	RequiredNum  int    `json:"requiredNum,string"`
	SubmitState  int    `json:"submitState,string"`
	ScoringCount int    `json:"scoringCount,string"`
	BackupCount  int    `json:"backupCount,string"`
}
type drawManageExpert struct {
	ExpertId       int64  `json:"expertId,string"`
	DeptId         int64  `json:"deptId,string"`
	Name           string `json:"name"`
	Sex            string `json:"sex"`
	PhoneNo        string `json:"phoneNo"`
	IdCard         string `json:"idCard"`
	ExpertCategory int    `json:"expertCategory,string"`
	DeptName       string `json:"deptName"`
	UnitName       string `json:"unitName"`
	Type           string `json:"type"`
	Title          string `json:"title"`
	InitEdu        int    `json:"initEdu,string"`
	InitMajor      string `json:"initMajor"`
	FinalEdu       int    `json:"finalEdu,string"`
	FinalMajor     string `json:"finalMajor"`
	GoodSubjects   string `json:"goodSubjects"`
	BankCard       string `json:"bankCard"`
	Intro          string `json:"intro"`
	Evaluation     string `json:"evaluation"`
	ExpertRole     int    `json:"expertRole,string"`
	SubmitState    int    `json:"submitState,string"`
}
type drawManageReq struct {
	JudgeTaskId    int64  `json:"judgeTaskId,string"`
	RulesId        int64  `json:"rulesId,string"`
	SubjectName    string `json:"subjectName"`
	RequiredNum    int    `json:"requiredNum"`
	ExpertId       int64  `json:"expertId,string"`
	ExpertRole     int    `json:"expertRole,string"`
	Name           string `json:"name"`
	ExpertCategory int    `json:"expertCategory,string"`
	SubmitState    int    `json:"submitState,string"`
}

func init() {
	route.AddGroup("/api/auth/ea_judge_draw_manage", "评委抽取管理").
		AddPOSTHandelPowerR("/tasks", drawManageTasks, true, cst.SYS_ROLE_KEY_ADMISSIONS, "可抽取任务").
		AddPOSTHandelPowerR("/rows", drawManageRows, true, cst.SYS_ROLE_KEY_ADMISSIONS, "轮次科目列表").
		AddPOSTHandelPowerR("/experts", drawManageExperts, true, cst.SYS_ROLE_KEY_ADMISSIONS, "可用专家").
		AddPOSTHandelPowerR("/save_selection", saveDrawSelection, true, cst.SYS_ROLE_KEY_ADMISSIONS, "暂存专家").
		AddPOSTHandelPowerR("/remove_selection", removeDrawSelection, true, cst.SYS_ROLE_KEY_ADMISSIONS, "移除专家").
		AddPOSTHandelPowerR("/submit", submitDrawSelection, true, cst.SYS_ROLE_KEY_ADMISSIONS, "提交抽取").
		AddPOSTHandelPowerR("/update_rule", updateDrawRule, true, cst.SYS_ROLE_KEY_ADMISSIONS, "修改科目人数").
		AddPOSTHandelPowerR("/export", exportDrawSelection, true, cst.SYS_ROLE_KEY_ADMISSIONS, "导出专家")
}

func bindDrawReq(c *gin.Context) (*drawManageReq, bool) {
	q := new(drawManageReq)
	if c.ShouldBindJSON(q) != nil {
		web.NewJsonResult(c).SetErrs(errs.PARAM_INVALID_ERR).Render()
		return q, false
	}
	return q, true
}
func drawManageTasks(c *gin.Context) {
	j := web.NewJsonResult(c)
	var list []drawManageTask
	ds.GetGDB().Raw("SELECT judge_task_id, CONCAT(task_year,'年',judge_task_name) task_name, judge_task_name, task_year, score_rounds, subject_count, experts_per_round, judge_task_type, task_start_time, task_end_time FROM ea_judge_task_info WHERE deleted=1 AND task_mode=1 AND publish_state=2 ORDER BY task_year DESC, judge_task_id DESC").Scan(&list)
	j.SetList(list).Render()
}
func ensureDrawRows(taskId int64, userId int64) error {
	var t drawManageTask
	r := ds.GetGDB().Raw("SELECT judge_task_id,score_rounds,subject_count,experts_per_round FROM ea_judge_task_info WHERE judge_task_id=? AND deleted=1 AND task_mode=1 AND publish_state=2", taskId).Scan(&t)
	if r.RowsAffected < 1 {
		return fmt.Errorf("任务不存在或尚未发布")
	}
	var n int64
	ds.GetGDB().Raw("SELECT COUNT(1) FROM ea_judge_draw_records WHERE judge_task_id=? AND deleted=1", taskId).Scan(&n)
	if n > 0 {
		return nil
	}
	for round := 1; round <= t.ScoreRounds; round++ {
		for subject := 0; subject < t.SubjectCount; subject++ {
			if e := ds.GetGDB().Exec("INSERT INTO ea_judge_draw_records(judge_task_id,draw_times,subject_name,required_num,submit_state,on_campus_total_num,on_campus_attend_num,on_campus_unattend_num,off_campus_total_num,off_campus_attend_num,off_campus_unattend_num,deleted,create_by,create_time,update_by,update_time) VALUES(?,?, '',?,1,0,0,0,0,0,0,1,?,NOW(),?,NOW())", taskId, round, t.ExpertsPerRound, userId, userId).Error; e != nil {
				return e
			}
		}
	}
	return nil
}
func drawManageRows(c *gin.Context) {
	q, ok := bindDrawReq(c)
	if !ok {
		return
	}
	j := web.NewJsonResult(c)
	if e := ensureDrawRows(q.JudgeTaskId, context.NewUserContext(c).GetUserId()); e != nil {
		j.SetErr(e.Error()).Render()
		return
	}
	var list []drawManageRow
	ds.GetGDB().Raw("SELECT r.rules_id,r.judge_task_id,r.draw_times,r.subject_name,r.required_num,r.submit_state,SUM(CASE WHEN s.expert_role=1 AND s.deleted=1 THEN 1 ELSE 0 END) scoring_count,SUM(CASE WHEN s.expert_role=2 AND s.deleted=1 THEN 1 ELSE 0 END) backup_count FROM ea_judge_draw_records r LEFT JOIN ea_judge_draw_selection s ON s.rules_id=r.rules_id AND s.deleted=1 WHERE r.judge_task_id=? AND r.deleted=1 GROUP BY r.rules_id ORDER BY r.draw_times,r.rules_id", q.JudgeTaskId).Scan(&list)
	j.SetList(list).Render()
}
func drawManageExperts(c *gin.Context) {
	q, ok := bindDrawReq(c)
	if !ok {
		return
	}
	j := web.NewJsonResult(c)
	var rule drawManageRow
	if ds.GetGDB().Raw("SELECT rules_id,judge_task_id,draw_times FROM ea_judge_draw_records WHERE rules_id=? AND deleted=1", q.RulesId).Scan(&rule).RowsAffected < 1 {
		j.SetErr("抽取记录不存在").Render()
		return
	}
	args := []interface{}{q.RulesId, rule.JudgeTaskId, rule.JudgeTaskId, rule.DrawTimes, q.RulesId}
	sql := "SELECT e.expert_id,e.dept_id,e.name,e.sex,e.phone_no,e.id_card,e.expert_category,COALESCE(d.dept_name,'') dept_name,e.unit_name,e.type,e.title,e.init_edu,e.init_major,e.final_edu,e.final_major,e.good_subjects,e.bank_card,e.intro,e.evaluation,COALESCE(cur.expert_role,0) expert_role,COALESCE(cur.submit_state,1) submit_state FROM ea_judge_experts e LEFT JOIN ea_dept_info d ON d.dept_id=e.dept_id LEFT JOIN ea_judge_draw_selection cur ON cur.rules_id=? AND cur.expert_id=e.expert_id AND cur.deleted=1 WHERE (COALESCE(cur.expert_role,0) IN (1,2) OR (e.deleted=1 AND e.allow_draw=1 AND e.audit_status=3 AND FIND_IN_SET((SELECT judge_task_type FROM ea_judge_task_info WHERE judge_task_id=?),e.type)>0)) AND NOT EXISTS(SELECT 1 FROM ea_judge_draw_selection used WHERE used.judge_task_id=? AND used.draw_times=? AND used.expert_id=e.expert_id AND used.deleted=1 AND used.rules_id<>?)"
	if strings.TrimSpace(q.Name) != "" {
		sql += " AND e.name LIKE ?"
		args = append(args, "%"+strings.TrimSpace(q.Name)+"%")
	}
	if q.ExpertCategory > 0 {
		sql += " AND e.expert_category=?"
		args = append(args, q.ExpertCategory)
	}
	if q.ExpertRole > 0 {
		sql += " AND COALESCE(cur.expert_role,0)=?"
		args = append(args, q.ExpertRole)
	}
	if q.SubmitState > 0 {
		sql += " AND COALESCE(cur.submit_state,1)=?"
		args = append(args, q.SubmitState)
	}
	sql += " ORDER BY e.expert_category,e.expert_id"
	var list []drawManageExpert
	ds.GetGDB().Raw(sql, args...).Scan(&list)
	j.SetList(list).Render()
}
func saveDrawSelection(c *gin.Context) {
	q, ok := bindDrawReq(c)
	if !ok {
		return
	}
	j := web.NewJsonResult(c)
	if q.RulesId < 1 || q.ExpertId < 1 || (q.ExpertRole != 1 && q.ExpertRole != 2) {
		j.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	var r drawManageRow
	ds.GetGDB().Raw("SELECT rules_id,judge_task_id,draw_times,required_num,submit_state FROM ea_judge_draw_records WHERE rules_id=? AND deleted=1", q.RulesId).Scan(&r)
	var scoring int64
	ds.GetGDB().Raw("SELECT COUNT(1) FROM ea_judge_draw_selection WHERE rules_id=? AND expert_role=1 AND deleted=1", q.RulesId).Scan(&scoring)
	if q.ExpertRole == 1 && scoring >= int64(r.RequiredNum) {
		j.SetErr("评分专家已达到设置数量").Render()
		return
	}
	var count int64
	ds.GetGDB().Raw("SELECT COUNT(1) FROM ea_judge_draw_selection WHERE judge_task_id=? AND draw_times=? AND expert_id=? AND deleted=1 AND rules_id<>?", r.JudgeTaskId, r.DrawTimes, q.ExpertId, q.RulesId).Scan(&count)
	if count > 0 {
		j.SetErr("同一轮次不可重复抽取该专家").Render()
		return
	}
	uid := context.NewUserContext(c).GetUserId()
	// 若之前已移除过该专家（deleted>1 表示已删除记录），则新插入记录的 replaced_from 指向被替换的原记录ID，
	// 与 ea_judge_expert_report.replaced_from 一致，便于回溯“删除->重新抽取”的关联。
	var replacedFrom int64
	ds.GetGDB().Raw("SELECT selection_id FROM ea_judge_draw_selection WHERE rules_id=? AND expert_id=? AND deleted>1 LIMIT 1", q.RulesId, q.ExpertId).Scan(&replacedFrom)
	e := ds.GetGDB().Exec("INSERT INTO ea_judge_draw_selection(rules_id,judge_task_id,draw_times,expert_id,expert_role,submit_state,deleted,replaced_from,create_by,create_time,update_by,update_time) VALUES(?,?,?,?,?,1,1,?,?,NOW(),?,NOW()) ON DUPLICATE KEY UPDATE expert_role=VALUES(expert_role),deleted=1,update_by=VALUES(update_by),update_time=NOW()", q.RulesId, r.JudgeTaskId, r.DrawTimes, q.ExpertId, q.ExpertRole, replacedFrom, uid, uid).Error
	if e != nil {
		j.SetErr(e.Error()).Render()
		return
	}
	j.Render()
}
func removeDrawSelection(c *gin.Context) {
	q, ok := bindDrawReq(c)
	if !ok {
		return
	}
	j := web.NewJsonResult(c)
	uid := context.NewUserContext(c).GetUserId()
	// 招办移除抽取专家：保留原记录（deleted 写入 Unix 时间戳作为删除标记），便于回溯；
	// 唯一键 (rules_id, expert_id, deleted) 含 deleted，后续重新抽取同一专家时会新增一条 deleted=1 的记录，
	// 而不是把原记录覆盖成 deleted=2，从而避免唯一键冲突并完整保留历史操作痕迹。
	ds.GetGDB().Exec("UPDATE ea_judge_draw_selection SET deleted=UNIX_TIMESTAMP()+selection_id,update_by=?,update_time=NOW() WHERE rules_id=? AND expert_id=? AND deleted=1", uid, q.RulesId, q.ExpertId)
	j.Render()
}
func updateDrawRule(c *gin.Context) {
	q, ok := bindDrawReq(c)
	if !ok {
		return
	}
	j := web.NewJsonResult(c)
	q.SubjectName = strings.TrimSpace(q.SubjectName)
	if q.RequiredNum < 1 {
		j.SetErr("抽取专家数必须大于0").Render()
		return
	}
	var r drawManageRow
	ds.GetGDB().Raw("SELECT rules_id,judge_task_id,draw_times FROM ea_judge_draw_records WHERE rules_id=? AND deleted=1", q.RulesId).Scan(&r)
	var selected int64
	ds.GetGDB().Raw("SELECT COUNT(1) FROM ea_judge_draw_selection WHERE rules_id=? AND expert_role=1 AND deleted=1", q.RulesId).Scan(&selected)
	if int64(q.RequiredNum) < selected {
		j.SetErr("抽取专家数不能小于当前已抽取评分专家数").Render()
		return
	}
	if q.SubjectName != "" {
		var dup int64
		ds.GetGDB().Raw("SELECT COUNT(1) FROM ea_judge_draw_records WHERE judge_task_id=? AND draw_times=? AND subject_name=? AND rules_id<>? AND deleted=1", r.JudgeTaskId, r.DrawTimes, q.SubjectName, q.RulesId).Scan(&dup)
		if dup > 0 {
			j.SetErr("同一轮次下科目名称不允许重复").Render()
			return
		}
	}
	ds.GetGDB().Exec("UPDATE ea_judge_draw_records SET subject_name=?,required_num=?,update_by=?,update_time=NOW() WHERE rules_id=?", q.SubjectName, q.RequiredNum, context.NewUserContext(c).GetUserId(), q.RulesId)
	j.Render()
}
func submitDrawSelection(c *gin.Context) {
	q, ok := bindDrawReq(c)
	if !ok {
		return
	}
	j := web.NewJsonResult(c)
	var r drawManageRow
	ds.GetGDB().Raw("SELECT rules_id,judge_task_id,draw_times,subject_name,required_num FROM ea_judge_draw_records WHERE rules_id=? AND deleted=1", q.RulesId).Scan(&r)
	if strings.TrimSpace(q.SubjectName) == "" {
		j.SetErr("科目名称必须填写").Render()
		return
	}
	q.RequiredNum = r.RequiredNum
	if q.RequiredNum < 1 {
		j.SetErr("请先设置抽取专家数").Render()
		return
	}
	var scoring int64
	ds.GetGDB().Raw("SELECT COUNT(1) FROM ea_judge_draw_selection WHERE rules_id=? AND expert_role=1 AND deleted=1", q.RulesId).Scan(&scoring)
	if scoring < 1 {
		j.SetErr("未设置评分专家不允许提交").Render()
		return
	}
	if scoring > int64(q.RequiredNum) {
		j.SetErr("评分专家数量超过设置数量").Render()
		return
	}
	var dup int64
	ds.GetGDB().Raw("SELECT COUNT(1) FROM ea_judge_draw_records WHERE judge_task_id=? AND draw_times=? AND subject_name=? AND rules_id<>? AND deleted=1", r.JudgeTaskId, r.DrawTimes, strings.TrimSpace(q.SubjectName), q.RulesId).Scan(&dup)
	if dup > 0 {
		j.SetErr("同一轮次下科目名称不允许重复").Render()
		return
	}
	uid := context.NewUserContext(c).GetUserId()
	ds.GetGDB().Exec("UPDATE ea_judge_draw_records SET subject_name=?,submit_state=2,update_by=?,update_time=NOW() WHERE rules_id=?", strings.TrimSpace(q.SubjectName), uid, q.RulesId)
	ds.GetGDB().Exec("UPDATE ea_judge_draw_selection SET submit_state=2,update_by=?,update_time=NOW() WHERE rules_id=? AND deleted=1", uid, q.RulesId)
	j.Render()
}
func exportDrawSelection(c *gin.Context) {
	q, ok := bindDrawReq(c)
	if !ok {
		return
	}
	if q.JudgeTaskId <= 0 {
		web.NewJsonResult(c).SetErr("请选择导出任务").Render()
		return
	}
	var dictRows []map[string]interface{}
	ds.GetGDB().Raw("SELECT dict_value,dict_label FROM sys_dict_data WHERE dict_type=? AND state='1' AND deleted=1", cst.BIZ_DICT_EXPERT_TYPE).Scan(&dictRows)
	expertTypeMap := make(map[string]string)
	for _, r := range dictRows {
		expertTypeMap[fmt.Sprint(r["dict_value"])] = fmt.Sprint(r["dict_label"])
	}

	var rows []map[string]interface{}
	ds.GetGDB().Raw("SELECT CONCAT(t.task_year,'年',t.judge_task_name) '任务名称',r.draw_times '评分轮次',r.subject_name '评分科目',CONCAT(DATE(t.task_start_time),'至',DATE(t.task_end_time)) '评分日期',e.name '姓名',e.sex '性别',e.phone_no '手机号',e.id_card '证件号',CASE e.expert_category WHEN 1 THEN '校内' ELSE '校外' END '专家类别',COALESCE(d.dept_name,e.unit_name) '所属学院/所在单位',e.type '专家类型',e.init_edu '初始学历',e.init_major '初始学历所学专业',e.final_edu '最终学历',e.final_major '最终学历所学专业',e.good_subjects '擅长科目',e.bank_card '银行卡号',e.intro '专家介绍',e.evaluation '专家评价',CASE s.expert_role WHEN 1 THEN '评分专家' ELSE '备用专家' END '专家模式',s.submit_state '提交状态' FROM ea_judge_draw_selection s JOIN ea_judge_draw_records r ON r.rules_id=s.rules_id JOIN ea_judge_task_info t ON t.judge_task_id=s.judge_task_id JOIN ea_judge_experts e ON e.expert_id=s.expert_id LEFT JOIN ea_dept_info d ON d.dept_id=e.dept_id WHERE s.judge_task_id=? AND s.deleted=1 ORDER BY r.draw_times,s.expert_role,s.selection_id", q.JudgeTaskId).Scan(&rows)
	// 数据库中 intro 字段为 text，map 转出来是 []byte，excelize 不识别而写入空字符串。
	// 统一将 []byte 转换为 string，保证专家介绍等长文本列能正常写出。
	for _, row := range rows {
		for k, v := range row {
			if b, ok := v.([]byte); ok {
				row[k] = string(b)
			}
		}
	}
	for _, row := range rows {
		if v, ok := row["专家类型"]; ok {
			parts := strings.Split(fmt.Sprint(v), ",")
			labels := make([]string, 0, len(parts))
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				if lb, exists := expertTypeMap[p]; exists {
					labels = append(labels, lb)
				} else {
					labels = append(labels, p)
				}
			}
			row["专家类型"] = strings.Join(labels, ",")
		}
		if v, ok := row["初始学历"]; ok {
			if n, err := strconv.Atoi(fmt.Sprint(v)); err == nil {
				row["初始学历"] = cst.GetEduTypeLabelByValue(n)
			}
		}
		if v, ok := row["最终学历"]; ok {
			if n, err := strconv.Atoi(fmt.Sprint(v)); err == nil {
				row["最终学历"] = cst.GetEduTypeLabelByValue(n)
			}
		}
		if v, ok := row["提交状态"]; ok {
			if fmt.Sprint(v) == "2" {
				row["提交状态"] = "已提交"
			} else {
				row["提交状态"] = "未提交"
			}
		}
	}
	f := excelize.NewFile()
	sheet := "评委专家"
	f.SetSheetName("Sheet1", sheet)
	headers := []string{"任务名称", "评分轮次", "评分科目", "评分日期", "姓名", "性别", "手机号", "证件号", "专家类别", "所属学院/所在单位", "专家类型", "初始学历", "初始学历所学专业", "最终学历", "最终学历所学专业", "擅长科目", "银行卡号", "专家介绍", "专家评价", "专家模式", "提交状态"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	for ri, row := range rows {
		for ci, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(ci+1, ri+2)
			f.SetCellValue(sheet, cell, row[h])
		}
	}
	buf, e := f.WriteToBuffer()
	if e != nil {
		web.NewJsonResult(c).SetErr(e.Error()).Render()
		return
	}
	web.DataPackageExcel(c, buf.Bytes())
}

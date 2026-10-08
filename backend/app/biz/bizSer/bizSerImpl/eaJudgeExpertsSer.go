package bizSerImpl

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"yx-exam-affair-gm/app/biz/bizDao"
	"yx-exam-affair-gm/app/biz/bizDao/bizDaoImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	"yx-exam-affair-gm/app/cst"
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"
	cmutils "yx-exam-affair-gm/app/utils"

	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds"
	"github.com/sealsee/web-base/public/ds/tx"
	"github.com/sealsee/web-base/public/utils/file/excel"
	"github.com/sealsee/web-base/public/utils/snowflake"
	"gorm.io/gorm"
)

/*
 * desc: 评委专家库服务层实现
 * author: init code
 * date: 2025/06/06
 */
type EaJudgeExpertsService struct{}

var eaJudgeExpertsService *EaJudgeExpertsService
var eaJudgeExpertsDao bizDao.IEaJudgeExpertsDao

func init() {
	eaJudgeExpertsService = &EaJudgeExpertsService{}
	eaJudgeExpertsDao = bizDaoImpl.NewEaJudgeExpertsDao()
}

func NewEaJudgeExpertsService() *EaJudgeExpertsService {
	return eaJudgeExpertsService
}

func (judgeExperts *EaJudgeExpertsService) GetEaJudgeExpertsById(id int64) *bizDo.EaJudgeExperts {
	return eaJudgeExpertsDao.GetById(id)
}

func (judgeExperts *EaJudgeExpertsService) GetLatestAuditHistory(id int64) *bizDo.EaJudgeExpertsAuditHistory {
	if id < 1 {
		return nil
	}
	history := new(bizDo.EaJudgeExpertsAuditHistory)
	res := ds.GetGDB().Where("expert_id = ? AND audit_status = ?", id, cst.EXPERT_AUDIT_PASSED).
		Order("change_time desc").First(history)
	if res.Error != nil || res.RowsAffected < 1 {
		return nil
	}
	return history
}

// GetNextPendingExpert 返回下一条待审核专家，用于招办连续审核。
// 查询条件 q 取自列表页的检索条件，保证连续审核的范围与当前筛选结果一致；
// q.ExpertId 为刚审核完的记录，优先向后取（expert_id > 当前值），取不到时再回绕到第一条。
func (judgeExperts *EaJudgeExpertsService) GetNextPendingExpert(q *bizDto.EaJudgeExpertsQuery) *bizDo.EaJudgeExperts {
	if q == nil {
		q = new(bizDto.EaJudgeExpertsQuery)
	}
	currentId := q.ExpertId

	buildQuery := func() *gorm.DB {
		db := ds.GetGDB().Model(&bizDo.EaJudgeExperts{}).
			Where("audit_status = ? AND deleted = ?", cst.EXPERT_AUDIT_PENDING, common.Normal)
		if name := strings.TrimSpace(q.Name); name != "" {
			db = db.Where("name LIKE ?", "%"+name+"%")
		}
		if idCard := strings.TrimSpace(q.IdCard); idCard != "" {
			db = db.Where("id_card LIKE ?", "%"+idCard+"%")
		}
		if phoneNo := strings.TrimSpace(q.PhoneNo); phoneNo != "" {
			db = db.Where("phone_no LIKE ?", "%"+phoneNo+"%")
		}
		if goodSubjects := strings.TrimSpace(q.GoodSubjects); goodSubjects != "" {
			db = db.Where("good_subjects LIKE ?", "%"+goodSubjects+"%")
		}
		if unitName := strings.TrimSpace(q.UnitName); unitName != "" {
			db = db.Where("unit_name LIKE ?", "%"+unitName+"%")
		}
		if title := strings.TrimSpace(q.Title); title != "" {
			db = db.Where("title LIKE ?", "%"+title+"%")
		}
		if sex := strings.TrimSpace(q.Sex); sex != "" {
			db = db.Where("sex = ?", sex)
		}
		if expertType := strings.TrimSpace(q.Type); expertType != "" {
			db = db.Where("FIND_IN_SET(?, type) > 0", expertType)
		}
		if q.DeptId != 0 {
			db = db.Where("dept_id = ?", q.DeptId)
		}
		if q.ExpertCategory > 0 {
			db = db.Where("expert_category = ?", q.ExpertCategory)
		}
		if q.InitEdu > 0 {
			db = db.Where("init_edu = ?", q.InitEdu)
		}
		if q.FinalEdu > 0 {
			db = db.Where("final_edu = ?", q.FinalEdu)
		}
		if q.AllowDraw > 0 {
			db = db.Where("allow_draw = ?", q.AllowDraw)
		}
		if q.IsEvaluated == cst.YES {
			db = db.Where("TRIM(COALESCE(evaluation, '')) <> ''")
		} else if q.IsEvaluated == cst.NO {
			db = db.Where("TRIM(COALESCE(evaluation, '')) = ''")
		}
		return db
	}

	expert := new(bizDo.EaJudgeExperts)
	var res *gorm.DB
	if currentId > 0 {
		res = buildQuery().Where("expert_id > ?", currentId).Order("expert_id asc").First(expert)
		if res.Error != nil || res.RowsAffected < 1 {
			expert = new(bizDo.EaJudgeExperts)
			res = buildQuery().Where("expert_id <> ?", currentId).Order("expert_id asc").First(expert)
		}
	} else {
		res = buildQuery().Order("expert_id asc").First(expert)
	}
	if res.Error != nil || res.RowsAffected < 1 {
		return nil
	}
	if dept := eaDeptInfoService.GetAllDeptMap()[expert.DeptId]; dept != nil {
		expert.DeptName = dept.DeptName
	}
	return expert
}

func (judgeExperts *EaJudgeExpertsService) ListEaJudgeExperts(q *bizDto.EaJudgeExpertsQuery, p *page.Page) []*bizDo.EaJudgeExperts {
	// 前端未传排序时的默认顺序：
	// - 禁止抽取管理（BannedOnly=true）：按"加入禁止抽取"的时间顺序展示（earliest create_time 升序，确保与用户新增/导入时的预期一致）。
	// - 其他列表：默认按主键 expert_id 升序，保证与连续审核(next_pending)的游标顺序一致。
	if q.GetOrders() == "" {
		if q.BannedOnly {
			q.AddOrderAsc("(SELECT MIN(b.create_time) FROM ea_judge_expert_ban b WHERE b.expert_id = ea_judge_experts.expert_id AND b.deleted = 1)")
		}
		q.AddOrderAsc("expert_id")
	}
	// 需要支持模糊搜索的字段
	q.AddLikeAll("name", q.Name)
	q.AddLikeAll("id_card", q.IdCard)
	q.AddLikeAll("phone_no", q.PhoneNo)
	q.AddLikeAll("good_subjects", q.GoodSubjects)
	q.AddLikeAll("unit_name", q.UnitName)
	q.AddLikeAll("title", q.Title)
	q.AddLikeAll("init_major", q.InitMajor)
	q.AddLikeAll("final_major", q.FinalMajor)
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	if strings.TrimSpace(q.Type) != "" {
		conditions = append(conditions, "FIND_IN_SET(?, type) > 0")
		args = append(args, strings.TrimSpace(q.Type))
		q.Type = ""
	}
	if q.IsEvaluated == cst.YES {
		conditions = append(conditions, "TRIM(COALESCE(evaluation, '')) <> ''")
	} else if q.IsEvaluated == cst.NO {
		conditions = append(conditions, "TRIM(COALESCE(evaluation, '')) = ''")
	}
	if q.BannedOnly {
		conditions = append(conditions, "EXISTS (SELECT 1 FROM ea_judge_expert_ban b WHERE b.expert_id = ea_judge_experts.expert_id AND b.deleted = 1)")
	}
	if q.ExcludeBanned {
		conditions = append(conditions, "NOT EXISTS (SELECT 1 FROM ea_judge_expert_ban b WHERE b.expert_id = ea_judge_experts.expert_id AND b.deleted = 1)")
	}
	var list []*bizDo.EaJudgeExperts
	if len(conditions) > 0 {
		list = eaJudgeExpertsDao.ListWithCondition(q, p, strings.Join(conditions, " AND "), args...)
	} else {
		list = eaJudgeExpertsDao.List(q, p)
	}
	deptMap := eaDeptInfoService.GetAllDeptMap()
	for _, item := range list {
		item.HasEvaluation = strings.TrimSpace(item.Evaluation) != ""
		if deptMap[item.DeptId] != nil {
			item.DeptName = deptMap[item.DeptId].DeptName
		}
	}
	return list
}

func (judgeExperts *EaJudgeExpertsService) ListAllEaJudgeExperts(q *bizDto.EaJudgeExpertsQuery) []*bizDo.EaJudgeExperts {
	listJudgeExpertAll := make([]*bizDo.EaJudgeExperts, 0)
	typeFilter := q.Type
	q.Type = ""
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE
	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		var list []*bizDo.EaJudgeExperts
		if strings.TrimSpace(typeFilter) != "" {
			list = eaJudgeExpertsDao.ListWithCondition(q, pageInfo, "FIND_IN_SET(?, type) > 0", strings.TrimSpace(typeFilter))
		} else {
			list = eaJudgeExpertsDao.List(q, pageInfo)
		}
		if len(list) == 0 {
			break
		}
		listJudgeExpertAll = append(listJudgeExpertAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}
	return listJudgeExpertAll
}

func (judgeExperts *EaJudgeExpertsService) InsertEaJudgeExperts(judgeExpertsDO *bizDo.EaJudgeExperts) (success bool, errInfo string) {
	if success, errInfo = validateJudgeExpert(judgeExpertsDO); !success {
		return
	}
	// 一个证件号只能维护一条专家数据
	exit := eaJudgeExpertsDao.Count(&bizDto.EaJudgeExpertsQuery{IdCard: judgeExpertsDO.IdCard}) > 0
	if exit {
		return false, "证件号已存在，一个证件号只能维护一条专家数据"
	}

	judgeExpertsDO.AllowDraw = cst.YES
	if judgeExpertsDO.DataFrom == cst.ADDBYADMISSIONS {
		judgeExpertsDO.AuditStatus = cst.EXPERT_AUDIT_PASSED
	} else {
		judgeExpertsDO.AuditStatus = cst.EXPERT_AUDIT_UNSUBMITTED
		if judgeExpertsDO.SubmitAction == cst.EXPERT_SUBMIT_ACTION_SUBMIT {
			judgeExpertsDO.AuditStatus = cst.EXPERT_AUDIT_PENDING
		} else if judgeExpertsDO.SubmitAction != "" && judgeExpertsDO.SubmitAction != cst.EXPERT_SUBMIT_ACTION_DRAFT {
			return false, "请选择暂存或提交审核"
		}
	}
	judgeExpertsDO.AuditRemark = ""
	judgeExpertsDO.Deleted = common.Normal

	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaJudgeExpertsDao.Insert(gtx, judgeExpertsDO)
	})
	return
}

func (judgeExperts *EaJudgeExpertsService) UpdateEaJudgeExpertsById(judgeExpertsDO *bizDo.EaJudgeExperts, sessionDeptId int64) (success bool, errInfo string) {
	dbExpert := eaJudgeExpertsDao.GetById(judgeExpertsDO.ExpertId)
	if dbExpert == nil {
		return false, "专家数据不存在"
	}
	isCollege := sessionDeptId > 0
	if isCollege {
		if dbExpert.DeptId != sessionDeptId {
			return false, "学院只能编辑本学院的专家"
		}
		judgeExpertsDO.DeptId = sessionDeptId
		judgeExpertsDO.ExpertCategory = cst.EXPERT_CATEGORY_INTERNAL
		judgeExpertsDO.UnitName = ""
	} else {
		if dbExpert.AuditStatus == cst.EXPERT_AUDIT_UNSUBMITTED {
			return false, "学院未提交的专家数据，招办不可编辑"
		}
	}
	if dbExpert.IdCard != judgeExpertsDO.IdCard {
		return false, "证件号不允许修改"
	}
	if success, errInfo = validateJudgeExpert(judgeExpertsDO); !success {
		return
	}
	changed := hasJudgeExpertChanged(dbExpert, judgeExpertsDO)
	if isCollege {
		switch judgeExpertsDO.SubmitAction {
		case cst.EXPERT_SUBMIT_ACTION_DRAFT:
			// 暂存时若内容未变化且当前非待审核状态（未提交/审核不通过），不允许保存；
			// 待审核状态下点暂存允许不修改内容，回退为未提交
			if !changed && dbExpert.AuditStatus != cst.EXPERT_AUDIT_PENDING {
				return false, "未进行信息修改"
			}
			judgeExpertsDO.AuditStatus = cst.EXPERT_AUDIT_UNSUBMITTED
		case cst.EXPERT_SUBMIT_ACTION_SUBMIT:
			// 未提交状态提交审核：允许不修改内容，直接更新为待审核
			if dbExpert.AuditStatus != cst.EXPERT_AUDIT_UNSUBMITTED {
				if !changed {
					return false, "未进行信息修改"
				}
			}
			judgeExpertsDO.AuditStatus = cst.EXPERT_AUDIT_PENDING
		default:
			return false, "请选择暂存或提交审核"
		}
	} else {
		judgeExpertsDO.AuditStatus = cst.EXPERT_AUDIT_PASSED
	}
	judgeExpertsDO.AuditRemark = ""
	uWhere := &bizDo.EaJudgeExperts{ExpertId: judgeExpertsDO.ExpertId}
	judgeExpertsDO.ExpertId = 0
	judgeExpertsDO.SetNullableCols("unit_name", "final_edu", "final_major", "audit_remark", "intro")
	tx.ExecGTx(func(gtx tx.Db) {
		if dbExpert.AuditStatus == cst.EXPERT_AUDIT_PASSED && changed {
			snapshot, err := json.Marshal(dbExpert)
			if err != nil {
				success = false
				return
			}
			history := &bizDo.EaJudgeExpertsAuditHistory{
				ExpertId: dbExpert.ExpertId, DeptId: dbExpert.DeptId,
				AuditStatus: dbExpert.AuditStatus, SnapshotData: string(snapshot),
				ChangeBy: judgeExpertsDO.UpdateBy, ChangeTime: time.Now(),
			}
			if result := gtx.Create(history); result.Error != nil || result.RowsAffected < 1 {
				success = false
				return
			}
		}
		success = eaJudgeExpertsDao.UpdateNew(gtx, uWhere, judgeExpertsDO) > 0
	})
	if !success {
		return false, "专家信息更新失败"
	}
	return true, ""
}

func (judgeExperts *EaJudgeExpertsService) SubmitForAudit(id int64, sessionDeptId int64, submitBy int64) (bool, string) {
	if sessionDeptId <= 0 {
		return false, "仅学院用户可提交审核"
	}
	dbExpert := eaJudgeExpertsDao.GetById(id)
	if dbExpert == nil || dbExpert.DeptId != sessionDeptId {
		return false, "学院只能提交本学院的专家"
	}
	if dbExpert.AuditStatus != cst.EXPERT_AUDIT_UNSUBMITTED {
		return false, "仅未提交状态的专家可提交审核"
	}
	uData := &bizDo.EaJudgeExperts{AuditStatus: cst.EXPERT_AUDIT_PENDING}
	uData.SetUpdateBy(submitBy)
	success := false
	tx.ExecGTx(func(gtx tx.Db) {
		success = eaJudgeExpertsDao.UpdateNew(gtx, &bizDo.EaJudgeExperts{ExpertId: id}, uData) > 0
	})
	if !success {
		return false, "提交审核失败"
	}
	return true, ""
}

// 学院批量提交专家审核：ids 非空时提交指定记录，ids 为空时提交本学院全部未提交数据
func (judgeExperts *EaJudgeExpertsService) SubmitForAuditBatch(ids []int64, sessionDeptId int64, submitBy int64) (int, string) {
	if sessionDeptId <= 0 {
		return 0, "仅学院用户可提交审核"
	}
	affected := 0
	tx.ExecGTx(func(gtx tx.Db) {
		targetIds := ids
		if len(targetIds) == 0 {
			var allIds []int64
			ds.GetGDB().Raw("SELECT expert_id FROM ea_judge_experts WHERE dept_id = ? AND audit_status = ? AND deleted = ?",
				sessionDeptId, cst.EXPERT_AUDIT_UNSUBMITTED, common.Normal).Scan(&allIds)
			targetIds = allIds
		}
		for _, id := range targetIds {
			if id <= 0 {
				continue
			}
			dbExpert := eaJudgeExpertsDao.GetById(id)
			if dbExpert == nil || dbExpert.DeptId != sessionDeptId || dbExpert.AuditStatus != cst.EXPERT_AUDIT_UNSUBMITTED {
				continue
			}
			uData := &bizDo.EaJudgeExperts{AuditStatus: cst.EXPERT_AUDIT_PENDING}
			uData.SetUpdateBy(submitBy)
			if eaJudgeExpertsDao.UpdateNew(gtx, &bizDo.EaJudgeExperts{ExpertId: id}, uData) > 0 {
				affected++
			}
		}
	})
	return affected, ""
}

func (judgeExperts *EaJudgeExpertsService) AuditEaJudgeExperts(id int64, action string, reason string, auditBy int64) (bool, string) {
	dbExpert := eaJudgeExpertsDao.GetById(id)
	if dbExpert == nil {
		return false, "专家数据不存在"
	}
	if dbExpert.AuditStatus != cst.EXPERT_AUDIT_PENDING {
		return false, "仅待审核状态的专家支持审核"
	}
	if action != cst.EXPERT_AUDIT_ACTION_PASS && action != cst.EXPERT_AUDIT_ACTION_REJECT {
		return false, "审核操作不合法"
	}
	reason = strings.TrimSpace(reason)
	if action == cst.EXPERT_AUDIT_ACTION_REJECT && reason == "" {
		return false, "审核不通过时必须填写不通过原因"
	}
	status := cst.EXPERT_AUDIT_PASSED
	if action == cst.EXPERT_AUDIT_ACTION_REJECT {
		status = cst.EXPERT_AUDIT_REJECTED
	}
	uData := &bizDo.EaJudgeExperts{AuditStatus: status, AuditRemark: reason}
	uData.SetUpdateBy(auditBy)
	uData.SetNullableCols("audit_remark")
	success := false
	tx.ExecGTx(func(gtx tx.Db) {
		success = eaJudgeExpertsDao.UpdateNew(gtx, &bizDo.EaJudgeExperts{ExpertId: id}, uData) > 0
	})
	if !success {
		return false, "审核操作失败"
	}
	return true, ""
}

func (judgeExperts *EaJudgeExpertsService) EvaluateEaJudgeExperts(id int64, evaluation string, evaluateBy int64) (bool, string) {
	dbExpert := eaJudgeExpertsDao.GetById(id)
	if dbExpert == nil || dbExpert.AuditStatus != cst.EXPERT_AUDIT_PASSED {
		return false, "仅审核通过的专家可以评价"
	}
	evaluation = strings.TrimSpace(evaluation)
	if utf8.RuneCountInString(evaluation) > 2000 {
		return false, "评价内容不能超过2000字符"
	}
	uData := &bizDo.EaJudgeExperts{Evaluation: evaluation}
	uData.SetUpdateBy(evaluateBy)
	uData.SetNullableCols("evaluation")
	success := false
	tx.ExecGTx(func(gtx tx.Db) {
		success = eaJudgeExpertsDao.UpdateNew(gtx, &bizDo.EaJudgeExperts{ExpertId: id}, uData) > 0
	})
	if !success {
		return false, "评价保存失败"
	}
	return true, ""
}

func (judgeExperts *EaJudgeExpertsService) ChangeAllowDraw(id int64, allowDraw int, updateBy int64) (bool, string) {
	if allowDraw != cst.YES && allowDraw != cst.NO {
		return false, "状态参数错误"
	}
	dbExpert := eaJudgeExpertsDao.GetById(id)
	if dbExpert == nil {
		return false, "专家数据不存在"
	}
	uData := &bizDo.EaJudgeExperts{AllowDraw: allowDraw}
	uData.SetUpdateBy(updateBy)
	success := false
	tx.ExecGTx(func(gtx tx.Db) {
		success = eaJudgeExpertsDao.UpdateNew(gtx, &bizDo.EaJudgeExperts{ExpertId: id}, uData) > 0
	})
	if !success {
		return false, "状态更新失败"
	}
	return true, ""
}

func (judgeExperts *EaJudgeExpertsService) BanEaJudgeExpert(id int64, updateBy int64) (bool, string) {
	expert := eaJudgeExpertsDao.GetById(id)
	if expert == nil {
		return false, "专家不存在"
	}
	if expert.AllowDraw != cst.YES {
		return false, "只能对启用状态的专家设置禁止抽取"
	}
	if expert.AuditStatus != cst.EXPERT_AUDIT_PASSED {
		return false, "只能对审核通过的专家设置禁止抽取"
	}
	var count int64
	ds.GetGDB().Model(&bizDo.EaJudgeExpertBan{}).Where("expert_id = ? AND deleted = ?", id, common.Normal).Count(&count)
	if count > 0 {
		return false, "专家已被禁止抽取，不允许重复添加"
	}
	ban := &bizDo.EaJudgeExpertBan{ExpertId: id, IdCard: expert.IdCard, Deleted: common.Normal, CreateBy: updateBy, CreateTime: time.Now(), UpdateBy: updateBy, UpdateTime: time.Now()}
	if res := ds.GetGDB().Create(ban); res.Error != nil || res.RowsAffected < 1 {
		return false, "禁止抽取失败"
	}
	return true, ""
}

func (judgeExperts *EaJudgeExpertsService) UnbanEaJudgeExperts(ids []int64, updateBy int64) (bool, string) {
	if len(ids) < 1 {
		return false, "请选择要删除的专家"
	}
	res := ds.GetGDB().Model(&bizDo.EaJudgeExpertBan{}).Where("expert_id IN ? AND deleted = ?", ids, common.Normal).
		Updates(map[string]interface{}{"deleted": common.Deleted, "update_by": updateBy, "update_time": time.Now()})
	if res.Error != nil {
		return false, "删除失败"
	}
	return true, ""
}

func (judgeExperts *EaJudgeExpertsService) UnbanAllEaJudgeExperts(updateBy int64) (bool, string) {
	res := ds.GetGDB().Model(&bizDo.EaJudgeExpertBan{}).Where("deleted = ?", common.Normal).
		Updates(map[string]interface{}{"deleted": common.Deleted, "update_by": updateBy, "update_time": time.Now()})
	if res.Error != nil {
		return false, "删除全部失败"
	}
	return true, ""
}

type HandlerImpBannedJudgeExpertsStruct struct {
	ExpertMap map[string]*bizDo.EaJudgeExperts
	BannedMap map[string]bool
	Ids       []int64
	Seen      map[string]bool
	Err       string
}

func (h *HandlerImpBannedJudgeExpertsStruct) Headers(headers []string) {
	if len(headers) != 1 || strings.TrimSpace(headers[0]) != "证件号" {
		h.Err = "表头数据错误，导入文件只能包含“证件号”一列"
	}
}

func (h *HandlerImpBannedJudgeExpertsStruct) Row(mp *map[string]string) excel.RunStatus {
	if h.Err != "" {
		return excel.Exit
	}
	m := *mp
	rowIdx := strconv.Itoa(cmutils.StrToInt(m["rid"]) + 1)
	idCard := strings.TrimSpace(m["证件号"])
	if idCard == "" {
		h.Err = "第" + rowIdx + "行，证件号不能为空"
		return excel.Exit
	}
	expert := h.ExpertMap[idCard]
	if expert == nil {
		h.Err = "第" + rowIdx + "行，专家不存在"
		return excel.Exit
	}
	if expert.AllowDraw != cst.YES {
		h.Err = "第" + rowIdx + "行，只能对启用状态的专家设置禁止抽取"
		return excel.Exit
	}
	if expert.AuditStatus != cst.EXPERT_AUDIT_PASSED {
		h.Err = "第" + rowIdx + "行，只能对审核通过的专家设置禁止抽取"
		return excel.Exit
	}
	if h.BannedMap[idCard] {
		h.Err = "第" + rowIdx + "行，专家已被禁止抽取，不允许重复添加"
		return excel.Exit
	}
	if h.Seen[idCard] {
		h.Err = "第" + rowIdx + "行，证件号重复"
		return excel.Exit
	}
	h.Seen[idCard] = true
	h.Ids = append(h.Ids, expert.ExpertId)
	return excel.Next
}

func (h *HandlerImpBannedJudgeExpertsStruct) After() {}

func (judgeExperts *EaJudgeExpertsService) ImportBannedEaJudgeExperts(q *bizDto.EaJudgeExpertsQuery) (bool, string) {
	experts := judgeExperts.ListAllEaJudgeExperts(new(bizDto.EaJudgeExpertsQuery))
	expertMap := make(map[string]*bizDo.EaJudgeExperts, len(experts))
	expertIdCard := make(map[int64]string, len(experts))
	for _, expert := range experts {
		// 同一证件号若存在多条记录，优先保留“启用”状态的专家
		if existing, ok := expertMap[expert.IdCard]; !ok || existing.AllowDraw != cst.YES {
			expertMap[expert.IdCard] = expert
		}
		expertIdCard[expert.ExpertId] = expert.IdCard
	}
	bannedMap := make(map[string]bool)
	var existingBans []bizDo.EaJudgeExpertBan
	ds.GetGDB().Where("deleted = ?", common.Normal).Find(&existingBans)
	for _, ban := range existingBans {
		if expert := eaJudgeExpertsDao.GetById(ban.ExpertId); expert != nil {
			bannedMap[expert.IdCard] = true
		}
	}
	handler := &HandlerImpBannedJudgeExpertsStruct{ExpertMap: expertMap, BannedMap: bannedMap, Seen: make(map[string]bool)}
	for _, url := range q.Files {
		excel.NewExcel().ImportWithUrl(url, handler)
		if handler.Err != "" {
			return false, handler.Err
		}
	}
	if len(handler.Ids) < 1 {
		return false, "没有导入数据"
	}
	bans := make([]*bizDo.EaJudgeExpertBan, 0, len(handler.Ids))
	for _, expertId := range handler.Ids {
		bans = append(bans, &bizDo.EaJudgeExpertBan{ExpertId: expertId, IdCard: expertIdCard[expertId], Deleted: common.Normal, CreateBy: q.OperId, CreateTime: time.Now(), UpdateBy: q.OperId, UpdateTime: time.Now()})
	}
	res := ds.GetGDB().CreateInBatches(bans, 100)
	if res.Error != nil || int(res.RowsAffected) != len(bans) {
		return false, "导入失败"
	}
	return true, fmt.Sprintf("成功导入 %d 条记录", len(handler.Ids))
}

func validateJudgeExpert(info *bizDo.EaJudgeExperts) (bool, string) {
	info.Name = strings.TrimSpace(info.Name)
	info.IdCard = strings.TrimSpace(info.IdCard)
	info.PhoneNo = strings.TrimSpace(info.PhoneNo)
	info.InitMajor = strings.TrimSpace(info.InitMajor)
	info.FinalMajor = strings.TrimSpace(info.FinalMajor)
	info.UnitName = strings.TrimSpace(info.UnitName)
	if utf8.RuneCountInString(info.Name) < 2 || utf8.RuneCountInString(info.Name) > 30 {
		return false, "姓名限制2-30字符"
	}
	if info.Sex != cst.MAN && info.Sex != cst.WOMAN {
		return false, "请选择正确的性别"
	}
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9()（）]{3,18}$`, info.IdCard); !matched {
		return false, "证件号必须为3-18位数字、大小写字母或括号"
	}
	if matched, _ := regexp.MatchString(`^1\d{10}$`, info.PhoneNo); !matched {
		return false, "手机号必须为1开头的11位数字"
	}
	if info.ExpertCategory != cst.EXPERT_CATEGORY_INTERNAL && info.ExpertCategory != cst.EXPERT_CATEGORY_EXTERNAL {
		return false, "请选择专家类别"
	}
	if info.ExpertCategory == cst.EXPERT_CATEGORY_INTERNAL {
		if info.DeptId <= 0 {
			return false, "校内专家必须选择所属学院"
		}
		info.UnitName = ""
	} else {
		if utf8.RuneCountInString(info.UnitName) < 1 || utf8.RuneCountInString(info.UnitName) > 30 {
			return false, "所在单位限制1-30字符"
		}
		info.DeptId = cst.BIZ_DEPTID_FOR_ADMISSIONS
	}
	if strings.TrimSpace(info.Type) == "" {
		return false, "专家类型不能为空"
	}
	if info.InitEdu <= 0 || utf8.RuneCountInString(info.InitMajor) < 1 || utf8.RuneCountInString(info.InitMajor) > 20 {
		return false, "初始学历及所学专业必填，专业限制1-20字符"
	}
	if info.FinalEdu > 0 && (utf8.RuneCountInString(info.FinalMajor) < 1 || utf8.RuneCountInString(info.FinalMajor) > 20) {
		return false, "选择最终学历后，所学专业必填且限制1-20字符"
	}
	if info.FinalEdu == 0 {
		info.FinalMajor = ""
	}
	info.Intro = strings.TrimSpace(info.Intro)
	if n := utf8.RuneCountInString(info.Intro); n > 1000 {
		return false, "专家介绍不能超过1000字符"
	}
	return true, ""
}

func hasJudgeExpertChanged(oldInfo, newInfo *bizDo.EaJudgeExperts) bool {
	return oldInfo.Name != newInfo.Name || oldInfo.Sex != newInfo.Sex || oldInfo.PhoneNo != newInfo.PhoneNo ||
		oldInfo.ExpertCategory != newInfo.ExpertCategory || oldInfo.DeptId != newInfo.DeptId ||
		oldInfo.UnitName != newInfo.UnitName || oldInfo.Type != newInfo.Type || oldInfo.Title != newInfo.Title ||
		oldInfo.InitEdu != newInfo.InitEdu || oldInfo.InitMajor != newInfo.InitMajor ||
		oldInfo.FinalEdu != newInfo.FinalEdu || oldInfo.FinalMajor != newInfo.FinalMajor ||
		oldInfo.GoodSubjects != newInfo.GoodSubjects || oldInfo.BankCard != newInfo.BankCard || oldInfo.Intro != newInfo.Intro
}

func (judgeExperts *EaJudgeExpertsService) DeleteEaJudgeExpertsById(id int64, deleteBy int64) (success bool) {

	// TODO 校验是否有被抽取的记录 不能删除

	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaJudgeExpertsDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

// 导出服务 -- start
type EaJudgeExpertsExportStruct struct {
	DeptMap       map[int64]string
	ExpertTypeMap map[string]string
	TitleMap      map[string]string
	Query         *bizDto.EaJudgeExpertsQuery
	// 专家类型为逗号分隔的多值存储，导出按包含匹配；在入口构造一次，供所有分页复用
	typeCond string
	typeArg  interface{}
}

func (h *EaJudgeExpertsExportStruct) Title() string {
	return "评委专家库导出"
}

func (h *EaJudgeExpertsExportStruct) HeaderColumn() []string {
	if h.Query.CollegeView {
		return []string{"专家ID,expert_id", "姓名,name", "性别,sex", "手机号,phone_no", "证件号,id_card", "专家类型,typeStr", "职称,titleStr", "初始学历,initEduStr", "所学专业,init_major", "最终学历,finalEduStr", "所学专业,final_major", "擅长科目,good_subjects", "银行卡,bank_card", "专家介绍,intro", "专家评价,evaluation"}
	}
	return []string{"专家ID,expert_id", "姓名,name", "证件号,id_card", "性别,sex", "手机号,phone_no", "专家类别,categoryStr", "所属学院/所在单位,deptOrUnit", "专家类型,typeStr", "职称,titleStr", "初始学历,initEduStr", "初始学历所学专业,init_major", "最终学历,finalEduStr", "最终学历所学专业,final_major", "擅长科目,good_subjects", "专家介绍,intro", "审核状态,auditStatusStr", "状态,allowDrawStr", "评价,evaluation"}
}

func (h *EaJudgeExpertsExportStruct) Rows(p *page.Page) []map[string]interface{} {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	if h.Query.BannedOnly {
		conditions = append(conditions, "EXISTS (SELECT 1 FROM ea_judge_expert_ban b WHERE b.expert_id = ea_judge_experts.expert_id AND b.deleted = 1)")
	}
	// 专家类型条件在入口已构造并固定到 typeCond，所有分页复用，避免分页循环时丢失
	if h.typeCond != "" {
		conditions = append(conditions, h.typeCond)
		args = append(args, h.typeArg)
	}
	var condition interface{}
	if len(conditions) > 0 {
		condition = strings.Join(conditions, " AND ")
	}
	result := eaJudgeExpertsDao.ListMapWithCondition(h.Query, p, condition, args...)
	for _, res := range result {
		// fmt.Println(res)
		// 数据库中 intro 字段为 text，map 转出来是 []byte，excelize 不识别而写入空字符串。
		// 统一将 []byte 转换为 string，保证专家介绍等长文本列能正常写出。
		for k, v := range res {
			if b, ok := v.([]byte); ok {
				res[k] = string(b)
			}
		}
		// 学院名称翻译
		if _, ok := res["expert_category"]; ok && mapIntValue(res["expert_category"]) == cst.EXPERT_CATEGORY_EXTERNAL {
			res["categoryStr"] = "校外"
			res["deptOrUnit"] = res["unit_name"]
		} else if _, ok := res["dept_id"]; ok {
			res["categoryStr"] = "校内"
			deptId := res["dept_id"].(int64)
			res["deptOrUnit"] = h.DeptMap[deptId]
		}
		if _, ok := res["type"]; ok {
			labels := make([]string, 0)
			for _, value := range strings.Split(fmt.Sprint(res["type"]), ",") {
				labels = append(labels, h.ExpertTypeMap[value])
			}
			res["typeStr"] = strings.Join(labels, ",")
		}
		if _, ok := res["title"]; ok {
			res["titleStr"] = h.TitleMap[fmt.Sprint(res["title"])]
		}
		if _, ok := res["init_edu"]; ok {
			res["initEduStr"] = cst.GetEduTypeLabelByValue(mapIntValue(res["init_edu"]))
		}
		if _, ok := res["final_edu"]; ok {
			res["finalEduStr"] = cst.GetEduTypeLabelByValue(mapIntValue(res["final_edu"]))
		}
		if _, ok := res["allow_draw"]; ok {
			res["allowDrawStr"] = cst.GetEnableDisableLabelByValue(mapIntValue(res["allow_draw"]))
		}
		if _, ok := res["audit_status"]; ok {
			res["auditStatusStr"] = expertAuditStatusText(mapIntValue(res["audit_status"]))
		}
		if _, ok := res["data_from"]; ok {
			switch mapIntValue(res["data_from"]) {
			case 1:
				res["dataFromStr"] = "学院"
			case 2:
				res["dataFromStr"] = "招办"
			}
		}
	}
	return result
}

func (h *EaJudgeExpertsExportStruct) Finish(url string) {
}

// 导出服务 -- end
func (judgeExperts *EaJudgeExpertsService) ExportEaJudgeExperts(q *bizDto.EaJudgeExpertsQuery) ([]byte, error) {
	// 查询所有学院信息
	deptMap := make(map[int64]string)
	departAll := eaDeptInfoService.ListAllEaDeptInfo(&bizDto.EaDeptInfoQuery{})
	for _, i := range departAll {
		deptMap[i.DeptId] = i.DeptName
	}
	expertTypeMap := make(map[string]string)
	titleMap := make(map[string]string)
	for _, v := range sysDicDataService.ListSysDictData(&userDto.SysDictDataQuery{DictType: cst.BIZ_DICT_EXPERT_TYPE, State: strconv.Itoa(cst.YES)}, &page.Page{CurPage: 1, PageSize: page.MAX_PAGE_SIZE}) {
		expertTypeMap[v.DictValue] = v.DictLabel
	}
	for _, v := range sysDicDataService.ListSysDictData(&userDto.SysDictDataQuery{DictType: cst.BIZ_DICT_EXPERT_TITLE, State: strconv.Itoa(cst.YES)}, &page.Page{CurPage: 1, PageSize: page.MAX_PAGE_SIZE}) {
		titleMap[v.DictValue] = v.DictLabel
	}
	h := &EaJudgeExpertsExportStruct{Query: q, DeptMap: deptMap, ExpertTypeMap: expertTypeMap, TitleMap: titleMap}
	// 专家类型支持多选存储（逗号分隔），导出时按包含匹配，与列表查询保持一致；
	// 在入口构造一次条件并清空 q.Type，避免产生 type = '本科' 的等值条件，也避免分页循环时条件丢失
	if typeFilter := strings.TrimSpace(q.Type); typeFilter != "" {
		h.typeCond = "FIND_IN_SET(?, type) > 0"
		h.typeArg = typeFilter
		q.Type = ""
	}
	return excel.NewExcel().ExportSync(h)
}

func expertAuditStatusText(status int) string {
	switch status {
	case cst.EXPERT_AUDIT_UNSUBMITTED:
		return "未提交"
	case cst.EXPERT_AUDIT_PENDING:
		return "待审核"
	case cst.EXPERT_AUDIT_PASSED:
		return "审核通过"
	case cst.EXPERT_AUDIT_REJECTED:
		return "审核不通过"
	default:
		return ""
	}
}

// 导入服务 -- start
type HandlerImpJudgeExpertsStruct struct {
	BizType         int
	BatchNo         string
	ExistExpertsMap map[string]*bizDo.EaJudgeExperts
	SubjectMap      map[string]string
	ExpertTypeMap   map[string]string
	TitleMap        map[string]string
	DeptMap         map[string]int64
	Query           *bizDto.EaJudgeExpertsQuery
	List            []*bizDo.EaJudgeExperts
	Msg             string
	Err             string
	ExcelHeader     []string
}

func (h *HandlerImpJudgeExpertsStruct) Headers(headers []string) {
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

func (h *HandlerImpJudgeExpertsStruct) Row(mp *map[string]string) excel.RunStatus {
	if h.Err != "" {
		return excel.Exit
	}
	m := *mp
	rowIdx := strconv.Itoa(cmutils.StrToInt(m["rid"]) + 1)
	experts2Insert := new(bizDo.EaJudgeExperts)

	name := strings.TrimSpace(m[h.ExcelHeader[0]])
	sex := strings.TrimSpace(m[h.ExcelHeader[1]])
	idCard := strings.TrimSpace(m[h.ExcelHeader[2]])
	phoneNo := strings.TrimSpace(m[h.ExcelHeader[3]])
	category := "校内"
	deptOrUnit := ""
	offset := 0
	if h.Query.DataFrom == cst.ADDBYADMISSIONS {
		category = strings.TrimSpace(m[h.ExcelHeader[4]])
		deptOrUnit = strings.TrimSpace(m[h.ExcelHeader[5]])
		offset = 2
	}
	title := strings.TrimSpace(m[h.ExcelHeader[4+offset]])
	expertTypes := strings.TrimSpace(m[h.ExcelHeader[5+offset]])
	initEdu := strings.TrimSpace(m[h.ExcelHeader[6+offset]])
	initMajor := strings.TrimSpace(m[h.ExcelHeader[7+offset]])
	finalEdu := strings.TrimSpace(m[h.ExcelHeader[8+offset]])
	finalMajor := strings.TrimSpace(m[h.ExcelHeader[9+offset]])
	goodSubjects := strings.TrimSpace(m[h.ExcelHeader[10+offset]])
	bankCard := strings.TrimSpace(m[h.ExcelHeader[11+offset]])
	intro := strings.TrimSpace(m[h.ExcelHeader[12+offset]])

	experts2Insert.Name = name
	experts2Insert.Sex = sex
	experts2Insert.IdCard = idCard
	experts2Insert.PhoneNo = phoneNo
	experts2Insert.Title = h.TitleMap[title]
	experts2Insert.InitEdu = cst.EduTypeMap[initEdu]
	experts2Insert.InitMajor = initMajor
	if finalEdu != "" {
		experts2Insert.FinalEdu = cst.EduTypeMap[finalEdu]
		experts2Insert.FinalMajor = finalMajor
	}
	experts2Insert.GoodSubjects = goodSubjects
	experts2Insert.BankCard = bankCard
	experts2Insert.Intro = intro

	if category == "校内" {
		experts2Insert.ExpertCategory = cst.EXPERT_CATEGORY_INTERNAL
		if h.Query.DataFrom == cst.ADDBYDEPT {
			experts2Insert.DeptId = h.Query.DeptId
		} else {
			experts2Insert.DeptId = h.DeptMap[deptOrUnit]
		}
	} else if category == "校外" && h.Query.DataFrom == cst.ADDBYADMISSIONS {
		experts2Insert.ExpertCategory = cst.EXPERT_CATEGORY_EXTERNAL
		experts2Insert.DeptId = cst.BIZ_DEPTID_FOR_ADMISSIONS
		experts2Insert.UnitName = deptOrUnit
	} else {
		h.Err = "第" + rowIdx + "行，学院只能导入本学院的校内专家"
		return excel.Exit
	}

	typeValues := make([]string, 0)
	for _, label := range strings.Split(expertTypes, ",") {
		value := h.ExpertTypeMap[strings.TrimSpace(label)]
		if value == "" {
			h.Err = "第" + rowIdx + "行，专家类型字典无法匹配"
			return excel.Exit
		}
		typeValues = append(typeValues, value)
	}
	experts2Insert.Type = strings.Join(typeValues, ",")
	for _, subject := range strings.Split(goodSubjects, ",") {
		if len(h.SubjectMap[strings.TrimSpace(subject)]) <= 0 {
			h.Err = "第" + rowIdx + "行，擅长科目字典无法匹配"
			return excel.Exit
		}
	}
	if title != "" && experts2Insert.Title == "" {
		h.Err = "第" + rowIdx + "行，职称字典无法匹配"
		return excel.Exit
	}
	if initEdu == "" || experts2Insert.InitEdu == 0 || (finalEdu != "" && experts2Insert.FinalEdu == 0) {
		h.Err = "第" + rowIdx + "行，学历字典无法匹配"
		return excel.Exit
	}
	if ok, msg := validateJudgeExpert(experts2Insert); !ok {
		h.Err = "第" + rowIdx + "行，" + msg
		return excel.Exit
	}
	experts2Insert.AllowDraw = cst.YES
	experts2Insert.DataFrom = h.Query.DataFrom
	if h.Query.DataFrom == cst.ADDBYADMISSIONS {
		experts2Insert.AuditStatus = cst.EXPERT_AUDIT_PASSED
	} else {
		experts2Insert.AuditStatus = cst.EXPERT_AUDIT_UNSUBMITTED
	}

	// 数据重复校验
	key := idCard
	if h.ExistExpertsMap[key] != nil {
		h.Err = fmt.Sprintf("第%v行，证件号已存在！[%v]", rowIdx, key)
		return excel.Exit
	}
	h.ExistExpertsMap[key] = experts2Insert

	// 创建人、时间
	experts2Insert.Deleted = common.Normal
	experts2Insert.SetCreateBy(h.Query.OperId)
	experts2Insert.SetUpdateBy(h.Query.OperId)
	// 放入list
	h.List = append(h.List, experts2Insert)
	return excel.Next
}

func (h *HandlerImpJudgeExpertsStruct) After() {
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
		eaJudgeExpertsDao.BatchInsert(gtx, h.List)
		sysDataInoutDao.Insert(gtx, dataInout)
	})
	if !success {
		h.Err = "导入失败，请检查是否与已有数据重复！"
	}
}

// 导入服务  --  end

// 导入
func (judgeExperts *EaJudgeExpertsService) ImportEaJudgeExperts(q *bizDto.EaJudgeExpertsQuery) (success bool, msg string) {

	// 使用雪花算法生成全局唯一id，作为本次导入操作的批次号
	currentBatchNo := strconv.FormatInt(snowflake.GenID(), 10)

	// 查询导入所需字典
	subjectMap := make(map[string]string)
	expertTypeMap := make(map[string]string)
	titleMap := make(map[string]string)
	subjectList := sysDicDataService.ListSysDictData(
		&userDto.SysDictDataQuery{DictType: cst.BIZ_DICT_SUBJECTS, State: strconv.Itoa(cst.YES)},
		&page.Page{CurPage: 1, PageSize: page.MAX_PAGE_SIZE})
	for _, v := range subjectList {
		subjectMap[v.DictLabel] = v.DictValue
	}
	for _, v := range sysDicDataService.ListSysDictData(
		&userDto.SysDictDataQuery{DictType: cst.BIZ_DICT_EXPERT_TYPE, State: strconv.Itoa(cst.YES)},
		&page.Page{CurPage: 1, PageSize: page.MAX_PAGE_SIZE}) {
		expertTypeMap[v.DictLabel] = v.DictValue
	}
	for _, v := range sysDicDataService.ListSysDictData(
		&userDto.SysDictDataQuery{DictType: cst.BIZ_DICT_EXPERT_TITLE, State: strconv.Itoa(cst.YES)},
		&page.Page{CurPage: 1, PageSize: page.MAX_PAGE_SIZE}) {
		titleMap[v.DictLabel] = v.DictValue
	}
	deptMap := make(map[string]int64)
	for id, dept := range eaDeptInfoService.GetAllDeptMap() {
		deptMap[dept.DeptName] = id
	}

	// 查询所有已有的专家库放入map
	expertsAll := judgeExperts.ListAllEaJudgeExperts(new(bizDto.EaJudgeExpertsQuery))
	dbExpertMap := make(map[string]*bizDo.EaJudgeExperts)
	for _, v := range expertsAll {
		dbExpertMap[v.IdCard] = v
	}

	// 遍历文件map  文件名->文件地址...
	excelHeader := []string{"姓名", "性别", "证件号", "手机号", "专家类别", "所属学院/所在单位", "职称", "专家类型", "初始学历", "初始学历所学专业", "最终学历", "最终学历所学专业", "擅长科目", "银行卡", "专家介绍"}
	if q.DataFrom == cst.ADDBYDEPT {
		excelHeader = []string{"姓名", "性别", "证件号", "手机号", "职称", "专家类型", "初始学历", "初始学历所学专业", "最终学历", "最终学历所学专业", "擅长科目", "银行卡", "专家介绍"}
	}
	handlerImpStruct := &HandlerImpJudgeExpertsStruct{BizType: int(cst.BIZ_JUDGE_EXPERTS), BatchNo: currentBatchNo,
		Query: q, SubjectMap: subjectMap, ExpertTypeMap: expertTypeMap, TitleMap: titleMap,
		DeptMap: deptMap, ExistExpertsMap: dbExpertMap, ExcelHeader: excelHeader}
	for _, v := range q.Files {
		// 根据url导入excel处理
		excel.NewExcel().ImportWithUrl(v, handlerImpStruct)
	}
	if len(handlerImpStruct.Err) > 0 {
		return false, handlerImpStruct.Err
	}
	return true, handlerImpStruct.Msg
}

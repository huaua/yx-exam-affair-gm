package bizSerImpl

import (
	"encoding/json"
	"strconv"
	"strings"
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

	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/ds"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/utils/file/excel"
	"github.com/sealsee/web-base/public/utils/snowflake"
	"gorm.io/gorm"
)

/*
 * desc: 教职工服务层实现
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchInfoService struct{}

var eaTchInfoService *EaTchInfoService
var eaTchInfoDao bizDao.IEaTchInfoDao

func init() {
	eaTchInfoService = &EaTchInfoService{}
	eaTchInfoDao = bizDaoImpl.NewEaTchInfoDao()
}

func NewEaTchInfoService() *EaTchInfoService {
	return eaTchInfoService
}

func (tchInfo *EaTchInfoService) GetEaTchInfoById(id int64) *bizDo.EaTchInfo {
	return eaTchInfoDao.GetById(id)
}

// GetNextPendingTchInfo returns the next enabled pending teacher for continuous admissions audit.
// 查询条件 q 来自列表页的检索条件，保证连续审核范围与当前筛选结果一致；
// q.TchId 表示刚审核完的那条记录，用于优先向后取（tch_id > 当前值），取不到时再从头回绕。
func (tchInfo *EaTchInfoService) GetNextPendingTchInfo(q *bizDto.EaTchInfoQuery) *bizDo.EaTchInfo {
	if q == nil {
		q = new(bizDto.EaTchInfoQuery)
	}
	currentId := q.TchId

	buildQuery := func() *gorm.DB {
		db := ds.GetGDB().Model(&bizDo.EaTchInfo{}).
			Where("audit_status = ? AND enabled_status = ? AND deleted = ?",
				cst.TCH_AUDIT_PENDING, cst.TCH_ENABLED, common.Normal)
		if name := strings.TrimSpace(q.TchName); name != "" {
			db = db.Where("tch_name LIKE ?", "%"+name+"%")
		}
		if jobNo := strings.TrimSpace(q.JobNo); jobNo != "" {
			db = db.Where("job_no LIKE ?", "%"+jobNo+"%")
		}
		if idCard := strings.TrimSpace(q.IdCard); idCard != "" {
			db = db.Where("id_card LIKE ?", "%"+idCard+"%")
		}
		if sex := strings.TrimSpace(q.Sex); sex != "" {
			db = db.Where("sex = ?", sex)
		}
		if duties := strings.TrimSpace(q.Duties); duties != "" {
			db = db.Where("duties LIKE ?", "%"+duties+"%")
		}
		if profOffice := strings.TrimSpace(q.ProfOffice); profOffice != "" {
			db = db.Where("prof_office = ?", profOffice)
		}
		if education := strings.TrimSpace(q.Education); education != "" {
			db = db.Where("education = ?", education)
		}
		if q.DeptId != 0 {
			db = db.Where("dept_id = ?", q.DeptId)
		}
		if q.InvigilationExperience > 0 {
			db = db.Where("invigilation_experience = ?", q.InvigilationExperience)
		}
		return db
	}

	data := new(bizDo.EaTchInfo)
	var res *gorm.DB
	if currentId > 0 {
		// 先取当前记录之后的一条，保持审核顺序与列表顺序一致
		res = buildQuery().Where("tch_id > ?", currentId).Order("tch_id asc").First(data)
		if res.Error != nil || res.RowsAffected < 1 {
			// 后面没有了，回绕到满足条件的第一条（跳过自身）
			data = new(bizDo.EaTchInfo)
			res = buildQuery().Where("tch_id <> ?", currentId).Order("tch_id asc").First(data)
		}
	} else {
		res = buildQuery().Order("tch_id asc").First(data)
	}
	if res.Error != nil || res.RowsAffected < 1 {
		return nil
	}
	if cst.BIZ_DEPTID_FOR_ADMISSIONS == data.DeptId {
		data.DeptName = "公共教室"
	} else if dept := eaDeptInfoService.GetAllDeptMap()[data.DeptId]; dept != nil {
		data.DeptName = dept.DeptName
	}
	return data
}

func (tchInfo *EaTchInfoService) GetLatestAuditHistory(tchId int64) *bizDo.EaTchInfoAuditHistory {
	return eaTchInfoDao.GetLatestAuditHistory(tchId)
}

func (tchInfo *EaTchInfoService) ListAuditHistory(tchId int64) []*bizDo.EaTchInfoAuditHistory {
	return eaTchInfoDao.ListAuditHistory(tchId)
}

func (tchInfo *EaTchInfoService) SubmitForAudit(tchId int64, deptId int64, submitBy int64) (bool, string) {
	if tchId <= 0 {
		return false, "参数错误，请检查数据"
	}
	if deptId <= 0 {
		return false, "仅学院用户可执行提交操作"
	}
	dbTchInfo := eaTchInfoDao.GetById(tchId)
	if dbTchInfo == nil {
		return false, "教职工信息不存在"
	}
	if dbTchInfo.DeptId != deptId {
		return false, "仅本学院的教职工数据可提交"
	}
	if dbTchInfo.AuditStatus != cst.TCH_AUDIT_UNSUBMITTED {
		return false, "仅未提交状态的教职工数据可提交至待审核"
	}
	uWhere := &bizDo.EaTchInfo{TchId: tchId}
	uData := &bizDo.EaTchInfo{
		AuditStatus: cst.TCH_AUDIT_PENDING,
	}
	uData.SetUpdateBy(submitBy)
	success := false
	tx.ExecGTx(func(gtx tx.Db) {
		success = eaTchInfoDao.UpdateNew(gtx, uWhere, uData) > 0
	})
	if !success {
		return false, "提交操作失败，请稍后再试"
	}
	return true, ""
}

// 学院批量提交教职工至待审核：ids 非空时提交指定记录，ids 为空时提交本学院全部未提交数据
func (tchInfo *EaTchInfoService) SubmitForAuditBatch(ids []int64, deptId int64, submitBy int64) (int, string) {
	if deptId <= 0 {
		return 0, "仅学院用户可执行提交操作"
	}
	affected := 0
	tx.ExecGTx(func(gtx tx.Db) {
		targetIds := ids
		if len(targetIds) == 0 {
			var allIds []int64
			ds.GetGDB().Raw("SELECT tch_id FROM ea_tch_info WHERE dept_id = ? AND audit_status = ? AND deleted = ?",
				deptId, cst.TCH_AUDIT_UNSUBMITTED, common.Normal).Scan(&allIds)
			targetIds = allIds
		}
		for _, tchId := range targetIds {
			if tchId <= 0 {
				continue
			}
			dbInfo := eaTchInfoDao.GetById(tchId)
			if dbInfo == nil || dbInfo.DeptId != deptId || dbInfo.AuditStatus != cst.TCH_AUDIT_UNSUBMITTED {
				continue
			}
			uWhere := &bizDo.EaTchInfo{TchId: tchId}
			uData := &bizDo.EaTchInfo{AuditStatus: cst.TCH_AUDIT_PENDING}
			uData.SetUpdateBy(submitBy)
			if eaTchInfoDao.UpdateNew(gtx, uWhere, uData) > 0 {
				affected++
			}
		}
	})
	return affected, ""
}

func (tchInfo *EaTchInfoService) ListInvigilationRecords(tchId int64) []*bizDo.EaTchArrangeRecord {
	if tchId < 1 {
		return nil
	}
	sql := `select MIN(d.tch_arrange_id) as tch_arrange_id, MIN(d.tch_report_id) as tch_report_id, d.tch_task_id, MAX(t.tch_task_name) as tch_task_name,
				d.room_report_id, MAX(r.room_name) as room_name,
				GROUP_CONCAT(DISTINCT r.arrange_no_str ORDER BY r.arrange_no_str SEPARATOR '/') as arrange_no_str,
				CAST(IFNULL(t.task_start_time, '1970-01-01 00:00:00') AS DATETIME) as task_start_time,
				CAST(IFNULL(t.task_start_time, '1970-01-01 00:00:00') AS DATETIME) as task_time
			from ea_arrange_tch d
			left join ea_tch_task_info t on t.tch_task_id = d.tch_task_id
			left join ea_arrange_room  r on r.room_report_id = d.room_report_id
			where d.deleted = 1 and d.tch_id = ?
			group by d.tch_task_id, d.room_report_id, t.task_start_time
			order by task_start_time desc, MIN(d.create_time) desc`
	list := make([]*bizDo.EaTchArrangeRecord, 0)
	res := ds.GetGDB().Raw(sql, tchId).Scan(&list)
	if res.Error != nil {
		panic(res.Error)
	}
	return list
}

func (tchInfo *EaTchInfoService) AuditEaTchInfo(tchId int64, action string, reason string, auditBy int64) (bool, string) {
	if tchId <= 0 {
		return false, "参数错误，请检查数据"
	}
	if action != cst.TCH_AUDIT_ACTION_PASS && action != cst.TCH_AUDIT_ACTION_REJECT {
		return false, "审核操作不合法"
	}
	dbTchInfo := eaTchInfoDao.GetById(tchId)
	if dbTchInfo == nil {
		return false, "教职工信息不存在"
	}
	if dbTchInfo.AuditStatus != cst.TCH_AUDIT_PENDING {
		return false, "仅待审核状态的教职工支持审核操作"
	}
	if action == cst.TCH_AUDIT_ACTION_REJECT {
		trimmed := strings.TrimSpace(reason)
		if len(trimmed) <= 0 {
			return false, "审核不通过时必须填写不通过原因"
		}
	}

	targetStatus := cst.TCH_AUDIT_PASSED
	if action == cst.TCH_AUDIT_ACTION_REJECT {
		targetStatus = cst.TCH_AUDIT_REJECTED
	}
	uWhere := &bizDo.EaTchInfo{TchId: tchId}
	uData := &bizDo.EaTchInfo{
		AuditStatus: targetStatus,
		AuditRemark: strings.TrimSpace(reason),
	}
	uData.SetUpdateBy(auditBy)
	success := false
	tx.ExecGTx(func(gtx tx.Db) {
		success = eaTchInfoDao.UpdateNew(gtx, uWhere, uData) > 0
	})
	if !success {
		return false, "审核操作失败，请稍后再试"
	}
	return true, ""
}

func (tchInfo *EaTchInfoService) ChangeEnabledStatus(tchId int64, enabledStatus int, optBy int64) (bool, string) {
	if tchId <= 0 {
		return false, "参数错误，请检查数据"
	}
	if enabledStatus != cst.TCH_ENABLED && enabledStatus != cst.TCH_DISABLED {
		return false, "启用状态参数错误"
	}
	dbTchInfo := eaTchInfoDao.GetById(tchId)
	if dbTchInfo == nil {
		return false, "教职工信息不存在"
	}
	if dbTchInfo.EnabledStatus == enabledStatus {
		return true, ""
	}
	uWhere := &bizDo.EaTchInfo{TchId: tchId}
	uData := &bizDo.EaTchInfo{EnabledStatus: enabledStatus}
	uData.SetUpdateBy(optBy)
	success := false
	tx.ExecGTx(func(gtx tx.Db) {
		success = eaTchInfoDao.UpdateNew(gtx, uWhere, uData) > 0
	})
	if !success {
		return false, "更新启用状态失败"
	}
	return true, ""
}

// ChangeEnabledStatusBatch 批量切换教职工的启用/禁用状态（招办）
func (tchInfo *EaTchInfoService) ChangeEnabledStatusBatch(tchIds []int64, enabledStatus int, optBy int64) (int, string) {
	if len(tchIds) == 0 {
		return 0, "请选择教职工"
	}
	if enabledStatus != cst.TCH_ENABLED && enabledStatus != cst.TCH_DISABLED {
		return 0, "启用状态参数错误"
	}
	affected := 0
	tx.ExecGTx(func(gtx tx.Db) {
		for _, tchId := range tchIds {
			if tchId <= 0 {
				continue
			}
			uWhere := &bizDo.EaTchInfo{TchId: tchId}
			uData := &bizDo.EaTchInfo{EnabledStatus: enabledStatus}
			uData.SetUpdateBy(optBy)
			affected += eaTchInfoDao.UpdateNew(gtx, uWhere, uData)
		}
	})
	return affected, ""
}

func (tchInfo *EaTchInfoService) ListEaTchInfo(q *bizDto.EaTchInfoQuery, p *page.Page) []*bizDo.EaTchInfo {
	q.AddLikeAll("tch_name", q.TchName)
	q.AddLikeAll("job_no", q.JobNo)
	list := eaTchInfoDao.List(q, p)
	if len(list) <= 0 {
		return list
	}

	deptMap := eaDeptInfoService.GetAllDeptMap()
	for _, i := range list {
		if cst.BIZ_DEPTID_FOR_ADMISSIONS == i.DeptId {
			i.DeptName = "公共教室"
		} else {
			if deptMap[i.DeptId] != nil {
				i.DeptName = deptMap[i.DeptId].DeptName
			}
		}
	}
	return list
}

func (tchInfo *EaTchInfoService) ListAllEaTchInfo(q *bizDto.EaTchInfoQuery, condition interface{}, args ...interface{}) []*bizDo.EaTchInfo {
	listAll := make([]*bizDo.EaTchInfo, 0)
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE
	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		list := eaTchInfoDao.ListWithCondition(q, pageInfo, condition, args...)
		if len(list) == 0 {
			break
		}
		listAll = append(listAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}
	return listAll
}

func (tchInfo *EaTchInfoService) InsertEaTchInfo(tchInfoDO *bizDo.EaTchInfo) (success bool, errInfo string) {
	if tchInfoDO.AuditStatus == 0 {
		// 招办新增的数据直接置为审核通过
		if tchInfoDO.DateFrom == cst.ADDBYADMISSIONS {
			tchInfoDO.AuditStatus = cst.TCH_AUDIT_PASSED
		} else {
			tchInfoDO.AuditStatus = cst.TCH_AUDIT_UNSUBMITTED
		}
	}
	if tchInfoDO.InvigilationExperience == 0 {
		tchInfoDO.InvigilationExperience = cst.TCH_NO_INVIGILATION_EXPERIENCE
	}
	if tchInfoDO.EnabledStatus == 0 {
		tchInfoDO.EnabledStatus = cst.TCH_ENABLED
	}

	success, errInfo = validateTchInfo(tchInfoDO)
	if !success {
		return
	}

	// 工号限制
	if eaTchInfoDao.Count(&bizDto.EaTchInfoQuery{JobNo: tchInfoDO.JobNo}) > 0 {
		return false, "工号已存在，请检查数据"
	}

	// 证件号限制
	if eaTchInfoDao.Count(&bizDto.EaTchInfoQuery{IdCard: tchInfoDO.IdCard}) > 0 {
		return false, "证件号已存在，请检查数据"
	}

	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaTchInfoDao.Insert(gtx, tchInfoDO)
	})
	return
}

/*
 * 字段长度限制
 * author: 蓝鲸
 * date：2024/09/04
 */
func validateTchInfo(tchInfoDO *bizDo.EaTchInfo) (success bool, errInfo string) {
	// 必填字段字段
	if len(tchInfoDO.TchName) <= 0 || len(tchInfoDO.Sex) <= 0 || len(tchInfoDO.JobNo) <= 0 ||
		len(tchInfoDO.IdCard) <= 0 || tchInfoDO.DeptId < 0 {
		return false, "姓名，性别，工号，证件号，所属学院为必填信息，请检查数据"
	}

	// 姓名限制
	if utf8.RuneCountInString(tchInfoDO.TchName) < 2 || utf8.RuneCountInString(tchInfoDO.TchName) > 30 {
		return false, "姓名限制2-30个字符"
	}

	// 性别校验
	if cst.MAN != tchInfoDO.Sex && cst.WOMAN != tchInfoDO.Sex {
		return false, "性别只能是男或女"
	}

	// 工号限制
	if utf8.RuneCountInString(tchInfoDO.JobNo) < 2 || utf8.RuneCountInString(tchInfoDO.JobNo) > 20 {
		return false, "工号限制2-20个字符"
	}
	if !cmutils.IsNumberAndLetter(tchInfoDO.JobNo) {
		return false, "工号只允许输入字母和数字"
	}

	// 证件号限制
	if utf8.RuneCountInString(tchInfoDO.IdCard) > 32 {
		return false, "证件号长度过长"
	}
	if !cmutils.IsNumberAndLetter(tchInfoDO.IdCard) {
		return false, "证件号只允许输入字母和数字"
	}

	// 职务限制
	if len(tchInfoDO.Duties) > 0 {
		if utf8.RuneCountInString(tchInfoDO.Duties) < 2 || utf8.RuneCountInString(tchInfoDO.Duties) > 15 {
			return false, "职务限制2-15个字符"
		}
	}

	// 银行卡限制
	if len(tchInfoDO.Bankcard) > 0 {
		if utf8.RuneCountInString(tchInfoDO.Bankcard) < 16 || utf8.RuneCountInString(tchInfoDO.Bankcard) > 19 {
			return false, "银行卡限制16-19个字符"
		}
		if !cmutils.IsNumberAndLetter(tchInfoDO.Bankcard) {
			return false, "银行卡只允许输入字母和数字"
		}
	}

	// 专业或行政限制
	if len(tchInfoDO.ProfOffice) > 0 {
		if utf8.RuneCountInString(tchInfoDO.ProfOffice) < 2 || utf8.RuneCountInString(tchInfoDO.ProfOffice) > 15 {
			return false, "专业或行政限制2-15个字符"
		}
	}

	// 专业背景限制
	if len(tchInfoDO.Education) > 0 {
		if utf8.RuneCountInString(tchInfoDO.Education) > 1000 {
			return false, "专业背景限制1000个字符"
		}
	}
	return true, ""
}

func (tchInfo *EaTchInfoService) UpdateEaTchInfoById(tchInfoDO *bizDo.EaTchInfo, deptId int64) (success bool, errInfo string) {
	if tchInfoDO.TchId <= 0 {
		return false, "参数错误，请检查数据"
	}
	dbTchInfo := eaTchInfoDao.GetById(tchInfoDO.TchId)
	if dbTchInfo == nil {
		return false, "参数错误，请检查数据"
	}

	isCollege := deptId > 0
	if isCollege {
		if dbTchInfo.DeptId != deptId {
			return false, "无权编辑其他学院的教职工"
		}
		if dbTchInfo.EnabledStatus != cst.TCH_ENABLED {
			return false, "该教职工已禁用，无法编辑"
		}
		tchInfoDO.DeptId = deptId
	} else {
		// 招办：禁止编辑学院未提交的教职工
		if dbTchInfo.AuditStatus == cst.TCH_AUDIT_UNSUBMITTED {
			return false, "学院未提交的教职工数据，招办不可编辑"
		}
		if tchInfoDO.DeptId <= 0 {
			tchInfoDO.DeptId = dbTchInfo.DeptId
		}
	}

	if dbTchInfo.IdCard != tchInfoDO.IdCard {
		return false, "证件号不允许修改"
	}

	success, errInfo = validateTchInfo(tchInfoDO)
	if !success {
		return
	}

	if dbTchInfo.JobNo != tchInfoDO.JobNo && eaTchInfoDao.Count(&bizDto.EaTchInfoQuery{JobNo: tchInfoDO.JobNo}) > 0 {
		return false, "工号已存在，请检查数据"
	}

	contentChanged := hasTchInfoContentChanged(dbTchInfo, tchInfoDO)
	statusChanged := false
	if isCollege {
		switch tchInfoDO.SubmitAction {
		case cst.TCH_SUBMIT_ACTION_DRAFT:
			// 暂存时若内容未变化且当前非待审核状态（未提交/审核不通过），不允许保存；
			// 待审核状态下点暂存允许不修改内容，回退为未提交
			if !contentChanged && dbTchInfo.AuditStatus != cst.TCH_AUDIT_PENDING {
				return false, "未进行信息修改"
			}
			tchInfoDO.AuditStatus = cst.TCH_AUDIT_UNSUBMITTED
		case cst.TCH_SUBMIT_ACTION_SUBMIT:
			// 未提交状态提交审核：允许不修改内容，直接更新为待审核
			if dbTchInfo.AuditStatus != cst.TCH_AUDIT_UNSUBMITTED {
				if !contentChanged {
					return false, "未进行信息修改"
				}
			}
			tchInfoDO.AuditStatus = cst.TCH_AUDIT_PENDING
		default:
			return false, "请选择暂存或提交操作"
		}
		statusChanged = dbTchInfo.AuditStatus != tchInfoDO.AuditStatus
	} else {
		// 招办编辑后默认重置为审核通过，并清空不通过原因
		tchInfoDO.AuditStatus = cst.TCH_AUDIT_PASSED
		tchInfoDO.AuditRemark = ""
		statusChanged = dbTchInfo.AuditStatus != tchInfoDO.AuditStatus || dbTchInfo.AuditRemark != tchInfoDO.AuditRemark
	}

	if !contentChanged && !statusChanged {
		return true, ""
	}

	uWhere := &bizDo.EaTchInfo{TchId: tchInfoDO.TchId}
	tchInfoDO.TchId = 0
	tchInfoDO.SetNullableCols("duties", "bankcard", "prof_office", "education")
	tx.ExecGTx(func(gtx tx.Db) {
		if dbTchInfo.AuditStatus == cst.TCH_AUDIT_PASSED && (contentChanged || statusChanged) {
			snapshot, err := json.Marshal(dbTchInfo)
			if err != nil {
				success = false
				return
			}
			history := &bizDo.EaTchInfoAuditHistory{
				TchId:        dbTchInfo.TchId,
				DeptId:       dbTchInfo.DeptId,
				AuditStatus:  dbTchInfo.AuditStatus,
				SnapshotData: string(snapshot),
				ChangeBy:     tchInfoDO.UpdateBy,
				ChangeTime:   time.Now(),
			}
			if result := gtx.Create(history); result.Error != nil || result.RowsAffected < 1 {
				success = false
				return
			}
		}
		success = eaTchInfoDao.UpdateNew(gtx, uWhere, tchInfoDO) > 0
	})
	return
}

func hasTchInfoContentChanged(oldInfo, newInfo *bizDo.EaTchInfo) bool {
	return oldInfo.TchName != newInfo.TchName ||
		oldInfo.Sex != newInfo.Sex ||
		oldInfo.JobNo != newInfo.JobNo ||
		oldInfo.DeptId != newInfo.DeptId ||
		oldInfo.Duties != newInfo.Duties ||
		oldInfo.Bankcard != newInfo.Bankcard ||
		oldInfo.ProfOffice != newInfo.ProfOffice ||
		oldInfo.Education != newInfo.Education
}

func (tchInfo *EaTchInfoService) DeleteEaTchInfoById(id int64, deleteBy int64, deptId int64) (success bool, errMsg errs.ERROR) {
	success = false
	dbTchInfo := eaTchInfoDao.GetById(id)
	if dbTchInfo == nil {
		errMsg = errs.PARAM_INVALID_ERR
		return
	}
	// 招办新增和导入的学院不可被学院编辑和删除
	if dbTchInfo.DateFrom == cst.ADDBYADMISSIONS && deptId > 0 {
		errMsg = hint.BIZ_TCH_CAN_NOT_EDIT_ERR
		return
	}

	success = tx.ExecGTx(func(gtx tx.GTx) {
		eaTchInfoDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

// 导出服务 -- start
type EaTchInfoExportStruct struct {
	Query *bizDto.EaTchInfoQuery
}

func (h *EaTchInfoExportStruct) Title() string {
	return "教职工导出"
}

func (h *EaTchInfoExportStruct) HeaderColumn() []string {
	if h.Query.CollegeMode {
		return []string{"姓名,tch_name", "性别,sex", "工号,job_no", "证件号,id_card", "职务,duties", "工商银行卡号,bankcard", "专业或行政,prof_office", "专业背景,education", "审核状态,audit_status_str", "状态,enabled_status_str"}
	}
	// 招办导出：使用所属学院名称（不显示ID），含工商银行卡号、审核状态、启用/禁用状态
	return []string{"姓名,tch_name", "性别,sex", "工号,job_no", "证件号,id_card", "所属学院,dept_name_str", "职务,duties", "工商银行卡号,bankcard", "专业或行政,prof_office", "专业背景,education", "审核状态,audit_status_str", "状态,enabled_status_str"}
}

func (h *EaTchInfoExportStruct) Rows(p *page.Page) []map[string]interface{} {
	rows := eaTchInfoDao.ListMapWithCondition(h.Query, p, nil, nil)
	deptMap := eaDeptInfoService.GetAllDeptMap()
	for _, row := range rows {
		// 审核状态文本
		row["audit_status_str"] = auditStatusText(mapIntValue(row["audit_status"]))
		// 招办导出：所属学院id转名称
		deptId := mapInt64Value(row["dept_id"])
		if deptMap[deptId] != nil {
			row["dept_name_str"] = deptMap[deptId].DeptName
		} else if deptId == cst.BIZ_DEPTID_FOR_ADMISSIONS {
			row["dept_name_str"] = "公共教室"
		} else {
			row["dept_name_str"] = ""
		}
		// 启用状态
		enabled := mapIntValue(row["enabled_status"])
		if enabled == cst.TCH_ENABLED {
			row["enabled_status_str"] = "启用"
		} else {
			row["enabled_status_str"] = "禁用"
		}
		// 监考经验
		exp := mapIntValue(row["invigilation_experience"])
		if exp == cst.TCH_HAS_INVIGILATION_EXPERIENCE {
			row["invigilation_experience_str"] = "有"
		} else {
			row["invigilation_experience_str"] = "无"
		}
	}
	return rows
}

func (h *EaTchInfoExportStruct) Finish(url string) {
}

// 导出服务 -- end
func (tchInfo *EaTchInfoService) ExportEaTchInfo(q *bizDto.EaTchInfoQuery) ([]byte, error) {
	q.AddLikeAll("tch_name", q.TchName)
	q.AddLikeAll("job_no", q.JobNo)
	return excel.NewExcel().ExportSync(&EaTchInfoExportStruct{q})
}

func mapInt64Value(value interface{}) int64 {
	switch v := value.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint8:
		return int64(v)
	case uint64:
		return int64(v)
	case []byte:
		n, _ := strconv.ParseInt(string(v), 10, 64)
		return n
	case string:
		n, _ := strconv.ParseInt(v, 10, 64)
		return n
	default:
		return 0
	}
}

func mapIntValue(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint8:
		return int(v)
	case uint64:
		return int(v)
	case []byte:
		n, _ := strconv.Atoi(string(v))
		return n
	case string:
		n, _ := strconv.Atoi(v)
		return n
	default:
		return 0
	}
}

func auditStatusText(status int) string {
	switch status {
	case cst.TCH_AUDIT_PENDING:
		return "待审核"
	case cst.TCH_AUDIT_PASSED:
		return "审核通过"
	case cst.TCH_AUDIT_REJECTED:
		return "审核不通过"
	default:
		return "未提交"
	}
}

func (tchInfo *EaTchInfoService) ImportEaTchInfo(q *bizDto.EaTchInfoQuery) (success bool, errInfo string) {
	// 使用雪花算法生成全局唯一id，作为本次导入操作的批次号
	currentBatchNo := strconv.FormatInt(snowflake.GenID(), 10)

	// 查询所有已有的教职工信息放入map
	classroomInfos := tchInfo.ListAllEaTchInfo(new(bizDto.EaTchInfoQuery), nil, nil)
	dbIdCardMap := make(map[string]*bizDo.EaTchInfo)
	dbJobNoMap := make(map[string]*bizDo.EaTchInfo)
	for _, v := range classroomInfos {
		dbIdCardMap[v.IdCard] = v
		dbJobNoMap[v.JobNo] = v
	}

	deptMap := make(map[string]*bizDo.EaDeptInfo)
	if !q.CollegeMode {
		departAll := eaDeptInfoService.ListAllEaDeptInfo(&bizDto.EaDeptInfoQuery{})
		for _, i := range departAll {
			deptMap[i.DeptName] = i
		}
	}

	excelHeader := []string{"姓名", "性别", "工号", "证件号", "所属学院", "职务", "工商银行卡号", "专业或行政", "专业背景"}
	if q.CollegeMode {
		excelHeader = []string{"姓名", "性别", "工号", "证件号", "职务", "工商银行卡号", "专业或行政", "专业背景"}
	}
	handlerImpStruct := &HandlerImpTchInfoStruct{BizType: int(cst.BIZ_ROOM_BASE_INFO), BatchNo: currentBatchNo,
		Query: q, ExistIdCardMap: dbIdCardMap, ExistJobNoMap: dbJobNoMap, DeptNameMap: deptMap, ExcelHeader: excelHeader}
	for _, v := range q.Files {
		// 根据url导入excel处理
		excel.NewExcel().ImportWithUrl(v, handlerImpStruct)
	}
	if len(handlerImpStruct.Err) > 0 {
		return false, handlerImpStruct.Err
	}
	return true, handlerImpStruct.Msg
}

// 导入服务 -- start
type HandlerImpTchInfoStruct struct {
	BizType        int
	BatchNo        string
	ExistIdCardMap map[string]*bizDo.EaTchInfo
	ExistJobNoMap  map[string]*bizDo.EaTchInfo
	DeptNameMap    map[string]*bizDo.EaDeptInfo
	Query          *bizDto.EaTchInfoQuery
	List           []*bizDo.EaTchInfo
	Msg            string
	Err            string
	ExcelHeader    []string
}

func (h *HandlerImpTchInfoStruct) Headers(headers []string) {
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

func (h *HandlerImpTchInfoStruct) Row(mp *map[string]string) excel.RunStatus {
	// 判断表头数据是否错误
	if h.Err != "" {
		return excel.Exit
	}
	m := *mp
	rowIdx := strconv.Itoa(cmutils.StrToInt(m["rid"]) + 1)
	tchInfo2Insert := new(bizDo.EaTchInfo)
	allBlank := true
	for _, header := range h.ExcelHeader {
		if strings.TrimSpace(m[header]) != "" {
			allBlank = false
			break
		}
	}
	if allBlank {
		return excel.Next
	}

	// 取值
	tchInfo2Insert.TchName = strings.TrimSpace(m[h.ExcelHeader[0]])
	tchInfo2Insert.Sex = strings.TrimSpace(m[h.ExcelHeader[1]])
	tchInfo2Insert.JobNo = strings.TrimSpace(m[h.ExcelHeader[2]])
	tchInfo2Insert.IdCard = strings.TrimSpace(m[h.ExcelHeader[3]])
	if h.Query.CollegeMode {
		tchInfo2Insert.DeptId = h.Query.DeptId
		tchInfo2Insert.Duties = strings.TrimSpace(m[h.ExcelHeader[4]])
		tchInfo2Insert.Bankcard = strings.TrimSpace(m[h.ExcelHeader[5]])
		tchInfo2Insert.ProfOffice = strings.TrimSpace(m[h.ExcelHeader[6]])
		tchInfo2Insert.Education = strings.TrimSpace(m[h.ExcelHeader[7]])
	} else {
		deptName := strings.TrimSpace(m[h.ExcelHeader[4]])
		tchInfo2Insert.Duties = strings.TrimSpace(m[h.ExcelHeader[5]])
		tchInfo2Insert.Bankcard = strings.TrimSpace(m[h.ExcelHeader[6]])
		tchInfo2Insert.ProfOffice = strings.TrimSpace(m[h.ExcelHeader[7]])
		tchInfo2Insert.Education = strings.TrimSpace(m[h.ExcelHeader[8]])

		if len(deptName) <= 0 {
			h.Err = "第" + rowIdx + "行，所属学院不能为空！"
			return excel.Exit
		}
		if h.DeptNameMap[deptName] == nil {
			h.Err = "第" + rowIdx + "行，所属学院不存在！"
			return excel.Exit
		}
		tchInfo2Insert.DeptId = h.DeptNameMap[deptName].DeptId
	}

	success, errInfo := validateTchInfo(tchInfo2Insert)
	if !success {
		h.Err = "第" + rowIdx + "行，" + errInfo
		return excel.Exit
	}

	// 数据重复校验
	if h.ExistJobNoMap[tchInfo2Insert.JobNo] != nil {
		h.Err = "第" + rowIdx + "行的工号已经存在！"
		return excel.Exit
	}
	h.ExistJobNoMap[tchInfo2Insert.JobNo] = tchInfo2Insert

	if h.ExistIdCardMap[tchInfo2Insert.IdCard] != nil {
		h.Err = "第" + rowIdx + "行的身份证号已经存在！"
		return excel.Exit
	}
	h.ExistIdCardMap[tchInfo2Insert.IdCard] = tchInfo2Insert

	// 创建人、时间
	tchInfo2Insert.DateFrom = cst.ADDBYADMISSIONS
	if h.Query.CollegeMode {
		tchInfo2Insert.DateFrom = cst.ADDBYDEPT
	}
	tchInfo2Insert.AuditStatus = cst.TCH_AUDIT_UNSUBMITTED
	tchInfo2Insert.InvigilationExperience = cst.TCH_NO_INVIGILATION_EXPERIENCE
	tchInfo2Insert.EnabledStatus = cst.TCH_ENABLED
	tchInfo2Insert.Deleted = common.Normal
	tchInfo2Insert.SetCreateBy(h.Query.OperId)
	tchInfo2Insert.SetUpdateBy(h.Query.OperId)
	// 放入list
	h.List = append(h.List, tchInfo2Insert)
	return excel.Next
}

func (h *HandlerImpTchInfoStruct) After() {
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
		eaTchInfoDao.BatchInsert(gtx, h.List)
		sysDataInoutDao.Insert(gtx, dataInout)
	})
	if !success {
		h.Err = "导入失败，请检查是否与已有数据重复！"
	}
}

// 导入服务  --  end

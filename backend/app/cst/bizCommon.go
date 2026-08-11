// Package cst
/**
 * @note:
 *
 * @author: 大猫
 * @date: 2023/6/18
 */
package cst

const (
	USER_TYPE_SYS = "00" // 系统用户

	BIZ_DEPTID_FOR_ADMISSIONS = -100 // 招办部门id

	BIZ_LEVEL string = "biz_level" // 字典-层级

	BIZ_DICT_SUBJECTS string = "biz_dict_subjects" // 字典-科目
	BIZ_DICT_EXPERT_TYPE string = "biz_expert_type" // 字典-专家类型
	BIZ_DICT_EXPERT_TITLE string = "biz_expert_title" // 字典-专家职称

	ALL string = "全部"
)

// config key
const (
	// 考试列表获取接口地址
	SYS_CONFIG_SYNC_EXAM_URL string = "sync_sch_exams_url"
	// 数据同步接口地址
	SYS_CONFIG_SYNC_APPLY_URL string = "sync_stu_profs_url"
	// 数据同步接口token
	SYS_CONFIG_SYNC_TOKEN string = "sync_stu_prof_token"
	// 数据同步接口每页处理记录数 1000以内
	SYS_CONFIG_SYNC_PAGE_SIZE string = "sync_stu_prof_page_size"
)

// redis key 数据同步状态 0-未开始 1-进行中 2-已完成 3-异常终止
const (
	RDS_SYNC_EXAM_KEY  string = "sync_key_exam_list"      // 院校考试列表
	RDS_SYNC_STATE_KEY string = "sync_key_stu_prof_state" // 数据同步缓存key
	RDS_SYNC_STATE_NOT string = "0"                       // 数据同步未开始
	RDS_SYNC_STATE_ING string = "1"                       // 数据同步进行中
	RDS_SYNC_STATE_OK  string = "2"                       // 数据同步已完成
	RDS_SYNC_STATE_ERR string = "3"                       // 数据同步异常终止
)

// 请求类型 POST/GET
const (
	POST string = "POST"
	GET  string = "GET"
)

// 是否常量
const (
	YES = 1
	NO  = 2
)

var YesNoMap = map[string]int{
	"是": 1,
	"否": 2,
}

func GetYesNoLabelByValue(val int) string {
	for label, v := range YesNoMap {
		if v == val {
			return label
		}
	}
	return ""
}

// 启用禁用常量
var EnableDisableMap = map[string]int{
	"启用": 1,
	"禁用": 2,
}

func GetEnableDisableLabelByValue(val int) string {
	for label, v := range EnableDisableMap {
		if v == val {
			return label
		}
	}
	return ""
}

// 是否常量-字符串
const (
	YES1 = "1"
	NO2  = "2"
)

// 是否发布
const (
	// 未发布
	UNPUBLISH = "1"
	// 已发布
	PUBLISH = "2"
)

// 是否确认
const (
	// 未确认
	UNCONFIRM = "1"
	// 确认
	CONFIRM = "2"
)

const (
	MAN   = "男"
	WOMAN = "女"
)

const (
	// 暂存
	TEMP_SAVE = "A"
	// 教师任务分配 征集状态 1-已上报,待审核（学院上报后状态）
	TCH_REPORT_PENDING_AUDIT = "1"
	// 教师任务分配 征集状态 2-待分配（招办审核通过）
	TCH_REPORT_PENDING_ASSIGN = "2"
	// 教师任务分配 征集状态 3-不通过（招办审核不通过）
	TCH_REPORT_REJECTED = "3"
	// 教师任务分配 征集状态 4-已分配（已分配到考场）
	TCH_REPORT_ASSIGNED = "4"
)

// 历史常量保留：用于 ea_arrange_room.assignment_state 等历史字段
// 值与新征集状态语义可能重合，但作用于不同表，请勿混用
const (
	// 已上报（ea_arrange_room.assignment_state: 未分配）
	ESCALATED = "1"
	// 已征集（ea_arrange_room.assignment_state: 已分配）
	ASSIGNED = "2"
)

const (
	// 学院新增
	ADDBYDEPT = 1
	// 招办新增
	ADDBYADMISSIONS = 2
)

const (
	TCH_AUDIT_UNSUBMITTED = 1
	TCH_AUDIT_PENDING     = 2
	TCH_AUDIT_PASSED      = 3
	TCH_AUDIT_REJECTED    = 4

	TCH_ENABLED  = 1
	TCH_DISABLED = 2

	TCH_HAS_INVIGILATION_EXPERIENCE = 1
	TCH_NO_INVIGILATION_EXPERIENCE  = 2

	TCH_SUBMIT_ACTION_DRAFT  = "draft"
	TCH_SUBMIT_ACTION_SUBMIT = "submit"

	// 招办审核操作
	TCH_AUDIT_ACTION_PASS   = "pass"
	TCH_AUDIT_ACTION_REJECT = "reject"
)

const (
	EXPERT_CATEGORY_INTERNAL = 1
	EXPERT_CATEGORY_EXTERNAL = 2

	EXPERT_AUDIT_UNSUBMITTED = 1
	EXPERT_AUDIT_PENDING     = 2
	EXPERT_AUDIT_PASSED      = 3
	EXPERT_AUDIT_REJECTED    = 4

	EXPERT_SUBMIT_ACTION_DRAFT  = "draft"
	EXPERT_SUBMIT_ACTION_SUBMIT = "submit"

	EXPERT_AUDIT_ACTION_PASS   = "pass"
	EXPERT_AUDIT_ACTION_REJECT = "reject"
)

const (
	// 未编排
	UNARRANGE = "1"
	// 已编排
	ARRANGED = "2"
)

const (
	// 未分配
	UNASSIGNMENT = "1"
	// 已分配
	ASSIGNMENT = "2"
)

package bizDto

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 评委专家库数据查询实体
 * author: init code
 * date: 2025/06/06
 */
type EaJudgeExpertsQuery struct {
	ExpertId     int64             `gorm:"primaryKey" form:"expertId" json:"expertId,string,omitempty"` // 主键 专家id
	Name         string            `form:"name" json:"name,omitempty"`                                  // 姓名
	Sex          string            `form:"sex" json:"sex,omitempty"`                                    // 性别（男或女）
	PhoneNo      string            `form:"phoneNo" json:"phoneNo,omitempty"`                            // 手机号
	IdCard       string            `form:"idCard" json:"idCard,omitempty"`                              // 证件号
	DeptId       int64             `form:"deptId" json:"deptId,string,omitempty"`                       // 专家来源-所属学院id（招办的属校外）
	Type         string            `form:"type" json:"type,omitempty"`                                  // 专家类型，逗号分隔字典值
	ExpertCategory int             `form:"expertCategory" json:"expertCategory,string,omitempty"`       // 专家类别
	UnitName     string            `form:"unitName" json:"unitName,omitempty"`                          // 所在单位
	Title        string            `form:"title" json:"title,omitempty"`                                // 职称
	InitEdu      int               `form:"initEdu" json:"initEdu,string,omitempty"`                     // 初始学历
	InitMajor    string            `form:"initMajor" json:"initMajor,omitempty"`                        // 初始学历所学专业
	FinalEdu     int               `form:"finalEdu" json:"finalEdu,string,omitempty"`                   // 最终学历
	FinalMajor   string            `form:"finalMajor" json:"finalMajor,omitempty"`                      // 最终学历所学专业
	GoodSubjects string            `form:"goodSubjects" json:"goodSubjects,omitempty"`                  // 擅长科目 字典(1-素描 2-速写 3-色彩) 多选拼接值
	BankCard     string            `form:"bankCard" json:"bankCard,omitempty"`                          // 银行卡号
	Intro        string            `form:"intro" json:"intro,omitempty"`                                // 专家介绍
	AllowDraw    int               `form:"allowDraw" json:"allowDraw,string,omitempty"`                 // 是否允许抽取 1-允许 2-禁止
	DataFrom     int               `form:"dataFrom" json:"dataFrom,string,omitempty"`                   // 数据来源 1-学院新增 2-招办新增
	AuditStatus  int               `form:"auditStatus" json:"auditStatus,string,omitempty"`             // 审核状态
	AuditRemark  string            `form:"auditRemark" json:"auditRemark,omitempty"`                    // 审核原因
	Evaluation   string            `form:"evaluation" json:"evaluation,omitempty"`                      // 评价
	IsEvaluated  int               `gorm:"-" form:"isEvaluated" json:"isEvaluated,string,omitempty"`   // 是否评价 1-是 2-否
	Remark       string            `form:"remark" json:"remark,omitempty"`                              // 备注
	Files        map[string]string `gorm:"-" form:"files" json:"files,omitempty"`                       // 导入文件 map[文件名]地址
	OperId       int64             `gorm:"-" form:"operId" json:"operId,omitempty"`                     // 操作人userid
	BizType      int               `gorm:"-" form:"bizType" json:"bizType,omitempty"`                   // 业务类型
	CollegeView  bool              `gorm:"-" form:"-" json:"-"`                                       // 学院端导出格式标记
	BannedOnly   bool              `gorm:"-" form:"-" json:"-"`                                       // 仅查询禁止抽取名单
	ExcludeBanned bool             `gorm:"-" form:"-" json:"-"`                                       // 排除禁止抽取名单
	basemodel.BaseEntityQuery
}

type EaJudgeExpertsAuditRequest struct {
	ExpertId int64  `json:"expertId,string,omitempty"`
	Action   string `json:"action,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

type EaJudgeExpertsEvaluationRequest struct {
	ExpertId   int64  `json:"expertId,string,omitempty"`
	Evaluation string `json:"evaluation,omitempty"`
}

type EaJudgeExpertsStatusRequest struct {
	ExpertId int64 `json:"expertId,string,omitempty"`
	AllowDraw int  `json:"allowDraw,string,omitempty"`
}

type EaJudgeExpertsIdsRequest struct {
	ExpertIds []int64 `json:"expertIds,omitempty"`
}

package bizDo

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 评委专家库数据实体
 * author: init code
 * date: 2025/06/06
 */
type EaJudgeExperts struct {
	ExpertId     int64  `gorm:"primaryKey" json:"expertId,string,omitempty"` // 主键 专家id
	Name         string `json:"name,omitempty"`                              // 姓名
	Sex          string `json:"sex,omitempty"`                               // 性别（男或女）
	PhoneNo      string `json:"phoneNo,omitempty"`                           // 手机号
	IdCard       string `json:"idCard,omitempty"`                            // 证件号
	DeptId       int64  `json:"deptId,string,omitempty"`                     // 专家来源-所属学院id（招办的属校外）
	Type         string `gorm:"type:varchar(64)" json:"type,omitempty"`     // 专家类型，多选，逗号分隔字典值
	ExpertCategory int  `json:"expertCategory,string,omitempty"`             // 专家类别 1-校内 2-校外
	UnitName     string `json:"unitName,omitempty"`                          // 校外专家所在单位
	Title        string `json:"title,omitempty"`                             // 职称（字典值）
	InitEdu      int    `json:"initEdu,string,omitempty"`                    // 初始学历
	InitMajor    string `json:"initMajor,omitempty"`                         // 初始学历所学专业
	FinalEdu     int    `json:"finalEdu,string,omitempty"`                   // 最终学历
	FinalMajor   string `json:"finalMajor,omitempty"`                        // 最终学历所学专业
	GoodSubjects string `json:"goodSubjects,omitempty"`                      // 擅长科目 字典(1-素描 2-速写 3-色彩) 多选拼接值
	BankCard     string `json:"bankCard,omitempty"`                          // 银行卡号
	Intro        string `json:"intro,omitempty"`                             // 专家介绍
	AllowDraw    int    `json:"allowDraw,string,omitempty"`                  // 是否允许抽取 1-允许 2-禁止
	DataFrom     int    `json:"dataFrom,string,omitempty"`                   // 数据来源 1-学院新增 2-招办新增
	AuditStatus  int    `gorm:"default:1" json:"auditStatus,string,omitempty"` // 审核状态
	AuditRemark  string `json:"auditRemark,omitempty"`                      // 审核不通过原因
	Evaluation   string `gorm:"type:varchar(2000)" json:"evaluation,omitempty"` // 招办评价
	SubmitAction string `gorm:"-" json:"submitAction,omitempty"`           // 学院操作 draft/submit
	DeptName     string `gorm:"-" json:"deptName,omitempty"`               // 所属学院名称
	HasEvaluation bool  `gorm:"-" json:"hasEvaluation,omitempty"`          // 是否已评价
	Remark       string `json:"remark,omitempty"`                            // 备注
	basemodel.BaseEntity
}

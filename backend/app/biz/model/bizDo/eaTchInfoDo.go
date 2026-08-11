package bizDo

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 教职工数据实体
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchInfo struct {
	TchId                  int64                `gorm:"primaryKey" json:"tchId,string,omitempty"`                 // 主键 教职工ID
	TchName                string               `json:"tchName,omitempty"`                                        // 姓名
	Sex                    string               `json:"sex,omitempty"`                                            // 性别（男或女）
	JobNo                  string               `json:"jobNo,omitempty"`                                          // 工号
	IdCard                 string               `json:"idCard,omitempty"`                                         // 证件号
	DeptId                 int64                `json:"deptId,string,omitempty"`                                  // 所属学院id
	Duties                 string               `json:"duties,omitempty"`                                         // 职务
	Bankcard               string               `json:"bankcard,omitempty"`                                       // 工商银行卡号
	ProfOffice             string               `json:"profOffice,omitempty"`                                     // 专业或行政
	Education              string               `json:"education,omitempty"`                                      // 专业背景
	AuditStatus            int                  `gorm:"default:1" json:"auditStatus,string,omitempty"`            // 审核状态 1-未提交 2-待审核 3-审核通过 4-审核不通过
	AuditRemark            string               `json:"auditRemark,omitempty"`                                    // 审核不通过原因
	InvigilationExperience int                  `gorm:"default:2" json:"invigilationExperience,string,omitempty"` // 监考经验 1-有 2-无
	EnabledStatus          int                  `gorm:"default:1" json:"enabledStatus,string,omitempty"`          // 启用状态 1-启用 2-禁用
	DateFrom               int                  `json:"dateFrom,omitempty"`                                       // 数据来源 1-学院新增 2-招办新增
	Remark                 string               `json:"remark,omitempty"`                                         // 备注
	SubmitAction           string               `gorm:"-" json:"submitAction,omitempty"`                          // 学院操作 draft-暂存 submit-提交
	ReportTaskTimeList     []basemodel.BaseTime `gorm:"-" json:"reportTchList,omitempty"`                         // 已上报征集详情列表
	DeptName               string               `gorm:"-" json:"deptName,omitempty"`                              // 学院名称
	basemodel.BaseEntity
}

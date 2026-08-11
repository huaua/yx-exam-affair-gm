package bizDto

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 教职工数据查询实体
 * author: 蓝鲸
 * date: 2024/09/04
 */
type EaTchInfoQuery struct {
	TchId                  int64             `gorm:"primaryKey" form:"tchId" json:"tchId,string,omitempty"`                 // 主键 教职工ID
	TchName                string            `form:"tchName" json:"tchName,omitempty"`                                      // 姓名
	Sex                    string            `form:"sex" json:"sex,omitempty"`                                              // 性别（男或女）
	JobNo                  string            `form:"jobNo" json:"jobNo,omitempty"`                                          // 工号
	IdCard                 string            `form:"idCard" json:"idCard,omitempty"`                                        // 证件号
	DeptId                 int64             `form:"deptId" json:"deptId,string,omitempty"`                                 // 所属学院id
	Duties                 string            `form:"duties" json:"duties,omitempty"`                                        // 职务
	Bankcard               string            `form:"bankcard" json:"bankcard,omitempty"`                                    // 工商银行卡号
	ProfOffice             string            `form:"profOffice" json:"profOffice,omitempty"`                                // 专业或行政
	Education              string            `form:"education" json:"education,omitempty"`                                  // 专业背景
	AuditStatus            int               `form:"auditStatus" json:"auditStatus,string,omitempty"`                       // 审核状态 1-未提交 2-待审核 3-审核通过 4-审核不通过
	InvigilationExperience int               `form:"invigilationExperience" json:"invigilationExperience,string,omitempty"` // 监考经验 1-有 2-无
	EnabledStatus          int               `form:"enabledStatus" json:"enabledStatus,string,omitempty"`                   // 启用状态 1-启用 2-禁用
	Remark                 string            `form:"remark" json:"remark,omitempty"`                                        // 备注
	Files                  map[string]string `gorm:"-" form:"files" json:"files,omitempty"`                                 // 导入文件 map[文件名]地址
	OperId                 int64             `gorm:"-" form:"operId" json:"operId,omitempty"`                               // 操作人userid
	BizType                int               `gorm:"-" form:"bizType" json:"bizType,omitempty"`                             //
	CollegeMode            bool              `gorm:"-" form:"collegeMode" json:"-"`                                         // 是否学院端请求
	basemodel.BaseEntityQuery
}

// 教职工审核请求
type EaTchInfoAuditReq struct {
	TchId  int64  `json:"tchId,string,omitempty"`
	Action string `json:"action,omitempty"` // pass-通过 reject-不通过
	Reason string `json:"reason,omitempty"` // 审核不通过原因（不通过时必填）
}

// 教职工启用/禁用切换请求
type EaTchInfoChangeEnabledReq struct {
	TchId         int64 `json:"tchId,string,omitempty"`
	EnabledStatus int   `json:"enabledStatus,string,omitempty"`
}

// 教职工批量启用/禁用切换请求
type EaTchInfoChangeEnabledBatchReq struct {
	TchIds       []int64 `json:"tchIds"`
	EnabledStatus int    `json:"enabledStatus,string,omitempty"`
}

// 教职工监考记录/审核历史通用请求
type EaTchInfoArrangeRecordReq struct {
	TchId int64 `json:"tchId,string,omitempty"`
}

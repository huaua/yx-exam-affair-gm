package syncDomain

import (
	"github.com/sealsee/web-base/public/basemodel"
	"github.com/sealsee/web-base/public/ds/tx"
)

// 用于openapi接口返回数据对象转换
type EaSyncStuProfDto struct {
	BaoKaoID          int64               `gorm:"primaryKey;column:BaoKaoID" json:"baoKaoID,omitempty"`        // 主键 报考专业ID
	YongHuID          int64               `gorm:"column:YongHuID" json:"yongHuID,omitempty"`                   // 用户ID
	KaoShengID        int64               `gorm:"column:KaoShengID" json:"kaoShengID,omitempty"`               // 考生ID
	ZhengJianLX       int                 `gorm:"column:ZhengJianLX" json:"zhengJianLX,omitempty"`             // 证件类型
	ShenFenZH         string              `gorm:"column:ShenFenZH" json:"shenFenZH,omitempty"`                 // 身份证号
	XingMing          string              `gorm:"column:XingMing" json:"xingMing,omitempty"`                   // 姓名
	ZhunKaoZH         string              `gorm:"column:ZhunKaoZH" json:"zhunKaoZH,omitempty"`                 // 准考证号
	XueXiaoID         int64               `gorm:"column:XueXiaoID" json:"xueXiaoID,omitempty"`                 // 学校ID
	XueXiaoMC         string              `gorm:"column:XueXiaoMC" json:"xueXiaoMC,omitempty"`                 // 学校名称
	KaoShiID          int64               `gorm:"column:KaoShiID" json:"kaoShiID,omitempty"`                   // 考试ID
	KaoShiMC          string              `gorm:"column:KaoShiMC" json:"kaoShiMC,omitempty"`                   // 考试名称
	KaoDianID         int64               `gorm:"column:KaoDianID" json:"kaoDianID,omitempty"`                 // 考点ID
	KaoDianMC         string              `gorm:"column:KaoDianMC" json:"kaoDianMC,omitempty"`                 // 考点名称
	ZhuanYeID         int64               `gorm:"column:ZhuanYeID" json:"zhuanYeID,omitempty"`                 // 专业ID
	ZhuanYeMC         string              `gorm:"column:ZhuanYeMC" json:"zhuanYeMC,omitempty"`                 // 专业名称
	RiChengID         int64               `gorm:"column:RiChengID" json:"riChengID,omitempty"`                 // 日程ID
	FuRiChengID       int64               `gorm:"column:fuRiChengID" json:"fuRiChengID,omitempty"`             // 父日程ID
	KaoShiRQ          basemodel.TimeStamp `gorm:"column:KaoShiRQ" json:"kaoShiRQ,omitempty"`                   // 考试日期
	KaoShiRQSM        string              `gorm:"column:KaoShiRQSM" json:"kaoShiRQSM,omitempty"`               // 考试日期说明
	KaoChangMC        string              `gorm:"column:KaoChangMC" json:"kaoChangMC,omitempty"`               // 考场名称
	GaoKaoSF          string              `gorm:"column:gaoKaoSF" json:"gaoKaoSF,omitempty"`                   // 高考省份
	BaoMingFei        float64             `gorm:"column:BaoMingFei" json:"baoMingFei,omitempty"`               // 报名费
	BaoKaoBZ          int                 `gorm:"column:BaoKaoBZ" json:"baoKaoBZ,omitempty"`                   // 报考标志 1-未提交 2-已提交 3-已生效 4-已关闭 5-作废
	QueRenFS          int                 `gorm:"column:QueRenFS" json:"queRenFS,omitempty"`                   // 确认方式 1-线上确认 2-线下确认
	QueRenSJ          basemodel.TimeStamp `gorm:"column:QueRenSJ" json:"queRenSJ,omitempty"`                   // 确认时间
	ShengFenMC        string              `gorm:"column:ShengFenMC" json:"shengFenMC,omitempty"`               // 高考省份名称
	XingBie           string              `gorm:"column:XingBie" json:"xingBie,omitempty"`                     // 性别
	TongXinDZ         string              `gorm:"column:TongXinDZ" json:"tongXinDZ,omitempty"`                 // 通信地址
	TongXinYB         string              `gorm:"column:TongXinYB" json:"tongXinYB,omitempty"`                 // 通信邮编
	ShouJi            string              `gorm:"column:ShouJi" json:"shouJi,omitempty"`                       // 本人手机
	KaoShengHao       string              `gorm:"column:KaoShengHao" json:"kaoShengHao,omitempty"`             // 考生号
	AuditFlag         int                 `gorm:"column:auditFlag" json:"auditFlag,omitempty"`                 // 申请审核状态 1-申请中 2-审核不通过 3-作废完成
	AuditDes          string              `gorm:"column:auditDes" json:"auditDes,omitempty"`                   // 审核意见
	PrintTime         basemodel.TimeStamp `gorm:"column:printTime" json:"printTime,omitempty"`                 // 准考证打印时间
	WenLiKe           string              `gorm:"column:WenLiKe" json:"wenLiKe,omitempty"`                     // 文理科
	VideoCommitFlag   int                 `gorm:"column:videoCommitFlag" json:"videoCommitFlag,omitempty"`     // 网考视频提交状态 0-待提交 1-已提交
	OtherPlatFlag     int                 `gorm:"column:otherPlatFlag" json:"otherPlatFlag,omitempty"`         // 报考数据来源：1-外部报考数据 2-缴费资格库
	FinishedExamTime  basemodel.TimeStamp `gorm:"column:finishedExamTime" json:"finishedExamTime,omitempty"`   // 考试结束时间
	Year              int                 `gorm:"column:year" json:"year,omitempty"`                           //	年份
	SyncTime          basemodel.BaseTime  `gorm:"column:syncTime" json:"syncTime,omitempty"`                   //	同步时间
	RenZhengZP        string              `gorm:"column:renZhengZP" json:"renZhengZP,omitempty"`               // 认证照片
	Gatlx             int                 `gorm:"column:gatlx" json:"gatlx,omitempty"`                         // 港澳台类型（证件类型为港澳台时必填）
	ZhengJianLXDesc   string              `gorm:"column:zhengJianLXDesc" json:"zhengJianLXDesc,omitempty"`     // 证件类型名称
	ChuShengRQ        string              `gorm:"column:chuShengRQ" json:"chuShengRQ,omitempty"`               // 出生日期
	YingWangJie       string              `gorm:"column:yingWangJie" json:"yingWangJie,omitempty"`             // 应往届
	MinZu             string              `gorm:"column:minZu" json:"minZu,omitempty"`                         // 民族
	XueLi             string              `gorm:"column:xueLi" json:"xueLi,omitempty"`                         // 学历
	ZhengZhiMM        string              `gorm:"column:zhengZhiMM" json:"zhengZhiMM,omitempty"`               // 政治面貌
	SingleSubjectName string              `gorm:"column:singleSubjectName" json:"singleSubjectName,omitempty"` // 首选科目
	SuoZaiHS          string              `gorm:"column:suoZaiHS" json:"suoZaiHS,omitempty"`                   // 专业课学习学校
	StuTypeStr        string              `gorm:"column:stuTypeStr" json:"stuTypeStr,omitempty"`               //考生类型：1-小学生 2-初中生 3-高中生 4-中专生 5-大学在校生
	SuoZaiXX          string              `gorm:"column:suoZaiXX" json:"suoZaiXX,omitempty"`                   // 文化课学习学校
	SubjectChoosen    string              `gorm:"column:subjectChoosen" json:"subjectChoosen,omitempty"`       // 选考科目
	SubjectReChoose   string              `gorm:"column:subjectReChoose" json:"subjectReChoose,omitempty"`     // 再选科目

	basemodel.BaseEntity
}

func (dto *EaSyncStuProfDto) TableName() string {
	return "ea_sync_stu_prof"
}

// 注意：拼装list里deleted要设为1正常
func (dto *EaSyncStuProfDto) BatchInsert(gtx tx.GTx, dtoList []*EaSyncStuProfDto) bool {
	res := gtx.CreateInBatches(dtoList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 数据同步进度信息
type SyncStateInfoCache struct {
	State     string `form:"state" json:"state,omitempty"`                // 当前状态 0-未开始 1-进行中 2-已完成 3-异常终止
	KaoShiID  int64  `form:"kaoShiID" json:"kaoShiID,string,omitempty"`   // 考试ID
	KaoShiMC  string `form:"kaoShiMC" json:"kaoShiMC,omitempty"`          // 考试名称
	TotalSize int    `form:"totalSize" json:"totalSize,string,omitempty"` // 总记录数
	TotalPage int    `form:"totalPage" json:"totalPage,string,omitempty"` // 总页数
	PageSize  int    `form:"pageSize" json:"pageSize,string,omitempty"`   // 每页数量
	CurPage   int    `form:"curPage" json:"curPage,string,omitempty"`     // 当前页
	Percent   string `form:"percent" json:"percent,omitempty"`            // 百分比
	ErrMsg    string `form:"errMsg" json:"errMsg,omitempty"`              // 异常信息
}

// 考试信息
type SyncExamInfoCache struct {
	KaoShiID int64  `form:"kaoShiID" json:"kaoShiID,omitempty"` // 考试ID
	KaoShiMC string `form:"kaoShiMC" json:"kaoShiMC,omitempty"` // 考试名称
	KaoShiND string `form:"kaoShiND" json:"kaoShiND,omitempty"` // 考试年度
}

type LDatas struct {
	List []*SyncExamInfoCache `json:"list"`
}

type PDatas struct {
	Page PageDatas `json:"page"`
}

type PageDatas struct {
	CurPage   int                 `json:"curPage"`
	PageSize  int                 `json:"pageSize"`
	TotalSize int                 `json:"totalSize"`
	TotalPage int                 `json:"totalPage"`
	DataList  []*EaSyncStuProfDto `json:"dataList"`
}

type PageResult struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Datas   PDatas `json:"datas"`
}

type ListResult struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Datas   LDatas `json:"datas"`
}

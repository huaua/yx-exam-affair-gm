package syncDomain

import (
	"github.com/sealsee/web-base/public/basemodel"
	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/ds"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/query"
	"github.com/sealsee/web-base/public/ds/tx"
)

/*
 * desc: 数据同步-考生报考专业数据
 * author: liyu
 * date: 2025/06/13
 */
type EaSyncStuProf struct {
	BaoKaoID          int64   `gorm:"primaryKey;column:BaoKaoID" json:"baoKaoID,string,omitempty"`    // 主键 报考专业ID
	YongHuID          int64   `gorm:"column:YongHuID" json:"yongHuID,string,omitempty"`               // 用户ID
	KaoShengID        int64   `gorm:"column:KaoShengID" json:"kaoShengID,string,omitempty"`           // 考生ID
	ZhengJianLX       int     `gorm:"column:ZhengJianLX" json:"zhengJianLX,string,omitempty"`         // 证件类型
	ShenFenZH         string  `gorm:"column:ShenFenZH" json:"shenFenZH,omitempty"`                    // 身份证号
	XingMing          string  `gorm:"column:XingMing" json:"xingMing,omitempty"`                      // 姓名
	ZhunKaoZH         string  `gorm:"column:ZhunKaoZH" json:"zhunKaoZH,omitempty"`                    // 准考证号
	XueXiaoID         int64   `gorm:"column:XueXiaoID" json:"xueXiaoID,string,omitempty"`             // 学校ID
	XueXiaoMC         string  `gorm:"column:XueXiaoMC" json:"xueXiaoMC,omitempty"`                    // 学校名称
	KaoShiID          int64   `gorm:"column:KaoShiID" json:"kaoShiID,string,omitempty"`               // 考试ID
	KaoShiMC          string  `gorm:"column:KaoShiMC" json:"kaoShiMC,omitempty"`                      // 考试名称
	KaoDianID         int64   `gorm:"column:KaoDianID" json:"kaoDianID,string,omitempty"`             // 考点ID
	KaoDianMC         string  `gorm:"column:KaoDianMC" json:"kaoDianMC,omitempty"`                    // 考点名称
	ZhuanYeID         int64   `gorm:"column:ZhuanYeID" json:"zhuanYeID,string,omitempty"`             // 专业ID
	ZhuanYeMC         string  `gorm:"column:ZhuanYeMC" json:"zhuanYeMC,omitempty"`                    // 专业名称
	RiChengID         int64   `gorm:"column:RiChengID" json:"riChengID,string,omitempty"`             // 日程ID
	FuRiChengID       int64   `gorm:"column:fuRiChengID" json:"fuRiChengID,string,omitempty"`         // 父日程ID
	KaoShiRQ          string  `gorm:"column:KaoShiRQ" json:"kaoShiRQ,omitempty"`                      // 考试日期
	KaoShiRQSM        string  `gorm:"column:KaoShiRQSM" json:"kaoShiRQSM,omitempty"`                  // 考试日期说明
	KaoChangMC        string  `gorm:"column:KaoChangMC" json:"KaoChangMC,string,omitempty"`           // 考场名称
	GaoKaoSF          string  `gorm:"column:gaoKaoSF" json:"gaoKaoSF,omitempty"`                      // 高考省份
	BaoMingFei        float64 `gorm:"column:BaoMingFei" json:"baoMingFei,string,omitempty"`           // 报名费
	BaoKaoBZ          int     `gorm:"column:BaoKaoBZ" json:"baoKaoBZ,string,omitempty"`               // 报考标志 1-未提交 2-已提交 3-已生效 4-已关闭 5-作废
	QueRenFS          int     `gorm:"column:QueRenFS" json:"queRenFS,string,omitempty"`               // 确认方式 1-线上确认 2-线下确认
	QueRenSJ          string  `gorm:"column:QueRenSJ" json:"queRenSJ,omitempty"`                      // 确认时间
	ShengFenMC        string  `gorm:"column:ShengFenMC" json:"shengFenMC,omitempty"`                  // 高考省份名称
	XingBie           string  `gorm:"column:XingBie" json:"xingBie,omitempty"`                        // 性别
	TongXinDZ         string  `gorm:"column:TongXinDZ" json:"tongXinDZ,omitempty"`                    // 通信地址
	TongXinYB         string  `gorm:"column:TongXinYB" json:"tongXinYB,omitempty"`                    // 通信邮编
	ShouJi            string  `gorm:"column:ShouJi" json:"shouJi,omitempty"`                          // 本人手机
	KaoShengHao       string  `gorm:"column:KaoShengHao" json:"kaoShengHao,omitempty"`                // 考生号
	AuditFlag         int     `gorm:"column:auditFlag" json:"auditFlag,string,omitempty"`             // 申请审核状态 1-申请中 2-审核不通过 3-作废完成
	AuditDes          string  `gorm:"column:auditDes" json:"auditDes,omitempty"`                      // 审核意见
	PrintTime         string  `gorm:"column:printTime" json:"printTime,omitempty"`                    // 准考证打印时间
	WenLiKe           string  `gorm:"column:WenLiKe" json:"wenLiKe,omitempty"`                        // 文理科
	VideoCommitFlag   int     `gorm:"column:videoCommitFlag" json:"videoCommitFlag,string,omitempty"` // 网考视频提交状态 0-待提交 1-已提交
	OtherPlatFlag     int     `gorm:"column:otherPlatFlag" json:"otherPlatFlag,string,omitempty"`     // 报考数据来源：1-外部报考数据 2-缴费资格库
	FinishedExamTime  string  `gorm:"column:finishedExamTime" json:"finishedExamTime,omitempty"`      // 考试结束时间
	Year              int     `gorm:"column:year" json:"year,string,omitempty"`                       //	年份
	SyncTime          string  `gorm:"column:syncTime" json:"syncTime,omitempty"`                      //	同步时间
	RenZhengZP        string  `gorm:"column:renZhengZP" json:"renZhengZP,omitempty"`                  // 认证照片
	Gatlx             int     `gorm:"column:gatlx" json:"gatlx,omitempty"`                            // 港澳台类型（证件类型为港澳台时必填）
	ZhengJianLXDesc   string  `gorm:"column:zhengJianLXDesc" json:"zhengJianLXDesc,omitempty"`        // 证件类型名称
	ChuShengRQ        string  `gorm:"column:chuShengRQ" json:"chuShengRQ,omitempty"`                  // 出生日期
	YingWangJie       string  `gorm:"column:yingWangJie" json:"yingWangJie,omitempty"`                // 应往届
	MinZu             string  `gorm:"column:minZu" json:"minZu,omitempty"`                            // 民族
	XueLi             string  `gorm:"column:xueLi" json:"xueLi,omitempty"`                            // 学历
	ZhengZhiMM        string  `gorm:"column:zhengZhiMM" json:"zhengZhiMM,omitempty"`                  // 政治面貌
	SingleSubjectName string  `gorm:"column:singleSubjectName" json:"singleSubjectName,omitempty"`    // 首选科目
	SuoZaiHS          string  `gorm:"column:suoZaiHS" json:"suoZaiHS,omitempty"`                      // 专业课学习学校
	StuTypeStr        string  `gorm:"column:stuTypeStr" json:"stuTypeStr,omitempty"`                  //考生类型：1-小学生 2-初中生 3-高中生 4-中专生 5-大学在校生
	SuoZaiXX          string  `gorm:"column:suoZaiXX" json:"suoZaiXX,omitempty"`                      // 文化课学习学校
	SubjectChoosen    string  `gorm:"column:subjectChoosen" json:"subjectChoosen,omitempty"`          // 选考科目
	SubjectReChoose   string  `gorm:"column:subjectReChoose" json:"subjectReChoose,omitempty"`        // 再选科目
	basemodel.BaseEntity
}

type EaSyncStuProfQuery struct {
	BaoKaoID  int64  `gorm:"primaryKey;column:BaoKaoID" json:"baoKaoID,string,omitempty"` // 报考专业ID
	KaoShiID  int64  `gorm:"column:KaoShiID" json:"kaoShiID,string,omitempty"`            // 考试ID
	KaoShiMC  string `gorm:"column:KaoShiMC" json:"kaoShiMC,omitempty"`                   // 考试名称
	KaoDianID int64  `gorm:"column:KaoDianID" json:"kaoDianID,string,omitempty"`          // 考点ID
	ZhuanYeID int64  `gorm:"column:ZhuanYeID" json:"zhuanYeID,string,omitempty"`          // 专业ID
	RiChengID int64  `gorm:"column:RiChengID" json:"riChengID,string,omitempty"`          // 日程ID
	YongHuID  int64  `gorm:"column:YongHuID" json:"yongHuID,string,omitempty"`            // 用户ID
	ShenFenZH string `gorm:"column:ShenFenZH" json:"shenFenZH,omitempty"`                 // 身份证号
	OperId    int64  `gorm:"-" form:"operId" json:"operId,omitempty"`                     // 操作人userid
	ReSync    bool   `gorm:"-" form:"reSync" json:"reSync,omitempty"`                     // 是否重新同步
	basemodel.BaseEntityQuery
}

func (syncAppStuprof *EaSyncStuProf) GetById(id int64) *EaSyncStuProf {
	if id < 1 {
		return nil
	}
	data := new(EaSyncStuProf)
	where := new(EaSyncStuProf)
	where.BaoKaoID = id
	where.Deleted = common.Normal
	res := ds.GetGDB().Find(data, where)
	if res.Error != nil {
		panic(res.Error)
	}
	if res.RowsAffected < 1 {
		return nil
	}
	return data
}

func (syncAppStuprof *EaSyncStuProf) Count(q *EaSyncStuProfQuery) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCount[EaSyncStuProfQuery, EaSyncStuProf](q)
}

func (syncAppStuprof *EaSyncStuProf) CountWithCondition(q *EaSyncStuProfQuery, condition interface{}, args ...interface{}) int {
	q.Deleted = common.Normal
	return query.ExecGetQueryCountWithCondition[EaSyncStuProf](q, condition, args...)
}

func (syncAppStuprof *EaSyncStuProf) List(q *EaSyncStuProfQuery, p *page.Page) []*EaSyncStuProf {
	q.Deleted = common.Normal
	return query.ExecQueryList[EaSyncStuProfQuery, EaSyncStuProf](q, p)
}

func (syncAppStuprof *EaSyncStuProf) ListWithCondition(q *EaSyncStuProfQuery, p *page.Page, condition interface{}, args ...interface{}) []*EaSyncStuProf {
	q.Deleted = common.Normal
	return query.ExecQueryListWithCondition[EaSyncStuProf](q, p, condition, args...)
}

func (syncAppStuprof *EaSyncStuProf) ListMapWithCondition(q *EaSyncStuProfQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any {
	q.Deleted = common.Normal
	return query.ExecQueryListMapWithCondition[EaSyncStuProf](q, p, condition, args...)
}

func (syncAppStuprof *EaSyncStuProf) Insert(gtx tx.GTx, syncStuProfDO *EaSyncStuProf) bool {
	syncStuProfDO.Deleted = common.Normal
	res := gtx.Create(syncStuProfDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (syncAppStuprof *EaSyncStuProf) BatchInsert(gtx tx.GTx, syncStuProfList []*EaSyncStuProf) bool {
	res := gtx.CreateInBatches(syncStuProfList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (syncAppStuprof *EaSyncStuProf) Update(gtx tx.GTx, uData *EaSyncStuProf, uWhere *EaSyncStuProf) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (syncAppStuprof *EaSyncStuProf) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(EaSyncStuProf)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(EaSyncStuProf)
	uWhere.BaoKaoID = id
	uWhere.Deleted = common.Normal
	return syncAppStuprof.Update(gtx, uData, uWhere)
}

func (syncAppStuprof *EaSyncStuProf) DeleteByKaoShiID(gtx tx.GTx, kaoShiID int64) bool {
	if kaoShiID == 0 {
		return false
	}
	uWhere := new(EaSyncStuProf)
	uWhere.KaoShiID = kaoShiID
	rlt := gtx.Where(uWhere).Delete(&EaSyncStuProf{})
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return int(rlt.RowsAffected) > 0
}

func (syncAppStuprof *EaSyncStuProf) UpdateNew(gtx tx.Db, uWhere *EaSyncStuProf, uData *EaSyncStuProf) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNew(uWhere, uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (syncAppStuprof *EaSyncStuProf) UpdateNewWithCondition(gtx tx.Db, uWhere *EaSyncStuProf, uData *EaSyncStuProf, condition interface{}, args ...interface{}) int {
	uWhere.Deleted = common.Normal
	res := gtx.UpdatesNewWithCondition(uWhere, uData, condition, args...)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

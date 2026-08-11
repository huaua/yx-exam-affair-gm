package syncApi

import (
	"fmt"
	"yx-exam-affair-gm/app/cst"
	syncDomain "yx-exam-affair-gm/app/sync/domain"
	syncService "yx-exam-affair-gm/app/sync/service"

	cacheUtil "github.com/sealsee/web-base/public/utils/cache"

	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/route"
	"github.com/sealsee/web-base/public/web"
)

var syncStuProfDomain syncDomain.EaSyncStuProf
var syncStuProfService syncService.EaSyncStuProfService

/*
 * desc: 数据同步-考生报考专业-api路由
 * author: liyu
 * date: 2025/06/13
 */
func init() {
	route.AddGroup("/api/auth/ea_sync_stu_prof", "同步考生报考专业").
		AddPOSTHandel("/get", getStuProf, true, "查询详情").
		AddPOSTHandel("/list", listStuProf, true, "查询报考列表").
		AddPOSTHandel("/get_exam_list", getExamList, true, "获取考试列表").
		AddPOSTHandel("/do_sync", doSyncStuProf, true, "执行同步操作").
		AddPOSTHandel("/get_sync_state", getSyncState, true, "获取同步状态")
}

func getStuProf(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	sessionDeptId := ctx.GetUser().DeptId
	if cst.BIZ_DEPTID_FOR_ADMISSIONS != sessionDeptId {
		json.SetErrs(errs.REFUSE_VISIT_ERR).Render()
		return
	}
	q := new(syncDomain.EaSyncStuProfQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.BaoKaoID < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	syncStuProfDB := syncStuProfDomain.GetById(q.BaoKaoID)
	json.SetObj(syncStuProfDB).Render()
}

func listStuProf(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	sessionDeptId := ctx.GetUser().DeptId
	if cst.BIZ_DEPTID_FOR_ADMISSIONS != sessionDeptId {
		json.SetErrs(errs.REFUSE_VISIT_ERR).Render()
		return
	}
	q := new(syncDomain.EaSyncStuProfQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	p := q.GetPage()
	list := syncStuProfDomain.List(q, p)
	json.SetPageList(list, p).Render()
}

// 通过openapi获取报名平台的考试列表
func getExamList(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	sessionDeptId := ctx.GetUser().DeptId
	if cst.BIZ_DEPTID_FOR_ADMISSIONS != sessionDeptId {
		json.SetErrs(errs.REFUSE_VISIT_ERR).Render()
		return
	}
	q := new(syncDomain.EaSyncStuProfQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	q.OperId = ctx.GetUserId()
	success, msg, examInfo := syncStuProfService.GetExamList(q)
	if !success {
		json.SetErr(msg).Render()
		return
	}
	json.SetList(examInfo).Render()
}

func doSyncStuProf(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	sessionDeptId := ctx.GetUser().DeptId
	if cst.BIZ_DEPTID_FOR_ADMISSIONS != sessionDeptId {
		json.SetErrs(errs.REFUSE_VISIT_ERR).Render()
		return
	}
	q := new(syncDomain.EaSyncStuProfQuery)
	if err := c.ShouldBindJSON(q); err != nil {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.KaoShiID < 1 {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	if q.KaoShiMC == "" {
		json.SetErrs(errs.PARAM_INVALID_ERR).Render()
		return
	}
	q.OperId = ctx.GetUserId()
	// 同步服务
	success, msg, stateInfo := syncStuProfService.SyncStuProfDoing(q)
	if !success {
		json.SetErr(msg).Render()
		return
	}
	json.SetObj(stateInfo).Render()
}

func getSyncState(c *gin.Context) {
	json := web.NewJsonResult(c)
	ctx := context.NewUserContext(c)
	sessionDeptId := ctx.GetUser().DeptId
	if cst.BIZ_DEPTID_FOR_ADMISSIONS != sessionDeptId {
		json.SetErrs(errs.REFUSE_VISIT_ERR).Render()
		return
	}
	var syncCache *syncDomain.SyncStateInfoCache
	var err error
	// 查询缓存状态
	if cacheUtil.Exists(cst.RDS_SYNC_STATE_KEY) {
		syncCache, err = cacheUtil.GetStruct(cst.RDS_SYNC_STATE_KEY, &syncDomain.SyncStateInfoCache{})
		if err != nil {
			json.SetErr(fmt.Sprintf("查询缓存异常[%v]", err)).Render()
			return
		}
	}
	if syncCache == nil {
		syncCache = &syncDomain.SyncStateInfoCache{State: cst.RDS_SYNC_STATE_NOT}
	}
	json.SetObj(syncCache).Render()
}

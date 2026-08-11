package syncService

import (
	"fmt"
	"strconv"
	"time"
	"yx-exam-affair-gm/app/cst"
	syncDomain "yx-exam-affair-gm/app/sync/domain"
	"yx-exam-affair-gm/app/user/userSer/userSerImpl"
	myutils "yx-exam-affair-gm/app/utils"

	"github.com/sealsee/web-base/public/basemodel"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
	cacheUtil "github.com/sealsee/web-base/public/utils/cache"
	"github.com/sealsee/web-base/public/utils/file/excel"
	"github.com/sealsee/web-base/public/utils/redis"
)

/*
 * desc: 数据同步-考生报考专业数据同步服务
 * author: liyu
 * date: 2025/06/13
 */
type EaSyncStuProfService struct{}

var syncStuProfDomain syncDomain.EaSyncStuProf
var syncStuProfDto syncDomain.EaSyncStuProfDto

var sysConfigService = userSerImpl.NewSysConfigService()

func (syncStuProfService *EaSyncStuProfService) GetExamList(q *syncDomain.EaSyncStuProfQuery) (success bool, errInfo string, examList []*syncDomain.SyncExamInfoCache) {

	examListResult, err := redis.GetStruct(cst.RDS_SYNC_EXAM_KEY, &syncDomain.ListResult{})
	if err == nil && examListResult != nil && !q.ReSync {
		return true, "", examListResult.Datas.List
	}

	url := sysConfigService.GetSysConfigByConfigKey(cst.SYS_CONFIG_SYNC_EXAM_URL).ConfigValue
	token := sysConfigService.GetSysConfigByConfigKey(cst.SYS_CONFIG_SYNC_TOKEN).ConfigValue
	url = fmt.Sprintf("%v?token=%v", url, token)
	// 调用openapi获取数据情况
	listResult, err := myutils.HttpRequest[syncDomain.ListResult](cst.POST, url, nil, nil)
	if err != nil {
		return false, "网络请求失败", nil
	}
	if !listResult.Success {
		return false, listResult.Message, nil
	}
	redis.SetStruct(cst.RDS_SYNC_EXAM_KEY, listResult, 24*time.Hour)

	return true, "", listResult.Datas.List
}

func (syncStuProfService *EaSyncStuProfService) SyncStuProfDoing(q *syncDomain.EaSyncStuProfQuery) (success bool, errInfo string, syncInfo *syncDomain.SyncStateInfoCache) {

	// 查询缓存状态
	if cacheUtil.Exists(cst.RDS_SYNC_STATE_KEY) {
		syncCache, err := cacheUtil.GetStruct(cst.RDS_SYNC_STATE_KEY, &syncDomain.SyncStateInfoCache{})
		if err != nil {
			return false, "缓存异常", nil
		}
		// 进行中、完成、失败，返回具体内容, 完成和失败弹框显示出重新同步按钮reSync
		if syncCache.State == cst.RDS_SYNC_STATE_ING ||
			(!q.ReSync && syncCache.KaoShiID == q.KaoShiID && syncCache.State == cst.RDS_SYNC_STATE_OK) ||
			(!q.ReSync && syncCache.State == cst.RDS_SYNC_STATE_ERR) {
			return true, "", syncCache
		}
	}
	// 设为同步中
	cacheUtil.SetStruct(cst.RDS_SYNC_STATE_KEY, &syncDomain.SyncStateInfoCache{State: cst.RDS_SYNC_STATE_ING, KaoShiID: q.KaoShiID}, 365*24*time.Hour)

	// 异步执行
	go func() {

		// 1.根据 q.ReSync = true 重新导入 2.接口增加kaoShiID条件过滤 3.接口查询空数据时不要删除旧数据 4.接口查询报考考生信息详情

		url := sysConfigService.GetSysConfigByConfigKey(cst.SYS_CONFIG_SYNC_APPLY_URL).ConfigValue
		token := sysConfigService.GetSysConfigByConfigKey(cst.SYS_CONFIG_SYNC_TOKEN).ConfigValue
		pageSize := myutils.StrToInt(sysConfigService.GetSysConfigByConfigKey(cst.SYS_CONFIG_SYNC_PAGE_SIZE).ConfigValue)

		curPage, pageSize, totalPage, totalSize := 1, pageSize, 0, 0
		result, err := dealHttpPostResultPage(url, token, q, curPage, pageSize, totalPage, totalSize, true)
		if err != nil {
			return
		}
		for result != nil && curPage < result.Datas.Page.TotalPage {
			curPage++
			result, err = dealHttpPostResultPage(url, token, q, curPage, pageSize, result.Datas.Page.TotalPage, result.Datas.Page.TotalSize, false)
			if err != nil {
				return
			}
		}
	}()

	// 返回当前的进度信息
	syncCache, _ := cacheUtil.GetStruct(cst.RDS_SYNC_STATE_KEY, &syncDomain.SyncStateInfoCache{})
	return true, "", syncCache
}

func dealHttpPostResultPage(url, token string, q *syncDomain.EaSyncStuProfQuery, curPage, pageSize, totalPage, totalSize int, delAll bool) (*syncDomain.PageResult, error) {
	cache := &syncDomain.SyncStateInfoCache{CurPage: curPage}
	cache.KaoShiID = q.KaoShiID
	cache.KaoShiMC = q.KaoShiMC
	cache.PageSize = pageSize
	cache.CurPage = curPage
	cache.TotalPage = totalPage
	cache.TotalSize = totalSize
	var percent string
	if totalPage > 0 {
		percent = strconv.FormatFloat(float64(curPage)/float64(totalPage)*100, 'f', 2, 64)
	}
	cache.Percent = percent

	param := fmt.Sprintf(`{"m":"","p":{"kaoShiID":%v,"curPage":%v,"pageSize":%v}}`, q.KaoShiID, curPage, pageSize)
	url = fmt.Sprintf("%v?token=%v&data=%v", url, token, param)
	// 调用openapi获取数据情况
	resultPage, err := myutils.HttpRequest[syncDomain.PageResult](cst.POST, url, nil, nil)
	if err != nil {
		cache.State = cst.RDS_SYNC_STATE_ERR
		cache.ErrMsg = fmt.Sprintf("接口异常[%v][%v]", url, err.Error())
		cacheUtil.SetStruct(cst.RDS_SYNC_STATE_KEY, cache, 365*24*time.Hour)
		return nil, err
	}
	if !resultPage.Success {
		cache.State = cst.RDS_SYNC_STATE_ERR
		cache.ErrMsg = fmt.Sprintf("同步失败[%v][%v]", resultPage.Code, resultPage.Message)
		cacheUtil.SetStruct(cst.RDS_SYNC_STATE_KEY, cache, 365*24*time.Hour)
		return nil, err
	}
	// fmt.Printf("success:%v, code:%v, msg:%v, totalSize:%v, totalPage:%v, curPage:%v \n",
	// 	resultPage.Success, resultPage.Code, resultPage.Message,
	// 	resultPage.Datas.Page.TotalSize, resultPage.Datas.Page.TotalPage, resultPage.Datas.Page.CurPage)

	// 更新缓存同步状态信息
	totalPage = resultPage.Datas.Page.TotalPage
	totalSize = resultPage.Datas.Page.TotalSize
	state := cst.RDS_SYNC_STATE_ING
	if totalPage > 0 {
		percent = strconv.FormatFloat(float64(curPage)/float64(totalPage)*100, 'f', 2, 64)
	}
	if curPage == totalPage {
		state = cst.RDS_SYNC_STATE_OK
	}
	cache.State = state
	cache.TotalPage = totalPage
	cache.TotalSize = totalSize
	cache.Percent = percent

	hasData := true
	// 如果没数据，直接返回成功，不要清数据
	if len(resultPage.Datas.Page.DataList) < 1 {
		// cache.State = cst.RDS_SYNC_STATE_ERR
		// cache.ErrMsg = "接口获取报考记录为空"
		// cacheUtil.SetStruct(cst.RDS_SYNC_STATE_KEY, cache, 365*24*time.Hour)
		// return nil, fmt.Errorf("[ksid:%v,page:%v] api get null", q.KaoShiID, curPage)
		hasData = false
		cache.State = cst.RDS_SYNC_STATE_OK
	}

	if hasData {
		// 遍历补充
		for _, data := range resultPage.Datas.Page.DataList {
			data.SyncTime = basemodel.BaseTime(time.Now())
			data.Deleted = cst.YES
			data.SetCreateBy(q.OperId)
		}

		// 数据库执行事务
		var success bool
		tx.ExecGTx(func(gtx tx.GTx) {
			// 批量删除考试下的所有报考记录, 首次查询时，根据传参判断是否需要删除
			if delAll {
				syncStuProfDomain.DeleteByKaoShiID(gtx, q.KaoShiID)
			}
			// 批量插入数据
			success = syncStuProfDto.BatchInsert(gtx, resultPage.Datas.Page.DataList)
		})
		if !success {
			cache.State = cst.RDS_SYNC_STATE_ERR
			cache.ErrMsg = "批量插入异常"
			cacheUtil.SetStruct(cst.RDS_SYNC_STATE_KEY, cache, 365*24*time.Hour)
			return nil, fmt.Errorf("[ksmc:%v,page:%v] batch insert err", q.KaoShiMC, curPage)
		}
	}

	// 更新进度缓存
	cacheUtil.SetStruct(cst.RDS_SYNC_STATE_KEY, cache, 365*24*time.Hour)
	return resultPage, nil
}

// 导出服务 -- start
type EaSyncStuProfExportStruct struct {
	Query *syncDomain.EaSyncStuProfQuery
}

func (h *EaSyncStuProfExportStruct) Title() string {
	return "同步报考记录导出"
}

func (h *EaSyncStuProfExportStruct) HeaderColumn() []string {
	return []string{"报考ID,BaoKaoID", "用户ID,YongHuID", "姓名,XingMing"}
}

func (h *EaSyncStuProfExportStruct) Rows(p *page.Page) []map[string]interface{} {

	return syncStuProfDomain.ListMapWithCondition(h.Query, p, nil, nil)
}

func (h *EaSyncStuProfExportStruct) Finish(url string) {
}

// 导出服务 -- end
func (syncStuProfService *EaSyncStuProfService) ExportSyncStuProf(q *syncDomain.EaSyncStuProfQuery) ([]byte, error) {
	return excel.NewExcel().ExportSync(&EaSyncStuProfExportStruct{q})
}

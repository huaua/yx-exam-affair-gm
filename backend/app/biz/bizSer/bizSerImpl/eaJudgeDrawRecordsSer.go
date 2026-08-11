package bizSerImpl

import (
	"yx-exam-affair-gm/app/biz/bizDao"
	"yx-exam-affair-gm/app/biz/bizDao/bizDaoImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
)

/*
 * desc: 评委抽取规则服务层实现
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeDrawRecordsService struct{}

var eaJudgeDrawRecordsService *EaJudgeDrawRecordsService
var eaJudgeDrawRecordsDao bizDao.IEaJudgeDrawRecordsDao

func init() {
	eaJudgeDrawRecordsService = &EaJudgeDrawRecordsService{}
	eaJudgeDrawRecordsDao = bizDaoImpl.NewEaJudgeDrawRecordsDao()
}

func NewEaJudgeDrawRecordsService() *EaJudgeDrawRecordsService {
	return eaJudgeDrawRecordsService
}

func (judgeDrawRecords *EaJudgeDrawRecordsService) ListAllEaJudgeDrawRecords(q *bizDto.EaJudgeDrawRecordsQuery) []*bizDo.EaJudgeDrawRecords {
	listJudgeExpertAll := make([]*bizDo.EaJudgeDrawRecords, 0)
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE
	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		list := eaJudgeDrawRecordsDao.List(q, pageInfo)
		if len(list) == 0 {
			break
		}
		listJudgeExpertAll = append(listJudgeExpertAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}
	return listJudgeExpertAll
}

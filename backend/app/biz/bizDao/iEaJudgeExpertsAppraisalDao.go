package bizDao

import (
    "yx-exam-affair-gm/app/biz/model/bizDo"
    "yx-exam-affair-gm/app/biz/model/bizDto"

    "github.com/sealsee/web-base/public/ds/page"
    "github.com/sealsee/web-base/public/ds/tx"
)

/*
 * desc: 评委评价数据层接口
 * author: init code
 * date: 2025/06/11
 */
type IEaJudgeExpertsAppraisalDao interface {

    /*
     * desc: 根据id获取数据
     * author: init code
     * date: 2025/06/11
     */
    GetById(id int64) *bizDo.EaJudgeExpertsAppraisal

    /*
     * desc: 查询满足条件的条数
     * author: init code
     * date: 2025/06/11
     */
    Count(q *bizDto.EaJudgeExpertsAppraisalQuery) int

    /*
     * desc: 查询满足条件的条数（支持特殊条件）
     * author: init code
     * date: 2025/06/11
     */
    CountWithCondition(q *bizDto.EaJudgeExpertsAppraisalQuery, condition interface{}, args ...interface{}) int

    /*
     * desc: 查询满足条件的记录列表
     * author: init code
     * date: 2025/06/11
     */
    List(q *bizDto.EaJudgeExpertsAppraisalQuery, p *page.Page) []*bizDo.EaJudgeExpertsAppraisal

    /*
     * desc: 查询满足条件的记录列表（支持特殊条件）
     * author: init code
     * date: 2025/06/11
     */
    ListWithCondition(q *bizDto.EaJudgeExpertsAppraisalQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaJudgeExpertsAppraisal

    /*
     * desc: 查询满足条件的记录列表map（支持特殊条件）
     * author: init code
     * date: 2025/06/11
     */
    ListMapWithCondition(q *bizDto.EaJudgeExpertsAppraisalQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any

    /*
     * desc: 新增记录
     * author: init code
     * date: 2025/06/11
     */
    Insert(gtx tx.GTx, judgeExpertsAppraisal *bizDo.EaJudgeExpertsAppraisal) bool

    /*
     * desc: 批量新增记录
     * author: init code
     * date: 2025/06/11
     */
    BatchInsert(gtx tx.GTx, judgeExpertsAppraisalList []*bizDo.EaJudgeExpertsAppraisal) bool

    /*
     * desc: 更新记录
     * author: init code
     * date: 2025/06/11
     */
    Update(gtx tx.GTx, uData *bizDo.EaJudgeExpertsAppraisal, uWhere *bizDo.EaJudgeExpertsAppraisal) bool

    /*
     * desc: 删除记录
     * author: init code
     * date: 2025/06/11
     */
    DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool

    /*
     * desc: 更新（支持清空）
     * author: init code
     * date: 2025/06/11
     */
    UpdateNew(gtx tx.Db, uWhere *bizDo.EaJudgeExpertsAppraisal, uData *bizDo.EaJudgeExpertsAppraisal) int

    /*
     * desc: 特殊条件更新（支持清空）
     * author: init code
     * date: 2025/06/11
     */
    UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaJudgeExpertsAppraisal, uData *bizDo.EaJudgeExpertsAppraisal, condition interface{}, args ...interface{}) int
}

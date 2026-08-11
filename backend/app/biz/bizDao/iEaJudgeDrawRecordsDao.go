package bizDao

import (
    "yx-exam-affair-gm/app/biz/model/bizDo"
    "yx-exam-affair-gm/app/biz/model/bizDto"

    "github.com/sealsee/web-base/public/ds/page"
    "github.com/sealsee/web-base/public/ds/tx"
)

/*
 * desc: 评委抽取规则数据层接口
 * author: init code
 * date: 2025/06/11
 */
type IEaJudgeDrawRecordsDao interface {

    /*
     * desc: 根据id获取数据
     * author: init code
     * date: 2025/06/11
     */
    GetById(id int64) *bizDo.EaJudgeDrawRecords

    /*
     * desc: 查询满足条件的条数
     * author: init code
     * date: 2025/06/11
     */
    Count(q *bizDto.EaJudgeDrawRecordsQuery) int

    /*
     * desc: 查询满足条件的条数（支持特殊条件）
     * author: init code
     * date: 2025/06/11
     */
    CountWithCondition(q *bizDto.EaJudgeDrawRecordsQuery, condition interface{}, args ...interface{}) int

    /*
     * desc: 查询满足条件的记录列表
     * author: init code
     * date: 2025/06/11
     */
    List(q *bizDto.EaJudgeDrawRecordsQuery, p *page.Page) []*bizDo.EaJudgeDrawRecords

    /*
     * desc: 查询满足条件的记录列表（支持特殊条件）
     * author: init code
     * date: 2025/06/11
     */
    ListWithCondition(q *bizDto.EaJudgeDrawRecordsQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaJudgeDrawRecords

    /*
     * desc: 查询满足条件的记录列表map（支持特殊条件）
     * author: init code
     * date: 2025/06/11
     */
    ListMapWithCondition(q *bizDto.EaJudgeDrawRecordsQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any

    /*
     * desc: 新增记录
     * author: init code
     * date: 2025/06/11
     */
    Insert(gtx tx.GTx, judgeDrawRecords *bizDo.EaJudgeDrawRecords) bool

    /*
     * desc: 批量新增记录
     * author: init code
     * date: 2025/06/11
     */
    BatchInsert(gtx tx.GTx, judgeDrawRecordsList []*bizDo.EaJudgeDrawRecords) bool

    /*
     * desc: 更新记录
     * author: init code
     * date: 2025/06/11
     */
    Update(gtx tx.GTx, uData *bizDo.EaJudgeDrawRecords, uWhere *bizDo.EaJudgeDrawRecords) bool

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
    UpdateNew(gtx tx.Db, uWhere *bizDo.EaJudgeDrawRecords, uData *bizDo.EaJudgeDrawRecords) int

    /*
     * desc: 特殊条件更新（支持清空）
     * author: init code
     * date: 2025/06/11
     */
    UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaJudgeDrawRecords, uData *bizDo.EaJudgeDrawRecords, condition interface{}, args ...interface{}) int
}

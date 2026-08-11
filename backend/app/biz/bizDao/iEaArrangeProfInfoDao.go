package bizDao

import (
    "yx-exam-affair-gm/app/biz/model/bizDo"
    "yx-exam-affair-gm/app/biz/model/bizDto"

    "github.com/sealsee/web-base/public/ds/page"
    "github.com/sealsee/web-base/public/ds/tx"
)

/*
 * desc: 编排专业数据层接口
 * author: init code
 * date: 2024/09/24
 */
type IEaArrangeProfInfoDao interface {

    /*
     * desc: 根据id获取数据
     * author: init code
     * date: 2024/09/24
     */
    GetById(id int64) *bizDo.EaArrangeProfInfo

    /*
     * desc: 查询满足条件的条数
     * author: init code
     * date: 2024/09/24
     */
    Count(q *bizDto.EaArrangeProfInfoQuery) int

    /*
     * desc: 查询满足条件的条数（支持特殊条件）
     * author: init code
     * date: 2024/09/24
     */
    CountWithCondition(q *bizDto.EaArrangeProfInfoQuery, condition interface{}, args ...interface{}) int

    /*
     * desc: 查询满足条件的记录列表
     * author: init code
     * date: 2024/09/24
     */
    List(q *bizDto.EaArrangeProfInfoQuery, p *page.Page) []*bizDo.EaArrangeProfInfo

    /*
     * desc: 查询满足条件的记录列表（支持特殊条件）
     * author: init code
     * date: 2024/09/24
     */
    ListWithCondition(q *bizDto.EaArrangeProfInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaArrangeProfInfo

    /*
     * desc: 查询满足条件的记录列表map（支持特殊条件）
     * author: init code
     * date: 2024/09/24
     */
    ListMapWithCondition(q *bizDto.EaArrangeProfInfoQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any

    /*
     * desc: 新增记录
     * author: init code
     * date: 2024/09/24
     */
    Insert(gtx tx.GTx, arrangeProfInfo *bizDo.EaArrangeProfInfo) bool

    /*
     * desc: 批量新增记录
     * author: init code
     * date: 2024/09/24
     */
    BatchInsert(gtx tx.GTx, arrangeProfInfoList []*bizDo.EaArrangeProfInfo) bool

    /*
     * desc: 更新记录
     * author: init code
     * date: 2024/09/24
     */
    Update(gtx tx.GTx, uData *bizDo.EaArrangeProfInfo, uWhere *bizDo.EaArrangeProfInfo) bool

    /*
     * desc: 删除记录
     * author: init code
     * date: 2024/09/24
     */
    DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool

    /*
     * desc: 更新（支持清空）
     * author: init code
     * date: 2024/09/24
     */
    UpdateNew(gtx tx.Db, uWhere *bizDo.EaArrangeProfInfo, uData *bizDo.EaArrangeProfInfo) int

    /*
     * desc: 特殊条件更新（支持清空）
     * author: init code
     * date: 2024/09/24
     */
    UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaArrangeProfInfo, uData *bizDo.EaArrangeProfInfo, condition interface{}, args ...interface{}) int
}

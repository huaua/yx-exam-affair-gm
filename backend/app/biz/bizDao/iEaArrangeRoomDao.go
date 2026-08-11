package bizDao

import (
    "yx-exam-affair-gm/app/biz/model/bizDo"
    "yx-exam-affair-gm/app/biz/model/bizDto"

    "github.com/sealsee/web-base/public/ds/page"
    "github.com/sealsee/web-base/public/ds/tx"
)

/*
 * desc: 考场编排详情数据层接口
 * author: init code
 * date: 2024/09/24
 */
type IEaArrangeRoomDao interface {

    /*
     * desc: 根据id获取数据
     * author: init code
     * date: 2024/09/24
     */
    GetById(id int64) *bizDo.EaArrangeRoom

    /*
     * desc: 查询满足条件的条数
     * author: init code
     * date: 2024/09/24
     */
    Count(q *bizDto.EaArrangeRoomQuery) int

    /*
     * desc: 查询满足条件的条数（支持特殊条件）
     * author: init code
     * date: 2024/09/24
     */
    CountWithCondition(q *bizDto.EaArrangeRoomQuery, condition interface{}, args ...interface{}) int

    /*
     * desc: 查询满足条件的记录列表
     * author: init code
     * date: 2024/09/24
     */
    List(q *bizDto.EaArrangeRoomQuery, p *page.Page) []*bizDo.EaArrangeRoom

    /*
     * desc: 查询满足条件的记录列表（支持特殊条件）
     * author: init code
     * date: 2024/09/24
     */
    ListWithCondition(q *bizDto.EaArrangeRoomQuery, p *page.Page, condition interface{}, args ...interface{}) []*bizDo.EaArrangeRoom

    /*
     * desc: 查询满足条件的记录列表map（支持特殊条件）
     * author: init code
     * date: 2024/09/24
     */
    ListMapWithCondition(q *bizDto.EaArrangeRoomQuery, p *page.Page, condition interface{}, args ...interface{}) []map[string]any

    /*
     * desc: 新增记录
     * author: init code
     * date: 2024/09/24
     */
    Insert(gtx tx.GTx, arrangeRoom *bizDo.EaArrangeRoom) bool

    /*
     * desc: 批量新增记录
     * author: init code
     * date: 2024/09/24
     */
    BatchInsert(gtx tx.GTx, arrangeRoomList []*bizDo.EaArrangeRoom) bool

    /*
     * desc: 更新记录
     * author: init code
     * date: 2024/09/24
     */
    Update(gtx tx.GTx, uData *bizDo.EaArrangeRoom, uWhere *bizDo.EaArrangeRoom) bool

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
    UpdateNew(gtx tx.Db, uWhere *bizDo.EaArrangeRoom, uData *bizDo.EaArrangeRoom) int

    /*
     * desc: 特殊条件更新（支持清空）
     * author: init code
     * date: 2024/09/24
     */
    UpdateNewWithCondition(gtx tx.Db, uWhere *bizDo.EaArrangeRoom, uData *bizDo.EaArrangeRoom, condition interface{}, args ...interface{}) int
}

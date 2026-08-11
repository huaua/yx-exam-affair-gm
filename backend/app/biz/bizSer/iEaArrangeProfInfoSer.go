package bizSer

import (
	"github.com/sealsee/web-base/public/errs"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/ds/page"
)

/*
 * desc: 编排专业服务层接口
 * author: init code
 * date: 2024/09/24
 */
type IEaArrangeProfInfoService interface {

	/*
	 * desc: 根据id获取数据
	 * author: init code
	 * date: 2024/09/24
	 */
	GetEaArrangeProfInfoById(id int64) *bizDo.EaArrangeProfInfo

	/*
	 * desc: 查询满足条件的记录列表
	 * author: init code
	 * date: 2024/09/24
	 */
	ListEaArrangeProfInfo(q *bizDto.EaArrangeProfInfoQuery, p *page.Page) []*bizDo.EaArrangeProfInfo

	/*
	 * desc: 插入数据
	 * author: init code
	 * date: 2024/09/24
	 */
	InsertEaArrangeProfInfo(arrangeProfInfo *bizDo.EaArrangeProfInfo) bool

	/*
	 * desc: 更新数据
	 * author: init code
	 * date: 2024/09/24
	 */
	UpdateEaArrangeProfInfoById(arrangeProfInfo *bizDo.EaArrangeProfInfo) bool

	/*
	 * desc: 取消编排数据
	 * author: 蓝鲸
	 * date: 2024/09/26
	 */
	CancelArrangeById(arrangeProfInfo *bizDo.EaArrangeProfInfo) (bool, errs.ERROR)

	/*
	 * desc: 删除数据
	 * author: init code
	 * date: 2024/09/24
	 */
	DeleteEaArrangeProfInfoById(id int64, deleteBy int64) bool

	/*
	 * desc: 导出数据
	 * author: init code
	 * date: 2024/09/24
	 */
	ExportEaArrangeProfInfo(q *bizDto.EaArrangeProfInfoQuery) ([]byte, error)

	/*
	 * desc: 导入数据
	 * author: 蓝鲸
	 * date: 2024/09/25
	 */
	ImportEaArrangeProfInfo(q *bizDto.EaArrangeProfInfoQuery) (bool, string)

	/*
	 * desc: 查询所有专业方向信息
	 * author: 蓝鲸
	 * date: 2024/09/25
	 */
	ListAllProfInfo(q *bizDto.EaArrangeProfInfoQuery) []*bizDo.EaArrangeProfInfo

	/*
	 * desc: 查询所有专业方向信息
	 * author: 蓝鲸
	 * date: 2024/09/25
	 */
	ListAvailableRoomInTask(q *bizDto.EaArrangeProfInfoQuery) []*bizDo.EaRoomTaskReport
}

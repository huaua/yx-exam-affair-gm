package bizSerImpl

import (
	"encoding/json"
	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/utils/snowflake"
	"strconv"
	"time"
	"unicode/utf8"
	"yx-exam-affair-gm/app/biz/bizDao"
	"yx-exam-affair-gm/app/biz/bizDao/bizDaoImpl"
	"yx-exam-affair-gm/app/biz/model/bizDo"
	"yx-exam-affair-gm/app/biz/model/bizDto"
	"yx-exam-affair-gm/app/cst"
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"
	"yx-exam-affair-gm/app/user/userDao"
	"yx-exam-affair-gm/app/user/userDao/userDaoImpl"
	"yx-exam-affair-gm/app/user/userSer/userSerImpl"
	cmutils "yx-exam-affair-gm/app/utils"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
	"github.com/sealsee/web-base/public/utils/file/excel"
)

/*
 * desc: 教室服务层实现
 * author: 蓝鲸
 * date: 2024/08/14
 */
type EaRoomInfoService struct{}

var eaRoomInfoService *EaRoomInfoService
var eaRoomInfoDao bizDao.IEaRoomInfoDao
var sysDicDataService *userSerImpl.SysDictDataService
var eaTaskRoomService *EaRoomTaskReportService
var sysDataInoutDao userDao.ISysDataInoutDao

func init() {
	eaRoomInfoService = &EaRoomInfoService{}
	eaRoomInfoDao = bizDaoImpl.NewEaRoomInfoDao()
	sysDataInoutDao = userDaoImpl.NewSysDataInoutDao()
	sysDicDataService = userSerImpl.NewSysDictDataService()
	eaTaskRoomService = &EaRoomTaskReportService{}
}

func NewEaRoomInfoService() *EaRoomInfoService {
	return eaRoomInfoService
}

func (roomInfo *EaRoomInfoService) GetEaClassroomInfoById(id int64) *bizDo.EaRoomInfo {
	return eaRoomInfoDao.GetById(id)
}

func (roomInfo *EaRoomInfoService) ListEaClassroomInfo(q *bizDto.EaRoomInfoQuery, p *page.Page) []*bizDo.EaRoomInfo {
	q.AddLikeAll("room_name", q.RoomName)
	list := eaRoomInfoDao.List(q, p)
	if len(list) <= 0 {
		return list
	}

	deptMap := eaDeptInfoService.GetAllDeptMap()
	for _, i := range list {
		if cst.BIZ_DEPTID_FOR_ADMISSIONS == i.DeptId {
			i.DeptName = "公共教室"
		} else {
			if deptMap[i.DeptId] != nil {
				i.DeptName = deptMap[i.DeptId].DeptName
			}
		}
	}
	return list
}

func (roomInfo *EaRoomInfoService) ListAllEaClassroomInfo(q *bizDto.EaRoomInfoQuery) []*bizDo.EaRoomInfo {
	listDepartAll := make([]*bizDo.EaRoomInfo, 0)
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE
	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		list := eaRoomInfoDao.List(q, pageInfo)
		if len(list) == 0 {
			break
		}
		listDepartAll = append(listDepartAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}
	return listDepartAll
}

func (roomInfo *EaRoomInfoService) InsertEaClassroomInfo(classroomInfoDO *bizDo.EaRoomInfo) (success bool, errInfo string) {
	if len(classroomInfoDO.RoomName) <= 0 || len(classroomInfoDO.LevelCode) <= 0 {
		return false, "教室名称、所属层次不能为空"
	}

	// 层级校验
	levelQuery := &userDto.SysDictDataQuery{DictType: cst.BIZ_LEVEL, DictValue: classroomInfoDO.LevelCode, State: strconv.Itoa(cst.YES)}
	levelQuery.PageSize = 1
	levelList := sysDicDataService.ListSysDictData(levelQuery, &page.Page{CurPage: 1, PageSize: 1})
	if len(levelList) <= 0 {
		return false, "所属层次不正确"
	}

	_, errInfo = _validateClassRoomInfo(classroomInfoDO)
	if len(errInfo) > 0 {
		return false, errInfo
	}

	// 检查名称和层次是否已经存在
	exit := eaRoomInfoDao.Count(&bizDto.EaRoomInfoQuery{RoomName: classroomInfoDO.RoomName, LevelCode: classroomInfoDO.LevelCode}) > 0
	if exit {
		return false, "请检查层次下是否已经存在教室名称"
	}

	classroomInfoDO.Deleted = common.Normal
	classroomInfoDO.State = common.OK
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaRoomInfoDao.Insert(gtx, classroomInfoDO)
	})
	return
}

func _validateClassRoomInfo(classroomInfoDO *bizDo.EaRoomInfo) (success bool, errInfo string) {

	// 组数和按位容量必须有一个大于0
	if classroomInfoDO.GroupNum <= 0 && classroomInfoDO.Capacity <= 0 {
		return false, "请输入按位容量"
	}

	// 如果组数不为0，则必须校验按组容量大于0
	if classroomInfoDO.GroupNum > 0 && classroomInfoDO.GroupCapacity <= 0 {
		return false, "请输入按组容量"
	}

	if classroomInfoDO.GroupNum <= 0 {
		classroomInfoDO.GroupCapacity = 0
	}

	if len(classroomInfoDO.RoomName) > 0 {
		if utf8.RuneCountInString(classroomInfoDO.RoomName) > 20 {
			return false, "教室名称长度不能超过20个字符"
		}
	}

	if len(classroomInfoDO.CampusName) > 0 {
		if utf8.RuneCountInString(classroomInfoDO.CampusName) < 2 || utf8.RuneCountInString(classroomInfoDO.CampusName) > 15 {
			return false, "所属校区允许只输入2-15个字符"
		}
	}
	return true, ""
}

func (roomInfo *EaRoomInfoService) UpdateEaClassroomInfoById(classroomInfoDO *bizDo.EaRoomInfo, onlyChangeStatus bool, deptId int64) (success bool, errInfo string) {
	if !onlyChangeStatus {
		_, errInfo = _validateClassRoomInfo(classroomInfoDO)
		if len(errInfo) > 0 {
			return false, errInfo
		}
	}

	classroomInfoInDB := roomInfo.GetEaClassroomInfoById(classroomInfoDO.RoomId)
	if classroomInfoInDB == nil {
		return false, "教室不存在"
	}

	if deptId != classroomInfoInDB.DeptId {
		return false, "对不起，您无权限修改该教室信息"
	}

	// 如果修改了名称，需要校验此名称是否已存在
	if len(classroomInfoDO.RoomName) > 0 && classroomInfoInDB.RoomName != classroomInfoDO.RoomName {
		// 检查名称和层次是否已经存在
		exit := eaRoomInfoDao.Count(&bizDto.EaRoomInfoQuery{RoomName: classroomInfoDO.RoomName, LevelCode: classroomInfoDO.LevelCode}) > 0
		if exit {
			return false, "请检查层次下是否已经存在教室名称"
		}
	}

	uWhere := &bizDo.EaRoomInfo{RoomId: classroomInfoDO.RoomId}
	if !onlyChangeStatus {
		classroomInfoDO.SetNullableCols("campus_name")
		classroomInfoDO.SetToZeroCols("group_num", "group_capacity", "capacity")
	}
	success = tx.ExecGTx(func(gtx tx.Db) {
		// 如果教室名称更新了，需要更新冗余的名称字段
		eaTaskRoomService.UpdateRoomNameById(classroomInfoDO, classroomInfoDO.UpdateBy)
		eaRoomInfoDao.UpdateNew(gtx, uWhere, classroomInfoDO)
	})
	return success, ""
}

func (roomInfo *EaRoomInfoService) DeleteEaClassroomInfoById(id int64, deleteBy int64) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = eaRoomInfoDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

// 导出服务 -- start
type EaClassroomInfoExportStruct struct {
	Query *bizDto.EaRoomInfoQuery
}

func (h *EaClassroomInfoExportStruct) Title() string {
	return "教室导出"
}

func (h *EaClassroomInfoExportStruct) HeaderColumn() []string {
	return []string{"教室ID,room_id", "教室名称,room_name", "所在校区,campus_name", "所属层次,level_code", "教室类型,room_type", "所属学院id,dept_id", "组数,group_num", "容量,capacity", "启用状态: 1-正常 2-停用,state", "删除标记：1-正常的；2-已删除的；,deleted", "备注,remark", "创建者,create_by", "创建时间,create_time", "更新者,update_by", "更新时间,update_time"}
}

func (h *EaClassroomInfoExportStruct) Rows(p *page.Page) []map[string]interface{} {
	return eaRoomInfoDao.ListMapWithCondition(h.Query, p, nil, nil)
}

func (h *EaClassroomInfoExportStruct) Finish(url string) {
}

// 导出服务 -- end
func (roomInfo *EaRoomInfoService) ExportEaClassroomInfo(q *bizDto.EaRoomInfoQuery) ([]byte, error) {
	return excel.NewExcel().ExportSync(&EaClassroomInfoExportStruct{q})
}

// 导入服务 -- start
type HandlerImpRoomInfoStruct struct {
	BizType          int
	BatchNo          string
	ExistRoomInfoMap map[string]*bizDo.EaRoomInfo
	LevelMap         map[string]string
	Query            *bizDto.EaRoomInfoQuery
	List             []*bizDo.EaRoomInfo
	Msg              string
	Err              string
	ExcelHeader      []string
}

func (h *HandlerImpRoomInfoStruct) Headers(headers []string) {
	if len(headers) != len(h.ExcelHeader) {
		h.Err += "表头数据错误!"
		return
	}
	for i, header := range headers {
		if header != h.ExcelHeader[i] {
			h.Err += "表头数据错误!"
		}
	}
}

func (h *HandlerImpRoomInfoStruct) Row(mp *map[string]string) excel.RunStatus {
	// 判断表头数据是否错误
	if h.Err != "" {
		return excel.Exit
	}
	m := *mp
	rowIdx := strconv.Itoa(cmutils.StrToInt(m["rid"]) + 1)
	RoomInfo22Insert := new(bizDo.EaRoomInfo)

	// 取值
	roomName := m[h.ExcelHeader[0]]
	campusName := m[h.ExcelHeader[1]]
	levelName := m[h.ExcelHeader[2]]
	groupNumStr := m[h.ExcelHeader[3]]
	groupCapacityStr := m[h.ExcelHeader[4]]
	capacity := m[h.ExcelHeader[5]]
	deptId := h.Query.DeptId

	if len(roomName) <= 0 || len(levelName) <= 0 {
		h.Err = "第" + rowIdx + "行，教室名称、所属层次不能为空！"
		return excel.Exit
	}
	if len(h.LevelMap[levelName]) <= 0 {
		h.Err = "第" + rowIdx + "行，所属层次有误"
		return excel.Exit
	}
	levelCode := h.LevelMap[levelName]

	RoomInfo22Insert.RoomName = roomName
	RoomInfo22Insert.CampusName = campusName
	RoomInfo22Insert.LevelCode = h.LevelMap[levelName]
	if len(groupNumStr) > 0 {
		if !cmutils.IsPositiveInteger(groupNumStr) {
			h.Err = "第" + rowIdx + "行，组数只能是正整数！"
			return excel.Exit
		}
		RoomInfo22Insert.GroupNum = cmutils.StrToInt(groupNumStr)
	}
	if len(groupCapacityStr) > 0 {
		if !cmutils.IsPositiveInteger(groupCapacityStr) {
			h.Err = "第" + rowIdx + "行，按组容量只能是正整数！"
			return excel.Exit
		}

		// bug[10274]设置了按组容量，但没有设置组数，也需要提示
		if RoomInfo22Insert.GroupNum <= 0 {
			h.Err = "第" + rowIdx + "行，设置了按组容量但未设置组数！"
			return excel.Exit
		}
		RoomInfo22Insert.GroupCapacity = cmutils.StrToInt(groupCapacityStr)
	}
	if len(capacity) > 0 {
		if !cmutils.IsPositiveInteger(capacity) {
			h.Err = "第" + rowIdx + "行，按位容量只能是正整数！"
			return excel.Exit
		}
		RoomInfo22Insert.Capacity = cmutils.StrToInt(capacity)
	}

	RoomInfo22Insert.DeptId = deptId
	success, errInfo := _validateClassRoomInfo(RoomInfo22Insert)
	if !success {
		h.Err = "第" + rowIdx + "行，" + errInfo
		return excel.Exit
	}

	// 数据重复校验
	key := roomName + "_" + levelCode
	if h.ExistRoomInfoMap[key] != nil {
		h.Err = "第" + rowIdx + "行，该层次下已经存在此教室名称！"
		return excel.Exit
	}
	h.ExistRoomInfoMap[key] = RoomInfo22Insert

	// 创建人、时间
	RoomInfo22Insert.Deleted = common.Normal
	RoomInfo22Insert.State = common.OK
	RoomInfo22Insert.SetCreateBy(h.Query.OperId)
	RoomInfo22Insert.SetUpdateBy(h.Query.OperId)
	// 放入list
	h.List = append(h.List, RoomInfo22Insert)
	return excel.Next
}

func (h *HandlerImpRoomInfoStruct) After() {
	if len(h.Err) > 0 {
		return
	}
	// 判断是否有导入数据
	if len(h.List) < 1 {
		h.Err = "没有导入数据！"
		return
	}
	h.Msg += "成功导入 " + strconv.Itoa(len(h.List)) + " 条记录!\n"
	// 拼装导入记录
	dataInout := new(userDo.SysDataInout)
	dataInout.OperType = common.OPER_TYPE_IMPORT_DATA
	dataInout.BizType = h.BizType
	dataInout.BatchNo = h.BatchNo
	// 压缩导入的多个文件
	if len(h.Query.Files) > 1 {
		dataInout.FileIn = cmutils.CompressZipFiles(h.Query.Files)
		dataInout.FileInName = "多文件"
	} else {
		for k, v := range h.Query.Files {
			dataInout.FileIn = v
			dataInout.FileInName = k
		}
	}
	// 传参json化
	h.Query.Files = nil // 文件不放在条件参数里记录
	p, _ := json.Marshal(h.Query)
	pm := string(p)
	if len(pm) > 1000 {
		dataInout.Param = pm[:1000] //截取1000个字符长度
	} else {
		dataInout.Param = pm
	}
	dataInout.OperTime = time.Now()
	dataInout.State = common.OK
	dataInout.Remark = h.Msg
	dataInout.SetCreateBy(h.Query.OperId)
	dataInout.SetUpdateBy(h.Query.OperId)
	// 数据库执行事务
	success := tx.ExecGTx(func(gtx tx.GTx) {
		// 最后批量插入导入的数据，并新增一条导入记录
		eaRoomInfoDao.BatchInsert(gtx, h.List)
		sysDataInoutDao.Insert(gtx, dataInout)
	})
	if !success {
		h.Err = "导入失败，请检查是否与已有数据重复！"
	}
}

// 导入服务  --  end

// 导入
func (roomInfo *EaRoomInfoService) ImportEaClassroomInfo(q *bizDto.EaRoomInfoQuery) (success bool, msg string) {

	// 使用雪花算法生成全局唯一id，作为本次导入操作的批次号
	currentBatchNo := strconv.FormatInt(snowflake.GenID(), 10)

	// 查询所有层次信息放入map
	levelMap := make(map[string]string)
	p := &page.Page{CurPage: 1, PageSize: page.MAX_PAGE_SIZE}
	dictDataQuery := &userDto.SysDictDataQuery{DictType: cst.BIZ_LEVEL, State: strconv.Itoa(cst.YES)}
	levelList := sysDicDataService.ListSysDictData(dictDataQuery, p)
	for _, v := range levelList {
		levelMap[v.DictLabel] = v.DictValue
	}

	// 查询所有已有的分数线信息放入map
	classroomInfos := roomInfo.ListAllEaClassroomInfo(new(bizDto.EaRoomInfoQuery))
	dbRoomMap := make(map[string]*bizDo.EaRoomInfo)
	for _, v := range classroomInfos {
		key := v.RoomName + "_" + v.LevelCode
		dbRoomMap[key] = v
	}

	// 遍历文件map  文件名->文件地址...
	excelHeader := []string{"教室名称", "所属校区", "层次", "组数", "按组容量", "按位容量"}
	handlerImpStruct := &HandlerImpRoomInfoStruct{BizType: int(cst.BIZ_ROOM_BASE_INFO), BatchNo: currentBatchNo,
		Query: q, LevelMap: levelMap, ExistRoomInfoMap: dbRoomMap, ExcelHeader: excelHeader}
	for _, v := range q.Files {
		// 根据url导入excel处理
		excel.NewExcel().ImportWithUrl(v, handlerImpStruct)
	}
	if len(handlerImpStruct.Err) > 0 {
		return false, handlerImpStruct.Err
	}
	return true, handlerImpStruct.Msg
}

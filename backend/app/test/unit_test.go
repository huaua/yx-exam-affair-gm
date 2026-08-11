package test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"yx-exam-affair-gm/app/biz/bizDao/bizDaoImpl"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/IOFile"
	"github.com/sealsee/web-base/public/basemodel"
	"github.com/sealsee/web-base/public/ds"
	"github.com/sealsee/web-base/public/setting"
	"github.com/sealsee/web-base/public/utils"
)

func TestCondition(t *testing.T) {
	type QueryDTO struct {
		Id     int64
		Name   string
		Age    int
		Sex    int
		Score  float32
		Score2 float64
	}
	queryDto := &QueryDTO{Name: "kkk", Sex: 1, Score: 0.00000, Score2: 0.0000}
	// 获取表达式参数
	condition, args := (&basemodel.BaseEntityQuery{}).
		AddEq("deleted", 1).
		AddEq("id", queryDto.Id).
		AddLikeAll("name", queryDto.Name).
		AddEq("age", queryDto.Age).
		AddEq("sex", queryDto.Sex).
		AddEq("score", queryDto.Score).
		AddEq("score2", queryDto.Score2).
		GetConditionArgs()
	fmt.Println(condition)
	fmt.Println(args...)
}

func TestErrFormat(t *testing.T) {
	// err := hint.BIZ_TCH_TASK_DETAIL_INFO_ERR.Format("2024-09-09", "xx学院", 0)
	// fmt.Print(err)

	curPage, totalPage := 1, 13
	percent := strconv.FormatFloat(float64(curPage)/float64(totalPage)*100, 'f', 2, 64)
	fmt.Println(percent)
}

func TestStr(t *testing.T) {
	col := ""
	gormTag := "primaryKey;column:BaoKaoID"
	gormTags := strings.Split(gormTag, ";")
	for i, v := range gormTags {
		fmt.Printf("i=%v, val=%v \n", i, v)
		if strings.Contains(v, "column:") {
			col = strings.Replace(v, "column:", "", -1)
			fmt.Println(col)
		}
	}
}

// 测试框架'容器'加载；测试复杂查询条件使用
func TestCtxFunc(t *testing.T) {
	// 加载容器
	setting.Init("../../config/config.yaml")
	_, cleanup, _ := ds.InitCompent(setting.Conf.Datasource)
	defer cleanup()
	utils.Init()
	IOFile.Init()

	// 容器实例
	tchTaskReportDao := bizDaoImpl.NewEaTchTaskReportDao()

	// 创建新的查询对象
	q := &bizDto.EaTchTaskReportQuery{}
	// q.AddIn("task_time", "2024-09-04 00:00:00", "2024-09-05 00:00:00")

	timeList := make([]any, 0)
	timeList = append(timeList, "2024-09-04 00:00:00", "2024-09-05 00:00:00")

	q.AddIn("task_time", timeList...)

	// count := tchTaskReportDao.Count(q)
	// fmt.Printf("-----> count: %v", count)

	list := tchTaskReportDao.List(q, q.GetPage())

	for _, d := range list {
		fmt.Printf("----->list1: 日期[%v], 名称[%v] \n", d.TaskTime, d.TchName)
	}

	// count2 := seDepartDao.CountWithCondition(q, "create_time > ? and update_time < ?", "2023-06-01 00:00:00", "2024-07-01 00:00:00")
	// fmt.Printf("-----> count2: %v", count2)

	// list2 := seDepartDao.ListWithCondition(q, p, "create_time > ? and update_time < ?", "2023-06-01 00:00:00", "2024-07-01 00:00:00")
	// for _, d := range list2 {
	// 	fmt.Printf("----->list2: 部门代码[%v], 名称[%v] \n", d.DepartCode, d.DepartName)
	// }

	// listmap3 := seDepartDao.ListMapWithCondition(q, p, "depart_code <> ? AND create_time > ? OR update_time < ?", "03", "2023-06-01 00:00:00", "2024-07-01 00:00:00")
	// for idx, mm := range listmap3 {
	// 	fmt.Printf("----->listmap3: [%v] 部门代码[%v], 名称[%v] \n", idx, mm["depart_code"], mm["depart_name"])
	// }

	// sysConfigService := userSerImpl.NewSysConfigService()
	// qConfig := &userDto.SysConfigQuery{}
	// qConfig.ConfigKey = "biz"
	// pConfig := qConfig.GetPage()
	// listConfig := sysConfigService.ListSysConfig(qConfig, pConfig)
	// for _, c := range listConfig {
	// 	fmt.Printf("----->listConfig: key[%v], name[%v] \n", c.ConfigKey, c.ConfigName)
	// }

	// sysConfigDao := userDaoImpl.NewSysConfigDao()
	// qConfig2 := &userDto.SysConfigQuery{}
	// qConfig2.AddNot("config_id", 2)
	// countConfig2 := sysConfigDao.CountWithCondition(qConfig2, "(config_name = ? OR config_key = ?)", "aa", "sys.index.skinName")
	// fmt.Printf("----->listConfig2 count: [%v]\n", countConfig2)

}

func TestRawSql(t *testing.T) {
	// 加载容器
	setting.Init("../../config/config.yaml")
	_, cleanup, _ := ds.InitCompent(setting.Conf.Datasource)
	defer cleanup()
	utils.Init()
	IOFile.Init()

	// var roomTaskReportService bizSerImpl.EaRoomTaskReportService
	// q := new(bizDto.EaRoomTaskReportQuery)
	// // q.RoomName = "梦园教学"
	// q.State = "2"
	// roomTaskReportList := roomTaskReportService.ListAllEaTaskRoomInfo(q)
	// for _, r := range roomTaskReportList {
	// 	fmt.Printf("----->raw res: [%v],[%v],[%v],[%v],[%v],[%v],[%v],[%v] \n", r.RoomId, r.RoomName, r.State, r.CampusName, r.LevelCode, r.GroupNum, r.GroupCapacity, r.Capacity)
	// }

	// sql := "select * from se_depart where deleted = ?"
	// result := query.RawSqlQueryList[bizDo.SeDepart](sql, 1, 10, 10)
	// 带分页
	// sql := "select * from sys_user where deleted = ? limit ? offset ?"
	// result := query.RawSqlQueryList[userDo.SysUser](sql, 1, 10, 10)
	// for _, d := range result {
	// 	fmt.Printf("----->raw res: 账号[%v], 姓名[%v] \n", d.UserName, d.NickName)
	// }

	// 复杂sql
	type UserDTO struct {
		UserName   string
		NickName   string
		DeptId     int64
		DeptName   string
		CreateTime basemodel.BaseTime
	}

	// sql2 := `select a.user_name,a.nick_name,a.dept_id,b.dept_name
	// 		from sys_user a left join ea_dept_info b on a.dept_id = b.dept_id
	// 		where `
	// 获取表达式参数
	user := &UserDTO{UserName: "kkk"}
	condition, args := (&basemodel.BaseEntityQuery{}).
		AddEq("a.deleted", 1).
		AddEq("create_time", user.CreateTime).
		AddLikeAll("a.nick_name", user.UserName).
		GetConditionArgs()
	fmt.Printf("%v, %v", condition, args)

	// p2 := page.NewPage()
	// p2.CurPage = 1
	// result2 := query.RawSqlQueryListWithPage[UserDTO](p2, sql2+condition, args...)
	// fmt.Printf("----->raw res2: 总行数: [%v], 总页数: [%v] \n", p2.TotalSize, p2.TotalPage)
	// for _, d := range result2 {
	// 	fmt.Printf("----->raw res2: 账号[%v],[%v], 部门[%v],[%v] \n", d.UserName, d.NickName, d.DeptId, d.DeptName)
	// }

	// 写表
	// updateSql := `update sys_user set remark = 'AAA' where user_id = ?`
	// updateSql := `update ea_classroom_info set capacity = capacity + 1 where room_id = ?`
	// tx.ExecGTx(func(gtx tx.Db) { gtx.Raw(updateSql, 9).Scan(nil) })

	// sql3 := `select count(1) from ea_room_task_report a
	// 		left join ea_room_info b on b.deleted = 1 and a.room_id = b.room_id
	// 		where a.deleted = 1 and a.task_id = ? and a.dept_id = ?`
	// enListNum := query.RawSqlQuery[int](sql3, 1, -100)
	// fmt.Printf("-----> raw enListNum: %v \n", enListNum)

	// 查询可使用的考场列表
	// sql4 := `select a.*, b.state as publish_state from ea_room_info a
	// 		left join ea_room_task_report b on b.deleted = 1 and a.room_id = b.room_id and b.task_id = ? and b.dept_id = ?
	// 		where a.deleted = 1 and a.dept_id = ? and a.level_code = ?
	// 		and a.room_id not in (
	// 			select c.room_id from ea_room_task_report c
	// 			where c.deleted = 1 and c.task_id <> ? and c.task_day = ?
	// 		)`
	// p4 := page.NewPage()
	// roomList := query.RawSqlQueryListWithPage[bizDto.EaRoomEnlistDto](p4, sql4, 1, -100, -100, 2, 1, "2024-12-25 00:00:00")
	// for _, room := range roomList {
	// 	fmt.Printf("----->raw roomList: 考场[%v],[%v],[%v],[%v] \n", room.RoomId, room.RoomName, room.DeptId, room.State)
	// }
	// fmt.Printf("-----> raw roomList: total: %v", p4.TotalSize)

}

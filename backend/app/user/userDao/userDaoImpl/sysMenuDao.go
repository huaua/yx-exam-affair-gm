package userDaoImpl

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/ds"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/query"
	"github.com/sealsee/web-base/public/ds/tx"
)

var sysMenuDaoImpl *SysMenuDaoImpl

func init() {
	sysMenuDaoImpl = &SysMenuDaoImpl{}
}

type SysMenuDaoImpl struct{}

func NewSysMenuDao() *SysMenuDaoImpl {
	return sysMenuDaoImpl
}

func (menu *SysMenuDaoImpl) GetById(id int64) *userDo.SysMenu {
	if id < 1 {
		return nil
	}
	data := new(userDo.SysMenu)
	where := new(userDo.SysMenu)
	where.MenuId = id
	where.Deleted = common.Normal
	res := ds.GetGDB().Find(data, where)
	if res.Error != nil {
		panic(res.Error)
	}
	if res.RowsAffected < 1 {
		return nil
	}
	return data
}

func (menu *SysMenuDaoImpl) Count(menuQuery *userDto.SysMenuQuery) int {
	menuQuery.Deleted = common.Normal
	return query.ExecGetQueryCount[userDto.SysMenuQuery, userDo.SysMenu](menuQuery)
}

func (menu *SysMenuDaoImpl) CountWithCondition(menuQuery *userDto.SysMenuQuery, condition string, args ...interface{}) int {
	menuQuery.Deleted = common.Normal
	return query.ExecGetQueryCountWithCondition[userDo.SysMenu](menuQuery, condition, args...)
}

func (menu *SysMenuDaoImpl) List(menuQuery *userDto.SysMenuQuery, page *page.Page) []*userDo.SysMenu {
	menuQuery.Deleted = common.Normal
	return query.ExecQueryList[userDto.SysMenuQuery, userDo.SysMenu](menuQuery, page)
}

func (menu *SysMenuDaoImpl) ListWithCondition(q *userDto.SysMenuQuery, page *page.Page, condition string, args ...interface{}) []*userDo.SysMenu {
	q.Deleted = common.Normal
	return query.ExecQueryListWithCondition[userDo.SysMenu](q, page, condition, args...)
}

func (menu *SysMenuDaoImpl) ListMenuTreeAll() []*userDo.SysMenu {
	sqlStr := `SELECT menu_id,menu_name,parent_id,order_num,path,component,query,is_frame,is_cache,menu_type,visible,state,perms,icon,create_by,create_time,update_by,update_time,remark FROM sys_menu WHERE deleted = 1 AND menu_type IN ('M', 'C') AND state = 1 ORDER BY parent_id, order_num`
	var data []*userDo.SysMenu
	rlt := ds.GetGDB().Raw(sqlStr).Scan(&data)
	if rlt.RowsAffected <= 0 {
		return nil
	}
	return data
}

func (menu *SysMenuDaoImpl) Insert(gtx tx.GTx, menuDO *userDo.SysMenu) bool {
	menuDO.Deleted = common.Normal
	res := gtx.Create(menuDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (menu *SysMenuDaoImpl) BatchInsert(gtx tx.GTx, menuList []*userDo.SysMenu) bool {
	res := gtx.CreateInBatches(menuList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (menu *SysMenuDaoImpl) Update(gtx tx.GTx, uData *userDo.SysMenu, uWhere *userDo.SysMenu) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (menu *SysMenuDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(userDo.SysMenu)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(userDo.SysMenu)
	uWhere.MenuId = id
	uWhere.Deleted = common.Normal
	return menu.Update(gtx, uData, uWhere)
}

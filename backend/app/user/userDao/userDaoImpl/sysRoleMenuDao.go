package userDaoImpl

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/query"
	"github.com/sealsee/web-base/public/ds/tx"
)

var sysRoleMenuDaoImpl *SysRoleMenuDaoImpl

func init() {
	sysRoleMenuDaoImpl = &SysRoleMenuDaoImpl{}
}

type SysRoleMenuDaoImpl struct{}

func NewSysRoleMenuDao() *SysRoleMenuDaoImpl {
	return sysRoleMenuDaoImpl
}

func (roleMenu *SysRoleMenuDaoImpl) Count(q *userDto.SysRoleMenuQuery) int {
	return query.ExecGetQueryCount[userDto.SysRoleMenuQuery, userDo.SysRoleMenu](q)
}

func (roleMenu *SysRoleMenuDaoImpl) List(q *userDto.SysRoleMenuQuery, page *page.Page) []*userDo.SysRoleMenu {
	return query.ExecQueryList[userDto.SysRoleMenuQuery, userDo.SysRoleMenu](q, page)
}

func (roleMenu *SysRoleMenuDaoImpl) ListWithCondition(q *userDto.SysRoleMenuQuery, page *page.Page, condition string, args ...interface{}) []*userDo.SysRoleMenu {
	return query.ExecQueryListWithCondition[userDo.SysRoleMenu](q, page, condition, args...)
}

func (roleMenu *SysRoleMenuDaoImpl) CheckMenuExistRole(menuId int64) int {
	if menuId < 1 {
		return 0
	}
	q := new(userDto.SysRoleMenuQuery)
	q.MenuId = menuId
	return query.ExecGetQueryCount[userDto.SysRoleMenuQuery, userDo.SysRoleMenu](q)
}

func (roleMenu *SysRoleMenuDaoImpl) BatchInsertRoleMenu(gtx tx.GTx, roleMenuList []*userDo.SysRoleMenu) int {
	res := gtx.CreateInBatches(roleMenuList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (roleMenu *SysRoleMenuDaoImpl) DeleteRoleMenuByRoleId(gtx tx.GTx, roleId int64) int {
	rlt := gtx.Where("role_id = ?", roleId).Delete(&userDo.SysRoleMenu{})
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return int(rlt.RowsAffected)
}

func (roleMenu *SysRoleMenuDaoImpl) BatchDeleteRoleMenu(gtx tx.GTx, roleIds []int64) int {
	rlt := gtx.Where("role_id IN ?", roleIds).Delete(&userDo.SysRoleMenu{})
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return int(rlt.RowsAffected)
}

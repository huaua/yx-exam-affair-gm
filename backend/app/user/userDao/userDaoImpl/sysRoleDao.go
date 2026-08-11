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

var sysRoleDaoImpl *SysRoleDaoImpl

func init() {
	sysRoleDaoImpl = &SysRoleDaoImpl{}
}

type SysRoleDaoImpl struct{}

func NewSysRoleDao() *SysRoleDaoImpl {
	return sysRoleDaoImpl
}

func (role *SysRoleDaoImpl) GetById(id int64) *userDo.SysRole {
	if id < 1 {
		return nil
	}
	data := new(userDo.SysRole)
	where := new(userDo.SysRole)
	where.RoleId = id
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

func (role *SysRoleDaoImpl) Count(roleQuery *userDto.SysRoleQuery) int {
	roleQuery.Deleted = common.Normal
	return query.ExecGetQueryCount[userDto.SysRoleQuery, userDo.SysRole](roleQuery)
}

func (r *SysRoleDaoImpl) CountWithCondition(role *userDto.SysRoleQuery, condition string, args ...interface{}) int {
	role.Deleted = common.Normal
	return query.ExecGetQueryCountWithCondition[userDo.SysRole](role, condition, args...)
}

func (role *SysRoleDaoImpl) ListAll() []*userDo.SysRole {
	roleQuery := &userDto.SysRoleQuery{}
	roleQuery.Deleted = common.Normal
	roleQuery.PageSize = page.MAX_PAGE_SIZE
	roleQuery.AddOrderAsc("role_sort")
	return query.ExecQueryList[userDto.SysRoleQuery, userDo.SysRole](roleQuery, roleQuery.GetPage())
}

func (role *SysRoleDaoImpl) List(roleQuery *userDto.SysRoleQuery, page *page.Page) []*userDo.SysRole {
	roleQuery.Deleted = common.Normal
	return query.ExecQueryList[userDto.SysRoleQuery, userDo.SysRole](roleQuery, page)
}

func (role *SysRoleDaoImpl) ListWithCondition(roleQuery *userDto.SysRoleQuery, page *page.Page, c string, args ...interface{}) []*userDo.SysRole {
	roleQuery.Deleted = common.Normal
	return query.ExecQueryListWithCondition[userDo.SysRole](roleQuery, page, c, args...)
}

func (role *SysRoleDaoImpl) Insert(gtx tx.GTx, roleDO *userDo.SysRole) bool {
	roleDO.Deleted = common.Normal
	res := gtx.Create(roleDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (role *SysRoleDaoImpl) BatchInsert(gtx tx.GTx, roleList []*userDo.SysRole) bool {
	res := gtx.CreateInBatches(roleList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (role *SysRoleDaoImpl) Update(gtx tx.GTx, uData *userDo.SysRole, uWhere *userDo.SysRole) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (role *SysRoleDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(userDo.SysRole)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(userDo.SysRole)
	uWhere.RoleId = id
	uWhere.Deleted = common.Normal
	return role.Update(gtx, uData, uWhere)
}

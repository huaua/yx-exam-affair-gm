package userDaoImpl

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/query"
	"github.com/sealsee/web-base/public/ds/tx"
)

var sysUserRoleDaoImpl *SysUserRoleDaoImpl

func init() {
	sysUserRoleDaoImpl = &SysUserRoleDaoImpl{}
}

type SysUserRoleDaoImpl struct{}

func NewSysUserRoleDao() *SysUserRoleDaoImpl {
	return sysUserRoleDaoImpl
}

func (userRole *SysUserRoleDaoImpl) List(userRoleQuery *userDto.SysUserRoleQuery, page *page.Page) []*userDo.SysUserRole {
	return query.ExecQueryList[userDto.SysUserRoleQuery, userDo.SysUserRole](userRoleQuery, page)
}

func (userRole *SysUserRoleDaoImpl) CountUserRoleByRoleId(roleId int64) int {
	if roleId < 1 {
		return 0
	}
	q := new(userDto.SysUserRoleQuery)
	q.RoleId = roleId
	return query.ExecGetQueryCount[userDto.SysUserRoleQuery, userDo.SysUserRole](q)
}

func (userRole *SysUserRoleDaoImpl) BatchInsertUserRole(gtx tx.GTx, userRoleList []*userDo.SysUserRole) int {
	res := gtx.CreateInBatches(userRoleList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return int(res.RowsAffected)
}

func (userRole *SysUserRoleDaoImpl) DeleteUserRoleInfo(gtx tx.GTx, userId int64, roleId int64) int {
	rlt := gtx.Where("user_id = ? AND role_id = ?", userId, roleId).Delete(&userDo.SysUserRole{})
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return int(rlt.RowsAffected)
}

func (userRole *SysUserRoleDaoImpl) DeleteUserRoleByUserId(gtx tx.GTx, userId int64) int {
	rlt := gtx.Where("user_id = ?", userId).Delete(&userDo.SysUserRole{})
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return int(rlt.RowsAffected)
}

func (userRole *SysUserRoleDaoImpl) BatchDeleteUserRole(gtx tx.GTx, userIds []int64) int {
	rlt := gtx.Where("user_id IN ?", userIds).Delete(&userDo.SysUserRole{})
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return int(rlt.RowsAffected)
}

func (userRole *SysUserRoleDaoImpl) DeleteUserRoleInfos(gtx tx.GTx, roleId int64, userIds []int64) int {
	rlt := gtx.Where("role_id = ? AND user_id IN ?", roleId, userIds).Delete(&userDo.SysUserRole{})
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return int(rlt.RowsAffected)
}

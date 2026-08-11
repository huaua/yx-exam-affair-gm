package userSerImpl

import (
	"strings"
	"yx-exam-affair-gm/app/hint"
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"
	"yx-exam-affair-gm/app/user/userDao"
	"yx-exam-affair-gm/app/user/userDao/userDaoImpl"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

type SysRoleService struct{}

var sysRoleService *SysRoleService
var sysRoleDao userDao.ISysRoleDao
var sysUserRoleDao userDao.ISysUserRoleDao

func init() {
	sysRoleService = &SysRoleService{}
	sysRoleDao = userDaoImpl.NewSysRoleDao()
	sysUserRoleDao = userDaoImpl.NewSysUserRoleDao()
}

func NewSysRoleService() *SysRoleService {
	return sysRoleService
}

func (role *SysRoleService) GetSysRoleById(id int64) *userDo.SysRole {
	return sysRoleDao.GetById(id)
}

func (role *SysRoleService) GetRoleIdsByUserId(userId int64) (roleIds []int64) {
	q := &userDto.SysUserRoleQuery{UserId: userId}
	q.PageSize = page.MAX_PAGE_SIZE
	userRoleList := sysUserRoleDao.List(q, q.GetPage())
	for _, ur := range userRoleList {
		roleIds = append(roleIds, ur.RoleId)
	}
	return
}

func (role *SysRoleService) GetRolePermission(userId int64) (roles []string) {
	// 超级管理员，暂时用固定ID=1判断，后续考虑安全性
	if userId == 1 {
		roles = append(roles, "admin")
	} else {
		// user_role 根据userId查roleIds
		q := &userDto.SysUserRoleQuery{UserId: userId}
		q.PageSize = page.MAX_PAGE_SIZE
		userRoleList := sysUserRoleDao.List(q, q.GetPage())
		roleIds := make([]int64, 0)
		for _, ur := range userRoleList {
			roleIds = append(roleIds, ur.RoleId)
		}
		// 查出该用户的所有角色串
		qRole := &userDto.SysRoleQuery{}
		qRole.PageSize = page.MAX_PAGE_SIZE
		roleList := sysRoleDao.ListWithCondition(qRole, qRole.GetPage(), "role_id IN ?", roleIds)
		for _, r := range roleList {
			roles = append(roles, strings.Split(strings.TrimSpace(r.RoleKey), ",")...)
		}
	}
	return
}

func (role *SysRoleService) GetSysRoleByUserId(userId int64) (roles []*userDo.SysRole) {
	// 超级管理员，暂时用固定ID=1判断，后续考虑安全性
	if userId == 1 {
		roles = sysRoleDao.ListAll()
	} else {
		// user_role 根据userId查roleIds
		q := &userDto.SysUserRoleQuery{UserId: userId}
		q.PageSize = page.MAX_PAGE_SIZE
		userRoleList := sysUserRoleDao.List(q, q.GetPage())
		roleIds := make([]int64, 0)
		for _, ur := range userRoleList {
			roleIds = append(roleIds, ur.RoleId)
		}
		// 查出该用户的所有角色串
		qRole := &userDto.SysRoleQuery{}
		qRole.PageSize = page.MAX_PAGE_SIZE
		qRole.AddOrderAsc("role_sort")
		roles = sysRoleDao.ListWithCondition(qRole, qRole.GetPage(), "role_id IN ?", roleIds)
	}
	return
}

func (role *SysRoleService) ListSysRole(roleQuery *userDto.SysRoleQuery, page *page.Page) []*userDo.SysRole {
	return sysRoleDao.List(roleQuery, page)
}

func (role *SysRoleService) InsertSysRole(r *userDo.SysRole) (success bool) {
	// 校验roleName roleKey是否重复
	q := &userDto.SysRoleQuery{}
	if sysRoleDao.CountWithCondition(q, "role_name = ? OR role_key = ?", r.RoleName, r.RoleKey) > 0 {
		panic(hint.SYS_ROLE_EXISTS)
	}
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysRoleDao.Insert(gtx, r)
		if success && len(r.MenuIds) > 0 {
			var roleMenuList []*userDo.SysRoleMenu
			for _, menuId := range r.MenuIds {
				roleMenuList = append(roleMenuList, &userDo.SysRoleMenu{RoleId: r.RoleId, MenuId: menuId})
			}
			success = sysRoleMenuDao.BatchInsertRoleMenu(gtx, roleMenuList) > 0
		}
	})
	return
}

func (role *SysRoleService) UpdateSysRoleById(roleDO *userDo.SysRole) (success bool) {
	uWhere := &userDo.SysRole{RoleId: roleDO.RoleId}
	roleDO.RoleId = 0
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysRoleDao.Update(gtx, roleDO, uWhere)
	})
	return
}

func (r *SysRoleService) UpdateSysRoleAndMenu(role *userDo.SysRole) (success bool) {
	roleId := role.RoleId
	// 校验roleName roleKey是否重复
	q := &userDto.SysRoleQuery{}
	if sysRoleDao.CountWithCondition(q, "role_id <> ? AND (role_name = ? OR role_key = ?)", roleId, role.RoleName, role.RoleKey) > 0 {
		panic(hint.SYS_ROLE_EXISTS)
	}
	// 处理角色菜单列表
	var roleMenuList []*userDo.SysRoleMenu
	if len(role.MenuIds) > 0 {
		for _, menuId := range role.MenuIds {
			roleMenuList = append(roleMenuList, &userDo.SysRoleMenu{RoleId: roleId, MenuId: menuId})
		}
	}
	// 处理角色信息
	uWhere := &userDo.SysRole{RoleId: roleId}
	role.RoleId = 0 // 更新字段去掉roleId
	// 数据库执行事务
	tx.ExecGTx(func(gtx tx.GTx) {
		// 更新用户角色，先删再加
		sysRoleMenuDao.DeleteRoleMenuByRoleId(gtx, roleId)
		if len(roleMenuList) > 0 {
			sysRoleMenuDao.BatchInsertRoleMenu(gtx, roleMenuList)
		}
		// 更新role
		success = sysRoleDao.Update(gtx, role, uWhere)
	})
	return
}

func (role *SysRoleService) UpdateRoleState(roleId int64, state string, oper int64) (success bool) {
	uWhere := &userDo.SysRole{RoleId: roleId}
	uSet := &userDo.SysRole{State: state}
	uSet.SetUpdateBy(oper)
	// 数据库执行事务
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysRoleDao.Update(gtx, uSet, uWhere)
	})
	return
}

func (role *SysRoleService) DeleteSysRoleById(id int64, deleteBy int64) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysRoleDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

package userDao

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

type ISysUserRoleDao interface {
	List(userRoleQuery *userDto.SysUserRoleQuery, page *page.Page) []*userDo.SysUserRole
	// 通过角色ID查询角色使用数量
	CountUserRoleByRoleId(roleId int64) int
	// 批量新增用户角色信息
	BatchInsertUserRole(gtx tx.GTx, userRoleList []*userDo.SysUserRole) int
	// 删除用户和角色关联信息
	DeleteUserRoleInfo(gtx tx.GTx, userId int64, roleId int64) int
	// 通过用户ID删除用户和角色关联
	DeleteUserRoleByUserId(gtx tx.GTx, userId int64) int
	// 批量删除用户和角色关联
	BatchDeleteUserRole(gtx tx.GTx, userIds []int64) int
	// 批量取消授权用户角色
	DeleteUserRoleInfos(gtx tx.GTx, roleId int64, userIds []int64) int
}

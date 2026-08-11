package userSerImpl

import (
	"yx-exam-affair-gm/app/biz/bizDao"
	"yx-exam-affair-gm/app/biz/bizDao/bizDaoImpl"
	"yx-exam-affair-gm/app/hint"
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"
	"yx-exam-affair-gm/app/user/userDao"
	"yx-exam-affair-gm/app/user/userDao/userDaoImpl"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
	"github.com/sealsee/web-base/public/utils/crypto"
	"github.com/sealsee/web-base/public/utils/file/excel"
)

type SysUserService struct{}

var sysUserService *SysUserService
var sysUserDao userDao.ISysUserDao

var eaDeptInfoDao bizDao.IEaDeptInfoDao

func init() {
	sysUserService = &SysUserService{}
	sysUserDao = userDaoImpl.NewSysUserDao()
	eaDeptInfoDao = bizDaoImpl.NewEaDeptInfoDao()
}

func NewSysUserService() *SysUserService {
	return sysUserService
}

func (user *SysUserService) GetSysUserById(id int64) *userDo.SysUser {
	return sysUserDao.GetById(id)
}

func (user *SysUserService) GetSysUserByName(username string) *userDo.SysUser {
	return sysUserDao.GetByName(username)
}

func (user *SysUserService) ListSysUser(userQuery *userDto.SysUserQuery, page *page.Page) []*userDo.SysUser {

	userList := sysUserDao.List(userQuery, page)

	// 增加学院名称
	deptMap := make(map[int64]string)
	for _, user := range userList {
		if user.DeptId > 0 {
			if deptMap[user.DeptId] != "" {
				user.DeptName = deptMap[user.DeptId]
			} else {
				dept := eaDeptInfoDao.GetById(user.DeptId)
				if dept != nil {
					user.DeptName = dept.DeptName
					deptMap[user.DeptId] = dept.DeptName
				}
			}
		}
	}

	return userList
}

func (u *SysUserService) InsertSysUser(user *userDo.SysUser) (success bool) {
	q := &userDto.SysUserQuery{UserName: user.UserName}
	if sysUserDao.Count(q) > 0 {
		panic(hint.SYS_USER_EXISTS)
	}
	user.Password = crypto.GeneralPwd(user.Password)
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysUserDao.Insert(gtx, user)
		if success && len(user.RoleIds) > 0 {
			var userRoleList []*userDo.SysUserRole
			for _, roleId := range user.RoleIds {
				userRoleList = append(userRoleList, &userDo.SysUserRole{UserId: user.UserId, RoleId: roleId})
			}
			success = sysUserRoleDao.BatchInsertUserRole(gtx, userRoleList) > 0
		}
	})
	return
}

func (u *SysUserService) UpdateSysUserAndRole(user *userDo.SysUser) (success bool) {
	userId := user.UserId
	if user.UserName != "" {
		// 判断username是否重复
		q := &userDto.SysUserQuery{UserName: user.UserName}
		if sysUserDao.CountWithCondition(q, "user_id <> ?", userId) > 0 {
			panic(hint.SYS_USER_EXISTS)
		}
	}
	// 处理用户角色列表
	var userRoleList []*userDo.SysUserRole
	if len(user.RoleIds) > 0 {
		for _, roleId := range user.RoleIds {
			userRoleList = append(userRoleList, &userDo.SysUserRole{UserId: userId, RoleId: roleId})
		}
	}
	// 处理用户信息
	uWhere := &userDo.SysUser{UserId: userId}
	user.UserId = 0 // 更新字段去掉userId
	// 数据库执行事务
	tx.ExecGTx(func(gtx tx.GTx) {
		// 更新用户角色，先删再加
		if len(userRoleList) > 0 {
			sysUserRoleDao.DeleteUserRoleByUserId(gtx, userId)
			sysUserRoleDao.BatchInsertUserRole(gtx, userRoleList)
		}
		// 更新user
		success = sysUserDao.Update(gtx, user, uWhere)
	})
	return
}

// 只更新用户表数据
func (u *SysUserService) UpdateUserSingle(user *userDo.SysUser) (success bool) {
	userId := user.UserId
	uWhere := &userDo.SysUser{UserId: userId}
	user.UserId = 0 // 更新字段去掉userId
	// 数据库执行事务
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysUserDao.Update(gtx, user, uWhere)
	})
	return
}

func (u *SysUserService) ResetPwd(userId int64, password string, oper int64) (success bool) {
	uWhere := &userDo.SysUser{UserId: userId}
	uSet := &userDo.SysUser{Password: crypto.GeneralPwd(password)}
	uSet.SetUpdateBy(oper)
	// 数据库执行事务
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysUserDao.Update(gtx, uSet, uWhere)
	})
	return
}

func (u *SysUserService) UpdateUserState(userId int64, state string, oper int64) (success bool) {
	uWhere := &userDo.SysUser{UserId: userId}
	uSet := &userDo.SysUser{State: state}
	uSet.SetUpdateBy(oper)
	// 数据库执行事务
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysUserDao.Update(gtx, uSet, uWhere)
	})
	return
}

func (u *SysUserService) UpdateUserProfile(user *userDo.SysUser) (success bool) {
	userId := user.UserId
	if user.UserName != "" {
		// 判断username是否重复
		q := &userDto.SysUserQuery{UserName: user.UserName}
		if sysUserDao.CountWithCondition(q, "user_id <> ?", userId) > 0 {
			panic(hint.SYS_USER_EXISTS)
		}
	}
	// 处理用户信息
	uWhere := &userDo.SysUser{UserId: userId}
	// 不要更新的字段
	user.UserId = 0
	user.Avatar = ""
	user.Password = ""
	user.DeptId = 0
	user.State = ""

	// 数据库执行事务
	tx.ExecGTx(func(gtx tx.GTx) {
		// 更新user
		success = sysUserDao.Update(gtx, user, uWhere)
	})
	return
}

func (u *SysUserService) UpdateUserAvatar(userId int64, avatar string, oper int64) (success bool) {
	// 处理用户信息
	uWhere := &userDo.SysUser{UserId: userId}
	uSet := &userDo.SysUser{Avatar: avatar}
	uSet.SetUpdateBy(oper)
	// 数据库执行事务
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysUserDao.Update(gtx, uSet, uWhere)
	})
	return
}

func (user *SysUserService) DeleteSysUserById(id int64, deleteBy int64) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysUserDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

// 注：仅支持1000条以内的简单同步导出 by codegen
func (user *SysUserService) ExportSysUser(userQuery *userDto.SysUserQuery) []byte {
	nPage := page.NewPage()
	nPage.PageSize = page.MAX_PAGE_SIZE
	dataList := sysUserDao.List(userQuery, nPage)
	row_title := []interface{}{"用户ID", "部门ID", "用户账号", "用户昵称", "用户类型", "用户邮箱", "手机号码", "用户性别", "头像地址", "密码", "帐号状态", "删除标志", "创建者", "创建时间", "更新者", "更新时间", "备注"}
	dataSlice := make([][]interface{}, 0)
	dataSlice = append(dataSlice, row_title)
	for _, v := range dataList {
		dataSlice = append(dataSlice, []interface{}{v.UserId, v.DeptId, v.UserName, v.NickName, v.UserType, v.Email, v.Phonenumber, v.Sex, v.Avatar, v.Password, v.State, v.Deleted, v.CreateBy, v.CreateTime, v.UpdateBy, v.UpdateTime, v.Remark})
	}
	return excel.ExportExcel(dataSlice)
}

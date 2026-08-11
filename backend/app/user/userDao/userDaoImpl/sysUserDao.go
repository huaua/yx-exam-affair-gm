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

var sysUserDaoImpl *SysUserDaoImpl

func init() {
	sysUserDaoImpl = &SysUserDaoImpl{}
}

type SysUserDaoImpl struct{}

func NewSysUserDao() *SysUserDaoImpl {
	return sysUserDaoImpl
}

func (user *SysUserDaoImpl) GetById(id int64) *userDo.SysUser {
	if id < 1 {
		return nil
	}
	data := new(userDo.SysUser)
	where := new(userDo.SysUser)
	where.UserId = id
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

func (user *SysUserDaoImpl) GetByName(username string) *userDo.SysUser {
	if username == "" {
		return nil
	}
	data := new(userDo.SysUser)
	where := new(userDo.SysUser)
	where.UserName = username
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

func (user *SysUserDaoImpl) Count(userQuery *userDto.SysUserQuery) int {
	userQuery.Deleted = common.Normal
	return query.ExecGetQueryCount[userDto.SysUserQuery, userDo.SysUser](userQuery)
}

func (user *SysUserDaoImpl) CountWithCondition(userQuery *userDto.SysUserQuery, condition string, args ...interface{}) int {
	userQuery.Deleted = common.Normal
	return query.ExecGetQueryCountWithCondition[userDo.SysUser](userQuery, condition, args...)
}

func (user *SysUserDaoImpl) List(userQuery *userDto.SysUserQuery, page *page.Page) []*userDo.SysUser {
	userQuery.Deleted = common.Normal
	return query.ExecQueryList[userDto.SysUserQuery, userDo.SysUser](userQuery, page)
}

func (user *SysUserDaoImpl) Insert(gtx tx.GTx, userDO *userDo.SysUser) bool {
	userDO.Deleted = common.Normal
	res := gtx.Create(userDO)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

// 注意：拼装list里deleted要设为1正常
func (user *SysUserDaoImpl) BatchInsert(gtx tx.GTx, userList []*userDo.SysUser) bool {
	res := gtx.CreateInBatches(userList, 100)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (user *SysUserDaoImpl) Update(gtx tx.GTx, uData *userDo.SysUser, uWhere *userDo.SysUser) bool {
	uWhere.Deleted = common.Normal
	res := gtx.Where(uWhere).Updates(uData)
	if res.Error != nil {
		panic(res.Error)
	}
	return res.RowsAffected > 0
}

func (user *SysUserDaoImpl) DeleteById(gtx tx.GTx, id int64, deleteBy int64) bool {
	uData := new(userDo.SysUser)
	uData.SetDeleteBy(deleteBy)
	uWhere := new(userDo.SysUser)
	uWhere.UserId = id
	uWhere.Deleted = common.Normal
	return user.Update(gtx, uData, uWhere)
}

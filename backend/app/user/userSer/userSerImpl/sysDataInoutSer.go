package userSerImpl

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"
	"yx-exam-affair-gm/app/user/userDao"
	"yx-exam-affair-gm/app/user/userDao/userDaoImpl"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
)

type SysDataInoutService struct{}

var sysDataInoutService *SysDataInoutService
var sysDataInoutDao userDao.ISysDataInoutDao

func init() {
	sysDataInoutService = &SysDataInoutService{}
	sysDataInoutDao = userDaoImpl.NewSysDataInoutDao()
}

func NewSysDataInoutService() *SysDataInoutService {
	return sysDataInoutService
}

func (dataInout *SysDataInoutService) GetSysDataInoutById(id int64) *userDo.SysDataInout {
	return sysDataInoutDao.GetById(id)
}

func (dataInout *SysDataInoutService) ListSysDataInout(dataInoutQuery *userDto.SysDataInoutQuery, page *page.Page) []*userDo.SysDataInout {
	return sysDataInoutDao.List(dataInoutQuery, page)
}

func (dataInout *SysDataInoutService) ListSysDataInoutWithCondition(dataInoutQuery *userDto.SysDataInoutQuery, page *page.Page, condition string, args ...interface{}) []*userDo.SysDataInout {
	return sysDataInoutDao.ListWithCondition(dataInoutQuery, page, condition, args...)
}

func (dataInout *SysDataInoutService) InsertSysDataInout(dataInoutDO *userDo.SysDataInout) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysDataInoutDao.Insert(gtx, dataInoutDO)
	})
	return
}

func (dataInout *SysDataInoutService) UpdateSysDataInoutById(dataInoutDO *userDo.SysDataInout) (success bool) {
	uWhere := &userDo.SysDataInout{Id: dataInoutDO.Id}
	dataInoutDO.Id = 0
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysDataInoutDao.Update(gtx, dataInoutDO, uWhere)
	})
	return
}

func (dataInout *SysDataInoutService) DeleteSysDataInoutById(id int64, deleteBy int64) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysDataInoutDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

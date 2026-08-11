package userSer

import (
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"

	"github.com/sealsee/web-base/public/ds/page"
)

type ISysConfigService interface {
	GetSysConfigById(id int64) *userDo.SysConfig
	ListSysConfig(config *userDto.SysConfigQuery, page *page.Page) []*userDo.SysConfig
	InsertSysConfig(config *userDo.SysConfig) bool
	UpdateSysConfigById(config *userDo.SysConfig) bool
	DeleteSysConfigById(id int64, deleteBy int64) bool
	ExportSysConfig(config *userDto.SysConfigQuery) []byte
	ExportSysConfigNew(config *userDto.SysConfigQuery) ([]byte, error)
	GetSysConfigByConfigKey(configKey string) *userDo.SysConfig
}

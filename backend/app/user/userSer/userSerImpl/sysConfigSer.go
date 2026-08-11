package userSerImpl

import (
	"yx-exam-affair-gm/app/hint"
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"
	"yx-exam-affair-gm/app/user/userDao"
	"yx-exam-affair-gm/app/user/userDao/userDaoImpl"

	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
	"github.com/sealsee/web-base/public/utils/file/excel"
)

type SysConfigService struct{}

var sysConfigService *SysConfigService
var sysConfigDao userDao.ISysConfigDao

func init() {
	sysConfigService = &SysConfigService{}
	sysConfigDao = userDaoImpl.NewSysConfigDao()
}

func NewSysConfigService() *SysConfigService {
	return sysConfigService
}

func (config *SysConfigService) GetSysConfigById(id int64) *userDo.SysConfig {
	return sysConfigDao.GetById(id)
}

func (config *SysConfigService) ListSysConfig(configQuery *userDto.SysConfigQuery, page *page.Page) []*userDo.SysConfig {
	configQuery.AddLikeAll("config_name", configQuery.ConfigName)
	configQuery.AddLikeAll("config_key", configQuery.ConfigKey)
	return sysConfigDao.List(configQuery, page)
}

func (config *SysConfigService) InsertSysConfig(configDO *userDo.SysConfig) (success bool) {
	// 校验 configName configKey 是否重复
	q := &userDto.SysConfigQuery{}
	if sysConfigDao.CountWithCondition(q, "config_name = ? OR config_key = ?", configDO.ConfigName, configDO.ConfigKey) > 0 {
		panic(hint.SYS_CONFIG_EXISTS)
	}
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysConfigDao.Insert(gtx, configDO)
	})
	return
}

func (config *SysConfigService) UpdateSysConfigById(configDO *userDo.SysConfig) (success bool) {
	if configDO.ConfigId < 1 {
		return false
	}
	// 校验 configName configKey 是否重复
	q := &userDto.SysConfigQuery{}
	if sysConfigDao.CountWithCondition(q, "config_id <> ? AND (config_name = ? OR config_key = ?)", configDO.ConfigId, configDO.ConfigName, configDO.ConfigKey) > 0 {
		panic(hint.SYS_CONFIG_EXISTS)
	}
	uWhere := &userDo.SysConfig{ConfigId: configDO.ConfigId}
	configDO.ConfigId = 0
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysConfigDao.Update(gtx, configDO, uWhere)
	})
	return
}

func (config *SysConfigService) DeleteSysConfigById(id int64, deleteBy int64) (success bool) {
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysConfigDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

// 注：仅支持1000条以内的简单同步导出 by codegen
func (config *SysConfigService) ExportSysConfig(configQuery *userDto.SysConfigQuery) []byte {
	nPage := page.NewPage()
	nPage.PageSize = page.MAX_PAGE_SIZE
	dataList := sysConfigDao.List(configQuery, nPage)
	row_title := []interface{}{"参数主键", "参数名称", "参数键名", "参数键值", "系统内置", "创建者", "备注"}
	dataSlice := make([][]interface{}, 0)
	dataSlice = append(dataSlice, row_title)
	for _, v := range dataList {
		dataSlice = append(dataSlice, []interface{}{v.ConfigId, v.ConfigName, v.ConfigKey, v.ConfigValue, v.ConfigType, v.Remark})
	}

	return excel.ExportExcel(dataSlice)
}

type HandlerStruct struct {
	Query *userDto.SysConfigQuery
}

func (h *HandlerStruct) Title() string {
	return "系统参数"
}

func (h *HandlerStruct) HeaderColumn() []string {
	return []string{"参数名,config_name", "参数键,config_key", "参数值,config_value"}
}

func (h *HandlerStruct) Rows(p *page.Page) []map[string]interface{} {
	return sysConfigDao.ListMap(h.Query, p)
}

func (h *HandlerStruct) Finish(url string) {
}

func (config *SysConfigService) ExportSysConfigNew(configQuery *userDto.SysConfigQuery) ([]byte, error) {
	return excel.NewExcel().ExportSync(&HandlerStruct{configQuery})
}

func (config *SysConfigService) GetSysConfigByConfigKey(configKey string) *userDo.SysConfig {
	// nPage := page.NewPage()
	// nPage.PageSize = page.EXPORT_PAGE_SIZE
	// configQuery := &userDto.SysConfigQuery{ConfigKey: cst.BIZ_ADMIT_YEAR}
	// list := sysConfigDao.List(configQuery, nPage)
	// if list != nil && len(list) == 1 {
	// 	return list[0]
	// }
	// // 没有数据和多条数据都返回nil
	return sysConfigDao.GetByKey(configKey)
}

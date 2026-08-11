package userDo

import "github.com/sealsee/web-base/public/basemodel"

type SysMenu struct {
	MenuId    int64      `gorm:"primaryKey" json:"menuId,string,omitempty"` // 主键 菜单ID
	MenuName  string     `json:"menuName,omitempty"`                        // 菜单名称
	ParentId  int64      `json:"parentId,string,omitempty"`                 // 父菜单ID
	OrderNum  int        `json:"orderNum,omitempty"`                        // 显示顺序
	Path      string     `json:"path,omitempty"`                            // 路由地址
	Component string     `json:"component,omitempty"`                       // 组件路径
	Query     string     `json:"query,omitempty"`                           // 路由参数
	IsFrame   int        `json:"isFrame,omitempty"`                         // 是否为外链（1-是 2-否）
	IsCache   int        `json:"isCache,omitempty"`                         // 是否缓存（1-缓存 2-不缓存）
	MenuType  string     `json:"menuType,omitempty"`                        // 菜单类型（M目录 C菜单 F按钮）
	Visible   int        `json:"visible,omitempty"`                         // 菜单状态（1-显示 2-隐藏）
	Perms     string     `json:"perms,omitempty"`                           // 权限标识
	Icon      string     `json:"icon,omitempty"`                            // 菜单图标
	State     string     `json:"state,omitempty"`                           // 菜单状态（1-正常 0-停用）
	Remark    string     `json:"remark,omitempty"`                          // 备注
	Children  []*SysMenu `gorm:"-" json:"children"`                         // 子菜单
	basemodel.BaseEntity
}

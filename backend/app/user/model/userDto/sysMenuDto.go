package userDto

import "github.com/sealsee/web-base/public/basemodel"

type SysMenuQuery struct {
	MenuId    int64  `gorm:"primaryKey" json:"menuId,string,omitempty"` // 主键 菜单ID
	MenuName  string `form:"menuName" json:"menuName,omitempty"`        // 菜单名称
	ParentId  int64  `form:"parentId" json:"parentId,omitempty"`        // 父菜单ID
	OrderNum  int    `form:"orderNum" json:"orderNum,omitempty"`        // 显示顺序
	Path      string `form:"path" json:"path,omitempty"`                // 路由地址
	Component string `form:"component" json:"component,omitempty"`      // 组件路径
	Query     string `form:"query" json:"query,omitempty"`              // 路由参数
	IsFrame   int    `form:"isFrame" json:"isFrame,omitempty"`          // 是否为外链（1是 2否）
	IsCache   int    `form:"isCache" json:"isCache,omitempty"`          // 是否缓存（1缓存 2不缓存）
	MenuType  string `form:"menuType" json:"menuType,omitempty"`        // 菜单类型（M目录 C菜单 F按钮）
	Visible   int    `form:"visible" json:"visible,omitempty"`          // 菜单状态（1显示 2隐藏）
	Perms     string `form:"perms" json:"perms,omitempty"`              // 权限标识
	Icon      string `form:"icon" json:"icon,omitempty"`                // 菜单图标
	State     string `form:"state" json:"state,omitempty"`              // 菜单状态（1正常 0停用）
	Remark    string `form:"remark" json:"remark,omitempty"`            // 备注
	basemodel.BaseEntityQuery
}

// 路由配置信息
type RouterVO struct {
	Name       string      `json:"name"`       // 路由名字
	Path       string      `json:"path"`       // 路由地址
	Hidden     bool        `json:"hidden"`     // 是否隐藏路由，当设置 true 的时候该路由不会再侧边栏出现
	Redirect   string      `json:"redirect"`   // 重定向地址，当设置 noRedirect 的时候该路由在面包屑导航中不可被点击
	Component  string      `json:"component"`  // 组件地址
	Query      string      `json:"query"`      // 路由参数：如 {"id": 1, "name": "aaa"}
	AlwaysShow bool        `json:"alwaysShow"` // 当你一个路由下面的 children 声明的路由大于1个时，自动会变成嵌套的模式--如组件页面
	Meta       *MetaVO     `json:"meta"`       // 其他元素
	Children   []*RouterVO `json:"children"`   // 子路由
}

// 路由显示信息
type MetaVO struct {
	Title   string `json:"title"`   // 设置该路由在侧边栏和面包屑中展示的名字
	Icon    string `json:"icon"`    // 设置该路由的图标，前端对应路径src/assets/icons/svg
	NoCache bool   `json:"noCache"` // 设置为true，则不会被 <keep-alive>缓存
	Link    string `json:"link"`    // 内链地址（http(s)://开头）
}

// Treeselect树结构
type TreeSelect struct {
	Id       int64         `json:"id"`
	Label    string        `json:"label"`
	Children []*TreeSelect `json:"children"`
}

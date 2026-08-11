package userSerImpl

import (
	"strings"
	"yx-exam-affair-gm/app/hint"
	"yx-exam-affair-gm/app/user/model/userDo"
	"yx-exam-affair-gm/app/user/model/userDto"
	"yx-exam-affair-gm/app/user/userDao"
	"yx-exam-affair-gm/app/user/userDao/userDaoImpl"

	"github.com/sealsee/web-base/public/cst/common"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/tx"
	"github.com/sealsee/web-base/public/utils/file/excel"
	"github.com/sealsee/web-base/public/utils/stringUtils"
)

type SysMenuService struct{}

var sysMenuService *SysMenuService
var sysMenuDao userDao.ISysMenuDao
var sysRoleMenuDao userDao.ISysRoleMenuDao

func init() {
	sysMenuService = &SysMenuService{}
	sysMenuDao = userDaoImpl.NewSysMenuDao()
	sysRoleMenuDao = userDaoImpl.NewSysRoleMenuDao()
}

func NewSysMenuService() *SysMenuService {
	return sysMenuService
}

func (menu *SysMenuService) GetSysMenuById(id int64) *userDo.SysMenu {
	return sysMenuDao.GetById(id)
}

func (menu *SysMenuService) ListSysMenu(menuQuery *userDto.SysMenuQuery, page *page.Page) []*userDo.SysMenu {
	return sysMenuDao.List(menuQuery, page)
}

func (menu *SysMenuService) ListSysMenuByUserId(userId int64) (list []*userDo.SysMenu) {
	// 系统管理员，直接返回所有菜单，暂时用id=1识别，后续考虑安全性 TODO
	if userId == 1 {
		list = listAllMenu()
	} else {
		// userId --> roleIds --> menu
		// user_role 根据userId查roleIds
		q := &userDto.SysUserRoleQuery{UserId: userId}
		q.PageSize = page.MAX_PAGE_SIZE
		userRoleList := sysUserRoleDao.List(q, q.GetPage())
		roleIds := make([]int64, 0)
		var hasAdminRole bool
		for _, ur := range userRoleList {
			// TODO 判断角色状态，如果停用则跳出
			if ur.RoleId == 1 {
				hasAdminRole = true
			}
			roleIds = append(roleIds, ur.RoleId)
		}
		if hasAdminRole {
			list = listAllMenu()
		} else {
			// role_menu 根据roleId查menuIds
			qRoleMenu := &userDto.SysRoleMenuQuery{}
			qRoleMenu.PageSize = page.MAX_PAGE_SIZE
			roleMenuList := sysRoleMenuDao.ListWithCondition(qRoleMenu, qRoleMenu.GetPage(), "role_id IN ?", roleIds)
			menuIds := make([]int64, 0)
			for _, rm := range roleMenuList {
				menuIds = append(menuIds, rm.MenuId)
			}
			// 根据menuId查list
			list = listAllMenuWithCondition("menu_id IN ?", menuIds)
		}
	}
	return
}

func (menu *SysMenuService) BuildMenuTreeSelect(menus []*userDo.SysMenu) (list []*userDto.TreeSelect) {
	return buildMenuTree(menus)
}

func (menu *SysMenuService) ListMenuTreeByUserId(userId int64) (list []*userDo.SysMenu, permissions []string) {
	// 系统管理员，直接返回所有菜单，暂时用id=1识别，后续考虑安全性
	if userId == 1 {
		list = listAllMenuWithCondition("state = ? AND menu_type IN ('M', 'C')", common.OK)
		permissions = append(permissions, "*:*:*")
	} else {
		// userId --> roleIds --> menu
		// user_role 根据userId查roleIds
		q := &userDto.SysUserRoleQuery{UserId: userId}
		q.PageSize = page.MAX_PAGE_SIZE
		userRoleList := sysUserRoleDao.List(q, q.GetPage())
		if len(userRoleList) > 0 {
			roleIds := make([]int64, 0)
			var hasAdminRole bool
			for _, ur := range userRoleList {
				// TODO 判断角色状态，如果停用则跳出
				if ur.RoleId == 1 {
					hasAdminRole = true
				}
				roleIds = append(roleIds, ur.RoleId)
			}
			if hasAdminRole {
				list = listAllMenuWithCondition("state = ? AND menu_type IN ('M', 'C')", common.OK)
				permissions = append(permissions, "*:*:*")
			} else {
				// role_menu 根据roleId查menuIds
				qRoleMenu := &userDto.SysRoleMenuQuery{}
				qRoleMenu.PageSize = page.MAX_PAGE_SIZE
				roleMenuList := sysRoleMenuDao.ListWithCondition(qRoleMenu, qRoleMenu.GetPage(), "role_id IN ?", roleIds)
				if len(roleMenuList) > 0 {
					menuIds := make([]int64, 0)
					for _, rm := range roleMenuList {
						menuIds = append(menuIds, rm.MenuId)
					}
					// 根据menuId查list
					mlist := listAllMenuWithCondition("state = ? AND menu_id IN ?", common.OK, menuIds)
					// 组装权限数组
					for _, m := range mlist {
						if strings.TrimSpace(m.Perms) != "" {
							permissions = append(permissions, strings.Split(strings.TrimSpace(m.Perms), ",")...)
						}
						// AND menu_type IN ('M', 'C')
						if m.MenuType == "M" || m.MenuType == "C" {
							list = append(list, m)
						}
					}
				}
			}
		}
	}
	if len(list) > 0 {
		list = getChildPerms(list, 0)
	}
	return
}

func (menu *SysMenuService) ListMenuIdByRoleId(roleId int64, menuCheckStrictly bool) []int64 {
	// role_menu 根据roleId查menuIds
	qRoleMenu := &userDto.SysRoleMenuQuery{RoleId: roleId}
	qRoleMenu.PageSize = page.MAX_PAGE_SIZE
	roleMenuList := sysRoleMenuDao.List(qRoleMenu, qRoleMenu.GetPage())
	menuIds := make([]int64, 0)
	for _, rm := range roleMenuList {
		menuIds = append(menuIds, rm.MenuId)
	}
	if menuCheckStrictly {
		// 根据menuId查list
		menus := listAllMenuWithCondition("menu_id IN ?", menuIds)
		parentIds := make(map[int64]int64, 0)
		for _, m := range menus {
			parentIds[m.ParentId] = m.MenuId
		}
		newMenuIds := make([]int64, 0)
		for _, i := range menuIds {
			if parentIds[i] < 1 {
				newMenuIds = append(newMenuIds, i)
			}
		}
		return newMenuIds
	}
	return menuIds
}

func listAllMenu() (list []*userDo.SysMenu) {
	qMenu := &userDto.SysMenuQuery{}
	qMenu.PageSize = page.MAX_PAGE_SIZE
	qMenu.AddOrderAsc("parent_id")
	qMenu.AddOrderAsc("order_num")
	list = sysMenuDao.List(qMenu, qMenu.GetPage())
	return
}

func listAllMenuWithCondition(condition string, args ...interface{}) (list []*userDo.SysMenu) {
	qMenu := &userDto.SysMenuQuery{}
	qMenu.PageSize = page.MAX_PAGE_SIZE
	qMenu.AddOrderAsc("parent_id")
	qMenu.AddOrderAsc("order_num")
	list = sysMenuDao.ListWithCondition(qMenu, qMenu.GetPage(), condition, args...)
	return
}

func (menu *SysMenuService) BuildRouters(menus []*userDo.SysMenu) (list []*userDto.RouterVO) {
	return buildMenuRouters(menus)
}

func (menu *SysMenuService) InsertSysMenu(menuDO *userDo.SysMenu) (success bool) {
	// 校验menuName是否重复
	q := &userDto.SysMenuQuery{MenuName: menuDO.MenuName}
	if sysMenuDao.Count(q) > 0 {
		panic(hint.SYS_MENU_EXISTS)
	}
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysMenuDao.Insert(gtx, menuDO)
	})
	return
}

func (menu *SysMenuService) UpdateSysMenuById(menuDO *userDo.SysMenu) (success bool) {
	if menuDO.MenuId == menuDO.ParentId {
		panic(hint.SYS_MENU_PARENTID_ERR)
	}
	// 校验menuName是否重复
	q := &userDto.SysMenuQuery{MenuName: menuDO.MenuName}
	if sysMenuDao.CountWithCondition(q, "menu_id <> ?", menuDO.MenuId) > 0 {
		panic(hint.SYS_MENU_EXISTS)
	}
	uWhere := &userDo.SysMenu{MenuId: menuDO.MenuId}
	menuDO.MenuId = 0
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysMenuDao.Update(gtx, menuDO, uWhere)
	})
	return
}

func (menu *SysMenuService) DeleteSysMenuById(id int64, deleteBy int64) (success bool) {
	// 查询是否有子菜单存在
	q := &userDto.SysMenuQuery{ParentId: id}
	if sysMenuDao.Count(q) > 0 {
		panic(hint.SYS_MENU_HAS_CHILD)
	}
	// 查询是否有角色分配该菜单
	qRoleMenu := &userDto.SysRoleMenuQuery{MenuId: id}
	if sysRoleMenuDao.Count(qRoleMenu) > 0 {
		panic(hint.SYS_MENU_HAS_ROLE)
	}
	tx.ExecGTx(func(gtx tx.GTx) {
		success = sysMenuDao.DeleteById(gtx, id, deleteBy)
	})
	return
}

// 注：仅支持1000条以内的简单同步导出 by codegen
func (menu *SysMenuService) ExportSysMenu(menuQuery *userDto.SysMenuQuery) []byte {
	nPage := page.NewPage()
	nPage.PageSize = page.MAX_PAGE_SIZE
	dataList := sysMenuDao.List(menuQuery, nPage)
	row_title := []interface{}{"菜单ID", "菜单名称", "父菜单ID", "显示顺序", "路由地址", "组件路径", "路由参数", "是否为外链", "是否缓存", "菜单类型", "菜单状态", "菜单状态", "权限标识", "菜单图标", "创建者", "创建时间", "更新者", "更新时间", "备注"}
	dataSlice := make([][]interface{}, 0)
	dataSlice = append(dataSlice, row_title)
	for _, v := range dataList {
		dataSlice = append(dataSlice, []interface{}{v.MenuId, v.MenuName, v.ParentId, v.OrderNum, v.Path, v.Component, v.Query, v.IsFrame, v.IsCache, v.MenuType, v.Visible, v.State, v.Perms, v.Icon, v.CreateBy, v.CreateTime, v.UpdateBy, v.UpdateTime, v.Remark})
	}
	return excel.ExportExcel(dataSlice)
}

// 根据菜单树组装路由结构
func buildMenuRouters(menus []*userDo.SysMenu) (list []*userDto.RouterVO) {
	for _, menu := range menus {
		router := new(userDto.RouterVO)
		router.Hidden = menu.Visible == 2
		router.Name = getRouteName(menu)
		router.Path = getRoutePath(menu)
		router.Component = getComponent(menu)
		router.Query = menu.Query

		meta := new(userDto.MetaVO)
		meta.Title = menu.MenuName
		meta.Icon = menu.Icon
		meta.NoCache = menu.IsCache == 2
		if isHttpLink(menu.Path) {
			meta.Link = menu.Path
		}
		router.Meta = meta

		hasChild := len(menu.Children) > 0
		if hasChild && menu.MenuType == "M" {
			router.AlwaysShow = true
			router.Redirect = "noRedirect"
			router.Children = buildMenuRouters(menu.Children)
		} else if isMenuFrame(menu) {
			router.Meta = nil
			childrenList := make([]*userDto.RouterVO, 0)
			children := new(userDto.RouterVO)
			children.Path = menu.Path
			children.Component = menu.Component
			children.Name = stringUtils.Capitalize(menu.Path)
			children.Meta = meta
			children.Query = menu.Query
			childrenList = append(childrenList, children)
			router.Children = childrenList
		} else if menu.ParentId == 0 && isInnerLink(menu) {
			meta1 := new(userDto.MetaVO)
			meta1.Title = menu.MenuName
			meta1.Icon = menu.Icon
			router.Meta = meta1
			router.Path = "/"
			childrenList := make([]*userDto.RouterVO, 0)
			children := new(userDto.RouterVO)
			routerPath := innerLinkReplaceEach(menu.Path)
			children.Path = routerPath
			children.Component = "InnerLink"
			children.Name = stringUtils.Capitalize(routerPath)
			meta2 := new(userDto.MetaVO)
			meta2.Title = menu.MenuName
			meta2.Icon = menu.Icon
			meta2.Link = menu.Path
			children.Meta = meta2
			childrenList = append(childrenList, children)
			router.Children = childrenList
		}
		list = append(list, router)
	}
	return
}

func buildMenuTree(menus []*userDo.SysMenu) (list []*userDto.TreeSelect) {
	mList := make([]*userDo.SysMenu, 0)
	mIdMap := make(map[int64]int64, 0)
	for _, m := range menus {
		mIdMap[m.MenuId] = m.MenuId
	}
	for _, m := range menus {
		// 如果是顶级节点，则遍历该父节点的所有子节点
		if mIdMap[m.ParentId] < 1 {
			recursionFn(menus, m)
			mList = append(mList, m)
		}
	}
	for _, m := range mList {
		list = append(list, convertMenuToTreeSelect(m))
	}
	return
}

func convertMenuToTreeSelect(menu *userDo.SysMenu) *userDto.TreeSelect {
	tree := &userDto.TreeSelect{}
	tree.Id = menu.MenuId
	tree.Label = menu.MenuName
	if len(menu.Children) > 0 {
		treeChilds := make([]*userDto.TreeSelect, 0)
		for _, c := range menu.Children {
			treeChilds = append(treeChilds, convertMenuToTreeSelect(c))
		}
		tree.Children = treeChilds
	}
	return tree
}

// 根据父菜单ID，获取所有子菜单
func getChildPerms(list []*userDo.SysMenu, parentId int64) (menus []*userDo.SysMenu) {
	for _, m := range list {
		if m.ParentId == parentId {
			recursionFn(list, m)
			menus = append(menus, m)
		}
	}
	return
}

// 递归列表
func recursionFn(list []*userDo.SysMenu, menu *userDo.SysMenu) {
	childList := getChildList(list, menu)
	menu.Children = childList
	for _, c := range childList {
		if hasChild(list, c) {
			recursionFn(list, c)
		}
	}
}

// 获得子菜单列表
func getChildList(list []*userDo.SysMenu, menu *userDo.SysMenu) (tlist []*userDo.SysMenu) {
	for _, t := range list {
		if t.ParentId == menu.MenuId {
			tlist = append(tlist, t)
		}
	}
	return
}

// 是否有子菜单
func hasChild(list []*userDo.SysMenu, menu *userDo.SysMenu) bool {
	clist := getChildList(list, menu)
	return len(clist) > 0
}

// 获取路由名称
func getRouteName(menu *userDo.SysMenu) string {
	routeName := stringUtils.Capitalize(menu.Path)
	if isMenuFrame(menu) {
		routeName = ""
	}
	return routeName
}

// 是否为菜单内部跳转
func isMenuFrame(menu *userDo.SysMenu) bool {
	return menu.ParentId == 0 && menu.MenuType == "C" && menu.IsFrame == 2
}

// 获取路由地址
func getRoutePath(menu *userDo.SysMenu) string {
	routePath := menu.Path
	// 内链打开外网方式
	if menu.ParentId != 0 && isInnerLink(menu) {
		routePath = innerLinkReplaceEach(routePath)
	}
	// 非外链并且是一级目录（类型为目录）
	if menu.ParentId == 0 && menu.MenuType == "M" && menu.IsFrame == 2 {
		routePath = "/" + menu.Path
		// 非外链并且是一级目录（类型为菜单）
	} else if isMenuFrame(menu) {
		routePath = "/"
	}
	return routePath
}

// 是否为内链组件
func isInnerLink(menu *userDo.SysMenu) bool {
	return menu.IsFrame == 2 && isHttpLink(menu.Path)
}

// 是否为parent_view组件
func isParentView(menu *userDo.SysMenu) bool {
	return menu.ParentId != 0 && menu.MenuType == "M"
}

// 是否为http(s)://开头
func isHttpLink(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

// 内链域名特殊字符替换
func innerLinkReplaceEach(path string) string {
	result := strings.Replace(path, "http://", "", -1)
	result = strings.Replace(result, "https://", "", -1)
	result = strings.Replace(result, "www.", "", -1)
	result = strings.Replace(result, ".", "/", -1)
	return result
}

// 获取组件信息
func getComponent(menu *userDo.SysMenu) string {
	component := "Layout"
	if menu.Component != "" && !isMenuFrame(menu) {
		component = menu.Component
	} else if menu.Component == "" && menu.ParentId != 0 && isInnerLink(menu) {
		component = "InnerLink"
	} else if menu.Component == "" && isParentView(menu) {
		component = "ParentView"
	}
	return component
}

// Package hint Package cst
/**
 * @note: 提示
 *
 * @author: 大猫
 * @date: 2023/5/9
 */
package hint

import "github.com/sealsee/web-base/public/errs"

var SYS_USER_EXISTS = errs.ERROR{"1000", "账号已存在"}
var SYS_OLD_PASSWORD_ERR = errs.ERROR{"1001", "修改密码失败，旧密码错误"}
var SYS_SAME_PASSWORD_ERR = errs.ERROR{"1002", "新密码不能与旧密码相同"}
var SYS_ROLE_EXISTS = errs.ERROR{"1003", "角色名称或权限字符已存在"}
var SYS_MENU_EXISTS = errs.ERROR{"1004", "菜单名称已存在"}
var SYS_MENU_PARENTID_ERR = errs.ERROR{"1005", "上级菜单不能选择自己"}
var SYS_MENU_HAS_CHILD = errs.ERROR{"1006", "存在子菜单，不允许删除"}
var SYS_MENU_HAS_ROLE = errs.ERROR{"1007", "菜单已分配，不允许删除"}
var SYS_CONFIG_EXISTS = errs.ERROR{"1008", "参数名称或参数键名已存在"}
var SYS_FILE_READ_FAIL = errs.ERROR{"1009", "文件读取失败！"}
var SYS_PASSWORD_FORMAT_ERR = errs.ERROR{"1010", "密码输入错误，应由6~18位字符组成，并且同时包含大小写字母、数字！"}
var SYS_LENGTH_ERR = errs.ERROR{"1011", "内容过长，请检查各个字段的输入！"}
var SYS_USERNAME_FORMAT_ERR = errs.ERROR{"1012", "账号输入错误，应由6~18位字符组成,可以包含数字和字母！"}
var SYS_NICKNAME_FORMAT_ERR = errs.ERROR{"1013", "姓名输入错误，应由2~15位字符组成！"}
var SYS_NO_PERMISSION = errs.ERROR{"1020", "无数据权限！"}
var SYS_USER_NO_ROLEKEY = errs.ERROR{"1021", "请选择账号权限！"}
var SYS_USER_NO_DEPARTMENT = errs.ERROR{"1022", "请选择所属学院！"}

var BIZ_TCH_TASK_NAME_LEN_ERR = errs.ERROR{"2000", "任务名称输入错误，应由2~30位字符组成！"}
var BIZ_TCH_TASK_TIME_EMPTY_ERR = errs.ERROR{"2001", "征集日期不能为空"}
var BIZ_TCH_TASK_TIME_DIFF_ERR = errs.ERROR{"2002", "征集日期时间跨度不能超过7天"}
var BIZ_TCH_TASK_DETAIL_EMPTY_ERR = errs.ERROR{"2003", "征集任务详情不能为空"}
var BIZ_TCH_TASK_DETAIL_INFO_ERR = errs.ERROR{"2004", "征集任务设置错误, 日期:[%v], 学院:[%v], 人数:[%v]"}
var BIZ_TCH_TASK_NAME_EXIST_ERR = errs.ERROR{"2005", "任务名称已存在！"}
var BIZ_TCH_TASK_STATE_PUBLISHED_ERR = errs.ERROR{"2006", "征集任务已发布，不能修改！"}
var BIZ_TCH_TASK_DETAIL_DEPT_DUPLICATE_ERR = errs.ERROR{"2007", "征集任务明细数据重复! 日期:[%v], 学院:[%v]"}
var BIZ_TCH_TASK_NOT_EXIST_ERR = errs.ERROR{"2008", "征集任务不存在"}
var BIZ_NEED_REPLACE_TCH_ERR = errs.ERROR{"2009", "请从未征集的教职工列表中选择替换"}
var BIZ_TCH_TASK_DETAIL_NOT_EXIST_ERR = errs.ERROR{"2010", "征集详情不存在"}
var BIZ_TCH_NOT_EXIST_ERR = errs.ERROR{"2011", "教职工信息不存在"}
var BIZ_TCH_HAS_REPORT_THIS_TIME_ERR = errs.ERROR{"2012", "教职工：%v(%v)在%v已经被上报了任务，请检查数据"}
var BIZ_TCH_CAN_NOT_EDIT_ERR = errs.ERROR{"2013", "招办新增和导入的教职工不可被学院编辑"}
var BIZ_TCH_JOB_NO_EXISTS_ERR = errs.ERROR{"2014", "工号已存在"}
var BIZ_TCH_ID_CARD_EXISTS_ERR = errs.ERROR{"2015", "证件号已存在"}

var BIZ_ARRANGE_ROOM_ADD_ERR = errs.ERROR{"2016", "当前任务下的专业方向已完成编排，请选择编辑"}
var BIZ_ARRANGE_ROOM_EDIT_ERR = errs.ERROR{"2017", "当前任务下的专业方向未完成编排，请先新增"}
var BIZ_NO_ARRANGE_ROOM_ERR = errs.ERROR{"2018", "您所选择的专业编排计划不存在"}
var BIZ_ARRANGE_ROOM_NO_FORMAT_ERR = errs.ERROR{"2019", "第%v行的考场编号格式不正确"}
var BIZ_HAS_TCH_ARRANGE_ERR = errs.ERROR{"2020", "您已经分配了教职工，无法取消编排"}
var BIZ_ARRANGE_ROOM_NO_REPEAT_ERR = errs.ERROR{"2021", "当前考场任务下已经存在编号为%v的编排考场"}
var BIZ_TCH_HAS_ARRANGED_ERR = errs.ERROR{"2022", "教职工：%v已经在当前任务下被分配，无法重复分配"}
var BIZ_ROOM_TASK_NOT_EXIST_ERR = errs.ERROR{"2023", "上报信息不存在"}
var BIZ_TASK_DAY_NOT_MATCH_ERR = errs.ERROR{"2024", "教职工%v上报的日期和考场征集任务时间不匹配"}
var BIZ_TCH_HAS_ARRANGED_CANCEL_ERR = errs.ERROR{"2025", "教职工已被编排，不允许取消上报"}
var BIZ_TCH_REPORT_NOT_EXIST_ERR = errs.ERROR{"2026", "教职工不存在：%v"}
var BIZ_TCH_AUDIT_REASON_EMPTY_ERR = errs.ERROR{"2027", "审核不通过时必须填写不通过原因"}
var BIZ_TCH_AUDIT_STATUS_INVALID_ERR = errs.ERROR{"2028", "当前教职工数据无需审核"}
var BIZ_TCH_UNSUBMITTED_CANNOT_EDIT_ERR = errs.ERROR{"2029", "学院未提交的教职工数据，招办不可编辑"}
var BIZ_TCH_AUDIT_NOT_PENDING_ERR = errs.ERROR{"2030", "仅待审核状态的教职工支持审核操作"}
var BIZ_TCH_AUDIT_REJECTED_REASON_EMPTY_ERR = errs.ERROR{"2031", "审核不通过原因不能为空"}

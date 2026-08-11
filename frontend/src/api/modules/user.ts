import { ReqPage, ResPage, User, Common } from "@/api/interface/index";
import { PROXY_TAG } from "@/api/config/servicePort";
import http from "@/api";

/**
 * @name 用户管理模块
 */

// ========================== 用户管理 =========================
// 获取用户列表
export const getUserList = (params: ReqPage) => {
  return http.post<ResPage<User.ResUserList>>(PROXY_TAG + `/api/auth/sys_user/list`, params);
};

// 切换用户状态
export const changeUserStatus = (params: { userId: string; state: string }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/sys_user/change_state`, params);
};

// 新增用户
export const addUser = (params: User.ReqAddUser) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/sys_user/add`, params);
};

// 编辑用户
export const editUser = (params: User.ReqEditUser) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/sys_user/edit`, params);
};

// 设置用户密码
export const setUserPassword = (params: { userId: string; password: string }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/sys_user/reset_pwd`, params);
};

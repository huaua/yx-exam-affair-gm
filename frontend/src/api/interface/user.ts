/**
 * @name 用户管理模块
 */

export namespace User {
  // ========================== 用户管理 =========================
  export interface ResUserList {
    userId: string; // 用户id
    userName: string; // 用户名
    nickName: string; // 用户昵称
    deptId: string; // 学院id
    deptName: string; // 学院名称
    state: string; // 状态 '1'-启用 '0'-禁用
    password?: string; // 密码
    roleKey?: string; // 角色代码： department-学院 admissions-招办
  }
  export interface ReqAddUser {
    userName: string; // 用户名
    nickName: string; // 用户昵称
    password: string; // 密码
    roleKey: string; // 角色代码： department-学院 admissions-招办
    deptId: string; // 所属学院id （可选，roleKey=department时必填）
  }
  export interface ReqEditUser {
    userName: string; // 用户名
    nickName: string; // 用户昵称
  }
}

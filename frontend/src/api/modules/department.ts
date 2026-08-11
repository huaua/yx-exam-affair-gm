import { ReqPage, ResPage, Department, Common } from "@/api/interface/index";
import { PROXY_TAG } from "@/api/config/servicePort";
import http from "@/api";

/**
 * @name 学院管理模块
 */

// ========================== 学院管理 =========================
// 获取学院列表
export const getDepartmentList = (params: ReqPage) => {
  return http.post<ResPage<Department.ResDepartmentList>>(PROXY_TAG + `/api/auth/ea_dept_info/list`, params, { cancel: false });
};

// 新增学院
export const addDepartment = (params: Department.ReqAddDepartment[]) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_dept_info/add`, params);
};

// 编辑学院
export const editDepartment = (params: Department.ReqEditDepartment) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_dept_info/edit`, params);
};

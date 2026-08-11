/**
 * @name 学院管理模块
 */
export namespace Department {
  // ========================== 学院管理 =========================
  export interface ReqDepartmentList {
    curPage?: number;
    pageSize?: number;
  }
  export interface ResDepartmentList {
    deptId: string; // 主键ID 所在学院ID
    deptCode: string; // 学院代码
    deptName: string; // 单位名称
    campusName: string; // 所在校区
  }
  export interface ReqAddDepartment {
    deptCode: string; // 学院代码
    deptName: string; // 单位名称
    campusName: string; // 所在校区
  }
  export interface ReqEditDepartment {
    deptId: string; // 主键ID 所在学院ID
    deptCode: string; // 学院代码
    deptName: string; // 单位名称
    campusName: string; // 所在校区
  }
}

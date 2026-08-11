import { ReqPage, ResPage, Arrange, Common } from "@/api/interface/index";
import { PROXY_TAG } from "@/api/config/servicePort";
import http from "@/api";

/**
 * @name 考场编排模块
 */

// ========================== 考场排考编排 =========================
// 获取考场编排列表
export const getClassroomArrangeList = (params: ReqPage<Arrange.ReqClassroomArrangeList>) => {
  return http.post<ResPage<Arrange.ResClassroomArrangeList>>(PROXY_TAG + `/api/auth/ea_arrange_prof_info/list`, params);
};

// 导入编排学院专业方向
export const importArrangeDeptProfDirection = (params: { roomTaskId: string; files: { [fileName: string]: string } }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_arrange_prof_info/import`, params);
};

// 获取可编排教室列表
export const getCanArrangeClassroomList = (params: { infoId: string }) => {
  return http.post<{ list: Arrange.ResCanArrangeClassroomList[] }>(
    PROXY_TAG + `/api/auth/ea_arrange_prof_info/list_available_room`,
    params
  );
};

// 获取编排详情
export const getArrangeDetail = (params: { infoId: string }) => {
  return http.post<{ obj: Arrange.ResArrangeDetail }>(PROXY_TAG + `/api/auth/ea_arrange_prof_info/get`, params);
};

// 获取编排详情列表
export const getArrangeDetailList = (params: { infoId: string }) => {
  return http.post<{ list: Arrange.ResArrangeDetailList[] }>(PROXY_TAG + `/api/auth/ea_arrange_room/list`, params);
};

// 保存编排详情列表
export const saveArrangeDetailList = (params: Arrange.ReqSaveArrangeDetailList) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_arrange_room/edit`, params);
};

// 取消编排
export const cacelArrange = (params: { infoId: string }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_arrange_prof_info/cancel_arrange`, params);
};

// ========================== 考场监考分配 =========================
// 获取监考分配列表
export const getTeacherAssignList = (params: ReqPage<Arrange.ReqTeacherAssignList>) => {
  return http.post<ResPage<Arrange.ResTeacherAssignList>>(PROXY_TAG + `/api/auth/ea_arrange_tch/list`, params);
};

// 保存监考老师分配
export const saveTeacherAssign = (params: Arrange.ReqSaveTeacherAssignItem[]) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_arrange_tch/edit`, params);
};

// 取消分配
export const cacelAssign = (params: { roomReportId: string }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_arrange_tch/del`, params);
};

// 取消单个监考老师分配
export const cancelSingleAssign = (params: { roomReportId: string; tchReportId: string }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_arrange_tch/del_single`, params);
};

// ========================== 编排查看 =========================
// 获取监考分配列表
export const getArrangeViewList = (params: ReqPage<Arrange.ReqArrangeViewList>) => {
  return http.post<ResPage<Arrange.ResArrangeViewList>>(PROXY_TAG + `/api/auth/ea_arrange_tch/list_by_prof`, params);
};

// 导出编排列表
export const exportArrangeInfo = (params: Arrange.ReqArrangeViewList) => {
  return http.download(PROXY_TAG + `/api/auth/ea_arrange_tch/export`, params);
};

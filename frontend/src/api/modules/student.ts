import { ReqPage, ResPage, Student } from "@/api/interface/index";
import { PROXY_TAG } from "@/api/config/servicePort";
import http from "@/api";

/**
 * @name 考生管理模块
 */

// ========================== 考生管理 =========================
// 获取考生列表
export const getStudentList = (params: ReqPage<Student.ReqStudentList>) => {
  return http.post<ResPage<Student.ResStudentList>>(PROXY_TAG + `/api/auth/ea_sync_stu_prof/list`, params);
};

// 检查考生信息同步状态
export const checkStudentInfoSyncState = (params: {}) => {
  return http.post<{ obj: Student.ResSyncStudentInfo }>(PROXY_TAG + `/api/auth/ea_sync_stu_prof/get_sync_state`, params, {
    loading: false
  });
};

// 获取考试列表
export const getExamList = (params: { reSync?: boolean }) => {
  return http.post<{ list: Student.ResExamList[] }>(PROXY_TAG + `/api/auth/ea_sync_stu_prof/get_exam_list`, params);
};

// 同步考生信息
export const syncStudentInfo = (params: Student.ReqSyncStudentInfo) => {
  return http.post<{ obj: Student.ResSyncStudentInfo }>(PROXY_TAG + `/api/auth/ea_sync_stu_prof/do_sync`, params);
};

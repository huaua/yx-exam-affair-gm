import { ReqPage, ResPage, Teacher, Common } from "@/api/interface/index";
import { PROXY_TAG } from "@/api/config/servicePort";
import http from "@/api";

/**
 * @name 任务征集模块
 */

// ========================== 教职工管理 =========================
// 获取教职工列表
export const getTeacherList = (params: ReqPage<Teacher.ReqTeacherList>) => {
  return http.post<ResPage<Teacher.ResTeacherList>>(PROXY_TAG + `/api/auth/ea_tch_info/list`, params);
};

// 新增教职工
export const addTeacher = (params: Teacher.ReqAddTeacher) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_info/add`, params);
};

// 编辑教职工
export const editTeacher = (params: Teacher.ReqEditTeacher) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_info/edit`, params);
};

// 删除教职工
export const deleteTeacher = (params: { tchId: string }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_info/del`, params);
};

// 导入教职工
export const importTeacher = (params: any) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_info/import`, params);
};

// 导出教职工
export const exportTeacher = (params: Partial<Teacher.ReqTeacherList>) => {
  return http.download(PROXY_TAG + `/api/auth/ea_tch_info/export`, params);
};

// ========================== 监考老师征集管理 =========================
// 获取监考老师征集管理列表
export const getTeacherRecruitManageList = (params: ReqPage<Teacher.ReqTeacherRecruitManageList>) => {
  return http.post<ResPage<Teacher.ResTeacherRecruitManageList>>(PROXY_TAG + `/api/auth/ea_tch_task_info/list`, params);
};

// 切换发布状态
export const changeTeacherRecruitManageStatus = (params: { tchTaskId: string; state: string }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_task_info/change_state`, params);
};

// 校验监考老师征集任务名称
export const checkTeacherRecruitTaskName = (params: Teacher.ReqCheckTeacherRecruitTaskName) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_task_info/check_info`, params);
};

// 新增监考老师征集任务
export const addTeacherRecruitTask = (params: Teacher.ReqAddTeacherRecruitTask) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_task_info/add`, params);
};

// 编辑监考老师征集任务
export const editTeacherRecruitTask = (params: Teacher.ReqEditTeacherRecruitTask) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_task_info/edit`, params);
};

// 编辑监考老师征集任务详情
export const editTeacherRecruitTaskDetail = (params: { tchTaskId: string }) => {
  return http.post<{ list: Teacher.ResEditTeacherRecruitTaskDetail[] }>(
    PROXY_TAG + `/api/auth/ea_tch_task_detail/list_by_day`,
    params
  );
};

// ========================== 监考老师征集查看 =========================
// 获取监考老师征集查看列表
export const getTeacherRecruitViewList = (params: ReqPage<Teacher.ReqTeacherRecruitViewList>) => {
  return http.post<ResPage<Teacher.ResTeacherRecruitViewList>>(PROXY_TAG + `/api/auth/ea_tch_task_report/list_by_task`, params);
};

// 导出监考老师征集查看列表
export const exportTeacherRecruitViewList = (params: { tchTaskId: string }) => {
  return http.download(PROXY_TAG + `/api/auth/ea_tch_task_report/export`, params);
};

// 审核征集状态（招办）
export const auditTeacherReport = (params: { reportId: string; action: "pass" | "reject" }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_task_report/audit`, params);
};

// 重置征集状态（招办）
export const resetTeacherReport = (params: { reportId: string }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_task_report/reset`, params);
};

// ========================== 监考老师上报 =========================
// 获取监考老师上报任务列表
export const getTeacherSumbitList = (params: {}) => {
  return http.post<{ list: Teacher.ResTeacherRecruitManageList[] }>(
    PROXY_TAG + `/api/auth/ea_tch_task_info/list_all_by_dept`,
    params
  );
};

// 获取可上报的监考老师列表
export const getCanSubmitTeacherList = (params: Teacher.ReqCanSubmitTeacherList) => {
  return http.post<{ list: Teacher.ResCanSubmitTeacherList[] }>(PROXY_TAG + `/api/auth/ea_tch_task_report/list_tch`, params);
};

// 获取当前任务的日期列表
export const getCurTaskDateList = (params: { tchTaskId: string; onlyQueryNeedReportTimeList: boolean }) => {
  return http.post<{ list: Teacher.ResCurTaskDateList[] }>(PROXY_TAG + `/api/auth/ea_tch_task_detail/list`, params);
};

// 缓存预上报的监考老师
export const cacheSubmitTeacher = (params: Teacher.ReqSubmitTeacher[]) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_task_report/temp_save`, params);
};

// 清除缓存预上报的监考老师
export const nocacheSubmitTeacher = (params: Teacher.ReqSubmitTeacher[]) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_task_report/temp_clear`, params);
};

// 获取缓存预上报的监考老师
export const getCacheSubmitTeacher = (params: { tchTaskId: string }) => {
  return http.post<{ list: Teacher.ResGetCacheSubmitTeacherItem[] }>(
    PROXY_TAG + `/api/auth/ea_tch_task_report/temp_list`,
    params
  );
};

// 上报监考老师
export const submitTeacher = (params: Teacher.ReqSubmitTeacher[]) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_task_report/add`, params);
};

// 追加监考老师
export const appendTeacher = (params: Teacher.ReqAppendTeacher) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_task_report/edit`, params);
};

// ========================== 监考老师上报查看 =========================
// 获取监考老师上报查看列表
export const getTeacherSubmitViewList = (params: ReqPage<Teacher.ReqTeacherRecruitManageList>) => {
  return http.post<ResPage<Teacher.ResTeacherSubmitViewList>>(PROXY_TAG + `/api/auth/ea_tch_task_report/list`, params);
};

// 校验是否可以取消上报
export const checkCancelSubmit = (params: { reportId: string }) => {
  return http.post<{ obj: boolean }>(PROXY_TAG + `/api/auth/ea_tch_task_report/check_cancel`, params);
};

// 取消上报
export const cancelSubmit = (params: { reportId: string }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_task_report/del`, params);
};

// 替换上报老师
export const replaceSubmitTeacher = (params: { reportId: string; taskDetailId: string; tchId: string }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_task_report/replace_report`, params);
};

// ========================== 教职工审核 =========================
// 审核教职工（招办）
export const auditTeacher = (params: Teacher.ReqAuditTeacher) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_info/audit`, params);
};

// 切换启用/禁用状态（招办）
export const changeTeacherEnabled = (params: Teacher.ReqChangeTchEnabled) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_info/change_enabled`, params);
};

// 批量切换启用/禁用状态（招办）
export const changeTeacherEnabledBatch = (params: Teacher.ReqChangeTchEnabledBatch) => {
  return http.post<{ affected: number }>(PROXY_TAG + `/api/auth/ea_tch_info/change_enabled_batch`, params);
};

// 学院提交教职工至待审核（仅未提交状态可提交）
// 获取下一条待审核教职工，params 传入列表页当前查询条件与当前记录 tchId，保证按筛选范围顺序流转
export const getNextPendingTeacher = (params: Record<string, any> = {}) => {
  return http.post<{ obj: Teacher.ResTeacherList | null }>(PROXY_TAG + `/api/auth/ea_tch_info/next_pending`, params);
};

export const submitTchForAudit = (params: { tchId: string }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_tch_info/submit`, params);
};

// 学院批量提交教职工至待审核（tchIds 为空时提交本学院全部未提交数据）
export const submitTchForAuditBatch = (params: { tchIds: number[] }) => {
  return http.post<{ affected: number }>(PROXY_TAG + `/api/auth/ea_tch_info/submit_batch`, params);
};

// 查询监考记录（按日期倒序）
export const getInvigilationRecords = (params: { tchId: string }) => {
  return http.post<{ list: Teacher.ResInvigilationRecord[] }>(PROXY_TAG + `/api/auth/ea_tch_info/arrange_records`, params);
};

// 查询教职工审核历史
export const getTeacherAuditHistory = (params: { tchId: string }) => {
  return http.post<{ list: Teacher.ResAuditHistory[] }>(PROXY_TAG + `/api/auth/ea_tch_info/audit_history`, params);
};

// 查询最近一次审核通过的快照
export const getLatestAuditHistory = (params: { tchId: string }) => {
  return http.post<{ obj: Teacher.ResAuditHistory | null }>(PROXY_TAG + `/api/auth/ea_tch_info/latest_audit_history`, params);
};

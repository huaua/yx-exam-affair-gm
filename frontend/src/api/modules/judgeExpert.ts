import { ReqPage, ResPage, JudgeExpert, Common } from "@/api/interface/index";
import { PROXY_TAG } from "@/api/config/servicePort";
import http from "@/api";

/**
 * @name 评委专家模块
 */

// ========================== 专家库管理 =========================
// 获取专家库列表
export const getExpertDatabaseList = (params: ReqPage<JudgeExpert.ReqExpertDatabaseList>) => {
  return http.post<ResPage<JudgeExpert.ResExpertDatabaseList>>(`${PROXY_TAG}/api/auth/ea_judge_experts/list`, params);
};

// 获取专家详情（含评价内容）
export const getExpertDetail = (params: { expertId: string }) => {
  return http.post<{ obj: JudgeExpert.ResExpertDatabaseList | null }>(`${PROXY_TAG}/api/auth/ea_judge_experts/get`, params);
};

// 新增专家
export const addExpert = (params: JudgeExpert.ReqAddExpert) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_experts/add`, params);
};

// 编辑专家
export const editExpert = (params: JudgeExpert.ReqEditExpert) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_experts/edit`, params);
};

export const submitExpertForAudit = (params: { expertId: string }) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_experts/submit`, params);
};

// 学院批量提交专家审核（expertIds 为空时提交本学院全部未提交数据）
export const submitExpertForAuditBatch = (params: { expertIds: number[] }) => {
  return http.post<{ affected: number }>(`${PROXY_TAG}/api/auth/ea_judge_experts/submit_batch`, params);
};

export const auditExpert = (params: JudgeExpert.ReqAuditExpert) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_experts/audit`, params);
};

// 获取下一条待审核专家，params 传入列表页当前查询条件与当前记录 expertId，保证按筛选范围顺序流转
export const getNextPendingExpert = (params: Record<string, any> = {}) => {
  return http.post<{ obj: JudgeExpert.ResExpertDatabaseList | null }>(
    `${PROXY_TAG}/api/auth/ea_judge_experts/next_pending`,
    params
  );
};

export const getLatestExpertAuditHistory = (params: { expertId: string }) => {
  return http.post<{ obj: JudgeExpert.ResExpertAuditHistory | null }>(
    `${PROXY_TAG}/api/auth/ea_judge_experts/latest_audit_history`,
    params
  );
};

export const evaluateExpert = (params: { expertId: string; evaluation: string }) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_experts/evaluate`, params);
};

export const changeExpertStatus = (params: { expertId: string; allowDraw: string }) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_experts/status`, params);
};

// 删除专家
export const deleteExpert = (params: { expertId: string }) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_experts/del`, params);
};

// 导入专家
export const importExpert = (params: any) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_experts/import`, params);
};

// 导出专家
export const exportExpert = (params: JudgeExpert.ReqExpertDatabaseList) => {
  return http.download(`${PROXY_TAG}/api/auth/ea_judge_experts/export`, params);
};

export const getBannedExpertList = (params: ReqPage<JudgeExpert.ReqExpertDatabaseList>) =>
  http.post<ResPage<JudgeExpert.ResExpertDatabaseList>>(`${PROXY_TAG}/api/auth/ea_judge_experts_ban/list`, params);
export const getBanExpertCandidates = (params: ReqPage<JudgeExpert.ReqExpertDatabaseList>) =>
  http.post<ResPage<JudgeExpert.ResExpertDatabaseList>>(`${PROXY_TAG}/api/auth/ea_judge_experts_ban/candidates`, params);
export const addBannedExpert = (params: { expertId: string }) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_experts_ban/add`, params);
export const deleteBannedExperts = (params: { expertIds: number[] }) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_experts_ban/del`, params);
export const deleteAllBannedExperts = () =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_experts_ban/del_all`, {});
export const importBannedExperts = (params: any) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_experts_ban/import`, params);
export const exportBannedExperts = (params: JudgeExpert.ReqExpertDatabaseList) =>
  http.download(`${PROXY_TAG}/api/auth/ea_judge_experts_ban/export`, params);

// ========================== 评委征集任务管理 =========================
// 获取评委征集管理列表
export const getJudgeRecruitManageList = (params: ReqPage<JudgeExpert.ReqJudgeRecruitManageList>) => {
  return http.post<ResPage<JudgeExpert.ResJudgeRecruitManageList>>(`${PROXY_TAG}/api/auth/ea_judge_task_info/list`, params);
};

// 新增评委征集任务
export const addJudgeRecruitTask = (params: JudgeExpert.ReqAddJudgeRecruitTask) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_task_info/add`, params);
};

// 编辑评委征集任务
export const editJudgeRecruitTask = (params: JudgeExpert.ReqEditJudgeRecruitTask) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_task_info/edit`, params);
};

export const appendJudgeTaskRounds = (params: { judgeTaskId: string; scoreRounds: string }) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_task_info/append_rounds`, params);

export const publishJudgeRecruitTask = (params: { judgeTaskId: string }) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_task_info/publish`, params);

export const getJudgeTaskRequirements = (params: { judgeTaskId: string }) =>
  http.post<{ list: JudgeExpert.JudgeTaskReportRequirement[] }>(
    `${PROXY_TAG}/api/auth/ea_judge_task_info/requirement/list`,
    params
  );
export const importJudgeTaskRequirements = (params: { judgeTaskId: string; files: Record<string, string> }) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_task_info/requirement/import`, params);
export const editJudgeTaskRequirement = (params: { requirementId: string; expertNum: string }, options = {}) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_task_info/requirement/edit`, params, options);
export const batchEditJudgeTaskRequirements = (params: { items: { requirementId: string; expertNum: string }[] }) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_task_info/requirement/batch_edit`, params);
export const deleteJudgeTaskRequirement = (params: { requirementId: string }) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_task_info/requirement/del`, params);

// 抽取评委
export const extractJudge = (params: { judgeTaskId: string }) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_task_info/extract`, params);
};

// 抽取评委按擅长科目
export const extractJudgeBySubject = (params: { judgeTaskId: string }) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_task_info/extract_by_subject`, params);
};

// 取消评委征集任务
export const cancelJudgeRecruitTask = (params: { judgeTaskId: string }) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_task_info/cancel_extract`, params);
};

// 删除评委征集任务
export const deleteJudgeRecruitTask = (params: { judgeTaskId: string }) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_task_info/del`, params);
};

// 查看评委征集任务
export const viewJudgeRecruitTask = (params: { judgeTaskId: string }) => {
  return http.post<any>(`${PROXY_TAG}/api/auth/ea_judge_task_info/get`, params);
};

// ========================== 评委抽取管理 ==========================
export const getDrawManageTasks = () => http.post<any>(`${PROXY_TAG}/api/auth/ea_judge_draw_manage/tasks`, {});
export const getDrawManageRows = (params: { judgeTaskId: string }) =>
  http.post<any>(`${PROXY_TAG}/api/auth/ea_judge_draw_manage/rows`, params);
export const getDrawManageExperts = (params: any) => http.post<any>(`${PROXY_TAG}/api/auth/ea_judge_draw_manage/experts`, params);
export const saveDrawManageSelection = (params: any) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_draw_manage/save_selection`, params);
export const removeDrawManageSelection = (params: any) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_draw_manage/remove_selection`, params);
export const submitDrawManageSelection = (params: any) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_draw_manage/submit`, params);
export const updateDrawManageRule = (params: any) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_draw_manage/update_rule`, params);
export const exportDrawManageExperts = (params: any) =>
  http.download(`${PROXY_TAG}/api/auth/ea_judge_draw_manage/export`, params);

// ========================== 评委抽取管理 ==========================
// ========================== 评委应答 =========================
// 获取所有评委征集任务列表
export const getAllJudgeRecruitTaskList = (params: {}) => {
  return http.post<ResPage<JudgeExpert.ResJudgeRecruitManageList>>(`${PROXY_TAG}/api/auth/ea_judge_task_info/list_all`, params);
};

// 获取评委应答列表
export const getJudgeResponseList = (params: ReqPage<JudgeExpert.ReqJudgeResponseList>) => {
  return http.post<{
    list: { attendNum: string; unattendNum: string; notSureNum: string; list: JudgeExpert.ResJudgeResponseList[] };
    page: any;
  }>(`${PROXY_TAG}/api/auth/ea_judge_draw_detail/list_by_role`, params);
};

// 应答评委应答
export const responseJudgeResponse = (params: { id: string; confirmState: string }) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_draw_detail/edit`, params);
};

// ========================== 学院评委专家上报 ==========================
export const getCollegeReportTasks = () =>
  http.post<{ list: JudgeExpert.CollegeReportTask[] }>(`${PROXY_TAG}/api/auth/ea_judge_expert_report/tasks`, {});
export const getCollegeReportDirections = (params: { judgeTaskId: string }) =>
  http.post<{ list: JudgeExpert.CollegeReportDirection[] }>(`${PROXY_TAG}/api/auth/ea_judge_expert_report/directions`, params);
export const getCollegeReportExperts = (params: any) =>
  http.post<{ list: JudgeExpert.CollegeReportExpert[] }>(`${PROXY_TAG}/api/auth/ea_judge_expert_report/experts`, params);
export const saveCollegeReportDraft = (params: { requirementId: string; expertId: string }) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_expert_report/save`, params);
export const removeCollegeReportDraft = (params: { reportId: string }) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_expert_report/remove`, params);
export const submitCollegeReport = (params: { requirementId: string }) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_expert_report/submit`, params);
export const getCollegeReportedExperts = (params: { requirementId: string }) =>
  http.post<{ list: JudgeExpert.CollegeReportExpert[] }>(`${PROXY_TAG}/api/auth/ea_judge_expert_report/reported`, params);
export const replaceCollegeReportedExpert = (params: { reportId: string; newExpertId: string }) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_expert_report/replace`, params);
export const exportCollegeReportedExperts = (params: { judgeTaskId: string }) =>
  http.download(`${PROXY_TAG}/api/auth/ea_judge_expert_report/export`, params);

// ========================== 评委上报审核（招办） ==========================
export const getReportAuditTasks = () =>
  http.post<{ list: JudgeExpert.CollegeReportTask[] }>(`${PROXY_TAG}/api/auth/ea_judge_expert_report/audit/tasks`, {});
export const getReportAuditDirections = (params: { judgeTaskId: string }) =>
  http.post<{ list: JudgeExpert.CollegeReportDirection[] }>(
    `${PROXY_TAG}/api/auth/ea_judge_expert_report/audit/directions`,
    params
  );
export const getReportAuditList = (params: JudgeExpert.ReqReportAuditList) =>
  http.post<{ list: JudgeExpert.CollegeReportExpert[] }>(`${PROXY_TAG}/api/auth/ea_judge_expert_report/audit/list`, params);
export const exportReportAudit = (params: JudgeExpert.ReqReportAuditList) =>
  http.download(`${PROXY_TAG}/api/auth/ea_judge_expert_report/audit/export`, params);
export const operateReportAudit = (params: { reportId: string; auditState: "2" | "3" }) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_expert_report/audit/operate`, params);
export const resetReportAudit = (params: { reportId: string }) =>
  http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_expert_report/audit/reset`, params);

// ========================== 评委查询 =========================
// 获取评委查询列表
export const getJudgeQueryList = (params: ReqPage<JudgeExpert.ReqJudgeResponseList>) => {
  return http.post<{
    list: { attendNum: string; unattendNum: string; notSureNum: string; list: JudgeExpert.ResJudgeResponseList[] };
    page: any;
  }>(`${PROXY_TAG}/api/auth/ea_judge_draw_detail/list`, params);
};

// 重置评委应答
export const restJudgeResponse = (params: { id: string }) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_draw_detail/reset`, params);
};

// ========================== 专家评价管理 =========================
// 获取专家评价列表
export const getExpertEvaluationList = (params: { expertId: string }) => {
  return http.post<ResPage<JudgeExpert.ResExpertEvaluationList>>(`${PROXY_TAG}/api/auth/ea_judge_experts_appraisal/list`, params);
};

// 获取专家可评价任务列表
export const getCanEvaluateJudgeTaskList = (params: { expertId: string }) => {
  return http.post<{ list: JudgeExpert.ResCanEvaluateJudgeTaskList[] }>(
    `${PROXY_TAG}/api/auth/ea_judge_draw_detail/list_by_expert`,
    params
  );
};

// 新增专家评价
export const addExpertEvaluation = (params: JudgeExpert.ReqAddExpertEvaluation) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_experts_appraisal/add`, params);
};

// 编辑专家评价
export const editExpertEvaluation = (params: JudgeExpert.ReqEditExpertEvaluation) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_experts_appraisal/edit`, params);
};

// 删除专家评价
export const deleteExpertEvaluation = (params: { appraisalId: string }) => {
  return http.post<Common.emptyObject>(`${PROXY_TAG}/api/auth/ea_judge_experts_appraisal/del`, params);
};

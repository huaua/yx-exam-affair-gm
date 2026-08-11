import { ReqPage, ResPage, Classroom, Common } from "@/api/interface/index";
import { PROXY_TAG } from "@/api/config/servicePort";
import http from "@/api";

/**
 * @name 考场管理模块
 */

// ========================== 教室管理 =========================
// 获取教室列表
export const getClassroomList = (params: ReqPage<Classroom.ReqClassroomList>) => {
  return http.post<ResPage<Classroom.ResClassroomList>>(PROXY_TAG + `/api/auth/ea_classroom_info/list`, params);
};

// 切换状态
export const changeClassroomStatus = (params: { roomId: string; state: string }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_classroom_info/change`, params);
};

// 新增教室
export const addClassroom = (params: Classroom.ReqAddClassroom) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_classroom_info/add`, params);
};

// 编辑教室
export const editClassroom = (params: Classroom.ReqEditClassroom) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_classroom_info/edit`, params);
};

// 导入教室
export const importClassroom = (params: any) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_classroom_info/import`, params);
};

// ========================== 考场征集管理 =========================
// 获取考场征集任务列表
export const getRecruitManageList = (params: ReqPage<Classroom.ReqRecruitManageList>) => {
  return http.post<ResPage<Classroom.ResRecruitManageList>>(PROXY_TAG + `/api/auth/ea_task_info/list`, params, { cancel: false });
};

// 切换发布状态
export const changeRecruitTaskStatus = (params: { taskId: string; state: string }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_task_info/change`, params);
};

// 新增考场征集任务
export const addRecruitTask = (params: Classroom.ReqAddRecruitTask) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_task_info/add`, params);
};

// 编辑考场征集任务
export const editRecruitTask = (params: Classroom.ReqEditRecruitTask) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_task_info/edit`, params);
};

// 确认征集-获取任务关联的考场列表
export const getClassroomListByTaskId = (params: Classroom.ReqClassroomListByTaskId) => {
  return http.post<{ list: Classroom.ResClassroomListByTaskId }>(PROXY_TAG + `/api/auth/ea_task_room_info/list`, params, {
    cancel: false
  });
};

// 批量确认征集（已上报→已征集，设置监考老师数量；不在列表中的已征集教室自动清除）
export const confirmRecruitBatch = (params: {
  taskId: string | number;
  roomList: { id: string | number; tchNum: string | number }[];
}) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_task_room_info/edit`, params, { showErrorMsg: false });
};

// 确认征集（单条：已上报→已征集）
export const confirmRecruit = (params: { id: string; tchNum: string | number }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_task_room_info/confirm`, params);
};

// 重置征集确认（已征集→已上报，清空监考老师数量）
export const resetRecruit = (params: { id: string }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_task_room_info/reset`, params);
};

// 批量重置征集确认（已征集→已上报，清空监考老师数量）
export const resetRecruitBatch = (params: { taskId: string | number; roomList: { id: string | number }[] }) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_task_room_info/reset_batch`, params);
};

// 导出任务教室征集关联（Excel）
export const exportTaskRoomInfo = (params: { taskId: string | number }) => {
  return http.download(PROXY_TAG + `/api/auth/ea_task_room_info/export`, params);
};

// ========================== 考场征集 =========================
// 获取考场征集列表
export const getClassroomRecruitList = (params: ReqPage<Classroom.ReqClassroomRecruitList>) => {
  return http.post<{ list: Classroom.ResClassroomRecruitList[]; obj: number }>(
    PROXY_TAG + `/api/auth/ea_task_room_info/list_all`,
    params
  );
};

// 获取所有已发布的任务列表（搜索任务用于下拉框）
export const getPublishedRecruitTaskAllList = (params = {}) => {
  return http.post<ResPage<Classroom.ResRecruitManageList>>(PROXY_TAG + `/api/auth/ea_task_info/list_all`, params);
};

// 上报考场
export const publishClassroom = (params: Classroom.ReqPublishClassroom) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_task_room_info/add`, params);
};

// 取消上报考场
export const unpublishClassroom = (params: Classroom.ReqPublishClassroom) => {
  return http.post<Common.emptyObject>(PROXY_TAG + `/api/auth/ea_task_room_info/del`, params);
};

// ========================== 使用考场管理 =========================

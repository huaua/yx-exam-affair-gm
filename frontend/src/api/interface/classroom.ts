/**
 * @name 考场管理模块
 */
export namespace Classroom {
  // ========================== 教室管理 =========================
  export interface ReqClassroomList {
    levelCode: string; // 层次
    roomName: string; // 教室名称
    deptId: string; // 所属学院id
    state: string; // 启用状态: 1-正常 2-停用
  }
  export interface ResClassroomList {
    roomId: string; // 主键ID 教室ID
    roomName: string; // 教室名称
    campusName: string; // 所在校区
    levelCode: string; // 所属层次
    groupNum: number; // 组数
    groupCapacity: number; // 按组容量
    capacity: number; // 按位容量
    deptId: string; // 所属学院id
    deptName: string; // 所属学院
    state: string; // 启用状态: 1-正常 2-停用
  }
  export interface ReqAddClassroom {
    roomName: string; // 教室名称
    campusName: string; // 所在校区
    levelCode: string; // 所属层次
    groupNum: number; // 组数
    capacity: number; // 容量
  }
  export interface ReqEditClassroom {
    roomName: string; // 教室名称
    campusName: string; // 所在校区
    levelCode: string; // 所属层次
    groupNum: number; // 组数
    capacity: number; // 容量
  }

  // ========================== 考场征集管理 =========================
  export interface ReqRecruitManageList {
    levelCode?: string; // 层次
    taskDay?: string; // 征集日期
    taskName?: string; // 任务名称
  }
  export interface ResRecruitManageList {
    taskId: string; // 主键ID 任务ID
    taskName: string; // 任务名称
    taskDay: string; // 征集日期
    levelCode: number; // 所属层次
    needNum: number; // 需要考场数量
    collectNum: number; // 已征集数量
    comfirmNum: number; // 已确认数量
    state: string; // 发布状态: 1-未发布 2-已发布
  }
  export interface ReqAddRecruitTask {
    taskName: string; // 任务名称
    taskDay: string; // 征集日期
    levelCode: string; // 所属层次
    needNum: string; // 需要考场数量
  }
  export interface ReqEditRecruitTask {
    taskId: string; // 主键ID，必填
    taskName: string; // 任务名称
    taskDay: string; // 征集日期
    levelCode: string; // 所属层次
    needNum: string; // 需要考场数量
  }
  export interface ReqClassroomListByTaskId {
    taskId: string; // 主键ID，必填
    roomState: string; // 教室状态 1-启用 2-禁用
    state: string; // 征集状态 1-待确认 2-已确认
    roomName?: string; // 教室名称
  }
  export interface ResClassroomListByTaskId {
    id: string;
    taskId: string;
    roomId: string;
    roomName: string; // 教室名称
    campusName: string; // 校区名称
    deptId: string;
    deptName?: string; // 所属学院
    roomState: string;
    state: string; // 征集状态 1-待确认(已上报) 2-已确认(已征集)
    levelCode: string; // 所属层次
    groupNum: number; // 组数
    groupCapacity: number; // 按组容量
    capacity: number; // 按位容量
    tchNum: number | string; // 监考老师数量
  }
  interface roomList {
    id: string;
    tchNum: string; // 监考老师数量
  }
  export interface ReqConfirmRecruit {
    taskId: string; // 主键ID，必填
    roomList: roomList[];
  }

  // ========================== 考场征集 =========================
  export interface ReqClassroomRecruitList {
    taskId: string; // 任务ID
  }
  export interface ResClassroomRecruitList {
    roomId: string;
    roomName: string; // 教室名称
    campusName: string; // 校区名称
    levelCode: string; // 层级（字典获取）
    groupNum: number; // 组数
    groupCapacity: number; // 按组容量
    capacity: number; // 按位容量
    publishState: string; // 上报状态 1-未上报 2-已上报，未确认 3-已上报，已确认
  }
  export interface ReqPublishClassroom {
    taskId: string; // 任务ID
    roomId: string;
  }
  // ========================== 使用考场管理 =========================
}

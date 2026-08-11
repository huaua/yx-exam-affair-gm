/**
 * @name 任务征集模块
 */

export namespace Teacher {
  // ========================== 教职工管理 =========================
  export interface ReqTeacherList {
    tchName: string; // 姓名
    jobNo: string; // 工号
    idCard: string; // 证件号（精准查询）
    sex: string; // 性别（男或女）
    auditStatus: string; // 审核状态 1-未提交 2-待审核 3-审核通过 4-审核不通过
    invigilationExperience: string; // 监考经验 1-有 2-无
    deptId: string; // 所属学院id
  }
  export interface ResTeacherList {
    tchId: string; // 主键ID 教职工ID
    tchName: string; // 姓名
    sex: string; // 性别（男或女）
    jobNo: string; // 工号
    idCard: string; // 证件号
    deptId: string; // 所属学院id
    deptName?: string; // 所属学院名称
    duties: string; // 职务
    bankcard: string; // 工商银行卡号
    profOffice: string; // 专业或行政
    education: string; // 专业背景
    auditStatus: string; // 审核状态 1-未提交 2-待审核 3-审核通过 4-审核不通过
    auditRemark: string; // 审核不通过原因
    invigilationExperience: string; // 监考经验 1-有 2-无
    enabledStatus: string; // 启用状态 1-启用 2-禁用
    dateFrom: number; // 1-学院 2-招办
  }
  export interface ReqAddTeacher {
    tchName: string; // 姓名
    sex: string; // 性别（男或女）
    jobNo: string; // 工号
    idCard: string; // 证件号
    deptId: string; // 所属学院id
    duties: string; // 职务
    bankcard: string; // 工商银行卡号
    profOffice: string; // 专业或行政
    education: string; // 专业背景
    submitAction?: "draft" | "submit"; // 学院操作：暂存/提交
  }
  export interface ReqEditTeacher {
    tchId: string; // 主键ID 教职工ID
    tchName: string; // 姓名
    sex: string; // 性别（男或女）
    jobNo: string; // 工号
    idCard: string; // 证件号
    deptId: string; // 所属学院id
    duties: string; // 职务
    bankcard: string; // 工商银行卡号
    profOffice: string; // 专业或行政
    education: string; // 专业背景
    submitAction?: "draft" | "submit"; // 学院操作：暂存/提交
  }
  export interface ReqAuditTeacher {
    tchId: string;
    action: "pass" | "reject";
    reason?: string;
  }
  export interface ReqChangeTchEnabled {
    tchId: string;
    enabledStatus: string; // 1-启用 2-禁用
  }
  export interface ReqChangeTchEnabledBatch {
    tchIds: number[]; // 教职工ID列表
    enabledStatus: string; // 1-启用 2-禁用
  }
  export interface ResAuditHistory {
    historyId: string;
    tchId: string;
    deptId: string;
    auditStatus: string;
    snapshotData: string;
    changeBy: string;
    changeTime: string;
  }
  export interface ResInvigilationRecord {
    tchArrangeId: string;
    tchReportId: string;
    tchTaskId: string;
    tchTaskName: string;
    taskTime: string;
    roomReportId: string;
    roomName: string;
    arrangeNoStr: string;
    taskStartTime: string;
  }

  // ========================== 监考老师征集管理 =========================
  export interface ReqTeacherRecruitManageList {
    tchTaskName?: string; // 考试任务名称
    state?: string; // 发布状态: 1-未发布 2-已发布
  }
  export interface ResTeacherRecruitManageList {
    tchTaskId: string; // 主键ID 教师征集任务ID
    tchTaskName: string; // 考试任务名称
    taskDetailId: string; // 教师征集详情ID
    taskTimeRemark: string; // 征集日期说明
    taskStartTime: string; // 征集开始时间
    taskEndTime: string; // 征集结束时间
    state: string; // 发布状态: 1-未发布 2-已发布
    needNum: string; // 总征集数量
    collectNum: string; // 已征集数量
    auditPassNum: string; // 审核通过数
  }
  export interface ReqCheckTeacherRecruitTaskName {
    tchTaskName: string; // 考试任务名称
    taskStartTime: string; // 征集开始时间
    taskEndTime: string; // 征集结束时间
  }
  interface TaskDetail {
    taskTime: string; // 征集日期, 必填
    deptId: string; // 学院id, 必填
    deptName: string; // 学院名称
    needNum: number; // 征集人数, 必填
    forceFlag: string; // 是否强制征集: 1-强制 2-不强制
  }
  export interface ReqAddTeacherRecruitTask {
    tchTaskName: string; // 考试任务名称
    taskStartTime: string; // 征集开始时间
    taskEndTime: string; // 征集结束时间
    tchTaskDetail: TaskDetail[]; // 征集详情list
  }
  export interface ReqEditTeacherRecruitTask {
    tchTaskId: string; // 考试id
    tchTaskDetail: TaskDetail[]; // 征集详情list
  }
  interface TaskDetailForEdit extends TaskDetail {
    taskDetailId: string; // 主键ID 教师征集详情ID
    tchTaskId: string; // 教师征集任务ID
    collectNum: number; // 已征集数量
  }
  export interface ResEditTeacherRecruitTaskDetail {
    day: string;
    list: TaskDetailForEdit[]; // 征集详情list
  }

  // ========================== 征集老师审核 =========================
  export interface ReqTeacherRecruitViewList {
    tchTaskId?: string; // 征集任务ID
    taskTime?: string; // 征集日期
    deptId?: string; // 学院id
    sex?: string; // 性别（男/女）
    state?: string; // 征集状态 1-已上报,待审核 2-待分配 3-不通过 4-已分配
    invigilationExperience?: string; // 监考经验 1-有 2-无
    jobNo?: string; // 工号
  }
  export interface ResTeacherRecruitViewList {
    tchId: string;
    tchName: string; // 姓名
    reportId: string; // 上报Id
    sex: string; // 性别
    jobNo: string; // 工号
    idCard?: string; // 证件号
    duties: string; // 职务
    profOffice: string; // 专业或行政
    education?: string; // 专业背景
    invigilationExperience?: string; // 监考经验 1-有 2-无
    invigilationCount?: number; // 监考次数
    taskTime: string; // 征集日期
    deptId: string; // 学院id
    deptName?: string; // 学院名称
    state: string; // 征集状态 1-已上报,待审核 2-待分配 3-不通过 4-已分配
    remark: string; // 备注
  }

  // ========================== 监考老师上报 =========================
  export interface ReqCanSubmitTeacherList {
    tchTaskId: string; // 征集任务ID
    taskTime?: string; // 征集日期
    tchNameOrJobNo?: string; // 姓名或者工号搜索条件（可选）
  }
  export interface ResCanSubmitTeacherList {
    tchId: string; // 主键ID 教职工ID
    tchName: string; // 姓名
    sex: string; // 性别（男或女）
    jobNo: string; // 工号
    reportTchList: string[]; // 已填报日期
    remark: string; // 备注
  }
  export interface ResCurTaskDateList {
    taskDetailId: string; // 教师征集详情ID
    taskTime: string; // 征集日期
    deptId: string; // 学院id
    deptName: string; // 学院名称
    collectNum: number; // 已征集数量
    needNum: number; // 征集人数
    forceFlag: string; // 是否强制征集: 1-强制 2-不强制
  }
  interface AddTeacherItem {
    taskDetailId: string; // 主键ID 教师征集详情ID
    tchTaskId: string; // 征集任务ID
    taskTime: string; // 征集日期
  }
  export interface ReqSubmitTeacher {
    tchId: string;
    tchName: string; // 姓名
    jobNo: string; // 工号
    reportTaskDetailList: AddTeacherItem[];
  }
  interface AppendTeacherItem {
    tchId: string; // 教师ID
    tchName: string;
    jobNo: string;
    sex: string;
  }
  export interface ReqAppendTeacher {
    taskDetailId: string;
    tchList: AppendTeacherItem[];
  }
  export interface ResGetCacheSubmitTeacherItem extends AppendTeacherItem {
    taskTimes: string[];
  }

  // ========================== 监考老师上报查看 =========================
  export interface ResTeacherSubmitViewList {
    reportId: string;
    tchTaskId: string;
    taskDetailId: string;
    tchName: string; // 姓名
    sex: string; // 性别
    jobNo: string; // 工号
    duties: string; // 职务
    profOffice: string; // 专业
    taskTime: string; // 征集日期
    invigilationExperience: string; // 监考经验 1-有 2-无
    state: string;
    remark: string; // 备注
  }
}

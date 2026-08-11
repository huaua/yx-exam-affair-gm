export namespace JudgeExpert {
  // ========================== 专家库管理 =========================
  export interface ReqExpertDatabaseList {
    name?: string;
    sex?: string;
    idCard?: string;
    expertCategory?: string;
    type?: string;
    unitName?: string;
    deptId?: string;
    title?: string;
    auditStatus?: string;
    allowDraw?: string;
    initEdu?: string;
    finalEdu?: string;
    goodSubjects?: string;
    isEvaluated?: string;
  }
  export interface ResExpertDatabaseList {
    expertId: string;
    name: string;
    sex: string;
    phoneNo: string;
    idCard: string;
    expertCategory: string;
    deptId: string;
    deptName?: string;
    unitName?: string;
    type: string;
    title: string;
    initEdu: string;
    initMajor: string;
    finalEdu?: string;
    finalMajor?: string;
    goodSubjects: string;
    bankCard?: string;
    intro?: string;
    allowDraw: string;
    auditStatus: string;
    auditRemark?: string;
    evaluation?: string;
    hasEvaluation?: boolean;
    dataFrom?: string;
  }
  export interface ReqAddExpert extends Omit<ResExpertDatabaseList, "expertId" | "deptName" | "hasEvaluation"> {
    submitAction?: "draft" | "submit";
  }
  export interface ReqEditExpert extends ReqAddExpert {
    expertId: string;
  }
  export interface ReqAuditExpert {
    expertId: string;
    action: "pass" | "reject";
    reason?: string;
  }
  export interface ResExpertAuditHistory {
    historyId: string;
    expertId: string;
    auditStatus: string;
    snapshotData: string;
    changeTime: string;
  }

  // ========================== 评委征集任务管理 =========================
  export interface ReqJudgeRecruitManageList {
    judgeTaskName?: string; // 任务名称
    judgeTaskType?: string; // 类型：1-博士 2-硕士 3-本科 4-附中
    state?: string; // 抽取状态：1-未开始 2-抽取中 3-已完成
  }
  export interface ResJudgeRecruitManageList {
    judgeTaskId: string; // 主键ID 任务ID
    judgeTaskName: string; // 任务名称
    taskYear: string;
    taskMode: string;
    scoreRounds: string;
    subjectCount: string;
    expertsPerRound: string;
    publishState: string;
    taskStartTime: string; // 开始时间
    taskEndTime: string; // 结束时间
    judgeTaskType: string; // 类型
    judgeTaskTypeStr: string; // 类型描述
    state: string; // 抽取状态
    onCampusTotalNum: string; // 校内评委需求量
    onCampusAttendNum: string; // 校内评委参加数量
    offCampusTotalNum: string; // 校外评委需求量
    offCampusAttendNum: string; // 校外评委参加数量
    reportCount?: string; // 上报数（学院上报模式下的已上报记录数）
    requirementCount?: string; // 需求数（直接抽取=轮次*科目*每轮专家数；上报模式=各方向专家需求总和）
    drawCount?: string; // 已抽取评委数（直接抽取模式下 ea_judge_draw_selection 评分专家数量）
  }

  export interface ReqAddJudgeRecruitTask {
    judgeTaskName: string; // 任务名称
    taskYear: string;
    taskMode: string;
    scoreRounds?: string;
    subjectCount?: string;
    expertsPerRound?: string;
    judgeTaskType: string; // 类型：1-博士 2-硕士 3-本科 4-附中
    taskStartTime: string; // 开始时间
    taskEndTime: string; // 结束时间
    drawSubject: string; // 是否按科目抽取 1-是 2-否
  }

  export interface ReqEditJudgeRecruitTask extends ReqAddJudgeRecruitTask {
    judgeTaskId: string; // 任务ID
  }

  export interface JudgeTaskReportRequirement {
    requirementId: string;
    judgeTaskId: string;
    directionCode: string;
    directionName: string;
    deptName: string;
    expertNum: string | number;
    reportedNum: string;
  }

  // ========================== 学院评委上报 / 评委上报审核 =========================
  export interface CollegeReportTask {
    judgeTaskId: string;
    judgeTaskName: string;
    taskName: string; // “任务年份”年“任务名称”
    taskYear: string;
    judgeTaskType: string;
    taskStartTime: string;
    taskEndTime: string;
  }

  export interface CollegeReportDirection {
    requirementId: string;
    judgeTaskId: string;
    directionCode: string;
    directionName: string;
    expertNum: string;
    reportedNum: string;
    approvedNum?: string;
    rejectedNum?: string;
    draftNum?: string;
  }

  export interface CollegeReportExpert extends ResExpertDatabaseList {
    reportId: string;
    requirementId: string;
    judgeTaskId: string;
    submitState: string;
    auditState: string;
    directionCode?: string;
    directionName?: string;
    taskName?: string;
    taskStartTime?: string;
    taskEndTime?: string;
  }

  export interface ReqReportAuditList {
    judgeTaskId: string;
    deptId?: string;
    directionCode?: string;
    name?: string;
    sex?: string;
    title?: string;
    auditState?: string;
  }

  // ========================== 评委应答 =========================
  export interface ReqJudgeResponseList {
    judgeTaskId?: string; // 任务id
    confirmState?: string; // 确认状态 1-未应答 2-确认参加 3-不参加
  }
  export interface ResJudgeResponseList {
    id: string; // 主键ID 任务ID
    expertName: string; // 专家姓名
    sex: string; // 性别
    phoneNo: string; // 手机号
    type: string; // 专家类型
    typeStr: string; // 专家类型描述
    goodSubjects: string; // 擅长科目名称
    taskTimeRemark: string; // 征集日期
    deptId: string; // 部门id
    deptName: string; // 学院名称
    drawTimes: string; // 抽取次数
    confirmState: string; // 状态 1-未应答 2-确认参加 3-不参加
  }
  // ========================== 评委查询 =========================
  // 同 评委应答 数据结构

  // ========================== 专家评价管理 =========================
  export interface ResExpertEvaluationList {
    appraisalId: string; // 评价ID
    judgeTaskId: string; // 任务ID
    judgeTaskName: string; // 任务名称
    taskTimeRemark: string; // 评分日期
    appraisalContent: string; // 评价内容
  }
  export interface ResCanEvaluateJudgeTaskList {
    judgeTaskId: string; // 任务ID
    judgeTaskName: string; // 任务名称
  }
  export interface ReqAddExpertEvaluation {
    expertId: string; // 专家ID
    judgeTaskId: string; // 任务名称
    appraisalContent: string; // 评价内容
  }
  export interface ReqEditExpertEvaluation {
    appraisalId: string; // 评价ID
    appraisalContent: string; // 评价内容
  }
}

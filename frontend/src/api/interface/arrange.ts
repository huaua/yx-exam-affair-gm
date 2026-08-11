/**
 * @name 考场编排模块
 */

export namespace Arrange {
  // ========================== 考场排考编排 =========================
  export interface ReqClassroomArrangeList {
    roomTaskId: string; // 考场任务ID
    deptName: string; // 学院名称
    profName: string; // 专业名称
    directionName: string; // 方向名称
    arrangeState: string; // 考场编排状态: 1-待编排 2-已编排
  }
  export interface ResClassroomArrangeList {
    infoId: string; // 主键ID 数据ID
    roomTaskId: string; // 考场任务ID
    taskDay: string; // 征集日期
    deptId: string; // 所属学院id
    deptName: string; // 学院名称
    profName: string; // 专业名称
    directionName: string; // 方向名称
    arrangeType: string; // 安排类型: 1-按组 2-按位
    stuNum: number; // 考生数量
    arrangeStuNum: number; // 已编排考生数量
    arrangeRoomNum: number; // 已编排考场数量
    arrangeState: string; // 考场编排状态: 1-待编排 2-已编排
  }
  export interface ResCanArrangeClassroomList {
    id: string;
    taskId: string;
    roomId: string;
    roomName: string;
    groupCapacity: string; // 按组容量
    deptId: string;
    roomState: string;
    state: string;
    levelCode: string;
    groupNum: string;
    capacity: string;
  }
  export interface ResArrangeDetail {
    infoId: string; // 主键ID 数据ID
    roomTaskId: string; // 考场任务ID
    roomTaskName: string; // 考场任务名称
    taskDay: string; // 征集日期
    deptId: string; // 所属学院id
    deptName: string; // 学院名称
    profName: string; // 专业名称
    directionName: string; // 方向名称
    arrangeType: string; // 安排类型: 1-按组 2-按位
    stuNum: string; // 考生数量
    arrangeStuNum: string; // 已编排考生数量
    arrangeRoomNum: string; // 已编排考场数量
    arrangeState: string; // 考场编排状态: 1-待编排 2-已编排
    maxArrangeRoomNoStr: string; // 最大考场编号字符串
    maxArrangeRoomNo: number; // 最大考场编号
  }
  export interface ResArrangeDetailList {
    roomArrangeId: string; // 主键ID 考场安排ID
    infoId: string; // 数据ID
    roomReportId: string; // 考场上报ID
    id?: string; // 同 roomReportId(前端冗余字段)
    roomTaskId: string; // 任务ID
    taskId?: string; // 同 roomTaskId(前端冗余字段)
    roomId: string; // 教室ID
    roomName: string; // 教室名称
    stuNum: string | number; // 考生数量
    roomCapacity: string; // 考场实际容量
    arrangeNoStr: string; // 考场安排编号(字符串)
    arrangeNo: string; // 考场安排编号(数值)
    needTchNum: string; // 需要监考老师数量
    tchNum?: string; // 同 needTchNum(前端冗余字段)
    groupNum: string; // 组数
    groupCapacity: string; // 按组容量
    capacity: string; // 按位容量
  }
  interface SaveArrangeDetailItem {
    roomReportId: string; // 考场上报ID
    roomTaskId: string; // 任务ID
    roomId: string; // 教室ID
    roomName: string; // 教室名称
    stuNum: string; // 考生数量
    roomCapacity: string; // 考场实际容量
    arrangeNoStr: string; // 考场安排编号(字符串)
    needTchNum: string; // 需要监考老师数量
  }
  export interface ReqSaveArrangeDetailList {
    infoId: string;
    arrangeRoomList: SaveArrangeDetailItem[];
  }

  // ========================== 考场监考分配 =========================
  export interface ReqTeacherAssignList {
    roomTaskId: string; // 考场任务ID
    roomArrangeId?: string; // 考场编排id 单条查询时必填
    assignmentState?: string; // 分配状态 1-未分配 2-已分配
    deptId?: string; // 学院id
    roomName?: string; // 教室名称
    profName?: string; // 专业名称
    directionName?: string; // 方向名称
  }
  export interface ResTeacherAssignList {
    roomReportId: string; // 考场上报ID
    roomId: string; // 教室ID
    roomName: string; // 教室名称
    roomArrangeId: string; // 考场编排id
    profDirection: string; // 专业方向
    taskDay: string; // 考场日期
    tchNum: string; // 老师数量
    arrangeTchReportId: string[]; // 编排的老师上报id
    arrangeTchName: string[]; // 编排的老师名称
    arrangeNoStr: string; // 考场安排编号(字符串)
    assignmentState: string;
  }
  export interface ReqSaveTeacherAssignItem {
    roomTaskId: string; // 考场任务ID
    roomReportId: string; // 考场上报ID
    tchReportIdList?: string[]; // 教职工上报id
  }
  // ========================== 编排查看 =========================
  export interface ReqArrangeViewList {
    roomTaskId: string; // 考场任务ID 必填
    deptId?: string; // 学院id
    profName?: string; // 专业名称 支持模糊
    directionName?: string; // 方向名称 支持模糊
    roomName?: string; // 教室名称 支持模糊
    tchName?: string; // 老师名称 支持模糊
  }
  export interface ResArrangeViewList {
    infoId: string; // 专业数据id
    deptId: string; // 学院id
    deptName: string; // 学院名称
    profName: string; // 专业名称
    directionName: string; // 方向名称
    arrangeType: string; // 安排类型: 1-按组 2-按位
    roomName: string; // 教室名称
    groupNum: string; // 组数
    stuNum: string; // 考生数量
    arrangeNoStr: string; // 考场安排编号(字符串)
    tchNum: string; // 需要的监考老师数量
    arrangeTchName: string; // 编排的老师名称
    assignmentState: string; // 分配状态 1-未分配 2-已分配
  }
}

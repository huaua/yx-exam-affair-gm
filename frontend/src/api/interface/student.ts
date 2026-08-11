/**
 * @name 考生管理模块
 */

export namespace Student {
  // ========================== 考生管理 =========================
  export interface ReqStudentList {
    kaoShiID?: string; // 考试id
    shengFenMC?: string; // 省份
    xingMing?: string; // 姓名
    shenFenZH?: string; // 身份证号
    zhunKaoZH?: string; // 准考证号
  }
  export interface ResStudentList {
    kaoShiID: string; // 考试id
    shengFenMC: string; // 省份
    xingMing: string; // 姓名
    shenFenZH: string; // 身份证号
    shouJi: string; // 手机号
    zhunKaoZH: string; // 准考证号
    zhuanYeMC: string; // 专业
    kaoShiRQ: string; // 考试时间
  }
  export interface ReqSyncStudentInfo {
    kaoShiID: string; // 考试id
    kaoShiMC: string; // 考试名称
  }
  export interface ResSyncStudentInfo {
    state: string; // 同步状态  0-未同步 1-同步中 2-完成 2-失败
    kaoShiID: string; // 考试id
    kaoShiMC: string; // 考试名称
    percent: string; // 进度
    errMsg?: string; // 同步失败消息
  }
  export interface ResExamList {
    kaoShiID: string; // 考试id
    kaoShiMC: string; // 考试名称
    kaoShiND: string; // 考试年度
  }
}

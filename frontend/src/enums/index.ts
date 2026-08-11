/**
 * @description：角色关键字
 */
export enum RoleKeyEnum {
  ADMISSIONS = "admissions",
  DEPARTMENT = "department"
}

/**
 * @description：评委任务状态
 */
export enum JudgeTaskState {
  NOT_STARTED = "1", // 未开始
  EXTRACTING = "2", // 抽取中
  COMPLETED = "3" // 已完成
}

/**
 * @description：评委应答状态
 */
export enum JudgeResponseState {
  NOT_RESPONDED = "1", // 待定
  ACCEPTED = "2", // 参加
  DECLINED = "3" // 不参加
}

/**
 * @description：考生信息同步状态
 */
export enum StudentSyncEnum {
  NOT_STARTED = "0", // 未开始
  IN_PROGRESS = "1", // 进行中
  SUCCESS = "2", // 成功
  FAILED = "3" // 失败
}

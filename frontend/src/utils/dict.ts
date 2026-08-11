// ? 系统全局字典

import { JudgeTaskState, JudgeResponseState } from "@/enums";

/**
 * @description：用户性别
 */
export const genderType = [
  { label: "男", value: "1" },
  { label: "女", value: "2" }
];

/**
 * @description：用户状态
 */
export const userStatus = [
  { label: "启用", value: 1, tagType: "success" },
  { label: "禁用", value: 0, tagType: "danger" }
];

/**
 * @description：征集状态
 */
export const recruitStatus = [
  { label: "未发布", value: "1" },
  { label: "已发布", value: "2" }
];

/**
 * @description：启用禁用状态
 */
export const onOrOffStatus = [
  { label: "启用", value: "1" },
  { label: "禁用", value: "2" }
];

/**
 * @description：招办教室
 */
export const admissionClassroomType = [{ label: "公共教室", value: "-100" }];

/**
 * @description：招办专家来源
 */
export const admissionExpertType = [{ label: "校外", value: "-100" }];

/**
 * @description：学历 (与后端保持一致)
 */
export const educationType = [
  { label: "本科", value: "1" },
  { label: "硕士", value: "2" },
  { label: "博士", value: "3" }
];

/**
 * @description：专家类型 (与后端保持一致)
 */
export const expertType = [
  { label: "本科", value: "1" },
  { label: "研究生", value: "2" },
  { label: "博士", value: "3" },
  { label: "附中", value: "4" }
];

/**
 * @description：评委抽取状态
 */
export const judgeDrawStatus = [
  { label: "未开始", value: JudgeTaskState.NOT_STARTED },
  { label: "征集中", value: JudgeTaskState.EXTRACTING },
  { label: "已完成", value: JudgeTaskState.COMPLETED }
];

/**
 * @description：评委应答状态
 */
export const judgeResponseStatus = [
  { label: "待定", value: JudgeResponseState.NOT_RESPONDED, tagType: "default" },
  { label: "参加", value: JudgeResponseState.ACCEPTED, tagType: "success" },
  { label: "不参加", value: JudgeResponseState.DECLINED, tagType: "danger" }
];

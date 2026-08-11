-- 学院评委专家上报功能数据库迁移。
-- 前置条件：已执行 update_260728_judge_recruit_task.sql 迁移脚本。
-- 历史保留设计：唯一键不再对“已删除(deleted=2)”记录做唯一约束，改用函数唯一索引
--   uk_active (requirement_id, expert_id, CASE WHEN deleted=1 THEN 1 ELSE NULL END)
-- 仅对“正常记录(deleted=1)”强制唯一（保证同一专家同一要求下只有一条有效上报，统计/显示不受影响），
-- 已删除记录映射到 NULL 被忽略，可保留同一专家的多次删除历史；delete_time 记录每次软删时间用于区分。

-- 根据学院名称补充上报要求中的学院ID。
UPDATE ea_judge_task_report_requirement r
JOIN ea_dept_info d ON d.dept_name = r.dept_name AND d.deleted = 1
SET r.dept_id = d.dept_id
WHERE r.deleted = 1 AND (r.dept_id IS NULL OR r.dept_id = 0);

-- 创建学院评委专家上报记录表。
CREATE TABLE IF NOT EXISTS ea_judge_expert_report (
  report_id bigint NOT NULL AUTO_INCREMENT COMMENT '上报记录ID',
  judge_task_id bigint NOT NULL COMMENT '评委征集任务ID',
  requirement_id bigint NOT NULL COMMENT '研究方向上报要求ID',
  dept_id bigint NOT NULL COMMENT '上报学院ID',
  expert_id bigint NOT NULL COMMENT '上报专家ID',
  submit_state tinyint NOT NULL DEFAULT 1 COMMENT '提交状态：1-未提交，2-已提交',
  audit_state tinyint NOT NULL DEFAULT 1 COMMENT '审核状态：1-待审核，2-审核通过，3-审核不通过',
  replaced_from bigint DEFAULT NULL COMMENT '被替换的原上报记录ID',
  deleted tinyint NOT NULL DEFAULT 1 COMMENT '删除标记：1-正常，2-已删除',
  delete_time datetime DEFAULT NULL COMMENT '删除时间，软删时记录，用于区分同一专家的多次删除历史',
  create_by bigint DEFAULT NULL COMMENT '创建人ID',
  create_time datetime DEFAULT NULL COMMENT '创建时间',
  update_by bigint DEFAULT NULL COMMENT '更新人ID',
  update_time datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (report_id),
  UNIQUE KEY uk_active (requirement_id, expert_id, (CASE WHEN deleted=1 THEN 1 ELSE NULL END)) COMMENT '仅对正常记录(deleted=1)强制唯一，软删记录可保留多次历史',
  KEY idx_task_dept_state (judge_task_id, dept_id, submit_state, audit_state, deleted),
  KEY idx_requirement_state (requirement_id, submit_state, audit_state, deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学院评委专家上报记录表';
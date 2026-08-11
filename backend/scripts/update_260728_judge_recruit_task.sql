ALTER TABLE `ea_judge_task_info`
  MODIFY COLUMN `judge_task_name` varchar(50) NOT NULL COMMENT '评委任务名称',
  ADD COLUMN `task_year` int NOT NULL DEFAULT 2026 COMMENT '任务年份' AFTER `judge_task_name`,
  ADD COLUMN `task_mode` tinyint NOT NULL DEFAULT 1 COMMENT '任务模式:1-直接抽取 2-学院上报' AFTER `task_year`,
  ADD COLUMN `score_rounds` int NOT NULL DEFAULT 1 COMMENT '评分轮次' AFTER `draw_times`,
  ADD COLUMN `subject_count` int NOT NULL DEFAULT 1 COMMENT '科目数量' AFTER `score_rounds`,
  ADD COLUMN `experts_per_round` int NOT NULL DEFAULT 1 COMMENT '每轮专家数' AFTER `subject_count`,
  ADD COLUMN `publish_state` tinyint NOT NULL DEFAULT 2 COMMENT '发布状态:1-未发布 2-已发布' AFTER `experts_per_round`,
  ADD INDEX `idx_judge_task_year_name` (`task_year`, `judge_task_name`, `deleted`);

UPDATE `ea_judge_task_info`
SET `task_year` = YEAR(`task_start_time`),
    `task_mode` = 1,
    `publish_state` = 2
WHERE `deleted` = 1;

CREATE TABLE `ea_judge_task_report_requirement` (
  `requirement_id` bigint NOT NULL AUTO_INCREMENT COMMENT '专家数量要求ID',
  `judge_task_id` bigint NOT NULL COMMENT '评委任务ID',
  `direction_code` varchar(64) NOT NULL COMMENT '方向编码',
  `direction_name` varchar(128) NOT NULL COMMENT '方向名称',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `dept_name` varchar(128) NOT NULL COMMENT '部门名称',
  `expert_num` int NOT NULL COMMENT '专家数量',
  `reported_num` int NOT NULL DEFAULT 0 COMMENT '已上报专家数',
  `deleted` tinyint NOT NULL DEFAULT 1 COMMENT '删除标记：1-正常的；2-已删除的；',
  `create_by` bigint DEFAULT NULL COMMENT '创建人',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新人',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`requirement_id`) COMMENT '主键',
   UNIQUE KEY `uk_task_direction_dept` (`judge_task_id`,`direction_code`,`dept_name`,`deleted`) COMMENT '唯一索引:评委任务ID+方向编码+部门名称+删除标志'
) ENGINE=InnoDB COMMENT='学院上报专家数量要求';

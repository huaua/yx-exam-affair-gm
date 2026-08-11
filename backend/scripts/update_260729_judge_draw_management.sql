-- 评委抽取管理：轮次科目配置与暂存专家
ALTER TABLE `ea_judge_draw_records`
  ADD COLUMN `required_num` int NOT NULL DEFAULT 0 COMMENT '评分专家计划数量' AFTER `subject_name`,
  ADD COLUMN `submit_state` tinyint NOT NULL DEFAULT 1 COMMENT '1-未提交 2-已提交' AFTER `required_num`;

CREATE TABLE IF NOT EXISTS `ea_judge_draw_selection` (
  `selection_id` bigint NOT NULL AUTO_INCREMENT,
  `rules_id` bigint NOT NULL,
  `judge_task_id` bigint NOT NULL,
  `draw_times` int NOT NULL,
  `expert_id` bigint NOT NULL,
  `expert_role` tinyint NOT NULL COMMENT '1-评分专家 2-备用专家',
  `submit_state` tinyint NOT NULL DEFAULT 1 COMMENT '1-未提交 2-已提交',
  `deleted` tinyint NOT NULL DEFAULT 1,
  `create_by` bigint DEFAULT NULL,
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_by` bigint DEFAULT NULL,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`selection_id`),
  UNIQUE KEY `uk_draw_selection_rule_expert` (`rules_id`,`expert_id`),
  KEY `idx_draw_selection_round` (`judge_task_id`,`draw_times`,`expert_id`,`deleted`)
) ENGINE=InnoDB COMMENT='评委抽取管理专家选择记录';

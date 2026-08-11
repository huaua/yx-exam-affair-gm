-- 20260727 招办角色专家管理调整

ALTER TABLE `ea_judge_experts`
  MODIFY COLUMN `type` varchar(64) NULL COMMENT '专家类型，多选，逗号分隔字典值',
  ADD COLUMN `expert_category` tinyint NOT NULL DEFAULT 1 COMMENT '专家类别 1-校内 2-校外' AFTER `dept_id`,
  ADD COLUMN `unit_name` varchar(30) NULL COMMENT '校外专家所在单位' AFTER `expert_category`,
  ADD COLUMN `title` varchar(32) NULL COMMENT '职称字典值' AFTER `type`,
  ADD COLUMN `init_major` varchar(20) NULL COMMENT '初始学历所学专业' AFTER `init_edu`,
  ADD COLUMN `final_major` varchar(20) NULL COMMENT '最终学历所学专业' AFTER `final_edu`,
  ADD COLUMN `audit_status` tinyint NOT NULL DEFAULT 1 COMMENT '审核状态 1-未提交 2-待审核 3-审核通过 4-审核不通过' AFTER `data_from`,
  ADD COLUMN `audit_remark` varchar(500) NULL COMMENT '审核不通过原因' AFTER `audit_status`,
  ADD COLUMN `evaluation` varchar(2000) NULL COMMENT '招办评价' AFTER `audit_remark`;

UPDATE `ea_judge_experts`
SET `expert_category` = IF(`dept_id` = -100, 2, 1),
    `type` = CASE `type` WHEN '1' THEN '3' WHEN '2' THEN '2' WHEN '3' THEN '1' WHEN '4' THEN '4' ELSE `type` END,
    `audit_status` = 3
WHERE `deleted` = 1;

UPDATE `ea_judge_task_info`
SET `judge_task_type` = CASE `judge_task_type` WHEN 1 THEN 3 WHEN 2 THEN 2 WHEN 3 THEN 1 WHEN 4 THEN 4 ELSE `judge_task_type` END
WHERE `deleted` = 1;

CREATE TABLE IF NOT EXISTS `ea_judge_experts_audit_history` (
  `history_id` bigint NOT NULL AUTO_INCREMENT,
  `expert_id` bigint NOT NULL COMMENT '专家ID',
  `dept_id` bigint NOT NULL COMMENT '所属学院ID',
  `audit_status` tinyint NOT NULL COMMENT '快照审核状态',
  `snapshot_data` longtext NOT NULL COMMENT '审核通过数据快照JSON',
  `change_by` bigint NULL COMMENT '修改人',
  `change_time` datetime NOT NULL COMMENT '修改时间',
  PRIMARY KEY (`history_id`),
  KEY `idx_expert_audit_history_expert` (`expert_id`, `change_time`)
) ENGINE=InnoDB COMMENT='专家审核通过历史快照';

INSERT INTO `sys_dict_type` (`dict_name`, `dict_type`, `state`, `deleted`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
SELECT '专家类型', 'biz_expert_type', '1', 1, 1, NOW(), 1, NOW(), '专家类型多选字典'
WHERE NOT EXISTS (SELECT 1 FROM `sys_dict_type` WHERE `dict_type` = 'biz_expert_type' AND `deleted` = 1);

INSERT INTO `sys_dict_data` (`dict_sort`, `dict_label`, `dict_value`, `dict_type`, `list_class`, `state`, `deleted`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
SELECT 1, '本科', '1', 'biz_expert_type', 'default', '1', 1, 1, NOW(), 1, NOW(), '专家类型'
WHERE NOT EXISTS (SELECT 1 FROM `sys_dict_data` WHERE `dict_type` = 'biz_expert_type' AND `dict_value` = '1' AND `deleted` = 1);
INSERT INTO `sys_dict_data` (`dict_sort`, `dict_label`, `dict_value`, `dict_type`, `list_class`, `state`, `deleted`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
SELECT 2, '研究生', '2', 'biz_expert_type', 'default', '1', 1, 1, NOW(), 1, NOW(), '专家类型'
WHERE NOT EXISTS (SELECT 1 FROM `sys_dict_data` WHERE `dict_type` = 'biz_expert_type' AND `dict_value` = '2' AND `deleted` = 1);
INSERT INTO `sys_dict_data` (`dict_sort`, `dict_label`, `dict_value`, `dict_type`, `list_class`, `state`, `deleted`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
SELECT 3, '博士', '3', 'biz_expert_type', 'default', '1', 1, 1, NOW(), 1, NOW(), '专家类型'
WHERE NOT EXISTS (SELECT 1 FROM `sys_dict_data` WHERE `dict_type` = 'biz_expert_type' AND `dict_value` = '3' AND `deleted` = 1);
INSERT INTO `sys_dict_data` (`dict_sort`, `dict_label`, `dict_value`, `dict_type`, `list_class`, `state`, `deleted`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
SELECT 4, '附中', '4', 'biz_expert_type', 'default', '1', 1, 1, NOW(), 1, NOW(), '专家类型'
WHERE NOT EXISTS (SELECT 1 FROM `sys_dict_data` WHERE `dict_type` = 'biz_expert_type' AND `dict_value` = '4' AND `deleted` = 1);

INSERT INTO `sys_dict_type` (`dict_name`, `dict_type`, `state`, `deleted`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
SELECT '专家职称', 'biz_expert_title', '1', 1, 1, NOW(), 1, NOW(), '专家职称字典'
WHERE NOT EXISTS (SELECT 1 FROM `sys_dict_type` WHERE `dict_type` = 'biz_expert_title' AND `deleted` = 1);

INSERT INTO `sys_dict_data` (`dict_sort`, `dict_label`, `dict_value`, `dict_type`, `list_class`, `state`, `deleted`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
SELECT 1, '教授', '1', 'biz_expert_title', 'default', '1', 1, 1, NOW(), 1, NOW(), '专家职称'
WHERE NOT EXISTS (SELECT 1 FROM `sys_dict_data` WHERE `dict_type` = 'biz_expert_title' AND `dict_value` = '1' AND `deleted` = 1);
INSERT INTO `sys_dict_data` (`dict_sort`, `dict_label`, `dict_value`, `dict_type`, `list_class`, `state`, `deleted`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
SELECT 2, '副教授', '2', 'biz_expert_title', 'default', '1', 1, 1, NOW(), 1, NOW(), '专家职称'
WHERE NOT EXISTS (SELECT 1 FROM `sys_dict_data` WHERE `dict_type` = 'biz_expert_title' AND `dict_value` = '2' AND `deleted` = 1);
INSERT INTO `sys_dict_data` (`dict_sort`, `dict_label`, `dict_value`, `dict_type`, `list_class`, `state`, `deleted`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
SELECT 3, '讲师', '3', 'biz_expert_title', 'default', '1', 1, 1, NOW(), 1, NOW(), '专家职称'
WHERE NOT EXISTS (SELECT 1 FROM `sys_dict_data` WHERE `dict_type` = 'biz_expert_title' AND `dict_value` = '3' AND `deleted` = 1);

-- 2026/07/23 教职工管理升级：招办审核流程、启用禁用按钮、监考经验回算
-- 适用已经存在 ea_tch_info、ea_tch_info_audit_history 表的库，未初始化请使用 update_240910.sql

-- 1.确保 ea_tch_info 字段存在（升级库可能缺失列）
ALTER TABLE `ea_tch_info`
  MODIFY COLUMN `audit_status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '审核状态 1-未提交 2-待审核 3-审核通过 4-审核不通过',
  MODIFY COLUMN `enabled_status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '启用状态 1-启用 2-禁用',
  MODIFY COLUMN `invigilation_experience` tinyint(1) NOT NULL DEFAULT 2 COMMENT '监考经验 1-有 2-无',
  ADD COLUMN IF NOT EXISTS `audit_remark` varchar(500) NULL COMMENT '审核不通过原因' AFTER `audit_status`;

-- 2.招办新增/导入的数据默认为审核通过状态
UPDATE `ea_tch_info`
SET `audit_status` = 3
WHERE `date_from` = 2 AND `audit_status` IN (1, 2);

-- 3.重新计算所有教职工的监考经验（依据现有监考分配记录）
UPDATE `ea_tch_info` AS `t`
SET `t`.`invigilation_experience` = 1
WHERE EXISTS (
  SELECT 1 FROM `ea_arrange_tch` AS `a`
  WHERE `a`.`tch_id` = `t`.`tch_id` AND `a`.`deleted` = 1
);

UPDATE `ea_tch_info` AS `t`
SET `t`.`invigilation_experience` = 2
WHERE NOT EXISTS (
  SELECT 1 FROM `ea_arrange_tch` AS `a`
  WHERE `a`.`tch_id` = `t`.`tch_id` AND `a`.`deleted` = 1
) AND (`t`.`invigilation_experience` IS NULL OR `t`.`invigilation_experience` <> 2);

-- 4.确保历史快照表存在
CREATE TABLE IF NOT EXISTS `ea_tch_info_audit_history` (
  `history_id` bigint NOT NULL AUTO_INCREMENT COMMENT '历史记录ID',
  `tch_id` bigint NOT NULL COMMENT '教职工ID',
  `dept_id` bigint NOT NULL COMMENT '所属学院ID',
  `audit_status` tinyint(1) NOT NULL COMMENT '修改前审核状态',
  `snapshot_data` longtext NOT NULL COMMENT '修改前完整数据JSON',
  `change_by` bigint NOT NULL COMMENT '修改人',
  `change_time` datetime NOT NULL COMMENT '修改时间',
  PRIMARY KEY (`history_id`) USING BTREE,
  KEY `idx_tch_history_tch_id` (`tch_id`) USING BTREE,
  KEY `idx_tch_history_dept_id` (`dept_id`) USING BTREE
) ENGINE=InnoDB COMMENT='教职工审核通过数据修改前快照';

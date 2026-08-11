-- 学院端教职工暂存、提交、审核对比流程升级。
ALTER TABLE `ea_tch_info`
  ADD COLUMN `audit_remark` varchar(500) NULL COMMENT '审核不通过原因' AFTER `audit_status`,
  ADD COLUMN `enabled_status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '启用状态 1-启用 2-禁用' AFTER `invigilation_experience`,
  ADD INDEX `idx_tch_enabled_status` (`enabled_status`);

CREATE TABLE `ea_tch_info_audit_history` (
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

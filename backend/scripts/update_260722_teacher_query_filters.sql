-- 已有数据库升级：教职工列表新增审核状态、监考经验筛选字段。
ALTER TABLE `ea_tch_info`
  ADD COLUMN `audit_status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '审核状态 1-未提交 2-待审核 3-审核通过 4-审核不通过' AFTER `education`,
  ADD COLUMN `invigilation_experience` tinyint(1) NOT NULL DEFAULT 2 COMMENT '监考经验 1-有 2-无' AFTER `audit_status`,
  ADD INDEX `idx_tch_audit_status` (`audit_status`),
  ADD INDEX `idx_tch_invigilation_experience` (`invigilation_experience`);

-- 依据已有的实际监考分配记录回填历史经验。
UPDATE `ea_tch_info` AS `t`
SET `t`.`invigilation_experience` = 1
WHERE EXISTS (
  SELECT 1
  FROM `ea_arrange_tch` AS `a`
  WHERE `a`.`tch_id` = `t`.`tch_id`
    AND `a`.`deleted` = 1
);

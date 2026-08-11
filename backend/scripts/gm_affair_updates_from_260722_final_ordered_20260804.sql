-- 国美考务数据库最终更新脚本（260722含及以后）
-- 生成日期：2026-08-04
--
-- 本文件不是机械拼接，已按表和字段依赖重新排序：
-- 1. 教职工查询字段 -> 学院流程字段 -> 最终审核数据修复。
-- 2. 专家管理字段 -> 专家查询索引。
-- 3. 禁止抽取专家最终建表 -> 兼容旧结构补列/回填。
-- 4. 评委征集任务及要求表 -> 要求表最终软删结构 -> 抽取管理。
-- 5. 学院专家上报使用最终 bigint deleted 表结构，删除被后续方案替代的旧表迁移段。
--
-- 适用基线：尚未执行 update_260722 及其后续脚本的数据库。
-- 执行前请备份数据库。


-- ============================================================================
-- SOURCE: update_260722_teacher_query_filters.sql
-- ============================================================================
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



-- ============================================================================
-- SOURCE: update_260722_teacher_department_workflow.sql
-- ============================================================================
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



-- ============================================================================
-- SOURCE: update_260723_teacher_audit_workflow.sql
-- ============================================================================
-- 2026/07/23 教职工管理升级：招办审核流程、启用禁用按钮、监考经验回算
-- 适用已经存在 ea_tch_info、ea_tch_info_audit_history 表的库，未初始化请使用 update_240910.sql

-- 1.确保 ea_tch_info 字段存在（升级库可能缺失列）
ALTER TABLE `ea_tch_info`
  MODIFY COLUMN `audit_status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '审核状态 1-未提交 2-待审核 3-审核通过 4-审核不通过',
  MODIFY COLUMN `enabled_status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '启用状态 1-启用 2-禁用',
  MODIFY COLUMN `invigilation_experience` tinyint(1) NOT NULL DEFAULT 2 COMMENT '监考经验 1-有 2-无';

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



-- ============================================================================
-- SOURCE: update_260723.sql
-- ============================================================================
-- 迁移脚本：ea_tch_task_report.state 字段枚举值变更
-- 旧值:  1-已上报 2-已征集/已分配
-- 新值:  1-已上报,待审核 2-待分配 3-不通过 4-已分配
-- 影响: 原 state='2'（已征集/已分配）→ 改为 '4'（已分配）
-- 日期: 2026-07-23
-- 注意: '2' 的新含义为"招办审核通过,待分配考场"，因此旧 '2' 必须先迁出

UPDATE `ea_tch_task_report` SET `state` = '4' WHERE `state` = '2';


-- ============================================================================
-- SOURCE: update_260727_expert_management.sql
-- ============================================================================
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



-- ============================================================================
-- SOURCE: update_260727_college_expert_management.sql
-- ============================================================================
-- 20260727 学院角色专家库管理调整

SET @schema_name = DATABASE();

SET @sql = IF(
  EXISTS (
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = @schema_name
      AND table_name = 'ea_judge_experts'
      AND index_name = 'idx_expert_id_card'
  ),
  'SELECT 1',
  'ALTER TABLE `ea_judge_experts` ADD INDEX `idx_expert_id_card` (`id_card`)'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql = IF(
  EXISTS (
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = @schema_name
      AND table_name = 'ea_judge_experts'
      AND index_name = 'idx_expert_dept_audit'
  ),
  'SELECT 1',
  'ALTER TABLE `ea_judge_experts` ADD INDEX `idx_expert_dept_audit` (`dept_id`, `audit_status`, `deleted`)'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;



-- ============================================================================
-- SOURCE: update_260728_expert_ban_management.sql
-- ============================================================================
-- 20260728 禁止抽取专家管理：建表（最终结构）
-- 说明：新增 id_card（证件号）与 deleted（1-正常 2-逻辑删除）字段；
--       不再对 expert_id 建唯一索引，因为逻辑删除后需允许对同一专家重新禁止，
--       唯一性改由业务层（deleted=1 计数）保证。
CREATE TABLE IF NOT EXISTS `ea_judge_expert_ban` (
  `ban_id` bigint NOT NULL AUTO_INCREMENT COMMENT '禁止抽取记录ID',
  `expert_id` bigint NOT NULL COMMENT '专家ID',
  `id_card` varchar(32) NOT NULL DEFAULT '' COMMENT '证件号',
  `create_by` bigint NULL COMMENT '创建人',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `deleted` int NOT NULL DEFAULT 1 COMMENT '1-正常 2-逻辑删除',
  PRIMARY KEY (`ban_id`)
) ENGINE=InnoDB COMMENT='禁止抽取专家名单';

-- 专家表查询索引（allow_draw + id_card），幂等创建
SET @index_exists = (
  SELECT COUNT(1)
  FROM information_schema.statistics
  WHERE table_schema = DATABASE()
    AND table_name = 'ea_judge_experts'
    AND index_name = 'idx_judge_experts_allow_draw_id_card'
);
SET @ddl = IF(
  @index_exists = 0,
  'CREATE INDEX idx_judge_experts_allow_draw_id_card ON ea_judge_experts (allow_draw, id_card)',
  'SELECT 1'
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ============================================================================
-- SOURCE: update_260729_expert_ban_idcard_logicdelete.sql
-- ============================================================================
-- 20260729 禁止抽取专家：新增证件号(id_card)与逻辑删除(deleted)字段，并移除 expert_id 唯一索引
-- 适用场景：已通过旧版 update_260728 建表（无 id_card/deleted 且带 uk_judge_expert_ban_expert 唯一索引）的数据库。
-- 本脚本幂等，可重复执行。

-- 1) 补齐 id_card 列
SET @col_idcard = (
  SELECT COUNT(1)
  FROM information_schema.columns
  WHERE table_schema = DATABASE()
    AND table_name = 'ea_judge_expert_ban'
    AND column_name = 'id_card'
);
SET @ddl = IF(
  @col_idcard = 0,
  'ALTER TABLE ea_judge_expert_ban ADD COLUMN id_card varchar(32) NOT NULL DEFAULT '''' COMMENT ''证件号''',
  'SELECT 1'
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =====================================================================

-- 2) 补齐 deleted 列
SET @col_deleted = (
  SELECT COUNT(1)
  FROM information_schema.columns
  WHERE table_schema = DATABASE()
    AND table_name = 'ea_judge_expert_ban'
    AND column_name = 'deleted'
);
SET @ddl = IF(
  @col_deleted = 0,
  'ALTER TABLE ea_judge_expert_ban ADD COLUMN deleted int NOT NULL DEFAULT 1 COMMENT ''删除标记：1-正常的；2-已删除的；''',
  'SELECT 1'
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =====================================================================

-- 3) 移除 expert_id 唯一索引（逻辑删除后需允许重新禁止同一专家，唯一性改由业务层保证）
SET @idx = (
  SELECT COUNT(1)
  FROM information_schema.statistics
  WHERE table_schema = DATABASE()
    AND table_name = 'ea_judge_expert_ban'
    AND index_name = 'uk_judge_expert_ban_expert'
);
SET @ddl = IF(
  @idx > 0,
  'ALTER TABLE ea_judge_expert_ban DROP INDEX uk_judge_expert_ban_expert',
  'SELECT 1'
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =====================================================================

-- 4) 回填已有禁止记录的证件号（来自专家表）
UPDATE ea_judge_expert_ban b
  JOIN ea_judge_experts e ON e.expert_id = b.expert_id
SET b.id_card = e.id_card
WHERE b.id_card = '';

-- 5) 存量记录标记为正常（deleted=1），保证历史数据可见
UPDATE ea_judge_expert_ban SET deleted = 1 WHERE deleted = 0 OR deleted IS NULL;



-- ============================================================================
-- SOURCE: update_260729_expert_ban_update_cols.sql
-- ============================================================================
-- 20260729 禁止抽取专家：新增更新者(update_by)与更新时间(update_time)字段
-- 适用场景：取消禁止抽取（逻辑删除）时记录是谁、何时操作的；新建/导入禁止记录时同样写入。
-- 本脚本幂等，可重复执行。

-- 1) 补齐 update_by 列
SET @col_update_by = (
  SELECT COUNT(1)
  FROM information_schema.columns
  WHERE table_schema = DATABASE()
    AND table_name = 'ea_judge_expert_ban'
    AND column_name = 'update_by'
);
SET @ddl = IF(
  @col_update_by = 0,
  'ALTER TABLE ea_judge_expert_ban ADD COLUMN update_by bigint NOT NULL DEFAULT 0 COMMENT ''更新者（谁取消禁止抽取/最后更新）''',
  'SELECT 1'
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =====================================================================

-- 2) 补齐 update_time 列
SET @col_update_time = (
  SELECT COUNT(1)
  FROM information_schema.columns
  WHERE table_schema = DATABASE()
    AND table_name = 'ea_judge_expert_ban'
    AND column_name = 'update_time'
);
SET @ddl = IF(
  @col_update_time = 0,
  'ALTER TABLE ea_judge_expert_ban ADD COLUMN update_time datetime NULL DEFAULT NULL COMMENT ''更新时间（取消禁止抽取/最后更新时间）''',
  'SELECT 1'
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 3) 回填历史记录：将尚未填写的记录的更新者/更新时间初始化为创建者/创建时间
UPDATE ea_judge_expert_ban
SET update_by   = create_by,
    update_time = create_time
WHERE update_by = 0 OR update_time IS NULL;



-- ============================================================================
-- SOURCE: update_260728_judge_recruit_task.sql
-- ============================================================================
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



-- ============================================================================
-- SOURCE: update_260803_judge_task_report_requirement_deleted.sql
-- ============================================================================
-- 学院评委专家上报任务方向要求表：软删字段改为 Unix 时间戳，解决唯一键冲突。
-- 与 ea_judge_expert_report 的修复策略保持一致：
--   未删除记录 deleted = 1（固定值，保证 uk_task_direction_dept 唯一键）；
--   已删除记录 deleted = UNIX_TIMESTAMP()（不同记录时间戳不同，避免重复）。

-- 1. 先修改字段类型为 bigint，否则后续 UPDATE 时间戳会超 tinyint 范围。
ALTER TABLE ea_judge_task_report_requirement MODIFY COLUMN deleted bigint NOT NULL DEFAULT 1 COMMENT '删除标记：1-正常的；大于1-已删除的（Unix时间戳）';

-- 2. 将已软删的旧数据（deleted=2）转换为唯一时间戳，避免唯一键冲突。
UPDATE ea_judge_task_report_requirement
SET deleted = UNIX_TIMESTAMP() + requirement_id
WHERE deleted = 2;

-- 3. 幂等重建唯一键（先删除若存在，再创建）。
SET @drop_sql = IF(
  EXISTS (
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'ea_judge_task_report_requirement'
      AND index_name = 'uk_task_direction_dept'
  ),
  'ALTER TABLE ea_judge_task_report_requirement DROP INDEX uk_task_direction_dept',
  'SELECT 1'
);
PREPARE drop_stmt FROM @drop_sql;
EXECUTE drop_stmt;
DEALLOCATE PREPARE drop_stmt;

SET @add_sql = IF(
  NOT EXISTS (
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'ea_judge_task_report_requirement'
      AND index_name = 'uk_task_direction_dept'
  ),
  'ALTER TABLE ea_judge_task_report_requirement ADD UNIQUE INDEX uk_task_direction_dept (judge_task_id, direction_code, dept_name, deleted)',
  'SELECT 1'
);
PREPARE add_stmt FROM @add_sql;
EXECUTE add_stmt;
DEALLOCATE PREPARE add_stmt;

-- ============================================================================
-- SOURCE: update_260729_judge_draw_management.sql
-- ============================================================================
-- 评委抽取管理：轮次科目配置与暂存专家
ALTER TABLE `ea_judge_draw_records`
  ADD COLUMN `required_num` int NOT NULL DEFAULT 0 COMMENT '评分专家计划数量' AFTER `subject_name`,
  ADD COLUMN `submit_state` tinyint NOT NULL DEFAULT 1 COMMENT '1-未提交 2-已提交' AFTER `required_num`;

CREATE TABLE IF NOT EXISTS `ea_judge_draw_selection` (
  `selection_id` bigint NOT NULL AUTO_INCREMENT COMMENT '选择记录ID',
  `rules_id` bigint NOT NULL COMMENT '抽取规则ID',
  `judge_task_id` bigint NOT NULL COMMENT '评委任务ID',
  `draw_times` int NOT NULL COMMENT '抽取次数',
  `expert_id` bigint NOT NULL COMMENT '专家ID',
  `expert_role` tinyint NOT NULL COMMENT '1-评分专家 2-备用专家',
  `submit_state` tinyint NOT NULL DEFAULT 1 COMMENT '1-未提交 2-已提交',
  `deleted` tinyint NOT NULL DEFAULT 1 COMMENT '删除标记：1-正常的；2-已删除的；',
  `create_by` bigint DEFAULT NULL COMMENT '创建人',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新人',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`selection_id`),
  UNIQUE KEY `uk_draw_selection_rule_expert` (`rules_id`,`expert_id`),
  KEY `idx_draw_selection_round` (`judge_task_id`,`draw_times`,`expert_id`,`deleted`)
) ENGINE=InnoDB COMMENT='评委抽取管理专家选择记录';



-- ============================================================================
-- SOURCE: update_260730_college_judge_report.sql
-- ============================================================================
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
  deleted bigint NOT NULL DEFAULT 1 COMMENT '是否删除：1-未删除，>1-已删除(时间戳)',
  create_by bigint DEFAULT NULL COMMENT '创建人ID',
  create_time datetime DEFAULT NULL COMMENT '创建时间',
  update_by bigint DEFAULT NULL COMMENT '更新人ID',
  update_time datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (report_id),
  UNIQUE KEY uk_active (requirement_id, expert_id, deleted) COMMENT '正常记录(deleted=1)唯一约束；已删记录时间戳不冲突',
  KEY idx_task_dept_state (judge_task_id, dept_id, submit_state, audit_state, deleted),
  KEY idx_requirement_state (requirement_id, submit_state, audit_state, deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学院评委专家上报记录表';



-- ============================================================================
-- SOURCE: update_260731_judge_expert_allow_draw_comment.sql
-- ============================================================================
-- 20260731 调整 ea_judge_experts.allow_draw 字段说明
-- 语义由"是否允许抽取"调整为"是否启用"，即专家库的启用/禁用开关。
-- 招办在专家库列表对专家做启用/禁用操作时即更新该字段；
-- 学院上报与招办抽取时，allow_draw=2（禁用）的专家不可被选择（复用于启用禁用控制，不再新增字段）。
-- 该 ALTER 仅修改列注释（元数据），可重复执行，无数据影响。

ALTER TABLE `ea_judge_experts`
  MODIFY COLUMN `allow_draw` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否启用 1-启用 2-禁用';



-- ============================================================================
-- SOURCE: update_260806_expert_intro_text.sql
-- ============================================================================
-- 20260806 专家介绍(intro)字段扩容
-- 问题：intro 原为 varchar(512)，utf8mb4 下最多 512 字符；前端允许输入 1000 字，
--       超长（>512 字符）写入时 STRICT_TRANS_TABLES 报 "Data too long for column 'intro'"，
--       被事务 defer-recover 静默吞掉，表现为保存成功但库里未更新。
-- 解决：改为 text（utf8mb4 下约可存 21845 字，远超前端 1000 字上限）。
-- 适用：线上已存在 ea_judge_experts 表的库，仅修改列类型，幂等可重复执行。

-- 1. 幂等修改 intro 列类型为 text（若已为 text 则等价 MODIFY，无影响）。
SET @col_intro_type = (
  SELECT DATA_TYPE
  FROM information_schema.columns
  WHERE table_schema = DATABASE()
    AND table_name = 'ea_judge_experts'
    AND column_name = 'intro'
);
SET @ddl = IF(
  @col_intro_type IS NULL OR @col_intro_type <> 'text',
  'ALTER TABLE ea_judge_experts MODIFY COLUMN intro text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT ''专家介绍''',
  'SELECT 1'
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =====================================================================

-- =====================================================================

-- =====================================================================
-- 评委抽取专家选择(removeDrawSelection) 删除改造：deleted 改为 bigint 时间戳软删
-- 目标：招办在评委抽取管理中移除专家时，不再 UPDATE 同一条记录的 deleted=2（覆盖原数据，无法回溯），
--       而是把原记录的 deleted 写为 Unix 时间戳作为删除标记（保留历史），后续重新抽取同一专家时
--       因唯一键含 deleted 会自然 INSERT 一条 deleted=1 的新记录（新增一条数据），避免主键/唯一键冲突。
-- 时间相关内容（对齐上方 ea_judge_expert_report 模板）：
--   * deleted：bigint，1-未删除，>1-已删除(记录删除时刻的 Unix 时间戳)，用于回溯删除时间；
--   * update_time：移除/删除时更新为 NOW()，记录删除操作发生的时间；
--   * replaced_from：被替换的原选择记录ID（重抽新增时记录来源），便于回溯“删除->重新新增”的关联。
-- 参考：学院角色专家上报 ea_judge_expert_report / ea_judge_task_report_requirement 的 deleted bigint 软删模式。
-- =====================================================================

-- 1. 修改字段类型为 bigint，否则后续 UPDATE 时间戳会超 tinyint 范围。
--    与上方 ea_judge_expert_report.deleted 保持一致：>1 即记录删除时刻的 Unix 时间戳，便于回溯删除时间。
ALTER TABLE ea_judge_draw_selection MODIFY COLUMN deleted bigint NOT NULL DEFAULT 1 COMMENT '是否删除：1-未删除，>1-已删除(时间戳)';

-- 1.1 新增 replaced_from（对齐上方 ea_judge_expert_report.replaced_from），记录被替换/删除后重新新增的来源记录ID，便于回溯。
SET @col_repfrom = (
  SELECT DATA_TYPE
  FROM information_schema.columns
  WHERE table_schema = DATABASE()
    AND table_name = 'ea_judge_draw_selection'
    AND column_name = 'replaced_from'
);
SET @ddl = IF(
  @col_repfrom IS NULL,
  "ALTER TABLE ea_judge_draw_selection ADD COLUMN replaced_from bigint DEFAULT NULL COMMENT '被替换的原选择记录ID（重抽新增时记录来源，便于回溯）'",
  'SELECT 1'
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 2. 将已软删的旧数据（deleted=2）转换为唯一时间戳，并回填 update_time 记录删除时间，避免唯一键冲突。
UPDATE ea_judge_draw_selection
SET deleted = UNIX_TIMESTAMP() + selection_id,
    update_time = NOW()
WHERE deleted = 2;

-- 3. 幂等重建唯一键（先删除若存在，再创建）。
SET @drop_sql = IF(
  EXISTS (
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'ea_judge_draw_selection'
      AND index_name = 'uk_draw_selection_rule_expert'
  ),
  'ALTER TABLE ea_judge_draw_selection DROP INDEX uk_draw_selection_rule_expert',
  'SELECT 1'
);
PREPARE drop_stmt FROM @drop_sql;
EXECUTE drop_stmt;
DEALLOCATE PREPARE drop_stmt;

SET @add_sql = IF(
  NOT EXISTS (
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'ea_judge_draw_selection'
      AND index_name = 'uk_draw_selection_rule_expert'
  ),
  'ALTER TABLE ea_judge_draw_selection ADD UNIQUE INDEX uk_draw_selection_rule_expert (rules_id, expert_id, deleted)',
  'SELECT 1'
);
PREPARE add_stmt FROM @add_sql;
EXECUTE add_stmt;
DEALLOCATE PREPARE add_stmt;

-- =====================================================================
-- 20260811 监考老师取消上报替换链路改造
-- 问题：原替换逻辑直接修改 ea_tch_task_report 原记录中的老师信息，导致原上报人、原审核状态等历史被覆盖，无法追溯。
-- 方案：原记录软删除，替换老师新增一条有效记录；新记录通过 replaced_from 指向被替换的原 report_id。
-- 说明：当前有效记录 deleted=1，历史记录 deleted=2；删除操作人和时间记录在 update_by、update_time。
-- =====================================================================

SET @col_tch_report_replaced_from = (
  SELECT DATA_TYPE
  FROM information_schema.columns
  WHERE table_schema = DATABASE()
    AND table_name = 'ea_tch_task_report'
    AND column_name = 'replaced_from'
);
SET @ddl = IF(
  @col_tch_report_replaced_from IS NULL,
  "ALTER TABLE ea_tch_task_report ADD COLUMN replaced_from bigint DEFAULT NULL COMMENT '被替换的原上报记录ID' AFTER state",
  'SELECT 1'
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @idx_tch_report_replaced_from = (
  SELECT COUNT(1)
  FROM information_schema.statistics
  WHERE table_schema = DATABASE()
    AND table_name = 'ea_tch_task_report'
    AND index_name = 'idx_tch_report_replaced_from'
);
SET @ddl = IF(
  @idx_tch_report_replaced_from = 0,
  'ALTER TABLE ea_tch_task_report ADD INDEX idx_tch_report_replaced_from (replaced_from)',
  'SELECT 1'
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

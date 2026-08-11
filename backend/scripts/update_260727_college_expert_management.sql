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

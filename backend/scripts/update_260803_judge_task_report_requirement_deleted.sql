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

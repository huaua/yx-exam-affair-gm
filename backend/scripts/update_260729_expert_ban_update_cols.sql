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

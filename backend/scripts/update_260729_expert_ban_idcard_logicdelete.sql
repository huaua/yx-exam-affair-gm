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

-- 4) 回填已有禁止记录的证件号（来自专家表）
UPDATE ea_judge_expert_ban b
  JOIN ea_judge_experts e ON e.expert_id = b.expert_id
SET b.id_card = e.id_card
WHERE b.id_card = '';

-- 5) 存量记录标记为正常（deleted=1），保证历史数据可见
UPDATE ea_judge_expert_ban SET deleted = 1 WHERE deleted = 0 OR deleted IS NULL;

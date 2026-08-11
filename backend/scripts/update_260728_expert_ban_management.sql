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

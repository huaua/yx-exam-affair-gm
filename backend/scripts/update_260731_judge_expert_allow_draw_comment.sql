-- 20260731 调整 ea_judge_experts.allow_draw 字段说明
-- 语义由“是否允许抽取”调整为“是否启用”，即专家库的启用/禁用开关。
-- 招办在专家库列表对专家做启用/禁用操作时即更新该字段；
-- 学院上报与招办抽取时，allow_draw=2（禁用）的专家不可被选择（复用于启用禁用控制，不再新增字段）。
-- 该 ALTER 仅修改列注释（元数据），可重复执行，无数据影响。

ALTER TABLE `ea_judge_experts`
  MODIFY COLUMN `allow_draw` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否启用 1-启用 2-禁用';

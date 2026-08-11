-- 考场上报记录表，冗余任务日期、删除考场状态字段
ALTER TABLE `ea_room_task_report` 
DROP COLUMN `room_state`,
ADD COLUMN `task_day` datetime NOT NULL COMMENT '征集日期' AFTER `task_id`;

-- 数据同步处理
update ea_room_task_report a, ea_room_task_info b set a.task_day = b.task_day where a.task_id = b.task_id

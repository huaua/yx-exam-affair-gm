ALTER TABLE `ea_room_task_report`
    ADD COLUMN `group_num` int DEFAULT NULL COMMENT '组数' AFTER `room_name`,
    ADD COLUMN `group_capacity` int DEFAULT NULL COMMENT '按组容量' AFTER `group_num`,
    ADD COLUMN `capacity` int DEFAULT NULL COMMENT '按位容量' AFTER `group_capacity`;

CREATE TABLE `ea_arrange_prof_info` (
                                        `info_id` bigint NOT NULL AUTO_INCREMENT COMMENT '数据ID',
                                        `room_task_id` bigint NOT NULL COMMENT '考场任务ID',
                                        `task_day` datetime NOT NULL COMMENT '征集日期',
                                        `dept_name` varchar(32) NOT NULL COMMENT '学院名称',
                                        `prof_name` varchar(64) NOT NULL COMMENT '专业名称',
                                        `direction_name` varchar(64) DEFAULT NULL COMMENT '方向名称',
                                        `arrange_type` char(1) DEFAULT '1' COMMENT '安排类型: 1-按组 2-按位',
                                        `stu_num` int DEFAULT '0' COMMENT '考生数量',
                                        `arrange_stu_num` int DEFAULT '0' COMMENT '已编排考生数量',
                                        `arrange_room_num` int DEFAULT '0' COMMENT '已编排考场数量',
                                        `arrange_state` char(1) DEFAULT '1' COMMENT '考场编排状态: 1-待编排 2-已编排',
                                        `deleted` tinyint(1) DEFAULT '1' COMMENT '删除标记：1-正常的；2-已删除的；',
                                        `remark` varchar(500) DEFAULT NULL COMMENT '备注',
                                        `create_by` bigint DEFAULT NULL COMMENT '创建者',
                                        `create_time` datetime DEFAULT NULL COMMENT '创建时间',
                                        `update_by` bigint DEFAULT NULL COMMENT '更新者',
                                        `update_time` datetime DEFAULT NULL COMMENT '更新时间',
                                        PRIMARY KEY (`info_id`) USING BTREE
) ENGINE=InnoDB COMMENT='编排专业信息表';


CREATE TABLE `ea_arrange_room` (
                                   `room_arrange_id` bigint NOT NULL AUTO_INCREMENT COMMENT '考场安排ID',
                                   `info_id` bigint NOT NULL COMMENT '数据ID',
                                   `arrange_type` char(1) NOT NULL COMMENT '安排类型: 1-按组 2-按位',
                                   `room_report_id` bigint NOT NULL COMMENT '考场上报ID',
                                   `room_task_id` bigint NOT NULL COMMENT '任务ID',
                                   `room_id` bigint NOT NULL COMMENT '教室ID',
                                   `room_name` varchar(32) NOT NULL COMMENT '教室名称',
                                   `stu_num` int DEFAULT '0' COMMENT '考生数量',
                                   `room_capacity` int DEFAULT '0' COMMENT '考场实际容量',
                                   `arrange_no_str` varchar(8) NOT NULL COMMENT '考场安排编号(字符串)',
                                   `arrange_no` int NOT NULL COMMENT '考场安排编号(数值)',
                                   `need_tch_num` int NOT NULL COMMENT '需要监考老师数量',
                                   `assignment_state` char(1) NOT NULL DEFAULT '1' COMMENT '教师分配状态 1-未分配 2-已分配',
                                   `deleted` tinyint(1) DEFAULT '1' COMMENT '删除标记：1-正常的；2-已删除的；',
                                   `remark` varchar(500) DEFAULT NULL COMMENT '备注',
                                   `create_by` bigint DEFAULT NULL COMMENT '创建者',
                                   `create_time` datetime DEFAULT NULL COMMENT '创建时间',
                                   `update_by` bigint DEFAULT NULL COMMENT '更新者',
                                   `update_time` datetime DEFAULT NULL COMMENT '更新时间',
                                   PRIMARY KEY (`room_arrange_id`) USING BTREE,
                                   KEY `idx_roomtaskid` (`room_task_id`) USING BTREE
) ENGINE=InnoDB COMMENT='考场编排详情表';


CREATE TABLE `ea_arrange_tch` (
                                  `tch_arrange_id` bigint NOT NULL AUTO_INCREMENT COMMENT '老师安排ID',
                                  `room_report_id` bigint NOT NULL COMMENT '考场上报ID',
                                  `room_task_id` bigint NOT NULL COMMENT '考场任务ID',
                                  `tch_report_id` bigint NOT NULL COMMENT '教职工上报ID',
                                  `tch_task_id` bigint NOT NULL COMMENT '教师征集任务ID',
                                  `tch_id` bigint NOT NULL COMMENT '教职工ID',
                                  `tch_name` varchar(32) NOT NULL COMMENT '教职工姓名（冗余）',
                                  `deleted` tinyint(1) DEFAULT '1' COMMENT '删除标记：1-正常的；2-已删除的；',
                                  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
                                  `create_by` bigint DEFAULT NULL COMMENT '创建者',
                                  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
                                  `update_by` bigint DEFAULT NULL COMMENT '更新者',
                                  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
                                  PRIMARY KEY (`tch_arrange_id`) USING BTREE
) ENGINE=InnoDB COMMENT='监考老师分配表';
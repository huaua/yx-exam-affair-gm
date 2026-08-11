ALTER TABLE `ea_classroom_info` CHANGE COLUMN `room_type` `group_capacity` int NULL COMMENT '按组容量' AFTER `group_num`,
    MODIFY COLUMN `capacity` int NULL DEFAULT 0 COMMENT '按位容量' AFTER `group_capacity`;

ALTER TABLE `ea_task_room_info` DROP COLUMN `room_type`;

ALTER TABLE `ea_classroom_info` RENAME TO `ea_room_info`;
ALTER TABLE `ea_task_info` RENAME TO `ea_room_task_info`;
ALTER TABLE `ea_task_room_info` RENAME TO `ea_room_task_report`;

CREATE TABLE `ea_tch_info` (
                               `tch_id` bigint NOT NULL AUTO_INCREMENT COMMENT '教职工ID',
                               `tch_name` varchar(32) NOT NULL COMMENT '姓名',
                               `sex` char(2) NOT NULL COMMENT '性别（男或女）',
                               `job_no` varchar(32) NOT NULL COMMENT '工号',
                               `id_card` varchar(32) NOT NULL COMMENT '证件号',
                               `dept_id` bigint NOT NULL COMMENT '所属学院id',
                               `duties` varchar(32) DEFAULT NULL COMMENT '职务',
                               `bankcard` varchar(32) DEFAULT NULL COMMENT '工商银行卡号',
                               `prof_office` varchar(32) DEFAULT NULL COMMENT '专业或行政',
                               `education` varchar(1024)  DEFAULT NULL COMMENT '专业背景',
                               `audit_status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '审核状态 1-未提交 2-待审核 3-审核通过 4-审核不通过',
                               `audit_remark` varchar(500) DEFAULT NULL COMMENT '审核不通过原因',
                               `invigilation_experience` tinyint(1) NOT NULL DEFAULT '2' COMMENT '监考经验 1-有 2-无',
                               `enabled_status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '启用状态 1-启用 2-禁用',
                               `date_from` tinyint(1) NOT NULL DEFAULT '1' COMMENT '数据来源 1-学院新增 2-招办新增',
                               `deleted` tinyint(1) DEFAULT '1' COMMENT '删除标记：1-正常的；2-已删除的；',
                               `remark` varchar(500) DEFAULT NULL COMMENT '备注',
                               `create_by` bigint DEFAULT NULL COMMENT '创建者',
                               `create_time` datetime DEFAULT NULL COMMENT '创建时间',
                               `update_by` bigint DEFAULT NULL COMMENT '更新者',
                               `update_time` datetime DEFAULT NULL COMMENT '更新时间',
                               PRIMARY KEY (`tch_id`) USING BTREE,
                               KEY `idx_dept` (`dept_id`) USING BTREE,
                               KEY `idx_idcard` (`id_card`) USING BTREE,
                               KEY `idx_jobno` (`job_no`) USING BTREE,
                               KEY `idx_tch_audit_status` (`audit_status`) USING BTREE,
                               KEY `idx_tch_invigilation_experience` (`invigilation_experience`) USING BTREE,
                               KEY `idx_tch_enabled_status` (`enabled_status`) USING BTREE
) ENGINE=InnoDB COMMENT='教职工信息表';

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

CREATE TABLE `ea_tch_task_info` (
                                    `tch_task_id` bigint NOT NULL AUTO_INCREMENT COMMENT '教师征集任务ID',
                                    `tch_task_name` varchar(32) NOT NULL COMMENT '考试任务名称',
                                    `task_time_remark` varchar(32) NOT NULL COMMENT '征集日期说明',
                                    `task_start_time` datetime NOT NULL COMMENT '征集开始时间',
                                    `task_end_time` datetime NOT NULL COMMENT '征集结束时间',
                                    `need_num` int NOT NULL DEFAULT '0' COMMENT '总征集数量',
                                    `collect_num` int NOT NULL DEFAULT '0' COMMENT '已征集数量',
                                    `state` char(1) NOT NULL DEFAULT '1' COMMENT '发布状态: 1-未发布 2-已发布',
                                    `deleted` tinyint(1) DEFAULT '1' COMMENT '删除标记：1-正常的；2-已删除的；',
                                    `remark` varchar(500) DEFAULT NULL COMMENT '备注',
                                    `create_by` bigint DEFAULT NULL COMMENT '创建者',
                                    `create_time` datetime DEFAULT NULL COMMENT '创建时间',
                                    `update_by` bigint DEFAULT NULL COMMENT '更新者',
                                    `update_time` datetime DEFAULT NULL COMMENT '更新时间',
                                    PRIMARY KEY (`tch_task_id`) USING BTREE,
                                    KEY `idx_state` (`state`) USING BTREE
) ENGINE=InnoDB COMMENT='监考老师征集信息表';

CREATE TABLE `ea_tch_task_detail` (
                                      `task_detail_id` bigint NOT NULL AUTO_INCREMENT COMMENT '教师征集详情ID',
                                      `tch_task_id` bigint NOT NULL COMMENT '教师征集任务ID',
                                      `task_time` datetime NOT NULL COMMENT '征集日期',
                                      `dept_id` bigint NOT NULL COMMENT '学院id',
                                      `need_num` int DEFAULT '0' COMMENT '征集人数',
                                      `collect_num` int DEFAULT '0' COMMENT '已征集数量',
                                      `force_flag` char(1) DEFAULT '1' COMMENT '是否强制征集: 1-强制 2-不强制',
                                      `deleted` tinyint(1) DEFAULT '1' COMMENT '删除标记：1-正常的；2-已删除的；',
                                      `remark` varchar(500) DEFAULT NULL COMMENT '备注',
                                      `create_by` bigint DEFAULT NULL COMMENT '创建者',
                                      `create_time` datetime DEFAULT NULL COMMENT '创建时间',
                                      `update_by` bigint DEFAULT NULL COMMENT '更新者',
                                      `update_time` datetime DEFAULT NULL COMMENT '更新时间',
                                      PRIMARY KEY (`task_detail_id`) USING BTREE,
                                      KEY `idx_task_time` (`tch_task_id`,`task_time`) USING BTREE,
                                      KEY `idx_dept` (`dept_id`) USING BTREE
) ENGINE=InnoDB COMMENT='监考老师征集任务明细表';

CREATE TABLE `ea_tch_task_report` (
                                      `report_id` bigint NOT NULL AUTO_INCREMENT COMMENT '上报ID',
                                      `tch_task_id` bigint NOT NULL COMMENT '征集任务ID',
                                      `task_detail_id` bigint NOT NULL COMMENT '教师征集详情ID',
                                      `task_time` datetime NOT NULL COMMENT '征集日期',
                                      `tch_id` bigint NOT NULL COMMENT '教职工ID',
                                      `tch_name` varchar(32) NOT NULL COMMENT '姓名',
                                      `job_no` varchar(32) NOT NULL COMMENT '工号',
                                      `dept_id` bigint NOT NULL COMMENT '学院id',
                                      `state` char(1) DEFAULT '1' COMMENT '状态: 1-已上报 2-已征集',
                                      `deleted` tinyint(1) DEFAULT '1' COMMENT '删除标记：1-正常的；2-已删除的；',
                                      `remark` varchar(500) DEFAULT NULL COMMENT '备注',
                                      `create_by` bigint DEFAULT NULL COMMENT '创建者',
                                      `create_time` datetime DEFAULT NULL COMMENT '创建时间',
                                      `update_by` bigint DEFAULT NULL COMMENT '更新者',
                                      `update_time` datetime DEFAULT NULL COMMENT '更新时间',
                                      PRIMARY KEY (`report_id`) USING BTREE,
                                      KEY `idx_time` (`task_time`) USING BTREE,
                                      KEY `idx_dept` (`dept_id`) USING BTREE
) ENGINE=InnoDB COMMENT='老师任务分配表';

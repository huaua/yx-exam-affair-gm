/*
 Navicat Premium Data Transfer

 Source Server         : 华为云
 Source Server Type    : MySQL
 Source Server Version : 80025
 Source Host           : 192.168.10.11:33106
 Source Schema         : yx_exam_affair_gm

 Target Server Type    : MySQL
 Target Server Version : 80025
 File Encoding         : 65001

 Date: 14/08/2024 09:36:18
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for ea_classroom_info
-- ----------------------------
DROP TABLE IF EXISTS `ea_classroom_info`;
CREATE TABLE `ea_classroom_info`  (
                                      `room_id` bigint NOT NULL AUTO_INCREMENT COMMENT '教室ID',
                                      `room_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '教室名称',
                                      `campus_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '所在校区',
                                      `level_code` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '所属层次',
                                      `room_type` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '教室类型',
                                      `dept_id` bigint NOT NULL COMMENT '所属学院id',
                                      `group_num` int NULL DEFAULT 0 COMMENT '组数',
                                      `capacity` int NULL DEFAULT 0 COMMENT '容量',
                                      `state` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '1' COMMENT '启用状态: 1-正常 2-停用',
                                      `deleted` tinyint(1) NULL DEFAULT 1 COMMENT '删除标记：1-正常的；2-已删除的；',
                                      `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注',
                                      `create_by` bigint NULL DEFAULT NULL COMMENT '创建者',
                                      `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
                                      `update_by` bigint NULL DEFAULT NULL COMMENT '更新者',
                                      `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
                                      PRIMARY KEY (`room_id`) USING BTREE,
                                      INDEX `idx_dept_leve`(`dept_id` ASC, `level_code` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '教室信息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ea_classroom_info
-- ----------------------------

-- ----------------------------
-- Table structure for ea_dept_info
-- ----------------------------
DROP TABLE IF EXISTS `ea_dept_info`;
CREATE TABLE `ea_dept_info`  (
                                 `dept_id` bigint NOT NULL AUTO_INCREMENT COMMENT '学院id',
                                 `dept_code` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '学院代码',
                                 `dept_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '学院名称',
                                 `campus_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '所在校区',
                                 `deleted` tinyint(1) NOT NULL COMMENT '是否删除：1-未删除；2-已删除；',
                                 `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注',
                                 `create_by` bigint NOT NULL COMMENT '创建人ID',
                                 `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                 `update_by` bigint NOT NULL COMMENT '更新人ID',
                                 `update_time` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                 PRIMARY KEY (`dept_id`) USING BTREE,
                                 UNIQUE INDEX `uk_code`(`dept_code` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '学院基础信息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ea_dept_info
-- ----------------------------

-- ----------------------------
-- Table structure for ea_task_info
-- ----------------------------
DROP TABLE IF EXISTS `ea_task_info`;
CREATE TABLE `ea_task_info`  (
                                 `task_id` bigint NOT NULL AUTO_INCREMENT COMMENT '任务ID',
                                 `task_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '任务名称',
                                 `task_day` datetime NOT NULL COMMENT '征集日期',
                                 `level_code` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '所属层次',
                                 `need_num` int NULL DEFAULT 0 COMMENT '需要考场数量',
                                 `collect_num` int NULL DEFAULT 0 COMMENT '已征集数量',
                                 `comfirm_num` int NULL DEFAULT 0 COMMENT '已确认数量',
                                 `state` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '1' COMMENT '发布状态: 1-未发布 2-已发布',
                                 `deleted` tinyint(1) NULL DEFAULT 1 COMMENT '删除标记：1-正常的；2-已删除的；',
                                 `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注',
                                 `create_by` bigint NULL DEFAULT NULL COMMENT '创建者',
                                 `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
                                 `update_by` bigint NULL DEFAULT NULL COMMENT '更新者',
                                 `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
                                 PRIMARY KEY (`task_id`) USING BTREE,
                                 INDEX `idx_day`(`task_day` ASC) USING BTREE,
                                 INDEX `idx_leve`(`level_code` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '征集任务表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ea_task_info
-- ----------------------------

-- ----------------------------
-- Table structure for ea_task_room_info
-- ----------------------------
DROP TABLE IF EXISTS `ea_task_room_info`;
CREATE TABLE `ea_task_room_info`  (
                                      `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
                                      `task_id` bigint NOT NULL COMMENT '任务ID',
                                      `room_id` bigint NOT NULL COMMENT '教室ID',
                                      `room_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '教室名称',
                                      `tch_num` int NULL DEFAULT 0 COMMENT '监考老师数量',
                                      `state` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '1' COMMENT '确认状态: 1-未确认 2-已确认',
                                      `deleted` tinyint(1) NULL DEFAULT 1 COMMENT '删除标记：1-正常的；2-已删除的；',
                                      `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '备注',
                                      `create_by` bigint NULL DEFAULT NULL COMMENT '创建者',
                                      `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
                                      `update_by` bigint NULL DEFAULT NULL COMMENT '更新者',
                                      `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
                                      PRIMARY KEY (`id`) USING BTREE,
                                      INDEX `idx_task`(`task_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '任务教室征集关联表表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ea_task_room_info
-- ----------------------------

-- ----------------------------
-- Table structure for sys_config
-- ----------------------------
DROP TABLE IF EXISTS `sys_config`;
CREATE TABLE `sys_config`  (
                               `config_id` int NOT NULL AUTO_INCREMENT COMMENT '参数主键',
                               `config_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '参数名称',
                               `config_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '参数键名',
                               `config_value` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '参数键值',
                               `config_type` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'N' COMMENT '系统内置（Y是 N否）',
                               `create_by` bigint NOT NULL DEFAULT 1 COMMENT '创建者',
                               `deleted` tinyint(1) NULL DEFAULT NULL COMMENT '删除标记 1-正常 2-已删除',
                               `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
                               `update_by` bigint NULL DEFAULT NULL COMMENT '更新者',
                               `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
                               `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注',
                               PRIMARY KEY (`config_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '参数配置表' ROW_FORMAT = COMPACT;

-- ----------------------------
-- Records of sys_config
-- ----------------------------

-- ----------------------------
-- Table structure for sys_data_inout
-- ----------------------------
DROP TABLE IF EXISTS `sys_data_inout`;
CREATE TABLE `sys_data_inout`  (
                                   `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
                                   `oper_type` tinyint(1) NOT NULL COMMENT '操作类型 1-导入 2-更新 3-导出',
                                   `biz_type` tinyint NOT NULL COMMENT '业务类型 1-高考改革数据 2-非高考改革数据 3-中专数据 4-专升本数据 5-第二学位数据',
                                   `batch_no` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '操作批次号',
                                   `param` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '参数信息',
                                   `file_in` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '上传文件地址，多个压缩成zip文件',
                                   `file_in_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '上传文件名称',
                                   `file_out` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '下载文件地址，多个压缩zip文件',
                                   `oper_time` datetime NULL DEFAULT NULL COMMENT '操作时间',
                                   `state` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '1' COMMENT '状态 1-成功 2-失败',
                                   `remark` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注',
                                   `deleted` tinyint(1) NULL DEFAULT 1 COMMENT '是否删除：1-未删除；2-已删除；',
                                   `create_by` bigint NOT NULL COMMENT '创建人ID',
                                   `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
                                   `update_by` bigint NULL DEFAULT NULL COMMENT '更新人ID',
                                   `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
                                   PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '导入导出记录表' ROW_FORMAT = COMPACT;

-- ----------------------------
-- Records of sys_data_inout
-- ----------------------------

-- ----------------------------
-- Table structure for sys_dict_data
-- ----------------------------
DROP TABLE IF EXISTS `sys_dict_data`;
CREATE TABLE `sys_dict_data`  (
                                  `dict_code` bigint NOT NULL AUTO_INCREMENT COMMENT '字典编码',
                                  `dict_sort` int NULL DEFAULT 0 COMMENT '字典排序',
                                  `dict_label` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '字典标签',
                                  `dict_value` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '字典键值',
                                  `dict_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '字典类型',
                                  `css_class` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '样式属性（其他样式扩展）',
                                  `list_class` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '表格回显样式',
                                  `is_default` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'N' COMMENT '是否默认（Y是 N否）',
                                  `state` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '1' COMMENT '状态: 1-正常 2-停用',
                                  `deleted` tinyint(1) NULL DEFAULT 1 COMMENT '删除标记：1-正常的；2-已删除的；',
                                  `create_by` bigint NOT NULL COMMENT '创建者',
                                  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
                                  `update_by` bigint NULL DEFAULT NULL COMMENT '更新者',
                                  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
                                  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注',
                                  PRIMARY KEY (`dict_code`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '字典数据表' ROW_FORMAT = COMPACT;

-- ----------------------------
-- Records of sys_dict_data
-- ----------------------------

-- ----------------------------
-- Table structure for sys_dict_type
-- ----------------------------
DROP TABLE IF EXISTS `sys_dict_type`;
CREATE TABLE `sys_dict_type`  (
                                  `dict_id` bigint NOT NULL AUTO_INCREMENT COMMENT '字典主键',
                                  `dict_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '字典名称',
                                  `dict_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '字典类型',
                                  `state` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '1' COMMENT '状态: 1-正常 2-停用',
                                  `deleted` tinyint(1) NULL DEFAULT 1 COMMENT '删除标记：1-正常的；2-已删除的；',
                                  `create_by` bigint NOT NULL COMMENT '创建者',
                                  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
                                  `update_by` bigint NULL DEFAULT NULL COMMENT '更新者',
                                  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
                                  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注',
                                  PRIMARY KEY (`dict_id`) USING BTREE,
                                  UNIQUE INDEX `dict_type`(`dict_type` ASC, `deleted` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '字典类型表' ROW_FORMAT = COMPACT;

-- ----------------------------
-- Records of sys_dict_type
-- ----------------------------

-- ----------------------------
-- Table structure for sys_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_menu`;
CREATE TABLE `sys_menu`  (
                             `menu_id` bigint NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
                             `menu_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '菜单名称',
                             `parent_id` bigint NULL DEFAULT 0 COMMENT '父菜单ID',
                             `order_num` int NULL DEFAULT 0 COMMENT '显示顺序',
                             `path` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '路由地址',
                             `component` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '组件路径',
                             `query` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '路由参数',
                             `is_frame` tinyint(1) NULL DEFAULT 2 COMMENT '是否为外链（1-是 2-否）',
                             `is_cache` tinyint(1) NULL DEFAULT 1 COMMENT '是否缓存（1-缓存 2-不缓存）',
                             `menu_type` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '菜单类型（M目录 C菜单 F按钮）',
                             `visible` tinyint(1) NULL DEFAULT 1 COMMENT '菜单状态（1-显示 2-隐藏）',
                             `perms` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '权限标识',
                             `icon` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '#' COMMENT '菜单图标',
                             `state` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '1' COMMENT '菜单状态（1-正常 2-停用）',
                             `deleted` tinyint(1) NULL DEFAULT 1 COMMENT '删除标记：1-正常的；2-已删除的；',
                             `create_by` bigint NOT NULL COMMENT '创建者',
                             `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
                             `update_by` bigint NULL DEFAULT NULL COMMENT '更新者',
                             `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
                             `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '备注',
                             PRIMARY KEY (`menu_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '菜单权限表' ROW_FORMAT = COMPACT;

-- ----------------------------
-- Records of sys_menu
-- ----------------------------
INSERT INTO `sys_menu` VALUES (1, '系统管理', 0, 1, 'sys', NULL, '', 2, 1, 'M', 1, '', 'system', '1', 1, 1, '2023-05-08 15:19:01', 1, '2023-05-15 11:08:43', '系统管理目录');
INSERT INTO `sys_menu` VALUES (100, '用户管理', 1, 1, 'user', 'system/user/index', '', 2, 1, 'C', 1, 'sys:user:list', 'user', '1', 1, 1, '2023-05-08 15:19:01', NULL, NULL, '用户管理菜单');
INSERT INTO `sys_menu` VALUES (101, '角色管理', 1, 2, 'role', 'system/role/index', '', 2, 1, 'C', 1, 'sys:role:list', 'peoples', '1', 1, 1, '2023-05-08 15:19:01', NULL, NULL, '角色管理菜单');
INSERT INTO `sys_menu` VALUES (102, '菜单管理', 1, 3, 'menu', 'system/menu/index', '', 2, 1, 'C', 1, 'sys:menu:list', 'tree-table', '1', 1, 1, '2023-05-08 15:19:01', NULL, NULL, '菜单管理菜单');
INSERT INTO `sys_menu` VALUES (105, '字典管理', 1, 6, 'dict', 'system/dict/index', '', 2, 1, 'C', 1, 'sys:dict:list', 'dict', '1', 1, 1, '2023-05-24 22:37:29', NULL, NULL, '字典管理菜单');
INSERT INTO `sys_menu` VALUES (106, '参数设置', 1, 7, 'config', 'system/config/index', '', 2, 1, 'C', 1, 'sys:config:list', 'edit', '1', 1, 1, '2023-05-24 22:37:29', NULL, NULL, '参数设置菜单');
INSERT INTO `sys_menu` VALUES (1000, '用户查询', 100, 1, '', '', '', 2, 1, 'F', 1, 'sys:user:query', '#', '1', 1, 1, '2023-05-08 15:19:01', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1001, '用户新增', 100, 2, '', '', '', 2, 1, 'F', 1, 'sys:user:add', '#', '1', 1, 1, '2023-05-08 15:19:01', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1002, '用户修改', 100, 3, '', '', '', 2, 1, 'F', 1, 'sys:user:edit', '#', '1', 1, 1, '2023-05-08 15:19:01', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1003, '用户删除', 100, 4, '', '', '', 2, 1, 'F', 1, 'sys:user:remove', '#', '1', 1, 1, '2023-05-08 15:19:01', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1004, '用户导出', 100, 5, '', '', '', 2, 1, 'F', 1, 'sys:user:export', '#', '1', 1, 1, '2023-05-08 15:19:01', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1005, '用户导入', 100, 6, '', '', '', 2, 1, 'F', 1, 'sys:user:import', '#', '1', 1, 1, '2023-05-08 15:19:01', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1006, '重置密码', 100, 7, '', '', '', 2, 1, 'F', 1, 'sys:user:resetPwd', '#', '1', 1, 1, '2023-05-08 15:19:01', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1007, '角色查询', 101, 1, '', '', '', 2, 1, 'F', 1, 'sys:role:query', '#', '1', 1, 1, '2023-05-08 15:19:01', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1008, '角色新增', 101, 2, '', '', '', 2, 1, 'F', 1, 'sys:role:add', '#', '1', 1, 1, '2023-05-08 15:19:01', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1009, '角色修改', 101, 3, '', '', '', 2, 1, 'F', 1, 'sys:role:edit', '#', '1', 1, 1, '2023-05-08 15:19:01', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1010, '角色删除', 101, 4, '', '', '', 2, 1, 'F', 1, 'sys:role:remove', '#', '1', 1, 1, '2023-05-08 15:19:01', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1011, '角色导出', 101, 5, '', '', '', 2, 1, 'F', 1, 'sys:role:export', '#', '1', 1, 1, '2023-05-08 15:19:02', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1012, '菜单查询', 102, 1, '', '', '', 2, 1, 'F', 1, 'sys:menu:query', '#', '1', 1, 1, '2023-05-08 15:19:02', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1013, '菜单新增', 102, 2, '', '', '', 2, 1, 'F', 1, 'sys:menu:add', '#', '1', 1, 1, '2023-05-08 15:19:02', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1014, '菜单修改', 102, 3, '', '', '', 2, 1, 'F', 1, 'sys:menu:edit', '#', '1', 1, 1, '2023-05-08 15:19:02', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1015, '菜单删除', 102, 4, '', '', '', 2, 1, 'F', 1, 'sys:menu:remove', '#', '1', 1, 1, '2023-05-08 15:19:02', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1025, '字典查询', 105, 1, '#', '', '', 2, 1, 'F', 1, 'sys:dict:query', '#', '1', 1, 1, '2023-05-24 22:37:29', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1026, '字典新增', 105, 2, '#', '', '', 2, 1, 'F', 1, 'sys:dict:add', '#', '1', 1, 1, '2023-05-24 22:37:29', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1027, '字典修改', 105, 3, '#', '', '', 2, 1, 'F', 1, 'sys:dict:edit', '#', '1', 1, 1, '2023-05-24 22:37:29', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1028, '字典删除', 105, 4, '#', '', '', 2, 1, 'F', 1, 'sys:dict:remove', '#', '1', 1, 1, '2023-05-24 22:37:29', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1029, '字典导出', 105, 5, '#', '', '', 2, 1, 'F', 1, 'sys:dict:export', '#', '1', 1, 1, '2023-05-24 22:37:29', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1030, '参数查询', 106, 1, '#', '', '', 2, 1, 'F', 1, 'sys:config:query', '#', '1', 1, 1, '2023-05-24 22:37:29', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1031, '参数新增', 106, 2, '#', '', '', 2, 1, 'F', 1, 'sys:config:add', '#', '1', 1, 1, '2023-05-24 22:37:29', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1032, '参数修改', 106, 3, '#', '', '', 2, 1, 'F', 1, 'sys:config:edit', '#', '1', 1, 1, '2023-05-24 22:37:29', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1033, '参数删除', 106, 4, '#', '', '', 2, 1, 'F', 1, 'sys:config:remove', '#', '1', 1, 1, '2023-05-24 22:37:29', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1034, '参数导出', 106, 5, '#', '', '', 2, 1, 'F', 1, 'sys:config:export', '#', '1', 1, 1, '2023-05-24 22:37:29', NULL, NULL, '');
INSERT INTO `sys_menu` VALUES (1039, '系统工具', 0, 2, 'tool', '', '', 2, 1, 'M', 1, '', 'tool', '1', 1, 1, '2023-06-14 10:46:47', 1, '2023-06-14 10:46:47', '');
INSERT INTO `sys_menu` VALUES (1040, '生成代码', 1039, 1, 'gen', 'tool/gen/index', '', 2, 1, 'C', 1, 'tool:gen:list', 'code', '1', 1, 1, '2023-06-14 10:48:02', 1, '2023-06-14 10:48:02', '');

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role`  (
                             `role_id` bigint NOT NULL AUTO_INCREMENT COMMENT '角色ID',
                             `role_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色名称',
                             `role_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色权限字符串',
                             `role_sort` int NOT NULL COMMENT '显示顺序',
                             `data_scope` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '1' COMMENT '数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）',
                             `menu_check_strictly` tinyint(1) NULL DEFAULT 1 COMMENT '菜单树选择项是否关联显示',
                             `dept_check_strictly` tinyint(1) NULL DEFAULT 1 COMMENT '部门树选择项是否关联显示',
                             `state` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色状态（1正常 2停用）',
                             `deleted` tinyint(1) NULL DEFAULT 1 COMMENT '删除标志（1代表存在 2代表删除）',
                             `create_by` bigint NOT NULL COMMENT '创建者',
                             `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
                             `update_by` bigint NULL DEFAULT NULL COMMENT '更新者',
                             `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
                             `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注',
                             PRIMARY KEY (`role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '角色信息表' ROW_FORMAT = COMPACT;

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO `sys_role` VALUES (1, '超级管理员', 'admin', 1, '1', 1, 1, '1', 1, 1, '2023-05-24 22:37:38', NULL, NULL, '超级管理员');
INSERT INTO `sys_role` VALUES (2, '招办管理角色', 'admissions', 2, '2', 1, 1, '1', 1, 1, '2023-05-24 22:37:38', 1, '2023-06-14 19:35:50', '招办管理员');
INSERT INTO `sys_role` VALUES (3, '学院管理角色', 'department', 3, '3', 1, 1, '1', 1, 1, '2023-05-24 22:37:38', 1, '2023-06-14 19:35:50', '学院管理员');

-- ----------------------------
-- Table structure for sys_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_menu`;
CREATE TABLE `sys_role_menu`  (
                                  `role_id` bigint NOT NULL COMMENT '角色ID',
                                  `menu_id` bigint NOT NULL COMMENT '菜单ID',
                                  PRIMARY KEY (`role_id`, `menu_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '角色和菜单关联表' ROW_FORMAT = COMPACT;

-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user`  (
                             `user_id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户ID',
                             `dept_id` bigint NULL DEFAULT NULL COMMENT '部门ID',
                             `user_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户账号',
                             `nick_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '用户昵称',
                             `user_type` varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '00' COMMENT '用户类型（00系统用户）',
                             `email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '用户邮箱',
                             `phonenumber` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '手机号码',
                             `sex` varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '用户性别: 男  女  未知',
                             `avatar` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '头像地址',
                             `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '密码',
                             `state` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '1' COMMENT '帐号状态: 1-正常 2-停用',
                             `deleted` tinyint(1) NULL DEFAULT 1 COMMENT '删除标记：1-正常的；2-已删除的；',
                             `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注',
                             `create_by` bigint NOT NULL COMMENT '创建者',
                             `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
                             `update_by` bigint NULL DEFAULT NULL COMMENT '更新者',
                             `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
                             PRIMARY KEY (`user_id`) USING BTREE,
                             UNIQUE INDEX `user_name`(`user_name` ASC, `deleted` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户信息表' ROW_FORMAT = COMPACT;

-- ----------------------------
-- Records of sys_user pwd:Csk002
-- ----------------------------
INSERT INTO `sys_user` VALUES (1, NULL, 'admin', '超级管理员', '00', 'admin@admin.com', '18888888888', '男', '', '$2a$14$RSvRR5h7BOp0It6D0OjgqO6xdcOHCUBe/owWDpO4/pwsZaNHTh9HO', '1', 1, '系统管理员', 1, '2021-08-15 12:02:05', 1, '2024-05-20 09:46:07');
INSERT INTO `sys_user` VALUES (2, NULL, 'zhaoban', '招办管理员', '00', '', '', '1', '', '$2a$14$RSvRR5h7BOp0It6D0OjgqO6xdcOHCUBe/owWDpO4/pwsZaNHTh9HO', '1', 1, '招生办管理员', 1, '2021-08-15 12:02:05', 1, '2023-06-14 18:03:43');
INSERT INTO `sys_user` VALUES (3, 1, 'xueyuan', '学院管理员', '00', '', '', '1', '', '$2a$14$RSvRR5h7BOp0It6D0OjgqO6xdcOHCUBe/owWDpO4/pwsZaNHTh9HO', '1', 1, '学院管理员', 1, '2021-08-15 12:02:05', 1, '2023-06-14 18:03:43');

-- ----------------------------
-- Table structure for sys_user_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_user_role`;
CREATE TABLE `sys_user_role`  (
                                  `user_id` bigint NOT NULL COMMENT '用户ID',
                                  `role_id` bigint NOT NULL COMMENT '角色ID',
                                  PRIMARY KEY (`user_id`, `role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户和角色关联表' ROW_FORMAT = COMPACT;

-- ----------------------------
-- Records of sys_user_role
-- ----------------------------
INSERT INTO `sys_user_role` (`user_id`, `role_id`) VALUES (2, 2);
INSERT INTO `sys_user_role` (`user_id`, `role_id`) VALUES (3, 3);


SET FOREIGN_KEY_CHECKS = 1;

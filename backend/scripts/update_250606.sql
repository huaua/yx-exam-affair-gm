-- 20250606 新增评委管理模块 评委专家库表 by liyu start
DROP TABLE IF EXISTS `ea_judge_experts`;
CREATE TABLE `ea_judge_experts`  (
  `expert_id` bigint NOT NULL AUTO_INCREMENT COMMENT '专家id',
  `name` varchar(32) NOT NULL COMMENT '姓名',
  `sex` char(2) NOT NULL COMMENT '性别（男或女）',
  `phone_no` varchar(32) NOT NULL COMMENT '手机号',
  `id_card` varchar(32) NOT NULL COMMENT '证件号',
  `dept_id` bigint NOT NULL COMMENT '专家来源-所属学院id（招办的属校外）',
  `type` tinyint(1) NULL DEFAULT NULL COMMENT '专家类型 字典(1-博士 2-硕士 3-本科 4-附中)',
  `init_edu` tinyint(1) NULL DEFAULT NULL COMMENT '初始学历 字典(1-本科 2-硕士 3-博士)',
  `final_edu` tinyint(1) NULL DEFAULT NULL COMMENT '最终学历 字典(1-本科 2-硕士 3-博士)',
  `good_subjects` varchar(256) NULL DEFAULT NULL COMMENT '擅长科目 字典(1-素描 2-速写 3-色彩) 多选拼接值',
  `bank_card` varchar(32) NULL DEFAULT NULL COMMENT '银行卡号',
  `intro` varchar(512) NULL DEFAULT NULL COMMENT '专家介绍',
  `allow_draw` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否允许抽取 1-允许 2-禁止',
  `data_from` tinyint(1) NOT NULL DEFAULT 1 COMMENT '数据来源 1-学院新增 2-招办新增',
  `deleted` tinyint(1) NULL DEFAULT 1 COMMENT '删除标记：1-正常的；2-已删除的；',
  `remark` varchar(500) NULL DEFAULT NULL COMMENT '备注',
  `create_by` bigint NULL DEFAULT NULL COMMENT '创建者',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint NULL DEFAULT NULL COMMENT '更新者',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`expert_id`) USING BTREE
) ENGINE = InnoDB COMMENT = '评委专家库';
-- 20250606 新增评委管理模块 评委专家库表 by liyu end

-- 20250613 擅长科目字典维护 by liyu start
INSERT INTO `sys_dict_type` (`dict_id`, `dict_name`, `dict_type`, `state`, `deleted`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`) VALUES (12, '科目', 'biz_dict_subjects', '1', 1, 1, '2025-06-13 10:47:48', 1, '2025-06-13 10:47:53', '评委专家擅长科目');

INSERT INTO `sys_dict_data` (`dict_code`, `dict_sort`, `dict_label`, `dict_value`, `dict_type`, `css_class`, `list_class`, `is_default`, `state`, `deleted`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`) VALUES (33, 1, '素描', '1', 'biz_dict_subjects', NULL, 'default', NULL, '1', 1, 1, '2025-06-13 10:49:25', 1, '2025-06-13 10:49:31', '擅长科目');
INSERT INTO `sys_dict_data` (`dict_code`, `dict_sort`, `dict_label`, `dict_value`, `dict_type`, `css_class`, `list_class`, `is_default`, `state`, `deleted`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`) VALUES (34, 2, '速写', '2', 'biz_dict_subjects', NULL, 'default', NULL, '1', 1, 1, '2025-06-13 10:50:10', 1, '2025-06-13 10:50:15', '擅长科目');
INSERT INTO `sys_dict_data` (`dict_code`, `dict_sort`, `dict_label`, `dict_value`, `dict_type`, `css_class`, `list_class`, `is_default`, `state`, `deleted`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`) VALUES (35, 3, '色彩', '3', 'biz_dict_subjects', NULL, 'default', NULL, '1', 1, 1, '2025-06-13 10:50:50', 1, '2025-06-13 10:50:55', '擅长科目');
-- 20250613 擅长科目字典维护 by liyu end

-- 20250613 同步考生报考专业数据表 by liyu start
DROP TABLE IF EXISTS `ea_sync_stu_prof`;
CREATE TABLE `ea_sync_stu_prof`  (
  `BaoKaoID` bigint NOT NULL COMMENT '报考专业ID',
  `YongHuID` bigint NULL DEFAULT NULL COMMENT '用户ID',
  `KaoShengID` bigint NULL DEFAULT NULL COMMENT '考生ID',
  `ZhengJianLX` tinyint(1) NULL DEFAULT NULL COMMENT '证件类型',
  `ShenFenZH` varchar(20) NULL DEFAULT NULL COMMENT '身份证号',
  `XingMing` varchar(100) NULL DEFAULT NULL COMMENT '姓名',
  `ZhunKaoZH` varchar(20) NULL DEFAULT NULL COMMENT '准考证号',
  `XueXiaoID` bigint NULL DEFAULT NULL COMMENT '学校ID',
  `XueXiaoMC` varchar(100) NULL DEFAULT NULL COMMENT '学校名称',
  `KaoShiID` bigint NULL DEFAULT NULL COMMENT '考试ID',
  `KaoShiMC` varchar(100) NULL DEFAULT NULL COMMENT '考试名称',
  `KaoDianID` bigint NULL DEFAULT NULL COMMENT '考点ID',
  `KaoDianMC` varchar(100) NULL DEFAULT NULL COMMENT '考点名称',
  `ZhuanYeID` bigint NULL DEFAULT NULL COMMENT '专业ID',
  `ZhuanYeMC` varchar(100) NULL DEFAULT NULL COMMENT '专业名称',
  `RiChengID` bigint NOT NULL COMMENT '日程ID',
  `fuRiChengID` bigint NULL DEFAULT NULL COMMENT '父日程ID',
  `KaoShiRQ` date NULL DEFAULT NULL COMMENT '考试日期',
  `KaoShiRQSM` varchar(64) NULL DEFAULT NULL COMMENT '考试日期说明',
  `KaoChangMC` varchar(128) NULL DEFAULT NULL COMMENT '考场名称',
  `gaoKaoSF` varchar(32) NULL DEFAULT NULL COMMENT '高考省份',
  `BaoMingFei` decimal(8, 2) NULL DEFAULT NULL COMMENT '报名费',
  `BaoKaoBZ` tinyint(1) NULL DEFAULT NULL COMMENT '报考标志 1-未提交 2-已提交 3-已生效 4-已关闭 5-作废',
  `QueRenFS` tinyint(1) NULL DEFAULT NULL COMMENT '确认方式 1-线上确认 2-线下确认',
  `QueRenSJ` datetime NULL DEFAULT NULL COMMENT '确认时间',
  `ShengFenMC` varchar(32) NULL DEFAULT NULL COMMENT '高考省份名称',
  `XingBie` varchar(32) NULL DEFAULT NULL COMMENT '性别',
  `TongXinDZ` varchar(512) NULL DEFAULT NULL COMMENT '通信地址',
  `TongXinYB` varchar(32) NULL DEFAULT NULL COMMENT '通信邮编',
  `ShouJi` varchar(256) NULL DEFAULT NULL COMMENT '本人手机',
  `KaoShengHao` varchar(32) NULL DEFAULT NULL COMMENT '考生号',
  `auditFlag` tinyint(1) NULL DEFAULT NULL COMMENT '申请审核状态 1-申请中 2-审核不通过 3-作废完成',
  `auditDes` varchar(128) NULL DEFAULT NULL COMMENT '审核意见',
  `printTime` datetime NULL DEFAULT NULL COMMENT '准考证打印时间',
  `WenLiKe` varchar(32) NULL DEFAULT NULL COMMENT '文理科',
  `videoCommitFlag` tinyint(1) NULL DEFAULT NULL COMMENT '网考视频提交状态 0-待提交 1-已提交',
  `otherPlatFlag` tinyint(1) NULL DEFAULT NULL COMMENT '报考数据来源：1-外部报考数据 2-缴费资格库',
  `finishedExamTime` datetime NULL DEFAULT NULL COMMENT '考试结束时间',
  `year` smallint NULL DEFAULT NULL COMMENT '年份',
  `syncTime` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '同步时间',
  `renZhengZP` varchar(256) NULL DEFAULT NULL COMMENT '认证照片',
  `gatlx` tinyint(1) NULL DEFAULT NULL COMMENT '港澳台类型',
  `zhengJianLXDesc` varchar(32) NULL DEFAULT NULL COMMENT '证件类型名称',
  `chuShengRQ` varchar(16) NULL DEFAULT NULL COMMENT '出生日期',
  `yingWangJie` varchar(16) NULL DEFAULT NULL COMMENT '应往届',
  `minZu` varchar(16) NULL DEFAULT NULL COMMENT '民族',
  `xueLi` varchar(16) NULL DEFAULT NULL COMMENT '学历',
  `zhengZhiMM` varchar(16) NULL DEFAULT NULL COMMENT '政治面貌',
  `singleSubjectName` varchar(128) NULL DEFAULT NULL COMMENT '首选科目',
  `suoZaiHS` varchar(128) NULL DEFAULT NULL COMMENT '专业课学习学校',
  `stuTypeStr` varchar(16) NULL DEFAULT NULL COMMENT '考生类型',
  `suoZaiXX` varchar(128) NULL DEFAULT NULL COMMENT '文化课学习学校',
  `subjectChoosen` varchar(128) NULL DEFAULT NULL COMMENT '选考科目',
  `subjectReChoose` varchar(128) NULL DEFAULT NULL COMMENT '再选科目',
  `deleted` tinyint(1) NOT NULL COMMENT '是否删除：1-未删除；2-已删除；',
  `create_by` bigint NOT NULL COMMENT '创建人ID',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint NULL DEFAULT NULL COMMENT '更新人ID',
  `update_time` datetime NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`BaoKaoID`) USING BTREE,
  INDEX `idx_shz_unique`(`ShenFenZH` ASC) USING BTREE
) ENGINE = InnoDB COMMENT = '考生报考专业同步数据表';
-- 20250613 同步考生报考专业数据表 by liyu end

-- 20250616 新增评委任务想搞模块 by 蓝鲸 start
CREATE TABLE `ea_judge_task_info` (
      `judge_task_id` bigint NOT NULL AUTO_INCREMENT COMMENT '评委征集任务ID',
      `judge_task_name` varchar(32) NOT NULL COMMENT '评委任务名称',
      `judge_task_type` tinyint(1) NOT NULL COMMENT '专家类型 1-博士 2-硕士 3-本科 4-附中',
      `task_start_time` datetime NOT NULL COMMENT '评分开始时间',
      `task_end_time` datetime NOT NULL COMMENT '评分结束时间',
      `draw_subject` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否按擅长科目抽取: 1-是 2-否',
      `state` tinyint(1) NOT NULL DEFAULT '1' COMMENT '抽取状态: 1-未开始 2-抽取中 3-已完成',
      `draw_times` int NOT NULL DEFAULT '0' COMMENT '抽取次数',
      `on_campus_total_num` int NOT NULL DEFAULT '0' COMMENT '校内评委需求量',
      `on_campus_attend_num` int NOT NULL DEFAULT '0' COMMENT '校内参加评委数',
      `on_campus_unattend_num` int NOT NULL DEFAULT '0' COMMENT '校内不参加评委数',
      `off_campus_total_num` int NOT NULL DEFAULT '0' COMMENT '校外评委需求量',
      `off_campus_attend_num` int NOT NULL DEFAULT '0' COMMENT '校外参加评委数',
      `off_campus_unattend_num` int NOT NULL DEFAULT '0' COMMENT '校外不参加评委数',
      `deleted` tinyint(1) DEFAULT '1' COMMENT '删除标记：1-正常的；2-已删除的；',
      `remark` varchar(500) DEFAULT NULL COMMENT '备注',
      `create_by` bigint DEFAULT NULL COMMENT '创建者',
      `create_time` datetime DEFAULT NULL COMMENT '创建时间',
      `update_by` bigint DEFAULT NULL COMMENT '更新者',
      `update_time` datetime DEFAULT NULL COMMENT '更新时间',
      PRIMARY KEY (`judge_task_id`) USING BTREE
) ENGINE=InnoDB COMMENT='评委征集任务信息表';

CREATE TABLE `ea_judge_draw_records` (
     `rules_id` bigint NOT NULL AUTO_INCREMENT COMMENT '规则ID',
     `judge_task_id` bigint NOT NULL COMMENT '任务ID',
     `draw_times` int NOT NULL DEFAULT '0' COMMENT '抽取次数',
     `subject_name` varchar(256) DEFAULT NULL COMMENT '科目名称',
     `on_campus_total_num` int NOT NULL DEFAULT '0' COMMENT '校内抽取数',
     `on_campus_attend_num` int NOT NULL DEFAULT '0' COMMENT '校内参加评委数',
     `on_campus_unattend_num` int NOT NULL DEFAULT '0' COMMENT '校内不参加评委数',
     `off_campus_total_num` int NOT NULL DEFAULT '0' COMMENT '校外抽取数',
     `off_campus_attend_num` int NOT NULL DEFAULT '0' COMMENT '校外参加评委数',
     `off_campus_unattend_num` int NOT NULL DEFAULT '0' COMMENT '校外不参加评委数',
     `deleted` tinyint(1) DEFAULT '1' COMMENT '删除标记：1-正常的；2-已删除的；',
     `remark` varchar(500) DEFAULT NULL COMMENT '备注',
     `create_by` bigint DEFAULT NULL COMMENT '创建者',
     `create_time` datetime DEFAULT NULL COMMENT '创建时间',
     `update_by` bigint DEFAULT NULL COMMENT '更新者',
     `update_time` datetime DEFAULT NULL COMMENT '更新时间',
     PRIMARY KEY (`rules_id`) USING BTREE
) ENGINE=InnoDB COMMENT='评委抽取规则表';

CREATE TABLE `ea_judge_draw_detail` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '关联ID',
    `expert_id` bigint NOT NULL COMMENT '专家ID',
    `expert_name` varchar(32) NOT NULL COMMENT '专家姓名',
    `dept_id` bigint NOT NULL COMMENT '专家部门id',
    `judge_task_id` bigint NOT NULL COMMENT '任务ID',
    `task_time_remark` varchar(64) NOT NULL COMMENT '评分日期',
    `draw_times` int NOT NULL DEFAULT '0' COMMENT '抽取次数',
    `draw_by_subject` tinyint(1) NOT NULL COMMENT '是否按擅长科目抽取 1-是 2-否',
    `subject_name` varchar(256) DEFAULT NULL COMMENT '科目名称',
    `confirm_state` char(1) NOT NULL DEFAULT '1' COMMENT '确认状态 1-未应答 2-确认参加 3-不参加',
    `deleted` tinyint(1) DEFAULT '1' COMMENT '删除标记：1-正常的；2-已删除的；',
    `remark` varchar(500) DEFAULT NULL COMMENT '备注',
    `create_by` bigint DEFAULT NULL COMMENT '创建者',
    `create_time` datetime DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint DEFAULT NULL COMMENT '更新者',
    `update_time` datetime DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB COMMENT='抽取评委明细表';

CREATE TABLE `ea_judge_experts_appraisal` (
    `appraisal_id` bigint NOT NULL AUTO_INCREMENT COMMENT '评价ID',
    `expert_id` bigint NOT NULL COMMENT '专家ID',
    `judge_task_id` bigint NOT NULL COMMENT '任务ID',
    `judge_task_name` varchar(32) NOT NULL COMMENT '评委任务名称',
    `task_time_remark` varchar(64) NOT NULL COMMENT '评分日期',
    `appraisal_content` varchar(256) NOT NULL COMMENT '评价内容',
    `deleted` tinyint(1) DEFAULT '1' COMMENT '删除标记：1-正常的；2-已删除的；',
    `remark` varchar(500) DEFAULT NULL COMMENT '备注',
    `create_by` bigint DEFAULT NULL COMMENT '创建者',
    `create_time` datetime DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint DEFAULT NULL COMMENT '更新者',
    `update_time` datetime DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (`appraisal_id`) USING BTREE
) ENGINE=InnoDB COMMENT='评委评价表';
-- 20250616 新增评委任务想搞模块 by 蓝鲸 end

-- 20250616 增加数据同步相关参数数据 by liyu start
INSERT INTO `sys_config` (`config_id`, `config_name`, `config_key`, `config_value`, `config_type`, `create_by`, `deleted`, `create_time`, `update_by`, `update_time`, `remark`) VALUES (100, '考试列表获取接口地址', 'sync_sch_exams_url', 'http://192.168.18.202:20800/api/o/school/apply/exam_list.htm', 'Y', 1, 1, '2025-06-16 14:28:00', NULL, NULL, '数据同步');
INSERT INTO `sys_config` (`config_id`, `config_name`, `config_key`, `config_value`, `config_type`, `create_by`, `deleted`, `create_time`, `update_by`, `update_time`, `remark`) VALUES (101, '报考专业同步接口地址', 'sync_stu_profs_url', 'http://192.168.18.202:20800/api/o/school/apply/stu_prof_list.htm', 'Y', 1, 1, '2025-06-16 14:28:00', NULL, NULL, '数据同步');
INSERT INTO `sys_config` (`config_id`, `config_name`, `config_key`, `config_value`, `config_type`, `create_by`, `deleted`, `create_time`, `update_by`, `update_time`, `remark`) VALUES (102, '数据同步接口token', 'sync_stu_prof_token', 'ECQuzvufLXfM9JfFC90ACweIyIQu0rSA', 'Y', 1, 1, '2025-06-16 14:28:00', NULL, NULL, '数据同步');
INSERT INTO `sys_config` (`config_id`, `config_name`, `config_key`, `config_value`, `config_type`, `create_by`, `deleted`, `create_time`, `update_by`, `update_time`, `remark`) VALUES (103, '数据同步接口每页处理数', 'sync_stu_prof_page_size', '100', 'Y', 1, 1, '2025-06-16 14:28:00', NULL, NULL, '数据同步');
-- 20250616 增加数据同步相关参数数据 by liyu start

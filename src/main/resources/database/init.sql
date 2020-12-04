
DROP TABLE EXISTS attach_file IF;
CREATE TABLE `attach_file` (
	`id` VARCHAR(32) NOT NULL COMMENT '主键' COLLATE 'utf8mb4_unicode_ci',
	`attach_file_name` VARCHAR(255) NOT NULL COMMENT '附件文件名' COLLATE 'utf8mb4_unicode_ci',
	`s3_url` VARCHAR(255) NOT NULL COMMENT 's3服务url' COLLATE 'utf8mb4_unicode_ci',
	`s3_bucket_name` VARCHAR(50) NULL DEFAULT NULL COMMENT 's3_bucket名称' COLLATE 'utf8mb4_unicode_ci',
	`s3_key_name` VARCHAR(50) NULL DEFAULT NULL COMMENT 's3_key名称' COLLATE 'utf8mb4_unicode_ci',
	`created_by` VARCHAR(255) NOT NULL COMMENT '创建人' COLLATE 'utf8mb4_unicode_ci',
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_by` VARCHAR(255) NOT NULL COMMENT '更新人' COLLATE 'utf8mb4_unicode_ci',
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '是否删除',
	PRIMARY KEY (`id`) USING BTREE
)
COMMENT='附件信息表'
COLLATE='utf8mb4_unicode_ci'
ENGINE=InnoDB
;

CREATE TABLE `form_info` (
	`id` VARCHAR(32) NOT NULL COMMENT '主键' COLLATE 'utf8mb4_unicode_ci',
	`form_temp_id` VARCHAR(32) NOT NULL COMMENT '表单模板id' COLLATE 'utf8mb4_unicode_ci',
	`name` VARCHAR(255) NOT NULL COMMENT '名称' COLLATE 'utf8mb4_unicode_ci',
	`type` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '类型',
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`created_by` VARCHAR(255) NOT NULL COMMENT '创建人' COLLATE 'utf8mb4_unicode_ci',
	`updated_by` VARCHAR(255) NOT NULL COMMENT '更新人' COLLATE 'utf8mb4_unicode_ci',
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '是否删除',
	PRIMARY KEY (`id`) USING BTREE
)
COMMENT='表单记录表'
COLLATE='utf8mb4_unicode_ci'
ENGINE=InnoDB
;
CREATE TABLE `form_item_info` (
	`id` VARCHAR(32) NOT NULL COMMENT '主键' COLLATE 'utf8mb4_unicode_ci',
	`form_id` VARCHAR(32) NOT NULL COMMENT '表单id' COLLATE 'utf8mb4_unicode_ci',
	`item_temp_id` VARCHAR(32) NOT NULL COMMENT '表单项模板id' COLLATE 'utf8mb4_unicode_ci',
	`name` VARCHAR(255) NOT NULL COMMENT '表单项名称' COLLATE 'utf8mb4_unicode_ci',
	`value` VARCHAR(255) NOT NULL COMMENT '表单项值' COLLATE 'utf8mb4_unicode_ci',
	`created_by` VARCHAR(255) NOT NULL COMMENT '创建人' COLLATE 'utf8mb4_unicode_ci',
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_by` VARCHAR(255) NOT NULL COMMENT '更新人' COLLATE 'utf8mb4_unicode_ci',
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '是否删除',
	PRIMARY KEY (`id`) USING BTREE
)
COMMENT='表单项记录表'
COLLATE='utf8mb4_unicode_ci'
ENGINE=InnoDB
;
CREATE TABLE `form_item_template` (
	`id` VARCHAR(32) NOT NULL COMMENT '主键' COLLATE 'utf8mb4_unicode_ci',
	`form_template_id` VARCHAR(32) NOT NULL COMMENT '表单模板' COLLATE 'utf8mb4_unicode_ci',
	`name` VARCHAR(50) NOT NULL COMMENT '名称' COLLATE 'utf8mb4_unicode_ci',
	`title` VARCHAR(50) NOT NULL COMMENT '标题' COLLATE 'utf8mb4_unicode_ci',
	`element_type` VARCHAR(50) NOT NULL DEFAULT 'text' COMMENT '元素类型' COLLATE 'utf8mb4_unicode_ci',
	`data_ci_id` VARCHAR(255) NOT NULL COMMENT 'ci数据id' COLLATE 'utf8mb4_unicode_ci',
	`data_filters` TEXT(65535) NULL DEFAULT NULL COMMENT 'ci数据检索条件' COLLATE 'utf8mb4_unicode_ci',
	`data_options` TEXT(65535) NULL DEFAULT NULL COMMENT '自定义数据源选项' COLLATE 'utf8mb4_unicode_ci',
	`is_public` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '是否通用',
	`required` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '必填',
	`is_edit` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '是否可编辑',
	`regular` TINYINT(2) NULL DEFAULT NULL COMMENT '正则表达式',
	`is_view` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '是否显示',
	`width` INT(11) NULL DEFAULT NULL COMMENT '长度',
	`def_value` VARCHAR(255) NULL DEFAULT NULL COMMENT '默认值' COLLATE 'utf8mb4_unicode_ci',
	`sort` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '排序',
	`created_by` VARCHAR(255) NOT NULL COMMENT '创建人' COLLATE 'utf8mb4_unicode_ci',
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_by` VARCHAR(255) NOT NULL COMMENT '更新人' COLLATE 'utf8mb4_unicode_ci',
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '是否删除',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `form_template_id` (`form_template_id`, `name`) USING BTREE
)
COMMENT='表单项模板表'
COLLATE='utf8mb4_unicode_ci'
ENGINE=InnoDB
;
CREATE TABLE `form_template` (
	`id` VARCHAR(32) NOT NULL COMMENT '主键' COLLATE 'utf8mb4_unicode_ci',
	`temp_id` VARCHAR(32) NOT NULL COMMENT '模板id' COLLATE 'utf8mb4_unicode_ci',
	`temp_type` VARCHAR(32) NOT NULL COMMENT '模板类型' COLLATE 'utf8mb4_unicode_ci',
	`name` VARCHAR(255) NOT NULL COMMENT '名称' COLLATE 'utf8mb4_unicode_ci',
	`target_entitys` VARCHAR(255) NOT NULL COMMENT '目标对象集' COLLATE 'utf8mb4_unicode_ci',
	`description` VARCHAR(512) NULL DEFAULT NULL COMMENT '描述' COLLATE 'utf8mb4_unicode_ci',
	`style` VARCHAR(50) NULL DEFAULT NULL COMMENT '表单风格' COLLATE 'utf8mb4_unicode_ci',
	`created_by` VARCHAR(255) NOT NULL COMMENT '创建人' COLLATE 'utf8mb4_unicode_ci',
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_by` VARCHAR(255) NOT NULL COMMENT '更新人' COLLATE 'utf8mb4_unicode_ci',
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '是否删除',
	PRIMARY KEY (`id`) USING BTREE
)
COMMENT='表单模板信息表'
COLLATE='utf8mb4_unicode_ci'
ENGINE=InnoDB
;
CREATE TABLE `request_info` (
	`id` VARCHAR(32) NOT NULL COMMENT '主键' COLLATE 'utf8mb4_unicode_ci',
	`form_id` VARCHAR(32) NOT NULL COMMENT '表单id' COLLATE 'utf8mb4_unicode_ci',
	`request_temp_id` VARCHAR(32) NOT NULL COMMENT '请求模板id' COLLATE 'utf8mb4_unicode_ci',
	`proc` VARCHAR(50) NOT NULL COMMENT '流程id' COLLATE 'utf8mb4_unicode_ci',
	`name` VARCHAR(255) NOT NULL COMMENT '名称' COLLATE 'utf8mb4_unicode_ci',
	`status` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '状态',
	`created_by` VARCHAR(255) NOT NULL COMMENT '创建人' COLLATE 'utf8mb4_unicode_ci',
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_by` VARCHAR(255) NOT NULL COMMENT '更新人' COLLATE 'utf8mb4_unicode_ci',
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '是否删除',
	PRIMARY KEY (`id`) USING BTREE
)
COMMENT='请求记录表'
COLLATE='utf8mb4_unicode_ci'
ENGINE=InnoDB
;
CREATE TABLE `request_template` (
	`id` VARCHAR(32) NOT NULL COMMENT '主键' COLLATE 'utf8mb4_unicode_ci',
	`request_temp_group` VARCHAR(32) NOT NULL COMMENT '模板组编号' COLLATE 'utf8mb4_unicode_ci',
	`proc_def_key` VARCHAR(255) NOT NULL COMMENT '流程编排key' COLLATE 'utf8mb4_unicode_ci',
	`proc_def_id` VARCHAR(255) NOT NULL COMMENT '流程编排id' COLLATE 'utf8mb4_unicode_ci',
	`proc_def_name` VARCHAR(255) NOT NULL COMMENT '流程编排名称' COLLATE 'utf8mb4_unicode_ci',
	`name` VARCHAR(255) NOT NULL COMMENT '名称' COLLATE 'utf8mb4_unicode_ci',
	`version` VARCHAR(50) NOT NULL DEFAULT '1' COMMENT '版本号' COLLATE 'utf8mb4_unicode_ci',
	`tags` VARCHAR(512) NULL DEFAULT NULL COMMENT '标签' COLLATE 'utf8mb4_unicode_ci',
	`status` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '状态',
	`created_by` VARCHAR(255) NOT NULL COMMENT '创建人' COLLATE 'utf8mb4_unicode_ci',
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_by` VARCHAR(255) NOT NULL COMMENT '更新人' COLLATE 'utf8mb4_unicode_ci',
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '是否删除',
	PRIMARY KEY (`id`) USING BTREE
)
COMMENT='请求模板信息表'
COLLATE='utf8mb4_unicode_ci'
ENGINE=InnoDB
;
CREATE TABLE `request_template_group` (
	`id` VARCHAR(32) NOT NULL COMMENT '主键' COLLATE 'utf8mb4_unicode_ci',
	`manage_role_id` VARCHAR(32) NULL DEFAULT NULL COMMENT '管理角色ID' COLLATE 'utf8mb4_unicode_ci',
	`manage_role_name` VARCHAR(255) NOT NULL COMMENT '管理角色名称' COLLATE 'utf8mb4_unicode_ci',
	`name` VARCHAR(255) NOT NULL COMMENT '名称' COLLATE 'utf8mb4_unicode_ci',
	`description` VARCHAR(512) NULL DEFAULT NULL COMMENT '描述' COLLATE 'utf8mb4_unicode_ci',
	`version` VARCHAR(50) NOT NULL DEFAULT '1' COMMENT '版本号' COLLATE 'utf8mb4_unicode_ci',
	`status` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '状态',
	`created_by` VARCHAR(255) NOT NULL COMMENT '创建人' COLLATE 'utf8mb4_unicode_ci',
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_by` VARCHAR(255) NOT NULL COMMENT '更新人' COLLATE 'utf8mb4_unicode_ci',
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '是否删除',
	PRIMARY KEY (`id`) USING BTREE
)
COMMENT='请求模板组信息表'
COLLATE='utf8mb4_unicode_ci'
ENGINE=InnoDB
;
CREATE TABLE `request_template_role` (
	`id` VARCHAR(32) NOT NULL COMMENT '主键' COLLATE 'utf8mb4_unicode_ci',
	`request_template_id` VARCHAR(32) NOT NULL COMMENT '请求模板编号' COLLATE 'utf8mb4_unicode_ci',
	`role_id` VARCHAR(32) NOT NULL COMMENT '角色编号' COLLATE 'utf8mb4_unicode_ci',
	`role_type` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '角色类型',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `request_template_id` (`request_template_id`, `role_id`, `role_type`) USING BTREE
)
COMMENT='请求模板角色关系表'
COLLATE='utf8mb4_unicode_ci'
ENGINE=InnoDB
;
CREATE TABLE `role_info` (
	`id` VARCHAR(32) NOT NULL COMMENT '主键' COLLATE 'utf8mb4_unicode_ci',
	`name` VARCHAR(255) NOT NULL COMMENT '角色名' COLLATE 'utf8mb4_unicode_ci',
	`display_name` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_unicode_ci',
	PRIMARY KEY (`id`) USING BTREE
)
COMMENT='角色信息表'
COLLATE='utf8mb4_unicode_ci'
ENGINE=InnoDB
;
CREATE TABLE `role_relation` (
	`id` VARCHAR(32) NOT NULL COMMENT '主键' COLLATE 'utf8mb4_unicode_ci',
	`role_id` VARCHAR(32) NOT NULL COMMENT '角色id' COLLATE 'utf8mb4_unicode_ci',
	`record_table` VARCHAR(255) NOT NULL COMMENT '记录所属表' COLLATE 'utf8mb4_unicode_ci',
	`record_id` VARCHAR(32) NOT NULL COMMENT '记录id' COLLATE 'utf8mb4_unicode_ci',
	`role_type` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '角色类型',
	PRIMARY KEY (`id`) USING BTREE
)
COMMENT='角色关联表 '
COLLATE='utf8mb4_unicode_ci'
ENGINE=InnoDB
;
CREATE TABLE `task_info` (
	`id` VARCHAR(32) NOT NULL COMMENT '主键' COLLATE 'utf8mb4_unicode_ci',
	`request_id` VARCHAR(32) NULL DEFAULT NULL COMMENT '请求id' COLLATE 'utf8mb4_unicode_ci',
	`request_no` VARCHAR(255) NULL DEFAULT NULL COMMENT '请求编号' COLLATE 'utf8mb4_unicode_ci',
	`parent_id` VARCHAR(32) NULL DEFAULT NULL COMMENT '父级任务id' COLLATE 'utf8mb4_unicode_ci',
	`task_temp_id` VARCHAR(32) NOT NULL COMMENT '任务模板id' COLLATE 'utf8mb4_unicode_ci',
	`form_id` VARCHAR(32) NULL DEFAULT NULL COMMENT '表单id' COLLATE 'utf8mb4_unicode_ci',
	`proc_node` VARCHAR(50) NOT NULL COMMENT '流程节点' COLLATE 'utf8mb4_unicode_ci',
	`callback_url` VARCHAR(255) NOT NULL COMMENT '回调url' COLLATE 'utf8mb4_unicode_ci',
	`callback_parameter` TEXT(65535) NOT NULL COMMENT '回调参数' COLLATE 'utf8mb4_unicode_ci',
	`name` VARCHAR(255) NOT NULL COMMENT '任务名称' COLLATE 'utf8mb4_unicode_ci',
	`reporter` VARCHAR(50) NULL DEFAULT NULL COMMENT '上报人' COLLATE 'utf8mb4_unicode_ci',
	`report_time` DATETIME NULL DEFAULT NULL COMMENT '上报时间',
	`emergency` VARCHAR(50) NULL DEFAULT NULL COMMENT '紧急程度' COLLATE 'utf8mb4_unicode_ci',
	`report_role` VARCHAR(50) NULL DEFAULT NULL COMMENT '上报角色' COLLATE 'utf8mb4_unicode_ci',
	`result` TEXT(65535) NULL DEFAULT NULL COMMENT '执行结果' COLLATE 'utf8mb4_unicode_ci',
	`description` VARCHAR(512) NULL DEFAULT NULL COMMENT '描述' COLLATE 'utf8mb4_unicode_ci',
	`attach_file_id` VARCHAR(32) NULL DEFAULT NULL COMMENT '附件id' COLLATE 'utf8mb4_unicode_ci',
	`status` TINYINT(2) NULL DEFAULT '0' COMMENT '状态',
	`version` VARCHAR(50) NULL DEFAULT NULL COMMENT '版本号' COLLATE 'utf8mb4_unicode_ci',
	`created_by` VARCHAR(255) NOT NULL COMMENT '创建人' COLLATE 'utf8mb4_unicode_ci',
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_by` VARCHAR(255) NOT NULL COMMENT '更新人' COLLATE 'utf8mb4_unicode_ci',
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '是否删除',
	PRIMARY KEY (`id`) USING BTREE
)
COMMENT='任务记录表'
COLLATE='utf8mb4_unicode_ci'
ENGINE=InnoDB
;
CREATE TABLE `task_template` (
	`id` VARCHAR(32) NOT NULL COMMENT '主键' COLLATE 'utf8mb4_unicode_ci',
	`proc_def_id` VARCHAR(255) NOT NULL COMMENT '流程编排id' COLLATE 'utf8mb4_unicode_ci',
	`proc_def_key` VARCHAR(255) NOT NULL COMMENT '流程编排key' COLLATE 'utf8mb4_unicode_ci',
	`proc_def_name` VARCHAR(255) NOT NULL COMMENT '流程编排名称' COLLATE 'utf8mb4_unicode_ci',
	`proc_node` VARCHAR(255) NOT NULL COMMENT '流程节点' COLLATE 'utf8mb4_unicode_ci',
	`name` VARCHAR(255) NOT NULL COMMENT '名称' COLLATE 'utf8mb4_unicode_ci',
	`description` VARCHAR(512) NULL DEFAULT NULL COMMENT '描述' COLLATE 'utf8mb4_unicode_ci',
	`created_by` VARCHAR(255) NOT NULL COMMENT '创建人' COLLATE 'utf8mb4_unicode_ci',
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`updated_by` VARCHAR(255) NOT NULL COMMENT '更新人' COLLATE 'utf8mb4_unicode_ci',
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0' COMMENT '是否删除',
	PRIMARY KEY (`id`) USING BTREE
)
COMMENT='任务模板信息表'
COLLATE='utf8mb4_unicode_ci'
ENGINE=InnoDB
;

-- 正在导出表  taskman.form_item_template 的数据：~12 rows (大约)
/*!40000 ALTER TABLE `form_item_template` DISABLE KEYS */;
INSERT INTO `form_item_template` (`id`, `form_template_id`, `name`, `title`, `element_type`, `data_ci_id`, `data_filters`, `data_options`, `is_public`, `required`, `is_edit`, `regular`, `is_view`, `width`, `def_value`, `sort`, `created_by`, `created_time`, `updated_by`, `updated_time`, `del_flag`) VALUES
	('1333228970392473601', '', '上报人', '上报人', 'text', '', NULL, NULL, 0, 0, 0, NULL, 0, NULL, NULL, 0, '', '2020-12-01 16:23:53', '', '2020-12-01 16:23:53', 0),
	('1333297882663247873', '', '任务名称', '任务名称', 'text', '', NULL, NULL, 0, 0, 0, NULL, 0, NULL, NULL, 0, '', '2020-12-01 16:23:53', '', '2020-12-01 16:23:53', 0),
	('1333303089673617409', '', '紧急程度', '紧急程度', 'select', '', NULL, NULL, 0, 0, 0, NULL, 0, NULL, NULL, 0, '', '2020-12-01 16:23:53', '', '2020-12-01 16:23:53', 0),
	('1333303897844670466', '', '任务描述', '任务描述', 'text', '', NULL, NULL, 0, 0, 0, NULL, 0, NULL, NULL, 0, '', '2020-12-01 16:23:53', '', '2020-12-01 16:23:53', 0),
	('1333304415006547970', '', '任务附件', '任务附件', 'file', '', NULL, NULL, 0, 0, 0, NULL, 0, NULL, NULL, 0, '', '2020-12-01 16:23:53', '', '2020-12-01 16:23:53', 0),
	('1333319171714420738', '', '处理结果', '处理结果', 'text', '', NULL, NULL, 0, 0, 0, NULL, 0, NULL, NULL, 0, '', '2020-12-01 16:23:54', '', '2020-12-01 16:23:54', 0),
	('1333324857626169346', '', '上报时间', '上报时间', 'date', '', NULL, NULL, 0, 0, 0, NULL, 0, NULL, NULL, 0, '', '2020-12-01 16:23:53', '', '2020-12-01 16:23:53', 0);
/*!40000 ALTER TABLE `form_item_template` ENABLE KEYS */;


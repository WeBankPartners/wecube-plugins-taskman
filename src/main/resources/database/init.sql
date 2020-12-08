SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS attach_file;
CREATE TABLE `attach_file` (
	`id` VARCHAR(32) NOT NULL,
	`attach_file_name` VARCHAR(255) NOT NULL,
	`s3_url` VARCHAR(255) NOT NULL,
	`s3_bucket_name` VARCHAR(50) NULL DEFAULT NULL,
	`s3_key_name` VARCHAR(50) NULL DEFAULT NULL,
	`created_by` VARCHAR(255) NOT NULL,
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_by` VARCHAR(255) NOT NULL,
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS form_info;
CREATE TABLE `form_info` (
	`id` VARCHAR(32) NOT NULL,
	`form_temp_id` VARCHAR(32) NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`type` TINYINT(2) NOT NULL DEFAULT '0',
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`created_by` VARCHAR(255) NOT NULL,
	`updated_by` VARCHAR(255) NOT NULL,
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS form_item_info;
CREATE TABLE `form_item_info` (
	`id` VARCHAR(32) NOT NULL,
	`form_id` VARCHAR(32) NOT NULL,
	`item_temp_id` VARCHAR(32) NOT NULL,
	`name` VARCHAR(255) NOT NULL,
	`value` VARCHAR(255) NOT NULL,
	`created_by` VARCHAR(255) NOT NULL,
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_by` VARCHAR(255) NOT NULL,
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS form_item_template;
CREATE TABLE `form_item_template` (
	`id` VARCHAR(32) NOT NULL,
	`form_template_id` VARCHAR(32) NOT NULL,
	`name` VARCHAR(50) NOT NULL,
	`title` VARCHAR(50) NOT NULL,
	`element_type` VARCHAR(50) NOT NULL DEFAULT 'text',
	`data_ci_id` VARCHAR(255) NOT NULL,
	`data_filters` TEXT(65535) NULL DEFAULT NULL,
	`data_options` TEXT(65535) NULL DEFAULT NULL,
	`is_public` TINYINT(2) NOT NULL DEFAULT '0',
	`required` TINYINT(2) NOT NULL DEFAULT '0',
	`is_edit` TINYINT(2) NOT NULL DEFAULT '0',
	`regular` TINYINT(2) NULL DEFAULT NULL ,
	`is_view` TINYINT(2) NOT NULL DEFAULT '0',
	`width` INT(11) NULL DEFAULT NULL,
	`def_value` VARCHAR(255) NULL DEFAULT NULL,
	`sort` TINYINT(2) NOT NULL DEFAULT '0',
	`created_by` VARCHAR(255) NOT NULL,
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_by` VARCHAR(255) NOT NULL,
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `form_template_id` (`form_template_id`, `name`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS form_template;
CREATE TABLE `form_template` (
	`id` VARCHAR(32) NOT NULL,
	`temp_id` VARCHAR(32) NOT NULL,
	`temp_type` VARCHAR(32) NOT NULL ,
	`name` VARCHAR(255) NOT NULL ,
	`target_entitys` VARCHAR(255) NOT NULL ,
	`description` VARCHAR(512) NULL DEFAULT NULL ,
	`style` VARCHAR(50) NULL DEFAULT NULL ,
	`created_by` VARCHAR(255) NOT NULL,
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_by` VARCHAR(255) NOT NULL,
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS request_info;
CREATE TABLE `request_info` (
	`id` VARCHAR(32) NOT NULL,
	`form_id` VARCHAR(32) NOT NULL,
	`request_temp_id` VARCHAR(32) NOT NULL,
	`proc` VARCHAR(50) NOT NULL,
	`name` VARCHAR(255) NOT NULL ,
	`status` TINYINT(2) NOT NULL DEFAULT '0' ,
	`created_by` VARCHAR(255) NOT NULL,
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_by` VARCHAR(255) NOT NULL,
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS request_template;
CREATE TABLE `request_template` (
	`id` VARCHAR(32) NOT NULL,
	`request_temp_group` VARCHAR(32) NOT NULL ,
	`proc_def_key` VARCHAR(255) NOT NULL ,
	`proc_def_id` VARCHAR(255) NOT NULL,
	`proc_def_name` VARCHAR(255) NOT NULL ,
	`name` VARCHAR(255) NOT NULL ,
	`version` VARCHAR(50) NOT NULL DEFAULT '1' ,
	`tags` VARCHAR(512) NULL DEFAULT NULL ,
	`status` TINYINT(2) NOT NULL DEFAULT '0' ,
	`created_by` VARCHAR(255) NOT NULL,
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_by` VARCHAR(255) NOT NULL,
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS request_template_group;
CREATE TABLE `request_template_group` (
	`id` VARCHAR(32) NOT NULL,
	`manage_role_id` VARCHAR(32) NULL DEFAULT NULL,
	`manage_role_name` VARCHAR(255) NOT NULL ,
	`name` VARCHAR(255) NOT NULL ,
	`description` VARCHAR(512) NULL DEFAULT NULL ,
	`version` VARCHAR(50) NOT NULL DEFAULT '1' ,
	`status` TINYINT(2) NOT NULL DEFAULT '0' ,
	`created_by` VARCHAR(255) NOT NULL,
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_by` VARCHAR(255) NOT NULL,
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS role_relation;
CREATE TABLE `role_relation` (
	`id` VARCHAR(32) NOT NULL,
	`record_table` VARCHAR(255) NOT NULL ,
	`record_id` VARCHAR(32) NOT NULL,
	`role_type` TINYINT(2) NOT NULL DEFAULT '0' ,
	`role_name` VARCHAR(50) NULL DEFAULT '0',
	`display_name` VARCHAR(50) NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS role_relation;
CREATE TABLE `task_info` (
	`id` VARCHAR(32) NOT NULL,
	`request_id` VARCHAR(32) NULL DEFAULT NULL,
	`request_no` VARCHAR(255) NULL DEFAULT NULL ,
	`parent_id` VARCHAR(32) NULL DEFAULT NULL,
	`task_temp_id` VARCHAR(32) NOT NULL,
	`form_id` VARCHAR(32) NULL DEFAULT NULL,
	`proc_node` VARCHAR(50) NOT NULL ,
	`callback_url` VARCHAR(255) NOT NULL COMMENT 'url',
	`callback_parameter` TEXT(65535) NOT NULL ,
	`name` VARCHAR(255) NOT NULL ,
	`reporter` VARCHAR(50) NULL DEFAULT NULL ,
	`report_time` DATETIME NULL DEFAULT NULL ,
	`emergency` VARCHAR(50) NULL DEFAULT NULL ,
	`report_role` VARCHAR(50) NULL DEFAULT NULL ,
	`result` TEXT(65535) NULL DEFAULT NULL ,
	`description` VARCHAR(512) NULL DEFAULT NULL ,
	`attach_file_id` VARCHAR(32) NULL DEFAULT NULL,
	`status` TINYINT(2) NULL DEFAULT '0' ,
	`version` VARCHAR(50) NULL DEFAULT NULL ,
	`created_by` VARCHAR(255) NOT NULL,
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_by` VARCHAR(255) NOT NULL,
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS task_template;
CREATE TABLE `task_template` (
	`id` VARCHAR(32) NOT NULL,
	`proc_def_id` VARCHAR(255) NOT NULL,
	`proc_def_key` VARCHAR(255) NOT NULL,
	`proc_def_name` VARCHAR(255) NOT NULL ,
	`proc_node` VARCHAR(255) NOT NULL ,
	`name` VARCHAR(255) NOT NULL ,
	`description` VARCHAR(512) NULL DEFAULT NULL ,
	`created_by` VARCHAR(255) NOT NULL,
	`created_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_by` VARCHAR(255) NOT NULL,
	`updated_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`del_flag` TINYINT(2) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;


INSERT INTO `form_item_template` (`id`, `form_template_id`, `name`, `title`, `element_type`, `data_ci_id`, `data_filters`, `data_options`, `is_public`, `required`, `is_edit`, `regular`, `is_view`, `width`, `def_value`, `sort`, `created_by`, `created_time`, `updated_by`, `updated_time`, `del_flag`) VALUES
	('1333228970392473601', '', '', '', 'text', '', NULL, NULL, 0, 0, 0, NULL, 0, NULL, NULL, 0, '', '2020-12-01 16:23:53', '', '2020-12-01 16:23:53', 0),
	('1333297882663247873', '', '', '', 'text', '', NULL, NULL, 0, 0, 0, NULL, 0, NULL, NULL, 0, '', '2020-12-01 16:23:53', '', '2020-12-01 16:23:53', 0),
	('1333303089673617409', '', '', '', 'select', '', NULL, NULL, 0, 0, 0, NULL, 0, NULL, NULL, 0, '', '2020-12-01 16:23:53', '', '2020-12-01 16:23:53', 0),
	('1333303897844670466', '', '', '', 'text', '', NULL, NULL, 0, 0, 0, NULL, 0, NULL, NULL, 0, '', '2020-12-01 16:23:53', '', '2020-12-01 16:23:53', 0),
	('1333304415006547970', '', '', '', 'file', '', NULL, NULL, 0, 0, 0, NULL, 0, NULL, NULL, 0, '', '2020-12-01 16:23:53', '', '2020-12-01 16:23:53', 0),
	('1333319171714420738', '', '', '', 'text', '', NULL, NULL, 0, 0, 0, NULL, 0, NULL, NULL, 0, '', '2020-12-01 16:23:54', '', '2020-12-01 16:23:54', 0),
	('1333324857626169346', '', '', '', 'date', '', NULL, NULL, 0, 0, 0, NULL, 0, NULL, NULL, 0, '', '2020-12-01 16:23:53', '', '2020-12-01 16:23:53', 0);



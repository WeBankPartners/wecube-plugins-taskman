SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `attach_file`;
CREATE TABLE IF NOT EXISTS `attach_file` (
  `id` varchar(64)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `s3_url` varchar(255)  NOT NULL COMMENT 's3url',
  `s3_bucket_name` varchar(255)  DEFAULT NULL COMMENT 's3_bucket',
  `s3_key_name` varchar(255)  DEFAULT NULL COMMENT 's3_key',
  `created_by` varchar(255)  NULL,
  `created_time` datetime NULL,
  `updated_by` varchar(255)  NULL,
  `updated_time` datetime NULL,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ;

DROP TABLE IF EXISTS `form_info`;
CREATE TABLE IF NOT EXISTS `form_info` (
  `id` varchar(32)  NOT NULL,
  `record_id` varchar(32)  NOT NULL,
  `form_template_id` varchar(50)  DEFAULT NULL,
  `name` varchar(255)  NOT NULL,
  `type` tinyint(2) NOT NULL DEFAULT '0',
  `created_time` datetime NULL,
  `created_by` varchar(255)   NULL,
  `updated_by` varchar(255)   NULL,
  `updated_time` datetime NULL,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  `form_type` varchar(45) NULL,
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ;

DROP TABLE IF EXISTS `form_item_info`;
CREATE TABLE IF NOT EXISTS `form_item_info` (
  `id` varchar(32)  NOT NULL,
  `form_id` varchar(32)  NOT NULL,
  `item_temp_id` varchar(32)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `value` varchar(255)  DEFAULT NULL,
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ;

DROP TABLE IF EXISTS `form_item_template`;
CREATE TABLE IF NOT EXISTS `form_item_template` (
  `id` varchar(64) NOT NULL,
  `form_template_id` varchar(64) NOT NULL,
  `input_type` varchar(64) default 'text',
  `name` varchar(64) NOT NULL,
  `description` varchar(255)  DEFAULT NULL,
  `default_value` varchar(255) DEFAULT NULL,
  `sort` int NOT NULL DEFAULT '0',
  `package_name` varchar(255) DEFAULT '0',
  `entity` varchar(255) DEFAULT '0',
  `attr_def_id` varchar(255) DEFAULT NULL,
  `attr_def_data_type` varchar(255) DEFAULT NULL,
  `element_type` varchar(50) NOT NULL DEFAULT 'text',
  `title` varchar(50) NOT NULL,
  `width` int(11) DEFAULT NULL,
  `ref_package_name` varchar(255) DEFAULT NULL,
  `ref_entity` varchar(255) DEFAULT NULL,
  `ref_filters` text COMMENT 'ciid',
  `data_options` text,
  `required` tinyint(2) NOT NULL DEFAULT '0',
  `regular` varchar(50) DEFAULT NULL,
  `is_edit` tinyint(2) NOT NULL DEFAULT '0',
  `is_view` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


INSERT INTO `form_item_template`
(`id`, `form_template_id`, `name`, `default_value`, `is_currency`, `sort`, `package_name`, `entity`, `attr_def_id`, `attr_def_data_type`, `element_type`, `title`, `width`, `ref_package_name`, `ref_entity`, `ref_filters`, `data_options`, `required`, `regular`, `is_edit`, `is_view`)
 VALUES
    ('1333228970392473601', '', 'reporter', NULL, 0, 0, '0', '0', NULL, NULL, 'text', '上报人', NULL, NULL, NULL, '', NULL, 0, NULL, 0, 0),
    ('1333303089673617409', '', 'emergency', NULL, 0, 0, '0', '0', NULL, NULL, 'select', '紧急程度', NULL, NULL, NULL, '', NULL, 0, NULL, 0, 0),
    ('1333304415006547970', '', 'attach_file_id', NULL, 0, 0, '0', '0', NULL, NULL, 'file', '附件', NULL, NULL, NULL, '', NULL, 0, NULL, 0, 0),
    ('1333319171714420738', '', 'result', NULL, 0, 0, '0', '0', NULL, NULL, 'text', '处理结果', NULL, NULL, NULL, '', NULL, 0, NULL, 0, 0),
    ('1333324857626169346', '', 'report_time', NULL, 0, 0, '0', '0', NULL, NULL, 'date', '上报时间', NULL, NULL, NULL, '', NULL, 0, NULL, 0, 0),
    ('1333324857626169347', '', 'due_data', NULL, 0, 0, '0', '0', NULL, NULL, 'date', '过期时间', NULL, NULL, NULL, '', NULL, 0, NULL, 0, 0);

DROP TABLE IF EXISTS `form_template`;
CREATE TABLE IF NOT EXISTS `form_template` (
  `id` varchar(32)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `description` varchar(512)  DEFAULT NULL,
  `role` varchar(64)  DEFAULT NULL,
  `created_by` varchar(255)   NULL,
  `created_time` datetime NULL,
  `updated_by` varchar(255)   NULL,
  `updated_time` datetime NULL,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ;

DROP TABLE IF EXISTS `request_info`;
CREATE TABLE IF NOT EXISTS `request_info` (
  `id` varchar(32)  NOT NULL,
  `request_temp_id` varchar(32)  NOT NULL,
  `proc_inst_key` varchar(50)  DEFAULT NULL,
  `name` varchar(255)  NOT NULL,
  `reporter` varchar(50)  DEFAULT NULL,
  `report_time` datetime DEFAULT NULL,
  `emergency` varchar(50)  DEFAULT NULL,
  `report_role` varchar(50)  DEFAULT NULL,
  `attach_file_id` varchar(32)  DEFAULT NULL,
  `status` varchar(50)  NOT NULL DEFAULT '0',
  `result` mediumtext ,
  `created_by` varchar(255)   NULL,
  `created_time` datetime NULL,
  `updated_by` varchar(255)   NULL,
  `updated_time` datetime NULL,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ;

DROP TABLE IF EXISTS `request_template`;
CREATE TABLE IF NOT EXISTS `request_template` (
  `id` varchar(64)  NOT NULL,
  `group` varchar(64)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `description` varchar(512)  DEFAULT NULL,
  `form_template` varchar(64) DEFAULT NULL,
  `tags` varchar(512)  DEFAULT NULL,
  `status` varchar(32) NOT NULL DEFAULT 'created',
  `package_name` varchar(255)  NULL,
  `entity_name` varchar(255)  NULL,
  `proc_def_key` varchar(255)  NOT NULL COMMENT 'key',
  `proc_def_id` varchar(255)  NOT NULL,
  `proc_def_name` varchar(255)  NOT NULL,
  `created_by` varchar(255)   NULL,
  `created_time` datetime NULL,
  `updated_by` varchar(255)   NULL,
  `updated_time` datetime NULL,
  `entity_attrs` text DEFAULT NULL,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ;

DROP TABLE IF EXISTS `request_template_group`;
CREATE TABLE IF NOT EXISTS `request_template_group` (
  `id` varchar(32)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `description` varchar(255)  DEFAULT NULL,
  `manage_role` varchar(64)  DEFAULT NULL,
  `created_by` varchar(255)   NULL,
  `created_time` datetime NULL,
  `updated_by` varchar(255)   NULL,
  `updated_time` datetime NULL,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ;

DROP TABLE IF EXISTS `role`;
CREATE TABLE IF NOT EXISTS `role` (
  `id` varchar(64)  NOT NULL,
  `display_name` varchar(64) NOT NULL,
  `update_time` datetime,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT=' ';

DROP TABLE IF EXISTS `task_info`;
CREATE TABLE IF NOT EXISTS `task_info` (
  `id` varchar(32)  NOT NULL,
  `request_id` varchar(32)  DEFAULT NULL,
  `request_no` varchar(255)  DEFAULT NULL,
  `parent_id` varchar(32)  DEFAULT NULL,
  `task_temp_id` varchar(32)  NOT NULL,
  `node_def_id` varchar(50)  NOT NULL,
  `node_name` varchar(50)  NOT NULL,
  `callback_url` varchar(255)  NOT NULL COMMENT 'url',
  `callback_parameter` text  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `reporter` varchar(50)  DEFAULT NULL,
  `report_time` datetime DEFAULT NULL,
  `emergency` varchar(50)  DEFAULT NULL,
  `report_role` varchar(50)  DEFAULT NULL,
  `result` text ,
  `description` varchar(512)  DEFAULT NULL,
  `attach_file_id` varchar(32)  DEFAULT NULL,
  `status` tinyint(2) DEFAULT '0',
  `version` varchar(50)  DEFAULT NULL,
  `created_by` varchar(255)   NULL,
  `created_time` datetime NULL,
  `updated_by` varchar(255)   NULL,
  `updated_time` datetime NULL,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ;

DROP TABLE IF EXISTS `task_template`;
CREATE TABLE IF NOT EXISTS `task_template` (
  `id` varchar(32)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `description` varchar(512)  DEFAULT NULL,
  `form_template` varchar(64) DEFAULT NULL,
  `request_template_id` varchar(64) DEFAULT NULL,
  `proc_def_id` varchar(255)  NOT NULL,
  `proc_def_key` varchar(255)  NOT NULL,
  `proc_def_name` varchar(255)  NOT NULL,
  `node_def_id` varchar(255)  NOT NULL,
  `node_name` varchar(255)  NOT NULL,
  `created_by` varchar(255)   NULL,
  `created_time` datetime NULL,
  `updated_by` varchar(255)   NULL,
  `updated_time` datetime NULL,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ;

SET FOREIGN_KEY_CHECKS = 1;
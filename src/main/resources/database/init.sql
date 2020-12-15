SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `attach_file`;
CREATE TABLE IF NOT EXISTS `attach_file` (
  `id` varchar(32)  NOT NULL,
  `attach_file_name` varchar(255)  NOT NULL,
  `s3_url` varchar(255)  NOT NULL COMMENT 's3url',
  `s3_bucket_name` varchar(50)  DEFAULT NULL COMMENT 's3_bucket',
  `s3_key_name` varchar(50)  DEFAULT NULL COMMENT 's3_key',
  `created_by` varchar(255)  NOT NULL,
  `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` varchar(255)  NOT NULL,
  `updated_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
  `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` varchar(255)  NOT NULL,
  `updated_by` varchar(255)  NOT NULL,
  `updated_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
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
  `id` varchar(32)  NOT NULL,
  `form_template_id` varchar(32)  NOT NULL,
  `attr_def_id` varchar(50)  DEFAULT NULL,
  `attr_data_type` varchar(50)  DEFAULT NULL,
  `name` varchar(50)  NOT NULL,
  `title` varchar(50)  NOT NULL,
  `package_name` varchar(50)  DEFAULT '0',
  `entity_name` varchar(50)  DEFAULT '0',
  `entity_filters` mediumtext  COMMENT 'ciid',
  `element_type` varchar(50)  NOT NULL DEFAULT 'text',
  `is_currency` tinyint(2) NOT NULL DEFAULT '0',
  `default_value` varchar(255)  DEFAULT NULL,
  `required` tinyint(2) NOT NULL DEFAULT '0',
  `is_edit` tinyint(2) NOT NULL DEFAULT '0',
  `regular` varchar(50)  DEFAULT NULL,
  `width` int(11) DEFAULT NULL,
  `sort` tinyint(2) NOT NULL DEFAULT '0',
  `data_options` text ,
  `is_view` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `form_template_id` (`form_template_id`,`name`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ;


INSERT IGNORE  INTO `form_item_template` (`id`, `form_template_id`, `attr_def_id`, `attr_data_type`, `name`, `title`, `package_name`, `entity_name`, `entity_filters`, `element_type`, `is_currency`, `default_value`, `required`, `is_edit`, `regular`, `width`, `sort`, `data_options`, `is_view`) VALUES
('1333228970392473601', '', NULL, NULL, '上报人', '上报人', '0', '0', '', 'text', 0, NULL, 0, 0, NULL, NULL, 0, NULL, 0),
('1333297882663247873', '', NULL, NULL, '任务名称', '任务名称', '0', '0', '', 'text', 0, NULL, 0, 0, NULL, NULL, 0, NULL, 0),
('1333303089673617409', '', NULL, NULL, '紧急程度', '紧急程度', '0', '0', '', 'select', 0, NULL, 0, 0, NULL, NULL, 0, NULL, 0),
('1333303897844670466', '', NULL, NULL, '任务描述', '任务描述', '0', '0', '', 'text', 0, NULL, 0, 0, NULL, NULL, 0, NULL, 0),
('1333304415006547970', '', NULL, NULL, '任务附件', '任务附件', '0', '0', '', 'file', 0, NULL, 0, 0, NULL, NULL, 0, NULL, 0),
('1333319171714420738', '', NULL, NULL, '处理结果', '处理结果', '0', '0', '', 'text', 0, NULL, 0, 0, NULL, NULL, 0, NULL, 0),
('1333324857626169346', '', NULL, NULL, '上报时间', '上报时间', '0', '0', '', 'date', 0, NULL, 0, 0, NULL, NULL, 0, NULL, 0);


DROP TABLE IF EXISTS `form_template`;
CREATE TABLE IF NOT EXISTS `form_template` (
  `id` varchar(32)  NOT NULL,
  `temp_id` varchar(32)  NOT NULL,
  `temp_type` varchar(32)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `description` varchar(512)  DEFAULT NULL,
  `target_entitys` varchar(255)  NOT NULL,
  `input_attr_def` text ,
  `output_attr_def` text ,
  `other_attr_def` text ,
  `style` varchar(50)  DEFAULT NULL,
  `created_by` varchar(255)  NOT NULL,
  `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` varchar(255)  NOT NULL,
  `updated_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
  `created_by` varchar(255)  NOT NULL,
  `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` varchar(255)  NOT NULL,
  `updated_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ;

DROP TABLE IF EXISTS `request_template`;
CREATE TABLE IF NOT EXISTS `request_template` (
  `id` varchar(32)  NOT NULL,
  `request_temp_group` varchar(32)  NOT NULL,
  `proc_def_key` varchar(255)  NOT NULL COMMENT 'key',
  `proc_def_id` varchar(255)  NOT NULL,
  `proc_def_name` varchar(255)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `version` varchar(50)  NOT NULL DEFAULT '1',
  `tags` varchar(512)  DEFAULT NULL,
  `status` tinyint(2) NOT NULL DEFAULT '0',
  `created_by` varchar(255)  NOT NULL,
  `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` varchar(255)  NOT NULL,
  `updated_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ;

DROP TABLE IF EXISTS `request_template_group`;
CREATE TABLE IF NOT EXISTS `request_template_group` (
  `id` varchar(32)  NOT NULL,
  `manage_role_id` varchar(32)  DEFAULT NULL,
  `manage_role_name` varchar(255)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `description` varchar(512)  DEFAULT NULL,
  `version` varchar(50)  NOT NULL DEFAULT '1',
  `status` tinyint(2) NOT NULL DEFAULT '0',
  `created_by` varchar(255)  NOT NULL,
  `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` varchar(255)  NOT NULL,
  `updated_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ;

DROP TABLE IF EXISTS `role_relation`;
CREATE TABLE IF NOT EXISTS `role_relation` (
  `id` varchar(32)  NOT NULL,
  `record_id` varchar(32)  NOT NULL,
  `role_type` tinyint(2) NOT NULL DEFAULT '0',
  `role_name` varchar(50)  DEFAULT '0',
  `display_name` varchar(50)  DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `ft_index_role_type` (`role_type`),
  FULLTEXT KEY `ft_index_role_name` (`role_name`)
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
  `created_by` varchar(255)  NOT NULL,
  `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` varchar(255)  NOT NULL,
  `updated_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ;

DROP TABLE IF EXISTS `task_template`;
CREATE TABLE IF NOT EXISTS `task_template` (
  `id` varchar(32)  NOT NULL,
  `proc_def_id` varchar(255)  NOT NULL,
  `proc_def_key` varchar(255)  NOT NULL COMMENT 'key',
  `proc_def_name` varchar(255)  NOT NULL,
  `node_def_id` varchar(255)  NOT NULL,
  `node_name` varchar(255)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `description` varchar(512)  DEFAULT NULL,
  `created_by` varchar(255)  NOT NULL,
  `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` varchar(255)  NOT NULL,
  `updated_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ;

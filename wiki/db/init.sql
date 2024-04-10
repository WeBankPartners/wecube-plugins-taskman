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
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT=' ';

DROP TABLE IF EXISTS `role`;
CREATE TABLE IF NOT EXISTS `role` (
  `id` varchar(64)  NOT NULL,
  `display_name` varchar(255) NOT NULL,
  `core_id` varchar(64) DEFAULT NULL,
  `email` varchar(64) DEFAULT NULL,
  `updated_time` datetime,
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT=' ';

DROP TABLE IF EXISTS `form_template`;
CREATE TABLE IF NOT EXISTS `form_template` (
  `id` varchar(64)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `description` varchar(512)  DEFAULT NULL,
  `created_by` varchar(255) DEFAULT NULL,
  `created_time` datetime,
  `updated_by` varchar(255) DEFAULT NULL,
  `updated_time` datetime,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT=' ';

DROP TABLE IF EXISTS `form_item_template`;
CREATE TABLE IF NOT EXISTS `form_item_template` (
  `id` varchar(64) NOT NULL,
  `form_template` varchar(64) DEFAULT NULL,
  `name` varchar(64) NOT NULL,
  `description` varchar(255)  DEFAULT NULL,
  `item_group` varchar(255)  DEFAULT NULL,
  `item_group_name` varchar(255)  DEFAULT NULL,
  `default_value` varchar(255) DEFAULT NULL,
  `sort` int NOT NULL DEFAULT '0',
  `package_name` varchar(255) DEFAULT NULL,
  `entity` varchar(255) DEFAULT NULL,
  `attr_def_id` varchar(255) DEFAULT NULL,
  `attr_def_name` varchar(255) DEFAULT NULL,
  `attr_def_data_type` varchar(255) DEFAULT NULL,
  `element_type` varchar(64) NOT NULL DEFAULT 'text',
  `title` varchar(64) DEFAULT NULL,
  `width` int(11) DEFAULT 80,
  `ref_package_name` varchar(255) DEFAULT NULL,
  `ref_entity` varchar(255) DEFAULT NULL,
  `data_options` text,
  `required` varchar(16) NOT NULL DEFAULT 'no',
  `regular` varchar(255) DEFAULT NULL,
  `is_edit` varchar(16) NOT NULL DEFAULT 'yes',
  `is_view` varchar(16) NOT NULL DEFAULT 'yes',
  `is_output` varchar(16) NOT NULL DEFAULT 'no',
  `in_display_name` varchar(16) NOT NULL DEFAULT 'no',
  `is_ref_inside` varchar(16) NOT NULL DEFAULT 'no',
  `multiple` varchar(16) NOT NULL DEFAULT 'N',
  PRIMARY KEY (`id`),
  CONSTRAINT `fore_form_item_template` FOREIGN KEY (`form_template`) REFERENCES `form_template` (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT=' ';

DROP TABLE IF EXISTS `form`;
CREATE TABLE IF NOT EXISTS `form` (
  `id` varchar(64)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `description` varchar(255)  DEFAULT NULL,
  `form_template` varchar(64)  DEFAULT NULL,
  `created_time` datetime NULL,
  `created_by` varchar(255)   NULL,
  `updated_by` varchar(255)   NULL,
  `updated_time` datetime NULL,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  CONSTRAINT `fore_form_form_template` FOREIGN KEY (`form_template`) REFERENCES `form_template` (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT=' ';

DROP TABLE IF EXISTS `form_item`;
CREATE TABLE IF NOT EXISTS `form_item` (
  `id` varchar(64)  NOT NULL,
  `form` varchar(64)  NOT NULL,
  `form_item_template` varchar(64)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `value` varchar(255)  DEFAULT NULL,
  `item_group` varchar(255)  DEFAULT NULL,
  `row_data_id` varchar(255)  DEFAULT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fore_form_item_form` FOREIGN KEY (`form`) REFERENCES `form` (`id`),
  CONSTRAINT `fore_form_item_ref_template` FOREIGN KEY (`form_item_template`) REFERENCES `form_item_template` (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT=' ';

DROP TABLE IF EXISTS `request_template_group`;
CREATE TABLE IF NOT EXISTS `request_template_group` (
  `id` varchar(64)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `description` varchar(255)  DEFAULT NULL,
  `manage_role` varchar(64)  DEFAULT NULL,
  `created_by` varchar(255)   NULL,
  `created_time` datetime NULL,
  `updated_by` varchar(255)   NULL,
  `updated_time` datetime NULL,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  CONSTRAINT `fore_template_group_role` FOREIGN KEY (`manage_role`) REFERENCES `role` (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT=' ';

DROP TABLE IF EXISTS `request_template`;
CREATE TABLE IF NOT EXISTS `request_template` (
  `id` varchar(64)  NOT NULL,
  `group` varchar(64)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `description` varchar(255)  DEFAULT NULL,
  `form_template` varchar(64) DEFAULT NULL,
  `tags` varchar(255)  DEFAULT NULL,
  `record_id` varchar(64) DEFAULT NULL,
  `version` varchar(64) DEFAULT NULL,
  `confirm_time` varchar(64) DEFAULT NULL,
  `status` varchar(32) NOT NULL DEFAULT 'created',
  `package_name` varchar(255)  NULL,
  `entity_name` varchar(255)  NULL,
  `proc_def_key` varchar(255)  NOT NULL COMMENT 'key',
  `proc_def_id` varchar(255)  NOT NULL,
  `proc_def_name` varchar(255)  NOT NULL,
  `expire_day` int  DEFAULT 7,
  `created_by` varchar(255)   NULL,
  `created_time` datetime NULL,
  `updated_by` varchar(255)   NULL,
  `updated_time` datetime NULL,
  `entity_attrs` text DEFAULT NULL,
  `handler` varchar(64)  DEFAULT NULL,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  CONSTRAINT `fore_request_template_form` FOREIGN KEY (`form_template`) REFERENCES `form_template` (`id`),
  CONSTRAINT `fore_request_template_group` FOREIGN KEY (`group`) REFERENCES `request_template_group` (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT=' ';

DROP TABLE IF EXISTS `request`;
CREATE TABLE IF NOT EXISTS `request` (
  `id` varchar(64)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `form` varchar(64) DEFAULT NULL,
  `request_template` varchar(64)  NOT NULL,
  `proc_instance_id` varchar(64)  DEFAULT NULL,
  `proc_instance_key` varchar(64)  DEFAULT NULL,
  `reporter` varchar(64)  DEFAULT NULL,
  `handler` varchar(64)  DEFAULT NULL,
  `report_time` datetime DEFAULT NULL,
  `emergency` int  DEFAULT 3,
  `report_role` text  DEFAULT NULL,
  `attach_file` varchar(64)  DEFAULT NULL,
  `status` varchar(64)  NOT NULL DEFAULT 'created',
  `cache` mediumtext ,
  `bind_cache` mediumtext ,
  `result` mediumtext ,
  `expire_time` varchar(32)  DEFAULT NULL,
  `expect_time` varchar(32)  DEFAULT NULL,
  `confirm_time` varchar(32)  DEFAULT NULL,
  `created_by` varchar(255)   NULL,
  `created_time` datetime NULL,
  `updated_by` varchar(255)   NULL,
  `updated_time` datetime NULL,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  CONSTRAINT `fore_request_form` FOREIGN KEY (`form`) REFERENCES `form` (`id`),
  CONSTRAINT `fore_request_ref_template` FOREIGN KEY (`request_template`) REFERENCES `request_template` (`id`),
  CONSTRAINT `fore_request_attach_file` FOREIGN KEY (`attach_file`) REFERENCES `attach_file` (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT=' ';

DROP TABLE IF EXISTS `request_template_role`;
CREATE TABLE IF NOT EXISTS `request_template_role` (
  `id` varchar(160)  NOT NULL,
  `request_template` varchar(64) NOT NULL,
  `role` varchar(64) NOT NULL,
  `role_type` varchar(64) NOT NULL DEFAULT 'MGMT',
  PRIMARY KEY (`id`),
  CONSTRAINT `fore_request_template` FOREIGN KEY (`request_template`) REFERENCES `request_template` (`id`),
  CONSTRAINT `fore_request_template_role` FOREIGN KEY (`role`) REFERENCES `role` (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT=' ';

DROP TABLE IF EXISTS `task_template`;
CREATE TABLE IF NOT EXISTS `task_template` (
  `id` varchar(64)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `description` varchar(255)  DEFAULT NULL,
  `form_template` varchar(64) DEFAULT NULL,
  `request_template` varchar(64) DEFAULT NULL,
  `node_id` varchar(255)  DEFAULT NULL,
  `node_def_id` varchar(255)  DEFAULT NULL,
  `node_name` varchar(255)  DEFAULT NULL,
  `expire_day` int  DEFAULT 7,
  `handler` varchar(64)  DEFAULT NULL,
  `created_by` varchar(255)   NULL,
  `created_time` datetime NULL,
  `updated_by` varchar(255)   NULL,
  `updated_time` datetime NULL,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  CONSTRAINT `fore_task_template_form` FOREIGN KEY (`form_template`) REFERENCES `form_template` (`id`),
  CONSTRAINT `fore_task_template_request` FOREIGN KEY (`request_template`) REFERENCES `request_template` (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT=' ';

DROP TABLE IF EXISTS `task_template_role`;
CREATE TABLE IF NOT EXISTS `task_template_role` (
  `id` varchar(160)  NOT NULL,
  `task_template` varchar(64) NOT NULL,
  `role` varchar(64) NOT NULL,
  `role_type` varchar(64) NOT NULL DEFAULT 'MGMT',
  PRIMARY KEY (`id`),
  CONSTRAINT `fore_task_role_ref_template` FOREIGN KEY (`task_template`) REFERENCES `task_template` (`id`),
  CONSTRAINT `fore_task_template_ref_role` FOREIGN KEY (`role`) REFERENCES `role` (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT=' ';

DROP TABLE IF EXISTS `task`;
CREATE TABLE IF NOT EXISTS `task` (
  `id` varchar(64)  NOT NULL,
  `name` varchar(255)  NOT NULL,
  `description` varchar(255)  DEFAULT NULL,
  `form` varchar(64) DEFAULT NULL,
  `attach_file` varchar(64)  DEFAULT NULL,
  `status` varchar(64) DEFAULT 'created',
  `version` varchar(64)  DEFAULT NULL,
  `request` varchar(64)  DEFAULT NULL,
  `parent` varchar(64)  DEFAULT NULL,
  `task_template` varchar(64)  DEFAULT NULL,
  `package_name` varchar(255)  DEFAULT NULL,
  `entity_name` varchar(255)  DEFAULT NULL,
  `proc_def_id` varchar(255)  DEFAULT NULL,
  `proc_def_key` varchar(255)  DEFAULT NULL,
  `proc_def_name` varchar(255)  DEFAULT NULL,
  `node_def_id` varchar(64)  DEFAULT NULL,
  `node_name` varchar(64)  DEFAULT NULL,
  `callback_url` varchar(255)  DEFAULT NULL,
  `callback_parameter` varchar(64)  DEFAULT NULL,
  `callback_data` text  DEFAULT NULL,
  `reporter` varchar(64)  DEFAULT NULL,
  `report_time` datetime DEFAULT NULL,
  `report_role` varchar(255)  DEFAULT NULL,
  `emergency` int  DEFAULT 3,
  `result` text ,
  `cache` mediumtext ,
  `callback_request_id` varchar(255) DEFAULT NULL,
  `handler` varchar(64)  DEFAULT NULL,
  `next_option` varchar(255)  DEFAULT NULL,
  `chose_option` varchar(64)  DEFAULT NULL,
  `expire_time` varchar(32)  DEFAULT NULL,
  `created_by` varchar(255)   NULL,
  `created_time` datetime NULL,
  `updated_by` varchar(255)   NULL,
  `updated_time` datetime NULL,
  `del_flag` tinyint(2) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  CONSTRAINT `fore_task_form` FOREIGN KEY (`form`) REFERENCES `form` (`id`),
  CONSTRAINT `fore_task_ref_request` FOREIGN KEY (`request`) REFERENCES `request` (`id`),
  CONSTRAINT `fore_task_attach_file` FOREIGN KEY (`attach_file`) REFERENCES `attach_file` (`id`),
  CONSTRAINT `fore_task_ref_template` FOREIGN KEY (`task_template`) REFERENCES `task_template` (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT=' ';

DROP TABLE IF EXISTS `operation_log`;
CREATE TABLE IF NOT EXISTS `operation_log` (
  `id` varchar(64)  NOT NULL,
  `request` varchar(64) DEFAULT NULL,
  `task` varchar(64) DEFAULT NULL,
  `operation` varchar(64) NOT NULL,
  `operator` varchar(64) NOT NULL,
  `op_time` varchar(64) NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fore_operation_log_request` FOREIGN KEY (`request`) REFERENCES `request` (`id`),
  CONSTRAINT `fore_operation_log_task` FOREIGN KEY (`task`) REFERENCES `task` (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT=' ';

SET FOREIGN_KEY_CHECKS = 1;

#@v0.1.0.5-begin@;
alter table attach_file drop column s3_url;
alter table attach_file add column request varchar(64) default null;
alter table attach_file add column task varchar(64) default null;
#@v0.1.0.5-end@;

#@v0.1.0.21-begin@;
alter table task add column notify_count int default 0;
#@v0.1.0.21-end@;

#@v0.1.1.1-begin@;
alter table request add column parent varchar(64) default null;
#@v0.1.1.1-end@;

#@v0.1.2.15-begin@;
alter table form_item modify column value text default null;
#@v0.1.2.15-end@;

#@v0.1.2.16-begin@;
alter table request add column `rollback_desc` text  DEFAULT null;
#@v0.1.2.16-end@;

#@v0.1.3.13-begin@;
alter table request add column `type` tinyint  NOT NULL DEFAULT 0 COMMENT '模板类型:0表示请求 1表示发布';
alter table request add column operator_obj varchar(255) default null DEFAULT NULL COMMENT '操作对象';
alter table request add column description varchar(255) default null COMMENT '发布描述';
alter table request add column role varchar(255) default null COMMENT '创建请求的role';
alter table request add column revoke_flag tinyint(2) default 0 COMMENT '是否撤回,0表示否,1表示撤回';

alter table request_template add column `type` tinyint  NOT NULL DEFAULT 0 COMMENT '模板类型:0表示请求 1表示发布';
alter table request_template add column operator_obj_type varchar(255) default null DEFAULT NULL COMMENT '操作对象类型';
alter table request_template add column parent_id varchar(64) default null DEFAULT NULL COMMENT '父类ID';

alter table task add column template_type tinyint  NOT NULL DEFAULT 0 COMMENT '模板类型:0表示请求 1表示发布';


DROP TABLE IF EXISTS `collect_template`;
CREATE TABLE `collect_template` (
  `id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `request_template` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0' COMMENT '请求模板id',
  `type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '模板类型:0表示请求 1表示发布',
  `user` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户',
  `role` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '收藏模板使用角色',
  `created_time` datetime DEFAULT NULL,
  `updated_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `account` (`user`(191)),
  KEY `fore_collect_request_template` (`request_template`),
  CONSTRAINT `fore_collect_request_template` FOREIGN KEY (`request_template`) REFERENCES `request_template` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收藏模板';
#@v0.1.3.13-end@;

#@v0.1.3.14-begin@;
alter table form_item_template add column  default_clear  varchar(16) default 'no' not null COMMENT '是否清空';
#@v0.1.3.14-end@;

#@v0.1.3.45-begin@;
alter table operation_log drop FOREIGN KEY fore_operation_log_request;
alter table operation_log drop FOREIGN KEY fore_operation_log_task;
alter table operation_log add column `request_template` varchar(64) DEFAULT null;
alter table operation_log add column `request_template_name` varchar(255) DEFAULT null;
alter table operation_log add column `content` text DEFAULT null;
alter table operation_log add column `uri` varchar(255) DEFAULT null;
alter table operation_log add column `request_name` varchar(255) DEFAULT null;
alter table operation_log add column `task_name` varchar(255) DEFAULT null;

alter table request add index request_status (status);
alter table request add index request_created_by (created_by(191));
alter table request add index request_handler (handler);
alter table request add index request_expect_time (expect_time);
alter table request add index request_report_time (report_time);
alter table request add index request_confirm_time (confirm_time);
alter table request add index request_created_time (created_time);
alter table request add index request_updated_time (updated_time);
alter table request add index request_template_type (`type`);
alter table request add index request_del_flag (del_flag);

alter table task add index task_status (status);
alter table task add index task_del_flag (del_flag);
alter table task add index task_created_time (created_time);
alter table task add index task_updated_time (updated_time);
alter table task add index task_template_type (template_type);
#@v0.1.3.45-end@;

#@v1.0.5-begin@;
alter table request_template add column approve_by varchar(64) default null COMMENT '发布确认人';
alter table request_template add column check_switch tinyint default 0 COMMENT '是否加入定版流程';
alter table request_template add column confirm_switch tinyint default 0 COMMENT '是否加入确认流程';
alter table request_template add column back_desc  text default null COMMENT '退回理由';
alter table request_template modify column proc_def_key  varchar(255) default null COMMENT '编排key';
alter table request_template modify column proc_def_id  varchar(255) default null COMMENT '编排id';
alter table request_template modify column proc_def_name  varchar(255) default null COMMENT '编排名称';
alter table request_template add column proc_def_version  varchar(64) default null COMMENT '编排版本';

alter table request add column custom_form_cache text default null COMMENT '自定义表单cache';
alter table request add column notes text default null COMMENT '请求确认备注';
alter table request add column task_approval_cache text default null COMMENT '任务审批cache';
alter table request add column complete_status  varchar(64) default null COMMENT '任务节点完成状态:任务已完成 complete,未完成 uncompleted';

alter table task_template add column sort  int default 0 COMMENT '任务模板序号';
alter table task_template add column handle_mode  varchar(255) default null COMMENT '处理模式：custom.单人自定义 any.协同 all.并行 admin.提交人角色管理员 auto.自动通过';
alter table task_template add column type  varchar(64) default null COMMENT '任务类型: check 定版, approve 审批, implement 执行类型, confirm 请求确认';


alter table task add column type varchar(64) default null COMMENT '任务类型:submit 提交 check 定版, approve 审批, implement执行类型, confirm 请求确认 revoke 撤回';
alter table task add column sort int default '0' COMMENT '任务序号';
alter table task add column task_result  varchar(64) default null COMMENT '处理结果:approve同意,deny拒绝,redraw打回,complete完成,uncompleted未完成';
alter table task add column confirm_result varchar(64) default null COMMENT '任务确认结果:任务已完成 complete 未完成 uncompleted';
alter table task add column request_created_time datetime default null COMMENT '请求创建时间';

alter table form_item_template add column ref_id varchar(64) default null COMMENT '引用ID';


ALTER TABLE task_template_role DROP FOREIGN KEY fore_task_role_ref_template;
ALTER TABLE task_template_role DROP FOREIGN KEY fore_task_template_ref_role;
ALTER TABLE task_template_role DROP INDEX fore_task_role_ref_template;
ALTER TABLE task_template_role DROP INDEX fore_task_template_ref_role;


ALTER TABLE form_item DROP FOREIGN KEY fore_form_item_form;
ALTER TABLE request DROP FOREIGN KEY fore_request_form;
ALTER TABLE task DROP FOREIGN KEY fore_task_form;
ALTER TABLE form_item DROP INDEX fore_form_item_form;
ALTER TABLE request DROP INDEX fore_request_form;
ALTER TABLE task DROP INDEX fore_task_form;

ALTER TABLE form_item_template DROP FOREIGN KEY fore_form_item_template;
ALTER TABLE form DROP FOREIGN KEY fore_form_form_template;
ALTER TABLE request_template DROP FOREIGN KEY fore_request_template_form;
ALTER TABLE task_template DROP FOREIGN KEY fore_task_template_form;
ALTER TABLE form_item_template DROP INDEX fore_form_item_template;
ALTER TABLE form DROP INDEX fore_form_form_template;
ALTER TABLE request_template DROP INDEX fore_request_template_form;
ALTER TABLE task_template DROP INDEX fore_task_template_form;

alter table form_item_template change form_template form_template_old varchar(64) default null comment '表单模板(废弃)';
alter table form_item change form form_old varchar(64) default null comment '表单(废弃)';


alter table form rename to form_old;
alter table form_template rename to form_template_old;

CREATE TABLE IF NOT EXISTS  `form_template` (
  `id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `request_template` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '请求模板',
  `task_template` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '任务模板',
  `item_group` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '表单组',
  `item_group_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '表单组名',
  `item_group_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '表单组类型:workflow 编排数据,optional 自选,custom 自定义,request 请求信息',
  `item_group_rule` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '新增一行规则',
  `item_group_sort` tinyint(4) DEFAULT '1' COMMENT '排序',
  `ref_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '引用ID',
  `request_form_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '请求表单类型:message 信息表单 data 数据表单',
  `del_flag` tinyint(4) DEFAULT '0' COMMENT '是否删除',
  `created_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fore_request_template_new` (`request_template`),
  KEY `fore_task_template_new` (`task_template`),
  KEY `fore_ref_id_new` (`ref_id`),
  CONSTRAINT `fore_request_template_new` FOREIGN KEY (`request_template`) REFERENCES `request_template` (`id`),
  CONSTRAINT `fore_task_template_new` FOREIGN KEY (`task_template`) REFERENCES `task_template` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='表单模板表';

CREATE TABLE IF NOT EXISTS `form` (
  `id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `request` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '请求ID',
  `task` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '任务ID',
  `form_template` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '表单模板',
  `data_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '数据行ID',
  `created_by` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '更新人',
  `created_time` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `fore_form_request` (`request`),
  KEY `fore_form_task` (`task`),
  KEY `fore_form_template_new` (`form_template`),
  CONSTRAINT `fore_form_request` FOREIGN KEY (`request`) REFERENCES `request` (`id`),
  CONSTRAINT `fore_form_task` FOREIGN KEY (`task`) REFERENCES `task` (`id`),
  CONSTRAINT `fore_form_template_new` FOREIGN KEY (`form_template`) REFERENCES `form_template` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='表单表';


alter table form_item_template add column form_template varchar(64) default null COMMENT '表单模板';
alter table form_item add column form varchar(64) default null COMMENT '所属表单';


alter table form_item add constraint fore_form_item_form foreign key(form) REFERENCES form(id);
alter table form_item_template add constraint fore_form_item_template foreign key(form_template) REFERENCES form_template(id);
alter table form add constraint fore_form_form_template foreign key(form_template) REFERENCES form_template(id);


CREATE TABLE IF NOT EXISTS `task_handle_template` (
  `id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `task_template` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '任务模板',
  `role` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '角色',
  `assign` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '分派方式:template.模板指定 custom.提交人指定',
  `handler_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '人员设置方式：template.模板指定 template_suggest.模板建议 custom.提交人指定 custom_suggest.提交人建议 system.组内系统分配 claim.组内主动认领',
  `handler` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '处理人',
  `handle_mode` varchar(255) default null COMMENT '处理模式：custom.单人自定义 any.协同 all.并行 admin.提交人角色管理员 auto.自动通过',
  `sort` int  DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `force_task_template_new` (`task_template`),
  CONSTRAINT `force_task_template_new` FOREIGN KEY (`task_template`) REFERENCES `task_template` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务模版处理表';


CREATE TABLE IF NOT EXISTS `task_handle` (
    `id` varchar(64)  NOT NULL,
    `task_handle_template` varchar(64) DEFAULT NULL,
    `task` varchar(64)  DEFAULT NULL,
    `role` varchar(64)  DEFAULT NULL,
    `handler` varchar(64)  DEFAULT NULL,
    `handler_type` varchar(255)  DEFAULT NULL,
    `handle_result` varchar(64)  DEFAULT NULL,
    `handle_status` varchar(64)  DEFAULT 'uncompleted',
    `parent_id` varchar(64)  DEFAULT NULL,
    `created_time` datetime  NULL,
    `updated_time` datetime  NULL,
    `result_desc` text  DEFAULT NULL,
    `change_reason` varchar(64)  DEFAULT NULL,
    `sort` int  DEFAULT '0',
    `latest_flag` tinyint  DEFAULT '1',
    PRIMARY KEY (`id`),
    CONSTRAINT `fore_task_handle_template` FOREIGN KEY (`task_handle_template`) REFERENCES `task_handle_template` (`id`),
    CONSTRAINT `fore_task_handle_task` FOREIGN KEY (`task`) REFERENCES `task` (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务处理表';


alter table form_item add column request varchar(64) DEFAULT NULL COMMENT '请求ID';
alter table form_item add column updated_time datetime DEFAULT NULL COMMENT '更新时间';
alter table form_item add column original_id varchar(64) DEFAULT NULL COMMENT '原值id';
alter table form_item add column task_handle varchar(64) DEFAULT NULL COMMENT '任务处理ID';
alter table form_item add column del_falg tinyint(2) DEFAULT '0' COMMENT '删除标识';

alter table form_item add constraint fore_form_item_request foreign key(request) REFERENCES request(id);
alter table form_item add constraint fore_form_item_task_handle foreign key(task_handle) REFERENCES task_handle(id);

alter table attach_file add column task_handle varchar(64) default null COMMENT '任务处理';
alter table attach_file add constraint fore_attach_file_task_handle foreign key(task_handle) REFERENCES task_handle(id);
alter table form_item_template add column routine_expression text default null COMMENT '表单项计算表达式';

alter table task add index task_request_created_time(request_created_time);
alter table task add index task_type (type);
alter table task_handle add index task_handle_latest_flag(latest_flag);
alter table task_handle add index task_handle_created_time(created_time);
alter table task_handle add index task_handle_updated_time(updated_time);
alter table task_handle add index task_handle_result(handle_result);
alter table task_handle add index task_handle_handler(handler);

CREATE TABLE IF NOT EXISTS `task_notify` (
    `id` varchar(64)  NOT NULL,
    `task` varchar(64)  DEFAULT NULL,
    `doing_notify_count` TINYINT  DEFAULT 0,
    `timeout_notify_count` TINYINT  DEFAULT 0,
    `err_msg` text  DEFAULT NULL,
    `updated_time` datetime  NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `task_notify_task_Id` FOREIGN KEY (`task`) REFERENCES `task` (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务通知表';
alter table task_notify add index task_notify_updated_time(updated_time);
alter table request_template_role DROP FOREIGN KEY fore_request_template_role;
alter table request_template_group DROP FOREIGN KEY fore_template_group_role;
#@v1.0.5-end@;

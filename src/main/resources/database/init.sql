DROP TABLE IF EXISTS role_info;
CREATE TABLE role_info(
    id VARCHAR(32) NOT NULL   COMMENT '主键' ,
    name VARCHAR(255) NOT NULL   COMMENT '角色名' ,
    PRIMARY KEY (id)
) COMMENT = '角色信息表';

DROP TABLE IF EXISTS template_group;
CREATE TABLE template_group(
    id VARCHAR(32) NOT NULL   COMMENT '主键' ,
    manage_role VARCHAR(50)    COMMENT '所属角色' ,
    name VARCHAR(128)    COMMENT '名称' ,
    description VARCHAR(512)    COMMENT '描述' ,
    version VARCHAR(50)   DEFAULT 1 COMMENT '版本号' ,
    status TINYINT(2)   DEFAULT 0 COMMENT '状态' ,
    created_by VARCHAR(255) NOT NULL   COMMENT '创建人' ,
    created_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间' ,
    updated_by VARCHAR(255) NOT NULL   COMMENT '更新人' ,
    updated_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间' ,
    del_flag TINYINT(2) NOT NULL  DEFAULT 0 COMMENT '是否删除' ,
    PRIMARY KEY (id)
) COMMENT = '模板组信息表 ';

DROP TABLE IF EXISTS temp_group_role;
CREATE TABLE temp_group_role(
    id VARCHAR(32) NOT NULL   COMMENT '主键' ,
    group_id VARCHAR(32)    COMMENT '模板组编号' ,
    role_id VARCHAR(32)    COMMENT '角色编号' ,
    type TINYINT(2) NOT NULL  DEFAULT 0 COMMENT '类型' ,
    PRIMARY KEY (id)
) COMMENT = '模板角色关系表 ';

DROP TABLE IF EXISTS form_template;
CREATE TABLE form_template(
    id VARCHAR(32) NOT NULL   COMMENT '主键' ,
    name VARCHAR(50)    COMMENT '名称' ,
    description VARCHAR(50)    COMMENT '描述' ,
    style VARCHAR(50)    COMMENT '表单风格' ,
    created_by VARCHAR(255) NOT NULL   COMMENT '创建人' ,
    created_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间' ,
    updated_by VARCHAR(255) NOT NULL   COMMENT '更新人' ,
    updated_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间' ,
    del_flag TINYINT(2) NOT NULL  DEFAULT 0 COMMENT '是否删除' ,
    PRIMARY KEY (id)
) COMMENT = '表单模板信息表 ';

DROP TABLE IF EXISTS request_template;
CREATE TABLE request_template(
    id VARCHAR(32) NOT NULL   COMMENT '主键' ,
    deal_role VARCHAR(32)    COMMENT '所属角色' ,
    manage_role VARCHAR(32)    COMMENT '管理角色' ,
    group_id VARCHAR(32)    COMMENT '模板组编号' ,
    form_temp_id VARCHAR(32)    COMMENT '表单模板编号' ,
    proc_def_key VARCHAR(255)    COMMENT '流程编排key' ,
    name VARCHAR(128)    COMMENT '名称' ,
    version VARCHAR(50)    COMMENT '版本号' ,
    status TINYINT(2)   DEFAULT 0 COMMENT '状态' ,
    created_by VARCHAR(255) NOT NULL   COMMENT '创建人' ,
    created_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间' ,
    updated_by VARCHAR(255) NOT NULL   COMMENT '更新人' ,
    updated_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间' ,
    del_flag TINYINT(2) NOT NULL  DEFAULT 0 COMMENT '是否删除' ,
    PRIMARY KEY (id)
) COMMENT = '请求模板信息表 ';

DROP TABLE IF EXISTS task_template;
CREATE TABLE task_template(
    id VARCHAR(32) NOT NULL   COMMENT '主键' ,
    form_temp_id VARCHAR(32)    COMMENT '表单模板编号' ,
    request_temp_id VARCHAR(32)    COMMENT '请求模板编号' ,
    deal_role VARCHAR(128)    COMMENT '所属角色' ,
    proc_node_temp VARCHAR(255)    COMMENT '流程节点模板' ,
    name VARCHAR(255) NOT NULL   COMMENT '名称' ,
    created_by VARCHAR(255) NOT NULL   COMMENT '创建人' ,
    created_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间' ,
    updated_by VARCHAR(255) NOT NULL   COMMENT '更新人' ,
    updated_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间' ,
    del_flag TINYINT(2) NOT NULL  DEFAULT 0 COMMENT '是否删除' ,
    PRIMARY KEY (id)
) COMMENT = '任务模板信息表 ';

DROP TABLE IF EXISTS form_item_template;
CREATE TABLE form_item_template(
    id VARCHAR(32) NOT NULL   COMMENT '主键' ,
    temp_id VARCHAR(50)    COMMENT '模板id' ,
    temp_type VARCHAR(50)    COMMENT '模板类型' ,
    name VARCHAR(50)    COMMENT '名称' ,
    title VARCHAR(50)    COMMENT '标题' ,
    type VARCHAR(50)    COMMENT '类型' ,
    required TINYINT(2)    COMMENT '必填' ,
    正则 VARCHAR(50)    COMMENT 'evl' ,
    is_view VARCHAR(50)    COMMENT '是否显示' ,
    is_edit VARCHAR(50)    COMMENT '是否可编辑' ,
    width VARCHAR(50)    COMMENT '长度' ,
    def_value VARCHAR(50)    COMMENT '默认值' ,
    sort VARCHAR(50)    COMMENT '排序' ,
    created_by VARCHAR(255) NOT NULL   COMMENT '创建人' ,
    created_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间' ,
    updated_by VARCHAR(255) NOT NULL   COMMENT '更新人' ,
    updated_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间' ,
    del_flag TINYINT(2) NOT NULL  DEFAULT 0 COMMENT '是否删除' ,
    PRIMARY KEY (id)
) COMMENT = '表单项模板表 ';

DROP TABLE IF EXISTS currency_item_temolate;
CREATE TABLE currency_item_temolate(
    id VARCHAR(32) NOT NULL   COMMENT '主键' ,
    status TINYINT(2)   DEFAULT 0 COMMENT '状态' ,
    created_by VARCHAR(255) NOT NULL   COMMENT '创建人' ,
    created_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间' ,
    updated_by VARCHAR(255) NOT NULL   COMMENT '更新人' ,
    updated_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间' ,
    del_flag TINYINT(2) NOT NULL  DEFAULT 0 COMMENT '是否删除' ,
    PRIMARY KEY (id)
) COMMENT = '通用表单项模板表 ';

DROP TABLE IF EXISTS form_info;
CREATE TABLE form_info(
    status TINYINT(2)   DEFAULT 0 COMMENT '状态' ,
    created_by VARCHAR(255) NOT NULL   COMMENT '创建人' ,
    created_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间' ,
    updated_by VARCHAR(255) NOT NULL   COMMENT '更新人' ,
    updated_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间' ,
    del_flag TINYINT(2) NOT NULL  DEFAULT 0 COMMENT '是否删除' 
) COMMENT = '表单记录表 ';

DROP TABLE IF EXISTS form_item_info;
CREATE TABLE form_item_info(
    status TINYINT(2)   DEFAULT 0 COMMENT '状态' ,
    created_by VARCHAR(255) NOT NULL   COMMENT '创建人' ,
    created_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间' ,
    updated_by VARCHAR(255) NOT NULL   COMMENT '更新人' ,
    updated_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间' ,
    del_flag TINYINT(2) NOT NULL  DEFAULT 0 COMMENT '是否删除' 
) COMMENT = '表单项记录表 ';

DROP TABLE IF EXISTS request_info;
CREATE TABLE request_info(
    status TINYINT(2)   DEFAULT 0 COMMENT '状态' ,
    created_by VARCHAR(255) NOT NULL   COMMENT '创建人' ,
    created_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间' ,
    updated_by VARCHAR(255) NOT NULL   COMMENT '更新人' ,
    updated_time DATETIME NOT NULL  DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间' ,
    del_flag TINYINT(2) NOT NULL  DEFAULT 0 COMMENT '是否删除' 
) COMMENT = '请求记录表 ';

DROP TABLE IF EXISTS task_info;
CREATE TABLE task_info(
    id VARCHAR(32) NOT NULL   COMMENT '主键' ,
    STATUS TINYINT(2)   DEFAULT 0 COMMENT '状态' ,
    VERSION VARCHAR(50)    COMMENT '版本号' ,
    CREATED_BY VARCHAR(255)    COMMENT '创建人' ,
    CREATED_TIME DATETIME   DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间' ,
    UPDATED_BY VARCHAR(255)    COMMENT '更新人' ,
    UPDATED_TIME DATETIME   DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间' ,
    DEL_FLAG TINYINT(2)   DEFAULT 0 COMMENT '是否删除' ,
    PRIMARY KEY (id)
) COMMENT = '任务记录表 ';


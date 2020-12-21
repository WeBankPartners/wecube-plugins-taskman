package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

@ApiModel(value = "AddTaskTemplateReq",description = "add TaskTemplate req")
public class TaskTemplateReq extends QueryRoleRelationBaseReq {
    @ApiModelProperty(value = "任务模板名称",required = false,dataType = "String")
    private String name;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }
}

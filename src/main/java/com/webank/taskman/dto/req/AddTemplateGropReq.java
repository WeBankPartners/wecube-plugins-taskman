package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

@ApiModel(value = "AddTemplateGropReq",description = "add TemplateGroup req")
public class AddTemplateGropReq {
    @ApiModelProperty(value = "名称",required = true,dataType = "String")
    private String name;

    @ApiModelProperty(value = "管理角色id",required = true,dataType = "String")
    private String manageRoleId;

    @ApiModelProperty(value = "描述",dataType = "String")
    private String description;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getManageRoleId() {
        return manageRoleId;
    }

    public void setManageRoleId(String manageRoleId) {
        this.manageRoleId = manageRoleId;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }
}

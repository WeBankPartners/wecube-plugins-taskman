package com.webank.taskman.dto;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

@ApiModel(value = "TemplateGroup传入对象",description = "addReq")
public class TemplateGroupCreateVO {

    @ApiModelProperty(value = "主键",hidden = true,required = false)
    private String id;

    @ApiModelProperty(value = "名称",required = true,dataType = "String")
    private String name;

    @ApiModelProperty(value = "管理角色id",required = true,dataType = "String")
    private String manageRoleId;

    @ApiModelProperty(value = "描述",dataType = "String")
    private String description;

    @ApiModelProperty(value = "创建人",hidden = true)
    private String createdBy;

    @ApiModelProperty(value = "更新人",hidden = true)
    private String updatedBy;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getManageRoleId() {
        return manageRoleId;
    }

    public void setManageRoleId(String manageRoleId) {
        this.manageRoleId = manageRoleId;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getCreatedBy() {
        return createdBy;
    }

    public void setCreatedBy(String createdBy) {
        this.createdBy = createdBy;
    }

    public String getUpdatedBy() {
        return updatedBy;
    }

    public void setUpdatedBy(String updatedBy) {
        this.updatedBy = updatedBy;
    }

    @Override
    public String toString() {
        return "TemplateGroupVO{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", description='" + description + '\'' +
                ", createdBy='" + createdBy + '\'' +
                ", updatedBy='" + updatedBy + '\'' +
                '}';
    }
}

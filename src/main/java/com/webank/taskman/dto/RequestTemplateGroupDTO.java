package com.webank.taskman.dto;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.util.StringJoiner;

@ApiModel(value = "TemplateGroup respone object", description = "TemplateGroupDTO")
public class RequestTemplateGroupDTO {

    @ApiModelProperty(value = "主键", position = 1)
    private String id;

    @ApiModelProperty(value = "名称", position = 2)
    private String name;

    @ApiModelProperty(value = "描述", position = 3)
    private String description;

    @ApiModelProperty(value = "版本号", position = 4)
    private String version;

    @ApiModelProperty(value = "状态", position = 5)
    private Integer status;

    @ApiModelProperty(value = "管理角色", position = 6)
    private String manageRoleId;

    @ApiModelProperty(value = "管理角色姓名", position = 6)
    private String manageRoleName;

    public String getId() {
        return id;
    }

    public RequestTemplateGroupDTO setId(String id) {
        this.id = id;
        return this;
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

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }

    public Integer getStatus() {
        return status;
    }

    public void setStatus(Integer status) {
        this.status = status;
    }

    public String getManageRoleName() {
        return manageRoleName;
    }

    public void setManageRoleName(String manageRoleName) {
        this.manageRoleName = manageRoleName;
    }

    public String getManageRoleId() {
        return manageRoleId;
    }

    public void setManageRoleId(String manageRoleId) {
        this.manageRoleId = manageRoleId;
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", RequestTemplateGroupDTO.class.getSimpleName() + "[", "]")
                .add("id='" + id + "'")
                .add("name='" + name + "'")
                .add("description='" + description + "'")
                .add("version='" + version + "'")
                .add("status=" + status)
                .add("manageRoleId='" + manageRoleId + "'")
                .add("manageRoleName='" + manageRoleName + "'")
                .toString();
    }
}

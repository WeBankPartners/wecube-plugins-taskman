package com.webank.taskman.dto.req;


import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import javax.validation.constraints.NotBlank;

@ApiModel(value = "SaveRequestTemplateGropReq",description = "" )
public class SaveRequestTemplateGropReq {

    @ApiModelProperty(value = "模板组ID",dataType = "String",position = 1)
    private String id;

    @NotBlank(message = "名称不能为空")
    @ApiModelProperty(value = "模板组名称",dataType = "String",position = 2)
    private String name;

    @NotBlank(message = "管理角色不能为空")
    @ApiModelProperty(value = "管理角色id",dataType = "String",position = 3)
    private String manageRoleId;

    @NotBlank(message = "角色姓名不能为空")
    @ApiModelProperty(value = "管理角色姓名",dataType = "String",position = 4)
    private String manageRoleName;

    @ApiModelProperty(value = "描述",dataType = "String",position =5)
    private String description;
    @ApiModelProperty(value = "版本号",dataType = "String",position =7)
    private String version;


    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

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

    public String getManageRoleName() {
        return manageRoleName;
    }

    public void setManageRoleName(String manageRoleName) {
        this.manageRoleName = manageRoleName;
    }

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }

    @Override
    public String toString() {
        return "SaveAndUpdateTemplateGropReq{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", manageRoleId='" + manageRoleId + '\'' +
                ", description='" + description + '\'' +
                ", manageRoleName='" + manageRoleName + '\'' +
                '}';
    }
}

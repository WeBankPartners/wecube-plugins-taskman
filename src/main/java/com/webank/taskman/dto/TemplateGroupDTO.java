package com.webank.taskman.dto;

import com.webank.taskman.domain.RoleInfo;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

@ApiModel(value = "TemplateGroup respone object", description = "TemplateGroupDTO")
public class TemplateGroupDTO {

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
    private RoleInfo manageRole;


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

    public RoleInfo getManageRole() {
        return manageRole;
    }

    public void setManageRole(RoleInfo manageRole) {
        this.manageRole = manageRole;
    }
}

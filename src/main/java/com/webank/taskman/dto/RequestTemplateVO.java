package com.webank.taskman.dto;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.util.Date;
@ApiModel(value = "RequestTemplate对象",description = "RequestTemplate")
public class RequestTemplateVO {

    @ApiModelProperty(value = "主键")
    private String id;

    @ApiModelProperty(value = "所属角色")
    private String dealRole;

    @ApiModelProperty(value = "管理角色")
    private String manageRole;

    @ApiModelProperty(value = "模板组编号")
    private String groupId;

    @ApiModelProperty(value = "表单模板编号")
    private String formTempId;

    @ApiModelProperty(value = "流程编排key")
    private String procDefKey;

    @ApiModelProperty(value = "名称")
    private String name;

    @ApiModelProperty(value = "版本号")
    private String version;

    @ApiModelProperty(value = "状态")
    private Integer status;

    @ApiModelProperty(value = "创建人")
    private String createdBy;

    @ApiModelProperty(value = "更新人")
    private String updatedBy;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getDealRole() {
        return dealRole;
    }

    public void setDealRole(String dealRole) {
        this.dealRole = dealRole;
    }

    public String getManageRole() {
        return manageRole;
    }

    public void setManageRole(String manageRole) {
        this.manageRole = manageRole;
    }

    public String getGroupId() {
        return groupId;
    }

    public void setGroupId(String groupId) {
        this.groupId = groupId;
    }

    public String getFormTempId() {
        return formTempId;
    }

    public void setFormTempId(String formTempId) {
        this.formTempId = formTempId;
    }

    public String getProcDefKey() {
        return procDefKey;
    }

    public void setProcDefKey(String procDefKey) {
        this.procDefKey = procDefKey;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
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
        return "RequestTemplateVO{" +
                "id='" + id + '\'' +
                ", dealRole='" + dealRole + '\'' +
                ", manageRole='" + manageRole + '\'' +
                ", groupId='" + groupId + '\'' +
                ", formTempId='" + formTempId + '\'' +
                ", procDefKey='" + procDefKey + '\'' +
                ", name='" + name + '\'' +
                ", version='" + version + '\'' +
                ", status=" + status +
                ", createdBy='" + createdBy + '\'' +
                ", updatedBy='" + updatedBy + '\'' +
                '}';
    }
}

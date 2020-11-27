package com.webank.taskman.dto;

import java.util.Date;

public class RequestTemplateVO {

    private String id;

    private String dealRole;

    private String manageRole;

    private String groupId;

    private String formTempId;

    private String procDefKey;

    private String name;

    private String version;

    private Integer status;

    private String createdBy;

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

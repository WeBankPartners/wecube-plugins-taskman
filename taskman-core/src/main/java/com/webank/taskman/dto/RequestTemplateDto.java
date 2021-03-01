package com.webank.taskman.dto;

import java.util.ArrayList;
import java.util.List;

public class RequestTemplateDto {

    private String id;

    private String requestTempGroup;

    private String requestTempGroupName;

    private String procDefKey;

    private String procDefId;

    private String procDefName;

    private String name;

    private String description;
    private String version;

    private String tags;

    private String status;

    private List<RoleDto> useRoles = new ArrayList<>();

    private List<RoleDto> manageRoles = new ArrayList<>();

    public String getId() {
        return id;
    }

    public RequestTemplateDto setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestTempGroup() {
        return requestTempGroup;
    }

    public RequestTemplateDto setRequestTempGroup(String requestTempGroup) {
        this.requestTempGroup = requestTempGroup;
        return this;
    }

    public String getRequestTempGroupName() {
        return requestTempGroupName;
    }

    public void setRequestTempGroupName(String requestTempGroupName) {
        this.requestTempGroupName = requestTempGroupName;
    }

    public String getProcDefKey() {
        return procDefKey;
    }

    public RequestTemplateDto setProcDefKey(String procDefKey) {
        this.procDefKey = procDefKey;
        return this;
    }

    public String getProcDefId() {
        return procDefId;
    }

    public RequestTemplateDto setProcDefId(String procDefId) {
        this.procDefId = procDefId;
        return this;
    }

    public String getProcDefName() {
        return procDefName;
    }

    public RequestTemplateDto setProcDefName(String procDefName) {
        this.procDefName = procDefName;
        return this;
    }

    public String getName() {
        return name;
    }

    public RequestTemplateDto setName(String name) {
        this.name = name;
        return this;
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

    public RequestTemplateDto setVersion(String version) {
        this.version = version;
        return this;
    }

    public String getTags() {
        return tags;
    }

    public RequestTemplateDto setTags(String tags) {
        this.tags = tags;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public RequestTemplateDto setStatus(String status) {
        this.status = status;
        return this;
    }

    public List<RoleDto> getUseRoles() {
        return useRoles;
    }

    public void setUseRoles(List<RoleDto> useRoles) {
        this.useRoles = useRoles;
    }

    public List<RoleDto> getManageRoles() {
        return manageRoles;
    }

    public void setManageRoles(List<RoleDto> manageRoles) {
        this.manageRoles = manageRoles;
    }

    @Override
    public String toString() {
        return "RequestTemplateDTO{" + "id='" + id + '\'' + ", requestTempGroup='" + requestTempGroup + '\''
                + ", requestTempGroupName='" + requestTempGroupName + '\'' + ", procDefKey='" + procDefKey + '\''
                + ", procDefId='" + procDefId + '\'' + ", procDefName='" + procDefName + '\'' + ", name='" + name + '\''
                + ", description='" + description + '\'' + ", version='" + version + '\'' + ", tags='" + tags + '\''
                + ", status='" + status + '\'' + ", useRoles=" + useRoles + ", manageRoles=" + manageRoles + '}';
    }
}

package com.webank.taskman.dto.req;

import java.util.List;

import com.webank.taskman.dto.RoleDto;

public class RequestTemplateSaveReqDto {
    private String id;
    private String requestTempGroup;
    private String procDefId;
    private String procDefKey;
    private String procDefName;
    private String name;
    private String description;
    private String tags;
    private List<RoleDto> useRoles;
    private List<RoleDto> manageRoles;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getRequestTempGroup() {
        return requestTempGroup;
    }

    public void setRequestTempGroup(String requestTempGroup) {
        this.requestTempGroup = requestTempGroup;
    }

    public String getProcDefId() {
        return procDefId;
    }

    public void setProcDefId(String procDefId) {
        this.procDefId = procDefId;
    }

    public String getProcDefKey() {
        return procDefKey;
    }

    public void setProcDefKey(String procDefKey) {
        this.procDefKey = procDefKey;
    }

    public String getProcDefName() {
        return procDefName;
    }

    public void setProcDefName(String procDefName) {
        this.procDefName = procDefName;
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

    public String getTags() {
        return tags;
    }

    public void setTags(String tags) {
        this.tags = tags;
    }

}

package com.webank.taskman.support.core.dto;

import com.fasterxml.jackson.annotation.JsonAlias;

public class RolesDataResponse {

    @JsonAlias("id")
    private String roleId;
    @JsonAlias("name")
    private String roleName;
    @JsonAlias("displayName")
    private String description;

    public RolesDataResponse() {
    }

    public RolesDataResponse(String roleId, String roleName, String description) {
        this.roleId = roleId;
        this.roleName = roleName;
        this.description = description;
    }

    public String getRoleId() {
        return roleId;
    }

    public void setRoleId(String roleId) {
        this.roleId = roleId;
    }

    public String getRoleName() {
        return roleName;
    }

    public void setRoleName(String roleName) {
        this.roleName = roleName;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    @Override
    public String toString() {
        return "RolesDataResponse [roleId=" + roleId + ", roleName=" + roleName + ", description=" + description + "]";
    }
}

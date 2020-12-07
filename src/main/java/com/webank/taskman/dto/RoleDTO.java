package com.webank.taskman.dto;

import com.fasterxml.jackson.annotation.JsonAlias;

public class RoleDTO {

    private String roleName;
    private String displayName;


    public String getRoleName() {
        return roleName;
    }

    public void setRoleName(String roleName) {
        this.roleName = roleName;
    }

    public String getDisplayName() {
        return displayName;
    }

    public void setDisplayName(String displayName) {
        this.displayName = displayName;
    }
}

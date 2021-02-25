package com.webank.taskman.dto;

public class RoleDto {

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

    @Override
    public String toString() {
        return "RoleDTO{" +
                "roleName='" + roleName + '\'' +
                ", displayName='" + displayName + '\'' +
                '}';
    }
}

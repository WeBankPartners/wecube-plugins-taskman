package com.webank.taskman.dto;

import java.util.StringJoiner;

public class RequestTemplateGroupDTO {

    private String id;

    private String name;

    private String description;

    private String version;

    private String status;

    private String manageRoleId;

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

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
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
        return new StringJoiner(", ", RequestTemplateGroupDTO.class.getSimpleName() + "[", "]").add("id='" + id + "'")
                .add("name='" + name + "'").add("description='" + description + "'").add("version='" + version + "'")
                .add("status=" + status).add("manageRoleId='" + manageRoleId + "'")
                .add("manageRoleName='" + manageRoleName + "'").toString();
    }
}

package com.webank.taskman.dto.req;

public class SaveRequestTemplateGropReq {

    private String id;
    private String name;
    private String manageRoleId;
    private String manageRoleName;
    private String description;
    private String version;
    private String Status;

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

    public String getStatus() {
        return Status;
    }

    public SaveRequestTemplateGropReq setStatus(String status) {
        Status = status;
        return this;
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

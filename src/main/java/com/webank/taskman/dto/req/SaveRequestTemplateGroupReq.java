package com.webank.taskman.dto.req;

public class SaveRequestTemplateGroupReq {

    private String id;
    private String name;
    private String manageRoleId;
    private String manageRoleName;
    private String description;
    private String version;
    private String status;

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
        return status;
    }

    public SaveRequestTemplateGroupReq setStatus(String status) {
        this.status = status;
        return this;
    }

    @Override
    public String toString() {
        StringBuilder builder = new StringBuilder();
        builder.append("SaveRequestTemplateGroupReq [id=");
        builder.append(id);
        builder.append(", name=");
        builder.append(name);
        builder.append(", manageRoleId=");
        builder.append(manageRoleId);
        builder.append(", manageRoleName=");
        builder.append(manageRoleName);
        builder.append(", description=");
        builder.append(description);
        builder.append(", version=");
        builder.append(version);
        builder.append(", status=");
        builder.append(status);
        builder.append("]");
        return builder.toString();
    }

}

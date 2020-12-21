package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.webank.taskman.base.BaseEntity;

import java.io.Serializable;

public class RequestTemplateGroup extends BaseEntity implements Serializable {

    private static final long serialVersionUID = 1L;

    
    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    
    private String manageRoleId;

    
    private String manageRoleName;

    
    private String name;

    private String description;

    private String version;

    private String status;

    public String getId() {
        return id;
    }

    public RequestTemplateGroup setId(String id) {
        this.id = id;
        return this;
    }

    public String getManageRoleId() {
        return manageRoleId;
    }

    public void setManageRoleId(String manageRoleId) {
        this.manageRoleId = manageRoleId;
    }

    public String getManageRoleName() {
        return manageRoleName;
    }

    public void setManageRoleName(String manageRoleName) {
        this.manageRoleName = manageRoleName;
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

    @Override
    public String toString() {
        return "RequestTemplateGroup{" +
        "id=" + id +
        ", manageRoleId=" + manageRoleId +
        ", manageRoleName=" + manageRoleName +
        ", name=" + name +
        ", description=" + description +
        ", version=" + version +
        ", status=" + status +
        ", createdBy=" + getCreatedBy() +
        ", createdTime=" + getCreatedTime() +
        ", updatedBy=" + getUpdatedBy() +
        ", updatedTime=" + getUpdatedTime() +
        ", delFlag=" + getDelFlag() +
        "}";
    }

    public  enum  RequestTemplateGroupStatus{

    }
}

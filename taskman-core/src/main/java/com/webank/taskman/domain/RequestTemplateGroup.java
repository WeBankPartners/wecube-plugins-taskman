package com.webank.taskman.domain;

import java.io.Serializable;

import org.springframework.util.StringUtils;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.core.conditions.update.UpdateWrapper;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.webank.taskman.base.BaseEntity;

public class RequestTemplateGroup extends BaseEntity implements Serializable {

    private static final long serialVersionUID = 1L;
    
    
    public static final String STATUS_AVAILABLE = "1";
    public static final String STATUS_NOT_AVAILABLE = "0";

    
    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;
    private String manageRoleId;
    private String manageRoleName;
    private String name;
    private String description;
    private String version;
    private String status;

    public RequestTemplateGroup() {
    }

    public RequestTemplateGroup(String id, String manageRoleId, String manageRoleName, String name, String description, String version, String status) {
        this.id = id;
        this.manageRoleId = manageRoleId;
        this.manageRoleName = manageRoleName;
        this.name = name;
        this.description = description;
        this.version = version;
        this.status = status;
    }

    @JsonIgnore
    public LambdaQueryWrapper<RequestTemplateGroup> getLambdaQueryWrapper() {
        return new LambdaQueryWrapper<RequestTemplateGroup>()
            .eq(!StringUtils.isEmpty(id), RequestTemplateGroup::getId, id)
            .eq(!StringUtils.isEmpty(manageRoleId), RequestTemplateGroup::getManageRoleId, manageRoleId)
            .eq(!StringUtils.isEmpty(manageRoleName), RequestTemplateGroup::getManageRoleName, manageRoleName)
            .like(!StringUtils.isEmpty(name), RequestTemplateGroup::getName, name)
            .like(!StringUtils.isEmpty(description), RequestTemplateGroup::getDescription, description)
            .eq(!StringUtils.isEmpty(version), RequestTemplateGroup::getVersion, version)
            .eq(!StringUtils.isEmpty(status), RequestTemplateGroup::getStatus, status)
        ;
    }

    @JsonIgnore
    public UpdateWrapper<RequestTemplateGroup> getUpdateWrapper() {
        UpdateWrapper<RequestTemplateGroup> wrapper = new UpdateWrapper<>();
        wrapper.lambda().eq(!StringUtils.isEmpty(id), RequestTemplateGroup::getId, id)
                .set(!StringUtils.isEmpty(id), RequestTemplateGroup::getId, id)
                .set(!StringUtils.isEmpty(manageRoleId), RequestTemplateGroup::getManageRoleId, manageRoleId)
                .set(!StringUtils.isEmpty(manageRoleName), RequestTemplateGroup::getManageRoleName, manageRoleName)
                .set(!StringUtils.isEmpty(name), RequestTemplateGroup::getName, name)
                .set(!StringUtils.isEmpty(description), RequestTemplateGroup::getDescription, description)
                .set(!StringUtils.isEmpty(version), RequestTemplateGroup::getVersion, version)
                .set(!StringUtils.isEmpty(status), RequestTemplateGroup::getStatus, status)
        ;
        return wrapper;
    }

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

    public RequestTemplateGroup setManageRoleId(String manageRoleId) {
        this.manageRoleId = manageRoleId;
        return this;
    }

    public String getManageRoleName() {
        return manageRoleName;
    }

    public RequestTemplateGroup setManageRoleName(String manageRoleName) {
        this.manageRoleName = manageRoleName;
        return this;
    }

    public String getName() {
        return name;
    }

    public RequestTemplateGroup setName(String name) {
        this.name = name;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public RequestTemplateGroup setDescription(String description) {
        this.description = description;
        return this;
    }

    public String getVersion() {
        return version;
    }

    public RequestTemplateGroup setVersion(String version) {
        this.version = version;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public RequestTemplateGroup setStatus(String status) {
        this.status = status;
        return this;
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
}
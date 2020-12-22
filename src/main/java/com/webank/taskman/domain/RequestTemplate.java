package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.webank.taskman.base.BaseEntity;
import org.springframework.util.StringUtils;

import java.io.Serializable;

public class RequestTemplate extends BaseEntity implements Serializable {

    private static final long serialVersionUID = 1L;


    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;
    private String requestTempGroup;
    private String procDefKey;
    private String procDefId;
    private String procDefName;
    private String packageName;
    private String entityName;
    private String name;
    private String description;
    private String version;
    private String tags;
    private String status;

    public RequestTemplate() {
    }

    public RequestTemplate(String id, String requestTempGroup,
                           String procDefKey, String procDefId, String name,
                           String description, String version, String tags, String status) {
        this.id = id;
        this.requestTempGroup = requestTempGroup;
        this.procDefKey = procDefKey;
        this.procDefId = procDefId;
        this.name = name;
        this.description = description;
        this.version = version;
        this.tags = tags;
        this.status = status;
    }

    @JsonIgnore
    public LambdaQueryWrapper getLambdaQueryWrapper() {
        return new LambdaQueryWrapper<RequestTemplate>()
                .eq(!StringUtils.isEmpty(id), RequestTemplate::getId, id)
                .eq(!StringUtils.isEmpty(requestTempGroup), RequestTemplate::getRequestTempGroup, requestTempGroup)
                .eq(!StringUtils.isEmpty(procDefKey), RequestTemplate::getProcDefKey, procDefKey)
                .eq(!StringUtils.isEmpty(procDefId), RequestTemplate::getProcDefId, procDefId)
                .eq(!StringUtils.isEmpty(procDefName), RequestTemplate::getProcDefName, procDefName)
                .eq(!StringUtils.isEmpty(packageName), RequestTemplate::getPackageName, packageName)
                .eq(!StringUtils.isEmpty(entityName), RequestTemplate::getEntityName, entityName)
                .like(!StringUtils.isEmpty(name), RequestTemplate::getName, name)
                .eq(!StringUtils.isEmpty(description), RequestTemplate::getDescription, description)
                .eq(!StringUtils.isEmpty(version), RequestTemplate::getVersion, version)
                .eq(!StringUtils.isEmpty(tags), RequestTemplate::getTags, tags)
                .eq(!StringUtils.isEmpty(status), RequestTemplate::getStatus, status)
                ;
    }

    public String getId() {
        return id;
    }

    public RequestTemplate setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestTempGroup() {
        return requestTempGroup;
    }

    public RequestTemplate setRequestTempGroup(String requestTempGroup) {
        this.requestTempGroup = requestTempGroup;
        return this;
    }

    public String getProcDefKey() {
        return procDefKey;
    }

    public RequestTemplate setProcDefKey(String procDefKey) {
        this.procDefKey = procDefKey;
        return this;
    }

    public String getProcDefId() {
        return procDefId;
    }

    public RequestTemplate setProcDefId(String procDefId) {
        this.procDefId = procDefId;
        return this;
    }

    public String getProcDefName() {
        return procDefName;
    }

    public RequestTemplate setProcDefName(String procDefName) {
        this.procDefName = procDefName;
        return this;
    }

    public String getPackageName() {
        return packageName;
    }

    public RequestTemplate setPackageName(String packageName) {
        this.packageName = packageName;
        return this;
    }

    public String getEntityName() {
        return entityName;
    }

    public RequestTemplate setEntityName(String entityName) {
        this.entityName = entityName;
        return this;
    }

    public String getName() {
        return name;
    }

    public RequestTemplate setName(String name) {
        this.name = name;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public RequestTemplate setDescription(String description) {
        this.description = description;
        return this;
    }

    public String getVersion() {
        return version;
    }

    public RequestTemplate setVersion(String version) {
        this.version = version;
        return this;
    }

    public String getTags() {
        return tags;
    }

    public RequestTemplate setTags(String tags) {
        this.tags = tags;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public RequestTemplate setStatus(String status) {
        this.status = status;
        return this;
    }

    @Override
    public String toString() {
        return "RequestTemplate{" +
        "id=" + id +
        ", requestTempGroup=" + requestTempGroup +
        ", procDefKey=" + procDefKey +
        ", procDefId=" + procDefId +
        ", procDefName=" + procDefName +
        ", name=" + name +
        ", version=" + version +
        ", tags=" + tags +
        ", status=" + status +
        ", createdBy=" + getCreatedBy() +
        ", createdTime=" + getCreatedTime() +
        ", updatedBy=" + getUpdatedBy() +
        ", updatedTime=" + getUpdatedTime() +
        ", delFlag=" + getDelFlag() +
        "}";
    }
}

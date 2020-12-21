package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.webank.taskman.base.BaseEntity;

import java.io.Serializable;

public class RequestTemplate extends BaseEntity implements Serializable {

    private static final long serialVersionUID = 1L;


    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    private String requestTempGroup;


    private String procDefKey;


    private String procDefId;


    private String procDefName;


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

    public void setRequestTempGroup(String requestTempGroup) {
        this.requestTempGroup = requestTempGroup;
    }

    public String getProcDefKey() {
        return procDefKey;
    }

    public void setProcDefKey(String procDefKey) {
        this.procDefKey = procDefKey;
    }

    public String getProcDefId() {
        return procDefId;
    }

    public void setProcDefId(String procDefId) {
        this.procDefId = procDefId;
    }

    public String getProcDefName() {
        return procDefName;
    }

    public void setProcDefName(String procDefName) {
        this.procDefName = procDefName;
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

    public void setDescription(String description) {
        this.description = description;
    }

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }

    public String getTags() {
        return tags;
    }

    public void setTags(String tags) {
        this.tags = tags;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }


    public static QueryWrapper<RequestTemplate> getQueryWrapper(String id,String requestTempGroup,String procDefKey,String procDefId,String name,
                        String description,String version,String tags,String status){
        QueryWrapper<RequestTemplate> queryWrapper = new QueryWrapper<>();
        queryWrapper.setEntity(new RequestTemplate(id,requestTempGroup,procDefKey,procDefId,name,description,version,tags,status));
        return queryWrapper;
    }

    public static QueryWrapper<RequestTemplate> getQueryWrapper(String name){
        QueryWrapper<RequestTemplate> queryWrapper = new QueryWrapper<>();
        RequestTemplate template = new RequestTemplate();
        template.setName(name);
        return getQueryWrapper(null,null,null,null,name,null,null,null,null);
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

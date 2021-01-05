package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.webank.taskman.base.BaseEntity;
import org.springframework.util.StringUtils;

import java.io.Serializable;

public class TaskTemplate extends BaseEntity implements Serializable {

    private static final long serialVersionUID = 1L;


    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    private String requestTemplateId;

    private String procDefId;

    private String procDefKey;

    private String procDefName;

    private String name;

    private String nodeDefId;

    private String nodeName;

    private String description;

    public TaskTemplate() {
    }

    public TaskTemplate(String id, String procDefId, String procDefKey, String procDefName, String name, String nodeDefId, String nodeName, String description) {
        this.id = id;
        this.procDefId = procDefId;
        this.procDefKey = procDefKey;
        this.procDefName = procDefName;
        this.name = name;
        this.nodeDefId = nodeDefId;
        this.nodeName = nodeName;
        this.description = description;
    }

    @JsonIgnore
    public LambdaQueryWrapper<TaskTemplate> getLambdaQueryWrapper() {
        return new LambdaQueryWrapper<TaskTemplate>()
            .eq(!StringUtils.isEmpty(id), TaskTemplate::getId, id)
            .eq(!StringUtils.isEmpty(requestTemplateId), TaskTemplate::getRequestTemplateId, requestTemplateId)
            .eq(!StringUtils.isEmpty(procDefId), TaskTemplate::getProcDefId, procDefId)
            .eq(!StringUtils.isEmpty(procDefKey), TaskTemplate::getProcDefKey, procDefKey)
            .eq(!StringUtils.isEmpty(procDefName), TaskTemplate::getProcDefName, procDefName)
            .like(!StringUtils.isEmpty(name), TaskTemplate::getName, name)
            .eq(!StringUtils.isEmpty(nodeDefId), TaskTemplate::getNodeDefId, nodeDefId)
            .eq(!StringUtils.isEmpty(nodeName), TaskTemplate::getNodeName, nodeName)
            .eq(!StringUtils.isEmpty(description), TaskTemplate::getDescription, description)
        ;
    }

    public String getId() {
        return id;
    }

    public TaskTemplate setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestTemplateId() {
        return requestTemplateId;
    }

    public TaskTemplate setRequestTemplateId(String requestTemplateId) {
        this.requestTemplateId = requestTemplateId;
        return this;
    }

    public String getProcDefId() {
        return procDefId;
    }

    public TaskTemplate setProcDefId(String procDefId) {
        this.procDefId = procDefId;
        return this;
    }

    public String getProcDefKey() {
        return procDefKey;
    }

    public TaskTemplate setProcDefKey(String procDefKey) {
        this.procDefKey = procDefKey;
        return this;
    }

    public String getProcDefName() {
        return procDefName;
    }

    public TaskTemplate setProcDefName(String procDefName) {
        this.procDefName = procDefName;
        return this;
    }

    public String getName() {
        return name;
    }

    public TaskTemplate setName(String name) {
        this.name = name;
        return this;
    }

    public String getNodeDefId() {
        return nodeDefId;
    }

    public TaskTemplate setNodeDefId(String nodeDefId) {
        this.nodeDefId = nodeDefId;
        return this;
    }

    public String getNodeName() {
        return nodeName;
    }

    public TaskTemplate setNodeName(String nodeName) {
        this.nodeName = nodeName;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public TaskTemplate setDescription(String description) {
        this.description = description;
        return this;
    }

    @Override
    public String toString() {
        return "TaskTemplate{" +
                "id='" + id + '\'' +
                ", procDefId='" + procDefId + '\'' +
                ", procDefKey='" + procDefKey + '\'' +
                ", procDefName='" + procDefName + '\'' +
                ", name='" + name + '\'' +
                ", nodeDefId='" + nodeDefId + '\'' +
                ", nodeName='" + nodeName + '\'' +
                ", description='" + description + '\'' +
                '}';
    }


}

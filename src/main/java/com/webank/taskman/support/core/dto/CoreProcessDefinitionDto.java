package com.webank.taskman.support.core.dto;

import java.util.List;

public class CoreProcessDefinitionDto {

    private String procDefId;
    private String procDefKey;
    private String procDefName;
    private String procDefVersion;
    private String status;
    private String procDefData;
    private String rootEntity;
    private String createdTime;

    private Object permissionToRole;
    private List<CoreTaskNodeInfo> taskNodeInfos;

    public String getProcDefId() {
        return procDefId;
    }

    public void setProcDefId(String procDefId) {
        this.procDefId = procDefId;
    }

    public String getProcDefKey() {
        return procDefKey;
    }

    public void setProcDefKey(String procDefKey) {
        this.procDefKey = procDefKey;
    }

    public String getProcDefName() {
        return procDefName;
    }

    public void setProcDefName(String procDefName) {
        this.procDefName = procDefName;
    }

    public String getProcDefVersion() {
        return procDefVersion;
    }

    public void setProcDefVersion(String procDefVersion) {
        this.procDefVersion = procDefVersion;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public String getProcDefData() {
        return procDefData;
    }

    public void setProcDefData(String procDefData) {
        this.procDefData = procDefData;
    }

    public String getRootEntity() {
        return rootEntity;
    }

    public void setRootEntity(String rootEntity) {
        this.rootEntity = rootEntity;
    }

    public String getCreatedTime() {
        return createdTime;
    }

    public void setCreatedTime(String createdTime) {
        this.createdTime = createdTime;
    }

    public Object getPermissionToRole() {
        return permissionToRole;
    }

    public void setPermissionToRole(Object permissionToRole) {
        this.permissionToRole = permissionToRole;
    }

    public List<CoreTaskNodeInfo> getTaskNodeInfos() {
        return taskNodeInfos;
    }

    public void setTaskNodeInfos(List<CoreTaskNodeInfo> taskNodeInfos) {
        this.taskNodeInfos = taskNodeInfos;
    }

    public CoreProcessDefinitionDto() {
    }

    public CoreProcessDefinitionDto(String procDefId, String procDefKey, String procDefName, String procDefVersion, String status, String procDefData, String rootEntity, String createdTime) {
        this.procDefId = procDefId;
        this.procDefKey = procDefKey;
        this.procDefName = procDefName;
        this.procDefVersion = procDefVersion;
        this.status = status;
        this.procDefData = procDefData;
        this.rootEntity = rootEntity;
        this.createdTime = createdTime;
    }

    @Override
    public String toString() {
        return "CoreProcessDefinitionDto [procDefId=" + procDefId + ", procDefKey=" + procDefKey + ", procDefName="
                + procDefName + ", procDefVersion=" + procDefVersion + ", status=" + status + ", procDefData="
                + procDefData + ", rootEntity=" + rootEntity + ", createdTime=" + createdTime + "]";
    }

}

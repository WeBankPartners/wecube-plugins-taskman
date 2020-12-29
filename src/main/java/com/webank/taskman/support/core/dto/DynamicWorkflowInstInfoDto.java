package com.webank.taskman.support.core.dto;

public class DynamicWorkflowInstInfoDto {
    private String id;
    private String procInstId;
    private String procDefId;
    private String procDefKey;
    private String status;
    public String getId() {
        return id;
    }
    public void setId(String id) {
        this.id = id;
    }
    public String getProcInstId() {
        return procInstId;
    }
    public void setProcInstId(String procInstId) {
        this.procInstId = procInstId;
    }
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
    public String getStatus() {
        return status;
    }
    public void setStatus(String status) {
        this.status = status;
    }
    
    

}

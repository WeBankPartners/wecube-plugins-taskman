package com.webank.taskman.dto;

public class CoreCancelTaskDTO {

    private String procInstId;
    private String taskNodeId;

    public String getProcInstId() {
        return procInstId;
    }

    public CoreCancelTaskDTO setProcInstId(String procInstId) {
        this.procInstId = procInstId;
        return this;
    }

    public String getTaskNodeId() {
        return taskNodeId;
    }

    public CoreCancelTaskDTO setTaskNodeId(String taskNodeId) {
        this.taskNodeId = taskNodeId;
        return this;
    }
}

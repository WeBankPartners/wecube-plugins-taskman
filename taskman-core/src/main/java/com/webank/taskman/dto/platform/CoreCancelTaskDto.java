package com.webank.taskman.dto.platform;

public class CoreCancelTaskDto {

    private String procInstId;
    private String taskNodeId;

    public String getProcInstId() {
        return procInstId;
    }

    public CoreCancelTaskDto setProcInstId(String procInstId) {
        this.procInstId = procInstId;
        return this;
    }

    public String getTaskNodeId() {
        return taskNodeId;
    }

    public CoreCancelTaskDto setTaskNodeId(String taskNodeId) {
        this.taskNodeId = taskNodeId;
        return this;
    }
}

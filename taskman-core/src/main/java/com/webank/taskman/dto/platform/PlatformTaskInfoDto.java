package com.webank.taskman.dto.platform;

public class PlatformTaskInfoDto {
    private String procInstId;
    private String taskName;
    private String taskTempId;
    private String taskNodeId;
    private String taskDescription;
    private String roleName;
    private String reporter;
    private String callbackUrl;
    private String callbackParameter;

    public String getProcInstId() {
        return procInstId;
    }

    public PlatformTaskInfoDto setProcInstId(String procInstId) {
        this.procInstId = procInstId;
        return this;
    }

    public String getTaskName() {
        return taskName;
    }

    public PlatformTaskInfoDto setTaskName(String taskName) {
        this.taskName = taskName;
        return this;
    }

    public String getTaskTempId() {
        return taskTempId;
    }

    public PlatformTaskInfoDto setTaskTempId(String taskTempId) {
        this.taskTempId = taskTempId;
        return this;
    }

    public String getTaskNodeId() {
        return taskNodeId;
    }

    public PlatformTaskInfoDto setTaskNodeId(String taskNodeId) {
        this.taskNodeId = taskNodeId;
        return this;
    }

    public String getTaskDescription() {
        return taskDescription;
    }

    public PlatformTaskInfoDto setTaskDescription(String taskDescription) {
        this.taskDescription = taskDescription;
        return this;
    }

    public String getRoleName() {
        return roleName;
    }

    public PlatformTaskInfoDto setRoleName(String roleName) {
        this.roleName = roleName;
        return this;
    }

    public String getReporter() {
        return reporter;
    }

    public PlatformTaskInfoDto setReporter(String reporter) {
        this.reporter = reporter;
        return this;
    }

    public String getCallbackUrl() {
        return callbackUrl;
    }

    public PlatformTaskInfoDto setCallbackUrl(String callbackUrl) {
        this.callbackUrl = callbackUrl;
        return this;
    }

    public String getCallbackParameter() {
        return callbackParameter;
    }

    public PlatformTaskInfoDto setCallbackParameter(String callbackParameter) {
        this.callbackParameter = callbackParameter;
        return this;
    }
}

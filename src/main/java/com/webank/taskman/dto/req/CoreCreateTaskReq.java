package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.util.List;

@ApiModel
public class CoreCreateTaskReq {

    @ApiModelProperty(value = "",required = true,position = 1)
    private String procDefId;

    @ApiModelProperty(value = "",required = true,position = 2)
    private String taskNodeId;

    @ApiModelProperty(value = "",required = true,position = 3)
    private String taskName;

    @ApiModelProperty(value = "",required = true,position = 4)
    private String callbackUrl;

    @ApiModelProperty(value = "",required = true,position = 6)
    private String callbackParameter;

    @ApiModelProperty(value = "",position = 7)
    private String taskDescription;

    @ApiModelProperty(value = "",position = 8)
    private String reporter;

    @ApiModelProperty(value = "",position = 9)
    private String roleName;

    @ApiModelProperty(value = "",required = true,position = 9)
    private String dueDate;

    @ApiModelProperty(value = "",required = true,position = 10)
    private List<CoreCreateTaskInputBeanReq> inputs;

    public String getProcDefId() {
        return procDefId;
    }

    public void setProcDefId(String procDefId) {
        this.procDefId = procDefId;
    }

    public String getTaskNodeId() {
        return taskNodeId;
    }

    public void setTaskNodeId(String taskNodeId) {
        this.taskNodeId = taskNodeId;
    }

    public String getTaskName() {
        return taskName;
    }

    public void setTaskName(String taskName) {
        this.taskName = taskName;
    }

    public String getCallbackUrl() {
        return callbackUrl;
    }

    public void setCallbackUrl(String callbackUrl) {
        this.callbackUrl = callbackUrl;
    }

    public String getCallbackParameter() {
        return callbackParameter;
    }

    public void setCallbackParameter(String callbackParameter) {
        this.callbackParameter = callbackParameter;
    }

    public String getTaskDescription() {
        return taskDescription;
    }

    public void setTaskDescription(String taskDescription) {
        this.taskDescription = taskDescription;
    }

    public String getReporter() {
        return reporter;
    }

    public void setReporter(String reporter) {
        this.reporter = reporter;
    }

    public String getRoleName() {
        return roleName;
    }

    public void setRoleName(String roleName) {
        this.roleName = roleName;
    }

    public String getDueDate() {
        return dueDate;
    }

    public void setDueDate(String dueDate) {
        this.dueDate = dueDate;
    }

    public List<CoreCreateTaskInputBeanReq> getInputs() {
        return inputs;
    }

    public void setInputs(List<CoreCreateTaskInputBeanReq> inputs) {
        this.inputs = inputs;
    }


}

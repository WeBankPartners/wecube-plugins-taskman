package com.webank.taskman.dto;

import java.util.List;

public class CoreCreateTaskDto {

    private String requestId;
    private String dueDate;
    private List<String> allowedOptions;
    private List<TaskInfoReq> inputs;

    public String getRequestId() {
        return requestId;
    }
    public CoreCreateTaskDto setRequestId(String requestId) {
        this.requestId = requestId;
        return this;
    }

    public List<String> getAllowedOptions() {
        return allowedOptions;
    }

    public void setAllowedOptions(List<String> allowedOptions) {
        this.allowedOptions = allowedOptions;
    }

    public String getDueDate() {
        return dueDate;
    }

    public void setDueDate(String dueDate) {
        this.dueDate = dueDate;
    }

    public List<TaskInfoReq> getInputs() {
        return inputs;
    }

    public void setInputs(List<TaskInfoReq> inputs) {
        this.inputs = inputs;
    }

    public static class TaskInfoReq {
        private String id;
        private String procInstId;
        private String taskName;
        private String taskTempId;
        private String taskNodeId;
        private String taskDescription;
        private String roleName;
        private String overTime;
        private String reporter;
        private String callbackUrl;
        private String callbackParameter;
        private List<FormItemBean> formItems;

        public String getId() {
            return id;
        }

        public TaskInfoReq setId(String id) {
            this.id = id;
            return this;
        }

        public String getProcInstId() {
            return procInstId;
        }

        public TaskInfoReq setProcInstId(String procInstId) {
            this.procInstId = procInstId;
            return this;
        }

        public String getTaskName() {
            return taskName;
        }

        public TaskInfoReq setTaskName(String taskName) {
            this.taskName = taskName;
            return this;
        }

        public String getTaskTempId() {
            return taskTempId;
        }

        public TaskInfoReq setTaskTempId(String taskTempId) {
            this.taskTempId = taskTempId;
            return this;
        }

        public String getTaskNodeId() {
            return taskNodeId;
        }

        public TaskInfoReq setTaskNodeId(String taskNodeId) {
            this.taskNodeId = taskNodeId;
            return this;
        }

        public String getTaskDescription() {
            return taskDescription;
        }

        public TaskInfoReq setTaskDescription(String taskDescription) {
            this.taskDescription = taskDescription;
            return this;
        }

        public String getRoleName() {
            return roleName;
        }

        public TaskInfoReq setRoleName(String roleName) {
            this.roleName = roleName;
            return this;
        }

        public String getOverTime() {
            return overTime;
        }

        public TaskInfoReq setOverTime(String overTime) {
            this.overTime = overTime;
            return this;
        }

        public String getReporter() {
            return reporter;
        }

        public TaskInfoReq setReporter(String reporter) {
            this.reporter = reporter;
            return this;
        }

        public String getCallbackUrl() {
            return callbackUrl;
        }

        public TaskInfoReq setCallbackUrl(String callbackUrl) {
            this.callbackUrl = callbackUrl;
            return this;
        }

        public String getCallbackParameter() {
            return callbackParameter;
        }

        public TaskInfoReq setCallbackParameter(String callbackParameter) {
            this.callbackParameter = callbackParameter;
            return this;
        }

        public List<FormItemBean> getFormItems() {
            return formItems;
        }

        public TaskInfoReq setFormItems(List<FormItemBean> formItems) {
            this.formItems = formItems;
            return this;
        }

    }

    public static class FormItemBean {
        private String itemId;
        private String dataId;
        private String key;
        private List<String> val;

        public String getItemId() {
            return itemId;
        }

        public FormItemBean setItemId(String itemId) {
            this.itemId = itemId;
            return this;
        }

        public String getDataId() {
            return dataId;
        }

        public FormItemBean setDataId(String dataId) {
            this.dataId = dataId;
            return this;
        }

        public String getKey() {
            return key;
        }

        public FormItemBean setKey(String key) {
            this.key = key;
            return this;
        }

        public List<String> getVal() {
            return val;
        }

        public FormItemBean setVal(List<String> val) {
            this.val = val;
            return this;
        }
    }
}

package com.webank.taskman.dto;

import java.util.List;

public class CoreCreateTaskDTO {


    /**
     * allowedOptions :
     * dueDate :
     * inputs : [{"formItems":[{"itemId":"999","key":"app_inst","val":["0047_0000000026","0047_0000000027"]}],"procDefId":"","taskNodeId":"","overTime":123456,"callbackParameter":"","callbackUrl":"","reporter":"","roleName":"","taskDescription":"","taskName":""}]
     */

    private List<String> allowedOptions;
    private String dueDate;
    private List<InputsBean> inputs;

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

    public List<InputsBean> getInputs() {
        return inputs;
    }

    public void setInputs(List<InputsBean> inputs) {
        this.inputs = inputs;
    }

    public static class InputsBean {
        /**
         * formItems : [{"itemId":"999","key":"app_inst","val":["0047_0000000026","0047_0000000027"]}]
         * procDefId :
         * taskNodeId :
         * overTime : 123456
         * callbackParameter :
         * callbackUrl :
         * reporter :
         * roleName :
         * taskDescription :
         * taskName :
         */

        private String procDefId;
        private String taskNodeId;
        private int overTime;
        private String callbackParameter;
        private String callbackUrl;
        private String reporter;
        private String roleName;
        private String taskDescription;
        private String taskName;
        private List<FormItemsBean> formItems;

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

        public int getOverTime() {
            return overTime;
        }

        public void setOverTime(int overTime) {
            this.overTime = overTime;
        }

        public String getCallbackParameter() {
            return callbackParameter;
        }

        public void setCallbackParameter(String callbackParameter) {
            this.callbackParameter = callbackParameter;
        }

        public String getCallbackUrl() {
            return callbackUrl;
        }

        public void setCallbackUrl(String callbackUrl) {
            this.callbackUrl = callbackUrl;
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

        public String getTaskDescription() {
            return taskDescription;
        }

        public void setTaskDescription(String taskDescription) {
            this.taskDescription = taskDescription;
        }

        public String getTaskName() {
            return taskName;
        }

        public void setTaskName(String taskName) {
            this.taskName = taskName;
        }

        public List<FormItemsBean> getFormItems() {
            return formItems;
        }

        public void setFormItems(List<FormItemsBean> formItems) {
            this.formItems = formItems;
        }

        public static class FormItemsBean {
            /**
             * itemId : 999
             * key : app_inst
             * val : ["0047_0000000026","0047_0000000027"]
             */

            private String itemId;
            private String key;
            private List<String> val;

            public String getItemId() {
                return itemId;
            }

            public void setItemId(String itemId) {
                this.itemId = itemId;
            }

            public String getKey() {
                return key;
            }

            public void setKey(String key) {
                this.key = key;
            }

            public List<String> getVal() {
                return val;
            }

            public void setVal(List<String> val) {
                this.val = val;
            }
        }
    }
}

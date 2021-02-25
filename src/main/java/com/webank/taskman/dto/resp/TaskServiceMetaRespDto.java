package com.webank.taskman.dto.resp;

import java.util.List;

public class TaskServiceMetaRespDto {


    private List<TaskServiceMetaFormItem> formItems;

    public List<TaskServiceMetaFormItem> getFormItems() {
        return formItems;
    }

    public void setFormItems(List<TaskServiceMetaFormItem> formItems) {
        this.formItems = formItems;
    }

    public static class TaskServiceMetaFormItem {

        private String itemId;
        private String key;
        private TaskServiceMetaValueDef valueDef;

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

        public TaskServiceMetaValueDef getValueDef() {
            return valueDef;
        }

        public void setValueDef(TaskServiceMetaValueDef valueDef) {
            this.valueDef = valueDef;
        }

        
    }
    
    public static class TaskServiceMetaValueDef {

        private String type;
        private String expr;

        public String getType() {
            return type;
        }

        public void setType(String type) {
            this.type = type;
        }

        public String getExpr() {
            return expr;
        }

        public void setExpr(String expr) {
            this.expr = expr;
        }
    }

}

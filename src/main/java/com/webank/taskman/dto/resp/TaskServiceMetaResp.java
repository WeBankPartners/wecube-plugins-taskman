package com.webank.taskman.dto.resp;

import io.swagger.annotations.ApiModelProperty;

import java.util.List;

public class TaskServiceMetaResp {


    private List<TaskServiceMetaFormItem> formItems;

    public List<TaskServiceMetaFormItem> getFormItems() {
        return formItems;
    }

    public void setFormItems(List<TaskServiceMetaFormItem> formItems) {
        this.formItems = formItems;
    }

    public static class TaskServiceMetaFormItem {

        @ApiModelProperty(value = "",position = 1)
        private String itemId;
        @ApiModelProperty(value = "",position = 2)
        private String key;
        @ApiModelProperty(value = "",position = 3)
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

        public static class TaskServiceMetaValueDef {

            @ApiModelProperty(value = "",position = 1)
            private String type;
            @ApiModelProperty(value = "",position = 2)
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

}

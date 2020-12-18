package com.webank.taskman.dto.resp;

public class TaskServiceMetaDTO {


    private String itemId;
    private String key;
    private ValueDefBean valueDef;

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

    public ValueDefBean getValueDef() {
        return valueDef;
    }

    public void setValueDef(ValueDefBean valueDef) {
        this.valueDef = valueDef;
    }

    public static class ValueDefBean {

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

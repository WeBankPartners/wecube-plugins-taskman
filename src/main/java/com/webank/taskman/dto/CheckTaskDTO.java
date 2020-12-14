package com.webank.taskman.dto;

import com.webank.taskman.dto.resp.FormInfoResq;

public class CheckTaskDTO {
    private FormInfoResq formInfoResq;
    /**
     * itemId : 999
     * key : app_inst
     * valueDef : {"type":"ref","packageName":"wecmdb","entity":"app_instance"}
     */

    private String itemId;
    private String key;
    private ValueDefBean valueDef;

    public FormInfoResq getFormInfoResq() {
        return formInfoResq;
    }

    public void setFormInfoResq(FormInfoResq formInfoResq) {
        this.formInfoResq = formInfoResq;
    }

    @Override
    public String toString() {
        return "CheckTaskDTO{" +
                "formInfoResq=" + formInfoResq +
                '}';
    }

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
        /**
         * type : ref
         * packageName : wecmdb
         * entity : app_instance
         */

        private String type;
        private String packageName;
        private String entity;

        public String getType() {
            return type;
        }

        public void setType(String type) {
            this.type = type;
        }

        public String getPackageName() {
            return packageName;
        }

        public void setPackageName(String packageName) {
            this.packageName = packageName;
        }

        public String getEntity() {
            return entity;
        }

        public void setEntity(String entity) {
            this.entity = entity;
        }
    }
}

package com.webank.taskman.dto.resp;

import java.util.List;

public class CreateTaskServiceMetaResp {


    private List<FormItemsBean> formItems;

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
         * valueDef : {"type":"ref","expr":"wecmdb:app_instance"}
         */

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
            /**
             * type : ref
             * expr : wecmdb:app_instance
             */

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
}

package com.webank.taskman.dto.resp;

public class CreateTaskServiceMetaResp {

    private String itemId;
    private String key;
    private CreateTaskServiceMetaValueDefResp valueDef;

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

    public CreateTaskServiceMetaValueDefResp getValueDef() {
        return valueDef;
    }

    public void setValueDef(CreateTaskServiceMetaValueDefResp valueDef) {
        this.valueDef = valueDef;
    }
}

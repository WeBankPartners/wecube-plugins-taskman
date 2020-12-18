package com.webank.taskman.dto.resp;

import io.swagger.annotations.ApiModelProperty;

public class TaskServiceMetaFormItem {


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
}

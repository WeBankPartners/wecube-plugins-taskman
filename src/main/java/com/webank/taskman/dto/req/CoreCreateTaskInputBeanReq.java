package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.util.List;

@ApiModel
public class CoreCreateTaskInputBeanReq {


    @ApiModelProperty(required = true)
    private String itemId;
    @ApiModelProperty(required = true,value = "The actual value is the key that interface [service-meta] returns the result")
    private List<String> key;

    public String getItemId() {
        return itemId;
    }

    public void setItemId(String itemId) {
        this.itemId = itemId;
    }

    public List<String> getKey() {
        return key;
    }

    public void setKey(List<String> key) {
        this.key = key;
    }
}

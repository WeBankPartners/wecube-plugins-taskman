package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModel;

@ApiModel
public class SaveAndFormItemfoReq {

    private String id;

    private String formId;

    private String itemTempId;

    private String name;

    private String value;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getFormId() {
        return formId;
    }

    public void setFormId(String formId) {
        this.formId = formId;
    }

    public String getItemTempId() {
        return itemTempId;
    }

    public void setItemTempId(String itemTempId) {
        this.itemTempId = itemTempId;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getValue() {
        return value;
    }

    public void setValue(String value) {
        this.value = value;
    }

    @Override
    public String toString() {
        return "SaveAndFormItemfoReq{" +
                "id='" + id + '\'' +
                ", formId='" + formId + '\'' +
                ", itemTempId='" + itemTempId + '\'' +
                ", name='" + name + '\'' +
                ", value='" + value + '\'' +
                '}';
    }
}

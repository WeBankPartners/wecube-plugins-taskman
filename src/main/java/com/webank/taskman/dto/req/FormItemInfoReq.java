package com.webank.taskman.dto.req;

public class FormItemInfoReq {

    private String itemTempId;

    private String isCurrency;

    private String name;

    private String value;

    public String getItemTempId() {
        return itemTempId;
    }

    public void setItemTempId(String itemTempId) {
        this.itemTempId = itemTempId;
    }

    public String getIsCurrency() {
        return isCurrency;
    }

    public void setIsCurrency(String isCurrency) {
        this.isCurrency = isCurrency;
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
        return "FormItemInfoReq{" +
                "itemTempId='" + itemTempId + '\'' +
                ", isCurrency='" + isCurrency + '\'' +
                ", name='" + name + '\'' +
                ", value='" + value + '\'' +
                '}';
    }
}

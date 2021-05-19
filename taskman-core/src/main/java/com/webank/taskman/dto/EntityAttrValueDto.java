package com.webank.taskman.dto;

public class EntityAttrValueDto {
    private String itemTempId;
    private String attrDefId;
    private String dataType;
    private String name;
    private Object dataValue;

    public String getItemTempId() {
        return itemTempId;
    }

    public EntityAttrValueDto setItemTempId(String itemTempId) {
        this.itemTempId = itemTempId;
        return this;
    }

    public String getAttrDefId() {
        return attrDefId;
    }

    public EntityAttrValueDto setAttrDefId(String attrDefId) {
        this.attrDefId = attrDefId;
        return this;
    }

    public String getDataType() {
        return dataType;
    }

    public EntityAttrValueDto setDataType(String dataType) {
        this.dataType = dataType;
        return this;
    }

    public String getName() {
        return name;
    }

    public EntityAttrValueDto setName(String name) {
        this.name = name;
        return this;
    }

    public Object getDataValue() {
        return dataValue;
    }

    public EntityAttrValueDto setDataValue(Object dataValue) {
        this.dataValue = dataValue;
        return this;
    }

    @Override
    public String toString() {
        StringBuilder builder = new StringBuilder();
        builder.append("EntityAttrValueDto [itemTempId=");
        builder.append(itemTempId);
        builder.append(", attrDefId=");
        builder.append(attrDefId);
        builder.append(", dataType=");
        builder.append(dataType);
        builder.append(", name=");
        builder.append(name);
        builder.append(", dataValue=");
        builder.append(dataValue);
        builder.append("]");
        return builder.toString();
    }
}

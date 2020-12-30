package com.webank.taskman.constant;

public enum TemplateTypeEnum {

    REQUEST("0","request"),
    TASK("1","task"),
    ;
    private String type;

    private String description;

    TemplateTypeEnum(String type, String description) {
        this.type = type;
        this.description = description;
    }

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getDescription() {
        return description;
    }

    public TemplateTypeEnum setDescription(String description) {
        this.description = description;
        return this;
    }
}

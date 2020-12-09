package com.webank.taskman.constant;

public enum TemplateTypeEnum {

    REQUEST(0,"request"),
    TASK(1,"task"),
    ;
    private int type;

    private String description;

    TemplateTypeEnum(int type, String description) {
        this.type = type;
        this.description = description;
    }

    public int getType() {
        return type;
    }

    public void setType(int type) {
        this.type = type;
    }


}

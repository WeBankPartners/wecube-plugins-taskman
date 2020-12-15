package com.webank.taskman.constant;

public enum TaskNodeTypeEnum {

    SSTN("SSTN","SSTN node "),
    SUTN("SUTN","SUTN node"),//  artificial
    ;
    private String type;

    private String description;

    TaskNodeTypeEnum(String type, String description) {
        this.type = type;
        this.description = description;
    }

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }


}

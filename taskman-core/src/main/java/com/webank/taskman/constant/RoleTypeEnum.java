package com.webank.taskman.constant;

public enum RoleTypeEnum {

    MANAGE_ROLE(0, "MANAGE_ROLE"), USE_ROLE(1, "USE_ROLE");
    
    private int type;

    private String description;

    RoleTypeEnum(int type, String description) {
        this.type = type;
        this.description = description;
    }

    public int getType() {
        return type;
    }

    public void setType(int type) {
        this.type = type;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

}

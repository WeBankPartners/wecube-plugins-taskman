package com.webank.taskman.constant;

public enum TaskNodeType {

    SSTN("SSTN", "SSTN node "), SUTN("SUTN", "SUTN node");
    
    private String type;

    private String description;

    TaskNodeType(String type, String description) {
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

    public void setDescription(String description) {
        this.description = description;
    }

}

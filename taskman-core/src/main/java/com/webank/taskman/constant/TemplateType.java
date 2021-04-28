package com.webank.taskman.constant;

public enum TemplateType {

    REQUEST("0", "request"), TASK("1", "task");
    private String type;

    private String description;

    TemplateType(String type, String description) {
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

    public TemplateType setDescription(String description) {
        this.description = description;
        return this;
    }
}

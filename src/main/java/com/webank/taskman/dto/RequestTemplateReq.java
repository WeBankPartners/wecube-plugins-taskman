package com.webank.taskman.dto;

import io.swagger.annotations.ApiModelProperty;

public class RequestTemplateReq {

    private String id;
    private String groupId;
    private String formTempId;
    private String name;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getGroupId() {
        return groupId;
    }

    public void setGroupId(String groupId) {
        this.groupId = groupId;
    }

    public String getFormTempId() {
        return formTempId;
    }

    public void setFormTempId(String formTempId) {
        this.formTempId = formTempId;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    @Override
    public String toString() {
        return "RequestTemplateReq{" +
                "id='" + id + '\'' +
                ", groupId='" + groupId + '\'' +
                ", formTempId='" + formTempId + '\'' +
                ", name='" + name + '\'' +
                '}';
    }
}

package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModel;

import java.util.List;

@ApiModel
public class SaveFormInfoAndFormItemInfoReq {

    private String id;

    private String recordId;

    private String formTemplateId;

    private String name;

    private Integer type;

    private List<SaveFormItemInfoReq> saveFormItemInfoReqs;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getRecordId() {
        return recordId;
    }

    public void setRecordId(String recordId) {
        this.recordId = recordId;
    }

    public String getFormTemplateId() {
        return formTemplateId;
    }

    public void setFormTemplateId(String formTemplateId) {
        this.formTemplateId = formTemplateId;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public Integer getType() {
        return type;
    }

    public void setType(Integer type) {
        this.type = type;
    }

    public List<SaveFormItemInfoReq> getSaveFormItemInfoReqs() {
        return saveFormItemInfoReqs;
    }

    public void setSaveFormItemInfoReqs(List<SaveFormItemInfoReq> saveFormItemInfoReqs) {
        this.saveFormItemInfoReqs = saveFormItemInfoReqs;
    }

    @Override
    public String toString() {
        return "SaveFormInfoAndFormItemInfoReq{" +
                "id='" + id + '\'' +
                ", recordId='" + recordId + '\'' +
                ", formTemplateId='" + formTemplateId + '\'' +
                ", name='" + name + '\'' +
                ", type=" + type +
                ", saveAndFormItemfoReqs=" + saveFormItemInfoReqs +
                '}';
    }
}

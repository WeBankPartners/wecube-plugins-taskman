package com.webank.taskman.dto.resp;

import java.util.List;

public class FormInfoResqDto {

    private String id;
    private String recordId;
    private String formTemplateId;
    private String name;
    private Integer type;
    private String formType;
    private List<FormItemInfoRespDto> formItemInfo;

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

    public String getFormType() {
        return formType;
    }

    public FormInfoResqDto setFormType(String formType) {
        this.formType = formType;
        return this;
    }

    public List<FormItemInfoRespDto> getFormItemInfo() {
        return formItemInfo;
    }

    public void setFormItemInfo(List<FormItemInfoRespDto> formItemInfo) {
        this.formItemInfo = formItemInfo;
    }
}

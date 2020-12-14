package com.webank.taskman.dto.resp;

import com.webank.taskman.domain.FormItemInfo;

import java.util.List;

public class SynthesisRequestInfoFormRequest {
    private String id;

    private String recordId;

    private String formTemplateId;

    private String name;

    private Integer type;

    private List<FormItemInfo> formItemInfo;

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

    public List<FormItemInfo> getFormItemInfo() {
        return formItemInfo;
    }

    public void setFormItemInfo(List<FormItemInfo> formItemInfo) {
        this.formItemInfo = formItemInfo;
    }
}

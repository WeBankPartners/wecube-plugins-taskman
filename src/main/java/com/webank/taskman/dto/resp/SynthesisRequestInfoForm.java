package com.webank.taskman.dto.resp;

import io.swagger.annotations.ApiModel;

import java.util.List;

@ApiModel
public class SynthesisRequestInfoForm {
    private String id;

    private String recordId;

    private String formTemplateId;

    private String name;

    private Integer type;

    private String emergency;

    private String description;

    private String rootEntity;

    private List<FormItemInfoResp> formItemInfo;

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

    public String getEmergency() {
        return emergency;
    }

    public void setEmergency(String emergency) {
        this.emergency = emergency;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getRootEntity() {
        return rootEntity;
    }

    public void setRootEntity(String rootEntity) {
        this.rootEntity = rootEntity;
    }

    public List<FormItemInfoResp> getFormItemInfo() {
        return formItemInfo;
    }

    public void setFormItemInfo(List<FormItemInfoResp> formItemInfo) {
        this.formItemInfo = formItemInfo;
    }
}

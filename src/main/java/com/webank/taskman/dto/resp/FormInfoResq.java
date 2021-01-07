package com.webank.taskman.dto.resp;

import com.webank.taskman.domain.FormItemInfo;
import io.swagger.annotations.ApiModelProperty;

import java.util.List;

public class FormInfoResq {

    @ApiModelProperty(value = "",position = 1)
    private String id;
    @ApiModelProperty(value = "",position = 2)
    private String recordId;
    @ApiModelProperty(value = "",position = 3)
    private String formTemplateId;
    @ApiModelProperty(value = "",position = 4)
    private String name;
    @ApiModelProperty(value = "",position = 5)
    private Integer type;
    @ApiModelProperty(value = "",position = 6)
    private String formType;
    @ApiModelProperty(value = "",position = 7)
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

    public String getFormType() {
        return formType;
    }

    public FormInfoResq setFormType(String formType) {
        this.formType = formType;
        return this;
    }

    public List<FormItemInfoResp> getFormItemInfo() {
        return formItemInfo;
    }

    public void setFormItemInfo(List<FormItemInfoResp> formItemInfo) {
        this.formItemInfo = formItemInfo;
    }
}

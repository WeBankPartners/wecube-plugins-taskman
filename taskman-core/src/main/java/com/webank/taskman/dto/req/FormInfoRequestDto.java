package com.webank.taskman.dto.req;

import java.util.List;

public class FormInfoRequestDto {

    private String name;

    private String description;

    private String inputAttrDef;

    private String outputAttrDef;

    private String otherAttrDef;

    private String formType;

    private List<FormItemInfoRequestDto> formItems;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getInputAttrDef() {
        return inputAttrDef;
    }

    public void setInputAttrDef(String inputAttrDef) {
        this.inputAttrDef = inputAttrDef;
    }

    public String getOutputAttrDef() {
        return outputAttrDef;
    }

    public void setOutputAttrDef(String outputAttrDef) {
        this.outputAttrDef = outputAttrDef;
    }

    public String getOtherAttrDef() {
        return otherAttrDef;
    }

    public void setOtherAttrDef(String otherAttrDef) {
        this.otherAttrDef = otherAttrDef;
    }

    public String getFormType() {
        return formType;
    }

    public FormInfoRequestDto setFormType(String formType) {
        this.formType = formType;
        return this;
    }

    public List<FormItemInfoRequestDto> getFormItems() {
        return formItems;
    }

    public void setFormItems(List<FormItemInfoRequestDto> formItems) {
        this.formItems = formItems;
    }

}

package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.util.List;

@ApiModel
public class SaveFormInfoReq {


    @ApiModelProperty(value = "",required = true,position = 1)
    private String name;

    @ApiModelProperty(value = "",position = 2)
    private String description;

    @ApiModelProperty(value = "",position = 3)
    private String inputAttrDef;

    @ApiModelProperty(value = "",position = 4)
    private String outputAttrDef;

    @ApiModelProperty(value = "",position = 5)
    private String otherAttrDef;

    @ApiModelProperty(value = "",position = 6)
    private String formType;

    private List<SaveFormItemInfoReq> formItems;

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

    public SaveFormInfoReq setFormType(String formType) {
        this.formType = formType;
        return this;
    }

    public List<SaveFormItemInfoReq> getFormItems() {
        return formItems;
    }

    public void setFormItems(List<SaveFormItemInfoReq> formItems) {
        this.formItems = formItems;
    }


}

package com.webank.taskman.dto.resp;

import io.swagger.annotations.ApiModelProperty;

import java.util.List;

public class FormTemplateResp {

    @ApiModelProperty(value = "",position = 1)
    private String id;

    @ApiModelProperty(value = "",position = 2)
    private String tempId;

    @ApiModelProperty(value = "",position = 3)
    private String tempType;

    @ApiModelProperty(value = "",position = 4)
    private String formType;

    @ApiModelProperty(value = "",position = 4)
    private String name;

    @ApiModelProperty(value = "",position = 5)
    private String description;

    @ApiModelProperty(value = "",position = 6)
    private String style;

    @ApiModelProperty(value = "",position = 7)
    private String targetEntitys;

    @ApiModelProperty(value = "",position = 8)
    private String inputAttrDef;

    @ApiModelProperty(value = "",position = 9)
    private String outputAttrDef;

    @ApiModelProperty(value = "",position = 10)
    private String otherAttrDef;

    @ApiModelProperty(value = "",position = 11)
    private List<FormItemTemplateResp> items;


    public String getId() {
        return id;
    }

    public FormTemplateResp setId(String id) {
        this.id = id;
        return this;
    }

    public String getTempId() {
        return tempId;
    }

    public void setTempId(String tempId) {
        this.tempId = tempId;
    }

    public String getTempType() {
        return tempType;
    }

    public void setTempType(String tempType) {
        this.tempType = tempType;
    }

    public String getFormType() {
        return formType;
    }

    public FormTemplateResp setFormType(String formType) {
        this.formType = formType;
        return this;
    }

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

    public String getStyle() {
        return style;
    }

    public void setStyle(String style) {
        this.style = style;
    }

    public String getTargetEntitys() {
        return targetEntitys;
    }

    public void setTargetEntitys(String targetEntitys) {
        this.targetEntitys = targetEntitys;
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

    public List<FormItemTemplateResp> getItems() {
        return items;
    }

    public void setItems(List<FormItemTemplateResp> items) {
        this.items = items;
    }

    @Override
    public String toString() {
        return "FormTemplateResp{" +
                "id='" + id + '\'' +
                ", tempId='" + tempId + '\'' +
                ", tempType='" + tempType + '\'' +
                ", name='" + name + '\'' +
                ", description='" + description + '\'' +
                ", style='" + style + '\'' +
                ", targetEntitys='" + targetEntitys + '\'' +
                ", inputAttrDef='" + inputAttrDef + '\'' +
                ", outputAttrDef='" + outputAttrDef + '\'' +
                ", otherAttrDef='" + otherAttrDef + '\'' +
                ", items=" + items +
                '}';
    }
}

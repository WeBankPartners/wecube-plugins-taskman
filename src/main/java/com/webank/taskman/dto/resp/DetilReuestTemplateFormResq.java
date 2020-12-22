package com.webank.taskman.dto.resp;

import com.webank.taskman.domain.FormItemTemplate;

import java.util.List;

public class DetilReuestTemplateFormResq {
    private String id;

    private String tempId;

    private String tempType;

    private String name;

    private String description;

    private String style;

    private String targetEntitys;

    private String inputAttrDef;

    private String outputAttrDef;

    private String otherAttrDef;

    private List<FormItemTemplate> formItemTemplateList;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
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

    public List<FormItemTemplate> getFormItemTemplateList() {
        return formItemTemplateList;
    }

    public void setFormItemTemplateList(List<FormItemTemplate> formItemTemplateList) {
        this.formItemTemplateList = formItemTemplateList;
    }
}

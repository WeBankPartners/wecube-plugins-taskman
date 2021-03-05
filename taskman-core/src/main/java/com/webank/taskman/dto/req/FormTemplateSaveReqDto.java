package com.webank.taskman.dto.req;

import java.util.LinkedList;
import java.util.List;

import com.webank.taskman.dto.RoleDto;

public class FormTemplateSaveReqDto {

    private String id;

    private String tempId;

    // 模板类型: 0.请求模板 1.任务模板
    // = "int",position = 3)
    private String tempType = "0";

    private String formType;

    private String name;

    private String description;

    private String style;
    private String targetEntities;

    private String inputAttrDef;
    private String outputAttrDef;
    private String otherAttrDef;

    private List<RoleDto> useRole;

    private List<FormItemTemplateSaveReqDto> formItems = new LinkedList<>();

    public FormTemplateSaveReqDto() {
    }

    public FormTemplateSaveReqDto(String tempId, String tempType) {
        this.tempId = tempId;
        this.tempType = tempType;
    }

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

    public String getFormType() {
        return formType;
    }

    public FormTemplateSaveReqDto setFormType(String formType) {
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

    public String getTargetEntities() {
        return targetEntities;
    }

    public void setTargetEntities(String targetEntities) {
        this.targetEntities = targetEntities;
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

    public List<RoleDto> getUseRole() {
        return useRole;
    }

    public void setUseRole(List<RoleDto> useRole) {
        this.useRole = useRole;
    }

    public List<FormItemTemplateSaveReqDto> getFormItems() {
        return formItems;
    }

    public void setFormItems(List<FormItemTemplateSaveReqDto> formItems) {
        this.formItems = formItems;
    }
}

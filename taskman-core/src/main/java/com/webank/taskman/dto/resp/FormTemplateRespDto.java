package com.webank.taskman.dto.resp;

import java.util.ArrayList;
import java.util.List;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.webank.taskman.dto.RoleDto;

@JsonInclude(JsonInclude.Include.NON_NULL)
public class FormTemplateRespDto {

    private String id;

    private String tempId;

    private String tempType;

    private String formType;

    private String name;

    private String description;

    private String style;

    private String targetEntities;

    private String inputAttrDef;
    private String outputAttrDef;
    private String otherAttrDef;

    private List<FormItemTemplateRespDto> items;

    private List<TaskTemplateRespDto> taskTemplates;
    
    private List<RoleDto> useRoles = new ArrayList<>();

    private List<RoleDto> manageRoles = new ArrayList<>();

    public String getId() {
        return id;
    }

    public FormTemplateRespDto setId(String id) {
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

    public FormTemplateRespDto setFormType(String formType) {
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

    public List<FormItemTemplateRespDto> getItems() {
        return items;
    }

    public void setItems(List<FormItemTemplateRespDto> items) {
        this.items = items;
    }

    @Override
    public String toString() {
        return "FormTemplateResp{" + "id='" + id + '\'' + ", tempId='" + tempId + '\'' + ", tempType='" + tempType
                + '\'' + ", name='" + name + '\'' + ", description='" + description + '\'' + ", style='" + style + '\''
                + ", targetEntities='" + targetEntities + '\'' + ", inputAttrDef='" + inputAttrDef + '\''
                + ", outputAttrDef='" + outputAttrDef + '\'' + ", otherAttrDef='" + otherAttrDef + '\'' + ", items="
                + items + '}';
    }

    public List<TaskTemplateRespDto> getTaskTemplates() {
        return taskTemplates;
    }

    public void setTaskTemplates(List<TaskTemplateRespDto> taskTemplates) {
        this.taskTemplates = taskTemplates;
    }

    public List<RoleDto> getUseRoles() {
        return useRoles;
    }

    public void setUseRoles(List<RoleDto> useRoles) {
        this.useRoles = useRoles;
    }

    public List<RoleDto> getManageRoles() {
        return manageRoles;
    }

    public void setManageRoles(List<RoleDto> manageRoles) {
        this.manageRoles = manageRoles;
    }
    
    
}
package com.webank.taskman.dto.req;

import java.util.LinkedList;
import java.util.List;

import com.webank.taskman.dto.RoleDTO;

public class SaveFormTemplateReq {

    private String id;

    private String tempId;

    // @ApiModelProperty(value = "模板类型(0.请求模板 1.任务模板)",required = true,dataType
    // = "int",position = 3)
    private String tempType = "0";

    private String formType;

    private String name;

    private String description;

    private String style;

    // TODO
    // FIXME misspelling
    private String targetEntitys;

    private String inputAttrDef;
    private String outputAttrDef;
    private String otherAttrDef;

    private List<RoleDTO> useRole;

    private List<SaveFormItemTemplateReq> formItems = new LinkedList<>();

    public SaveFormTemplateReq() {
    }

    public SaveFormTemplateReq(String tempId, String tempType) {
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

    public SaveFormTemplateReq setFormType(String formType) {
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

    public List<RoleDTO> getUseRole() {
        return useRole;
    }

    public void setUseRole(List<RoleDTO> useRole) {
        this.useRole = useRole;
    }

    public List<SaveFormItemTemplateReq> getFormItems() {
        return formItems;
    }

    public void setFormItems(List<SaveFormItemTemplateReq> formItems) {
        this.formItems = formItems;
    }
}

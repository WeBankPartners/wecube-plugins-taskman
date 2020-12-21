package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.webank.taskman.base.BaseEntity;

import java.io.Serializable;

public class FormTemplate extends BaseEntity implements Serializable {

    private static final long serialVersionUID = 1L;


    @TableId(value = "id", type = IdType.ASSIGN_ID)
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


    public FormTemplate() {
    }

    public FormTemplate(String id, String tempId, String tempType) {
        this.id = id;
        this.tempId = tempId;
        this.tempType = tempType;
    }

    public FormTemplate(String id, String tempId) {
        this.id = id;
        this.tempId = tempId;
    }

    public FormTemplate(String id, String tempId, String tempType, String name, String description) {
        this.id = id;
        this.tempId = tempId;
        this.tempType = tempType;
        this.name = name;
        this.description = description;
    }
    public FormTemplate(String id, String tempId, String tempType, String name, String description, String style) {
        this.id = id;
        this.tempId = tempId;
        this.tempType = tempType;
        this.name = name;
        this.description = description;
        this.style = style;
    }

    public FormTemplate(String id, String tempId, String tempType, String name, String description, String style, String targetEntitys, String inputAttrDef, String outputAttrDef, String otherAttrDef) {
        this.id = id;
        this.tempId = tempId;
        this.tempType = tempType;
        this.name = name;
        this.description = description;
        this.style = style;
        this.targetEntitys = targetEntitys;
        this.inputAttrDef = inputAttrDef;
        this.outputAttrDef = outputAttrDef;
        this.otherAttrDef = otherAttrDef;
    }

    public  static QueryWrapper<FormTemplate> getQueryWrapper(String tempId,String tempType){
        QueryWrapper<FormTemplate> queryWrapper = new QueryWrapper<>();
        FormTemplate formTemplate = new FormTemplate(null,tempId,tempType);
        queryWrapper.setEntity(formTemplate);
        return queryWrapper;
    }
    public  static QueryWrapper<FormTemplate> getQueryWrapper(String id, String tempId, String tempType, String name, String description){
        QueryWrapper<FormTemplate> queryWrapper = new QueryWrapper<>();
        FormTemplate formTemplate = new FormTemplate(id,tempId,tempType);
        queryWrapper.setEntity(formTemplate);
        queryWrapper.like("name",name);
        queryWrapper.like("description",description);
        return queryWrapper;
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

    public static long getSerialVersionUID() {
        return serialVersionUID;
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

    @Override
    public String toString() {
        return "FormTemplate{" +
                "id='" + id + '\'' +
                ", tempId='" + tempId + '\'' +
                ", tempType='" + tempType + '\'' +
                ", name='" + name + '\'' +
                ", description='" + description + '\'' +
                ", targetEntitys='" + targetEntitys + '\'' +
                ", style='" + style + '\'' +
                '}';
    }
}

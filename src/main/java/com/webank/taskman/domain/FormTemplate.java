package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.webank.taskman.base.BaseEntity;
import org.springframework.util.StringUtils;

import java.io.Serializable;
import java.util.StringJoiner;

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

    @JsonIgnore
    public LambdaQueryWrapper<FormTemplate> getLambdaQueryWrapper() {
        return new LambdaQueryWrapper<FormTemplate>()
                .eq(!StringUtils.isEmpty(id), FormTemplate::getId, id)
                .eq(!StringUtils.isEmpty(tempId), FormTemplate::getTempId, tempId)
                .eq(!StringUtils.isEmpty(tempType), FormTemplate::getTempType, tempType)
                .like(!StringUtils.isEmpty(name), FormTemplate::getName, name)
                .like(!StringUtils.isEmpty(description), FormTemplate::getDescription, description)
                .like(!StringUtils.isEmpty(style), FormTemplate::getStyle, style)
                .eq(!StringUtils.isEmpty(targetEntitys), FormTemplate::getTargetEntitys, targetEntitys)
                .eq(!StringUtils.isEmpty(inputAttrDef), FormTemplate::getInputAttrDef, inputAttrDef)
                .eq(!StringUtils.isEmpty(outputAttrDef), FormTemplate::getOutputAttrDef, outputAttrDef)
                .eq(!StringUtils.isEmpty(otherAttrDef), FormTemplate::getOtherAttrDef, otherAttrDef)
                ;
    }

    public String getId() {
        return id;
    }

    public FormTemplate setId(String id) {
        this.id = id;
        return this;
    }

    public String getTempId() {
        return tempId;
    }

    public FormTemplate setTempId(String tempId) {
        this.tempId = tempId;
        return this;
    }

    public String getTempType() {
        return tempType;
    }

    public FormTemplate setTempType(String tempType) {
        this.tempType = tempType;
        return this;
    }

    public String getName() {
        return name;
    }

    public FormTemplate setName(String name) {
        this.name = name;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public FormTemplate setDescription(String description) {
        this.description = description;
        return this;
    }

    public String getStyle() {
        return style;
    }

    public FormTemplate setStyle(String style) {
        this.style = style;
        return this;
    }

    public String getTargetEntitys() {
        return targetEntitys;
    }

    public FormTemplate setTargetEntitys(String targetEntitys) {
        this.targetEntitys = targetEntitys;
        return this;
    }

    public String getInputAttrDef() {
        return inputAttrDef;
    }

    public FormTemplate setInputAttrDef(String inputAttrDef) {
        this.inputAttrDef = inputAttrDef;
        return this;
    }

    public String getOutputAttrDef() {
        return outputAttrDef;
    }

    public FormTemplate setOutputAttrDef(String outputAttrDef) {
        this.outputAttrDef = outputAttrDef;
        return this;
    }

    public String getOtherAttrDef() {
        return otherAttrDef;
    }

    public FormTemplate setOtherAttrDef(String otherAttrDef) {
        this.otherAttrDef = otherAttrDef;
        return this;
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", FormTemplate.class.getSimpleName() + "[", "]")
                .add("id='" + id + "'")
                .add("tempId='" + tempId + "'")
                .add("tempType='" + tempType + "'")
                .add("name='" + name + "'")
                .add("description='" + description + "'")
                .add("style='" + style + "'")
                .add("targetEntitys='" + targetEntitys + "'")
                .add("inputAttrDef='" + inputAttrDef + "'")
                .add("outputAttrDef='" + outputAttrDef + "'")
                .add("otherAttrDef='" + otherAttrDef + "'")
                .toString();
    }
}

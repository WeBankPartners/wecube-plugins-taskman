package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.fasterxml.jackson.annotation.JsonIgnore;
import io.swagger.annotations.ApiModelProperty;
import org.springframework.util.StringUtils;

import java.io.Serializable;

public class FormItemTemplate implements Serializable {

    private static final long serialVersionUID = 1L;

    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;
    private String tempId;
    private String formTemplateId;
    private String name;
    private String defaultValue;
    private Integer isCurrency;
    private Integer sort;

    private String packageName;
    private String entity;
    private String entityId;
    private String attrDefId;
    private String attrDefDataType;
    private String routineExp;

    private String elementType;
    private String title;
    private Integer width;

    private String refEntity;
    private String refPackageName;
    private String refFilters;
    private String dataOptions;

    private Integer required;
    private String regular;
    private Integer isEdit;
    private Integer isView;

    public FormItemTemplate() {
    }

    public FormItemTemplate(String id) {
        this.id = id;
    }

    public FormItemTemplate(String id, String formTemplateId) {
        this.id = id;
        this.formTemplateId = formTemplateId;
    }

    @JsonIgnore
    public LambdaQueryWrapper<FormItemTemplate> getLambdaQueryWrapper() {
        return new LambdaQueryWrapper<FormItemTemplate>()
                .eq(!StringUtils.isEmpty(id), FormItemTemplate::getId, id)
                .eq(!StringUtils.isEmpty(tempId), FormItemTemplate::getTempId, tempId)
                .eq(!StringUtils.isEmpty(formTemplateId), FormItemTemplate::getFormTemplateId, formTemplateId)
                .eq(!StringUtils.isEmpty(name), FormItemTemplate::getName, name)
                .eq(!StringUtils.isEmpty(defaultValue), FormItemTemplate::getDefaultValue, defaultValue)
                .eq(!StringUtils.isEmpty(isCurrency), FormItemTemplate::getIsCurrency, isCurrency)
                .eq(!StringUtils.isEmpty(sort), FormItemTemplate::getSort, sort)
                .eq(!StringUtils.isEmpty(packageName), FormItemTemplate::getPackageName, packageName)
                .eq(!StringUtils.isEmpty(entity), FormItemTemplate::getEntity, entity)
                .eq(!StringUtils.isEmpty(entityId), FormItemTemplate::getEntityId, entityId)
                .eq(!StringUtils.isEmpty(attrDefId), FormItemTemplate::getAttrDefId, attrDefId)
                .eq(!StringUtils.isEmpty(attrDefDataType), FormItemTemplate::getAttrDefDataType, attrDefDataType)
                .eq(!StringUtils.isEmpty(elementType), FormItemTemplate::getElementType, elementType)
                .eq(!StringUtils.isEmpty(title), FormItemTemplate::getTitle, title)
                .eq(!StringUtils.isEmpty(width), FormItemTemplate::getWidth, width)
                .eq(!StringUtils.isEmpty(refEntity), FormItemTemplate::getRefEntity, refEntity)
                .eq(!StringUtils.isEmpty(refPackageName), FormItemTemplate::getRefPackageName, refPackageName)
                .eq(!StringUtils.isEmpty(refFilters), FormItemTemplate::getRefFilters, refFilters)
                .eq(!StringUtils.isEmpty(dataOptions), FormItemTemplate::getDataOptions, dataOptions)
                .eq(!StringUtils.isEmpty(required), FormItemTemplate::getRequired, required)
                .eq(!StringUtils.isEmpty(regular), FormItemTemplate::getRegular, regular)
                .eq(!StringUtils.isEmpty(isEdit), FormItemTemplate::getIsEdit, isEdit)
                .eq(!StringUtils.isEmpty(isView), FormItemTemplate::getIsView, isView);
    }

    public static long getSerialVersionUID() {
        return serialVersionUID;
    }

    public String getId() {
        return id;
    }

    public FormItemTemplate setId(String id) {
        this.id = id;
        return this;
    }

    public String getTempId() {
        return tempId;
    }

    public FormItemTemplate setTempId(String tempId) {
        this.tempId = tempId;
        return this;
    }


    public String getFormTemplateId() {
        return formTemplateId;
    }

    public FormItemTemplate setFormTemplateId(String formTemplateId) {
        this.formTemplateId = formTemplateId;
        return this;
    }

    public String getName() {
        return name;
    }

    public FormItemTemplate setName(String name) {
        this.name = name;
        return this;
    }

    public String getDefaultValue() {
        return defaultValue;
    }

    public FormItemTemplate setDefaultValue(String defaultValue) {
        this.defaultValue = defaultValue;
        return this;
    }

    public Integer getIsCurrency() {
        return isCurrency;
    }

    public FormItemTemplate setIsCurrency(Integer isCurrency) {
        this.isCurrency = isCurrency;
        return this;
    }

    public Integer getSort() {
        return sort;
    }

    public FormItemTemplate setSort(Integer sort) {
        this.sort = sort;
        return this;
    }

    public String getPackageName() {
        return packageName;
    }

    public FormItemTemplate setPackageName(String packageName) {
        this.packageName = packageName;
        return this;
    }

    public String getEntity() {
        return entity;
    }

    public FormItemTemplate setEntity(String entity) {
        this.entity = entity;
        return this;
    }

    public String getEntityId() {
        return entityId;
    }

    public FormItemTemplate setEntityId(String entityId) {
        this.entityId = entityId;
        return this;
    }

    public String getAttrDefId() {
        return attrDefId;
    }

    public FormItemTemplate setAttrDefId(String attrDefId) {
        this.attrDefId = attrDefId;
        return this;
    }

    public String getAttrDefDataType() {
        return attrDefDataType;
    }

    public FormItemTemplate setAttrDefDataType(String attrDefDataType) {
        this.attrDefDataType = attrDefDataType;
        return this;
    }

    public String getRoutineExp() {
        return routineExp;
    }

    public FormItemTemplate setRoutineExp(String routineExp) {
        this.routineExp = routineExp;
        return this;
    }

    public String getElementType() {
        return elementType;
    }

    public FormItemTemplate setElementType(String elementType) {
        this.elementType = elementType;
        return this;
    }

    public String getTitle() {
        return title;
    }

    public FormItemTemplate setTitle(String title) {
        this.title = title;
        return this;
    }

    public Integer getWidth() {
        return width;
    }

    public FormItemTemplate setWidth(Integer width) {
        this.width = width;
        return this;
    }

    public String getRefEntity() {
        return refEntity;
    }

    public FormItemTemplate setRefEntity(String refEntity) {
        this.refEntity = refEntity;
        return this;
    }

    public String getRefPackageName() {
        return refPackageName;
    }

    public FormItemTemplate setRefPackageName(String refPackageName) {
        this.refPackageName = refPackageName;
        return this;
    }

    public String getRefFilters() {
        return refFilters;
    }

    public FormItemTemplate setRefFilters(String refFilters) {
        this.refFilters = refFilters;
        return this;
    }

    public String getDataOptions() {
        return dataOptions;
    }

    public FormItemTemplate setDataOptions(String dataOptions) {
        this.dataOptions = dataOptions;
        return this;
    }

    public Integer getRequired() {
        return required;
    }

    public FormItemTemplate setRequired(Integer required) {
        this.required = required;
        return this;
    }

    public String getRegular() {
        return regular;
    }

    public FormItemTemplate setRegular(String regular) {
        this.regular = regular;
        return this;
    }

    public Integer getIsEdit() {
        return isEdit;
    }

    public FormItemTemplate setIsEdit(Integer isEdit) {
        this.isEdit = isEdit;
        return this;
    }

    public Integer getIsView() {
        return isView;
    }

    public FormItemTemplate setIsView(Integer isView) {
        this.isView = isView;
        return this;
    }
}

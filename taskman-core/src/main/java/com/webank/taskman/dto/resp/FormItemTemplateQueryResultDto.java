package com.webank.taskman.dto.resp;

import com.fasterxml.jackson.annotation.JsonInclude;

@JsonInclude(JsonInclude.Include.NON_NULL)
public class FormItemTemplateQueryResultDto {

    private String id;
    private String tempId;
    private String formType;
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

    public String getId() {
        return id;
    }

    public FormItemTemplateQueryResultDto setId(String id) {
        this.id = id;
        return this;
    }

    public String getTempId() {
        return tempId;
    }

    public FormItemTemplateQueryResultDto setTempId(String tempId) {
        this.tempId = tempId;
        return this;
    }

    public String getFormType() {
        return formType;
    }

    public FormItemTemplateQueryResultDto setFormType(String formType) {
        this.formType = formType;
        return this;
    }

    public String getFormTemplateId() {
        return formTemplateId;
    }

    public FormItemTemplateQueryResultDto setFormTemplateId(String formTemplateId) {
        this.formTemplateId = formTemplateId;
        return this;
    }

    public String getName() {
        return name;
    }

    public FormItemTemplateQueryResultDto setName(String name) {
        this.name = name;
        return this;
    }

    public String getDefaultValue() {
        return defaultValue;
    }

    public FormItemTemplateQueryResultDto setDefaultValue(String defaultValue) {
        this.defaultValue = defaultValue;
        return this;
    }

    public Integer getIsCurrency() {
        return isCurrency;
    }

    public FormItemTemplateQueryResultDto setIsCurrency(Integer isCurrency) {
        this.isCurrency = isCurrency;
        return this;
    }

    public Integer getSort() {
        return sort;
    }

    public FormItemTemplateQueryResultDto setSort(Integer sort) {
        this.sort = sort;
        return this;
    }

    public String getPackageName() {
        return packageName;
    }

    public FormItemTemplateQueryResultDto setPackageName(String packageName) {
        this.packageName = packageName;
        return this;
    }

    public String getEntity() {
        return entity;
    }

    public FormItemTemplateQueryResultDto setEntity(String entity) {
        this.entity = entity;
        return this;
    }

    public String getEntityId() {
        return entityId;
    }

    public FormItemTemplateQueryResultDto setEntityId(String entityId) {
        this.entityId = entityId;
        return this;
    }

    public String getAttrDefId() {
        return attrDefId;
    }

    public FormItemTemplateQueryResultDto setAttrDefId(String attrDefId) {
        this.attrDefId = attrDefId;
        return this;
    }

    public String getAttrDefDataType() {
        return attrDefDataType;
    }

    public FormItemTemplateQueryResultDto setAttrDefDataType(String attrDefDataType) {
        this.attrDefDataType = attrDefDataType;
        return this;
    }

    public String getRoutineExp() {
        return routineExp;
    }

    public FormItemTemplateQueryResultDto setRoutineExp(String routineExp) {
        this.routineExp = routineExp;
        return this;
    }

    public String getElementType() {
        return elementType;
    }

    public FormItemTemplateQueryResultDto setElementType(String elementType) {
        this.elementType = elementType;
        return this;
    }

    public String getTitle() {
        return title;
    }

    public FormItemTemplateQueryResultDto setTitle(String title) {
        this.title = title;
        return this;
    }

    public Integer getWidth() {
        return width;
    }

    public FormItemTemplateQueryResultDto setWidth(Integer width) {
        this.width = width;
        return this;
    }

    public String getRefEntity() {
        return refEntity;
    }

    public FormItemTemplateQueryResultDto setRefEntity(String refEntity) {
        this.refEntity = refEntity;
        return this;
    }

    public String getRefPackageName() {
        return refPackageName;
    }

    public FormItemTemplateQueryResultDto setRefPackageName(String refPackageName) {
        this.refPackageName = refPackageName;
        return this;
    }

    public String getRefFilters() {
        return refFilters;
    }

    public FormItemTemplateQueryResultDto setRefFilters(String refFilters) {
        this.refFilters = refFilters;
        return this;
    }

    public String getDataOptions() {
        return dataOptions;
    }

    public FormItemTemplateQueryResultDto setDataOptions(String dataOptions) {
        this.dataOptions = dataOptions;
        return this;
    }

    public Integer getRequired() {
        return required;
    }

    public FormItemTemplateQueryResultDto setRequired(Integer required) {
        this.required = required;
        return this;
    }

    public String getRegular() {
        return regular;
    }

    public FormItemTemplateQueryResultDto setRegular(String regular) {
        this.regular = regular;
        return this;
    }

    public Integer getIsEdit() {
        return isEdit;
    }

    public FormItemTemplateQueryResultDto setIsEdit(Integer isEdit) {
        this.isEdit = isEdit;
        return this;
    }

    public Integer getIsView() {
        return isView;
    }

    public FormItemTemplateQueryResultDto setIsView(Integer isView) {
        this.isView = isView;
        return this;
    }
}

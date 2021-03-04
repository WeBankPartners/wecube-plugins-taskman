package com.webank.taskman.dto.resp;

import com.fasterxml.jackson.annotation.JsonInclude;

@JsonInclude(JsonInclude.Include.NON_NULL)
public class FormItemTemplateRespDto {

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

    public FormItemTemplateRespDto setId(String id) {
        this.id = id;
        return this;
    }

    public String getTempId() {
        return tempId;
    }

    public FormItemTemplateRespDto setTempId(String tempId) {
        this.tempId = tempId;
        return this;
    }

    public String getFormType() {
        return formType;
    }

    public FormItemTemplateRespDto setFormType(String formType) {
        this.formType = formType;
        return this;
    }

    public String getFormTemplateId() {
        return formTemplateId;
    }

    public FormItemTemplateRespDto setFormTemplateId(String formTemplateId) {
        this.formTemplateId = formTemplateId;
        return this;
    }

    public String getName() {
        return name;
    }

    public FormItemTemplateRespDto setName(String name) {
        this.name = name;
        return this;
    }

    public String getDefaultValue() {
        return defaultValue;
    }

    public FormItemTemplateRespDto setDefaultValue(String defaultValue) {
        this.defaultValue = defaultValue;
        return this;
    }

    public Integer getIsCurrency() {
        return isCurrency;
    }

    public FormItemTemplateRespDto setIsCurrency(Integer isCurrency) {
        this.isCurrency = isCurrency;
        return this;
    }

    public Integer getSort() {
        return sort;
    }

    public FormItemTemplateRespDto setSort(Integer sort) {
        this.sort = sort;
        return this;
    }

    public String getPackageName() {
        return packageName;
    }

    public FormItemTemplateRespDto setPackageName(String packageName) {
        this.packageName = packageName;
        return this;
    }

    public String getEntity() {
        return entity;
    }

    public FormItemTemplateRespDto setEntity(String entity) {
        this.entity = entity;
        return this;
    }

    public String getEntityId() {
        return entityId;
    }

    public FormItemTemplateRespDto setEntityId(String entityId) {
        this.entityId = entityId;
        return this;
    }

    public String getAttrDefId() {
        return attrDefId;
    }

    public FormItemTemplateRespDto setAttrDefId(String attrDefId) {
        this.attrDefId = attrDefId;
        return this;
    }

    public String getAttrDefDataType() {
        return attrDefDataType;
    }

    public FormItemTemplateRespDto setAttrDefDataType(String attrDefDataType) {
        this.attrDefDataType = attrDefDataType;
        return this;
    }

    public String getRoutineExp() {
        return routineExp;
    }

    public FormItemTemplateRespDto setRoutineExp(String routineExp) {
        this.routineExp = routineExp;
        return this;
    }

    public String getElementType() {
        return elementType;
    }

    public FormItemTemplateRespDto setElementType(String elementType) {
        this.elementType = elementType;
        return this;
    }

    public String getTitle() {
        return title;
    }

    public FormItemTemplateRespDto setTitle(String title) {
        this.title = title;
        return this;
    }

    public Integer getWidth() {
        return width;
    }

    public FormItemTemplateRespDto setWidth(Integer width) {
        this.width = width;
        return this;
    }

    public String getRefEntity() {
        return refEntity;
    }

    public FormItemTemplateRespDto setRefEntity(String refEntity) {
        this.refEntity = refEntity;
        return this;
    }

    public String getRefPackageName() {
        return refPackageName;
    }

    public FormItemTemplateRespDto setRefPackageName(String refPackageName) {
        this.refPackageName = refPackageName;
        return this;
    }

    public String getRefFilters() {
        return refFilters;
    }

    public FormItemTemplateRespDto setRefFilters(String refFilters) {
        this.refFilters = refFilters;
        return this;
    }

    public String getDataOptions() {
        return dataOptions;
    }

    public FormItemTemplateRespDto setDataOptions(String dataOptions) {
        this.dataOptions = dataOptions;
        return this;
    }

    public Integer getRequired() {
        return required;
    }

    public FormItemTemplateRespDto setRequired(Integer required) {
        this.required = required;
        return this;
    }

    public String getRegular() {
        return regular;
    }

    public FormItemTemplateRespDto setRegular(String regular) {
        this.regular = regular;
        return this;
    }

    public Integer getIsEdit() {
        return isEdit;
    }

    public FormItemTemplateRespDto setIsEdit(Integer isEdit) {
        this.isEdit = isEdit;
        return this;
    }

    public Integer getIsView() {
        return isView;
    }

    public FormItemTemplateRespDto setIsView(Integer isView) {
        this.isView = isView;
        return this;
    }
}

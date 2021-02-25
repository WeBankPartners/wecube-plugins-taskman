package com.webank.taskman.dto.req;

import java.util.StringJoiner;

public class FormItemTemplateSaveReqDto {

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

    private String elementType;
    private String title;
    private Integer width;

    private String refEntity;
    private String refPackageName;
    private String refFilters;
    private String routineExp;

    private String dataOptions;

    private Integer required;
    private String regular;
    private Integer isEdit;
    private Integer isView;

    public String getFormTemplateId() {
        return formTemplateId;
    }

    public FormItemTemplateSaveReqDto setFormTemplateId(String formTemplateId) {
        this.formTemplateId = formTemplateId;
        return this;
    }

    public String getName() {
        return name;
    }

    public FormItemTemplateSaveReqDto setName(String name) {
        this.name = name;
        return this;
    }

    public String getDefaultValue() {
        return defaultValue;
    }

    public FormItemTemplateSaveReqDto setDefaultValue(String defaultValue) {
        this.defaultValue = defaultValue;
        return this;
    }

    public Integer getIsCurrency() {
        return isCurrency;
    }

    public FormItemTemplateSaveReqDto setIsCurrency(Integer isCurrency) {
        this.isCurrency = isCurrency;
        return this;
    }

    public Integer getSort() {
        return sort;
    }

    public FormItemTemplateSaveReqDto setSort(Integer sort) {
        this.sort = sort;
        return this;
    }

    public String getPackageName() {
        return packageName;
    }

    public FormItemTemplateSaveReqDto setPackageName(String packageName) {
        this.packageName = packageName;
        return this;
    }

    public String getEntity() {
        return entity;
    }

    public FormItemTemplateSaveReqDto setEntity(String entity) {
        this.entity = entity;
        return this;
    }

    public String getEntityId() {
        return entityId;
    }

    public FormItemTemplateSaveReqDto setEntityId(String entityId) {
        this.entityId = entityId;
        return this;
    }

    public String getAttrDefId() {
        return attrDefId;
    }

    public FormItemTemplateSaveReqDto setAttrDefId(String attrDefId) {
        this.attrDefId = attrDefId;
        return this;
    }

    public String getAttrDefDataType() {
        return attrDefDataType;
    }

    public FormItemTemplateSaveReqDto setAttrDefDataType(String attrDefDataType) {
        this.attrDefDataType = attrDefDataType;
        return this;
    }

    public String getElementType() {
        return elementType;
    }

    public FormItemTemplateSaveReqDto setElementType(String elementType) {
        this.elementType = elementType;
        return this;
    }

    public String getTitle() {
        return title;
    }

    public FormItemTemplateSaveReqDto setTitle(String title) {
        this.title = title;
        return this;
    }

    public Integer getWidth() {
        return width;
    }

    public FormItemTemplateSaveReqDto setWidth(Integer width) {
        this.width = width;
        return this;
    }

    public String getRefEntity() {
        return refEntity;
    }

    public FormItemTemplateSaveReqDto setRefEntity(String refEntity) {
        this.refEntity = refEntity;
        return this;
    }

    public String getRefPackageName() {
        return refPackageName;
    }

    public FormItemTemplateSaveReqDto setRefPackageName(String refPackageName) {
        this.refPackageName = refPackageName;
        return this;
    }

    public String getRefFilters() {
        return refFilters;
    }

    public FormItemTemplateSaveReqDto setRefFilters(String refFilters) {
        this.refFilters = refFilters;
        return this;
    }

    public String getRoutineExp() {
        return routineExp;
    }

    public FormItemTemplateSaveReqDto setRoutineExp(String routineExp) {
        this.routineExp = routineExp;
        return this;
    }

    public String getDataOptions() {
        return dataOptions;
    }

    public FormItemTemplateSaveReqDto setDataOptions(String dataOptions) {
        this.dataOptions = dataOptions;
        return this;
    }

    public Integer getRequired() {
        return required;
    }

    public FormItemTemplateSaveReqDto setRequired(Integer required) {
        this.required = required;
        return this;
    }

    public String getRegular() {
        return regular;
    }

    public FormItemTemplateSaveReqDto setRegular(String regular) {
        this.regular = regular;
        return this;
    }

    public Integer getIsEdit() {
        return isEdit;
    }

    public FormItemTemplateSaveReqDto setIsEdit(Integer isEdit) {
        this.isEdit = isEdit;
        return this;
    }

    public Integer getIsView() {
        return isView;
    }

    public FormItemTemplateSaveReqDto setIsView(Integer isView) {
        this.isView = isView;
        return this;
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", FormItemTemplateSaveReqDto.class.getSimpleName() + "[", "]")
                .add("formTemplateId='" + formTemplateId + "'").add("name='" + name + "'")
                .add("defaultValue='" + defaultValue + "'").add("isCurrency=" + isCurrency).add("sort=" + sort)
                .add("packageName='" + packageName + "'").add("entity='" + entity + "'")
                .add("attrDefId='" + attrDefId + "'").add("attrDefDataType='" + attrDefDataType + "'")
                .add("elementType='" + elementType + "'").add("title='" + title + "'").add("width=" + width)
                .add("refEntity='" + refEntity + "'").add("refPackageName='" + refPackageName + "'")
                .add("refFilters='" + refFilters + "'").add("dataOptions='" + dataOptions + "'")
                .add("required=" + required).add("regular='" + regular + "'").add("isEdit=" + isEdit)
                .add("isView=" + isView).toString();
    }
}

package com.webank.taskman.dto.resp;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

@ApiModel
public class FormItemTemplateResp {

    @ApiModelProperty(value = "",position = 1)
    private String id;
    @ApiModelProperty(value = "",position = 2)
    private String tempId;
    private String formType;
    @ApiModelProperty(value = "",position = 2)
    private String formTemplateId;
    @ApiModelProperty(value = "",position = 3)
    private String name;
    @ApiModelProperty(value = "",position = 4)
    private String defaultValue;
    @ApiModelProperty(value = "",position = 5)
    private Integer isCurrency;
    @ApiModelProperty(value = "",position = 6)
    private Integer sort;

    @ApiModelProperty(value = "",position = 7)
    private String packageName;
    @ApiModelProperty(value = "",position = 8)
    private String entity;
    @ApiModelProperty(value = "",position = 8)
    private String entityId;
    @ApiModelProperty(value = "",position = 9)
    private String attrDefId;
    @ApiModelProperty(value = "",position = 10)
    private String attrDefDataType;
    @ApiModelProperty(value = "",position = 10)
    private String routineExp;

    @ApiModelProperty(value = "",position = 11)
    private String elementType;
    @ApiModelProperty(value = "",position = 12)
    private String title;
    @ApiModelProperty(value = "",position = 13)
    private Integer width;

    @ApiModelProperty(value = "",position = 14)
    private String refEntity;
    @ApiModelProperty(value = "",position = 15)
    private String refPackageName;
    @ApiModelProperty(value = "",position = 16)
    private String refFilters;
    @ApiModelProperty(value = "",position = 17)
    private String dataOptions;

    @ApiModelProperty(value = "",position = 18)
    private Integer required;
    @ApiModelProperty(value = "",position = 19)
    private String regular;
    @ApiModelProperty(value = "",position = 20)
    private Integer isEdit;
    @ApiModelProperty(value = "",position = 21)
    private Integer isView;

    public String getId() {
        return id;
    }

    public FormItemTemplateResp setId(String id) {
        this.id = id;
        return this;
    }

    public String getTempId() {
        return tempId;
    }

    public FormItemTemplateResp setTempId(String tempId) {
        this.tempId = tempId;
        return this;
    }

    public String getFormType() {
        return formType;
    }

    public FormItemTemplateResp setFormType(String formType) {
        this.formType = formType;
        return this;
    }

    public String getFormTemplateId() {
        return formTemplateId;
    }

    public FormItemTemplateResp setFormTemplateId(String formTemplateId) {
        this.formTemplateId = formTemplateId;
        return this;
    }

    public String getName() {
        return name;
    }

    public FormItemTemplateResp setName(String name) {
        this.name = name;
        return this;
    }

    public String getDefaultValue() {
        return defaultValue;
    }

    public FormItemTemplateResp setDefaultValue(String defaultValue) {
        this.defaultValue = defaultValue;
        return this;
    }

    public Integer getIsCurrency() {
        return isCurrency;
    }

    public FormItemTemplateResp setIsCurrency(Integer isCurrency) {
        this.isCurrency = isCurrency;
        return this;
    }

    public Integer getSort() {
        return sort;
    }

    public FormItemTemplateResp setSort(Integer sort) {
        this.sort = sort;
        return this;
    }

    public String getPackageName() {
        return packageName;
    }

    public FormItemTemplateResp setPackageName(String packageName) {
        this.packageName = packageName;
        return this;
    }

    public String getEntity() {
        return entity;
    }

    public FormItemTemplateResp setEntity(String entity) {
        this.entity = entity;
        return this;
    }

    public String getEntityId() {
        return entityId;
    }

    public FormItemTemplateResp setEntityId(String entityId) {
        this.entityId = entityId;
        return this;
    }

    public String getAttrDefId() {
        return attrDefId;
    }

    public FormItemTemplateResp setAttrDefId(String attrDefId) {
        this.attrDefId = attrDefId;
        return this;
    }

    public String getAttrDefDataType() {
        return attrDefDataType;
    }

    public FormItemTemplateResp setAttrDefDataType(String attrDefDataType) {
        this.attrDefDataType = attrDefDataType;
        return this;
    }

    public String getRoutineExp() {
        return routineExp;
    }

    public FormItemTemplateResp setRoutineExp(String routineExp) {
        this.routineExp = routineExp;
        return this;
    }

    public String getElementType() {
        return elementType;
    }

    public FormItemTemplateResp setElementType(String elementType) {
        this.elementType = elementType;
        return this;
    }

    public String getTitle() {
        return title;
    }

    public FormItemTemplateResp setTitle(String title) {
        this.title = title;
        return this;
    }

    public Integer getWidth() {
        return width;
    }

    public FormItemTemplateResp setWidth(Integer width) {
        this.width = width;
        return this;
    }

    public String getRefEntity() {
        return refEntity;
    }

    public FormItemTemplateResp setRefEntity(String refEntity) {
        this.refEntity = refEntity;
        return this;
    }

    public String getRefPackageName() {
        return refPackageName;
    }

    public FormItemTemplateResp setRefPackageName(String refPackageName) {
        this.refPackageName = refPackageName;
        return this;
    }

    public String getRefFilters() {
        return refFilters;
    }

    public FormItemTemplateResp setRefFilters(String refFilters) {
        this.refFilters = refFilters;
        return this;
    }

    public String getDataOptions() {
        return dataOptions;
    }

    public FormItemTemplateResp setDataOptions(String dataOptions) {
        this.dataOptions = dataOptions;
        return this;
    }

    public Integer getRequired() {
        return required;
    }

    public FormItemTemplateResp setRequired(Integer required) {
        this.required = required;
        return this;
    }

    public String getRegular() {
        return regular;
    }

    public FormItemTemplateResp setRegular(String regular) {
        this.regular = regular;
        return this;
    }

    public Integer getIsEdit() {
        return isEdit;
    }

    public FormItemTemplateResp setIsEdit(Integer isEdit) {
        this.isEdit = isEdit;
        return this;
    }

    public Integer getIsView() {
        return isView;
    }

    public FormItemTemplateResp setIsView(Integer isView) {
        this.isView = isView;
        return this;
    }
}

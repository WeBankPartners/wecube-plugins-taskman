package com.webank.taskman.dto.resp;

import io.swagger.annotations.ApiModelProperty;

public class CoreTaskFormItemResp {

    @ApiModelProperty(value = "",position = 1)
    private String id;
    @ApiModelProperty(value = "",position = 2)
    private String formTemplateId;
    @ApiModelProperty(value = "",position = 3)
    private String attrDefId;
    @ApiModelProperty(value = "",position = 4)
    private String attrDataType;
    @ApiModelProperty(value = "",position = 5)
    private String name;
    @ApiModelProperty(value = "",position = 6)
    private String title;
    @ApiModelProperty(value = "",position = 7)
    private String elementType;
    @ApiModelProperty(value = "",position = 8)
    private String defaultValue;
    @ApiModelProperty(value = "",position = 9)
    private Integer required;
    @ApiModelProperty(value = "",position = 10)
    private Integer isEdit;
    @ApiModelProperty(value = "",position = 11)
    private Integer regular;
    @ApiModelProperty(value = "",position = 12)
    private Integer width;
    @ApiModelProperty(value = "",position = 13)
    private Integer sort;
    @ApiModelProperty(value = "",position = 14)
    private String entityId;
    @ApiModelProperty(value = "",position = 15)
    private String entityFilters;
    @ApiModelProperty(value = "",position = 16)
    private String dataOptions;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getFormTemplateId() {
        return formTemplateId;
    }

    public void setFormTemplateId(String formTemplateId) {
        this.formTemplateId = formTemplateId;
    }

    public String getAttrDefId() {
        return attrDefId;
    }

    public void setAttrDefId(String attrDefId) {
        this.attrDefId = attrDefId;
    }

    public String getAttrDataType() {
        return attrDataType;
    }

    public void setAttrDataType(String attrDataType) {
        this.attrDataType = attrDataType;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getElementType() {
        return elementType;
    }

    public void setElementType(String elementType) {
        this.elementType = elementType;
    }

    public String getDefaultValue() {
        return defaultValue;
    }

    public void setDefaultValue(String defaultValue) {
        this.defaultValue = defaultValue;
    }

    public Integer getRequired() {
        return required;
    }

    public void setRequired(Integer required) {
        this.required = required;
    }

    public Integer getIsEdit() {
        return isEdit;
    }

    public void setIsEdit(Integer isEdit) {
        this.isEdit = isEdit;
    }

    public Integer getRegular() {
        return regular;
    }

    public void setRegular(Integer regular) {
        this.regular = regular;
    }

    public Integer getWidth() {
        return width;
    }

    public void setWidth(Integer width) {
        this.width = width;
    }

    public Integer getSort() {
        return sort;
    }

    public void setSort(Integer sort) {
        this.sort = sort;
    }

    public String getEntityId() {
        return entityId;
    }

    public void setEntityId(String entityId) {
        this.entityId = entityId;
    }

    public String getEntityFilters() {
        return entityFilters;
    }

    public void setEntityFilters(String entityFilters) {
        this.entityFilters = entityFilters;
    }

    public String getDataOptions() {
        return dataOptions;
    }

    public void setDataOptions(String dataOptions) {
        this.dataOptions = dataOptions;
    }
}

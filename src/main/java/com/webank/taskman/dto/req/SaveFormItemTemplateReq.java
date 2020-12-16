package com.webank.taskman.dto.req;


import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import javax.validation.constraints.NotBlank;

@ApiModel
public class SaveFormItemTemplateReq {


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
    @ApiModelProperty(value = "",position = 9)
    private String attrDefId;
    @ApiModelProperty(value = "",position = 10)
    private String attrDefDataType;

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

    public String getFormTemplateId() {
        return formTemplateId;
    }

    public void setFormTemplateId(String formTemplateId) {
        this.formTemplateId = formTemplateId;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDefaultValue() {
        return defaultValue;
    }

    public void setDefaultValue(String defaultValue) {
        this.defaultValue = defaultValue;
    }

    public Integer getIsCurrency() {
        return isCurrency;
    }

    public void setIsCurrency(Integer isCurrency) {
        this.isCurrency = isCurrency;
    }

    public Integer getSort() {
        return sort;
    }

    public void setSort(Integer sort) {
        this.sort = sort;
    }

    public String getPackageName() {
        return packageName;
    }

    public void setPackageName(String packageName) {
        this.packageName = packageName;
    }

    public String getEntity() {
        return entity;
    }

    public void setEntity(String entity) {
        this.entity = entity;
    }

    public String getAttrDefId() {
        return attrDefId;
    }

    public void setAttrDefId(String attrDefId) {
        this.attrDefId = attrDefId;
    }

    public String getAttrDefDataType() {
        return attrDefDataType;
    }

    public void setAttrDefDataType(String attrDefDataType) {
        this.attrDefDataType = attrDefDataType;
    }

    public String getElementType() {
        return elementType;
    }

    public void setElementType(String elementType) {
        this.elementType = elementType;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public Integer getWidth() {
        return width;
    }

    public void setWidth(Integer width) {
        this.width = width;
    }

    public String getRefEntity() {
        return refEntity;
    }

    public void setRefEntity(String refEntity) {
        this.refEntity = refEntity;
    }

    public String getRefPackageName() {
        return refPackageName;
    }

    public void setRefPackageName(String refPackageName) {
        this.refPackageName = refPackageName;
    }

    public String getRefFilters() {
        return refFilters;
    }

    public void setRefFilters(String refFilters) {
        this.refFilters = refFilters;
    }

    public String getDataOptions() {
        return dataOptions;
    }

    public void setDataOptions(String dataOptions) {
        this.dataOptions = dataOptions;
    }

    public Integer getRequired() {
        return required;
    }

    public void setRequired(Integer required) {
        this.required = required;
    }

    public String getRegular() {
        return regular;
    }

    public void setRegular(String regular) {
        this.regular = regular;
    }

    public Integer getIsEdit() {
        return isEdit;
    }

    public void setIsEdit(Integer isEdit) {
        this.isEdit = isEdit;
    }

    public Integer getIsView() {
        return isView;
    }

    public void setIsView(Integer isView) {
        this.isView = isView;
    }
}

package com.webank.taskman.dto.req;


import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.util.StringJoiner;

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

    public SaveFormItemTemplateReq setFormTemplateId(String formTemplateId) {
        this.formTemplateId = formTemplateId;
        return this;
    }

    public String getName() {
        return name;
    }

    public SaveFormItemTemplateReq setName(String name) {
        this.name = name;
        return this;
    }

    public String getDefaultValue() {
        return defaultValue;
    }

    public SaveFormItemTemplateReq setDefaultValue(String defaultValue) {
        this.defaultValue = defaultValue;
        return this;
    }

    public Integer getIsCurrency() {
        return isCurrency;
    }

    public SaveFormItemTemplateReq setIsCurrency(Integer isCurrency) {
        this.isCurrency = isCurrency;
        return this;
    }

    public Integer getSort() {
        return sort;
    }

    public SaveFormItemTemplateReq setSort(Integer sort) {
        this.sort = sort;
        return this;
    }

    public String getPackageName() {
        return packageName;
    }

    public SaveFormItemTemplateReq setPackageName(String packageName) {
        this.packageName = packageName;
        return this;
    }

    public String getEntity() {
        return entity;
    }

    public SaveFormItemTemplateReq setEntity(String entity) {
        this.entity = entity;
        return this;
    }

    public String getAttrDefId() {
        return attrDefId;
    }

    public SaveFormItemTemplateReq setAttrDefId(String attrDefId) {
        this.attrDefId = attrDefId;
        return this;
    }

    public String getAttrDefDataType() {
        return attrDefDataType;
    }

    public SaveFormItemTemplateReq setAttrDefDataType(String attrDefDataType) {
        this.attrDefDataType = attrDefDataType;
        return this;
    }

    public String getElementType() {
        return elementType;
    }

    public SaveFormItemTemplateReq setElementType(String elementType) {
        this.elementType = elementType;
        return this;
    }

    public String getTitle() {
        return title;
    }

    public SaveFormItemTemplateReq setTitle(String title) {
        this.title = title;
        return this;
    }

    public Integer getWidth() {
        return width;
    }

    public SaveFormItemTemplateReq setWidth(Integer width) {
        this.width = width;
        return this;
    }

    public String getRefEntity() {
        return refEntity;
    }

    public SaveFormItemTemplateReq setRefEntity(String refEntity) {
        this.refEntity = refEntity;
        return this;
    }

    public String getRefPackageName() {
        return refPackageName;
    }

    public SaveFormItemTemplateReq setRefPackageName(String refPackageName) {
        this.refPackageName = refPackageName;
        return this;
    }

    public String getRefFilters() {
        return refFilters;
    }

    public SaveFormItemTemplateReq setRefFilters(String refFilters) {
        this.refFilters = refFilters;
        return this;
    }

    public String getDataOptions() {
        return dataOptions;
    }

    public SaveFormItemTemplateReq setDataOptions(String dataOptions) {
        this.dataOptions = dataOptions;
        return this;
    }

    public Integer getRequired() {
        return required;
    }

    public SaveFormItemTemplateReq setRequired(Integer required) {
        this.required = required;
        return this;
    }

    public String getRegular() {
        return regular;
    }

    public SaveFormItemTemplateReq setRegular(String regular) {
        this.regular = regular;
        return this;
    }

    public Integer getIsEdit() {
        return isEdit;
    }

    public SaveFormItemTemplateReq setIsEdit(Integer isEdit) {
        this.isEdit = isEdit;
        return this;
    }

    public Integer getIsView() {
        return isView;
    }

    public SaveFormItemTemplateReq setIsView(Integer isView) {
        this.isView = isView;
        return this;
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", SaveFormItemTemplateReq.class.getSimpleName() + "[", "]")
                .add("formTemplateId='" + formTemplateId + "'")
                .add("name='" + name + "'")
                .add("defaultValue='" + defaultValue + "'")
                .add("isCurrency=" + isCurrency)
                .add("sort=" + sort)
                .add("packageName='" + packageName + "'")
                .add("entity='" + entity + "'")
                .add("attrDefId='" + attrDefId + "'")
                .add("attrDefDataType='" + attrDefDataType + "'")
                .add("elementType='" + elementType + "'")
                .add("title='" + title + "'")
                .add("width=" + width)
                .add("refEntity='" + refEntity + "'")
                .add("refPackageName='" + refPackageName + "'")
                .add("refFilters='" + refFilters + "'")
                .add("dataOptions='" + dataOptions + "'")
                .add("required=" + required)
                .add("regular='" + regular + "'")
                .add("isEdit=" + isEdit)
                .add("isView=" + isView)
                .toString();
    }
}

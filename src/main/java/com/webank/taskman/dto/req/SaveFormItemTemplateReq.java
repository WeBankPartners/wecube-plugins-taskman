package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

@ApiModel
public class SaveFormItemTemplateReq {

    @ApiModelProperty(value = "",required = true)
    private String name;

    @ApiModelProperty(value = "",required = true)
    private String title;

    @ApiModelProperty(value = "",required = true)
    private String elementType;

    @ApiModelProperty(value = "",required = true)
    private Integer required;

    private Integer regular;

    private Integer isView;

    private Integer isEdit;

    private Integer width;

    private String defValue;

    @ApiModelProperty(value = "",required = true)
    private Integer sort;

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


    public Integer getRequired() {
        return required;
    }

    public void setRequired(Integer required) {
        this.required = required;
    }

    public Integer getRegular() {
        return regular;
    }

    public void setRegular(Integer regular) {
        this.regular = regular;
    }

    public Integer getIsView() {
        return isView;
    }

    public void setIsView(Integer isView) {
        this.isView = isView;
    }

    public Integer getIsEdit() {
        return isEdit;
    }

    public void setIsEdit(Integer isEdit) {
        this.isEdit = isEdit;
    }

    public Integer getWidth() {
        return width;
    }

    public void setWidth(Integer width) {
        this.width = width;
    }

    public String getDefValue() {
        return defValue;
    }

    public void setDefValue(String defValue) {
        this.defValue = defValue;
    }

    public Integer getSort() {
        return sort;
    }

    public void setSort(Integer sort) {
        this.sort = sort;
    }

}

package com.webank.taskman.dto;

import io.swagger.annotations.ApiModel;

@ApiModel
public class FormItemTemplateDTO {

    private String id;
    private String name;
    private String title;
    private String elementType;
    private Integer isPublic;
    private Integer required;
    private Integer regular;
    private Integer isView;
    private Integer isEdit;
    private Integer width;
    private String defValue;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
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

    public Integer getIsPublic() {
        return isPublic;
    }

    public void setIsPublic(Integer isPublic) {
        this.isPublic = isPublic;
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

    @Override
    public String toString() {
        return "FormItemTemplateDTO{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", title='" + title + '\'' +
                ", elementType='" + elementType + '\'' +
                ", isPublic=" + isPublic +
                ", required=" + required +
                ", regular=" + regular +
                ", isView=" + isView +
                ", isEdit=" + isEdit +
                ", width=" + width +
                ", defValue='" + defValue + '\'' +
                '}';
    }
}

package com.webank.taskman.dto.resp;

import io.swagger.annotations.ApiModel;

@ApiModel
public class FormItemTemplateResp {
    private String elementType;
    private String title;
    private Integer width;
    private Integer isEdit;
    private Integer isView;

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

    @Override
    public String toString() {
        return "FormItemTemplateResp{" +
                "elementType='" + elementType + '\'' +
                ", title='" + title + '\'' +
                ", width=" + width +
                ", isEdit=" + isEdit +
                ", isView=" + isView +
                '}';
    }
}

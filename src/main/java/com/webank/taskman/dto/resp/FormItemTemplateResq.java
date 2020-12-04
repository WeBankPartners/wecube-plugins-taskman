package com.webank.taskman.dto.resp;


import io.swagger.annotations.ApiModel;

@ApiModel
public class FormItemTemplateResq {
    private String id;

    private String formTemplateId;

    private String name;

    private String title;

    private String elementType;

    private String dataCiId;

    private Integer isPublic;

    private Integer width;

    private String defValue;

    private Integer sort;

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

    public String getDataCiId() {
        return dataCiId;
    }

    public void setDataCiId(String dataCiId) {
        this.dataCiId = dataCiId;
    }

    public Integer getIsPublic() {
        return isPublic;
    }

    public void setIsPublic(Integer isPublic) {
        this.isPublic = isPublic;
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

    @Override
    public String toString() {
        return "FormItemTemplateResq{" +
                "id='" + id + '\'' +
                ", formTemplateId='" + formTemplateId + '\'' +
                ", name='" + name + '\'' +
                ", title='" + title + '\'' +
                ", elementType='" + elementType + '\'' +
                ", dataCiId='" + dataCiId + '\'' +
                ", isPublic=" + isPublic +
                ", width=" + width +
                ", defValue='" + defValue + '\'' +
                ", sort=" + sort +
                '}';
    }
}

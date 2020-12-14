package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.webank.taskman.base.BaseEntity;
import io.swagger.annotations.ApiModelProperty;

import javax.validation.constraints.NotBlank;
import java.io.Serializable;

public class FormItemTemplate extends BaseEntity implements Serializable {

    private static final long serialVersionUID = 1L;

    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;
    private String formTemplateId;
    private String attrDefId;
    private String attrDataType;
    private String name;
    private String title;
    private String elementType;
    private String defaultValue;
    private Integer required;
    private Integer isEdit;
    private String regular;
    private Integer width;
    private Integer sort;
    private String entityId;
    private String entityFilters;
    private String dataOptions;
    private Integer isView;
    private Integer isPublic;

    public static long getSerialVersionUID() {
        return serialVersionUID;
    }

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

    public String getRegular() {
        return regular;
    }

    public void setRegular(String regular) {
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

    public Integer getIsView() {
        return isView;
    }

    public void setIsView(Integer isView) {
        this.isView = isView;
    }

    public Integer getIsPublic() {
        return isPublic;
    }

    public void setIsPublic(Integer isPublic) {
        this.isPublic = isPublic;
    }

    @Override
    public String toString() {
        return "FormItemTemplate{" +
                "id='" + id + '\'' +
                ", formTemplateId='" + formTemplateId + '\'' +
                ", attrDefId='" + attrDefId + '\'' +
                ", attrDataType='" + attrDataType + '\'' +
                ", name='" + name + '\'' +
                ", title='" + title + '\'' +
                ", elementType='" + elementType + '\'' +
                ", defaultValue='" + defaultValue + '\'' +
                ", required=" + required +
                ", isEdit=" + isEdit +
                ", regular=" + regular +
                ", width=" + width +
                ", sort=" + sort +
                ", entityId='" + entityId + '\'' +
                ", entityFilters='" + entityFilters + '\'' +
                ", dataOptions='" + dataOptions + '\'' +
                ", isView=" + isView +
                ", isPublic=" + isPublic +
                ", createdBy=" + getCreatedBy() +
                ", createdTime=" + getCreatedTime() +
                ", updatedBy=" + getUpdatedBy() +
                ", updatedTime=" + getUpdatedTime() +
                ", delFlag=" + getDelFlag() +
                '}';
    }
}

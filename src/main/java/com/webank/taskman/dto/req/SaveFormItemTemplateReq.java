package com.webank.taskman.dto.req;


import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import javax.validation.constraints.NotBlank;

@ApiModel
public class SaveFormItemTemplateReq {

    @ApiModelProperty(value = "Id",dataType = "String",position = 1)
    private String id;

    @NotBlank(message = "表单模板ID不能为空")
    @ApiModelProperty(value = "表单模板ID",dataType = "String",position = 2)
    private String formTemplateId;

    @ApiModelProperty(value = "",dataType = "String",position = 3)
    private String attrDefId;
    @ApiModelProperty(value = "",dataType = "String",position = 3)
    private String attrDataType;

    @NotBlank(message = "名称不能为空")
    @ApiModelProperty(value = "名称",dataType = "String",position = 4)
    private String name;

    @NotBlank(message = "标题不能为空")
    @ApiModelProperty(value = "标题",dataType = "String",position = 4)
    private String title;

    @NotBlank(message = "元素类型不能为空")
    @ApiModelProperty(value = "element_type",dataType = "String",position = 5)
    private String elementType;


    @ApiModelProperty(value = "默认值",dataType = "Integer",position = 6)
    private String defaultValue;

    @ApiModelProperty(value = "必填选项",dataType = "Integer",position = 7)
    private Integer required;

    @ApiModelProperty(value = "是否可编辑",dataType = "Integer",position = 8)
    private Integer isEdit;


    @ApiModelProperty(value = "正则表达式",dataType = "Integer",position = 9)
    private Integer regular;

    @ApiModelProperty(value = "长度",dataType = "Integer",position = 10)
    private Integer width;

    @ApiModelProperty(value = "排序",dataType = "Integer",position = 11)
    private Integer sort;

    @NotBlank(message = "ci数据id不能为空")
    @ApiModelProperty(value = "entityId",dataType = "String",position = 12)
    private String entityId;

    @ApiModelProperty(value = "ci数据检索条件",dataType = "String",position = 13)
    private String entityFilters;

    @ApiModelProperty(value = "自定义数据源选项",dataType = "String",position = 14)
    private String dataOptions;


    @ApiModelProperty(value = "是否显示",dataType = "Integer",hidden = true)
    private Integer isView;


    @ApiModelProperty(value = "是否通用",dataType = "Integer",hidden = true)
    private Integer isPublic;

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
}

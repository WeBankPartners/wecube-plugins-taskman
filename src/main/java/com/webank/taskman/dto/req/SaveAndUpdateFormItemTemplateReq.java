package com.webank.taskman.dto.req;


import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;

@ApiModel
public class SaveAndUpdateFormItemTemplateReq {
    @ApiModelProperty(value = "Id",dataType = "String")
    private String id;

    @NotBlank(message = "名称不能为空")
    @ApiModelProperty(value = "名称",dataType = "String")
    private String name;

    @NotBlank(message = "表单模板ID不能为空")
    @ApiModelProperty(value = "表单模板ID",dataType = "String")
    private String formTemplateId;

    @NotBlank(message = "标题不能为空")
    @ApiModelProperty(value = "标题",dataType = "String")
    private String title;

    @NotBlank(message = "元素类型不能为空")
    @ApiModelProperty(value = "element_type",dataType = "String")
    private String elementType;

    @NotBlank(message = "ci数据id不能为空")
    @ApiModelProperty(value = "data_ci_id",dataType = "String")
    private String dataCiId;

    @ApiModelProperty(value = "ci数据检索条件",dataType = "String")
    private String dataFilters;

    @ApiModelProperty(value = "自定义数据源选项",dataType = "String")
    private String dataOptions;

    @ApiModelProperty(value = "是否通用",dataType = "Integer")
    private Integer isPublic;


    @ApiModelProperty(value = "必填选项",dataType = "Integer")
    private Integer required;


    @ApiModelProperty(value = "是否可编辑",dataType = "Integer")
    private Integer isEdit;

    @ApiModelProperty(value = "正则表达式",dataType = "Integer")
    private Integer regular;


    @ApiModelProperty(value = "是否显示",dataType = "Integer")
    private Integer isView;

    @ApiModelProperty(value = "长度",dataType = "Integer")
    private Integer width;

    @ApiModelProperty(value = "默认值",dataType = "Integer")
    private String defValue;

    @ApiModelProperty(value = "排序",dataType = "Integer")
    private Integer sort;

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

    public String getFormTemplateId() {
        return formTemplateId;
    }

    public void setFormTemplateId(String formTemplateId) {
        this.formTemplateId = formTemplateId;
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

    public String getDataFilters() {
        return dataFilters;
    }

    public void setDataFilters(String dataFilters) {
        this.dataFilters = dataFilters;
    }

    public String getDataOptions() {
        return dataOptions;
    }

    public void setDataOptions(String dataOptions) {
        this.dataOptions = dataOptions;
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

    public Integer getIsView() {
        return isView;
    }

    public void setIsView(Integer isView) {
        this.isView = isView;
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
        return "SaveAndUpdateFormItemTemplateReq{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", formTemplateId='" + formTemplateId + '\'' +
                ", title='" + title + '\'' +
                ", elementType='" + elementType + '\'' +
                ", dataCiId='" + dataCiId + '\'' +
                ", dataFilters='" + dataFilters + '\'' +
                ", dataOptions='" + dataOptions + '\'' +
                ", isPublic=" + isPublic +
                ", required=" + required +
                ", isEdit=" + isEdit +
                ", regular=" + regular +
                ", isView=" + isView +
                ", width=" + width +
                ", defValue='" + defValue + '\'' +
                ", sort=" + sort +
                '}';
    }
}

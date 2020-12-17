package com.webank.taskman.dto.req;


import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import javax.validation.constraints.NotBlank;

@ApiModel
public class SelectFormItemTemplateReq {
    @ApiModelProperty(value = "Id",dataType = "String")
    private String id;

    @ApiModelProperty(value = "名称",dataType = "String")
    private String name;

    @ApiModelProperty(value = "表单模板ID",dataType = "String")
    private String formTemplateId;

    @ApiModelProperty(value = "标题",dataType = "String")
    private String title;

    @ApiModelProperty(value = "element_type",dataType = "String")
    private String elementType;

    @ApiModelProperty(value = "data_ci_id",dataType = "String")
    private String dataCiId;

    @ApiModelProperty(value = "ci数据检索条件",dataType = "String")
    private String dataFilters;

    @ApiModelProperty(value = "自定义数据源选项",dataType = "String")
    private String dataOptions;

    @ApiModelProperty(value = "是否通用",dataType = "Integer")
    private Integer isPublic;

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

    public Integer getSort() {
        return sort;
    }

    public void setSort(Integer sort) {
        this.sort = sort;
    }

    @Override
    public String toString() {
        return "SelectFormItemTemplateReq{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", formTemplateId='" + formTemplateId + '\'' +
                ", title='" + title + '\'' +
                ", elementType='" + elementType + '\'' +
                ", dataCiId='" + dataCiId + '\'' +
                ", dataFilters='" + dataFilters + '\'' +
                ", dataOptions='" + dataOptions + '\'' +
                ", isPublic=" + isPublic +
                ", sort=" + sort +
                '}';
    }
}

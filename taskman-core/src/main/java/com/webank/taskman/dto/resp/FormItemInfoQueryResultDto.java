package com.webank.taskman.dto.resp;

import com.fasterxml.jackson.annotation.JsonInclude;

@JsonInclude(JsonInclude.Include.NON_NULL)
public class FormItemInfoQueryResultDto {
    private String id;

    private String value;

    private String itemTempId;

    private String elementType;

    private String title;

    private String name;

    private Integer width;

    private Integer isEdit;

    private Integer isView;

    private Integer sort;

    private String dataOptions;

    private String routineExp;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getValue() {
        return value;
    }

    public void setValue(String value) {
        this.value = value;
    }

    public String getItemTempId() {
        return itemTempId;
    }

    public void setItemTempId(String itemTempId) {
        this.itemTempId = itemTempId;
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

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
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

    public Integer getSort() {
        return sort;
    }

    public void setSort(Integer sort) {
        this.sort = sort;
    }

    public String getDataOptions() {
        return dataOptions;
    }

    public void setDataOptions(String dataOptions) {
        this.dataOptions = dataOptions;
    }

    public String getRoutineExp() {
        return routineExp;
    }

    public FormItemInfoQueryResultDto setRoutineExp(String routineExp) {
        this.routineExp = routineExp;
        return this;
    }
}

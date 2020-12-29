package com.webank.taskman.dto.resp;

import com.webank.taskman.domain.FormItemInfo;

import java.util.List;

public class ProcessingTasksResp {
    private String id;

    private String recordId;

    private String formTemplateId;

    private String name;

    private Integer type;

    private List<FormItemInfo> formItemInfoList;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getRecordId() {
        return recordId;
    }

    public void setRecordId(String recordId) {
        this.recordId = recordId;
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

    public Integer getType() {
        return type;
    }

    public void setType(Integer type) {
        this.type = type;
    }

    public List<FormItemInfo> getFormItemInfoList() {
        return formItemInfoList;
    }

    public void setFormItemInfoList(List<FormItemInfo> formItemInfoList) {
        this.formItemInfoList = formItemInfoList;
    }

    @Override
    public String toString() {
        return "ProcessingTasksResp{" +
                "id='" + id + '\'' +
                ", recordId='" + recordId + '\'' +
                ", formTemplateId='" + formTemplateId + '\'' +
                ", name='" + name + '\'' +
                ", type=" + type +
                ", formItemInfoList=" + formItemInfoList +
                '}';
    }
}

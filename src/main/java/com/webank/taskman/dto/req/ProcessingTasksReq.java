package com.webank.taskman.dto.req;


import java.util.List;

public class ProcessingTasksReq {
    private String recordId;

    private String result;

    private List<FormItemInfoReq> formItemInfoList;

    public String getRecordId() {
        return recordId;
    }

    public void setRecordId(String recordId) {
        this.recordId = recordId;
    }

    public String getResult() {
        return result;
    }

    public void setResult(String result) {
        this.result = result;
    }

    public List<FormItemInfoReq> getFormItemInfoList() {
        return formItemInfoList;
    }

    public void setFormItemInfoList(List<FormItemInfoReq> formItemInfoList) {
        this.formItemInfoList = formItemInfoList;
    }

    @Override
    public String toString() {
        return "ProcessingTasksReq{" +
                "recordId='" + recordId + '\'' +
                ", formItemInfoList=" + formItemInfoList +
                '}';
    }
}

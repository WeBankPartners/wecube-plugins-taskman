package com.webank.taskman.dto.req;


import java.util.List;

public class ProcessingTasksReq {

    public static final String RESULT_SUCCESSFUL="Successful/Approved";
    public static final String RESULT_FAILED="Failed/Rejected";


    private String recordId;
    private String result;
    private String resultMessage;

    private List<FormItemInfoRequestDto> formItemInfoList;

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

    public String getResultMessage() {
        return resultMessage;
    }

    public ProcessingTasksReq setResultMessage(String resultMessage) {
        this.resultMessage = resultMessage;
        return this;
    }

    public List<FormItemInfoRequestDto> getFormItemInfoList() {
        return formItemInfoList;
    }

    public void setFormItemInfoList(List<FormItemInfoRequestDto> formItemInfoList) {
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

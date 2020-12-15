package com.webank.taskman.dto.resp;

import com.webank.taskman.domain.FormItemInfo;

import java.util.List;

public class SynthesisRequestInfoResp {
    private String id;


    private String requestTempId;


    private String procInstKey;


    private String name;


    private String status;

    private List<FormItemInfo> formItemInfo;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getRequestTempId() {
        return requestTempId;
    }

    public void setRequestTempId(String requestTempId) {
        this.requestTempId = requestTempId;
    }

    public String getProcInstKey() {
        return procInstKey;
    }

    public void setProcInstKey(String procInstKey) {
        this.procInstKey = procInstKey;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public List<FormItemInfo> getFormItemInfo() {
        return formItemInfo;
    }

    public void setFormItemInfo(List<FormItemInfo> formItemInfo) {
        this.formItemInfo = formItemInfo;
    }
}

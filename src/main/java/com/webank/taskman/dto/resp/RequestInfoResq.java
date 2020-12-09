package com.webank.taskman.dto.resp;

import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.domain.FormItemInfo;

import java.util.List;

public class RequestInfoResq {
    private String id;

    private String requestTempId;

    private String procInstKey;

    private String name;

    private Integer status;

    private FormInfoResq formInfoResq;

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

    public Integer getStatus() {
        return status;
    }

    public void setStatus(Integer status) {
        this.status = status;
    }

    public FormInfoResq getFormInfoResq() {
        return formInfoResq;
    }

    public void setFormInfoResq(FormInfoResq formInfoResq) {
        this.formInfoResq = formInfoResq;
    }
}
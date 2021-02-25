package com.webank.taskman.support.core.dto;

import java.util.LinkedHashMap;
import java.util.List;

public class CoreResponse<DATATYPE> {

    private String status;
    private String message;
    private DATATYPE data;

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public DATATYPE getData() {
        return data;

    }

    @SuppressWarnings("unchecked")
    public void setData(Object data) {
        this.data = (DATATYPE) data;
    }

    @Override
    public String toString() {
        return "CoreResponse{" + "status='" + status + '\'' + ", message='" + message + '\'' + ", data="
                + (null != data ? data.getClass() : "[]") +
                // ", data=" + data.toString() +
                '}';
    }

    public static class DefaultCoreResponse extends CoreResponse<Object> {
    }

    @SuppressWarnings("rawtypes")
    public static class LinkedHashMapResponse extends CoreResponse<LinkedHashMap> {
    }

    public static class ListDataResponse extends CoreResponse<List<Object>> {
    }

}
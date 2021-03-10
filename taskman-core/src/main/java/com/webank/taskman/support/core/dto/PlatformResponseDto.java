package com.webank.taskman.support.core.dto;

import java.util.LinkedHashMap;
import java.util.List;

public class PlatformResponseDto<DATATYPE> {

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


    public static class DefaultCoreResponse extends PlatformResponseDto<Object> {
    }

    @SuppressWarnings("rawtypes")
    public static class LinkedHashMapResponse extends PlatformResponseDto<LinkedHashMap> {
    }

    public static class ListDataResponse extends PlatformResponseDto<List<Object>> {
    }

    @Override
    public String toString() {
        StringBuilder builder = new StringBuilder();
        builder.append("PlatformResponseDto [status=");
        builder.append(status);
        builder.append(", message=");
        builder.append(message);
        builder.append(", data=");
        builder.append(data);
        builder.append("]");
        return builder.toString();
    }

}
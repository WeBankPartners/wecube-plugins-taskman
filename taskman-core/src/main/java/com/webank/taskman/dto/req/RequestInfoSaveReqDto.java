package com.webank.taskman.dto.req;

import java.util.List;

public class RequestInfoSaveReqDto {

    private String id;

    private String requestTempId;

    private String rootEntity;

    private String emergency;

    private String name;

    private String description;

    private String status;

    private List<FormItemInfoRequestDto> formItems;

    public String getId() {
        return id;
    }

    public RequestInfoSaveReqDto setId(String id) {
        this.id = id;
        return this;
    }

    public String getRequestTempId() {
        return requestTempId;
    }

    public RequestInfoSaveReqDto setRequestTempId(String requestTempId) {
        this.requestTempId = requestTempId;
        return this;
    }

    public String getRootEntity() {
        return rootEntity;
    }

    public RequestInfoSaveReqDto setRootEntity(String rootEntity) {
        this.rootEntity = rootEntity;
        return this;
    }

    public String getEmergency() {
        return emergency;
    }

    public RequestInfoSaveReqDto setEmergency(String emergency) {
        this.emergency = emergency;
        return this;
    }

    public String getName() {
        return name;
    }

    public RequestInfoSaveReqDto setName(String name) {
        this.name = name;
        return this;
    }

    public String getDescription() {
        return description;
    }

    public RequestInfoSaveReqDto setDescription(String description) {
        this.description = description;
        return this;
    }

    public String getStatus() {
        return status;
    }

    public RequestInfoSaveReqDto setStatus(String status) {
        this.status = status;
        return this;
    }

    public List<FormItemInfoRequestDto> getFormItems() {
        return formItems;
    }

    public RequestInfoSaveReqDto setFormItems(List<FormItemInfoRequestDto> formItems) {
        this.formItems = formItems;
        return this;
    }
}

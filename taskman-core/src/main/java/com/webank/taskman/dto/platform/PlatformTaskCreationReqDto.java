package com.webank.taskman.dto.platform;

import java.util.List;

public class PlatformTaskCreationReqDto {

    private String requestId;
    private String dueDate;
    private List<String> allowedOptions;
    private List<PlatformTaskInfoDto> inputs;

    public String getRequestId() {
        return requestId;
    }

    public PlatformTaskCreationReqDto setRequestId(String requestId) {
        this.requestId = requestId;
        return this;
    }

    public List<String> getAllowedOptions() {
        return allowedOptions;
    }

    public void setAllowedOptions(List<String> allowedOptions) {
        this.allowedOptions = allowedOptions;
    }

    public String getDueDate() {
        return dueDate;
    }

    public void setDueDate(String dueDate) {
        this.dueDate = dueDate;
    }

    public List<PlatformTaskInfoDto> getInputs() {
        return inputs;
    }

    public void setInputs(List<PlatformTaskInfoDto> inputs) {
        this.inputs = inputs;
    }

}

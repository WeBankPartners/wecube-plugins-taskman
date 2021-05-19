package com.webank.taskman.support.platform.dto;

import java.util.ArrayList;
import java.util.List;

import com.fasterxml.jackson.annotation.JsonInclude;

@JsonInclude(JsonInclude.Include.NON_NULL)
public class PlatformPluginRequestDto {
    private String requestId;
    private String dueDate;
    private List<String> allowedOptions;
    private List<CreateTaskRequestInputDto> inputs;

    public PlatformPluginRequestDto withInputs(List<CreateTaskRequestInputDto> inputs) {
        this.inputs = inputs;
        return this;
    }

    public PlatformPluginRequestDto withRequestId(String requestId) {
        this.requestId = requestId;
        return this;
    }

    public PlatformPluginRequestDto withAllowedOptions(List<String> allowedOptions) {
        if (this.allowedOptions == null) {
            this.allowedOptions = new ArrayList<>();
        }
        this.allowedOptions.addAll(allowedOptions);
        return this;
    }

    public PlatformPluginRequestDto withDueDate(String dueDate) {
        this.dueDate = dueDate;
        return this;
    }

    public String getRequestId() {
        return requestId;
    }

    public void setRequestId(String requestId) {
        this.requestId = requestId;
    }

    public List<CreateTaskRequestInputDto> getInputs() {
        return inputs;
    }

    public void setInputs(List<CreateTaskRequestInputDto> inputs) {
        this.inputs = inputs;
    }

    public List<String> getAllowedOptions() {
        return allowedOptions;
    }

    public void setAllowedOptions(List<String> allowedOptions) {
        this.allowedOptions = allowedOptions;
    }

    public void setDueDate(String dueDate) {
        this.dueDate = dueDate;
    }

    public String getDueDate() {
        return dueDate;
    }

    @Override
    public String toString() {
        StringBuilder builder = new StringBuilder();
        builder.append("PluginRequest [requestId=");
        builder.append(requestId);
        builder.append(", dueDate=");
        builder.append(dueDate);
        builder.append(", allowedOptions=");
        builder.append(allowedOptions);
        builder.append(", inputs=");
        builder.append(inputs);
        builder.append("]");
        return builder.toString();
    }

}

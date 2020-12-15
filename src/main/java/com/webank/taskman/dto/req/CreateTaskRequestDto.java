package com.webank.taskman.dto.req;

import java.util.List;

public class CreateTaskRequestDto {
    private String requestId;
    private List<CreateTaskRequestInputDto> inputs;
    private String dueDate;

    public List<CreateTaskRequestInputDto> getInputs() {
        return inputs;
    }

    public void setInputs(List<CreateTaskRequestInputDto> inputs) {
        this.inputs = inputs;
    }

    public String getRequestId() {
        return requestId;
    }

    public void setRequestId(String requestId) {
        this.requestId = requestId;
    }

    public String getDueDate() { return dueDate; }

    public void setDueDate(String dueDate) { this.dueDate = dueDate; }


}

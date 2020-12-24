package com.webank.taskman.dto.resp;

import io.swagger.annotations.ApiModel;

@ApiModel
public class TaskInfoGetResp {
    private String id;
    private String status;
    private String reporter;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public String getReporter() {
        return reporter;
    }

    public void setReporter(String reporter) {
        this.reporter = reporter;
    }

    @Override
    public String toString() {
        return "TaskInfoGetResp{" +
                "id='" + id + '\'' +
                ", status='" + status + '\'' +
                ", reporter='" + reporter + '\'' +
                '}';
    }
}

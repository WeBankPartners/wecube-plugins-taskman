package com.webank.taskman.dto.req;


import io.swagger.annotations.ApiModel;

@ApiModel
public class SelectTaskInfoReq {

    private String id;

    private String taskTempId;

    private String name;

    private Integer status;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getTaskTempId() {
        return taskTempId;
    }

    public void setTaskTempId(String taskTempId) {
        this.taskTempId = taskTempId;
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

    @Override
    public String toString() {
        return "SelectTaskInfoReq{" +
                "id='" + id + '\'' +
                ", taskTempId='" + taskTempId + '\'' +
                ", name='" + name + '\'' +
                ", status=" + status +
                '}';
    }
}

package com.webank.taskman.dto.resp;


import io.swagger.annotations.ApiModel;


@ApiModel
public class TaskInfoInstanceResp {
    private String id;
    private String name;
    private String status;
    private String result;

    private FormInfoResq taskFormResq;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
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

    public String getResult() {
        return result;
    }

    public void setResult(String result) {
        this.result = result;
    }

    public FormInfoResq getTaskFormResq() {
        return taskFormResq;
    }

    public void setTaskFormResq(FormInfoResq taskFormResq) {
        this.taskFormResq = taskFormResq;
    }

    @Override
    public String toString() {
        return "TaskInfoInstanceResp{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", status='" + status + '\'' +
                ", result='" + result + '\'' +
                ", taskFormResq=" + taskFormResq +
                '}';
    }
}

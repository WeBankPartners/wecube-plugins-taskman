package com.webank.taskman.dto.resp;

public class TaskInfoInstanceRespDto {
    private String id;
    private String name;
    private String status;
    private String result;

    private FormInfoResqDto taskFormResq;

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

    public FormInfoResqDto getTaskFormResq() {
        return taskFormResq;
    }

    public void setTaskFormResq(FormInfoResqDto taskFormResq) {
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

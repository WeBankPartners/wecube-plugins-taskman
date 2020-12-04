package com.webank.taskman.dto.resp;


import io.swagger.annotations.ApiModel;

@ApiModel
public class TaskTemplateSVResp {
    private String id;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    @Override
    public String toString() {
        return "TaskTemplateSVResp{" +
                "id='" + id + '\'' +
                '}';
    }
}

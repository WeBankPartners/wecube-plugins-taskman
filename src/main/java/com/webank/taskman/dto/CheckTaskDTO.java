package com.webank.taskman.dto;

import com.webank.taskman.dto.resp.FormInfoResq;

public class CheckTaskDTO {
    private FormInfoResq formInfoResq;

    public FormInfoResq getFormInfoResq() {
        return formInfoResq;
    }

    public void setFormInfoResq(FormInfoResq formInfoResq) {
        this.formInfoResq = formInfoResq;
    }

    @Override
    public String toString() {
        return "CheckTaskDTO{" +
                "formInfoResq=" + formInfoResq +
                '}';
    }
}

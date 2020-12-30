package com.webank.taskman.constant;

public enum  StatusCodeEnum {

    SAVE_ERROR("1000","Save error:%s","Save error:%s"),
    PARAM_ISNULL("1001","param is null ","param is null"),
    NOT_FOUND_RECORD("3007","Record does not exist ","Record does not exist"),
    ;
    private String code;

    private String message;

    private String description;


    StatusCodeEnum(String code, String message, String description) {
        this.code = code;
        this.message = message;
        this.description = description;
    }

    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public StatusCodeEnum setErrorMsg(String msg){
        this.message = String.format(message,msg);
        return this;
    }
}

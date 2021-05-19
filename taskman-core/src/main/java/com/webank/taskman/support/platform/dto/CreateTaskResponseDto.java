package com.webank.taskman.support.platform.dto;

public class CreateTaskResponseDto {
    public final static String STATUS_OK = "0";
    public final static String STATUS_ERROR = "1";
    private String callbackParameter;
    private String errorCode;
    private String errorMessage;
    private Object output;

    public String getCallbackParameter() {
        return callbackParameter;
    }

    public void setCallbackParameter(String callbackParameter) {
        this.callbackParameter = callbackParameter;
    }

    public String getErrorCode() {
        return errorCode;
    }

    public void setErrorCode(String errorCode) {
        this.errorCode = errorCode;
    }

    public String getErrorMessage() {
        return errorMessage;
    }

    public void setErrorMessage(String errorMessage) {
        this.errorMessage = errorMessage;
    }

    public Object getOutput() {
        return output;
    }

    public void setOutput(Object output) {
        this.output = output;
    }

}

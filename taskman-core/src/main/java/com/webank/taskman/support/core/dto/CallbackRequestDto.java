package com.webank.taskman.support.core.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

public class CallbackRequestDto {
    @JsonProperty(value = "resultCode")
    private String resultCode;
    @JsonProperty(value = "resultMessage")
    private String resultMessage;
    private CallbackRequestResultDataDto results;

    public String getResultCode() {
        return resultCode;
    }

    public void setResultCode(String resultCode) {
        this.resultCode = resultCode;
    }

    public String getResultMessage() {
        return resultMessage;
    }

    public void setResultMessage(String resultMessage) {
        this.resultMessage = resultMessage;
    }

    public CallbackRequestResultDataDto getResults() {
        return results;
    }

    public void setResults(CallbackRequestResultDataDto results) {
        this.results = results;
    }

    @Override
    public String toString() {
        return "CallbackRequestDto [resultCode=" + resultCode + ", resultMessage=" + resultMessage + ", results="
                + results + "]";
    }

    public static class CallbackRequestResultDataDto {
        private String requestId;
        private List<CallbackRequestOutputsDto> outputs;

        public String getRequestId() {
            return requestId;
        }

        public void setRequestId(String requestId) {
            this.requestId = requestId;
        }

        public List<CallbackRequestOutputsDto> getOutputs() {
            return outputs;
        }

        public void setOutputs(List<CallbackRequestOutputsDto> outputs) {
            this.outputs = outputs;
        }

        @Override
        public String toString() {
            return "CallbackRequestResultDataDto [requestId=" + requestId + ", outputs=" + outputs + "]";
        }

    }
    public static class CallbackRequestOutputsDto {
        public static final String ERROR_CODE_SUCCESSFUL = "0";
        public static final String ERROR_CODE_FAILED = "1";

        private String errorCode;
        private String errorMessage;
        private String comment;
        private String callbackParameter;

        public CallbackRequestOutputsDto() {
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

        public String getComment() {
            return comment;
        }

        public void setComment(String comment) {
            this.comment = comment;
        }

        public CallbackRequestOutputsDto(String errorCode, String errorMessage, String comment, String callbackParameter) {
            super();
            this.errorCode = errorCode;
            this.errorMessage = errorMessage;
            this.comment = comment;
            this.callbackParameter = callbackParameter;
        }

        public String getCallbackParameter() {
            return callbackParameter;
        }

        public void setCallbackParameter(String callbackParameter) {
            this.callbackParameter = callbackParameter;
        }

        @Override
        public String toString() {
            return "CallbackRequestOutputsDto [errorCode=" + errorCode + ", errorMessage=" + errorMessage + ", comment="
                    + comment + ", callbackParameter=" + callbackParameter + "]";
        }
    }
}

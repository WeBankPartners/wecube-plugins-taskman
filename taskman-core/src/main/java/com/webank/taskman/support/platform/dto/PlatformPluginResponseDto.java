package com.webank.taskman.support.platform.dto;

import java.util.List;

import com.fasterxml.jackson.annotation.JsonProperty;

public class PlatformPluginResponseDto {
    public static final String RESULT_CODE_OK = "0";
    public static final String RESULT_CODE_FAIL = "1";

    @JsonProperty("resultCode")
    private String resultCode;
    @JsonProperty("resultMessage")
    private String resultMessage;
    @JsonProperty("results")
    private ResultData resultData;

    public List<Object> getOutputs() {
        return (resultData == null) ? null : resultData.getOutputs();
    }

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

    public ResultData getResultData() {
        return resultData;
    }

    public void setResultData(ResultData resultData) {
        this.resultData = resultData;
    }
    
    @Override
    public String toString() {
        return "PluginResponse [resultCode=" + resultCode + ", resultMessage=" + resultMessage + ", resultData="
                + resultData + "]";
    }

    public static class ResultData {
        private List<Object> outputs;

        public List<Object> getOutputs() {
            return outputs;
        }

        public void setOutputs(List<Object> outputs) {
            this.outputs = outputs;
        }

        public ResultData() {
        }

        public ResultData(List<Object> outputs) {
            super();
            this.outputs = outputs;
        }

        @Override
        public String toString() {
            return "ResultData [outputs=" + outputs + "]";
        }
    }
}

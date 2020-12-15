package com.webank.taskman.dto;

import java.util.List;

public class CoreCreateTaskResp {


    /**
     * resultCode :
     * resultMessage :
     * results : {"outputs":[{"callbackParameter":"","errorCode":"","errorMessage":"","output":{}}]}
     */

    private String resultCode;
    private String resultMessage;
    private ResultsBean results;

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

    public ResultsBean getResults() {
        return results;
    }

    public void setResults(ResultsBean results) {
        this.results = results;
    }

    public static class ResultsBean {
        private List<OutputsBean> outputs;

        public List<OutputsBean> getOutputs() {
            return outputs;
        }

        public void setOutputs(List<OutputsBean> outputs) {
            this.outputs = outputs;
        }

        public static class OutputsBean {
            /**
             * callbackParameter :
             * errorCode :
             * errorMessage :
             * output : {}
             */

            private String callbackParameter;
            private String errorCode;
            private String errorMessage;
            private OutputBean output;

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

            public OutputBean getOutput() {
                return output;
            }

            public void setOutput(OutputBean output) {
                this.output = output;
            }

            public static class OutputBean {
                private String status ="OK";
                private String message;
                private String code="200";
                private String data;

                public String getStatus() {
                    return status;
                }

                public void setStatus(String status) {
                    this.status = status;
                }

                public String getMessage() {
                    return message;
                }

                public void setMessage(String message) {
                    this.message = message;
                }

                public String getCode() {
                    return code;
                }

                public void setCode(String code) {
                    this.code = code;
                }

                public String getData() {
                    return data;
                }

                public void setData(String data) {
                    this.data = data;
                }
            }
        }
    }
}

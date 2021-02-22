package com.webank.taskman.base;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.StringJoiner;

import com.webank.taskman.constant.BizCodeEnum;
import com.webank.taskman.constant.StatusCodeEnum;

public class JsonResponse implements Serializable {
    /**
     * 
     */
    private static final long serialVersionUID = -3311641576380673663L;
    public final static String STATUS_OK = "OK";
    public final static String STATUS_ERROR = "ERROR";

    private String status;
    private String message;
    private Integer code;
    private String codeMessage;
    private Object data;

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

    public Object getData() {
        return data;
    }

    public void setData(Object data) {
        this.data = null != data ? data : (Object) new ArrayList<Object>();
    }

    public Integer getCode() {
        return code;
    }

    public void setCode(Integer code) {
        this.code = code;
    }

    public String getCodeMessage() {
        return codeMessage;
    }

    public void setCodeMessage(String codeMessage) {
        this.codeMessage = codeMessage;
    }

    public JsonResponse withData(Object data) {
        setData(data);
        return this;
    }

    public JsonResponse() {
    }

    public JsonResponse(String status, String message, Integer code, String codeMessage, Object data) {
        this.status = status;
        this.message = message;
        this.code = code;
        this.codeMessage = codeMessage;
        setData(data);
    }

    public static JsonResponse okay() {
        return new JsonResponse(STATUS_OK, "Success", 200, "", null);
    }

    public static JsonResponse okayWithData(Object data) {
        return okay().withData(data);
    }

    public static JsonResponse customError(String errorMessage) {
        return new JsonResponse(STATUS_ERROR, errorMessage, null, null, null);
    }

    public static JsonResponse customError(BizCodeEnum bizCodeEnum, Object data) {
        return new JsonResponse(STATUS_ERROR, null, bizCodeEnum.getCode(), bizCodeEnum.getMessage(), data);
    }

    public static JsonResponse customError(Integer code, String codeMessage, Object data) {
        return new JsonResponse(STATUS_ERROR, null, code, codeMessage, data);
    }

    public static JsonResponse customError(Integer code, String codeMessage, String errorMessage, Object data) {
        return new JsonResponse(STATUS_ERROR, codeMessage, code, errorMessage, data);
    }

    public static JsonResponse customError(BizCodeEnum bizCodeEnum, String errorMessage) {
        return new JsonResponse(STATUS_ERROR, bizCodeEnum.getMessage(), bizCodeEnum.getCode(), errorMessage, null);
    }

    public static JsonResponse customError(BizCodeEnum bizCodeEnum, String message, String errorMessage) {
        return new JsonResponse(STATUS_ERROR, message, bizCodeEnum.getCode(), errorMessage, null);
    }

    public static JsonResponse customError(StatusCodeEnum statusCodeEnum) {
        return new JsonResponse(STATUS_ERROR, statusCodeEnum.getMessage(), Integer.valueOf(statusCodeEnum.getCode()),
                statusCodeEnum.getMessage(), null);
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", JsonResponse.class.getSimpleName() + "[", "]").add("status='" + status + "'")
                .add("message='" + message + "'").add("code=" + code).add("codeMessage='" + codeMessage + "'")
                .add("data=" + data.getClass().getTypeName()).toString();
    }
}

package com.webank.taskman.dto;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

@ApiModel(value = "",description = "JsonResponse")
public class JsonResponse<T> {
	public final static String STATUS_OK = "OK";
	public final static String STATUS_ERROR = "ERROR";

	@ApiModelProperty(value = "状态")
	private String status;
	@ApiModelProperty(value = "消息")
	private String message;
	@ApiModelProperty(value = "自定义状态码")
	private Integer code;
	@ApiModelProperty(value = "错误消息")
	private String codeMessage;

	@ApiModelProperty(value = "数据")
	private T data;

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

	public T getData() {
		return data;
	}

	public void setData(T data) {
		this.data = data;
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

	public JsonResponse withData(T data) {
		this.data = data;
		return this;
	}

	public static JsonResponse okay() {
		JsonResponse result = new JsonResponse();
		result.setStatus(STATUS_OK);
		result.setMessage("Success");
		result.setData(200);
		return result;
	}


	public static JsonResponse okayWithData(Object data) {
		return okay().withData(data);
	}

	public static JsonResponse error(String errorMessage) {
		JsonResponse result = new JsonResponse();
		result.setStatus(STATUS_ERROR);
		result.setMessage(errorMessage);
		return result;
	}

	public static JsonResponse customError(Integer code, String codeMessage, Object data) {
		JsonResponse result = new JsonResponse();
		result.setCode(code);
		result.setCodeMessage(codeMessage);
		result.setData(data);
		return result;
	}


}

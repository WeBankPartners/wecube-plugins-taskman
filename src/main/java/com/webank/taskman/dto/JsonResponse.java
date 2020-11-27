package com.webank.taskman.dto;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

@ApiModel(value = "返回对象",description = "JsonResponse")
public class JsonResponse {
	public final static String STATUS_OK = "OK";
	public final static String STATUS_ERROR = "ERROR";

	@ApiModelProperty(value = "状态")
	private String status;
	@ApiModelProperty(value = "消息")
	private String message;
	@ApiModelProperty(value = "数据")
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
		this.data = data;
	}

	public JsonResponse withData(Object data) {
		this.data = data;
		return this;
	}

	public static JsonResponse okay() {
		JsonResponse result = new JsonResponse();
		result.setStatus(STATUS_OK);
		result.setMessage("Success");
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
}

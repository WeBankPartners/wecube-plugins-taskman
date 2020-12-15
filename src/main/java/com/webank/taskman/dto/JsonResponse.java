package com.webank.taskman.dto;

import com.webank.taskman.constant.BizCodeEnum;
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

	public JsonResponse() {
	}

	public JsonResponse(String status, String message, Integer code, String codeMessage, T data) {
		this.status = status;
		this.message = message;
		this.code = code;
		this.codeMessage = codeMessage;
		this.data = data;
	}

	public static JsonResponse okay() {
		return new JsonResponse(STATUS_OK,"Success",null,null,200);
	}

	public static JsonResponse okayWithData(Object data) {
		return okay().withData(data);
	}

	public static JsonResponse customError(String errorMessage) {
		return new JsonResponse(STATUS_ERROR,null,null,null,null);
	}
	public static JsonResponse customError(BizCodeEnum bizCodeEnum, Object data) {
		return new JsonResponse(STATUS_ERROR,null,bizCodeEnum.getCode(),bizCodeEnum.getMessage(),data);
	}
	public static JsonResponse customError(Integer code, String codeMessage, Object data) {
		return new JsonResponse(STATUS_ERROR,null,code,codeMessage,data);
	}
	public static JsonResponse customError(Integer code, String codeMessage,String errorMessage,Object data) {
		return new JsonResponse(STATUS_ERROR,errorMessage,code,codeMessage,data);
	}
	public static JsonResponse customError(BizCodeEnum bizCodeEnum, String errorMessage) {
		return new JsonResponse(STATUS_ERROR,errorMessage,bizCodeEnum.getCode(),bizCodeEnum.getMessage(),null);
	}


}

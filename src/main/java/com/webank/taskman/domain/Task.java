package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import java.util.Date;
import com.baomidou.mybatisplus.annotation.TableId;
import java.io.Serializable;

/**
 * <p>
 * 
 * </p>
 *
 * @author ${author}
 * @since 2020-11-23
 */
public class Task implements Serializable {

    private static final long serialVersionUID = 1L;

    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    private String serviceRequestId;

    private String callbackUrl;

    private String name;

    private String reporter;

    private Date reportTime;

    private String operatorRole;

    private String operator;

    private Date operateTime;

    private String inputParameters;

    private String description;

    private String result;

    private String resultMessage;

    private String status;

    private String requestId;

    private String callbackParameter;

    private String allowedOptions;

    private Date overTime;

    private String dueDate;


    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getServiceRequestId() {
        return serviceRequestId;
    }

    public void setServiceRequestId(String serviceRequestId) {
        this.serviceRequestId = serviceRequestId;
    }

    public String getCallbackUrl() {
        return callbackUrl;
    }

    public void setCallbackUrl(String callbackUrl) {
        this.callbackUrl = callbackUrl;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getReporter() {
        return reporter;
    }

    public void setReporter(String reporter) {
        this.reporter = reporter;
    }

    public Date getReportTime() {
        return reportTime;
    }

    public void setReportTime(Date reportTime) {
        this.reportTime = reportTime;
    }

    public String getOperatorRole() {
        return operatorRole;
    }

    public void setOperatorRole(String operatorRole) {
        this.operatorRole = operatorRole;
    }

    public String getOperator() {
        return operator;
    }

    public void setOperator(String operator) {
        this.operator = operator;
    }

    public Date getOperateTime() {
        return operateTime;
    }

    public void setOperateTime(Date operateTime) {
        this.operateTime = operateTime;
    }

    public String getInputParameters() {
        return inputParameters;
    }

    public void setInputParameters(String inputParameters) {
        this.inputParameters = inputParameters;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getResult() {
        return result;
    }

    public void setResult(String result) {
        this.result = result;
    }

    public String getResultMessage() {
        return resultMessage;
    }

    public void setResultMessage(String resultMessage) {
        this.resultMessage = resultMessage;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public String getRequestId() {
        return requestId;
    }

    public void setRequestId(String requestId) {
        this.requestId = requestId;
    }

    public String getCallbackParameter() {
        return callbackParameter;
    }

    public void setCallbackParameter(String callbackParameter) {
        this.callbackParameter = callbackParameter;
    }

    public String getAllowedOptions() {
        return allowedOptions;
    }

    public void setAllowedOptions(String allowedOptions) {
        this.allowedOptions = allowedOptions;
    }

    public Date getOverTime() {
        return overTime;
    }

    public void setOverTime(Date overTime) {
        this.overTime = overTime;
    }

    public String getDueDate() {
        return dueDate;
    }

    public void setDueDate(String dueDate) {
        this.dueDate = dueDate;
    }

    @Override
    public String toString() {
        return "Task{" +
        "id=" + id +
        ", serviceRequestId=" + serviceRequestId +
        ", callbackUrl=" + callbackUrl +
        ", name=" + name +
        ", reporter=" + reporter +
        ", reportTime=" + reportTime +
        ", operatorRole=" + operatorRole +
        ", operator=" + operator +
        ", operateTime=" + operateTime +
        ", inputParameters=" + inputParameters +
        ", description=" + description +
        ", result=" + result +
        ", resultMessage=" + resultMessage +
        ", status=" + status +
        ", requestId=" + requestId +
        ", callbackParameter=" + callbackParameter +
        ", allowedOptions=" + allowedOptions +
        ", overTime=" + overTime +
        ", dueDate=" + dueDate +
        "}";
    }
}

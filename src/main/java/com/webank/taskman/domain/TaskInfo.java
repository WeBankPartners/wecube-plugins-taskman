package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;

import java.io.Serializable;
import java.util.Date;

/**
 * <p>
 * 任务记录表
 * </p>
 *
 * @author ${author}
 * @since 2020-11-27
 */
public class TaskInfo implements Serializable {

    private static final long serialVersionUID = 1L;

    /**
     * 主键
     */
    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    /**
     * 父级任务id
     */
    private String parentId;

    /**
     * 任务模板id
     */
    private String taskTempId;

    /**
     * 请求id
     */
    private String requestId;

    /**
     * 流程节点
     */
    private String procNode;

    /**
     * 处理角色
     */
    private String dealRole;

    /**
     * 任务名称
     */
    private String name;

    /**
     * 请求编号
     */
    private String requestNo;

    /**
     * 回调url
     */
    private String callbackUrl;

    /**
     * 回调参数
     */
    private String callbackParameter;

    /**
     * 上报人
     */
    private String reporter;

    /**
     * 上报角色
     */
    private String reportRole;

    /**
     * 上报时间
     */
    private Date reportTime;

    /**
     * 执行结果
     */
    private String result;

    /**
     * 紧急程度
     */
    private String emergency;

    /**
     * 描述
     */
    private String description;

    /**
     * 附件id
     */
    private String attachFileId;

    /**
     * 状态
     */
    @TableField("STATUS")
    private Integer status;

    /**
     * 版本号
     */
    @TableField("VERSION")
    private String version;

    /**
     * 创建人
     */
    @TableField("CREATED_BY")
    private String createdBy;

    /**
     * 创建时间
     */
    @TableField("CREATED_TIME")
    private Date createdTime;

    /**
     * 更新人
     */
    @TableField("UPDATED_BY")
    private String updatedBy;

    /**
     * 更新时间
     */
    @TableField("UPDATED_TIME")
    private Date updatedTime;

    /**
     * 是否删除
     */
    @TableField("DEL_FLAG")
    private Integer delFlag;


    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getParentId() {
        return parentId;
    }

    public void setParentId(String parentId) {
        this.parentId = parentId;
    }

    public String getTaskTempId() {
        return taskTempId;
    }

    public void setTaskTempId(String taskTempId) {
        this.taskTempId = taskTempId;
    }

    public String getRequestId() {
        return requestId;
    }

    public void setRequestId(String requestId) {
        this.requestId = requestId;
    }

    public String getProcNode() {
        return procNode;
    }

    public void setProcNode(String procNode) {
        this.procNode = procNode;
    }

    public String getDealRole() {
        return dealRole;
    }

    public void setDealRole(String dealRole) {
        this.dealRole = dealRole;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getRequestNo() {
        return requestNo;
    }

    public void setRequestNo(String requestNo) {
        this.requestNo = requestNo;
    }

    public String getCallbackUrl() {
        return callbackUrl;
    }

    public void setCallbackUrl(String callbackUrl) {
        this.callbackUrl = callbackUrl;
    }

    public String getCallbackParameter() {
        return callbackParameter;
    }

    public void setCallbackParameter(String callbackParameter) {
        this.callbackParameter = callbackParameter;
    }

    public String getReporter() {
        return reporter;
    }

    public void setReporter(String reporter) {
        this.reporter = reporter;
    }

    public String getReportRole() {
        return reportRole;
    }

    public void setReportRole(String reportRole) {
        this.reportRole = reportRole;
    }

    public Date getReportTime() {
        return reportTime;
    }

    public void setReportTime(Date reportTime) {
        this.reportTime = reportTime;
    }

    public String getResult() {
        return result;
    }

    public void setResult(String result) {
        this.result = result;
    }

    public String getEmergency() {
        return emergency;
    }

    public void setEmergency(String emergency) {
        this.emergency = emergency;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getAttachFileId() {
        return attachFileId;
    }

    public void setAttachFileId(String attachFileId) {
        this.attachFileId = attachFileId;
    }

    public Integer getStatus() {
        return status;
    }

    public void setStatus(Integer status) {
        this.status = status;
    }

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }

    public String getCreatedBy() {
        return createdBy;
    }

    public void setCreatedBy(String createdBy) {
        this.createdBy = createdBy;
    }

    public Date getCreatedTime() {
        return createdTime;
    }

    public void setCreatedTime(Date createdTime) {
        this.createdTime = createdTime;
    }

    public String getUpdatedBy() {
        return updatedBy;
    }

    public void setUpdatedBy(String updatedBy) {
        this.updatedBy = updatedBy;
    }

    public Date getUpdatedTime() {
        return updatedTime;
    }

    public void setUpdatedTime(Date updatedTime) {
        this.updatedTime = updatedTime;
    }

    public Integer getDelFlag() {
        return delFlag;
    }

    public void setDelFlag(Integer delFlag) {
        this.delFlag = delFlag;
    }

    @Override
    public String toString() {
        return "TaskInfo{" +
        "id=" + id +
        ", parentId=" + parentId +
        ", taskTempId=" + taskTempId +
        ", requestId=" + requestId +
        ", procNode=" + procNode +
        ", dealRole=" + dealRole +
        ", name=" + name +
        ", requestNo=" + requestNo +
        ", callbackUrl=" + callbackUrl +
        ", callbackParameter=" + callbackParameter +
        ", reporter=" + reporter +
        ", reportRole=" + reportRole +
        ", reportTime=" + reportTime +
        ", result=" + result +
        ", emergency=" + emergency +
        ", description=" + description +
        ", attachFileId=" + attachFileId +
        ", status=" + status +
        ", version=" + version +
        ", createdBy=" + createdBy +
        ", createdTime=" + createdTime +
        ", updatedBy=" + updatedBy +
        ", updatedTime=" + updatedTime +
        ", delFlag=" + delFlag +
        "}";
    }
}

package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import java.util.Date;
import com.baomidou.mybatisplus.annotation.TableId;
import java.io.Serializable;

/**
 * <p>
 * 附件信息表
 * </p>
 *
 * @author ${author}
 * @since 2020-11-26
 */
public class AttachFile implements Serializable {

    private static final long serialVersionUID = 1L;

    /**
     * 主键
     */
    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    /**
     * 附件文件名
     */
    private String attachFileName;

    /**
     * s3服务url
     */
    private String s3Url;

    /**
     * s3_bucket名称
     */
    private String s3BucketName;

    /**
     * s3_key名称
     */
    private String s3KeyName;

    /**
     * 创建人
     */
    private String createdBy;

    /**
     * 创建时间
     */
    private Date createdTime;

    /**
     * 更新人
     */
    private String updatedBy;

    /**
     * 更新时间
     */
    private Date updatedTime;

    /**
     * 是否删除
     */
    private Integer delFlag;


    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getAttachFileName() {
        return attachFileName;
    }

    public void setAttachFileName(String attachFileName) {
        this.attachFileName = attachFileName;
    }

    public String gets3Url() {
        return s3Url;
    }

    public void sets3Url(String s3Url) {
        this.s3Url = s3Url;
    }

    public String gets3BucketName() {
        return s3BucketName;
    }

    public void sets3BucketName(String s3BucketName) {
        this.s3BucketName = s3BucketName;
    }

    public String gets3KeyName() {
        return s3KeyName;
    }

    public void sets3KeyName(String s3KeyName) {
        this.s3KeyName = s3KeyName;
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
        return "AttachFile{" +
        "id=" + id +
        ", attachFileName=" + attachFileName +
        ", s3Url=" + s3Url +
        ", s3BucketName=" + s3BucketName +
        ", s3KeyName=" + s3KeyName +
        ", createdBy=" + createdBy +
        ", createdTime=" + createdTime +
        ", updatedBy=" + updatedBy +
        ", updatedTime=" + updatedTime +
        ", delFlag=" + delFlag +
        "}";
    }
}

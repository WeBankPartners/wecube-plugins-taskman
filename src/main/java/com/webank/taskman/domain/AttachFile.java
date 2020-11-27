package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableLogic;

import java.io.Serializable;
import java.util.Date;

public class AttachFile implements Serializable {

    private static final long serialVersionUID = 1L;

    
    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    
    private String attachFileName;

    
    private String s3Url;

    
    private String s3BucketName;

    
    private String s3KeyName;

    
    private String createdBy;

    
    private Date createdTime;

    
    private String updatedBy;

    
    private Date updatedTime;


    @TableLogic
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

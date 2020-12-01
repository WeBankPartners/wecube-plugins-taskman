package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.webank.taskman.base.BaseEntity;

import java.io.Serializable;

public class AttachFile  extends BaseEntity implements Serializable {

    private static final long serialVersionUID = 1L;

    
    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    
    private String attachFileName;

    
    private String s3Url;

    
    private String s3BucketName;

    
    private String s3KeyName;

    
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


    @Override
    public String toString() {
        return "AttachFile{" +
        "id=" + id +
        ", attachFileName=" + attachFileName +
        ", s3Url=" + s3Url +
        ", s3BucketName=" + s3BucketName +
        ", s3KeyName=" + s3KeyName +
        ", createdBy=" + getCreatedBy() +
        ", createdTime=" + getCreatedTime() +
        ", updatedBy=" + getUpdatedBy() +
        ", updatedTime=" + getUpdatedTime() +
        ", delFlag=" + getDelFlag() +
        "}";
    }
}

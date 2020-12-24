package com.webank.taskman.domain;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.webank.taskman.base.BaseEntity;
import org.springframework.util.StringUtils;

import java.io.Serializable;
import java.util.StringJoiner;

public class AttachFile  extends BaseEntity implements Serializable {

    private static final long serialVersionUID = 1L;

    @TableId(value = "id", type = IdType.ASSIGN_ID)
    private String id;

    private String recordId;

    private String attachFileName;


    private String s3Url;


    private String s3BucketName;


    private String s3KeyName;

    public AttachFile() {
    }

    public AttachFile(String attachFileName, String s3Url, String s3BucketName, String s3KeyName) {
        this.attachFileName = attachFileName;
        this.s3Url = s3Url;
        this.s3BucketName = s3BucketName;
        this.s3KeyName = s3KeyName;
    }

    @JsonIgnore
    public LambdaQueryWrapper getLambdaQueryWrapper() {
        return new LambdaQueryWrapper<AttachFile>()
            .eq(!StringUtils.isEmpty(id), AttachFile::getId, id)
            .eq(!StringUtils.isEmpty(recordId), AttachFile::getRecordId, recordId)
            .eq(!StringUtils.isEmpty(attachFileName), AttachFile::getAttachFileName, attachFileName)
            .eq(!StringUtils.isEmpty(s3Url), AttachFile::getS3Url, s3Url)
            .eq(!StringUtils.isEmpty(s3BucketName), AttachFile::getS3BucketName, s3BucketName)
            .eq(!StringUtils.isEmpty(s3KeyName), AttachFile::getS3KeyName, s3KeyName);
    }


    public String getId() {
        return id;
    }

    public AttachFile setId(String id) {
        this.id = id;
        return this;
    }

    public String getRecordId() {
        return recordId;
    }

    public AttachFile setRecordId(String recordId) {
        this.recordId = recordId;
        return this;
    }

    public String getAttachFileName() {
        return attachFileName;
    }

    public AttachFile setAttachFileName(String attachFileName) {
        this.attachFileName = attachFileName;
        return this;
    }

    public String getS3Url() {
        return s3Url;
    }

    public AttachFile setS3Url(String s3Url) {
        this.s3Url = s3Url;
        return this;
    }

    public String getS3BucketName() {
        return s3BucketName;
    }

    public AttachFile setS3BucketName(String s3BucketName) {
        this.s3BucketName = s3BucketName;
        return this;
    }

    public String getS3KeyName() {
        return s3KeyName;
    }

    public AttachFile setS3KeyName(String s3KeyName) {
        this.s3KeyName = s3KeyName;
        return this;
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", AttachFile.class.getSimpleName() + "[", "]")
                .add("id='" + id + "'")
                .add("recordId='" + recordId + "'")
                .add("attachFileName='" + attachFileName + "'")
                .add("s3Url='" + s3Url + "'")
                .add("s3BucketName='" + s3BucketName + "'")
                .add("s3KeyName='" + s3KeyName + "'")
                .toString();
    }
}

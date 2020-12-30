package com.webank.taskman.base;

import com.baomidou.mybatisplus.annotation.FieldFill;
import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableLogic;
import com.webank.taskman.commons.AuthenticationContextHolder;

import java.util.Date;

public class BaseEntity<T> {

    private String createdBy;

    private Date createdTime;


    private String updatedBy;

    @TableField(fill = FieldFill.UPDATE)
    private Date updatedTime;

    @TableLogic
    private Integer delFlag;

    public String getCreatedBy() {
        return createdBy;
    }

    public <T extends BaseEntity> T  setCreatedBy(String createdBy) {
        this.createdBy = createdBy;
        return (T)this;
    }

    public Date getCreatedTime() {
        return createdTime;
    }

    public <T extends BaseEntity> T setCreatedTime(Date createdTime) {
        this.createdTime = createdTime;
        return (T)this;
    }

    public String getUpdatedBy() {
        return updatedBy;
    }

    public <T extends BaseEntity> T setUpdatedBy(String updatedBy) {
        this.updatedBy = updatedBy;
        return (T)this;
    }

    public Date getUpdatedTime() {
        return updatedTime;

    }

    public <T extends BaseEntity> T setUpdatedTime(Date updatedTime) {
        this.updatedTime = updatedTime;
        return (T)this;
    }

    public Integer getDelFlag() {
        return delFlag;
    }

    public <T extends BaseEntity> T setDelFlag(Integer delFlag) {
        this.delFlag = delFlag;
        return (T)this;
    }

    public BaseEntity setCurrenUserName(BaseEntity entity, String Id) {
//        if(StringUtils.isEmpty(Id)){
//            entity.setCreatedBy(AuthenticationContextHolder.getCurrentUsername());
//        }
        entity.setCreatedBy(AuthenticationContextHolder.getCurrentUsername());
        entity.setUpdatedBy(AuthenticationContextHolder.getCurrentUsername());

        return entity;
    }
}

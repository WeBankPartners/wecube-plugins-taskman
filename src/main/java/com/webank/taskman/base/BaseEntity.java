package com.webank.taskman.base;

import com.baomidou.mybatisplus.annotation.TableLogic;
import com.webank.taskman.commons.AuthenticationContextHolder;

import java.util.Date;

public class BaseEntity<T> {

    private String createdBy;

    private Date createdTime;


    private String updatedBy;


    private Date updatedTime;

    @TableLogic
    private Integer delFlag;

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

    public BaseEntity setCurrenUserName(BaseEntity entity, String Id) {
//        if(StringUtils.isEmpty(Id)){
//            entity.setCreatedBy(AuthenticationContextHolder.getCurrentUsername());
//        }
        entity.setCreatedBy(AuthenticationContextHolder.getCurrentUsername());
        entity.setUpdatedBy(AuthenticationContextHolder.getCurrentUsername());

        return entity;
    }
}

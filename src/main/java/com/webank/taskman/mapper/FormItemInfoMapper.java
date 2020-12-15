package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.webank.taskman.domain.FormItemInfo;

import java.util.List;


public interface FormItemInfoMapper extends BaseMapper<FormItemInfo> {

    List<FormItemInfo> selectFormItemInfo(String id);

}

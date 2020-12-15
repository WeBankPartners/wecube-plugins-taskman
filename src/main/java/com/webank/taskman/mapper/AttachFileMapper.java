package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.webank.taskman.domain.AttachFile;

import java.util.List;
import java.util.Map;


public interface AttachFileMapper extends BaseMapper<AttachFile> {

    List<Map<String,Object>> getList();
}

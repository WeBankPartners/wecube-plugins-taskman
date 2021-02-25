package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.resp.FormItemInfoRespDto;
import org.apache.ibatis.annotations.Param;

import java.util.List;


public interface FormItemTemplateMapper extends BaseMapper<FormItemTemplate> {

    void deleteRequestTemplateByIDMapper(String id);

    int deleteByDomain(FormItemTemplate formItemTemplate);

    List<FormItemInfoRespDto> getCreateTaskServiceMeta(@Param("procInstId")String procInstId, @Param("nodeDefId") String nodeDefId);

    List<FormItemInfoRespDto> selectDetail(String id);
}

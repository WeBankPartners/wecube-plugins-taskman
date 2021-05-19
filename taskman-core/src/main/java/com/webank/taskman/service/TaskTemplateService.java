package com.webank.taskman.service;


import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.base.LocalPageableQueryResult;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.req.TemplateQueryReqDto;
import com.webank.taskman.dto.req.TaskTemplateSaveReqDto;
import com.webank.taskman.dto.resp.TaskTemplateByRoleRespDto;
import com.webank.taskman.dto.resp.TaskTemplateRespDto;

public interface TaskTemplateService extends IService<TaskTemplate> {

    TaskTemplateRespDto saveTaskTemplateByReq(TaskTemplateSaveReqDto taskTemplateReq);

    TaskTemplateRespDto taskTemplateDetail(String id);

    LocalPageableQueryResult<TaskTemplateByRoleRespDto> selectTaskTemplatePage(Integer page, Integer pageSize, TemplateQueryReqDto req);
}

package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.req.TemplateQueryReqDto;
import com.webank.taskman.dto.req.TaskTemplateSaveReqDto;
import com.webank.taskman.dto.resp.TaskTemplateByRoleRespDto;
import com.webank.taskman.dto.resp.TaskTemplateRespDto;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

import java.util.List;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface TaskTemplateConverter extends BaseConverter<TaskTemplateRespDto, TaskTemplate> {

    TaskTemplate toEntityBySaveReq(TaskTemplateSaveReqDto req);

    TaskTemplate toEntityByQueryReq(TemplateQueryReqDto req);

    List<TaskTemplateByRoleRespDto> toRoleRespList(List<TaskTemplate> list);
}
package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;
import com.webank.taskman.dto.resp.TaskTemplateByRoleResp;
import com.webank.taskman.dto.resp.TaskTemplateResp;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

import java.util.List;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface TaskTemplateConverter extends BaseConverter<TaskTemplateResp, TaskTemplate> {

    TaskTemplate toEntityBySaveReq(SaveTaskTemplateReq req);

    List<TaskTemplateByRoleResp> toRoleRespList(List<TaskTemplate> list);
}
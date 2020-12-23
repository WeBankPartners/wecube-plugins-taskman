package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.req.SaveRequestTemplateReq;
import com.webank.taskman.dto.resp.DetailRequestTemplateResq;
import com.webank.taskman.dto.RequestTemplateDTO;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface RequestTemplateConverter extends BaseConverter<RequestTemplateDTO, RequestTemplate> {

    SaveRequestTemplateReq entityToSaveReq(RequestTemplate requestTemplate);

    RequestTemplate saveReqToEntity(SaveRequestTemplateReq req);

    DetailRequestTemplateResq detailRequest(RequestTemplate requestTemplate);



}
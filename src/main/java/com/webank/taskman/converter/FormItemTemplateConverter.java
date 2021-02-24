package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.req.SaveFormItemTemplateReq;
import com.webank.taskman.dto.FormItemTemplateDto;
import com.webank.taskman.dto.resp.FormItemTemplateResp;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

import java.util.List;

@Mapper(componentModel = "spring", uses = {}, unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface FormItemTemplateConverter extends BaseConverter<FormItemTemplateDto, FormItemTemplate> {

    FormItemTemplate toEntityBySaveReq(SaveFormItemTemplateReq req);

    FormItemTemplateResp toRespByEntity(FormItemTemplate formItemTemplate);

    List<FormItemTemplateResp> toRespByEntity( List<FormItemTemplate> formItemTemplate);


}

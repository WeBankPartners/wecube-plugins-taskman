package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.req.SaveFormItemTemplateReq;
import com.webank.taskman.dto.FormItemTemplateDTO;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring", uses = {}, unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface FormItemTemplateConverter extends BaseConverter<FormItemTemplateDTO, FormItemTemplate> {

    FormItemTemplate toEntityBySaveReq(SaveFormItemTemplateReq req);

}

package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.FormItemTemplateDTO;
import com.webank.taskman.dto.req.SaveAndUpdateFormItemTemplateReq;
import com.webank.taskman.dto.resp.FormItemTemplateResq;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

import java.util.List;

@Mapper(componentModel = "spring", uses = {}, unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface FormItemTemplateConverter extends BaseConverter<FormItemTemplateResq, FormItemTemplate> {

    FormItemTemplate addOrUpdateDomain(SaveAndUpdateFormItemTemplateReq req);

}

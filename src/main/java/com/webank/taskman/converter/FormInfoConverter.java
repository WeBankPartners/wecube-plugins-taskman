package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.dto.req.SaveFormInfoReq;
import com.webank.taskman.dto.resp.FormInfoResq;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface FormInfoConverter extends BaseConverter<FormInfoResq, FormInfo> {

    FormInfo saveReqToEntity(SaveFormInfoReq saveFormInfoReq);


}

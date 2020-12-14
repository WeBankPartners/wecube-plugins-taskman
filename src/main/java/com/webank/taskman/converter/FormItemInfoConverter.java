package com.webank.taskman.converter;
import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.dto.req.SaveFormItemInfoReq;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface FormItemInfoConverter extends BaseConverter<SaveFormItemInfoReq, FormItemInfo> {


}

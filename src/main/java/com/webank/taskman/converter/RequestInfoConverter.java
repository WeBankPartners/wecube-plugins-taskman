package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.RequestInfo;
import com.webank.taskman.dto.CreateTaskDto;
import com.webank.taskman.dto.RequestInfoDto;
import com.webank.taskman.dto.req.SaveRequestInfoReq;
import com.webank.taskman.dto.resp.RequestInfoInstanceResq;
import com.webank.taskman.dto.resp.RequestInfoResq;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface RequestInfoConverter extends BaseConverter<RequestInfoDto, RequestInfo> {

    RequestInfo reqToDomain(SaveRequestInfoReq req);

    RequestInfo createDtoToDomain(CreateTaskDto dto);

    RequestInfoResq toResp(RequestInfo requestInfo);

    RequestInfoInstanceResq toInstanceResp(RequestInfo requestInfo);
}

package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.RequestInfo;
import com.webank.taskman.dto.CreateTaskDto;
import com.webank.taskman.dto.RequestInfoDto;
import com.webank.taskman.dto.req.RequestInfoSaveReqDto;
import com.webank.taskman.dto.resp.RequestInfoInstanceResqDto;
import com.webank.taskman.dto.resp.RequestInfoResqDto;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface RequestInfoConverter extends BaseConverter<RequestInfoDto, RequestInfo> {

    RequestInfo reqToDomain(RequestInfoSaveReqDto req);

    RequestInfo createDtoToDomain(CreateTaskDto dto);

    RequestInfoResqDto toResp(RequestInfo requestInfo);

    RequestInfoInstanceResqDto toInstanceResp(RequestInfo requestInfo);
}

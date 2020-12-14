package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.dto.resp.SynthesisRequestInfoFormRequest;
import com.webank.taskman.dto.resp.SynthesisTaskInfoFormTask;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;


@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface SynthesisTaskInfoFormTaskConverter extends BaseConverter<SynthesisTaskInfoFormTask, FormInfo> {

}
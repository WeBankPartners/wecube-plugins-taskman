package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.dto.req.ProcessingTasksReq;
import com.webank.taskman.dto.req.SaveFormInfoReq;
import com.webank.taskman.dto.resp.FormInfoResq;
import com.webank.taskman.dto.resp.ProcessingTasksResp;
import com.webank.taskman.dto.resp.RequestFormResq;
import com.webank.taskman.dto.resp.TaskFormResq;
import org.mapstruct.Mapper;
import org.mapstruct.ReportingPolicy;

@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface FormInfoConverter extends BaseConverter<FormInfoResq, FormInfo> {

    TaskFormResq toTaskFormResq(FormInfo formInfo);

    RequestFormResq toRequestFormResq(FormInfo formInfo);
}

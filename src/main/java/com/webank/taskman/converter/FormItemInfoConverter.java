package com.webank.taskman.converter;
import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.dto.CoreCreateTaskDTO.TaskInfoReq.FormItemBean;
import com.webank.taskman.dto.req.SaveFormItemInfoReq;
import com.webank.taskman.dto.resp.FormItemInfoResp;
import com.webank.taskman.dto.resp.TaskServiceMetaResp.TaskServiceMetaFormItem;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Mappings;
import org.mapstruct.ReportingPolicy;

import java.util.List;
@Mapper(componentModel = "spring",uses = {},unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface FormItemInfoConverter extends BaseConverter<FormItemInfoResp, FormItemInfo> {

    FormItemInfo toEntityByReq(SaveFormItemInfoReq req);

    List<FormItemInfo> toEntityByReqs(List<SaveFormItemInfoReq> req);


    @Mappings({
            @Mapping(target = "itemId",source = "id"),
            @Mapping(target = "key",source = "name"),
            @Mapping(target = "valueDef.type",source = "elementType"),
            @Mapping(target = "valueDef.expr",source = "value"),
    })
    TaskServiceMetaFormItem respToServiceMeta(FormItemInfoResp resp);

    List<TaskServiceMetaFormItem> respToServiceMetas(List<FormItemInfoResp> resp);

    @Mappings({
            @Mapping(target = "itemTempId",source = "itemId"),
            @Mapping(target = "name",source = "key"),
            @Mapping(target = "value",expression=
                    "java(bean.getVal().stream().collect(java.util.stream.Collectors.joining(\",\")))"),
    })
    FormItemInfo toEntityByBean(FormItemBean bean);

    List<FormItemInfo> toEntityByBeans(List<FormItemBean> bean);
}

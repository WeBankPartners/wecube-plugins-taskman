package com.webank.taskman.converter;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.dto.CoreCreateTaskDto.FormItemBean;
import com.webank.taskman.dto.CreateTaskDto.EntityAttrValueDto;
import com.webank.taskman.dto.req.FormItemInfoRequestDto;
import com.webank.taskman.dto.resp.FormItemInfoRespDto;
import com.webank.taskman.dto.resp.TaskServiceMetaRespDto.TaskServiceMetaFormItem;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.ReportingPolicy;
import java.util.List;

@Mapper(componentModel = "spring", uses = {}, unmappedTargetPolicy = ReportingPolicy.IGNORE)
public interface FormItemInfoConverter extends BaseConverter<FormItemInfoRespDto, FormItemInfo> {

    FormItemInfo toEntityByReq(FormItemInfoRequestDto req);

    List<FormItemInfo> toEntityByReq(List<FormItemInfoRequestDto> req);

    @Mapping(target = "itemId", source = "id")
    @Mapping(target = "key", source = "name")
    @Mapping(target = "valueDef.type", source = "elementType")
    @Mapping(target = "valueDef.expr", source = "value")
    TaskServiceMetaFormItem respToServiceMeta(FormItemInfoRespDto resp);

    List<TaskServiceMetaFormItem> respToServiceMeta(List<FormItemInfoRespDto> resp);

    @Mapping(target = "itemTempId", source = "itemId")
    @Mapping(target = "name", source = "key")
    @Mapping(target = "value", expression = "java(bean.getVal().stream().collect(java.util.stream.Collectors.joining(\",\")))")
    FormItemInfo toEntityByBean(FormItemBean bean);

    List<FormItemInfo> toEntityByBean(List<FormItemBean> bean);

    @Mapping(target = "value", expression = "java(attrValueDto.getDataValue()+\"\")")
    FormItemInfo toEntityByAttrValue(EntityAttrValueDto attrValueDto);

    List<FormItemInfo> toEntityByAttrValue(List<EntityAttrValueDto> attrValueDtos);

}

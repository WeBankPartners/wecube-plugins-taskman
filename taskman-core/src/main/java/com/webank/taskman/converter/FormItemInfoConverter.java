package com.webank.taskman.converter;

import java.util.ArrayList;
import java.util.List;

import org.springframework.beans.BeanUtils;
import org.springframework.stereotype.Service;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.dto.EntityAttrValueDto;
import com.webank.taskman.dto.req.FormItemInfoRequestDto;
import com.webank.taskman.dto.resp.FormItemInfoQueryResultDto;

@Service
public class FormItemInfoConverter implements BaseConverter<FormItemInfoQueryResultDto, FormItemInfo> {

    public FormItemInfo convertToFormItemInfoByReq(FormItemInfoRequestDto dto){
        FormItemInfo info = new FormItemInfo();
        BeanUtils.copyProperties(dto, info);
        
        return info;
    }

    public List<FormItemInfo> convertToFormItemInfosByReq(List<FormItemInfoRequestDto> dtos){
        if(dtos == null){
            return null;
        }
        
        List<FormItemInfo> infos = new ArrayList<>();
        for(FormItemInfoRequestDto dto : dtos ){
            FormItemInfo info = convertToFormItemInfoByReq(dto);
            infos.add(info);
        }
        
        return infos;
    }

//    @Mapping(target = "itemId", source = "id")
//    @Mapping(target = "key", source = "name")
//    @Mapping(target = "valueDef.type", source = "elementType")
//    @Mapping(target = "valueDef.expr", source = "value")
//    TaskServiceMetaFormItem respToServiceMeta(FormItemInfoRespDto resp);

//    List<TaskServiceMetaFormItem> respToServiceMeta(List<FormItemInfoRespDto> resp);

//    @Mapping(target = "itemTempId", source = "itemId")
//    @Mapping(target = "name", source = "key")
//    @Mapping(target = "value", expression = "java(bean.getVal().stream().collect(java.util.stream.Collectors.joining(\",\")))")
//    FormItemInfo toEntityByBean(FormItemBean bean);
//
//    List<FormItemInfo> toEntityByBean(List<FormItemBean> bean);

    public FormItemInfo convertToFormItemInfo(EntityAttrValueDto dto) {
        FormItemInfo info = new FormItemInfo();
        // info.setFormId();
        info.setItemTempId(dto.getItemTempId());
        // info.setId()
        info.setName(dto.getName());
        info.setValue(String.valueOf(dto.getDataValue()));
        // info.setRecordId(recordId)

        return info;
    }

    public List<FormItemInfo> convertToFormItemInfos(List<EntityAttrValueDto> attrValueDtos) {
        if (attrValueDtos == null) {
            return null;
        }

        List<FormItemInfo> itemInfos = new ArrayList<>();
        for (EntityAttrValueDto dto : attrValueDtos) {
            FormItemInfo info = new FormItemInfo();
            // info.setFormId();
            info.setItemTempId(dto.getItemTempId());
            // info.setId()
            info.setName(dto.getName());
            info.setValue(String.valueOf(dto.getDataValue()));
            // info.setRecordId(recordId)

            itemInfos.add(info);
        }

        return itemInfos;
    }

    @Override
    public FormItemInfo convertToEntity(FormItemInfoQueryResultDto dto) {
        if (dto == null) {
            return null;
        }

        FormItemInfo entity = new FormItemInfo();
        BeanUtils.copyProperties(dto, entity);
        return entity;
    }

    @Override
    public FormItemInfoQueryResultDto convertToDto(FormItemInfo entity) {
        if (entity == null) {
            return null;
        }

        FormItemInfoQueryResultDto dto = new FormItemInfoQueryResultDto();
        BeanUtils.copyProperties(entity, dto);
        return dto;
    }

    @Override
    public List<FormItemInfo> convertToEntities(List<FormItemInfoQueryResultDto> dtos) {
        if (dtos == null) {
            return null;
        }

        List<FormItemInfo> entities = new ArrayList<>();
        for (FormItemInfoQueryResultDto dto : dtos) {
            FormItemInfo entity = convertToEntity(dto);
            entities.add(entity);
        }
        return entities;
    }

    @Override
    public List<FormItemInfoQueryResultDto> convertToDtos(List<FormItemInfo> entities) {
        if (entities == null) {
            return null;
        }

        List<FormItemInfoQueryResultDto> dtos = new ArrayList<>();
        for (FormItemInfo entity : entities) {
            FormItemInfoQueryResultDto dto = convertToDto(entity);
            dtos.add(dto);
        }
        return dtos;
    }

}

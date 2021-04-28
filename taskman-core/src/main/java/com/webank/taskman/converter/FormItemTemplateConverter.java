package com.webank.taskman.converter;

import java.util.ArrayList;
import java.util.List;

import org.springframework.beans.BeanUtils;
import org.springframework.stereotype.Service;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.dto.FormItemTemplateDto;
import com.webank.taskman.dto.req.FormItemTemplateSaveReqDto;
import com.webank.taskman.dto.resp.FormItemTemplateQueryResultDto;

@Service
public class FormItemTemplateConverter implements BaseConverter<FormItemTemplateDto, FormItemTemplate> {

    @Override
    public FormItemTemplate convertToEntity(FormItemTemplateDto dto) {
        if (dto == null) {
            return null;
        }

        FormItemTemplate entity = new FormItemTemplate();
        BeanUtils.copyProperties(dto, entity);
        return entity;
    }

    @Override
    public FormItemTemplateDto convertToDto(FormItemTemplate entity) {
        if (entity == null) {
            return null;
        }

        FormItemTemplateDto dto = new FormItemTemplateDto();
        BeanUtils.copyProperties(entity, dto);
        return dto;
    }

    @Override
    public List<FormItemTemplate> convertToEntities(List<FormItemTemplateDto> dtos) {
        if (dtos == null) {
            return null;
        }

        List<FormItemTemplate> entities = new ArrayList<>();
        for (FormItemTemplateDto dto : dtos) {
            FormItemTemplate entity = convertToEntity(dto);
            entities.add(entity);
        }
        return entities;
    }

    @Override
    public List<FormItemTemplateDto> convertToDtos(List<FormItemTemplate> entities) {
        if (entities == null) {
            return null;
        }

        List<FormItemTemplateDto> dtos = new ArrayList<>();
        for (FormItemTemplate entity : entities) {
            FormItemTemplateDto dto = convertToDto(entity);
            dtos.add(dto);
        }
        return dtos;
    }

    public FormItemTemplate convertToFormItemTemplate(FormItemTemplateSaveReqDto dto) {
        if (dto == null) {
            return null;
        }
        FormItemTemplate entity = new FormItemTemplate();
        BeanUtils.copyProperties(dto, entity);

        return entity;
    }

    public FormItemTemplateDto toRespByEntity(FormItemTemplate entity) {
        if (entity == null) {
            return null;
        }

        FormItemTemplateDto dto = new FormItemTemplateDto();
        BeanUtils.copyProperties(entity, dto);
        return dto;
    }

    public List<FormItemTemplateQueryResultDto> convertToFormItemTemplateQueryResultDtos(
            List<FormItemTemplate> entities) {
        if (entities == null) {
            return null;
        }

        List<FormItemTemplateQueryResultDto> dtos = new ArrayList<>();
        for (FormItemTemplate entity : entities) {
            FormItemTemplateQueryResultDto dto = new FormItemTemplateQueryResultDto();
            BeanUtils.copyProperties(entity, dto);

            dtos.add(dto);
        }

        return dtos;
    }

}

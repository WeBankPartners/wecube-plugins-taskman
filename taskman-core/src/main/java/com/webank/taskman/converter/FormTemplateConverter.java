package com.webank.taskman.converter;

import java.util.List;

import org.springframework.stereotype.Service;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.dto.req.FormTemplateSaveReqDto;
import com.webank.taskman.dto.resp.FormTemplateRespDto;

@Service
public class FormTemplateConverter implements BaseConverter<FormTemplateRespDto, FormTemplate> {

    public FormTemplate reqToDomain(FormTemplateSaveReqDto req){
        return null;
    }

    @Override
    public FormTemplate convertToEntity(FormTemplateRespDto dto) {
        // TODO Auto-generated method stub
        return null;
    }

    @Override
    public FormTemplateRespDto convertToDto(FormTemplate entity) {
        // TODO Auto-generated method stub
        return null;
    }

    @Override
    public List<FormTemplate> convertToEntities(List<FormTemplateRespDto> dtos) {
        // TODO Auto-generated method stub
        return null;
    }

    @Override
    public List<FormTemplateRespDto> convertToDtos(List<FormTemplate> entities) {
        // TODO Auto-generated method stub
        return null;
    }

}
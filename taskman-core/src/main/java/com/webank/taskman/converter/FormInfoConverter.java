package com.webank.taskman.converter;

import java.util.List;

import org.springframework.stereotype.Service;

import com.webank.taskman.base.BaseConverter;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.dto.resp.FormInfoResqDto;

@Service
public class FormInfoConverter implements BaseConverter<FormInfoResqDto, FormInfo> {

    @Override
    public FormInfo convertToEntity(FormInfoResqDto dto) {
        // TODO Auto-generated method stub
        return null;
    }

    @Override
    public FormInfoResqDto convertToDto(FormInfo entity) {
        // TODO Auto-generated method stub
        return null;
    }

    @Override
    public List<FormInfo> convertToEntities(List<FormInfoResqDto> dtos) {
        // TODO Auto-generated method stub
        return null;
    }

    @Override
    public List<FormInfoResqDto> convertToDtos(List<FormInfo> entities) {
        // TODO Auto-generated method stub
        return null;
    }

}

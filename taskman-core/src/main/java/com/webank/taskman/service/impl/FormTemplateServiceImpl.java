package com.webank.taskman.service.impl;

import java.util.Date;
import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.StringUtils;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.core.conditions.update.UpdateWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.base.LocalPageableQueryResult;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.RecordDeleteFlag;
import com.webank.taskman.constant.RoleType;
import com.webank.taskman.constant.TemplateType;
import com.webank.taskman.converter.FormItemTemplateConverter;
import com.webank.taskman.converter.FormTemplateConverter;
import com.webank.taskman.converter.TaskTemplateConverter;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.RoleDto;
import com.webank.taskman.dto.req.FormTemplateSaveReqDto;
import com.webank.taskman.dto.resp.FormItemTemplateQueryResultDto;
import com.webank.taskman.dto.resp.FormTemplateQueryResultDto;
import com.webank.taskman.dto.resp.TaskTemplateRespDto;
import com.webank.taskman.mapper.FormTemplateMapper;
import com.webank.taskman.service.FormItemTemplateService;
import com.webank.taskman.service.FormTemplateService;
import com.webank.taskman.service.RoleRelationService;
import com.webank.taskman.service.TaskTemplateService;

@Service
public class FormTemplateServiceImpl extends ServiceImpl<FormTemplateMapper, FormTemplate>
        implements FormTemplateService {

    @Autowired
    private FormTemplateMapper formTemplateMapper;
    @Autowired
    private FormTemplateConverter formTemplateConverter;

    @Autowired
    private FormItemTemplateService formItemTemplateService;

    @Autowired
    private FormItemTemplateConverter formItemTemplateConverter;

    @Autowired
    private TaskTemplateService taskTemplateService;

    @Autowired
    private TaskTemplateConverter taskTemplateConverter;
    
    @Autowired
    private RoleRelationService roleRelationService;
    
    @Override
    public LocalPageableQueryResult<FormTemplateQueryResultDto> selectFormTemplate(Integer page, Integer pageSize,
            FormTemplateSaveReqDto req) {
        Page<FormTemplate> pageInfo = new Page<>(page, pageSize);
        LambdaQueryWrapper<FormTemplate> formTemplateQueryWrapper = formTemplateConverter.reqToDomain(req)
                .getLambdaQueryWrapper();
        IPage<FormTemplate> iPage = formTemplateMapper.selectPage(pageInfo, formTemplateQueryWrapper);
        List<FormTemplateQueryResultDto> formTemplateResps = formTemplateConverter.convertToDtos(iPage.getRecords());

        LocalPageableQueryResult<FormTemplateQueryResultDto> queryResponse = new LocalPageableQueryResult<>(iPage.getSize(), iPage.getCurrent(),
                iPage.getSize(), formTemplateResps);
        return queryResponse;
    }

    @Override
    public void deleteFormTemplate(String id) {
        if (StringUtils.isEmpty(id)) {
            throw new TaskmanRuntimeException("Form template parameter cannot be empty.");
        }
        UpdateWrapper<FormTemplate> wrapper = new UpdateWrapper<>();
        wrapper.lambda().eq(FormTemplate::getId, id).set(FormTemplate::getDelFlag, RecordDeleteFlag.Deleted.ordinal())
                .set(FormTemplate::getUpdatedTime, new Date());
        ;
        formTemplateMapper.update(null, wrapper);
    }

    @Override
    public FormTemplateQueryResultDto detailFormTemplate(FormTemplateSaveReqDto req) {

        LambdaQueryWrapper<FormTemplate> formTemplateQueryWrapper = formTemplateConverter.reqToDomain(req)
                .getLambdaQueryWrapper();
        FormTemplate formTemplateEntity = formTemplateMapper.selectOne(formTemplateQueryWrapper);
        FormTemplateQueryResultDto formTemplateRespDto = formTemplateConverter.convertToDto(formTemplateEntity);

        if (formTemplateRespDto == null) {
            return null;
        }
        FormItemTemplate formItemTemplateCritetia = new FormItemTemplate()
                .setFormTemplateId(formTemplateRespDto.getId());
        LambdaQueryWrapper<FormItemTemplate> formItemTemplateQueryWrapper = formItemTemplateCritetia
                .getLambdaQueryWrapper();

        List<FormItemTemplate> formItemTemplateEntities = formItemTemplateService.list(formItemTemplateQueryWrapper);
        List<FormItemTemplateQueryResultDto> formItemTemplateDtos = formItemTemplateConverter
                .convertToFormItemTemplateQueryResultDtos(formItemTemplateEntities);
        formTemplateRespDto.setItems(formItemTemplateDtos);
        
        
        RoleRelation useRoleCriteria = new RoleRelation();
        useRoleCriteria.setRecordId(formTemplateEntity.getTempId());
        useRoleCriteria.setRoleType(RoleType.USE_ROLE.getType());
        
        LambdaQueryWrapper<RoleRelation> useRoleQueryWrapper = useRoleCriteria.getLambdaQueryWrapper();
        List<RoleRelation> useRoleRelations = roleRelationService.list(useRoleQueryWrapper);
        
        if(useRoleRelations != null){
            for(RoleRelation useRole : useRoleRelations){
                RoleDto roleDto = new RoleDto();
                roleDto.setRoleName(useRole.getRoleName());
                roleDto.setDisplayName(useRole.getDisplayName());
                
                formTemplateRespDto.getUseRoles().add(roleDto);
            }
        }
        
        
        RoleRelation mgmtRoleCriteria = new RoleRelation();
        mgmtRoleCriteria.setRecordId(formTemplateEntity.getTempId());
        mgmtRoleCriteria.setRoleType(RoleType.MANAGE_ROLE.getType());
        
        LambdaQueryWrapper<RoleRelation> mgmtRoleQueryWrapper = mgmtRoleCriteria.getLambdaQueryWrapper();
        List<RoleRelation> mgmtRoleRelations = roleRelationService.list(mgmtRoleQueryWrapper);
        
        if(mgmtRoleRelations != null){
            for(RoleRelation mgmtRole : mgmtRoleRelations){
                RoleDto roleDto = new RoleDto();
                roleDto.setRoleName(mgmtRole.getRoleName());
                roleDto.setDisplayName(mgmtRole.getDisplayName());
                
                formTemplateRespDto.getManageRoles().add(roleDto);
            }
        }

        if (TemplateType.REQUEST.getType().equals(req.getTempType())) {

            LambdaQueryWrapper<TaskTemplate> taskTemplateQueryWrapper = new TaskTemplate()
                    .setRequestTemplateId(req.getTempId()).getLambdaQueryWrapper();
            List<TaskTemplate> taskTemplateEntities = taskTemplateService.list(taskTemplateQueryWrapper);
            List<TaskTemplateRespDto> taskTemplateDtos = taskTemplateConverter.convertToDtos(taskTemplateEntities);
            if (taskTemplateDtos != null) {
                formTemplateRespDto.setTaskTemplates(taskTemplateDtos);
            }

        }
        return formTemplateRespDto;
    }

    @Override
    @Transactional
    public FormTemplateQueryResultDto saveOrUpdateFormTemplate(FormTemplateSaveReqDto formTemplateReq) {
        FormTemplate formTemplate = formTemplateConverter.reqToDomain(formTemplateReq);

        formTemplate.setName(StringUtils.isEmpty(formTemplate.getName()) ? "" : formTemplate.getName());
        formTemplate.setCreatedBy(AuthenticationContextHolder.getCurrentUsername());
        formTemplate.setUpdatedBy(AuthenticationContextHolder.getCurrentUsername());

        saveOrUpdate(formTemplate);

        formItemTemplateService.deleteByDomain(new FormItemTemplate().setFormTemplateId(formTemplate.getId()));

        formTemplateReq.getFormItems().stream().forEach(item -> {
            FormItemTemplate formItemTemplate = formItemTemplateConverter.convertToFormItemTemplate(item);
            formItemTemplate.setFormTemplateId(formTemplate.getId());
            formItemTemplate.setTempId(formTemplate.getTempId());
            formItemTemplateService.save(formItemTemplate);
        });
        return new FormTemplateQueryResultDto().setId(formTemplate.getId());
    }

    @Override
    public FormTemplateQueryResultDto queryDetailByTemplate(Integer tempType, String tempId) {
        
        FormTemplate formTemplateEntity = getOne(new FormTemplate().setTempId(tempId).setTempType(tempType + "").getLambdaQueryWrapper());
        return formTemplateConverter
                .convertToDto(formTemplateEntity);
    }
}

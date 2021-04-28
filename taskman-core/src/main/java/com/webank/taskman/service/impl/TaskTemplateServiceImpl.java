package com.webank.taskman.service.impl;

import java.util.List;

import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.base.PageableQueryResult;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.RoleType;
import com.webank.taskman.constant.TemplateType;
import com.webank.taskman.converter.TaskTemplateConverter;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.RoleDto;
import com.webank.taskman.dto.req.TemplateQueryReqDto;
import com.webank.taskman.dto.req.FormTemplateSaveReqDto;
import com.webank.taskman.dto.req.TaskTemplateSaveReqDto;
import com.webank.taskman.dto.resp.TaskTemplateByRoleRespDto;
import com.webank.taskman.dto.resp.TaskTemplateRespDto;
import com.webank.taskman.mapper.TaskTemplateMapper;
import com.webank.taskman.service.FormTemplateService;
import com.webank.taskman.service.RequestTemplateService;
import com.webank.taskman.service.RoleRelationService;
import com.webank.taskman.service.TaskTemplateService;

@Service
public class TaskTemplateServiceImpl extends ServiceImpl<TaskTemplateMapper, TaskTemplate>
        implements TaskTemplateService {

    @Autowired
    private RoleRelationService roleRelationService;

    @Autowired
    private FormTemplateService formTemplateService;

    @Autowired
    private TaskTemplateConverter taskTemplateConverter;

    @Autowired
    private TaskTemplateMapper taskTemplateMapper;

    @Autowired
    private RequestTemplateService requestTemplateService;

    @Override
    @Transactional
    public TaskTemplateRespDto saveTaskTemplateByReq(TaskTemplateSaveReqDto req) {
        RequestTemplate requestTemplate = requestTemplateService.getById(req.getTempId());
        if (requestTemplate == null) {
            throw new TaskmanRuntimeException("The Request Template does not exist! ID:" + req.getTempId());
        }

        TaskTemplate taskTemplate = taskTemplateConverter.convertToTaskTemplate(req);
        taskTemplate.setRequestTemplateId(requestTemplate.getId());
        taskTemplate.setUpdatedBy(AuthenticationContextHolder.getCurrentUsername());
        taskTemplate.setCreatedBy(AuthenticationContextHolder.getCurrentUsername());
        saveOrUpdate(taskTemplate);

        String taskTemplateId = taskTemplate.getId();
        roleRelationService.saveRoleRelationByTemplate(taskTemplateId, req.getUseRoles(), req.getManageRoles());

        // Bad smells
        FormTemplateSaveReqDto formTemplateReqDto = req.getForm();
        formTemplateReqDto.setTempId(taskTemplateId);
        formTemplateReqDto.setTempType(TemplateType.TASK.getType());
        formTemplateService.saveOrUpdateFormTemplate(formTemplateReqDto);

        TaskTemplateRespDto taskTemplateRespDto = new TaskTemplateRespDto();
        taskTemplateRespDto.setId(taskTemplateId);
        return taskTemplateRespDto;
    }

    @Override
    public TaskTemplateRespDto taskTemplateDetail(String id) {
        TaskTemplate taskTemplate = taskTemplateMapper.selectById(id);
        TaskTemplateRespDto taskTemplateResp = taskTemplateConverter.convertToDto(taskTemplate);
        List<RoleRelation> relations = roleRelationService
                .list(new RoleRelation().setRecordId(id).getLambdaQueryWrapper());
        relations.stream().forEach(roleRelation -> {
            RoleDto roleDTO = new RoleDto();
            roleDTO.setRoleName(roleRelation.getRoleName());
            roleDTO.setDisplayName(roleRelation.getRoleName());
            if (RoleType.USE_ROLE.getType() == roleRelation.getRoleType()) {
                taskTemplateResp.getUseRoles().add(roleDTO);
            } else {
                taskTemplateResp.getManageRoles().add(roleDTO);
            }
        });
        return taskTemplateResp;
    }

    @Override
    public PageableQueryResult<TaskTemplateByRoleRespDto> selectTaskTemplatePage(Integer page, Integer pageSize,
            TemplateQueryReqDto req) {
        String inSql = TemplateQueryReqDto.getEqUseRole();

        LambdaQueryWrapper<TaskTemplate> queryWrapper = taskTemplateConverter.convertToTaskTemplate(req)
                .getLambdaQueryWrapper().inSql(!StringUtils.isEmpty(inSql), TaskTemplate::getId, inSql);

        IPage<TaskTemplate> iPage = taskTemplateMapper.selectPage(new Page<>(page, pageSize), queryWrapper);
        List<TaskTemplateByRoleRespDto> list = taskTemplateConverter
                .convertToTaskTemplateByRoleRespDtos(iPage.getRecords());

        PageableQueryResult<TaskTemplateByRoleRespDto> queryResponse = new PageableQueryResult<>(iPage.getTotal(),
                iPage.getCurrent(), iPage.getSize(), list);

        return queryResponse;
    }

}

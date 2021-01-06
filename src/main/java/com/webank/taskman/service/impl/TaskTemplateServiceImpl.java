package com.webank.taskman.service.impl;


import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.RoleTypeEnum;
import com.webank.taskman.constant.TemplateTypeEnum;
import com.webank.taskman.converter.TaskTemplateConverter;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.RoleDTO;
import com.webank.taskman.dto.req.QueryRoleRelationBaseReq;
import com.webank.taskman.dto.req.QueryTemplateReq;
import com.webank.taskman.dto.req.SaveFormTemplateReq;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;
import com.webank.taskman.dto.resp.TaskTemplateByRoleResp;
import com.webank.taskman.dto.resp.TaskTemplateResp;
import com.webank.taskman.mapper.TaskTemplateMapper;
import com.webank.taskman.service.FormTemplateService;
import com.webank.taskman.service.RequestTemplateService;
import com.webank.taskman.service.RoleRelationService;
import com.webank.taskman.service.TaskTemplateService;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Service
public class TaskTemplateServiceImpl extends ServiceImpl<TaskTemplateMapper, TaskTemplate> implements TaskTemplateService {


    @Autowired
    RoleRelationService roleRelationService;

    @Autowired
    FormTemplateService formTemplateService;

    @Autowired
    TaskTemplateConverter taskTemplateConverter;

    @Autowired
    TaskTemplateMapper taskTemplateMapper;

    @Autowired
    RequestTemplateService requestTemplateService;


    @Override
    @Transactional
    public TaskTemplateResp saveTaskTemplateByReq(SaveTaskTemplateReq req) throws TaskmanRuntimeException {
        RequestTemplate requestTemplate = requestTemplateService.getById(req.getTempId());
        if(null == requestTemplate){
            throw  new TaskmanRuntimeException("The RequestTemplate is not exists! id:"+req.getTempId());
        }
        TaskTemplate taskTemplate = taskTemplateConverter.toEntityBySaveReq(req);
        taskTemplate.setRequestTemplateId(req.getTempId());
        taskTemplate.setCurrenUserName(taskTemplate, taskTemplate.getId());
        saveOrUpdate(taskTemplate);
        String taskTemplateId = taskTemplate.getId();
        roleRelationService.saveRoleRelationByTemplate(taskTemplateId, req.getUseRoles(), req.getManageRoles());

        SaveFormTemplateReq formTemplateReq = req.getForm();
        formTemplateReq.setTempId(taskTemplateId);
        formTemplateReq.setTempType(TemplateTypeEnum.TASK.getType());
        formTemplateService.saveFormTemplateByReq(req.getForm());

        TaskTemplateResp taskTemplateResp = new TaskTemplateResp();
        taskTemplateResp.setId(taskTemplateId);
        return taskTemplateResp;
    }

    @Override
    public TaskTemplateResp taskTemplateDetail(String id) {
        TaskTemplate taskTemplate = taskTemplateMapper.selectById(id);
        TaskTemplateResp taskTemplateResp = taskTemplateConverter.toDto(taskTemplate);
        List<RoleRelation> relations = roleRelationService.list(
                new RoleRelation().setRecordId(id).getLambdaQueryWrapper()
        );
        relations.stream().forEach(roleRelation -> {
            RoleDTO roleDTO = new RoleDTO();
            roleDTO.setRoleName(roleRelation.getRoleName());
            roleDTO.setDisplayName(roleRelation.getRoleName());
            if (RoleTypeEnum.USE_ROLE.getType() == roleRelation.getRoleType()) {
                taskTemplateResp.getUseRoles().add(roleDTO);
            } else {
                taskTemplateResp.getManageRoles().add(roleDTO);
            }
        });
        return taskTemplateResp;
    }

    @Override
    public QueryResponse<TaskTemplateByRoleResp> selectTaskTemplatePage(Integer page, Integer pageSize, QueryTemplateReq req) {
        String inSql = req.getEqUseRole();
        LambdaQueryWrapper<TaskTemplate> queryWrapper = taskTemplateConverter.toEntityByQueryReq(req)
                .getLambdaQueryWrapper().inSql(!StringUtils.isEmpty(inSql),TaskTemplate::getId,inSql);
        IPage<TaskTemplate> iPage = taskTemplateMapper.selectPage(new Page<>(page, pageSize),queryWrapper);
        List<TaskTemplateByRoleResp> list = taskTemplateConverter.toRoleRespList(iPage.getRecords());
        QueryResponse<TaskTemplateByRoleResp> queryResponse = new QueryResponse<>(iPage.getTotal(), iPage.getCurrent(), iPage.getSize(),list);
        return queryResponse;
    }


}

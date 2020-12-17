package com.webank.taskman.service.impl;


import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.conditions.update.UpdateWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.constant.RoleTypeEnum;
import com.webank.taskman.constant.TemplateTypeEnum;
import com.webank.taskman.converter.SynthesisTaskTemplateConverter;
import com.webank.taskman.converter.TaskTemplateConverter;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.PageInfo;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.RoleDTO;
import com.webank.taskman.dto.req.QueryRoleRelationBaseReq;
import com.webank.taskman.dto.req.SaveFormTemplateReq;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;
import com.webank.taskman.dto.resp.TaskTemplateByRoleResp;
import com.webank.taskman.dto.resp.TaskTemplateResp;
import com.webank.taskman.mapper.TaskTemplateMapper;
import com.webank.taskman.service.FormTemplateService;
import com.webank.taskman.service.RoleRelationService;
import com.webank.taskman.service.TaskTemplateService;
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
    SynthesisTaskTemplateConverter synthesisTaskTemplateConverter;

    @Override
    @Transactional
    public TaskTemplateResp saveTaskTemplateByReq(SaveTaskTemplateReq req) {
        TaskTemplate taskTemplate = taskTemplateConverter.toEntityBySaveReq(req);
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
    public void deleteTaskTemplateByIDService(String id) {
        UpdateWrapper<TaskTemplate> wrapper = new UpdateWrapper<>();
        wrapper.eq("id", id).set("del_flag", 1);
        taskTemplateMapper.update(null, wrapper);
    }

    @Override
    public List<TaskTemplateResp> selectTaskTemplateAll() {
        List<TaskTemplate> taskTemplates = taskTemplateMapper.selectList(null);
        List<TaskTemplateResp> svResps = taskTemplateConverter.toDto(taskTemplates);
        for (TaskTemplateResp svResp : svResps) {
            List<RoleRelation> roles = roleRelationService.list(new QueryWrapper<RoleRelation>()
                    .eq("record_id", svResp.getId()));
            roles.stream().forEach(roleRelation -> {
                RoleDTO roleDTO = new RoleDTO();
                roleDTO.setRoleName(roleRelation.getRoleName());
                roleDTO.setDisplayName(roleRelation.getRoleName());
                if (RoleTypeEnum.USE_ROLE.getType() == roleRelation.getRoleType()) {
                    svResp.getUseRoles().add(roleDTO);
                } else {
                    svResp.getManageRoles().add(roleDTO);
                }
            });
        }
        return svResps;
    }

    @Override
    public TaskTemplateResp selectTaskTemplateOne(String id) {
        TaskTemplate taskTemplate = taskTemplateMapper.selectById(id);
        TaskTemplateResp taskTemplateResp = taskTemplateConverter.toDto(taskTemplate);
        List<RoleRelation> relations = roleRelationService.list(new QueryWrapper<RoleRelation>()
                .eq("record_id", id));
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
    public QueryResponse<TaskTemplateByRoleResp> selectTaskTemplateByRole(Integer page, Integer pageSize) {
        String currentUserRolesToString = AuthenticationContextHolder.getCurrentUserRolesToString();
        QueryRoleRelationBaseReq req = new QueryRoleRelationBaseReq();
        req.setSourceTableFix("rt");
        req.setUseRoleName(currentUserRolesToString);
        String sql = req.getConditionSql();

        IPage<TaskTemplate> iPage = taskTemplateMapper.selectSynthesisRequestTemple(new Page<TaskTemplate>(page, pageSize), sql);

        List<TaskTemplateByRoleResp> srt = synthesisTaskTemplateConverter.toDto(iPage.getRecords());

        QueryResponse<TaskTemplateByRoleResp> queryResponse = new QueryResponse<>();
        queryResponse.setPageInfo(new PageInfo(iPage.getTotal(), iPage.getCurrent(), iPage.getSize()));
        queryResponse.setContents(srt);
        return queryResponse;
    }


}

package com.webank.taskman.service.impl;


import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.constant.RoleTypeEnum;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.RoleDTO;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;
import com.webank.taskman.dto.resp.TaskTemplateResp;
import com.webank.taskman.mapper.TaskTemplateMapper;
import com.webank.taskman.service.RoleRelationService;
import com.webank.taskman.service.TaskTemplateService;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.BeanUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;


@Service
public class TaskTemplateServiceImpl extends ServiceImpl<TaskTemplateMapper, TaskTemplate> implements TaskTemplateService {

    @Autowired
    TaskTemplateMapper taskTemplateMapper;

    @Autowired
    RoleRelationService roleRelationService;

    private static final String TABLE_NAME = "task_template";

    @Override
    public TaskTemplateResp saveTaskTemplate(SaveTaskTemplateReq taskTemplateReq) {
        TaskTemplate taskTemplate = new TaskTemplate();
        BeanUtils.copyProperties(taskTemplateReq, taskTemplate);
        String currentUsername = AuthenticationContextHolder.getCurrentUsername();
        taskTemplate.setUpdatedBy(currentUsername);
        if (StringUtils.isEmpty(taskTemplate.getId())) {
            taskTemplate.setCreatedBy(currentUsername);
        }
        saveOrUpdate(taskTemplate);
        String taskTemplateID = taskTemplate.getId();
        roleRelationService.deleteByTemplate(TABLE_NAME, taskTemplateID);
        saveRoleRelation(taskTemplateID, RoleTypeEnum.USE_ROLE.getType(), taskTemplateReq.getUseRoles());
        saveRoleRelation(taskTemplateID, RoleTypeEnum.MANAGE_ROLE.getType(), taskTemplateReq.getManageRoles());

        TaskTemplateResp taskTemplateResp = new TaskTemplateResp();
        taskTemplateResp.setId(taskTemplateID);
        return taskTemplateResp;
    }

    public void saveRoleRelation(String taskTemplateID, int roleType, List<RoleDTO> roles) {
        roles.stream().forEach(role -> {
            if (!org.springframework.util.StringUtils.isEmpty(role.getRoleName()) && !org.springframework.util.StringUtils.isEmpty(role.getDisplayName())) {
                roleRelationService.save(new RoleRelation(
                        TABLE_NAME, taskTemplateID, roleType, role.getRoleName(), role.getDisplayName()));
            }
        });
    }

}

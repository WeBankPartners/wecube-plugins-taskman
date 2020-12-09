package com.webank.taskman.service.impl;


import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.constant.RoleTypeEnum;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.RoleDTO;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;
import com.webank.taskman.dto.resp.TaskTemplateResp;
import com.webank.taskman.mapper.TaskTemplateMapper;
import com.webank.taskman.service.FormTemplateService;
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
    RoleRelationService roleRelationService;

    @Autowired
    FormTemplateService formTemplateService;


    @Override
    public TaskTemplateResp saveTaskTemplateByReq(SaveTaskTemplateReq req) {
        TaskTemplate taskTemplate = new TaskTemplate();
        BeanUtils.copyProperties(req, taskTemplate);
        String currentUsername = AuthenticationContextHolder.getCurrentUsername();
        taskTemplate.setUpdatedBy(currentUsername);
        if (StringUtils.isEmpty(taskTemplate.getId())) {
            taskTemplate.setCreatedBy(currentUsername);
        }
        saveOrUpdate(taskTemplate);
        String taskTemplateId = taskTemplate.getId();
        roleRelationService.saveRoleRelationByTemplate(taskTemplateId,req.getUseRoles(),req.getManageRoles());

        formTemplateService.saveFormTemplateByReq(req.getForm());

        TaskTemplateResp taskTemplateResp = new TaskTemplateResp();
        taskTemplateResp.setId(taskTemplateId);
        return new TaskTemplateResp().setId(taskTemplateId);
    }


}

package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.converter.TaskTemplateConverter;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.PageInfo;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;
import com.webank.taskman.dto.req.SelectTaskTemplateRep;
import com.webank.taskman.dto.resp.TaskTemplateResp;
import com.webank.taskman.mapper.TaskTemplateMapper;
import com.webank.taskman.service.TaskTemplateService;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.BeanUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.List;


@Service
public class TaskTemplateServiceImpl extends ServiceImpl<TaskTemplateMapper, TaskTemplate> implements TaskTemplateService {

    @Autowired
    TaskTemplateMapper taskTemplateMapper;

    @Autowired
    TaskTemplateConverter taskTemplateConverter;

    @Override
    public TaskTemplate saveTaskTemplate(SaveTaskTemplateReq taskTemplateReq) {
        TaskTemplate taskTemplate = new TaskTemplate();
        BeanUtils.copyProperties(taskTemplateReq, taskTemplate);

        if (StringUtils.isEmpty(taskTemplate.getId())) {
            taskTemplate.setCreatedBy("11");
            taskTemplate.setUpdatedBy("22");
            taskTemplateMapper.insert(taskTemplate);
            return taskTemplateMapper.selectById(taskTemplate);
        }
        if (!StringUtils.isEmpty(taskTemplate.getId())) {
            taskTemplate.setUpdatedTime(new Date());
            taskTemplateMapper.updateById(taskTemplate);
            return taskTemplateMapper.selectById(taskTemplate);
        }
        return null;
    }

    @Override
    public void deleteTaskTemplateByIDService(String id) {
        taskTemplateMapper.deleteTaskTemplateByIDMapper(id);
    }

    @Override
    public QueryResponse<TaskTemplateResp> selectTaskTemplate(Integer currentPage, Integer pageSize, SelectTaskTemplateRep rep) {
        Page<TaskTemplate> page = new Page<>(currentPage, pageSize);
        QueryWrapper<TaskTemplate> wrapper = new QueryWrapper<>();
        if (!StringUtils.isEmpty(rep.getId())) {
            wrapper.eq("id", rep.getId());
        }
        if (!StringUtils.isEmpty(rep.getProcDefId())) {
            wrapper.eq("proc_def_id", rep.getProcDefId());
        }
        if (!StringUtils.isEmpty(rep.getProcDefKey())) {
            wrapper.eq("proc_def_key", rep.getProcDefKey());
        }
        if (!StringUtils.isEmpty(rep.getProcDefName())) {
            wrapper.like("proc_def_name", rep.getProcDefName());
        }
        if (!StringUtils.isEmpty(rep.getProcNode())) {
            wrapper.eq("proc_node", rep.getProcNode());
        }
        if (!StringUtils.isEmpty(rep.getName())) {
            wrapper.like("name", rep.getName());
        }
        IPage<TaskTemplate> iPage = taskTemplateMapper.selectPage(page, wrapper);
        List<TaskTemplate> records = iPage.getRecords();
        List<TaskTemplateResp> taskTemplateResps = taskTemplateConverter.toDto(records);

        QueryResponse<TaskTemplateResp> queryResponse = new QueryResponse<>();
        PageInfo pageInfo = new PageInfo();
        pageInfo.setStartIndex(iPage.getCurrent());
        pageInfo.setPageSize(iPage.getSize());
        pageInfo.setTotalRows(iPage.getTotal());
        queryResponse.setPageInfo(pageInfo);
        queryResponse.setContents(taskTemplateResps);
        return queryResponse;
    }

    @Override
    public List<TaskTemplateResp> selectTaskTemplateAll() {
        List<TaskTemplate> taskTemplates = taskTemplateMapper.selectList(null);
        List<TaskTemplateResp> taskTemplateResps = taskTemplateConverter.toDto(taskTemplates);
        return taskTemplateResps;
    }

    @Override
    public TaskTemplateResp selectTaskTemplateOne(String id) {
        TaskTemplate taskTemplate = taskTemplateMapper.selectById(id);
        TaskTemplateResp taskTemplateResp=new TaskTemplateResp();
        BeanUtils.copyProperties(taskTemplate,taskTemplateResp);
        return taskTemplateResp;
    }
}

package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.conditions.update.UpdateWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanException;
import com.webank.taskman.converter.RequestTemplateConverter;
import com.webank.taskman.converter.RoleInfoConverTer;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.domain.RequestTemplateRole;
import com.webank.taskman.domain.RoleInfo;
import com.webank.taskman.dto.PageInfo;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.SaveRequestTemplateReq;
import com.webank.taskman.dto.resp.RequestTemplateResp;
import com.webank.taskman.mapper.RequestTemplateMapper;
import com.webank.taskman.mapper.RequestTemplateRoleMapper;
import com.webank.taskman.service.RequestTemplateRoleService;
import com.webank.taskman.service.RequestTemplateService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.util.ArrayList;
import java.util.List;


@Service
public class RequestTemplateServiceImpl extends ServiceImpl<RequestTemplateMapper, RequestTemplate> implements RequestTemplateService {

    @Autowired
    RequestTemplateMapper requestTemplateMapper;

    @Autowired
    RequestTemplateConverter requestTemplateConverter;

    @Autowired
    RequestTemplateRoleService requestTemplateRoleService;

    @Autowired
    RoleInfoConverTer roleInfoConverTer;

    @Autowired
    RequestTemplateRoleMapper requestTemplateRoleMapper;

    @Override
    public RequestTemplateResp saveRequestTemplate(SaveRequestTemplateReq requestTemplateReq) {
        RequestTemplate requestTemplate = requestTemplateConverter.reqToDomain(requestTemplateReq);
        if (null == requestTemplateReq.getRoleIds()) {
            throw new TaskmanException("roleIds is null");
        }
        String currentUsername = AuthenticationContextHolder.getCurrentUsername();
        requestTemplate.setUpdatedBy(currentUsername);
        if (StringUtils.isEmpty(requestTemplate.getId())) {
            requestTemplate.setCreatedBy(currentUsername);
        }
        saveOrUpdate(requestTemplate);
        List<RequestTemplateRole> roles = new ArrayList<>();
        List<RequestTemplateRole> managenmentRoles = new ArrayList<>();
        String requestTemplateId = requestTemplate.getId();

        requestTemplateRoleService.deleteByRequestTemplate(requestTemplateId);
        for (String roleId : requestTemplateReq.getRoleIds()) {
            roles.add(new RequestTemplateRole(requestTemplateId, roleId, 1));
        }

        for (String managenmentRole : requestTemplateReq.getManagementRole()) {
            managenmentRoles.add(new RequestTemplateRole(requestTemplateId, managenmentRole, 0));
        }
        requestTemplateRoleService.saveBatch(roles);
        requestTemplateRoleService.saveBatch(managenmentRoles);

        return new RequestTemplateResp().setId(requestTemplateId);

    }

    @Override
    public void deleteRequestTemplateService(String id) throws Exception {
        if (id == "") {
            throw new Exception("Request template parameter cannot be ID");
        }
        UpdateWrapper<RequestTemplate> wrapper = new UpdateWrapper<>();
        wrapper.eq("id", id).set("del_flag", 1);
        requestTemplateMapper.update(null, wrapper);
    }

    @Override
    public QueryResponse<RequestTemplateResp> selectAllequestTemplateService(Integer current, Integer limit, SaveRequestTemplateReq req) throws Exception {
        Page<RequestTemplate> page = new Page<>(current, limit);
        QueryWrapper<RequestTemplate> wrapper = new QueryWrapper<>();
        wrapper.eq(!StringUtils.isEmpty(req.getId()), "id", req.getId());
        wrapper.eq(!StringUtils.isEmpty(req.getRequestTempGroup()), "request_temp_group", req.getRequestTempGroup());
        wrapper.eq(!StringUtils.isEmpty(req.getProcDefKey()), "proc_def_key", req.getProcDefKey());
        wrapper.eq(!StringUtils.isEmpty(req.getProcDefId()), "proc_def_id", req.getProcDefId());
        wrapper.eq(!StringUtils.isEmpty(req.getProcDefName()), "proc_def_name", req.getProcDefName());
        wrapper.eq(!StringUtils.isEmpty(req.getName()), "name", req.getName());
        IPage<RequestTemplate> iPage = requestTemplateMapper.selectPage(page, wrapper);
        List<RequestTemplate> records = iPage.getRecords();
        List<RequestTemplateResp> requestTemplateDTOS = requestTemplateConverter.toDto(records);
        for (RequestTemplateResp requestTemplateDTO : requestTemplateDTOS) {
            RequestTemplateRole requestTemplateRole1=new RequestTemplateRole();
            List<RequestTemplateRole> roles= requestTemplateRoleMapper.selectList(new QueryWrapper<RequestTemplateRole>().eq("request_template_id",requestTemplateDTO.getId()).eq("role_type",1));
            List<RequestTemplateRole> management= requestTemplateRoleMapper.selectList(new QueryWrapper<RequestTemplateRole>().eq("request_template_id",requestTemplateDTO.getId()).eq("role_type",0));
            List<RoleInfo> rolesId=roleInfoConverTer.toDto(roles);
            List<RoleInfo> managementRole=roleInfoConverTer.toDto(management);
            requestTemplateDTO.setRoleIds(rolesId);
            requestTemplateDTO.setManagementRole(managementRole);
        }

        QueryResponse<RequestTemplateResp> queryResponse = new QueryResponse<>();
        PageInfo pageInfo = new PageInfo();
        pageInfo.setStartIndex(iPage.getCurrent());
        pageInfo.setPageSize(iPage.getSize());
        pageInfo.setTotalRows(iPage.getTotal());
        queryResponse.setPageInfo(pageInfo);
        queryResponse.setContents(requestTemplateDTOS);
        return queryResponse;
    }

    @Override
    public RequestTemplateResp detailRequestTemplate(String id) throws Exception {
        RequestTemplate requestTemplate = requestTemplateMapper.selectById(id);
        RequestTemplateResp formTemplateResp = requestTemplateConverter.toDto(requestTemplate);
        return formTemplateResp;
    }


}

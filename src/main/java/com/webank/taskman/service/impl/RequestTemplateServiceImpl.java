package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.Wrapper;
import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.conditions.segments.MergeSegments;
import com.baomidou.mybatisplus.core.conditions.update.UpdateWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanException;
import com.webank.taskman.constant.RoleTypeEnum;
import com.webank.taskman.converter.RequestTemplateConverter;
import com.webank.taskman.converter.RoleInfoConverTer;
import com.webank.taskman.converter.RoleRelationConverter;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.domain.RequestTemplateRole;
import com.webank.taskman.domain.RoleInfo;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.dto.PageInfo;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.RoleDTO;
import com.webank.taskman.dto.req.SaveRequestTemplateReq;
import com.webank.taskman.dto.resp.RequestTemplateResp;
import com.webank.taskman.mapper.RequestTemplateMapper;
import com.webank.taskman.mapper.RequestTemplateRoleMapper;
import com.webank.taskman.service.RequestTemplateRoleService;
import com.webank.taskman.service.RequestTemplateService;
import com.webank.taskman.service.RoleRelationService;
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
    RoleRelationService roleRelationService;

    @Autowired
    RoleRelationConverter roleRelationConverter;


    private static final String TABLE_NAME = "request_template";
    @Override
    public RequestTemplateResp saveRequestTemplate(SaveRequestTemplateReq requestTemplateReq) {
        RequestTemplate requestTemplate = requestTemplateConverter.reqToDomain(requestTemplateReq);
        String currentUsername = AuthenticationContextHolder.getCurrentUsername();
        requestTemplate.setUpdatedBy(currentUsername);
        if (StringUtils.isEmpty(requestTemplate.getId())) {
            requestTemplate.setCreatedBy(currentUsername);
        }
        saveOrUpdate(requestTemplate);
        String requestTemplateId = requestTemplate.getId();
        roleRelationService.deleteByTemplate(TABLE_NAME,requestTemplateId);
        requestTemplateReq.getUseRoles().stream().forEach(useRole-> roleRelationService.save( new RoleRelation(
                TABLE_NAME,requestTemplateId, RoleTypeEnum.USE_ROLE.getType(),useRole.getRoleName(),useRole.getDisplayName())));

        requestTemplateReq.getUseRoles().stream().forEach(manageRole-> roleRelationService.save( new RoleRelation(
                TABLE_NAME,requestTemplateId, RoleTypeEnum.MANAGE_ROLE.getType(),manageRole.getRoleName(),manageRole.getDisplayName())));

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


        // 分页实现
        IPage<RequestTemplate> iPage = requestTemplateMapper.selectPageByParam(page,null);

//        IPage<RequestTemplate> iPage = requestTemplateMapper.selectPage(page, wrapper);


        List<RequestTemplate> records = iPage.getRecords();
        List<RequestTemplateResp> respList = requestTemplateConverter.toDto(records);

        for (RequestTemplateResp resp : respList) {
            List<RoleRelation> roles = roleRelationService.list(
                    new QueryWrapper<RoleRelation>().eq("record_table",TABLE_NAME).eq("record_id",resp.getId()));
            roles.stream().forEach(role->{
                RoleDTO roleDTO = roleRelationConverter.toDto(role);
                if(RoleTypeEnum.USE_ROLE.getType() == role.getRoleType()){
                    resp.getUseRoles().add(roleDTO);
                }else{
                    resp.getManageRoles().add(roleDTO);
                }
            });
        }
        QueryResponse<RequestTemplateResp> queryResponse = new QueryResponse<>();
        PageInfo pageInfo = new PageInfo();
        pageInfo.setStartIndex(iPage.getCurrent());
        pageInfo.setPageSize(iPage.getSize());
        pageInfo.setTotalRows(iPage.getTotal());
        queryResponse.setPageInfo(pageInfo);
        queryResponse.setContents(respList);
        return queryResponse;
    }

    @Override
    public RequestTemplateResp detailRequestTemplate(String id) throws Exception {
        RequestTemplate requestTemplate = requestTemplateMapper.selectById(id);
        RequestTemplateResp formTemplateResp = requestTemplateConverter.toDto(requestTemplate);
        return formTemplateResp;
    }


}

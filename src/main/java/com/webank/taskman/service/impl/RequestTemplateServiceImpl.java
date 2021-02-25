package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.update.UpdateWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.github.pagehelper.PageHelper;
import com.github.pagehelper.PageInfo;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.RoleTypeEnum;
import com.webank.taskman.constant.StatusEnum;
import com.webank.taskman.constant.TemplateTypeEnum;
import com.webank.taskman.converter.FormItemTemplateConverter;
import com.webank.taskman.converter.FormTemplateConverter;
import com.webank.taskman.converter.RequestTemplateConverter;
import com.webank.taskman.converter.RoleRelationConverter;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.dto.RequestTemplateDto;
import com.webank.taskman.dto.RoleDto;
import com.webank.taskman.dto.req.RequestTemplateQueryReqDto;
import com.webank.taskman.dto.req.RequestTemplateSaveReqDto;
import com.webank.taskman.dto.resp.FormTemplateRespDto;
import com.webank.taskman.dto.resp.RequestTemplateRespDto;
import com.webank.taskman.mapper.FormItemTemplateMapper;
import com.webank.taskman.mapper.FormTemplateMapper;
import com.webank.taskman.mapper.RequestTemplateMapper;
import com.webank.taskman.service.RequestTemplateService;
import com.webank.taskman.service.RoleRelationService;
import com.webank.taskman.support.core.PlatformCoreServiceRestClient;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.StringUtils;

import java.util.*;

@Service
public class RequestTemplateServiceImpl extends ServiceImpl<RequestTemplateMapper, RequestTemplate>
        implements RequestTemplateService {

    @Autowired
    private RequestTemplateMapper requestTemplateMapper;

    @Autowired
    private RequestTemplateConverter requestTemplateConverter;

    @Autowired
    private RoleRelationService roleRelationService;

    @Autowired
    private RoleRelationConverter roleRelationConverter;

    @Autowired
    private FormTemplateConverter formTemplateConverter;

    @Autowired
    private FormTemplateMapper formTemplateMapper;

    @Autowired
    private FormItemTemplateMapper formItemTemplateMapper;

    @Autowired
    private PlatformCoreServiceRestClient coreServiceStub;

    @Autowired
    private FormItemTemplateConverter formItemTemplateConverter;

    @Override
    @Transactional
    public RequestTemplateDto saveRequestTemplate(RequestTemplateSaveReqDto req) {
        List<RoleDto> roles = roleRelationConverter
                .rolesDataResponseToDtoList(coreServiceStub.getAllAuthRolesOfCurrentUser());
        Set<RoleDto> ts = new HashSet<RoleDto>(roles);
        ts.addAll(req.getUseRoles());
        req.setUseRoles(new ArrayList<>(ts));
        RequestTemplate requestTemplate = requestTemplateConverter.saveReqToEntity(req);
        String currentUsername = AuthenticationContextHolder.getCurrentUsername();
        requestTemplate.setUpdatedBy(currentUsername);
        if (StringUtils.isEmpty(requestTemplate.getId())) {
            requestTemplate.setCreatedBy(currentUsername);
        }
        saveOrUpdate(requestTemplate);
        roleRelationService.saveRoleRelationByTemplate(requestTemplate.getId(), req.getUseRoles(),
                req.getManageRoles());
        return new RequestTemplateDto().setId(requestTemplate.getId());

    }

    @Override
    public void deleteRequestTemplateService(String id) throws TaskmanRuntimeException {
        if (StringUtils.isEmpty(id)) {
            throw new TaskmanRuntimeException("Request template parameter cannot be ID");
        }
        UpdateWrapper<RequestTemplate> wrapper = new UpdateWrapper<>();
        wrapper.lambda().eq(RequestTemplate::getId, id).set(RequestTemplate::getDelFlag, StatusEnum.ENABLE.ordinal())
                .set(RequestTemplate::getUpdatedTime, new Date());
        ;
        requestTemplateMapper.update(null, wrapper);
    }

    @Override
    public QueryResponse<RequestTemplateDto> selectRequestTemplatePage(Integer pageNum, Integer pageSize,
            RequestTemplateQueryReqDto req) {
        req.queryCurrentUserRoles();
        PageHelper.startPage(pageNum, pageSize);
        PageInfo<RequestTemplateDto> pages = new PageInfo(requestTemplateMapper.selectDTOListByParam(req));

        for (RequestTemplateDto resp : pages.getList()) {
            List<RoleRelation> roles = roleRelationService
                    .list(new RoleRelation().setRecordId(resp.getId()).getLambdaQueryWrapper());
            roles.stream().forEach(role -> {
                RoleDto roleDTO = roleRelationConverter.toDto(role);
                if (RoleTypeEnum.USE_ROLE.getType() == role.getRoleType()) {
                    resp.getUseRoles().add(roleDTO);
                } else {
                    resp.getManageRoles().add(roleDTO);
                }
            });
        }
        QueryResponse<RequestTemplateDto> queryResponse = new QueryResponse(pages.getTotal(), pageNum.longValue(),
                pageSize.longValue(), pages.getList());
        return queryResponse;
    }

    @Override
    public RequestTemplateRespDto detailRequestTemplate(String id) {
        RequestTemplateRespDto requestTemplateResp = requestTemplateConverter
                .toRespByEntity(requestTemplateMapper.selectById(id));
        FormTemplateRespDto formTemplateResp = formTemplateConverter
                .toDto(formTemplateMapper.selectOne(new FormTemplate().setTempId(requestTemplateResp.getId())
                        .setTempType(TemplateTypeEnum.REQUEST.getType()).getLambdaQueryWrapper()));
        List<FormItemTemplate> items = formItemTemplateMapper
                .selectList(new FormItemTemplate().setFormTemplateId(formTemplateResp.getId()).getLambdaQueryWrapper());
        formTemplateResp.setItems(formItemTemplateConverter.toRespByEntity(items));
        requestTemplateResp.setFormTemplateResp(formTemplateResp);
        return requestTemplateResp;
    }

}

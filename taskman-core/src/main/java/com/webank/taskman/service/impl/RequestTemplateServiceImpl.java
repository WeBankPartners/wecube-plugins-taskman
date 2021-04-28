package com.webank.taskman.service.impl;

import java.util.ArrayList;
import java.util.Collections;
import java.util.Date;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.core.conditions.update.UpdateWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.github.pagehelper.PageHelper;
import com.github.pagehelper.PageInfo;
import com.webank.taskman.base.PageableQueryResult;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.RoleType;
import com.webank.taskman.constant.GenernalStatus;
import com.webank.taskman.constant.TemplateType;
import com.webank.taskman.converter.FormItemTemplateConverter;
import com.webank.taskman.converter.FormTemplateConverter;
import com.webank.taskman.converter.RequestTemplateConverter;
import com.webank.taskman.converter.RoleRelationConverter;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.dto.RequestTemplateDto;
import com.webank.taskman.dto.RoleDto;
import com.webank.taskman.dto.req.RequestTemplateQueryReqDto;
import com.webank.taskman.dto.req.RequestTemplateSaveReqDto;
import com.webank.taskman.dto.req.RoleRelationBaseQueryReqDto;
import com.webank.taskman.dto.resp.FormTemplateRespDto;
import com.webank.taskman.dto.resp.RequestTemplateRespDto;
import com.webank.taskman.mapper.FormItemTemplateMapper;
import com.webank.taskman.mapper.FormTemplateMapper;
import com.webank.taskman.mapper.RequestTemplateMapper;
import com.webank.taskman.service.RequestTemplateGroupService;
import com.webank.taskman.service.RequestTemplateService;
import com.webank.taskman.service.RoleRelationService;
import com.webank.taskman.support.core.PlatformCoreServiceRestClient;

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
    
    @Autowired
    private RequestTemplateGroupService requestTemplateGroupService;

    @Override
    @Transactional
    public RequestTemplateDto saveRequestTemplate(RequestTemplateSaveReqDto req) {
        List<RoleDto> roles = roleRelationConverter
                .rolesDataResponseToDtoList(coreServiceStub.getAllAuthRolesOfCurrentUser());

        Set<String> toAddUseRoles = new HashSet<String>();
        if (roles != null) {
            for (RoleDto roleDto : roles) {
                if (StringUtils.isNoneBlank(roleDto.getRoleName())) {
                    toAddUseRoles.add(roleDto.getRoleName());
                }
            }
        }

        List<RoleDto> inputRoleDtos = req.getUseRoles();
        if (inputRoleDtos != null) {
            for (RoleDto roleDto : inputRoleDtos) {
                if (StringUtils.isNoneBlank(roleDto.getRoleName())) {
                    toAddUseRoles.add(roleDto.getRoleName());
                }
            }
        }

        req.setUseRoles(new ArrayList<>(inputRoleDtos));
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
        wrapper.lambda().eq(RequestTemplate::getId, id).set(RequestTemplate::getDelFlag, GenernalStatus.ENABLE.ordinal())
                .set(RequestTemplate::getUpdatedTime, new Date());
        ;
        requestTemplateMapper.update(null, wrapper);
    }

    @Override
    public PageableQueryResult<RequestTemplateDto> selectRequestTemplatePage(Integer pageNum, Integer pageSize,
            RequestTemplateQueryReqDto req) {
        req.queryCurrentUserRoles();
        PageHelper.startPage(pageNum, pageSize);
        PageInfo<RequestTemplateDto> pages = new PageInfo<>(requestTemplateMapper.selectDTOListByParam(req));

        for (RequestTemplateDto resp : pages.getList()) {
            List<RoleRelation> roles = roleRelationService
                    .list(new RoleRelation().setRecordId(resp.getId()).getLambdaQueryWrapper());
            roles.stream().forEach(role -> {
                RoleDto roleDTO = roleRelationConverter.convertToDto(role);
                if (RoleType.USE_ROLE.getType() == role.getRoleType()) {
                    resp.getUseRoles().add(roleDTO);
                } else {
                    resp.getManageRoles().add(roleDTO);
                }
            });
        }
        PageableQueryResult<RequestTemplateDto> queryResponse = new PageableQueryResult<>(pages.getTotal(), pageNum.longValue(),
                pageSize.longValue(), pages.getList());
        return queryResponse;
    }

    @Override
    public RequestTemplateRespDto detailRequestTemplate(String id) {
        RequestTemplateRespDto requestTemplateResp = requestTemplateConverter
                .toRespByEntity(requestTemplateMapper.selectById(id));
        FormTemplateRespDto formTemplateResp = formTemplateConverter
                .convertToDto(formTemplateMapper.selectOne(new FormTemplate().setTempId(requestTemplateResp.getId())
                        .setTempType(TemplateType.REQUEST.getType()).getLambdaQueryWrapper()));
        List<FormItemTemplate> items = formItemTemplateMapper
                .selectList(new FormItemTemplate().setFormTemplateId(formTemplateResp.getId()).getLambdaQueryWrapper());
        formTemplateResp.setItems(formItemTemplateConverter.toRespByEntity(items));
        requestTemplateResp.setFormTemplateResp(formTemplateResp);
        return requestTemplateResp;
    }

    /**
     * 
     */
    @Override
    public RequestTemplateDto releaseRequestTemplate(RequestTemplateSaveReqDto req) {
        if (StringUtils.isBlank(req.getId())) {
            throw new TaskmanRuntimeException("Request template ID should provide.");
        }
        RequestTemplate requestTemplate = this.getOne(new RequestTemplate().setId(req.getId()).getLambdaQueryWrapper());
        if (requestTemplate == null) {
            throw new TaskmanRuntimeException("Request template does not exist.");
        }
        String status = RequestTemplate.STATUS_UNRELEASED.equals(requestTemplate.getStatus())
                ? RequestTemplate.STATUS_RELEASED : RequestTemplate.STATUS_UNRELEASED;
        requestTemplate.setStatus(status);
        requestTemplate.setUpdatedBy(AuthenticationContextHolder.getCurrentUsername());
        requestTemplate.setUpdatedTime(new Date());
        this.updateById(requestTemplate);

        RequestTemplateDto respDto = new RequestTemplateDto().setId(requestTemplate.getId())
                .setStatus(requestTemplate.getStatus());
        return respDto;
    }

    @Override
    public List<RequestTemplateDto> fetchAvailableRequestTemplates() {
        RequestTemplate requestTemplate = new RequestTemplate().setStatus(RequestTemplate.STATUS_RELEASED);
        String inSql = RoleRelationBaseQueryReqDto.getEqUseRole();
        LambdaQueryWrapper<RequestTemplate> queryWrapper = requestTemplate.getLambdaQueryWrapper()
                .inSql(!StringUtils.isEmpty(inSql), RequestTemplate::getId, RoleRelationBaseQueryReqDto.getEqUseRole());
        
        List<RequestTemplate> requestTemplateEntities = this.list(queryWrapper);
        
        if(requestTemplateEntities == null ) {
            return Collections.emptyList();
        }
        
        List<RequestTemplateDto> retRequestTemplateDtos = requestTemplateConverter.convertToDtos(requestTemplateEntities);
        
        for(RequestTemplateDto retDto : retRequestTemplateDtos) {
            if(StringUtils.isNotBlank(retDto.getRequestTempGroup())) {
                RequestTemplateGroup rtGroupEntity = requestTemplateGroupService.getById(retDto.getRequestTempGroup());
                if(rtGroupEntity != null) {
                    retDto.setRequestTempGroupName(rtGroupEntity.getName());
                }
            }
        }
        return retRequestTemplateDtos;
    }

}

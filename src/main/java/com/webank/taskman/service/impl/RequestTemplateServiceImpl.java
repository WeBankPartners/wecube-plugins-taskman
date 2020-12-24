package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.update.UpdateWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.github.pagehelper.PageHelper;
import com.github.pagehelper.PageInfo;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.RoleTypeEnum;
import com.webank.taskman.converter.FormTemplateConverter;
import com.webank.taskman.converter.RequestTemplateConverter;
import com.webank.taskman.converter.RoleRelationConverter;
import com.webank.taskman.domain.FormItemTemplate;
import com.webank.taskman.domain.FormTemplate;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.dto.RequestTemplateDTO;
import com.webank.taskman.dto.RoleDTO;
import com.webank.taskman.dto.req.QueryRequestTemplateReq;
import com.webank.taskman.dto.req.SaveRequestTemplateReq;
import com.webank.taskman.dto.resp.DetailRequestTemplateResq;
import com.webank.taskman.mapper.FormItemTemplateMapper;
import com.webank.taskman.mapper.FormTemplateMapper;
import com.webank.taskman.mapper.RequestTemplateMapper;
import com.webank.taskman.service.RequestTemplateService;
import com.webank.taskman.service.RoleRelationService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.StringUtils;

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

    @Autowired
    FormTemplateConverter formTemplateConverter;

    @Autowired
    FormTemplateMapper formTemplateMapper;

    @Autowired
    FormItemTemplateMapper formItemTemplateMapper;


    @Override
    @Transactional
    public RequestTemplateDTO saveRequestTemplate(SaveRequestTemplateReq req) {
        RequestTemplate requestTemplate = requestTemplateConverter.saveReqToEntity(req);
        String currentUsername = AuthenticationContextHolder.getCurrentUsername();
        requestTemplate.setUpdatedBy(currentUsername);
        if (StringUtils.isEmpty(requestTemplate.getId())) {
            requestTemplate.setCreatedBy(currentUsername);
        }
        saveOrUpdate(requestTemplate);
        roleRelationService.saveRoleRelationByTemplate(requestTemplate.getId(), req.getUseRoles(), req.getManageRoles());
        return new RequestTemplateDTO().setId(requestTemplate.getId());

    }


    @Override
    public void deleteRequestTemplateService(String id) throws TaskmanRuntimeException {
        if (StringUtils.isEmpty(id)) {
            throw new TaskmanRuntimeException("Request template parameter cannot be ID");
        }
        UpdateWrapper<RequestTemplate> wrapper = new UpdateWrapper<>();
        wrapper.lambda().eq(RequestTemplate::getId, id).set(RequestTemplate::getDelFlag, 1);
        requestTemplateMapper.update(null, wrapper);
    }

    @Override
    public QueryResponse<RequestTemplateDTO> selectRequestTemplatePage(Integer pageNum, Integer pageSize, QueryRequestTemplateReq req) {
        req.setSourceTableFix("rt");
        StringBuffer useRole = new StringBuffer().append(StringUtils.isEmpty(req.getUseRoleName()) ? "" : req.getUseRoleName() + ",");
        req.setUseRoleName(useRole.append(AuthenticationContextHolder.getCurrentUserRolesToString()).toString());

        PageHelper.startPage(pageNum,pageSize);
        PageInfo<RequestTemplateDTO> pages = new PageInfo(requestTemplateMapper.selectDTOListByParam(req));

        for (RequestTemplateDTO resp : pages.getList()) {
            List<RoleRelation> roles = roleRelationService.list(new RoleRelation().setRecordId(resp.getId()).getLambdaQueryWrapper());
            roles.stream().forEach(role -> {
                RoleDTO roleDTO = roleRelationConverter.toDto(role);
                if (RoleTypeEnum.USE_ROLE.getType() == role.getRoleType()) {
                    resp.getUseRoles().add(roleDTO);
                } else {
                    resp.getManageRoles().add(roleDTO);
                }
            });
        }
        QueryResponse<RequestTemplateDTO> queryResponse = new QueryResponse(pages.getTotal(),pageNum.longValue(),pageSize.longValue(),pages.getList());
        return queryResponse;
    }

    @Override
    public DetailRequestTemplateResq detailRequestTemplate(String id) {
        DetailRequestTemplateResq detailRequestTemplateResq = requestTemplateConverter.detailRequest(requestTemplateMapper.selectById(id));
        detailRequestTemplateResq.setDetilReuestTemplateFormResq(
                formTemplateConverter.detailForm(
                        formTemplateMapper.selectOne(
                                new FormTemplate()
                                        .setTempId(detailRequestTemplateResq.getId())
                                        .setTempType("0").getLambdaQueryWrapper()
                        )));
        detailRequestTemplateResq.getDetilReuestTemplateFormResq().setFormItemTemplateList(
                formItemTemplateMapper.selectList(
                        new FormItemTemplate().setFormTemplateId(detailRequestTemplateResq.getDetilReuestTemplateFormResq().getId()).getLambdaQueryWrapper()
                ));
        return detailRequestTemplateResq;
    }


    @Override
    public List<RequestTemplateDTO> selectDTOListByParam(QueryRequestTemplateReq req) {
        List<RequestTemplateDTO> list = this.getBaseMapper().selectDTOListByParam(req);
        return list;
    }


}

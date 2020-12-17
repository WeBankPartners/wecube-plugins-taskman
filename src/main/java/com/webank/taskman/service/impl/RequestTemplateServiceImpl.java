package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.conditions.update.UpdateWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.constant.RoleTypeEnum;
import com.webank.taskman.converter.RequestTemplateConverter;
import com.webank.taskman.converter.RoleRelationConverter;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.domain.RoleRelation;
import com.webank.taskman.dto.PageInfo;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.RoleDTO;
import com.webank.taskman.dto.req.QueryRequestTemplateReq;
import com.webank.taskman.dto.req.SaveRequestTemplateReq;
import com.webank.taskman.dto.resp.RequestTemplateResp;
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


    private static final String TABLE_NAME = "request_template";

    @Override
    @Transactional
    public RequestTemplateResp saveRequestTemplate(SaveRequestTemplateReq req) {
        RequestTemplate requestTemplate = requestTemplateConverter.reqToDomain(req);
        String currentUsername = AuthenticationContextHolder.getCurrentUsername();
        requestTemplate.setUpdatedBy(currentUsername);
        if (StringUtils.isEmpty(requestTemplate.getId())) {
            requestTemplate.setCreatedBy(currentUsername);
        }
        saveOrUpdate(requestTemplate);
        roleRelationService.saveRoleRelationByTemplate(requestTemplate.getId(), req.getUseRoles(),req.getManageRoles());
        return new RequestTemplateResp().setId(requestTemplate.getId());

    }


    @Override
    public void deleteRequestTemplateService(String id) throws Exception {
        if (StringUtils.isEmpty(id)) {
            throw new Exception("Request template parameter cannot be ID");
        }
        UpdateWrapper<RequestTemplate> wrapper = new UpdateWrapper<>();
        wrapper.eq("id", id).set("del_flag", 1);
        requestTemplateMapper.update(null, wrapper);
    }

    @Override
    public QueryResponse<RequestTemplateResp> selectRequestTemplatePage(Integer current, Integer limit, QueryRequestTemplateReq req) throws Exception {
        req.setSourceTableFix("rt");
        IPage<RequestTemplate> iPage = requestTemplateMapper.selectPageByParam(new Page<>(current, limit),req);
        List<RequestTemplateResp> respList = requestTemplateConverter.toDto(iPage.getRecords());
        for (RequestTemplateResp resp : respList) {
            List<RoleRelation> roles = roleRelationService.list(new QueryWrapper<RoleRelation>().eq("record_id",resp.getId()));
            roles.stream().forEach(role->
            {
                RoleDTO roleDTO = roleRelationConverter.toDto(role);
                if(RoleTypeEnum.USE_ROLE.getType() == role.getRoleType()){
                    resp.getUseRoles().add(roleDTO);
                }else{
                    resp.getManageRoles().add(roleDTO);
                }
            });
        }
        QueryResponse<RequestTemplateResp> queryResponse = new QueryResponse<>();
        queryResponse.setPageInfo(new PageInfo(iPage.getTotal(),iPage.getCurrent(),iPage.getSize()));
        queryResponse.setContents(respList);
        return queryResponse;
    }

    @Override
    public RequestTemplateResp detailRequestTemplate(String id) throws Exception {
//        RequestTemplateResp formTemplateResp =
//        requestTemplateConverter.toDto(requestTemplateMapper.selectById(id));
        return requestTemplateConverter.toDto(requestTemplateMapper.selectById(id));
    }

    @Override
    public List<RequestTemplate> selectAvailableRequest(QueryRequestTemplateReq req) {
        return this.baseMapper.selectListByParam(req);
    }


}

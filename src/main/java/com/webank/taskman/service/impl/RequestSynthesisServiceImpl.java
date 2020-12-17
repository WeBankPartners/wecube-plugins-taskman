package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.converter.RequestTemplateConverter;
import com.webank.taskman.converter.SynthesisRequestInfoFormRequestConverter;
import com.webank.taskman.converter.SynthesisRequestInfoRespConverter;
import com.webank.taskman.converter.SynthesisRequestTemplateConverter;
import com.webank.taskman.domain.*;
import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.PageInfo;
import com.webank.taskman.dto.QueryResponse;
import com.webank.taskman.dto.req.QueryRoleRelationBaseReq;
import com.webank.taskman.dto.req.SynthesisRequestInfoReq;
import com.webank.taskman.dto.resp.SynthesisRequestInfoFormRequest;
import com.webank.taskman.dto.resp.SynthesisRequestInfoResp;
import com.webank.taskman.dto.resp.SynthesisRequestTempleResp;
import com.webank.taskman.mapper.FormInfoMapper;
import com.webank.taskman.mapper.FormItemInfoMapper;
import com.webank.taskman.mapper.RequestInfoMapper;
import com.webank.taskman.mapper.RequestTemplateMapper;
import com.webank.taskman.service.RequestSynthesisService;
import com.webank.taskman.utils.JsonUtils;
import com.webank.taskman.utils.UnderlineToCameUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@Service
public class RequestSynthesisServiceImpl extends ServiceImpl<RequestTemplateMapper, RequestTemplate> implements RequestSynthesisService {

    @Autowired
    RequestTemplateMapper requestTemplateMapper;

    @Autowired
    RequestInfoMapper requestInfoMapper;

    @Autowired
    FormItemInfoMapper formItemInfoMapper;

    @Autowired
    FormInfoMapper formInfoMapper;

    @Autowired
    SynthesisRequestTemplateConverter synthesisRequestTemplateConverter;

    @Autowired
    SynthesisRequestInfoRespConverter synthesisRequestInfoRespConverter;

    @Autowired
    SynthesisRequestInfoFormRequestConverter synthesisRequestInfoFormRequestConverter;

    @Override
    public QueryResponse<SynthesisRequestTempleResp> selectSynthesisRequestTempleService(Integer current, Integer limit) throws Exception {
        String currentUserRolesToString = AuthenticationContextHolder.getCurrentUserRolesToString();
        QueryRoleRelationBaseReq req = new QueryRoleRelationBaseReq();
        req.setSourceTableFix("rt");
        req.setUseRoleName(currentUserRolesToString);
        String sql = req.getConditionSql();

        IPage<RequestTemplate> iPage = requestTemplateMapper.selectSynthesisRequestTemple(new Page<>(current, limit),sql);
        List<SynthesisRequestTempleResp> srt=synthesisRequestTemplateConverter.toDto(iPage.getRecords());

        QueryResponse<SynthesisRequestTempleResp> queryResponse = new QueryResponse<>();
        queryResponse.setPageInfo(new PageInfo(iPage.getTotal(),iPage.getCurrent(),iPage.getSize()));
        queryResponse.setContents(srt);
        return queryResponse;
    }

    @Override
    public QueryResponse<SynthesisRequestInfoResp> selectSynthesisRequestInfoService(Integer current, Integer limit, SynthesisRequestInfoReq req) throws Exception {
        String currentUserRolesToString = AuthenticationContextHolder.getCurrentUserRolesToString();
        List<Map<String,Object>> list = new ArrayList<>();
        req.setRoleName(currentUserRolesToString);
        IPage<RequestInfo> iPage = requestInfoMapper.selectSynthesisRequestInfo(new Page<>(current, limit),req);
        List<SynthesisRequestInfoResp> srt=synthesisRequestInfoRespConverter.toDto(iPage.getRecords());

        QueryResponse<SynthesisRequestInfoResp> queryResponse = new QueryResponse<>();
        queryResponse.setPageInfo(new PageInfo(iPage.getTotal(),iPage.getCurrent(),iPage.getSize()));
        queryResponse.setContents(srt);

        return queryResponse;
    }

    @Override
    public SynthesisRequestInfoFormRequest selectSynthesisRequestInfoFormService(String id) throws Exception {

        FormInfo formInfo=formInfoMapper.selectOne(new QueryWrapper<FormInfo>().eq("record_id",id));
        if (null==formInfo||"".equals(formInfo)){
            throw new Exception("The request details do not exist");
        }
        List<FormItemInfo> formItemInfos=formItemInfoMapper.selectList(new QueryWrapper<FormItemInfo>().eq("form_id",formInfo.getId()));
        SynthesisRequestInfoFormRequest srt=synthesisRequestInfoFormRequestConverter.toDto(formInfo);
        srt.setFormItemInfo(formItemInfos);

        return srt;
    }

}

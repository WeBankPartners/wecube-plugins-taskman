package com.webank.taskman.service.impl;

import java.util.ArrayList;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.base.QueryResponse.PageInfo;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.StatusEnum;
import com.webank.taskman.converter.EntityAttrValueConverter;
import com.webank.taskman.converter.FormItemInfoConverter;
import com.webank.taskman.converter.RequestInfoConverter;
import com.webank.taskman.domain.FormItemInfo;
import com.webank.taskman.domain.RequestInfo;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.CreateTaskDto;
import com.webank.taskman.dto.CreateTaskDto.EntityAttrValueDto;
import com.webank.taskman.dto.CreateTaskDto.EntityValueDto;
import com.webank.taskman.dto.req.RequestInfoQueryReqDto;
import com.webank.taskman.dto.resp.RequestInfoResqDto;
import com.webank.taskman.mapper.RequestInfoMapper;
import com.webank.taskman.service.FormInfoService;
import com.webank.taskman.service.FormItemInfoService;
import com.webank.taskman.service.FormTemplateService;
import com.webank.taskman.service.RequestInfoService;
import com.webank.taskman.service.RequestTemplateService;
import com.webank.taskman.support.core.PlatformCoreServiceRestClient;
import com.webank.taskman.support.core.dto.DynamicEntityValueDto;
import com.webank.taskman.support.core.dto.DynamicEntityValueDto.DynamicEntityAttrValueDto;
import com.webank.taskman.support.core.dto.DynamicTaskNodeBindInfoDto;
import com.webank.taskman.support.core.dto.DynamicWorkflowInstCreationInfoDto;
import com.webank.taskman.support.core.dto.DynamicWorkflowInstInfoDto;

@Service
public class RequestInfoServiceImpl extends ServiceImpl<RequestInfoMapper, RequestInfo> implements RequestInfoService {

    @Autowired
    private FormInfoService formInfoService;
    @Autowired
    private PlatformCoreServiceRestClient platformCoreServiceRestClient;
    @Autowired
    private FormItemInfoService formItemInfoService;
    @Autowired
    private RequestTemplateService requestTemplateService;

    @Autowired
    private RequestInfoMapper requestInfoMapper;
    @Autowired
    private RequestInfoConverter requestInfoConverter;
    @Autowired
    private FormItemInfoConverter formItemInfoConverter;

    @Autowired
    private FormTemplateService formTemplateService;
    @Autowired
    private EntityAttrValueConverter entityAttrValueConverter;

    @Override
    public QueryResponse<RequestInfoResqDto> selectRequestInfoPage(Integer current, Integer limit,
            RequestInfoQueryReqDto req) {
        req.queryCurrentUserRoles();
        IPage<RequestInfoResqDto> iPage = requestInfoMapper.selectRequestInfo(new Page<>(current, limit), req);
        QueryResponse<RequestInfoResqDto> queryResponse = new QueryResponse<>();
        queryResponse.setPageInfo(new PageInfo(iPage.getTotal(), iPage.getCurrent(), iPage.getSize()));
        queryResponse.setContents(iPage.getRecords());
        return queryResponse;
    }

    @Override
    public RequestInfoResqDto selectDetail(String id) {
        RequestInfo requestInfo = requestInfoMapper.selectOne(new RequestInfo().setId(id).getLambdaQueryWrapper());
        RequestInfoResqDto requestInfoResq = requestInfoConverter.toResp(requestInfo);
        requestInfoResq.setFormItemInfos(formItemInfoService.returnDetail(id));
        return requestInfoResq;
    }

    /**
     * Submit new request
     */
    @Override
    @Transactional
    public RequestInfoResqDto createNewRequestInfo(CreateTaskDto reqDto) {
        RequestInfo requestInfoEntity = buildRequestInfoEntity(reqDto);

        saveOrUpdate(requestInfoEntity);

        reqDto.setId(requestInfoEntity.getId());
        doSaveRequestFormInfo(reqDto);

        // remotely invoke platform service to create new process instance.
        DynamicWorkflowInstInfoDto dynamicWorkflowInstInfoDto = createNewRemoteWorkflowInstance(reqDto);

        if (dynamicWorkflowInstInfoDto == null) {
            log.error("Remotely create new process instance failed due to null response.");
            throw new TaskmanRuntimeException("Remotely create new process instance failed!");
        }

        if (StatusEnum.InProgress.name().equals(dynamicWorkflowInstInfoDto.getStatus())) {
            requestInfoEntity.setProcInstId(dynamicWorkflowInstInfoDto.getProcInstKey());
            requestInfoEntity.setStatus(dynamicWorkflowInstInfoDto.getStatus());
            requestInfoEntity.setUpdatedTime(new Date());
            updateById(requestInfoEntity);
        }
        return requestInfoConverter.toResp(requestInfoEntity);
    }

    private void doSaveRequestFormInfo(CreateTaskDto reqDto) {
        List<FormItemInfo> formItems = new ArrayList<>();
        if (reqDto.getEntities() != null) {
            for (EntityValueDto entityDto : reqDto.getEntities()) {
                List<EntityAttrValueDto> attrValues = entityDto.getAttrValues();
                for (EntityAttrValueDto attrValueDto : attrValues) {
                    FormItemInfo fi = new FormItemInfo();
                    fi.setItemTempId(attrValueDto.getItemTempId());
                    fi.setName(attrValueDto.getName());
                    // TODO

                    formItems.add(fi);
                }
            }
        }

        // reqDto.getEntities().forEach(e -> {
        // items.addAll(formItemInfoConverter.toEntityByAttrValue(e.getAttrValues()));
        // });
        formInfoService.saveFormInfoAndFormItems(formItems, reqDto.getRequestTempId(), reqDto.getId());
    }

    public void saveRequestFormInfo(CreateTaskDto req) {
        List<FormItemInfo> items = new ArrayList<>();
        req.getEntities().stream().forEach(e -> {
            items.addAll(formItemInfoConverter.toEntityByAttrValue(e.getAttrValues()));
        });
        formInfoService.saveFormInfoAndFormItems(items, req.getRequestTempId(), req.getId());
    }

    @Override
    public DynamicWorkflowInstInfoDto createNewRemoteWorkflowInstance(CreateTaskDto reqDto) {
        RequestTemplate requestTemplate = requestTemplateService.getById(reqDto.getRequestTempId());
        String rootEntityId = reqDto.getRootEntity();

        DynamicWorkflowInstCreationInfoDto creationInfoDto = new DynamicWorkflowInstCreationInfoDto();
        creationInfoDto.setProcDefId(requestTemplate.getProcDefId());
        creationInfoDto.setProcDefKey(requestTemplate.getProcDefKey());

        // TODO determine root entity
        DynamicEntityValueDto rootEntityDto = new DynamicEntityValueDto();
        rootEntityDto.setOid(rootEntityId);
        rootEntityDto.setDataId(rootEntityId);

        creationInfoDto.setRootEntityValue(rootEntityDto);

        // TODO determine node bindings
        List<DynamicTaskNodeBindInfoDto> taskNodeBindInfos = buildTaskNodeBindings(reqDto);
        creationInfoDto.setTaskNodeBindInfos(taskNodeBindInfos);

        // TODO
        DynamicWorkflowInstInfoDto dto = platformCoreServiceRestClient.createNewWorkflowInstance(creationInfoDto);
        return dto;
    }

    private List<DynamicTaskNodeBindInfoDto> buildTaskNodeBindings(CreateTaskDto reqDto) {
        List<DynamicTaskNodeBindInfoDto> taskNodeBindInfos = new ArrayList<>();

        List<EntityValueDto> entityValueDtos = reqDto.getEntities();
        if (entityValueDtos == null || entityValueDtos.isEmpty()) {
            return taskNodeBindInfos;
        }

        Map<String, DynamicTaskNodeBindInfoDto> nodeIdAndTaskNodeBindMap = new HashMap<>();
        Map<String, DynamicEntityValueDto> oidAndEntityValueMap = new HashMap<>();
        
        for(EntityValueDto inputEntityDto : entityValueDtos){
            DynamicTaskNodeBindInfoDto taskNodeBindInfoDto = nodeIdAndTaskNodeBindMap.get(inputEntityDto.getNodeId());
            if(taskNodeBindInfoDto == null){
                taskNodeBindInfoDto = new DynamicTaskNodeBindInfoDto();
                taskNodeBindInfoDto.setNodeId(inputEntityDto.getNodeId());
                taskNodeBindInfoDto.setNodeDefId(inputEntityDto.getNodeDefId());
                
                nodeIdAndTaskNodeBindMap.put(inputEntityDto.getNodeId(), taskNodeBindInfoDto);
            }
            
            DynamicEntityValueDto entityValueDto = oidAndEntityValueMap.get(inputEntityDto.getOid());
            if(entityValueDto == null){
                entityValueDto = new DynamicEntityValueDto();
                entityValueDto.setOid(inputEntityDto.getOid());
                entityValueDto.setDataId(inputEntityDto.getDataId());
//                entityValueDto.setEntityDefId(inputEntityDto.get);
                entityValueDto.setEntityName(inputEntityDto.getEntityName());
                entityValueDto.setPackageName(inputEntityDto.getPackageName());
                
                oidAndEntityValueMap.put(inputEntityDto.getOid(), entityValueDto);
            }
            
            if(inputEntityDto.getPreviousOids() != null){
                for(String preOid : inputEntityDto.getPreviousOids()){
                    entityValueDto.addPreviousOid(preOid);
                }
            }
            
            if(inputEntityDto.getSucceedingOids() != null){
                for(String oid : inputEntityDto.getSucceedingOids()){
                    entityValueDto.addSucceedingOid(oid);
                }
            }
            
            if(inputEntityDto.getAttrValues() != null){
                for(EntityAttrValueDto inputAttrValueDto : inputEntityDto.getAttrValues()){
                    DynamicEntityAttrValueDto attrValueDto = new DynamicEntityAttrValueDto();
                    attrValueDto.setAttrDefId(inputAttrValueDto.getAttrDefId());
                    attrValueDto.setAttrName(inputAttrValueDto.getName());
                    attrValueDto.setDataType(inputAttrValueDto.getDataType());
                    attrValueDto.setDataValue(inputAttrValueDto.getDataValue());
                    
                    entityValueDto.addAttrValue(attrValueDto);
                }
            }
            
            taskNodeBindInfoDto.addBoundEntityValue(entityValueDto);
        }
        
        taskNodeBindInfos.addAll(nodeIdAndTaskNodeBindMap.values());

        return taskNodeBindInfos;
    }

    private RequestInfo buildRequestInfoEntity(CreateTaskDto reqDto) {
        RequestInfo requestInfo = new RequestInfo();
        requestInfo.setEmergency(reqDto.getEmergency());
        requestInfo.setDescription(reqDto.getDescription());
        requestInfo.setName(reqDto.getName());
        requestInfo.setRootEntity(reqDto.getRootEntity());
        requestInfo.setId(reqDto.getId());
        requestInfo.setRequestTempId(reqDto.getRequestTempId());

        requestInfo.setCreatedBy(AuthenticationContextHolder.getCurrentUsername());
        requestInfo.setUpdatedBy(AuthenticationContextHolder.getCurrentUsername());
        requestInfo.setReporter(AuthenticationContextHolder.getCurrentUsername());
        requestInfo.setReportTime(new Date());
        requestInfo.setReportRole(AuthenticationContextHolder.getCurrentUserRolesToString());

        return requestInfo;
    }

}

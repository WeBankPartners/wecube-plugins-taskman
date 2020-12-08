package com.webank.taskman.controller.x200;


import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.converter.RequestTemplateGroupConverter;
import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.req.SaveTemplateGropReq;
import com.webank.taskman.service.RequestTemplateGroupService;
import io.swagger.annotations.ApiOperation;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;

@RequestMapping("/v2/request")
public class V2RequestController {

    @Autowired
    RequestTemplateGroupService requestTemplateGroupService;

    @Autowired
    RequestTemplateGroupConverter requestTemplateGroupConverter;
    
    @PostMapping("/template/group/save")
    @ApiOperationSupport(order = 21)
    @ApiOperation(value = "save RequestTemplateGroup", notes = "")
    public JsonResponse createTemplateGroup(
            @RequestBody SaveTemplateGropReq req) throws Exception {
        if(StringUtils.isEmpty(req.getName())){
            return  JsonResponse.error(" manageRoleId is null");
        }
        if(StringUtils.isEmpty(req.getName())){
            return  JsonResponse.error(" name is null");
        }
        //requestTemplateGroupService.addTemplateGroup(requestTemplateGroupConverter.addReqToDomain(req));
        return JsonResponse.okay();
    }


}

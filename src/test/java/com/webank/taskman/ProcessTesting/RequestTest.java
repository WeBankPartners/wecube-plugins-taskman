package com.webank.taskman.ProcessTesting;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.webank.taskman.TmallApplicationTests;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.constant.RoleTypeEnum;
import com.webank.taskman.constant.StatusEnum;
import com.webank.taskman.converter.RequestTemplateConverter;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.RequestTemplateDTO;
import com.webank.taskman.dto.req.QueryRequestInfoReq;
import com.webank.taskman.dto.req.QueryRoleRelationBaseReq;
import com.webank.taskman.dto.resp.RequestInfoResq;
import com.webank.taskman.service.RequestInfoService;
import com.webank.taskman.service.RequestTemplateService;
import org.junit.Test;
import org.springframework.beans.factory.annotation.Autowired;

import java.util.List;

public class RequestTest extends TmallApplicationTests {

    @Autowired
    RequestInfoService requestInfoService;

    @Autowired
    RequestTemplateConverter requestTemplateConverter;

    @Autowired
    RequestTemplateService requestTemplateService;

    @Test
    public void selectRequestTest() {
        Integer page = 1;
        Integer pageSize = 10;
        QueryRequestInfoReq req = new QueryRequestInfoReq();
        req.setName("");
        req.setStatus("");
        req.setEmergency("");
        req.setReporter("");
        QueryResponse<RequestInfoResq> list = requestInfoService.selectRequestInfoService(page, pageSize, req);
        for (RequestInfoResq content : list.getContents()) {
            System.out.println(content.toString());
        }
    }

    @Test
    public void RequestToReportTest() {
        //select template
        RequestTemplate query = new RequestTemplate().setStatus(StatusEnum.RELEASED.toString());
        QueryWrapper<RequestTemplate> queryWrapper = new QueryWrapper<>();
        queryWrapper.eq("status",StatusEnum.RELEASED.toString())
                .inSql("id",String.format(QueryRoleRelationBaseReq.QUERY_BY_ROLE_SQL,
                        RoleTypeEnum.USE_ROLE.getType(),
                        AuthenticationContextHolder.getCurrentUserRolesToString()));
        List<RequestTemplateDTO> list=requestTemplateConverter.toDto(requestTemplateService.list(queryWrapper));
        for (RequestTemplateDTO requestTemplateDTO : list) {
            System.out.println(requestTemplateDTO.toString());
        }

    }
}

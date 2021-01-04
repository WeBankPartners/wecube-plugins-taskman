package com.webank.taskman.QueryList;

import com.webank.taskman.TmallApplicationTests;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.dto.req.QueryRequestInfoReq;
import com.webank.taskman.dto.resp.RequestInfoResq;
import com.webank.taskman.service.RequestInfoService;
import com.webank.taskman.utils.GsonUtil;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;

public class RequestInfoListTest extends TmallApplicationTests {
    private static final Logger log = LoggerFactory.getLogger(RequestInfoListTest.class);
    @Autowired
    RequestInfoService requestInfoService;

    @Test
    public void RequestInfoListTest(){
        Integer page=1;
        Integer pageSize=2;
        QueryRequestInfoReq req=new QueryRequestInfoReq();
        log.info("Received request parameters:{}", GsonUtil.GsonString(req) );
        QueryResponse<RequestInfoResq> list = requestInfoService.selectRequestInfoPage(page, pageSize,req);
        for (RequestInfoResq content : list.getContents()) {
        }
    }
}

package com.webank.taskman.QueryList;

import org.junit.Test;
import org.springframework.beans.factory.annotation.Autowired;

import com.webank.taskman.TmallApplicationTests;
import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.converter.TaskInfoConverter;
import com.webank.taskman.dto.TaskInfoDTO;
import com.webank.taskman.dto.req.QueryTaskInfoReq;
import com.webank.taskman.service.TaskInfoService;

public class TaskInfoListTest extends TmallApplicationTests {

    @Autowired
    TaskInfoService taskInfoService;

    @Autowired
    TaskInfoConverter taskInfoConverter;

    @Test
    public void taskInfoListTest(){
        Integer page=1;
        Integer pageSize=2;
        QueryTaskInfoReq req=new QueryTaskInfoReq();
        QueryResponse<TaskInfoDTO> queryResponse = taskInfoService.selectTaskInfo(page, pageSize,req);
        JsonResponse jsonResponse=new JsonResponse();
        jsonResponse.setData(queryResponse);
    }
}

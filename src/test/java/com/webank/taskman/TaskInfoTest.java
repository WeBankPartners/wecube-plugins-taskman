package com.webank.taskman;

import com.webank.taskman.base.QueryResponse;
import com.webank.taskman.dto.TaskInfoDTO;
import com.webank.taskman.dto.req.QueryTaskInfoReq;
import com.webank.taskman.service.TaskInfoService;
import com.webank.taskman.utils.GsonUtil;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;

public class TaskInfoTest extends TmallApplicationTests {

    private static final Logger log = LoggerFactory.getLogger(TaskInfoTest.class);
    @Autowired
    TaskInfoService taskInfoService;

    @Test
    public void  selectTaskInfo(){
        int page = 1;
        int pageSize = 10;
        QueryTaskInfoReq req = new QueryTaskInfoReq();
        QueryResponse<TaskInfoDTO> response = taskInfoService.selectTaskInfo(page,pageSize,req);
        log.info("query result:{}", GsonUtil.GsonString(response));
    }
}

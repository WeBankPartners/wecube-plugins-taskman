package com.webank.taskman.controller;

import static com.webank.taskman.base.JsonResponse.okayWithData;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.base.LocalPageableQueryResult;
import com.webank.taskman.dto.CreateTaskDto;
import com.webank.taskman.dto.req.RequestInfoQueryReqDto;
import com.webank.taskman.dto.resp.RequestInfoQueryResultDto;
import com.webank.taskman.service.RequestInfoService;

@RestController
@RequestMapping("/v1/request")
public class RequestInfoController {
    @Autowired
    private RequestInfoService requestInfoService;
    
    /**
     * Submit new request
     * 
     * @param req
     * @return
     */
    @PostMapping("/save")
    public JsonResponse createNewRequestInfo(@RequestBody CreateTaskDto req) {
        return okayWithData(requestInfoService.createNewRequestInfo(req));
    }

    @PostMapping("/search/{page}/{page-size}")
    public JsonResponse searchRequestInfos(@PathVariable("page") Integer page,
            @PathVariable("page-size") Integer pageSize, @RequestBody(required = false) RequestInfoQueryReqDto req) {
        LocalPageableQueryResult<RequestInfoQueryResultDto> list = requestInfoService.searchRequestInfos(page, pageSize, req);
        return okayWithData(list);
    }

    @GetMapping("/details/{id}")
    public JsonResponse fetchRequestInfoDetail(@PathVariable("id") String id) {
        RequestInfoQueryResultDto requestInfoResq = requestInfoService.fetchRequestInfoDetail(id);
        return okayWithData(requestInfoResq);
    }
}

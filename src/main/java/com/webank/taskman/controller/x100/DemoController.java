package com.webank.taskman.controller.x100;

import com.github.xiaoymin.knife4j.annotations.ApiOperationSupport;
import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.resp.SynthesisRequestTempleResp;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.HashSet;
import java.util.Set;

@Api(tags = {"6、 RequestORTask Synthesis API"})
@RestController
@RequestMapping("/v1/demo")
public class DemoController {


    private static final Logger log = LoggerFactory.getLogger(DemoController.class);

    @ApiOperationSupport(order = 1)
    @GetMapping("/set")
    @ApiOperation(value = "request-Synthesis-search")
    public JsonResponse<Set<SynthesisRequestTempleResp>> selectRequestSynthesis(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize") @PathVariable("pageSize") Integer pageSize)
            throws Exception {
        Set<SynthesisRequestTempleResp> set = new HashSet<>();
        set.add(new SynthesisRequestTempleResp());
        return JsonResponse.okayWithData(set);
    }

}

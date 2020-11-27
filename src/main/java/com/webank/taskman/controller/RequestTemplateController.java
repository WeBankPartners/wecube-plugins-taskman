package com.webank.taskman.controller;


import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.RequestTemplateVO;
import com.webank.taskman.service.RequestTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

/**
 * <p>
 * 请求模板信息表  前端控制器
 * </p>
 *
 * @author ${author}
 * @since 2020-11-26
 */
@RestController
@RequestMapping("/request-template")
@Api(tags = {"RequestTemplate接口使用"})
public class RequestTemplateController {
    @Autowired
    RequestTemplateService requestTemplateService;

    @PostMapping("/addRequestTemplate")
    @ApiOperation(value = "增加请求模板",notes = "需要传入RequestTemplateVO对象")
    public JsonResponse addRequestTemplate(@RequestBody RequestTemplateVO requestTemplateVO)throws Exception{
        requestTemplateService.AddRequestTemplate(requestTemplateVO);
        return JsonResponse.okay();
    }

    @PostMapping("/updateRequestTemplate")
    @ApiOperation(value = "修改请求模板",notes = "需要传入RequestTemplateVO对象")
    public JsonResponse updateRequestTemplate(@RequestBody RequestTemplateVO requestTemplateVO)throws Exception{
        requestTemplateService.updateRequestTemplate(requestTemplateVO);
        return JsonResponse.okay();
    }

    @GetMapping("/deleteRequestTemplate/{id}")
    @ApiOperation(value = "删除请求模板",notes = "需要传入RequestTemplateVO对象")
    public JsonResponse deleteRequestTemplate(@PathVariable("id") String id)throws Exception{
        requestTemplateService.deleteRequestTemplate(id);
        return JsonResponse.okay();
    }

    @PostMapping("/selectRequestTemplate")
    @ApiOperation(value = "查询请求模板",notes = "需要传入RequestTemplateVO对象")
    public JsonResponse selectRequestTemplate(@RequestBody RequestTemplateVO requestTemplateVO){
        return JsonResponse.okayWithData(requestTemplateService.selectRequestTemplate(requestTemplateVO));
    }

}


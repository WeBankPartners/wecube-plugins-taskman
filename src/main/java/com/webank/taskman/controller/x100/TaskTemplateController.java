package com.webank.taskman.controller.x100;


import com.webank.taskman.dto.*;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;
import com.webank.taskman.dto.resp.TaskTemplateResp;
import com.webank.taskman.service.TaskTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.BindingResult;
import org.springframework.validation.ObjectError;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;


@RestController
@RequestMapping("/v1/task/template")
@Api(tags = {"5、  TaskTemplate inteface API"})
public class TaskTemplateController {

    @Autowired
    private TaskTemplateService taskTemplateService;

    //TODO implemented   insert or update
    @PostMapping("/save")
    @ApiOperation(value = "saveTaskTemplate", notes = "Need to pass in object: ")
    public JsonResponse createTaskTemplate(@Valid @RequestBody SaveTaskTemplateReq taskTemplateReq, BindingResult bindingResult) throws Exception {
        if (bindingResult.hasErrors()) {
            for (ObjectError error : bindingResult.getAllErrors()) {
                return JsonResponse.okayWithData(error.getDefaultMessage());
            }
        }
        TaskTemplateResp taskTemplateResp = taskTemplateService.saveTaskTemplateByReq(taskTemplateReq);
        return JsonResponse.okayWithData(taskTemplateResp);
    }

}


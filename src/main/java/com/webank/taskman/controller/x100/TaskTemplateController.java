package com.webank.taskman.controller.x100;


import com.webank.taskman.domain.TaskTemplate;
import com.webank.taskman.dto.*;
import com.webank.taskman.dto.req.SaveTaskTemplateReq;
import com.webank.taskman.dto.req.SelectTaskTemplateRep;
import com.webank.taskman.dto.resp.TaskTemplateResp;
import com.webank.taskman.dto.resp.TaskTemplateSVResp;
import com.webank.taskman.service.TaskTemplateService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.annotations.ApiParam;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.BindingResult;
import org.springframework.validation.ObjectError;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;
import java.util.List;


@RestController
@RequestMapping("/v1/task/template")
@Api(tags = {"5、  TaskTemplate inteface API"})
public class TaskTemplateController {

    @Autowired
    private TaskTemplateService taskTemplateService;

    //TODO implemented   insert or update
    @PostMapping("/save")
    @ApiOperation(value = "add OR update TaskTemplate", notes = "Need to pass in object: ")
    public JsonResponse createTaskTemplate(@Valid @RequestBody SaveTaskTemplateReq taskTemplateReq, BindingResult bindingResult) throws Exception {
        if (bindingResult.hasErrors()) {
            for (ObjectError error : bindingResult.getAllErrors()) {
                return JsonResponse.okayWithData(error.getDefaultMessage());
            }
        }
        TaskTemplate taskTemplate = taskTemplateService.saveTaskTemplate(taskTemplateReq);
        TaskTemplateSVResp taskTemplateSVResp = new TaskTemplateSVResp();
        taskTemplateSVResp.setId(taskTemplate.getId());
        return JsonResponse.okayWithData(taskTemplateSVResp);
    }

    @DeleteMapping("/delete/{id}")
    @ApiOperation(value = "delete TaskTemplate", notes = "需要传入id")
    public JsonResponse deleteTaskTemplateByID(@PathVariable("id") String id) throws Exception {
        taskTemplateService.deleteTaskTemplateByIDService(id);
        return JsonResponse.okay();
    }

    @PostMapping("/search/{page}/{pageSize}")
    @ApiOperation(value = "search TaskTemplatePage ")
    public JsonResponse<QueryResponse<TaskTemplateResp>> selectTaskTemplatePage(
            @ApiParam(name = "page") @PathVariable("page") Integer page,
            @ApiParam(name = "pageSize") @PathVariable("pageSize") Integer pageSize,
            @RequestBody(required = false) SelectTaskTemplateRep rep
    ) throws Exception {
        QueryResponse<TaskTemplateResp> queryResponse = taskTemplateService.selectTaskTemplate(page, pageSize, rep);

        return JsonResponse.okayWithData(queryResponse);
    }

    @GetMapping("/search")
    @ApiOperation(value = "search TaskTemplateAll ")
    public JsonResponse selectTaskTemplateAll() throws Exception {
        List<TaskTemplateResp> taskTemplateRespList = taskTemplateService.selectTaskTemplateAll();
        return JsonResponse.okayWithData(taskTemplateRespList);
    }


    @GetMapping("/detail/{id}")
    @ApiOperation(value = "detail TaskTemplate", notes = "需要传入id")
    public JsonResponse detail(@PathVariable("id") String id) throws Exception {
        TaskTemplateResp templateResp= taskTemplateService.selectTaskTemplateOne(id);
        return JsonResponse.okayWithData(templateResp);
    }
}


package com.webank.taskman.controller;



import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.TemplateGroupVO;
import com.webank.taskman.service.TemplateGroupService;
import io.swagger.annotations.ApiImplicitParam;
import io.swagger.annotations.ApiImplicitParams;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
@RestController
@RequestMapping("/template-group")
public class TemplateGroupController {
    @Autowired
    TemplateGroupService templateGroupService;

    @PostMapping("/save")
    public JsonResponse createTemplateGroup(
            @RequestBody TemplateGroupVO templateGroupVO) throws Exception {
        templateGroupService.createTemplateGroupService(templateGroupVO);
        return JsonResponse.okay();
    }

    @ApiImplicitParams({
            @ApiImplicitParam(name = "templateGroupVO",value ="创建模板请求对象",readOnly = true )
    })

    @PostMapping("edit")
    public JsonResponse updateTemplateGroup(
            @RequestBody TemplateGroupVO templateGroupVO) throws Exception {
        templateGroupService.updateTemplateGroupService(templateGroupVO);
        return JsonResponse.okay();
    }

    @GetMapping("/selectAll")
    public JsonResponse selectAllTemplateGroup() throws Exception {
        return JsonResponse.okayWithData(templateGroupService.selectAllTemplateGroupService());
    }

    @GetMapping("/delete/{id}")
    public JsonResponse deleteTemplateGroupByID(@PathVariable("id") String id) throws Exception {
        templateGroupService.deleteTemplateGroupByIDService(id);
        return JsonResponse.okay();
    }
}


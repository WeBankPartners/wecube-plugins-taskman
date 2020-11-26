package com.webank.taskman.controller;



import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.TemplateGroupVO;
import com.webank.taskman.service.TemplateGroupService;
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


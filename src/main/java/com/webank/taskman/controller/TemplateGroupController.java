package com.webank.taskman.controller;



import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.dto.TemplateGroupVO;
import com.webank.taskman.service.TemplateGroupService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/template-group")
@Api(tags = {"TemplateGroup接口使用"})
public class TemplateGroupController {
    @Autowired
    TemplateGroupService templateGroupService;

    @PostMapping("/save")
    @ApiOperation(value = "增加模板组",notes = "需要传入templateGroupVO对象")
    public JsonResponse createTemplateGroup(
            @RequestBody TemplateGroupVO templateGroupVO) throws Exception {
        templateGroupService.createTemplateGroupService(templateGroupVO);
        return JsonResponse.okay();
    }

    @PostMapping("edit")
    @ApiOperation(value = "修改模板组",notes = "需要传入templateGroupVO对象")
    public JsonResponse updateTemplateGroup(
            @RequestBody TemplateGroupVO templateGroupVO) throws Exception {
        templateGroupService.updateTemplateGroupService(templateGroupVO);
        return JsonResponse.okay();
    }

    @GetMapping("/selectAll")
    @ApiOperation(value = "获取模板组数据",notes = "无需传入值")
    public JsonResponse selectAllTemplateGroup() throws Exception {
        return JsonResponse.okayWithData(templateGroupService.selectAllTemplateGroupService());
    }

    @GetMapping("/delete/{id}")
    @ApiOperation(value = "删除模板组",notes = "需要传入id")
    public JsonResponse deleteTemplateGroupByID(@PathVariable("id") String id) throws Exception {
        System.out.println("-------------------"+id);
        templateGroupService.deleteTemplateGroupByIDService(id);
        return JsonResponse.okay();
    }
}


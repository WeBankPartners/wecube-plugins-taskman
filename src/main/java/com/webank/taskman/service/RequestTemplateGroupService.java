package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.RequestTemplateGroup;
import com.webank.taskman.dto.TemplateGroupDTO;
import com.webank.taskman.dto.TemplateGroupVO;

import java.util.List;


public interface RequestTemplateGroupService extends IService<RequestTemplateGroup> {

    void createTemplateGroupService(TemplateGroupVO templateGroupVO) throws Exception;

    void updateTemplateGroupService(TemplateGroupVO templateGroupVO) throws Exception;

    List<TemplateGroupDTO> selectAllTemplateGroupService() throws Exception;

    void deleteTemplateGroupByIDService(String id) throws Exception;
}

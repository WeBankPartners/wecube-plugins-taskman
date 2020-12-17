package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.TemplateGroup;
import com.webank.taskman.dto.TemplateGroupDTO;
import com.webank.taskman.dto.TemplateGroupVO;

import java.util.List;

/**
 * <p>
 * 模板组信息表  服务类
 * </p>
 *
 * @author ${author}
 * @since 2020-11-26
 */
public interface TemplateGroupService extends IService<TemplateGroup> {

    void createTemplateGroupService(TemplateGroupVO templateGroupVO) throws Exception;

    void updateTemplateGroupService(TemplateGroupVO templateGroupVO) throws Exception;

    List<TemplateGroupDTO> selectAllTemplateGroupService() throws Exception;
}

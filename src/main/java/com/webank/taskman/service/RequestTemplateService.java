package com.webank.taskman.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.RequestTemplateDTO;
import com.webank.taskman.dto.RequestTemplateVO;

import java.util.List;
/**
 * <p>
 * 请求模板信息表  服务类
 * </p>
 *
 * @author ${author}
 * @since 2020-11-26
 */
public interface RequestTemplateService extends IService<RequestTemplate> {

    void AddRequestTemplate(RequestTemplateVO requestTemplateVO) throws Exception;

    void updateRequestTemplate(RequestTemplateVO requestTemplateVO) throws Exception;

    void deleteRequestTemplate(String id) throws Exception;

    List<RequestTemplateDTO> selectRequestTemplate(RequestTemplateVO requestTemplateVO);

}

package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.conditions.update.LambdaUpdateWrapper;
import com.webank.taskman.converter.RequestTemplateConverter;
import com.webank.taskman.domain.RequestTemplate;
import com.webank.taskman.dto.RequestTemplateDTO;
import com.webank.taskman.dto.RequestTemplateVO;
import com.webank.taskman.mapper.RequestTemplateMapper;
import com.webank.taskman.service.RequestTemplateService;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.List;
/**
 * <p>
 * 请求模板信息表  服务实现类
 * </p>
 *
 * @author ${author}
 * @since 2020-11-26
 */
@Service
public class RequestTemplateServiceImpl extends ServiceImpl<RequestTemplateMapper, RequestTemplate> implements RequestTemplateService {
    @Autowired
    RequestTemplateMapper requestTemplateMapper;

    @Autowired
    RequestTemplateConverter requestTemplateConverter;

    @Override
    public void AddRequestTemplate(RequestTemplateVO requestTemplateVO)  throws Exception {
        if (requestTemplateVO == null) {
            throw new Exception("Template group objects cannot be empty");
        }
        RequestTemplate requestTemplate=requestTemplateConverter.voToDomain(requestTemplateVO);
        requestTemplateMapper.insert(requestTemplate);
    }

    @Override
    public void updateRequestTemplate(RequestTemplateVO requestTemplateVO) throws Exception {
        if (requestTemplateVO == null) {
            throw new Exception("Template group objects cannot be empty");
        }
        RequestTemplate requestTemplate=requestTemplateConverter.voToDomain(requestTemplateVO);
        SimpleDateFormat df = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        requestTemplate.setUpdatedTime(df.parse(df.format(new Date())));
        requestTemplateMapper.updateById(requestTemplate);
    }

    @Override
    public void deleteRequestTemplate(String id) throws Exception {
        LambdaUpdateWrapper<RequestTemplate> lambdaUpdateWrapper = new LambdaUpdateWrapper<>();
        lambdaUpdateWrapper.eq(RequestTemplate::getId, id).set(RequestTemplate::getDelFlag, 1);
        int update = requestTemplateMapper.update(null, lambdaUpdateWrapper);
        if (update != 1) {
            throw new Exception("Template group deletion failed");
        }
    }

    @Override
    public List<RequestTemplateDTO> selectRequestTemplate(RequestTemplateVO requestTemplateVO){
        QueryWrapper<RequestTemplate> queryWrapper = new QueryWrapper<>();
        queryWrapper.like(requestTemplateVO.getId() != null, "id", requestTemplateVO.getId());
        List<RequestTemplate> list=requestTemplateMapper.selectList(queryWrapper);
        return requestTemplateConverter.toDto(list);
    }
}

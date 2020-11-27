package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.domain.RequestInfo;
import com.webank.taskman.mapper.RequestInfoMapper;
import com.webank.taskman.service.RequestInfoService;
import org.springframework.stereotype.Service;


@Service
public class RequestInfoServiceImpl extends ServiceImpl<RequestInfoMapper, RequestInfo> implements RequestInfoService {

}

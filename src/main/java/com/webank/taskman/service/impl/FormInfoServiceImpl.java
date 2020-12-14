package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.domain.FormInfo;
import com.webank.taskman.dto.req.SaveFormInfoReq;
import com.webank.taskman.mapper.FormInfoMapper;
import org.springframework.stereotype.Service;


@Service
public class FormInfoServiceImpl extends ServiceImpl<FormInfoMapper, FormInfo> implements com.webank.taskman.service.FormInfoService {

    @Override
    public void saveFormInfoByReq(SaveFormInfoReq req) {

    }
}

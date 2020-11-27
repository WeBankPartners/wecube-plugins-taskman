package com.webank.taskman.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.webank.taskman.domain.RoleInfo;
import com.webank.taskman.mapper.RoleInfoMapper;
import com.webank.taskman.service.RoleInfoService;
import org.springframework.stereotype.Service;


@Service
public class RoleInfoServiceImpl extends ServiceImpl<RoleInfoMapper, RoleInfo> implements RoleInfoService {

}

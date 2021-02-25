package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.webank.taskman.domain.RequestInfo;
import com.webank.taskman.dto.req.RequestInfoQueryReqDto;
import com.webank.taskman.dto.resp.RequestInfoResqDto;
import org.apache.ibatis.annotations.Param;


public interface RequestInfoMapper extends BaseMapper<RequestInfo> {

    IPage<RequestInfoResqDto> selectRequestInfo(Page page, @Param("param") RequestInfoQueryReqDto req);

}

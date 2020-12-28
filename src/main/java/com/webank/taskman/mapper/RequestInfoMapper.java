package com.webank.taskman.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.baomidou.mybatisplus.core.metadata.IPage;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.webank.taskman.domain.RequestInfo;
import com.webank.taskman.dto.RequestInfoDTO;
import com.webank.taskman.dto.req.QueryRequestInfoReq;
import com.webank.taskman.dto.req.SaveRequestInfoReq;
import com.webank.taskman.dto.req.SynthesisRequestInfoReq;
import com.webank.taskman.dto.resp.RequestInfoResq;
import org.apache.ibatis.annotations.Param;


public interface RequestInfoMapper extends BaseMapper<RequestInfo> {

    IPage<RequestInfoResq> selectRequestInfo(Page page, @Param("param") QueryRequestInfoReq req);

    IPage<RequestInfo> selectSynthesisRequestInfo(Page page, @Param("param") SynthesisRequestInfoReq req);
}

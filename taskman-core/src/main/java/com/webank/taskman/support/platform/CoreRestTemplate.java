package com.webank.taskman.support.platform;

import java.util.HashMap;
import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestTemplate;

import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.commons.AppProperties.ServiceTaskmanProperties;
import com.webank.taskman.support.platform.dto.PlatformResponseDto;
import com.webank.taskman.support.platform.dto.PlatformResponseDto.DefaultCoreResponse;
import com.webank.taskman.utils.GsonUtil;
import com.webank.taskman.utils.JsonUtils;

@Component
public class CoreRestTemplate {
    private static final Logger log = LoggerFactory.getLogger(CoreRestTemplate.class);

    @Autowired
    private RestTemplate restTemplate;

    @Autowired
    private ServiceTaskmanProperties taskmanProperties;
    
    //TODO to refactor

    public String get(String targetUrl) throws PlatformRemoteCallException {
        log.info("{} About to call {} ", taskmanProperties.getVersion(), targetUrl);
        try {
            DefaultCoreResponse jsonResponse = restTemplate.getForObject(targetUrl, DefaultCoreResponse.class);
            log.info("Core response: {}", jsonResponse);
            validateJsonResponse(jsonResponse, true);
            return GsonUtil.GsonString(jsonResponse.getData());
        } catch (Exception e) {
            throw e;
        }
    }

    @SuppressWarnings("unchecked")
    public <D, R extends PlatformResponseDto<?>> D get(String targetUrl, Class<R> responseType)
            throws PlatformRemoteCallException {
        log.info("{} About to call {} ", taskmanProperties.getVersion(), targetUrl);
        try {
            R jsonResponse = restTemplate.getForObject(targetUrl, responseType);
            log.info("Core response: {}", jsonResponse);
            validateJsonResponse(jsonResponse);
            return (D) jsonResponse.getData();
        } catch (Exception e) {
            throw e;
        }
    }

    @SuppressWarnings("unchecked")
    public <D, R extends PlatformResponseDto<?>> D get(String targetUrl, Class<R> responseType, String paramJsonStr)
            throws PlatformRemoteCallException {
        log.info("{} About to call {} ", taskmanProperties.getVersion(), targetUrl);
        Object uriVariable = paramJsonStr;
        try {
            Map<String, Object> map = new HashMap<>();
            uriVariable = JsonUtils.toObject(paramJsonStr, map.getClass());
        } catch (Exception e) {
            log.error("paramJsonStr is not json: {} ", targetUrl);
        }
        R jsonResponse = restTemplate.getForObject(targetUrl, responseType, uriVariable);
        log.info("Core response: {} ", jsonResponse);
        validateJsonResponse(jsonResponse);
        return (D) jsonResponse.getData();
    }

    @SuppressWarnings("unchecked")
    public <D, R extends PlatformResponseDto<?>> D postForResponse(String targetUrl, Object postObject, Class<R> responseType)
            throws PlatformRemoteCallException {
        log.info("{}About to POST {} with postObject {}", taskmanProperties.getVersion(), targetUrl,
                postObject.toString());
        R jsonResponse = restTemplate.postForObject(targetUrl, postObject, responseType);
        log.info("Core response: {} ", GsonUtil.GsonString(jsonResponse));
        validateJsonResponse(jsonResponse, false);
        return (D) jsonResponse.getData();
    }

    private void validateJsonResponse(PlatformResponseDto<?> jsonResponse) throws PlatformRemoteCallException {
        validateJsonResponse(jsonResponse, true);
    }

    private void validateJsonResponse(PlatformResponseDto<?> jsonResponse, boolean dataRequired)
            throws PlatformRemoteCallException {
        if (null == jsonResponse) {
            throw new PlatformRemoteCallException("Call WeCube-Core failed due to no response.");
        }
        if (!JsonResponse.STATUS_OK.equalsIgnoreCase(jsonResponse.getStatus())) {
            throw new PlatformRemoteCallException("Core Error: " + jsonResponse.getMessage(), jsonResponse);
        }
        if (dataRequired && null == jsonResponse.getData()) {
            throw new PlatformRemoteCallException("Call WeCube-Core failed due to unexpected empty response.",
                    jsonResponse);
        }
    }

    public RestTemplate getRestTemplate() {
        return restTemplate;
    }

}

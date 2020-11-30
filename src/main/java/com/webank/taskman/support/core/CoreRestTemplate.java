package com.webank.taskman.support.core;

import com.webank.taskman.dto.JsonResponse;
import com.webank.taskman.support.core.dto.CoreResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestTemplate;

@Component
public class CoreRestTemplate {

    private static final Logger log = LoggerFactory.getLogger(CoreRestTemplate.class);

    @Autowired
    private RestTemplate restTemplate;

    @SuppressWarnings("unchecked")
    public <D, R extends CoreResponse> D get(String targetUrl, Class<R> responseType) throws CoreRemoteCallException {
        log.info("About to call {} ", targetUrl);
        R jsonResponse = restTemplate.getForObject(targetUrl, responseType);
        log.info("Core response: {} ", jsonResponse);
        validateJsonResponse(jsonResponse);
        return (D) jsonResponse.getData();
    }

    public <D, R extends CoreResponse> D postForResponse(String targetUrl, Class<R> responseType)
            throws CoreRemoteCallException {
        return postForResponse(targetUrl, null, responseType);
    }

    @SuppressWarnings("unchecked")
    public <D, R extends CoreResponse> D postForResponse(String targetUrl, Object postObject, Class<R> responseType)
            throws CoreRemoteCallException {
        log.info("About to POST {} with postObject {}", targetUrl, postObject.toString());
        R jsonResponse = restTemplate.postForObject(targetUrl, postObject, responseType);
        log.info("Core response: {} ", jsonResponse.toString());
        validateJsonResponse(jsonResponse, false);
        return (D) jsonResponse.getData();
    }

    private void validateJsonResponse(CoreResponse jsonResponse) throws CoreRemoteCallException {
        validateJsonResponse(jsonResponse, true);
    }

    private void validateJsonResponse(CoreResponse jsonResponse, boolean dataRequired) throws CoreRemoteCallException {
        if (jsonResponse == null) {
            throw new CoreRemoteCallException("Call WeCube-Core failed due to no response.");
        }
        if (!JsonResponse.STATUS_OK.equalsIgnoreCase(jsonResponse.getStatus())) {
            throw new CoreRemoteCallException("Core Error: " + jsonResponse.getMessage(), jsonResponse);
        }
        if (dataRequired && jsonResponse.getData() == null) {
            throw new CoreRemoteCallException("Call WeCube-Core failed due to unexpected empty response.",
                    jsonResponse);
        }
    }

}

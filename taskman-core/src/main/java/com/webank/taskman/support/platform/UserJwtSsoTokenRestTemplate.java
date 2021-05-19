package com.webank.taskman.support.platform;

import java.io.IOException;
import java.net.URI;

import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.client.ClientHttpRequest;
import org.springframework.web.client.RestTemplate;

import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.commons.AuthenticationContextHolder.AuthenticatedUser;

/**
 * 
 * @author gavin
 *
 */
public class UserJwtSsoTokenRestTemplate extends RestTemplate{

    private static final Logger log = LoggerFactory.getLogger(UserJwtSsoTokenRestTemplate.class);

    @Override
    protected ClientHttpRequest createRequest(URI url, HttpMethod method) throws IOException {
        
        ClientHttpRequest req = super.createRequest(url, method);
        
        AuthenticatedUser loginUser = AuthenticationContextHolder.getCurrentUser();
        if(loginUser != null && StringUtils.isNotBlank(loginUser.getToken())){
            if (log.isDebugEnabled()) {
                log.debug("request {} with access token:{}", url.toString(), loginUser.getToken());
            }

            req.getHeaders().add(HttpHeaders.AUTHORIZATION,
                    loginUser.getToken());
        }
        return req;
    }
    
    
    
}

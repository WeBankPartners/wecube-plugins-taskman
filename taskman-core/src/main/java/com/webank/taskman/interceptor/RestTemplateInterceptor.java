package com.webank.taskman.interceptor;

import java.io.IOException;

import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpRequest;
import org.springframework.http.client.ClientHttpRequestExecution;
import org.springframework.http.client.ClientHttpRequestInterceptor;
import org.springframework.http.client.ClientHttpResponse;
import org.springframework.stereotype.Component;

import com.webank.taskman.commons.AuthenticationContextHolder;

@Component
public class RestTemplateInterceptor implements ClientHttpRequestInterceptor {

    @Override
    public ClientHttpResponse intercept(HttpRequest request, byte[] body, ClientHttpRequestExecution execution)
            throws IOException {
        HttpHeaders headers = request.getHeaders();
        if (AuthenticationContextHolder.getCurrentUser() != null
                && !AuthenticationContextHolder.getCurrentUser().getToken().isEmpty()) {
            headers.add("Authorization", AuthenticationContextHolder.getCurrentUser().getToken());
        }
        return execution.execute(request, body);
    }
}

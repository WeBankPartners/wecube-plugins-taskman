package com.webank.taskman.interceptor;

import com.webank.taskman.commons.AuthenticationContextHolder;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.security.Principal;

@Component
public class AuthenticationRequestContextInterceptor implements HandlerInterceptor {
    private static final String AUTHORIZATION = "Authorization";
    public static final String REQ_ATTR_KEY_CURRENT_USER = "REQ_ATTR_KEY_CURRENT_USER";

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) {
        Principal userPrincipal = request.getUserPrincipal();
        if (userPrincipal != null) {
            AuthenticationContextHolder.AuthenticatedUser currentUser = new AuthenticationContextHolder.AuthenticatedUser(userPrincipal.getName());

            currentUser.withAuthorities(extractAuthorities(userPrincipal));
            currentUser.setToken(request.getHeader(AUTHORIZATION));
            AuthenticationContextHolder.setAuthenticatedUser(currentUser);

            request.setAttribute(REQ_ATTR_KEY_CURRENT_USER, currentUser);
        }

        return true;
    }

    private String[] extractAuthorities(Principal userPrincipal) {
        String[] authorities = new String[0];
        if (userPrincipal instanceof UsernamePasswordAuthenticationToken) {
            authorities = ((UsernamePasswordAuthenticationToken) userPrincipal).getAuthorities().stream()
                    .map(GrantedAuthority::toString).toArray(String[]::new);
        }
        return authorities;
    }

    @Override
    public void afterCompletion(HttpServletRequest request, HttpServletResponse response, Object handler, Exception ex)
            throws Exception {
        AuthenticationContextHolder.clearCurrentUser();
        request.removeAttribute(REQ_ATTR_KEY_CURRENT_USER);
    }
}

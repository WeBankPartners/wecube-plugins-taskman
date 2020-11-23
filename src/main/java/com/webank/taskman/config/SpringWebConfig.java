package com.webank.taskman.config;

import com.webank.taskman.authclient.filter.Http401AuthenticationEntryPoint;
import com.webank.taskman.authclient.filter.JwtClientConfig;
import com.webank.taskman.authclient.filter.JwtSsoBasedAuthenticationFilter;
import com.webank.taskman.commons.AppProperties;
import com.webank.taskman.interceptor.AuthenticationRequestContextInterceptor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.method.configuration.EnableGlobalMethodSecurity;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configuration.WebSecurityConfigurerAdapter;
import org.springframework.web.servlet.config.annotation.EnableWebMvc;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.ResourceHandlerRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;
import springfox.documentation.swagger2.annotations.EnableSwagger2;

import javax.servlet.Filter;

@Configuration
@EnableWebMvc
@EnableSwagger2
@EnableWebSecurity
@EnableGlobalMethodSecurity(jsr250Enabled = true, prePostEnabled = true, securedEnabled = true)
public class SpringWebConfig extends WebSecurityConfigurerAdapter implements WebMvcConfigurer {

    @Autowired
    private AuthenticationRequestContextInterceptor authenticationRequestContextInterceptor;
    @Autowired
    private AppProperties.ServiceManagementProperties serviceManagementProperties;

    @Override
    public void addResourceHandlers(ResourceHandlerRegistry registry) {
        registry.addResourceHandler("swagger-ui.html").addResourceLocations("classpath:/META-INF/resources/");
        registry.addResourceHandler("/webjars/**").addResourceLocations("classpath:/META-INF/resources/webjars/");

        registry.addResourceHandler("/css/**").addResourceLocations("classpath:/static/css/");
        registry.addResourceHandler("/fonts/**").addResourceLocations("classpath:/static/fonts/");
        registry.addResourceHandler("/img/**").addResourceLocations("classpath:/static/img/");
        registry.addResourceHandler("/js/**").addResourceLocations("classpath:/static/js/");
        registry.addResourceHandler("/favicon.ico").addResourceLocations("classpath:/static/favicon.ico");
    }

    @Override
    public void addInterceptors(InterceptorRegistry registry) {
        registry.addInterceptor(authenticationRequestContextInterceptor).addPathPatterns("/**");
        WebMvcConfigurer.super.addInterceptors(registry);
    }

    @Override
    protected void configure(HttpSecurity http) throws Exception {
        configureLocalAuthentication(http);
    }

    protected void configureLocalAuthentication(HttpSecurity http) throws Exception {
        http.authorizeRequests() //
                .antMatchers("/index.html").permitAll() //
                .antMatchers("/workflow/**").permitAll() //
                .antMatchers("/swagger-ui.html/**", "/swagger-resources/**").permitAll()//
                .antMatchers("/webjars/**").permitAll() //
                .antMatchers("/v2/api-docs").permitAll() //
                .antMatchers("/csrf").permitAll() //
                .antMatchers("/**").permitAll() //
                .anyRequest().authenticated() //
                .and()//
                .addFilter(jwtSsoBasedAuthenticationFilter())//
                .csrf()//
                .disable() //
                .exceptionHandling() //
                .authenticationEntryPoint(new Http401AuthenticationEntryPoint()); //
    }

    protected Filter jwtSsoBasedAuthenticationFilter() throws Exception {
        JwtClientConfig config = new JwtClientConfig();
        config.setSigningKey(serviceManagementProperties.getJwtSigningKey());
        JwtSsoBasedAuthenticationFilter f = new JwtSsoBasedAuthenticationFilter(authenticationManager(), config);
        return (Filter) f;
    }
}

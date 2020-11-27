
package com.webank.taskman.config;

import com.fasterxml.classmate.TypeResolver;
import com.github.xiaoymin.swaggerbootstrapui.annotations.EnableSwaggerBootstrapUI;
import com.github.xiaoymin.swaggerbootstrapui.model.OrderExtensions;
import com.github.xiaoymin.swaggerbootstrapui.service.SpringAddtionalModelService;
import com.google.common.collect.Lists;
import com.webank.taskman.constant.SecurityConsts;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Import;
import springfox.bean.validators.configuration.BeanValidatorPluginsConfiguration;
import springfox.documentation.builders.ApiInfoBuilder;
import springfox.documentation.builders.PathSelectors;
import springfox.documentation.builders.RequestHandlerSelectors;
import springfox.documentation.service.*;
import springfox.documentation.spi.DocumentationType;
import springfox.documentation.spi.service.contexts.SecurityContext;
import springfox.documentation.spring.web.plugins.Docket;
import springfox.documentation.swagger2.annotations.EnableSwagger2;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpSession;
import java.util.List;

@Configuration
@EnableSwagger2
@EnableSwaggerBootstrapUI
@Import(BeanValidatorPluginsConfiguration.class)
public class SwaggerConfiguration {

    private static final String TOKEN_NAME="BearerToken";

    @Autowired
    SpringAddtionalModelService springAddtionalModelService;


    private final TypeResolver typeResolver;

    @Autowired
    public SwaggerConfiguration(TypeResolver typeResolver) {
        this.typeResolver = typeResolver;
    }


    @Bean(value = "groupRestApi")
    public Docket groupRestApi() {
        return new Docket(DocumentationType.SWAGGER_2)
                .apiInfo(groupApiInfo())
                .groupName("V1.0.0")
                .select()
                .apis(RequestHandlerSelectors.basePackage("com.webank.taskman.controller.x100"))
                .paths(PathSelectors.any())
                .build()
                .ignoredParameterTypes(HttpServletRequest.class)
                .ignoredParameterTypes(HttpSession.class)
//                .additionalModels(typeResolver.resolve(DeveloperApiInfo.class))
                .extensions(Lists.newArrayList(new OrderExtensions(2)))
                .securityContexts(Lists.newArrayList(securityContext()))
                .securitySchemes(Lists.<SecurityScheme>newArrayList(apiKey()));
    }

    private ApiInfo groupApiInfo(){
        return new ApiInfoBuilder()
                .title("Wecube-plugins-Taskman RESTFul API ")
                .description("<div style='font-size:14px;color:red;'>Wecube-plugins-Taskman RESTFul API</div>")
                .version("V1.0")
                .build();
    }

    private ApiKey apiKey() {
      //  HttpHeaders.AUTHORIZATION
        return new ApiKey(TOKEN_NAME, SecurityConsts.REQUEST_AUTH_HEADER, SecurityConsts.REQUEST_ITEM_NAME);
    }

    private SecurityContext securityContext() {
        return SecurityContext.builder()
                .securityReferences(defaultAuth())
                .forPaths(PathSelectors.regex("/.*"))
                .build();
    }

    List<SecurityReference> defaultAuth() {
        AuthorizationScope authorizationScope = new AuthorizationScope("global", "accessEverything");
        AuthorizationScope[] authorizationScopes = new AuthorizationScope[1];
        authorizationScopes[0] = authorizationScope;
        return Lists.newArrayList(new SecurityReference(TOKEN_NAME, authorizationScopes));
    }


}


package com.webank.taskman.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import springfox.documentation.builders.ApiInfoBuilder;
import springfox.documentation.builders.PathSelectors;
import springfox.documentation.builders.RequestHandlerSelectors;
import springfox.documentation.spi.DocumentationType;
import springfox.documentation.spring.web.plugins.Docket;
import springfox.documentation.swagger2.annotations.EnableSwagger2WebMvc;

@Configuration
@EnableSwagger2WebMvc
public class Knife4jConfiguration {

    @Bean(name="apiV1")
    public Docket defaultApi() {
        Docket docket=new Docket(DocumentationType.SWAGGER_2)
                .apiInfo(new ApiInfoBuilder()
                        .description("# Wecube-plugins-Taskman RESTful APIs")
                        .termsOfServiceUrl("http://www.xx.com/")
                        .version("1.0")
                        .build())
                .groupName("1.0.0")
                .select()
                .apis(RequestHandlerSelectors.basePackage("com.webank.taskman.controller.x100"))
                .paths(PathSelectors.any())
                .build();
        return docket;
    }


}

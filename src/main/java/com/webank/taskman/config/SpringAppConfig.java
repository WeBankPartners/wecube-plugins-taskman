package com.webank.taskman.config;

import com.webank.taskman.commons.AppProperties;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;

@Configuration
@EnableConfigurationProperties({ AppProperties.class, AppProperties.HttpClientProperties.class, AppProperties.ServiceTaskmanProperties.class })
@ComponentScan({ "com.webank.taskman.service" })
public class SpringAppConfig {

}

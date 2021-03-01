package com.webank.taskman.config;

import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Configuration;

import com.webank.taskman.commons.AppProperties;

@Configuration
@EnableConfigurationProperties({ AppProperties.class, AppProperties.HttpClientProperties.class, AppProperties.ServiceTaskmanProperties.class })
public class SpringAppConfig {

}

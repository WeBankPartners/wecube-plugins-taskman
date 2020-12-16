package com.webank.taskman;

import com.webank.taskman.utils.JsonUtils;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.transaction.annotation.EnableTransactionManagement;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;


@SpringBootApplication
@EnableTransactionManagement(proxyTargetClass = true)
public class TaskManApplication {
    public static void main(String[] args) throws IOException {
        SpringApplication.run(TaskManApplication.class, args);
//        test();
    }

}

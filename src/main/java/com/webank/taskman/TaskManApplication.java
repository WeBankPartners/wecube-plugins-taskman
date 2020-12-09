package com.webank.taskman;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.transaction.annotation.EnableTransactionManagement;

import java.io.IOException;




@SpringBootApplication
@EnableTransactionManagement
public class TaskManApplication {
    public static void main(String[] args) throws IOException {
        SpringApplication.run(TaskManApplication.class, args);
    }


}

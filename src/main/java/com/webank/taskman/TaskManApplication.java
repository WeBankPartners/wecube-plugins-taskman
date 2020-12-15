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

    public static void test() throws IOException {
        String json = "[{\"id\": \"1\",\"requestTempId\": \"1\",\"rootEntity\": \"22\",\"name\": \"33\",\"description\": \"44\",\"formItems\": [{\"name\": \"roleName\",\"value\": \"ADMIN\"},{\"name\": \"reporter\",\"value\": \"test测试工程师\"},{\"name\": \"age\",\"value\": 123}]}]";
        List<Map<String,Object>> list = new ArrayList<>();
        list = JsonUtils.toObject(json,list.getClass());
        System.out.println("after:"+list);
        list.stream().forEach(map->{
            ((List<Map<String,String>>)map.get("formItems")).stream().forEach(item->{
                map.put(item.get("name"),item.get("value"));
                map.remove("formItems");
            });
        });
        System.out.println("befor:"+list);
    }

}

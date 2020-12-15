package com.webank.taskman.dto;

import com.webank.taskman.utils.JsonUtils;
import springfox.documentation.spring.web.json.Json;

public class ParamJson {

    public static String documentJSON(){
        Json json1 = new Json("{\"type\":\"select\",\"width\":\"\"\"}");
        return JsonUtils.toJsonString(json1);
    }

    public static void main(String[] args) {
        System.out.println(documentJSON());
    }
}

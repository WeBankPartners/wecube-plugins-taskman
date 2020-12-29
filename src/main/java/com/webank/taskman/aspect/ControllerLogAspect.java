package com.webank.taskman.aspect;

import com.webank.taskman.utils.DateUtils;
import com.webank.taskman.utils.GsonUtil;
import org.aspectj.lang.JoinPoint;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.*;
import org.aspectj.lang.reflect.MethodSignature;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;
import org.springframework.web.context.request.RequestContextHolder;
import org.springframework.web.context.request.ServletRequestAttributes;
import org.springframework.web.servlet.HandlerMapping;

import javax.servlet.http.HttpServletRequest;
import java.io.UnsupportedEncodingException;
import java.net.URLDecoder;
import java.util.Map;

@Aspect
@Component
@ConditionalOnProperty(value = "service.taskman.aspect.log.controller")
public class ControllerLogAspect {

    private static final Logger log = LoggerFactory.getLogger(ControllerLogAspect.class);

    ThreadLocal<Long> startTime = new ThreadLocal<>();

    static final String pCutStr = "execution(* com.webank.taskman.controller.**.*.*(..))";

    public ControllerLogAspect(){
        log.info("Create log AOP instance:ControllerLogAspect...");
    }
    @Pointcut(value = pCutStr)
    public void logPointcut() {
    }
    @Around("logPointcut()")
    public Object doAround(ProceedingJoinPoint joinPoint) throws Throwable {
        try {
            HttpServletRequest request = getHttpServletRequest();
            String tragetClassName = joinPoint.getSignature().getDeclaringTypeName();
            String method = request.getMethod();
            String methodName = joinPoint.getSignature().getName();
            String uri = request.getAttribute(HandlerMapping.BEST_MATCHING_PATTERN_ATTRIBUTE).toString();
            StringBuffer logs = new StringBuffer();
            String bestMatchingPattern = request.getAttribute(HandlerMapping.BEST_MATCHING_PATTERN_ATTRIBUTE).toString();
            String jsonKey = "\n\t\"%s\":";
            String josnValue = "\"%s\"";
            logs.append(String.format("==========================Receive Request: [%s] start==========================",bestMatchingPattern)).append("\n{");
            logs.append(String.format(jsonKey,"URI")).append(uri);
            logs.append(String.format(jsonKey,"RequestMethod")).append(String.format(josnValue,method)).append(",");
            logs.append(String.format(jsonKey,"className")).append(String.format(josnValue,tragetClassName)).append(",");
            logs.append(String.format(jsonKey,"inteface")).append(String.format(josnValue,methodName)).append(",");
            Map pathVariables = (Map) request.getAttribute(HandlerMapping.URI_TEMPLATE_VARIABLES_ATTRIBUTE);
            if(null != pathVariables){
                logs.append(String.format(jsonKey,"pathParam")).append(String.format(josnValue,GsonUtil.GsonString(pathVariables))).append(",");
            }
            String queryString = request.getQueryString();
            if(!StringUtils.isEmpty(request.getQueryString())){
                logs.append(String.format(jsonKey,"queryParam")).append(String.format(josnValue,queryString)).append(",");
            }
            Object[] args = joinPoint.getArgs();
            if(!methodName.contains("S3") &&"POST".equals(method) && null != args && args.length > 0 ) {
                logs.append(String.format(jsonKey,"body")).append(String.format(josnValue,GsonUtil.GsonString(args[args.length - 1]))).append(",");
            }
            logs.append("\n}");
            String logContent = URLDecoder.decode(logs.toString(), "UTF-8");
            log.info(logContent);
            Object result = joinPoint.proceed();
            if(null == result) {
                return null;
            }
            log.info("result is Class：{}",result.getClass().getName());
            return result;
        } catch (Throwable e) {
            log.error("error：{}"+e.getMessage());
            throw e;
        }
    }


    @Before(value =pCutStr)
    public void beforMehhod() {
        startTime.set(System.currentTimeMillis());
    }

    @AfterReturning(returning="result",value = pCutStr)
    public void afterMehhod(JoinPoint joinPoint, Object result) {
        long end = System.currentTimeMillis();
        long total = end - startTime.get();
        HttpServletRequest request = getHttpServletRequest();
        String methodName = joinPoint.getSignature().getName();
        String bestMatchingPattern = request.getAttribute(HandlerMapping.BEST_MATCHING_PATTERN_ATTRIBUTE).toString();
        log.info("The total execution time of method [{}] is：{}",methodName, DateUtils.formatLongToTimeStr(total));
        log.info("==========================Response Request: [{}] complete=========================",bestMatchingPattern);
    }

    private HttpServletRequest getHttpServletRequest() {
        ServletRequestAttributes servletRequestAttributes = (ServletRequestAttributes) RequestContextHolder.getRequestAttributes();
        return servletRequestAttributes.getRequest();
    }

    private void getparams(Object[] args,StringBuffer logs) throws UnsupportedEncodingException {
        HttpServletRequest request = getHttpServletRequest();
        String method = request.getMethod();

        switch (method){
            case "POST":
                logs.append("body:").append(GsonUtil.GsonString(args[args.length-1]));
                break;
            case "GET":
                logs.append("queryParam:").append(request.getQueryString());
                break;
            default:
                break;
        }
    }

}

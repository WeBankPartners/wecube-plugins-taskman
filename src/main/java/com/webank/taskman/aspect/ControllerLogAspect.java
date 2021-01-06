package com.webank.taskman.aspect;

import com.webank.taskman.utils.DateUtils;
import com.webank.taskman.utils.GsonUtil;
import org.aspectj.lang.JoinPoint;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.*;
import org.aspectj.lang.reflect.MethodSignature;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;
import org.springframework.web.context.request.RequestContextHolder;
import org.springframework.web.context.request.ServletRequestAttributes;
import org.springframework.web.servlet.HandlerMapping;

import javax.servlet.http.HttpServletRequest;
import java.net.URLDecoder;
import java.util.Map;

@Aspect
@Component
@ConditionalOnProperty(value = "service.taskman.aspect-log-controller")
public class ControllerLogAspect {

    private static final Logger log = LoggerFactory.getLogger(ControllerLogAspect.class);

    ThreadLocal<Long> startTime = new ThreadLocal<>();

    static final String pCutStr = "execution(* com.webank.taskman.controller.**.*.*(..))";

    @Value(value = "${service.taskman.aspect-log-controller-printResult:false}")
    private boolean printResult = false;
    public ControllerLogAspect(){
        log.info("Create log AOP instance:ControllerLogAspect...");
    }

    @Pointcut(value = pCutStr)
    public void logPointcut() {}

    @Around("logPointcut()")
    public <T> Object doAround(ProceedingJoinPoint joinPoint) throws Throwable {
        StringBuffer logs = new StringBuffer();
        HttpServletRequest request = getHttpServletRequest();
        String bestMatchingPattern = request.getAttribute(HandlerMapping.BEST_MATCHING_PATTERN_ATTRIBUTE).toString();
        try {

            String tragetClassName = joinPoint.getSignature().getDeclaringTypeName();
            String requestMethod = request.getMethod();
            String methodName = joinPoint.getSignature().getName();
            String uri = request.getAttribute(HandlerMapping.BEST_MATCHING_PATTERN_ATTRIBUTE).toString();

            String jsonKey = "\n\t\"%s\":";
            String josnValue = "\"%s\"";
            logs.append(String.format("========Receive Request: [%s] start========",bestMatchingPattern)).append("\n{");
            logs.append(String.format(jsonKey,"URI")).append(String.format(josnValue,uri)).append(",");
            logs.append(String.format(jsonKey,"RequestMethod")).append(String.format(josnValue,requestMethod)).append(",");
            logs.append(String.format(jsonKey,"className")).append(String.format(josnValue,tragetClassName)).append(",");
            logs.append(String.format(jsonKey,"inteface")).append(String.format(josnValue,methodName)).append(",");
            Map pathVariables = (Map) request.getAttribute(HandlerMapping.URI_TEMPLATE_VARIABLES_ATTRIBUTE);
            if(null != pathVariables){
                logs.append(String.format(jsonKey,"pathParam")).append(GsonUtil.GsonString(pathVariables)).append(",");
            }
            String queryString = request.getQueryString();
            if(!StringUtils.isEmpty(request.getQueryString())){
                logs.append(String.format(jsonKey,"queryParam")).append(String.format(josnValue,queryString)).append(",");
            }
            Object[] args = joinPoint.getArgs();
            if(!methodName.contains("S3") &&"POST".equals(requestMethod) && null != args && args.length > 0 ) {
                logs.append(String.format(jsonKey,"body")).append(GsonUtil.GsonString(args[args.length - 1])).append(",");
            }
            logs.append(String.format(jsonKey,"returnClass")).append(String.format(josnValue,
                    getResultClass(((MethodSignature)joinPoint.getSignature()).getMethod().getGenericReturnType().getTypeName())
            )).append(",");
            Object result = joinPoint.proceed();
            if(null == result) {
                return null;
            }
            if(printResult){
                logs.append(String.format(jsonKey,"respone")).append(result);
            }
            logs.append("\n}");
            String logContent = URLDecoder.decode(logs.toString(), "UTF-8");
            log.info(logContent);
            log.info("The total execution time of method [{}] is：{}",
                    joinPoint.getSignature().getName(),
                    DateUtils.formatLongToTimeStr(System.currentTimeMillis() - startTime.get()));
            log.info("========Response Request: [{}] complete========",bestMatchingPattern);
            return result;
        } catch (Throwable e) {
            logs.append("error:").append(e.getMessage());
            logs.append("\n}");
            String logContent = URLDecoder.decode(logs.toString(), "UTF-8");
            log.info(logContent);
            log.info("========Response Request: [{}] complete========",bestMatchingPattern);
            throw e;
        }
    }

    private  String getResultClass(String name){

        String[] typeNames = name.split("<");
        for(int i=0;i<typeNames.length;i++){
            String str = typeNames[i].substring(0,typeNames[i].lastIndexOf(".")+1);
            name = name.replace(str,"");
        }
        return name;
    }

    @Before(value =pCutStr)
    public void beforMehhod() {
        startTime.set(System.currentTimeMillis());
    }

    @AfterReturning(returning="result",value = pCutStr)
    public void afterMehhod(JoinPoint joinPoint, Object result) {

    }

    private HttpServletRequest getHttpServletRequest() {
        ServletRequestAttributes servletRequestAttributes = (ServletRequestAttributes) RequestContextHolder.getRequestAttributes();
        return servletRequestAttributes.getRequest();
    }

}

package com.webank.taskman.handler;

import com.webank.taskman.base.JsonResponse;
import com.webank.taskman.commons.TaskmanException;
import com.webank.taskman.commons.TaskmanRuntimeException;
import com.webank.taskman.constant.BizCodeEnum;
import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.MessageSource;
import org.springframework.dao.DuplicateKeyException;
import org.springframework.http.HttpStatus;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.sql.SQLIntegrityConstraintViolationException;
import java.util.HashMap;
import java.util.Locale;
import java.util.Map;

/**
 * @author gavin
 */
@RestControllerAdvice
public class GlobalExceptionHandler {
    private static final Logger log = LoggerFactory.getLogger(GlobalExceptionHandler.class);

    public static final String MSG_ERR_CODE_PREFIX = "servicemgmt.msg.errcode.";

    public static final Locale DEF_LOCALE = Locale.ENGLISH;
    public static final String SQL_Exception = "Exception: ";

    @Autowired
    private MessageSource messageSource;

    @ExceptionHandler(TaskmanRuntimeException.class)
    @ResponseBody
    public JsonResponse handleWecubeException(HttpServletRequest request, final Exception e,
                                              HttpServletResponse response) {
        String errMsg = String.format("Processing failed cause by %s:%s", e.getClass().getSimpleName(),
                e.getMessage() == null ? "" : e.getMessage());
        log.error(errMsg + "\n", e.getMessage());
        TaskmanRuntimeException wecubeError = (TaskmanRuntimeException) e;

        return JsonResponse.customError(determineI18nErrorMessage(request, wecubeError));
    }

    @ExceptionHandler(value = RuntimeException.class)
    public JsonResponse handleException(RuntimeException e) {
        String errMsg = String.format("Processing failed cause by %s:%s", e.getClass().getSimpleName(),
                e.getMessage() == null ? "" : e.getMessage());
        log.error("错误异常:{}", e);
        String err =  e.getMessage();
        if(e instanceof DuplicateKeyException){
            err = err.substring(err.indexOf(SQL_Exception)).replace(SQL_Exception,"");
            err = err.substring(0,err.indexOf("###")-3);
        }
        err = String.format("Processing failed cause by :%s",err);
        return JsonResponse.customError(BizCodeEnum.RUNTIME_EXCEPTION, err,err);
    }


    @ExceptionHandler(value = MethodArgumentNotValidException.class)
    @ResponseStatus(HttpStatus.BAD_REQUEST)
    public JsonResponse handleException(MethodArgumentNotValidException e) {
        log.error("数据校验出现问题: {}, 异常类型: {}", e.getMessage(), e.getClass());
        BindingResult bindingResult = e.getBindingResult();

        Map<String, String> errorMap = new HashMap<>();
        bindingResult.getFieldErrors().forEach((fieldError -> {
            errorMap.put(fieldError.getField(), fieldError.getDefaultMessage());
        }));

        return JsonResponse.customError(BizCodeEnum.VAILD_EXCEPTION, errorMap);
    }

    @ExceptionHandler(value = Exception.class)
    @ResponseBody
    public JsonResponse defaultErrorHandler(HttpServletRequest req, Exception e) throws Exception {
        log.error("GlobalExceptionHandler: RequestHost {} invokes url {} ERROR: {}", req.getRemoteHost(),
                req.getRequestURL(), e.getMessage());
        if (e instanceof TaskmanException) {
            TaskmanException exception = (TaskmanException) e;
            return JsonResponse.customError(exception.getStatusCodeEnum());

        }
        return JsonResponse.customError(e.getMessage());
    }

    private String determineI18nErrorMessage(HttpServletRequest request, TaskmanRuntimeException e) {
        Locale locale = request.getLocale();
        if (locale == null) {
            locale = DEF_LOCALE;
        }
        if (StringUtils.isNoneBlank(e.getErrorCode())) {
            String msgCode = MSG_ERR_CODE_PREFIX + e.getErrorCode();
            String msg = "";
            try {
                msg = messageSource.getMessage(msgCode, e.getArgs(), locale);
            }catch (Exception ex){
                msg = e.getMessageKey();
                log.error("not find key:{}",msgCode);
            }
            return msg;
        } else {
            return e.getMessage();
        }
    }

}

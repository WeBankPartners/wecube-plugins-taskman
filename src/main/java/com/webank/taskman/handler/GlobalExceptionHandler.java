package com.webank.taskman.handler;

import com.webank.taskman.commons.TaskmanException;
import com.webank.taskman.constant.BizCodeEnum;
import com.webank.taskman.dto.JsonResponse;
import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.MessageSource;
import org.springframework.http.HttpStatus;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.util.HashMap;
import java.util.Locale;
import java.util.Map;

/**
 * 
 * @author gavin
 *
 */
@RestControllerAdvice
public class GlobalExceptionHandler {
    private static final Logger log = LoggerFactory.getLogger(GlobalExceptionHandler.class);

    public static final String MSG_ERR_CODE_PREFIX = "servicemgmt.msg.errcode.";

    public static final Locale DEF_LOCALE = Locale.ENGLISH;

    @Autowired
    private MessageSource messageSource;

    @ExceptionHandler(TaskmanException.class)
    @ResponseBody
    public JsonResponse handleWecubeException(HttpServletRequest request, final Exception e,
                                              HttpServletResponse response) {
        String errMsg = String.format("Processing failed cause by %s:%s", e.getClass().getSimpleName(),
                e.getMessage() == null ? "" : e.getMessage());
        log.error(errMsg + "\n", e);
        TaskmanException wecubeError = (TaskmanException) e;

        return JsonResponse.error(determineI18nErrorMessage(request, wecubeError));
    }

    @ExceptionHandler(value = RuntimeException.class)
    public JsonResponse handleException(RuntimeException e) {
        log.error("错误异常{}", e);

        return JsonResponse.customError(BizCodeEnum.RUNTIME_EXCEPTION.getCode(),
                BizCodeEnum.RUNTIME_EXCEPTION.getMessage(),
                null);
    }

    @ExceptionHandler(value = MethodArgumentNotValidException.class)
    @ResponseStatus(HttpStatus.BAD_REQUEST)
    public JsonResponse handleException(MethodArgumentNotValidException e) {
        log.error("数据校验出现问题{}, 异常类型: {}", e.getMessage(), e.getClass());
        BindingResult bindingResult = e.getBindingResult();

        Map<String, String> errorMap = new HashMap<>();
        bindingResult.getFieldErrors().forEach((fieldError -> {
            errorMap.put(fieldError.getField(), fieldError.getDefaultMessage());
        }));

        return JsonResponse.customError(BizCodeEnum.VAILD_EXCEPTION.getCode(),
                BizCodeEnum.VAILD_EXCEPTION.getMessage(),errorMap);
    }

    @ExceptionHandler(value = Exception.class)
    @ResponseBody
    public JsonResponse defaultErrorHandler(HttpServletRequest req, Exception e) throws Exception {
        log.error("errors occurred:", e);
        log.error("GlobalExceptionHandler: RequestHost {} invokes url {} ERROR: {}", req.getRemoteHost(),
                req.getRequestURL(), e.getMessage());
        return JsonResponse.error(e.getMessage());
    }

    private String determineI18nErrorMessage(HttpServletRequest request, TaskmanException e) {
        Locale locale = request.getLocale();
        if (locale == null) {
            locale = DEF_LOCALE;
        }
        if (StringUtils.isNoneBlank(e.getErrorCode())) {
            String msgCode = MSG_ERR_CODE_PREFIX + e.getErrorCode();
            String msg = messageSource.getMessage(msgCode, e.getArgs(), locale);
            return msg;
        } else {
            return e.getMessage();
        }
    }

}

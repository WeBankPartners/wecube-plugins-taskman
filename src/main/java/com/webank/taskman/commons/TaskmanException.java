package com.webank.taskman.commons;

import com.webank.taskman.constant.StatusCodeEnum;

public class TaskmanException extends RuntimeException {

    /**
     * 
     */
    private static final long serialVersionUID = -1397731977791044934L;

    private String errorCode;
    private String messageKey;
    private boolean applyMessage = false;
    private Object[] args;

    public TaskmanException(String message) {
        super(message);
        this.applyMessage = true;
    }

    public TaskmanException(String message, Throwable cause) {
        super(message, cause);
        this.applyMessage = true;
    }

    public TaskmanException(StatusCodeEnum errorCode) {
        this.errorCode = errorCode.getCode();
        this.messageKey = errorCode.getMessage();
        this.applyMessage = true;
    }
    public TaskmanException(StatusCodeEnum errorCode, String message) {
        super(message);
        this.errorCode = errorCode.getCode();
        this.applyMessage = true;
    }

    public TaskmanException(String errorCode, String message, Object... objects) {
        this(errorCode, null, message, null, false, objects);
    }

    private TaskmanException(String errorCode, String messageKey, String message, Throwable cause,
                             boolean applyMessage, Object... objects) {
        super(message, cause);
        this.errorCode = errorCode;
        this.messageKey = messageKey;
        if (objects != null && (this.args == null)) {
            this.args = new Object[objects.length];
            int index = 0;
            for (Object object : objects) {
                this.args[index] = object;
                index++;
            }
        }
        this.applyMessage = applyMessage;
    }

    public TaskmanException withErrorCode(String errorCode) {
        this.errorCode = errorCode;
        return this;
    }

    public TaskmanException withErrorCode(String errorCode, Object... objects) {
        this.errorCode = errorCode;
        if (objects != null && (this.args == null)) {
            this.args = new Object[objects.length];
            int index = 0;
            for (Object object : objects) {
                this.args[index] = object;
                index++;
            }
        }
        return this;
    }

    public TaskmanException withErrorCodeAndArgs(String errorCode, Object[] objects) {
        this.errorCode = errorCode;
        if (objects != null && (this.args == null)) {
            this.args = new Object[objects.length];
            int index = 0;
            for (Object object : objects) {
                this.args[index] = object;
                index++;
            }
        }
        return this;
    }

    public TaskmanException withMessageKey(String msgKey) {
        this.messageKey = msgKey;
        return this;
    }

    public TaskmanException withMessageKey(String msgKey, Object... objects) {
        this.messageKey = msgKey;
        if (objects != null && (this.args == null)) {
            this.args = new Object[objects.length];
            int index = 0;
            for (Object object : objects) {
                this.args[index] = object;
                index++;
            }
        }

        return this;
    }

    public String getErrorCode() {
        return errorCode;
    }

    public String getMessageKey() {
        return messageKey;
    }

    public boolean isApplyMessage() {
        return applyMessage;
    }

    public Object[] getArgs() {
        return args;
    }
}

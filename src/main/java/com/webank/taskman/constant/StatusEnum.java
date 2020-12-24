package com.webank.taskman.constant;

public enum  StatusEnum {
    //  Default state, the general value is 0 ,Examples:del_flag =0
    DEFAULT,
    // Available identification or Record deleted ,Examples:del_flag =1
    ENABLE,

    DISABLE,

    UNRELEASED,

    RELEASED,

    ALREADY_RECEIVED,

    UNCLAIMED,
    ;
    @Override
    public String toString(){
       return this.name().toLowerCase();
    }

    public static void main(String[] args) {
        System.out.println(StatusEnum.ENABLE);
    }
}

package com.webank.taskman.constant;

public enum  GenernalStatus {
    //  Default state, the general value is 0 ,Examples:del_flag =0
//    DEFAULT,
    // Available identification or Record deleted ,Examples:del_flag =1
    ENABLE,

    DISABLE,

//    UNRELEASED,
//
//    RELEASED,

    ALREADY_RECEIVED,

    UNCLAIMED,

    InProgress, // createNewWorkflowInstance is success!
    Processed,  // task is Processed!

    SUSPENSION, // taskInfo is canl
    ;
    @Override
    public String toString(){
       return this.name().toLowerCase();
    }

}
package com.webank.taskman.dto.req;

import io.swagger.annotations.ApiModelProperty;
import org.apache.commons.lang3.StringUtils;

public class QueryRoleRelationBaseReq {


    @ApiModelProperty(value = "使用角色",required = false,position = 109)
    private String useRoleName;
    @ApiModelProperty(value = "管理角色",required = false,position = 110)
    private String manageRoleName;

    @ApiModelProperty(hidden = true)
    private String sourceTableFix;


    @ApiModelProperty(hidden = true)
    private Integer roleType;

    @ApiModelProperty(hidden = true)
    private String roleName;

    @ApiModelProperty(hidden = true)
    private String conditionSql;

    public QueryRoleRelationBaseReq() {
    }

    public QueryRoleRelationBaseReq(String useRoleName, String manageRoleName, String sourceTableFix) {
        this.useRoleName = useRoleName;
        this.manageRoleName = manageRoleName;
        this.sourceTableFix = sourceTableFix;
    }

    public String getUseRoleName() {
        return useRoleName;
    }

    public void setUseRoleName(String useRoleName) {
        this.useRoleName = useRoleName;
    }


    public String getManageRoleName() {
        return manageRoleName;
    }

    public void setManageRoleName(String manageRoleName) {
        this.manageRoleName = manageRoleName;
    }

    public String getSourceTableFix() {
        return sourceTableFix;
    }

    public void setSourceTableFix(String sourceTableFix) {
        this.sourceTableFix = sourceTableFix;
    }

    public Integer getRoleType() {
        boolean useRole = false;
        boolean manageRole = false;
        if(!StringUtils.isEmpty(this.useRoleName)){
            useRole = true;
            this.roleType = 0;
            this.roleName = getUseRoleName();
        }
        if(!StringUtils.isEmpty(this.getManageRoleName())){
            manageRole = true;
            this.roleType = 1;
            this.roleName = getManageRoleName();
        }
        roleType = (useRole && manageRole) ? 2 :roleType;
        return roleType;
    }

    public void setRoleType(Integer roleType) {
        this.roleType = roleType;
    }

    public String getRoleName() {
        return roleName;
    }

    public void setRoleName(String roleName) {
        this.roleName = roleName;
    }
    protected static final  String NOT_ALL = "(SELECT COUNT(1) FROM role_relation rr WHERE rt.id = %s.record_id AND rr.role_type =%s AND MATCH(rr.role_name) AGAINST('%s') ) > 0";
    protected static final String ALL = "(SELECT COUNT(1) FROM role_relation rr WHERE rt.id = %s.record_id AND (rr.role_type =0 AND MATCH(rr.role_name) AGAINST('%s')) OR (rr.role_type =1 AND MATCH(rr.role_name) AGAINST('%s')))) > 0";

    public String getConditionSql() {
        if(StringUtils.isEmpty(conditionSql));
        switch (getRoleType()){
            case 0:
            case 1:
                conditionSql = String.format(NOT_ALL,sourceTableFix,getRoleType(),getRoleName());
                break;
            default:
                conditionSql = String.format(ALL,sourceTableFix,getUseRoleName(),getManageRoleName());
                break;

        }
        return conditionSql;
    }

    public void setConditionSql(String conditionSql) {
        this.conditionSql = conditionSql;
    }


    public static void main(String[] args) {
        QueryRoleRelationBaseReq req =  new QueryRoleRelationBaseReq();
        req.setSourceTableFix("rtf");
        req.setUseRoleName("APP_ARC,PRD_OPS");
        req.setManageRoleName("");
        System.out.println(req.getConditionSql());
        req.setUseRoleName("");
        req.setManageRoleName("APP_ARC,PRD_OPS");
        System.out.println(req.getConditionSql());
        req.setUseRoleName("APP_ARC");
        System.out.println(req.getConditionSql());
    }
}

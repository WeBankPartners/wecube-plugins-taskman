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
        if(!StringUtils.isEmpty(useRoleName)){
            this.roleType = !StringUtils.isEmpty(this.manageRoleName) ?2:0;
            this.roleName = useRoleName;
        }
        this.useRoleName = useRoleName;
    }


    public String getManageRoleName() {
        return manageRoleName;
    }

    public void setManageRoleName(String manageRoleName) {
        if(!StringUtils.isEmpty(manageRoleName)){
            this.roleType = !StringUtils.isEmpty(this.useRoleName) ?2:1;
            this.roleName = manageRoleName;
        }

        this.manageRoleName = manageRoleName;
    }

    public String getSourceTableFix() {
        return sourceTableFix;
    }

    public void setSourceTableFix(String sourceTableFix) {
        this.sourceTableFix = sourceTableFix;
    }

    public Integer getRoleType() {
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
    protected static final  String NOT_ALL = "(SELECT COUNT(1) FROM role_relation rr WHERE  %s.id =rr.record_id AND rr.role_type =%s AND MATCH(rr.role_name) AGAINST('%s') ) > 0";
    protected static final String ALL = "(SELECT COUNT(1) FROM role_relation rr WHERE %s.id = rr.record_id AND (rr.role_type =0 AND MATCH(rr.role_name) AGAINST('%s')) OR (rr.role_type =1 AND MATCH(rr.role_name) AGAINST('%s')))) > 0";

    public String getConditionSql() {
        if(StringUtils.isEmpty(conditionSql));
        switch (this.roleType){
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
        QueryRequestTemplateReq req = new QueryRequestTemplateReq();
        req.setSourceTableFix("rt");
        req.setUseRoleName("SUOPER");
        System.out.println(req.getConditionSql());
    }
}

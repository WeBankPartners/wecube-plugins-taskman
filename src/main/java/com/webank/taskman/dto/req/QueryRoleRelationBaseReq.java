package com.webank.taskman.dto.req;

import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.constant.RoleTypeEnum;
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
            this.conditionSql = null;
            this.roleType = !StringUtils.isEmpty(this.manageRoleName) ?2:1;
            this.roleName = useRoleName;
        }
        this.useRoleName = useRoleName;
    }


    public String getManageRoleName() {
        return manageRoleName;
    }

    public void setManageRoleName(String manageRoleName) {
        if(!StringUtils.isEmpty(manageRoleName)){
            this.conditionSql = null;
            this.roleType = !StringUtils.isEmpty(this.useRoleName) ?2:0;
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
    protected static final  String NOT_ALL = "%s.id in(SELECT rr.record_id FROM role_relation rr WHERE rr.role_type =%s AND MATCH(rr.role_name) AGAINST('%s'))";
    protected static final String ALL = "%s.id in (SELECT COUNT(1) FROM role_relation rr WHERE (rr.role_type =0 AND MATCH(rr.role_name) AGAINST('%s')) OR (rr.role_type =1 AND MATCH(rr.role_name) AGAINST('%s')))";
    public static final String QUERY_BY_ROLE_SQL = "SELECT rr.record_id FROM role_relation rr WHERE rr.role_type =%s AND MATCH(rr.role_name) AGAINST('%s')";

    public String getConditionSql() {
        if(StringUtils.isEmpty(conditionSql) && null != getRoleType()) {
            switch (this.roleType) {
                case 0:
                case 1:
                    conditionSql = String.format(NOT_ALL, sourceTableFix,getRoleType(), getRoleName());
                    break;
                default:
                    conditionSql = String.format(ALL, sourceTableFix,getUseRoleName(), getManageRoleName());
                    break;

            }
        }
        return conditionSql;
    }
    public void setConditionSql(String conditionSql) {
        this.conditionSql = conditionSql;
    }

    @ApiModelProperty(hidden = true)
    public void setEqUseRole(String tableFix){
        this.setSourceTableFix(tableFix);
        StringBuffer useRole = new StringBuffer();
        String queryUseRole = this.getUseRoleName();
        useRole.append(org.springframework.util.StringUtils.isEmpty(queryUseRole) ? "" : queryUseRole + ",");
        useRole.append(AuthenticationContextHolder.getCurrentUserRolesToString());
        this.setUseRoleName(useRole.toString());
    }

    @ApiModelProperty(hidden = true)
    public static String getEqUseRole(){
        String  inSql = String.format(QueryRoleRelationBaseReq.QUERY_BY_ROLE_SQL, RoleTypeEnum.USE_ROLE.getType(),
                AuthenticationContextHolder.getCurrentUserRolesToString());
        return inSql;
    }
}

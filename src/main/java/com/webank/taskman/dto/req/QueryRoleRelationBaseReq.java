package com.webank.taskman.dto.req;

import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.constant.RoleTypeEnum;
import io.swagger.annotations.ApiModelProperty;
import org.springframework.util.StringUtils;

import java.util.Arrays;
import java.util.HashSet;
import java.util.Set;

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
        this.conditionSql = null;
        if(!StringUtils.isEmpty(useRoleName)){
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
    protected static final  String NOT_ALL = "%s.id in(SELECT rr.record_id FROM role_relation rr WHERE rr.role_type =%s AND rr.role_name IN('%s'))";
    protected static final String ALL = "%s.id in (SELECT rr.record_id FROM role_relation rr WHERE (rr.role_type =0 AND rr.role_name IN('%s')) OR (rr.role_type =1  AND rr.role_name IN('%s')))";
    public static final String QUERY_BY_ROLE_SQL = "SELECT rr.record_id FROM role_relation rr WHERE rr.role_type =%s AND rr.role_name IN('%s')";

    public String getConditionSql() {
        if( null != getRoleType()) {
            switch (this.roleType) {
                case 0:
                case 1:
                    conditionSql = String.format(NOT_ALL, sourceTableFix,getRoleType(), getRoleNameStr(getRoleName()));
                    break;
                default:
                    conditionSql = String.format(ALL, sourceTableFix, getRoleNameStr(getManageRoleName()),getRoleNameStr(getUseRoleName()));
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
        queryCurrentUserRoles();
    }

    @ApiModelProperty(hidden = true)
    public static String getEqUseRole(){
        String roles = String.join("','",AuthenticationContextHolder.getCurrentUserRoles());
        return StringUtils.isEmpty(roles) ? null :String.format(QueryRoleRelationBaseReq.QUERY_BY_ROLE_SQL, RoleTypeEnum.USE_ROLE.getType(),roles);
    }

    public void queryCurrentUserRoles(){
        Set<String> sets = new HashSet<>();
        if(!StringUtils.isEmpty(getUseRoleName()) ){
            sets.addAll(Arrays.asList(getUseRoleName().split(",")));
        }
        if(null != AuthenticationContextHolder.getCurrentUserRoles()){
            sets.addAll( AuthenticationContextHolder.getCurrentUserRoles());
        }

        setUseRoleName(String.join(",",sets));
    }
    public String getRoleNameStr(String roles){
        Set<String> sets = new HashSet<>();
        if(!StringUtils.isEmpty(roles)){
            sets.addAll(Arrays.asList(roles.split(",")));
           return String.join("','",sets);
        }
        return "";
    }
}

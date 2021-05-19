package com.webank.taskman.dto.req;

import com.webank.taskman.commons.AuthenticationContextHolder;
import com.webank.taskman.constant.RoleType;
import org.springframework.util.StringUtils;

import java.util.Arrays;
import java.util.HashSet;
import java.util.Set;

public class RoleRelationBaseQueryReqDto {
    private String useRoleName;
    private String manageRoleName;
    private Integer roleType;
    private String roleName;
    private String conditionSql;

    public RoleRelationBaseQueryReqDto() {
    }

    public RoleRelationBaseQueryReqDto(String useRoleName, String manageRoleName) {
        this.useRoleName = useRoleName;
        this.manageRoleName = manageRoleName;
    }

    public String getUseRoleName() {
        return useRoleName;
    }

    public void setUseRoleName(String useRoleName) {
        this.conditionSql = null;
        if (!StringUtils.isEmpty(useRoleName)) {
            this.roleType = !StringUtils.isEmpty(this.manageRoleName) ? 2 : 1;
            this.roleName = useRoleName;
        }
        this.useRoleName = useRoleName;
    }

    public String getManageRoleName() {
        return manageRoleName;
    }

    public void setManageRoleName(String manageRoleName) {
        if (!StringUtils.isEmpty(manageRoleName)) {
            this.conditionSql = null;
            this.roleType = !StringUtils.isEmpty(this.useRoleName) ? 2 : 0;
            this.roleName = manageRoleName;
        }

        this.manageRoleName = manageRoleName;
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

    protected static final String NOT_ALL = "SELECT rr.record_id FROM role_relation rr WHERE rr.role_type =%s AND rr.role_name IN('%s')";
    protected static final String ALL = "SELECT rr.record_id FROM role_relation rr WHERE (rr.role_type =0 AND rr.role_name IN('%s')) OR (rr.role_type =1  AND rr.role_name IN('%s'))";
    public static final String QUERY_BY_ROLE_SQL = "SELECT rr.record_id FROM role_relation rr WHERE rr.role_type =%s AND rr.role_name IN('%s')";

    public String getConditionSql() {
        if (null != getRoleType()) {
            switch (this.roleType) {
            case 0:
            case 1:
                conditionSql = String.format(NOT_ALL, getRoleType(), getRoleNameStr(getRoleName()));
                break;
            default:
                conditionSql = String.format(ALL, getRoleNameStr(getManageRoleName()),
                        getRoleNameStr(getUseRoleName()));
                break;

            }
        }
        return conditionSql;
    }

    public void setConditionSql(String conditionSql) {
        this.conditionSql = conditionSql;
    }

    public void setEqUseRole(String tableFix) {
        queryCurrentUserRoles();
    }

    public static String getEqUseRole() {
        String roles = String.join("','", AuthenticationContextHolder.getCurrentUserRoles());
        return StringUtils.isEmpty(roles) ? null
                : String.format(RoleRelationBaseQueryReqDto.QUERY_BY_ROLE_SQL, RoleType.USE_ROLE.getType(), roles);
    }

    public void queryCurrentUserRoles() {
        Set<String> sets = new HashSet<>();
        if (!StringUtils.isEmpty(getUseRoleName())) {
            sets.addAll(Arrays.asList(getUseRoleName().split(",")));
        }
        if (null != AuthenticationContextHolder.getCurrentUserRoles()) {
            sets.addAll(AuthenticationContextHolder.getCurrentUserRoles());
        }

        setUseRoleName(String.join(",", sets));
    }

    public String getRoleNameStr(String roles) {
        Set<String> sets = new HashSet<>();
        if (!StringUtils.isEmpty(roles)) {
            sets.addAll(Arrays.asList(roles.split(",")));
            return String.join("','", sets);
        }
        return "";
    }
}

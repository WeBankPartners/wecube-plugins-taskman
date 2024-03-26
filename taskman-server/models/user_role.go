package models

type SimpleLocalRoleDto struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	DisplayName   string `json:"displayName"`
	Email         string `json:"email"`
	Status        string `json:"status"`        // Deleted, NotDeleted
	Administrator string `json:"administrator"` // 角色管理员
}

type UserDto struct {
	ID                string `json:"id"`
	UserName          string `json:"username"`
	Password          string `json:"password"`
	AuthType          string `json:"authType"` // LOCAL,UM
	RoleAdministrator string `json:"roleAdministrator"`
	Email             string `json:"email"`
}

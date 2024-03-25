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
	ID       string `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
	AuthType string `json:"authType"` // LOCAL,UM
}

type SimpleLocalUserDto struct {
	ID                string                `json:"id"`
	Username          string                `json:"username"`
	Password          string                `json:"password"`
	NativeName        string                `json:"nativeName"`
	Title             string                `json:"title"`
	EmailAddr         string                `json:"emailAddr"`
	OfficeTelNo       string                `json:"officeTelNo"`
	CellPhoneNo       string                `json:"cellPhoneNo"`
	Department        string                `json:"department"`
	EnglishName       string                `json:"englishName"`
	Active            bool                  `json:"active"`
	Blocked           bool                  `json:"blocked"`
	Deleted           bool                  `json:"deleted"`
	AuthSource        string                `json:"authSource"`
	AuthContext       string                `json:"authContext"`
	Roles             []*SimpleLocalRoleDto `json:"roles"`
	RoleAdministrator string                `json:"roleAdministrator"`
}

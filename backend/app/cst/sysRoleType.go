package cst

const (
	SYS_ROLE_KEY_ADMISSIONS = "admissions" // 角色-招办

	SYS_ROLE_KEY_DEPARTMENT = "department" // 角色-学院
)

var RoleMap = map[string]int64{

	SYS_ROLE_KEY_ADMISSIONS: 2,

	SYS_ROLE_KEY_DEPARTMENT: 3,
}

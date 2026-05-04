package authz

import "blog-server/entity"

func FromEntityRole(role entity.UserRole) Role {
	switch role {
	case entity.UserRoleAdmin:
		return RoleAdmin
	default:
		return RoleReader
	}
}

func FromString(role string) Role {
	switch role {
	case "admin":
		return RoleAdmin
	default:
		return RoleReader
	}
}

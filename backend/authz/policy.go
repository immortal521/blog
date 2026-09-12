package authz

import (
	"blog-server/entity"
)

// RolePermissions maps user roles to their permissions.
// Keys are entity.UserRole values: entity.UserRoleAdmin, entity.UserRoleReader
var rolePermissions = map[entity.UserRole][]Permission{
	entity.UserRoleAdmin: {
		{ResourcePost, ActionCreate},
		{ResourcePost, ActionRead},
		{ResourcePost, ActionUpdate},
		{ResourcePost, ActionDelete},

		{ResourceLink, ActionCreate},
		{ResourceLink, ActionUpdate},
		{ResourceLink, ActionDelete},
	},

	entity.UserRoleReader: {
		{ResourcePost, ActionRead},
	},
}

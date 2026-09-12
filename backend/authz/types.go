package authz

// Role type is now alias to entity.UserRole
// Use entity.UserRole.Values() for ["reader", "admin"]
// RoleAdmin = entity.UserRoleAdmin, RoleReader = entity.UserRoleReader

type Resource string

const (
	ResourcePost Resource = "post"
	ResourceLink Resource = "link"
)

type Action string

const (
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

type Permission struct {
	Resource Resource
	Action   Action
}

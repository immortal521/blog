package authz

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleReader Role = "reader"
)

func (Role) Values() []string {
	return []string{
		string(RoleReader),
		string(RoleAdmin),
	}
}

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

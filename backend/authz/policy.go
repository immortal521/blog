package authz

var rolePermissions = map[Role][]Permission{
	RoleAdmin: {
		{ResourcePost, ActionCreate},
		{ResourcePost, ActionRead},
		{ResourcePost, ActionUpdate},
		{ResourcePost, ActionDelete},

		{ResourceLink, ActionCreate},
		{ResourceLink, ActionUpdate},
		{ResourceLink, ActionDelete},
	},

	RoleReader: {
		{ResourcePost, ActionRead},
	},
}

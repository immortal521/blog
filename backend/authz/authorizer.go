package authz

import (
	"context"
	"errors"

	"blog-server/pkg/errx"
	"blog-server/repository"
)

type OwnerChecker interface {
	IsOwner(ctx context.Context, userID uint, resourceID uint) (bool, error)
}

type Authorizer struct {
	ownerCheckers map[Resource]OwnerChecker
}

func NewAuthorizer(
	postRepo repository.PostRepo,
	linkRepo repository.LinkRepo,
) *Authorizer {
	return &Authorizer{
		ownerCheckers: map[Resource]OwnerChecker{
			ResourcePost: postRepo,
			ResourceLink: linkRepo,
		},
	}
}

func can(role Role, resource Resource, action Action) bool {
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}

	for _, p := range perms {
		if p.Resource == resource && p.Action == action {
			return true
		}
	}
	return false
}

func ErrForbidden() error {
	return errx.New(errx.CodeForbidden, errors.New("permission denied"))
}

func (a *Authorizer) Authorize(
	ctx context.Context,
	userID uint,
	role Role,
	resource Resource,
	action Action,
	resourceID *uint,
) error {
	// RBAC
	if !can(role, resource, action) {
		return ErrForbidden()
	}

	if action == ActionUpdate || action == ActionDelete {

		if role == RoleAdmin {
			return nil
		}

		if resourceID == nil {
			return ErrForbidden()
		}

		checker, ok := a.ownerCheckers[resource]
		if !ok {
			return ErrForbidden()
		}

		okOwner, err := checker.IsOwner(ctx, userID, *resourceID)
		if err != nil {
			return errx.New(errx.CodeInternalError, err)
		}

		if !okOwner {
			return ErrForbidden()
		}
	}

	return nil
}

func (a *Authorizer) CanCreatePost(ctx context.Context, userID uint, role Role) error {
	return a.Authorize(ctx, userID, role, ResourcePost, ActionCreate, nil)
}

func (a *Authorizer) CanUpdatePost(ctx context.Context, userID uint, role Role, postID uint) error {
	return a.Authorize(ctx, userID, role, ResourcePost, ActionUpdate, &postID)
}

func (a *Authorizer) CanDeletePost(ctx context.Context, userID uint, role Role, postID uint) error {
	return a.Authorize(ctx, userID, role, ResourcePost, ActionDelete, &postID)
}

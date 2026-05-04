package repository

import (
	"context"

	"blog-server/datastore"
	"blog-server/ent"
	"blog-server/ent/user"
	"blog-server/entity"
	"blog-server/mapper"
	"blog-server/pkg/errx"

	"github.com/google/uuid"
)

// UserRepo handles persistence and read operations for user domain data.
type UserRepo interface {
	Create(ctx context.Context, user *entity.User, hashPassword string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUUID(ctx context.Context, uuidStr string) (bool, error)
	ExistsByID(ctx context.Context, id uint) (bool, error)
	GetAuthByEmail(ctx context.Context, email string) (*entity.UserAuth, error)
	GetAuthByID(ctx context.Context, id uint) (*entity.UserAuth, error)
}

type userRepo struct {
	store *datastore.DataStore
}

// NewUserRepo creates a new UserRepo instance.
func NewUserRepo(ds *datastore.DataStore) UserRepo {
	return &userRepo{store: ds}
}

// baseQuery returns a pre-configured UserQuery with soft-delete filter applied.
func (r *userRepo) baseQuery(ctx context.Context) *ent.UserQuery {
	return r.store.Client(ctx).User.
		Query().
		Where(user.DeletedAtIsNil())
}

func (r *userRepo) exists(ctx context.Context, predicate func(*ent.UserQuery)) (bool, error) {
	q := r.baseQuery(ctx)
	predicate(q)
	exists, err := q.Exist(ctx)
	if err != nil {
		return false, errx.New(errx.CodeInternalError, err)
	}
	return exists, nil
}

func (r *userRepo) getOne(
	ctx context.Context,
	predicate func(*ent.UserQuery),
) (*ent.User, error) {
	q := r.baseQuery(ctx)
	predicate(q)

	u, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errx.New(errx.CodeNotFound, err)
		}
		return nil, errx.New(errx.CodeInternalError, err)
	}
	return u, nil
}

// Create persists a new user with the provided hashed password.
//
// The input entity is not mutated. Avatar is optional and only set if non-nil.
// Returns the stored entity representation or an error if the operation fails.
func (r *userRepo) Create(ctx context.Context, u *entity.User, hashPassword string) (*entity.User, error) {
	c := r.store.Client(ctx).User.Create().
		SetUUID(u.UUID).
		SetEmail(u.Email).
		SetPassword(hashPassword).
		SetRole(u.Role).
		SetUsername(u.Username)

	if u.Avatar != nil {
		c.SetAvatar(*u.Avatar)
	}

	created, err := c.Save(ctx)
	if err != nil {
		return nil, errx.New(errx.CodeInternalError, err)
	}

	return mapper.ToUser(created), nil
}

// GetByEmail retrieves a user by email.
//
// Returns (nil, nil) if no active (non-soft-deleted) user exists with the given email.
// Errors unrelated to absence are returned as-is.
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	u, err := r.getOne(ctx, func(q *ent.UserQuery) {
		q.Where(user.EmailEQ(email))
	})
	if err != nil {
		return nil, err
	}
	return mapper.ToUser(u), nil
}

// ExistsByID reports whether a non-soft-deleted user exists with the given ID.
//
// Returns false with a non-nil error if the query fails.
func (r *userRepo) ExistsByID(ctx context.Context, id uint) (bool, error) {
	return r.exists(ctx, func(q *ent.UserQuery) {
		q.Where(user.IDEQ(id))
	})
}

// ExistsByUUID reports whether a non-soft-deleted user exists with the given UUID.
//
// Returns false with a non-nil error if the query fails.
func (r *userRepo) ExistsByUUID(ctx context.Context, uuidStr string) (bool, error) {
	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		return false, errx.New(errx.CodeInternalError, err)
	}

	return r.exists(ctx, func(q *ent.UserQuery) {
		q.Where(user.UUIDEQ(uuid))
	})
}

// ExistsByEmail reports whether a non-soft-deleted user exists with the given email.
//
// Returns false with a non-nil error if the query fails.
func (r *userRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.exists(ctx, func(q *ent.UserQuery) {
		q.Where(user.EmailEQ(email))
	})
}

// GetAuthByEmail retrieves authentication fields for a user identified by email.
//
// The returned struct includes only fields required for authentication workflows.
// Returns (nil, nil) if no active (non-soft-deleted) user exists with the given email.
// Errors unrelated to absence are returned as-is.
func (r *userRepo) GetAuthByEmail(ctx context.Context, email string) (*entity.UserAuth, error) {
	u, err := r.getOne(ctx, func(q *ent.UserQuery) {
		q.Where(
			user.EmailEQ(email),
		).Select(
			user.FieldID,
			user.FieldPassword,
			user.FieldRole,
		)
	})
	if err != nil {
		return nil, err
	}
	return &entity.UserAuth{
		ID:       u.ID,
		Password: u.Password,
		Role:     u.Role,
	}, nil
}

func (r *userRepo) GetAuthByID(ctx context.Context, id uint) (*entity.UserAuth, error) {
	u, err := r.getOne(ctx, func(q *ent.UserQuery) {
		q.Where(
			user.IDEQ(id),
		).Select(
			user.FieldID,
			user.FieldPassword,
			user.FieldRole,
		)
	})
	if err != nil {
		return nil, err
	}
	return &entity.UserAuth{
		ID:       u.ID,
		Password: u.Password,
		Role:     u.Role,
	}, nil
}

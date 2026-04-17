package repository

import (
	"context"

	"blog-server/datastore"
	"blog-server/ent"
	"blog-server/ent/user"
	"blog-server/entity"
	"blog-server/mapper"

	"github.com/google/uuid"
)

// UserRepo defines persistence operations for user entities and authentication data.
//
// Implementations treat soft-deleted records as non-existent. Methods that return
// an entity pointer return (nil, nil) when no matching record is found.
type UserRepo interface {
	// Create persists a new user with a pre-hashed password.
	//
	// The caller is responsible for providing a valid UUID and a securely hashed password.
	// Optional fields may be omitted via zero values or nil pointers.
	// Returns the stored user representation or an error if persistence fails.
	Create(ctx context.Context, user *entity.User, hashPassword string) (*entity.User, error)

	// GetByEmail retrieves a user by email.
	//
	// Returns (nil, nil) if no active (non-soft-deleted) user exists with the given email.
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// ExistsByEmail reports whether a non-soft-deleted user exists with the given email.
	//
	// Returns false with a non-nil error if the existence check fails.
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByUUID reports whether a non-soft-deleted user exists with the given UUID.
	//
	// Returns false with a non-nil error if the existence check fails.
	ExistsByUUID(ctx context.Context, uuid uuid.UUID) (bool, error)

	// ExistsByID reports whether a non-soft-deleted user exists with the given ID.
	//
	// Returns false with a non-nil error if the existence check fails.
	ExistsByID(ctx context.Context, id uint) (bool, error)

	// GetAuthByEmail retrieves authentication data by email.
	//
	// Only a subset of fields required for authentication is returned.
	// Returns (nil, nil) if no active (non-soft-deleted) user exists with the given email.
	GetAuthByEmail(ctx context.Context, email string) (*entity.UserAuth, error)
}

type userRepo struct {
	ds *datastore.DataStore
}

// ExistsByID reports whether a non-soft-deleted user exists with the given ID.
//
// Returns false with a non-nil error if the query fails.
func (r *userRepo) ExistsByID(ctx context.Context, id uint) (bool, error) {
	exists, err := r.ds.Client(ctx).User.
		Query().
		Where(
			user.IDEQ(id),
			user.DeletedAtIsNil(),
		).
		Exist(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// ExistsByUUID reports whether a non-soft-deleted user exists with the given UUID.
//
// Returns false with a non-nil error if the query fails.
func (r *userRepo) ExistsByUUID(ctx context.Context, uuid uuid.UUID) (bool, error) {
	exists, err := r.ds.Client(ctx).User.
		Query().
		Where(
			user.UUIDEQ(uuid),
			user.DeletedAtIsNil(),
		).
		Exist(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// GetAuthByEmail retrieves authentication fields for a user identified by email.
//
// The returned struct includes only fields required for authentication workflows.
// Returns (nil, nil) if no active (non-soft-deleted) user exists with the given email.
// Errors unrelated to absence are returned as-is.
func (r *userRepo) GetAuthByEmail(ctx context.Context, email string) (*entity.UserAuth, error) {
	u, err := r.ds.Client(ctx).User.
		Query().
		Where(
			user.EmailEQ(email),
			user.DeletedAtIsNil(),
		).
		Select(
			user.FieldUUID,
			user.FieldPassword,
			user.FieldRole,
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &entity.UserAuth{
		UUID:     u.UUID,
		Password: u.Password,
		Role:     u.Role,
	}, nil
}

// NewUserRepo constructs a UserRepo backed by the provided datastore.
//
// The returned implementation is safe for concurrent use if the underlying
// datastore client is concurrency-safe.
func NewUserRepo(ds *datastore.DataStore) UserRepo {
	return &userRepo{ds: ds}
}

// Create persists a new user with the provided hashed password.
//
// The input entity is not mutated. Avatar is optional and only set if non-nil.
// Returns the stored entity representation or an error if the operation fails.
func (r *userRepo) Create(ctx context.Context, u *entity.User, hashPassword string) (*entity.User, error) {
	c := r.ds.Client(ctx).User.
		Create().
		SetUUID(u.UUID).
		SetEmail(u.Email).
		SetPassword(hashPassword).
		SetRole(u.Role).
		SetUsername(u.Username)

	if u.Avatar != nil {
		c.SetAvatar(*u.Avatar)
	}

	user, err := c.Save(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.ToUser(user), nil
}

// GetByEmail retrieves a user by email.
//
// Returns (nil, nil) if no active (non-soft-deleted) user exists with the given email.
// Errors unrelated to absence are returned as-is.
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := r.ds.Client(ctx).User.
		Query().
		Where(
			user.EmailEQ(email),
			user.DeletedAtIsNil(),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ToUser(user), nil
}

// ExistsByEmail reports whether a non-soft-deleted user exists with the given email.
//
// Returns false with a non-nil error if the query fails.
func (r *userRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := r.ds.Client(ctx).User.
		Query().
		Where(
			user.EmailEQ(email),
			user.DeletedAtIsNil(),
		).
		Exist(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

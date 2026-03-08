package repository

import (
	"context"
	"errors"

	"blog-server/database"
	"blog-server/entity"
	"blog-server/errs"

	"gorm.io/gorm"
)

// IUserRepo defines the interface for user data access operations
type IUserRepo interface {
	GetUserByEmail(ctx context.Context, db database.DB, email string) (*entity.User, error)
	CreateUser(ctx context.Context, db database.DB, user *entity.User) (*entity.User, error)
	ExistsByEmail(ctx context.Context, db database.DB, email string) (bool, error)
	ExistsByID(ctx context.Context, db database.DB, id uint) (bool, error)
	ExistsByUUID(ctx context.Context, db database.DB, uuid string) (bool, error)
	GetRoleByUUID(ctx context.Context, db database.DB, uuid string) (*entity.UserRole, error)
}

type userRepo struct{}

// NewUserRepo creates a new user repository instance
func NewUserRepo() IUserRepo {
	return &userRepo{}
}

// CreateUser creates a new user in the database
func (u *userRepo) CreateUser(ctx context.Context, db database.DB, user *entity.User) (*entity.User, error) {
	gdb := unwrapDB(db)
	result := gorm.WithResult()
	err := gorm.G[entity.User](gdb, result).Create(ctx, user)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, errs.New(errs.CodeUserAlreadyExists, "user already exists", err)
	}
	if err != nil {
		return nil, errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email address
func (u *userRepo) GetUserByEmail(ctx context.Context, db database.DB, email string) (*entity.User, error) {
	gdb := unwrapDB(db)
	user, err := gorm.G[*entity.User](gdb).
		Where("email = ?", email).
		Where("deleted_at IS NULL").
		Take(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.New(errs.CodeUserNotFound, "user not found", err)
	}
	if err != nil {
		return nil, errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return user, nil
}

// ExistsByEmail checks if a user exists by email
func (u *userRepo) ExistsByEmail(ctx context.Context, db database.DB, email string) (bool, error) {
	return existsBy[entity.User](ctx, db, "email", email)
}

// ExistsByID checks if a user exists by ID
func (u *userRepo) ExistsByID(ctx context.Context, db database.DB, id uint) (bool, error) {
	return existsBy[entity.User](ctx, db, "id", id)
}

// ExistsByUUID checks if a user exists by UUID
func (u *userRepo) ExistsByUUID(ctx context.Context, db database.DB, uuid string) (bool, error) {
	return existsBy[entity.User](ctx, db, "uuid", uuid)
}

// GetRoleByUUID retrieves the user role by UUID
func (u *userRepo) GetRoleByUUID(ctx context.Context, db database.DB, uuid string) (*entity.UserRole, error) {
	gdb := unwrapDB(db)
	role, err := gorm.G[*entity.UserRole](gdb).
		Where("uuid = ?", uuid).
		Where("deleted_at IS NULL").
		Take(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.New(errs.CodeUserNotFound, "user not found", err)
	}
	if err != nil {
		return nil, errs.New(errs.CodeDatabaseError, "database error", err)
	}
	return role, nil
}

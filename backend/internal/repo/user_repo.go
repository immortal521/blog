package repo

import (
	"blog-server/internal/entity"
	"blog-server/pkg/errs"
	"context"
	"errors"

	"gorm.io/gorm"
)

type IUserRepo interface {
	GetUserByEmail(ctx context.Context, db *gorm.DB, email string) (*entity.User, error)
	CreateUser(ctx context.Context, db *gorm.DB, user *entity.User) error
	ExistsByEmail(ctx context.Context, db *gorm.DB, email string) (bool, error)
	ExistsByID(ctx context.Context, db *gorm.DB, id uint) (bool, error)
}

type userRepo struct{}

func (u *userRepo) CreateUser(ctx context.Context, db *gorm.DB, user *entity.User) error {
	panic("unimplemented")
}

func (u *userRepo) GetUserByEmail(ctx context.Context, db *gorm.DB, email string) (*entity.User, error) {
	user, err := gorm.G[*entity.User](db).Where("email = ?", email).Take(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepo) ExistsByEmail(ctx context.Context, db *gorm.DB, email string) (bool, error) {
	_, err := gorm.G[*entity.User](db).Where("email = ?", email).Take(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *userRepo) ExistsByID(ctx context.Context, db *gorm.DB, id uint) (bool, error) {
	_, err := gorm.G[*entity.User](db).Where("id = ?", id).Take(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func NewUserRepo() IUserRepo {
	return &userRepo{}
}

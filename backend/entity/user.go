package entity

import (
	"time"

	"blog-server/utils"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleReader UserRole = "reader"
	UserRoleAdmin  UserRole = "admin"
)

func (UserRole) Values() []string {
	return []string{
		string(UserRoleReader),
		string(UserRoleAdmin),
	}
}

type User struct {
	ID   uint
	UUID uuid.UUID

	Username string
	Email    string
	Avatar   *string

	Role UserRole

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserAuth struct {
	ID       uint
	Password string
	Role     UserRole
}

func GenerateUsername() string {
	return "user_" + utils.RandomString(6, "abcdefghijklmnopqrstuvwxyz0123456789")
}
